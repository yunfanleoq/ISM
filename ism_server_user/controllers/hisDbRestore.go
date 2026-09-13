/*
 * 历史库还原：列出 data/hisdbbackup、taosdump -i 导入、删除备份目录、上传 zip。
 */
package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"ISMServer/utils/errmsg"
	"ISMServer/utils/hisdbbackup"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
)

type hisBackupItem struct {
	FileName   string `json:"FileName"`
	CreateTime string `json:"CreateTime"`
	FilePath   string `json:"FilePath"`
	FileSize   string `json:"FileSize"`
	Valid      bool   `json:"Valid"`
}

type tdengineConn struct {
	host string
	port string
	user string
	pass string
}

func loadHistoryTdengine() (dbType int, conn tdengineConn, err error) {
	historyConf, err := config.NewConfig("ini", "conf/historyData.conf")
	if err != nil {
		return 0, conn, fmt.Errorf("无法读取 conf/historyData.conf")
	}
	dbType, _ = historyConf.Int("historyrecorddbtype")
	if dbType == 0 {
		dbType, _ = historyConf.Int("HistoryRecordDbType")
	}
	conn.host, _ = historyConf.String("TDengine::TDengineHost")
	conn.port, _ = historyConf.String("TDengine::TDenginePort")
	conn.user, _ = historyConf.String("TDengine::UserName")
	conn.pass, _ = historyConf.String("TDengine::PassWord")
	if conn.host == "" {
		conn.host = "127.0.0.1"
	}
	if conn.port == "" {
		conn.port = "6030"
	}
	if conn.user == "" {
		conn.user = "root"
	}
	if conn.pass == "" {
		conn.pass = "taosdata"
	}
	return dbType, conn, nil
}

func dumpEndpoint(viaDocker bool, conn tdengineConn) (host, port string) {
	if viaDocker {
		return "127.0.0.1", "6030"
	}
	return conn.host, conn.port
}

func runTaosSQL(viaDocker bool, container, host, port, user, pass, sql string) ([]byte, error) {
	if viaDocker {
		return runCmdTimeout(30*time.Second, "docker", "exec", container, "taos",
			"-h", host, "-P", port, "-u", user, taosdumpPasswordArg(pass), "-s", sql)
	}
	taosBin := "taos"
	if p, lookErr := exec.LookPath("taos"); lookErr == nil {
		taosBin = p
	}
	return runCmdTimeout(30*time.Second, taosBin,
		"-h", host, "-P", port, "-u", user, taosdumpPasswordArg(pass), "-s", sql)
}

func (c *HisDbOptController) GetHisBackUpList() {
	result := map[string]interface{}{
		"code": 0,
		"list": []hisBackupItem{},
	}
	root, err := hisdbbackup.AbsRoot()
	if err != nil {
		result["code"] = -1
		result["msg"] = "无法访问备份目录: " + err.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		result["code"] = -1
		result["msg"] = err.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	var list []hisBackupItem
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		abs := filepath.Join(root, e.Name())
		info, infoErr := e.Info()
		createTime := ""
		if infoErr == nil {
			createTime = info.ModTime().Format("2006-01-02 15:04:05")
		}
		list = append(list, hisBackupItem{
			FileName:   e.Name(),
			CreateTime: createTime,
			FilePath:   filepath.Join(hisdbbackup.Root, e.Name()),
			FileSize:   formatFileSize(hisdbbackup.DirSize(abs)),
			Valid:      hisdbbackup.IsValidTaosdumpDir(abs),
		})
	}
	result["list"] = list
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *HisDbOptController) HisDbRestore() {
	result := map[string]interface{}{
		"code": -1,
		"msg":  "",
	}
	var getParams struct {
		DbFilePath string `json:"DbFilePath"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &getParams); err != nil {
		result["msg"] = "参数不是 JSON"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	absTarget, err := hisdbbackup.SafePath(getParams.DbFilePath)
	if err != nil {
		result["msg"] = "非法备份路径: " + err.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	st, err := os.Stat(absTarget)
	if err != nil || !st.IsDir() {
		result["msg"] = "备份目录不存在"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	if !hisdbbackup.IsValidTaosdumpDir(absTarget) {
		result["msg"] = "不是有效的 taosdump 备份（目录内看不到 dbs.sql）。请指到同时有 dbs.sql 和 taosdump- 文件夹的那一层。"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	dbType, conn, err := loadHistoryTdengine()
	if err != nil {
		result["msg"] = err.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	if dbType != 2 {
		result["msg"] = "当前历史库不是 TDengine，无法用此方式还原。业务库请走「数据库管理」页。"
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

	dumpHost, dumpPort := dumpEndpoint(viaDocker, conn)
	dropSQL := "DROP DATABASE IF EXISTS " + tdengineHistoryDB + ";"
	dropOut, dropErr := runTaosSQL(viaDocker, tdContainer, dumpHost, dumpPort, conn.user, conn.pass, dropSQL)
	if dropErr != nil {
		detail := strings.TrimSpace(string(dropOut))
		if detail == "" {
			detail = dropErr.Error()
		}
		result["msg"] = "删除旧历史库失败: " + detail
		logs.Error("HisDbRestore DROP failed: %s", detail)
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	importArgs := []string{
		"-h", dumpHost,
		"-P", dumpPort,
		"-u", conn.user,
		taosdumpPasswordArg(conn.pass),
		"-i", absTarget,
	}
	var out []byte
	var runErr error
	if viaDocker {
		containerDir := "/tmp/ism_hisrestore_" + time.Now().Format("20060102150405")
		mkdirOut, mkdirErr := runCmdTimeout(30*time.Second, "docker", "exec", tdContainer, "mkdir", "-p", containerDir)
		if mkdirErr != nil {
			detail := strings.TrimSpace(string(mkdirOut))
			if detail == "" {
				detail = mkdirErr.Error()
			}
			result["msg"] = fmt.Sprintf("docker 内创建还原目录失败: 容器 %s: %s", dockerInspectName(tdContainer), detail)
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		cpOut, cpErr := runCmdTimeout(15*time.Minute, "docker", "cp", absTarget+"/.", tdContainer+":"+containerDir)
		if cpErr != nil {
			detail := strings.TrimSpace(string(cpOut))
			if detail == "" {
				detail = cpErr.Error()
			}
			_ = exec.Command("docker", "exec", tdContainer, "rm", "-rf", containerDir).Run()
			result["msg"] = "docker cp 备份到容器失败: " + detail
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		importArgs[len(importArgs)-1] = containerDir
		dockerArgs := append([]string{"exec", tdContainer, taosdump}, importArgs...)
		logs.Info("HisDbRestore docker taosdump: docker exec %s %s %s", tdContainer, taosdump, taosdumpArgsForLog(importArgs))
		out, runErr = runCmdTimeout(15*time.Minute, "docker", dockerArgs...)
		_ = exec.Command("docker", "exec", tdContainer, "rm", "-rf", containerDir).Run()
	} else {
		logs.Info("HisDbRestore taosdump: %s %s", taosdump, taosdumpArgsForLog(importArgs))
		out, runErr = runCmdTimeout(15*time.Minute, taosdump, importArgs...)
	}
	msg := strings.TrimSpace(string(out))
	if runErr != nil {
		if msg == "" {
			msg = runErr.Error()
		}
		result["msg"] = "taosdump 还原失败: " + msg
		logs.Error("HisDbRestore taosdump failed: %v", msg)
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "还原了历史库 "+filepath.Base(absTarget), errmsg.JournalLevelInfo, c.Ctx.Input)
	result["code"] = 0
	result["msg"] = "还原成功"
	result["path"] = filepath.Join(hisdbbackup.Root, filepath.Base(absTarget))
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *HisDbOptController) HisDbDeleteBackup() {
	result := map[string]interface{}{"code": errmsg.ERROR, "msg": ""}
	var getParams struct {
		DbFilePath string `json:"DbFilePath"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &getParams); err != nil {
		result["code"] = errmsg.NOTJSON
		result["msg"] = "参数不是 JSON"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	absTarget, err := hisdbbackup.SafePath(getParams.DbFilePath)
	if err != nil {
		result["msg"] = "非法备份路径: " + err.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	st, err := os.Stat(absTarget)
	if err != nil || !st.IsDir() {
		result["msg"] = "备份目录不存在"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	if err := os.RemoveAll(absTarget); err != nil {
		result["msg"] = err.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除历史库备份 "+filepath.Base(absTarget), errmsg.JournalLevelInfo, c.Ctx.Input)
	result["code"] = errmsg.SUCCSECODE
	c.Data["json"] = result
	c.ServeJSON()
}
