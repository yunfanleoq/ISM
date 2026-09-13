package controllers

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"ISMServer/utils/hisdbbackup"
)

func (c *HisDbOptController) HisDbBackupUpload() {
	result := map[string]interface{}{"Code": -1, "msg": ""}
	_, h, err := c.GetFile("file")
	if err != nil || h == nil {
		result["msg"] = "未收到上传文件"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	ext := strings.ToLower(filepath.Ext(h.Filename))
	if ext != ".zip" {
		result["Code"] = -2
		result["msg"] = "只接受 zip（把含 dbs.sql 的备份目录打成 zip）"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	root, err := hisdbbackup.AbsRoot()
	if err != nil {
		result["msg"] = err.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	tmpZip := filepath.Join(root, ".upload_"+time.Now().Format("20060102150405")+".zip")
	if saveErr := c.SaveToFile("file", tmpZip); saveErr != nil {
		result["Code"] = -4
		result["msg"] = "保存上传文件失败: " + saveErr.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	defer os.Remove(tmpZip)

	destName := "ISM_TDengine_Upload_" + time.Now().Format("2006-01-02_15-04-05")
	dest := filepath.Join(root, destName)
	if unzipErr := hisdbbackup.UnzipSafe(tmpZip, dest); unzipErr != nil {
		_ = os.RemoveAll(dest)
		result["msg"] = "解压失败: " + unzipErr.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	dest = hisdbbackup.FlattenUploadedDump(dest)
	if !hisdbbackup.IsValidTaosdumpDir(dest) {
		_ = os.RemoveAll(dest)
		result["msg"] = "zip 内看不到有效 taosdump 备份（需要 dbs.sql）。请打包同时含 dbs.sql 和 taosdump- 文件夹的那一层。"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	result["Code"] = 0
	result["code"] = 0
	result["msg"] = "已加入还原列表，不会自动覆盖历史库"
	result["path"] = filepath.Join(hisdbbackup.Root, filepath.Base(dest))
	c.Data["json"] = result
	c.ServeJSON()
}
