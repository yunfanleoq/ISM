/*
 * 历史库备份：优先 TDengine taosdump；不可用时明确返回提示。
 */
package controllers

import (
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

func resolveTaosdump() (bin string, viaDocker bool, container string) {
	if p, err := exec.LookPath("taosdump"); err == nil {
		return p, false, ""
	}
	container = os.Getenv("TD_CONTAINER")
	if container == "" {
		container = "tdengine"
	}
	if exec.Command("docker", "exec", container, "taosdump", "-V").Run() == nil {
		return "taosdump", true, container
	}
	return "", false, ""
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
		result["msg"] = "当前历史库不是 TDengine，暂不支持此备份方式"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	taosdump, viaDocker, tdContainer := resolveTaosdump()
	if taosdump == "" {
		result["msg"] = "未找到 taosdump：请安装 TDengine 客户端，或确保 docker 容器 tdengine 内可用 taosdump（TD_CONTAINER 可改容器名）。"
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

	if viaDocker {
		containerDir := "/tmp/" + distName
		mkdirCmd := exec.Command("docker", "exec", tdContainer, "mkdir", "-p", containerDir)
		if mkdirErr := mkdirCmd.Run(); mkdirErr != nil {
			logs.Error("HisDbBackUp docker mkdir failed: %v", mkdirErr)
			result["msg"] = "docker 内创建备份目录失败: " + mkdirErr.Error()
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		dockerArgs := []string{"exec", "-i", tdContainer, taosdump,
			"-h", host, "-P", port, "-u", user, "-p", pass,
			"-o", containerDir, "-A",
		}
		cmd := exec.Command("docker", dockerArgs...)
		out, runErr := cmd.CombinedOutput()
		msg := strings.TrimSpace(string(out))
		if runErr != nil {
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
	} else {
		args := []string{"-h", host, "-P", port, "-u", user, "-p", pass, "-o", hostOut, "-A"}
		cmd := exec.Command(taosdump, args...)
		out, runErr := cmd.CombinedOutput()
		msg := strings.TrimSpace(string(out))
		if runErr != nil {
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
