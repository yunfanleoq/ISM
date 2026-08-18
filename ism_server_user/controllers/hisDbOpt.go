/**
 * 历史库备份：优先 TDengine taosdump；不可用时明确返回提示。
 */
package controllers

import (
	"ISMServer/utils/errmsg"
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
)

const HisSavePath string = "data/hisdbbackup/"

type HisDbOptController struct {
	beego.Controller
}

func ensureHisBackupDir() error {
	if _, err := os.Stat(HisSavePath); os.IsNotExist(err) {
		return os.MkdirAll(HisSavePath, os.ModePerm)
	}
	return nil
}

func zipDir(srcDir, zipPath string) error {
	zf, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zf.Close()
	zw := zip.NewWriter(zf)
	defer zw.Close()
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return err
		}
		rel, relErr := filepath.Rel(srcDir, path)
		if relErr != nil {
			return relErr
		}
		w, createErr := zw.Create(rel)
		if createErr != nil {
			return createErr
		}
		f, openErr := os.Open(path)
		if openErr != nil {
			return openErr
		}
		_, copyErr := io.Copy(w, f)
		_ = f.Close()
		return copyErr
	})
}

func resolveTaosdump() (bin string, viaDocker bool, container string) {
	if p, err := exec.LookPath("taosdump"); err == nil {
		return p, false, ""
	}
	container = strings.TrimSpace(os.Getenv("TD_CONTAINER"))
	if container == "" {
		container = "tdengine"
	}
	if _, err := exec.LookPath("docker"); err != nil {
		return "", false, ""
	}
	out, err := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", container).CombinedOutput()
	if err != nil || !strings.Contains(strings.TrimSpace(string(out)), "true") {
		return "", false, ""
	}
	if exec.Command("docker", "exec", container, "taosdump", "-V").Run() == nil {
		return "taosdump", true, container
	}
	return "", false, ""
}

func (c *HisDbOptController) HisDbBackUp() {
	result := map[string]interface{}{"code": errmsg.ERROR}
	_ = ensureHisBackupDir()

	HistoryConf, err := config.NewConfig("ini", "conf/historyData.conf")
	if err != nil {
		result["msg"] = "无法读取 conf/historyData.conf"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	dbType, _ := HistoryConf.Int("HistoryRecordDbType")
	if dbType != 2 {
		result["code"] = -2
		result["msg"] = fmt.Sprintf("当前历史库类型=%d，暂仅支持 TDengine(2) 备份；请在配置页确认库类型或使用外部工具备份", dbType)
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	taosdump, viaDocker, tdContainer := resolveTaosdump()
	if taosdump == "" {
		result["code"] = -3
		result["msg"] = "未找到 taosdump：请安装 TDengine 客户端，或确保 docker 容器 tdengine 内可用 taosdump（TD_CONTAINER 可改容器名）。"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	host, _ := HistoryConf.String("TDengineHost")
	if host == "" {
		host = "127.0.0.1"
	}
	// 容器内备份应连本机 native 端口
	port, _ := HistoryConf.String("TDenginePort")
	if viaDocker {
		port = "6030"
		host = "127.0.0.1"
	} else if port == "" {
		port = "6030"
	}
	user, _ := HistoryConf.String("TDengineUserName")
	if user == "" {
		user = "root"
	}
	pass, _ := HistoryConf.String("TDenginePassWord")
	if pass == "" {
		pass = "taosdata"
	}

	now := time.Now().Format("2006-01-02_15-04-05")
	distName := "ISM_TDengine_Backup_" + now
	outDir := filepath.Join(HisSavePath, distName)
	_ = os.MkdirAll(outDir, os.ModePerm)

	args := []string{
		"-h", host, "-P", port, "-u", user, "-p", pass,
		"-o", outDir, "-a",
	}
	var cmd *exec.Cmd
	if viaDocker {
		// 容器内 /tmp 再 docker cp；taosdump 要求 -o 目录已存在
		contOut := "/tmp/" + distName
		mkdirCmd := exec.Command("docker", "exec", "-i", tdContainer, "mkdir", "-p", contOut)
		if mkdirErr := mkdirCmd.Run(); mkdirErr != nil {
			_ = os.RemoveAll(outDir)
			result["msg"] = "docker 创建备份目录失败: " + mkdirErr.Error()
			logs.Error("HisDbBackUp docker mkdir failed: %v", mkdirErr)
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		dockerArgs := []string{"exec", "-i", tdContainer, taosdump,
			"-h", host, "-P", port, "-u", user, "-p", pass,
			"-o", contOut, "-a",
		}
		cmd = exec.Command("docker", dockerArgs...)
		var stderr strings.Builder
		cmd.Stderr = &stderr
		if runErr := cmd.Run(); runErr != nil {
			_ = os.RemoveAll(outDir)
			_, _ = exec.Command("docker", "exec", "-i", tdContainer, "rm", "-rf", contOut).CombinedOutput()
			msg := strings.TrimSpace(stderr.String())
			if msg == "" {
				msg = runErr.Error()
			}
			if strings.Contains(msg, "is not exist") {
				result["msg"] = "docker taosdump 失败（输出目录不存在）: " + msg
			} else {
				result["msg"] = "docker exec taosdump 失败: " + msg
			}
			logs.Error("HisDbBackUp docker taosdump failed: %v", msg)
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		_ = os.RemoveAll(outDir)
		cp := exec.Command("docker", "cp", tdContainer+":"+contOut, HisSavePath)
		if cpErr := cp.Run(); cpErr != nil {
			result["msg"] = "docker cp 备份目录失败: " + cpErr.Error()
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		_, _ = exec.Command("docker", "exec", "-i", tdContainer, "rm", "-rf", contOut).CombinedOutput()
	} else {
		cmd = exec.Command(taosdump, args...)
		var stderr strings.Builder
		cmd.Stderr = &stderr
		if runErr := cmd.Run(); runErr != nil {
			_ = os.RemoveAll(outDir)
			msg := strings.TrimSpace(stderr.String())
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

	zipPath := outDir + ".zip"
	if zipErr := zipDir(outDir, zipPath); zipErr != nil {
		result["msg"] = "打包失败: " + zipErr.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	_ = os.RemoveAll(outDir)

	info, _ := os.Stat(zipPath)
	var size int64
	if info != nil {
		size = info.Size()
	}
	result["code"] = errmsg.SUCCSECODE
	result["fileName"] = filepath.Base(zipPath)
	result["fileSize"] = size
	result["msg"] = "ok"
	if viaDocker {
		result["via"] = "docker:" + tdContainer
	}
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *HisDbOptController) GetHisBackUpList() {
	_ = ensureHisBackupDir()
	type FileInfo struct {
		FileName   string
		FilePath   string
		CreateTime string
		FileSize   string
	}
	var list []FileInfo
	_ = filepath.Walk(HisSavePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".zip") && !strings.HasSuffix(strings.ToLower(info.Name()), ".sql") {
			return nil
		}
		list = append(list, FileInfo{
			FileName:   info.Name(),
			FilePath:   path,
			CreateTime: info.ModTime().Format("2006-01-02 15:04:05"),
			FileSize:   formatFileSize(info.Size()),
		})
		return nil
	})
	c.Data["json"] = map[string]interface{}{"code": errmsg.SUCCSECODE, "list": list}
	c.ServeJSON()
}

func (c *HisDbOptController) HisDbDown() {
	type ReqStu struct {
		DbFilePath string `json:"DbFilePath"`
	}
	var getParams ReqStu
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &getParams); err != nil || getParams.DbFilePath == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"code": errmsg.NOTJSON}
		c.ServeJSON()
		return
	}
	absSave, _ := filepath.Abs(HisSavePath)
	absTarget, err := filepath.Abs(getParams.DbFilePath)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"code": errmsg.ERROR}
		c.ServeJSON()
		return
	}
	rel, err := filepath.Rel(absSave, absTarget)
	if err != nil || strings.HasPrefix(rel, "..") {
		c.Data["json"] = map[string]interface{}{"code": -9, "msg": "path outside backup dir"}
		c.ServeJSON()
		return
	}
	f, err := os.Open(absTarget)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"code": errmsg.ERROR}
		c.ServeJSON()
		return
	}
	defer f.Close()
	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+filepath.Base(absTarget))
	_, _ = io.Copy(c.Ctx.ResponseWriter, f)
}
