/*
 * 历史库备份：优先 TDengine taosdump；不可用时明确返回提示。
 */
package controllers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type HisDbOptController struct {
	beego.Controller
}

// tdengineHistoryDB 与 dealWithHistoryData 建库名一致，备份只 dump 业务历史库。
const tdengineHistoryDB = "ISMHistoryDb"

func taosdumpPasswordArg(pass string) string {
	return "-p" + pass
}

func buildTaosdumpArgs(host, port, user, pass, outDir string) []string {
	return []string{
		"-h", host,
		"-P", port,
		"-u", user,
		taosdumpPasswordArg(pass),
		"-D", tdengineHistoryDB,
		"-o", outDir,
	}
}

func taosdumpArgsForLog(args []string) string {
	parts := make([]string, len(args))
	copy(parts, args)
	for i, a := range parts {
		if strings.HasPrefix(a, "-p") && len(a) > 2 {
			parts[i] = "-p***"
		}
	}
	return strings.Join(parts, " ")
}

func runCmdTimeout(timeout time.Duration, name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return out, fmt.Errorf("超时 %s（这是 TDengine 历史库备份，不是业务库备份；业务库请走数据库管理页）", timeout)
	}
	return out, err
}

func dockerInspectName(container string) string {
	out, err := runCmdTimeout(8*time.Second, "docker", "inspect", "-f", "{{.State.Status}}", container)
	status := strings.TrimSpace(string(out))
	if err != nil {
		if status == "" {
			status = err.Error()
		}
		return fmt.Sprintf("%s (inspect失败: %s)", container, status)
	}
	return fmt.Sprintf("%s (status=%s)", container, status)
}

func resolveTaosdump() (bin string, viaDocker bool, container string) {
	if p, err := exec.LookPath("taosdump"); err == nil {
		return p, false, ""
	}
	container = os.Getenv("TD_CONTAINER")
	if container == "" {
		container = "tdengine"
	}
	out, err := runCmdTimeout(8*time.Second, "docker", "exec", container, "taosdump", "-V")
	if err == nil {
		return "taosdump", true, container
	}
	logs.Warn("HisDbBackUp docker taosdump -V failed container=%s: %v %s", container, err, strings.TrimSpace(string(out)))
	return "", false, container
}

func (c *HisDbOptController) HisDbBackUp() {
	result := map[string]interface{}{
		"code": -1,
		"msg":  "",
		"path": "",
	}
	historyConf, err := config.NewConfig("ini", "conf/historyData.conf")
	if err != nil {
		result["msg"] = "无法读取 conf/historyData.conf"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	dbType, _ := historyConf.Int("historyrecorddbtype")
	if dbType == 0 {
		dbType, _ = historyConf.Int("HistoryRecordDbType")
	}
	if dbType != 2 {
		result["msg"] = "当前历史库不是 TDengine，暂不支持此备份方式。业务库（OceanBase/MariaDB）请走「数据库管理」页备份，不要用本按钮。"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	taosdump, viaDocker, tdContainer := resolveTaosdump()
	if taosdump == "" {
		hint := "未找到本机 taosdump"
		if tdContainer != "" {
			hint = fmt.Sprintf("%s；docker 容器 %s", hint, dockerInspectName(tdContainer))
		}
		result["msg"] = hint + "。请安装 TDengine 客户端，或确保容器内可用 taosdump（TD_CONTAINER 可改容器名）。"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	host, _ := historyConf.String("TDengine::TDengineHost")
	port, _ := historyConf.String("TDengine::TDenginePort")
	user, _ := historyConf.String("TDengine::UserName")
	pass, _ := historyConf.String("TDengine::PassWord")
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "6030"
	}
	if user == "" {
		user = "root"
	}
	if pass == "" {
		pass = "taosdata"
	}

	stamp := time.Now().Format("2006-08-02_15-04-05")
	distName := "ISM_TDengine_Backup_" + stamp
	hostOut := filepath.Join("data", "hisdbbackup", distName)
	if mkErr := os.MkdirAll(hostOut, 0755); mkErr != nil {
		result["msg"] = "创建备份目录失败: " + mkErr.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	dumpHost := host
	dumpPort := port
	if viaDocker {
		// 容器内 taosdump 走 native 6030，REST 6041 不可用
		dumpHost = "127.0.0.1"
		dumpPort = "6030"
	}

	if viaDocker {
		containerDir := "/tmp/" + distName
		mkdirOut, mkdirErr := runCmdTimeout(30*time.Second, "docker", "exec", "-T", tdContainer, "mkdir", "-p", containerDir)
		if mkdirErr != nil {
			detail := strings.TrimSpace(string(mkdirOut))
			if detail == "" {
				detail = mkdirErr.Error()
			}
			logs.Error("HisDbBackUp docker mkdir failed container=%s: %v %s", tdContainer, mkdirErr, detail)
			result["msg"] = fmt.Sprintf("docker 内创建备份目录失败: 容器 %s: %s", dockerInspectName(tdContainer), detail)
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		dumpArgs := buildTaosdumpArgs(dumpHost, dumpPort, user, pass, containerDir)
		dockerArgs := append([]string{"exec", "-T", tdContainer, taosdump}, dumpArgs...)
		logs.Info("HisDbBackUp docker taosdump: docker exec -T %s %s %s", tdContainer, taosdump, taosdumpArgsForLog(dumpArgs))
		out, runErr := runCmdTimeout(15*time.Minute, "docker", dockerArgs...)
		msg := strings.TrimSpace(string(out))
		if runErr != nil {
			if msg == "" {
				msg = runErr.Error()
			}
			if strings.Contains(msg, "is not exist") || strings.Contains(msg, "not exist") {
				result["msg"] = "docker taosdump 失败（输出目录不存在）: " + msg
			} else {
				result["msg"] = "docker exec taosdump 失败: " + msg
			}
			logs.Error("HisDbBackUp docker taosdump failed: %v", msg)
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		_ = exec.Command("docker", "cp", fmt.Sprintf("%s:%s/.", tdContainer, containerDir), hostOut).Run()
		_ = exec.Command("docker", "exec", "-T", tdContainer, "rm", "-rf", containerDir).Run()
	} else {
		args := buildTaosdumpArgs(dumpHost, dumpPort, user, pass, hostOut)
		logs.Info("HisDbBackUp taosdump: %s %s", taosdump, taosdumpArgsForLog(args))
		out, runErr := runCmdTimeout(15*time.Minute, taosdump, args...)
		msg := strings.TrimSpace(string(out))
		if runErr != nil {
			if msg == "" {
				msg = runErr.Error()
			}
			result["msg"] = "taosdump 失败: " + msg
			logs.Error("HisDbBackUp taosdump failed: %v", msg)
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
	}

	result["code"] = 0
	result["msg"] = "备份成功"
	result["path"] = hostOut
	c.Data["json"] = result
	c.ServeJSON()
}
