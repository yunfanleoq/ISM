/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:56:24
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package controllers

import (
	"ISMServer/models"
	protocolCommonFunc "ISMServer/protocol/commFunc"
	protocolCommon "ISMServer/protocol/common"
	ISMScript "ISMServer/task/ISMScript"
	staticDataTask "ISMServer/task/staticData"
	"ISMServer/utils/errmsg"
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/config"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/xormplus/xorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbOptController struct {
	beego.Controller
}

const SavePath string = "data/dbbackup/"

// MysqlSQLDump ...
func MysqlSQLDump(host, port, dbname, user, password, char, backupfilePath, zippwd string, tables []string, ProjectID string) (string, error) {
	db, err := xorm.NewEngine("mysql", user+":"+password+"@("+host+":"+port+")/"+dbname+"?charset="+char)
	defer db.Close()
	if err != nil {
		return "", err
	}

	nowtime := time.Now().Format("2006-01-02_15-04-05")
	dist := "Mysql_Backup_" + nowtime + ".sql"
	db.ProjectUUid = ProjectID
	err = db.DumpAllToFile(backupfilePath+"/"+dist, tables)
	if err != nil {
		return "", err
	}
	return dist, err
}

func SqliteSQLDump(dbpath, backupfilePath, zippwd string, tables []string, ProjectID string) (string, error) {
	db, err := xorm.NewEngine("sqlite3", dbpath)
	defer db.Close()
	if err != nil {
		return "", err
	}

	nowtime := time.Now().Format("2006-01-02_15-04-05")
	dist := "Sqlite3_Backup_" + nowtime + ".sql"
	db.ProjectUUid = ProjectID
	err = db.DumpAllToFile(backupfilePath+"/"+dist, tables)
	if err != nil {
		return "", err
	}
	return dist, err
}
func MysqlSQLTables(host, port, dbname, user, password, char string) ([]string, error) {
	var tablesList []string
	db, err := xorm.NewEngine("mysql", user+":"+password+"@("+host+":"+port+")/"+dbname+"?charset="+char)
	defer db.Close()
	if err != nil {
		return tablesList, err
	}
	err = db.Ping()
	if err != nil {
		return tablesList, err
	}
	tables, _ := db.DBMetas()
	for _, table := range tables {
		tablesList = append(tablesList, table.Name)
	}
	return tablesList, nil
}
func SqliteSQLTables(dbpath string) ([]string, error) {
	var tablesList []string
	db, err := xorm.NewEngine("sqlite3", dbpath)
	defer db.Close()
	if err != nil {
		return tablesList, err
	}
	tables, _ := db.DBMetas()
	for _, table := range tables {
		tablesList = append(tablesList, table.Name)
	}
	return tablesList, err
}
func BackProjectData(ProjectID string) int {
	var code int
	Tables := GetTablesListFunc()
	DbType, _ := config.Int("dbtype")
	if DbType == 1 {
		_, err := SqliteSQLDump("data/db/ism.db", SavePath, "", Tables, ProjectID)
		if err != nil {
			code = -2
		}
	} else if DbType == 4 {
		oceabaseuser, _ := config.String("oceanbaseuser")
		oceabasepwd, _ := config.String("oceanbasepwd")
		oceabasehost, _ := config.String("oceanbasehost")
		oceabaseport, _ := config.String("oceanbaseport")
		oceabasedbname, _ := config.String("oceanbasedbname")
		_, err := MysqlSQLDump(oceabasehost, oceabaseport, oceabasedbname, oceabaseuser, oceabasepwd, "utf8mb4", SavePath, "", Tables, ProjectID)
		if err != nil {
			code = -2
		}
	} else {
		mysqluser, _ := config.String("mysqluser")
		mysqlpwd, _ := config.String("mysqlpwd")
		mysqlhost, _ := config.String("mysqlhost")
		mysqlport, _ := config.String("mysqlport")
		mysqldbname, _ := config.String("mysqldbname")

		_, err := MysqlSQLDump(mysqlhost, mysqlport, mysqldbname, mysqluser, mysqlpwd, "utf8", SavePath, "", Tables, ProjectID)
		if err != nil {
			code = -2
		}
	}
	return code
}
func GetTablesListFunc() []string {
	var results string
	var getTablesList = make([]string, 0)
	dbtype, _ := config.Int("dbtype")
	if dbtype == 1 {
		rows2, _ := models.Db.Raw("select name from sqlite_master where type='table' order by name").Rows()
		defer rows2.Close()
		for rows2.Next() {
			rows2.Scan(&results)
			getTablesList = append(getTablesList, results)
		}
	} else if dbtype == 0 || dbtype == 4 {
		rows2, _ := models.Db.Raw("show tables;").Rows()
		defer rows2.Close()
		for rows2.Next() {
			rows2.Scan(&results)
			getTablesList = append(getTablesList, results)
		}
	}
	if len(getTablesList) == 0 {
		var tableslist string = "alarm_notice,alarm_trigger,custom_data,device_real_data,devices_alarm_list,devices_history_data_list,devices_model,devices_support_list,display_model_layer,display_models,modbus_devices_data_model,modbus_devices_register_group,monitor_list,opcua_devices_data_model,project_lists,project_user,project_video_list,roles_list,snmp_devices_data_model,static_data,system_data_model,system_imge,system_journal,user"
		getTablesList = strings.Split(tableslist, ",")
	}
	return getTablesList
}
func (c *DbOptController) GetTablesList() {

	result := map[string]interface{}{
		"code": nil,
		"list": GetTablesListFunc(),
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *DbOptController) DbBackUp() {

	type TablesStu struct {
		Tables    []string `json:"tables"`
		ProjectID string   `json:"ProjectID"`
	}
	var code int = 0
	var getParams TablesStu

	_, errMk := os.Stat(SavePath)

	if os.IsNotExist(errMk) {
		os.Mkdir(SavePath, os.ModePerm)
	}

	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &getParams)
	if err != nil {
		code = errmsg.NOTJSON
	} else {
		DbType, _ := config.Int("dbtype")
		if DbType == 1 {
			_, err := SqliteSQLDump("data/db/ism.db", SavePath, "", getParams.Tables, "")
			if err != nil {
				code = -2
			}
		} else if DbType == 4 {
			oceabaseuser, _ := config.String("oceanbaseuser")
			oceabasepwd, _ := config.String("oceanbasepwd")
			oceabasehost, _ := config.String("oceanbasehost")
			oceabaseport, _ := config.String("oceanbaseport")
			oceabasedbname, _ := config.String("oceanbasedbname")
			_, err := MysqlSQLDump(oceabasehost, oceabaseport, oceabasedbname, oceabaseuser, oceabasepwd, "utf8mb4", SavePath, "", getParams.Tables, "")
			if err != nil {
				code = -2
			}
		} else {
			mysqluser, _ := config.String("mysqluser")
			mysqlpwd, _ := config.String("mysqlpwd")
			mysqlhost, _ := config.String("mysqlhost")
			mysqlport, _ := config.String("mysqlport")
			mysqldbname, _ := config.String("mysqldbname")

			_, err := MysqlSQLDump(mysqlhost, mysqlport, mysqldbname, mysqluser, mysqlpwd, "utf8", SavePath, "", getParams.Tables, "")
			if err != nil {
				code = -2
			}
		}

	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func formatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
func (c *DbOptController) GetBackUpList() {

	type backUpList struct {
		FileName   string `json:"FileName"`
		CreateTime string `json:"CreateTime"`
		FilePath   string `json:"FilePath"`
		FileSize   string `json:"FileSize"`
	}
	var getBackUpList []backUpList

	var err error
	var code int
	err = filepath.Walk(SavePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			var files backUpList
			files.FileName = info.Name()
			files.CreateTime = info.ModTime().Format("2006-01-02 15:04:05")
			files.FilePath = path
			files.FileSize = formatFileSize(info.Size())
			getBackUpList = append(getBackUpList, files)
		}
		return nil
	})

	if err != nil {
		code = -1
	}
	result := map[string]interface{}{
		"code": code,
		"list": getBackUpList,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func SqliteSQLImport(dbpath string, dbName string) int {
	db, err := xorm.NewEngine("sqlite3", dbpath)
	defer db.Close()
	if err != nil {
		return -1
	}
	_, err1 := db.ImportFile(dbName)
	if err1 != nil {
		return -3
	}
	return 0
}
func ZipFiles(filename string, files []string, oldform, newform string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer newZipFile.Close()
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 把files添加到zip中
	for _, file := range files {
		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.Base(file)
		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}

func MysqlSQLImport(host, port, dbname, user, password, char string, dbName string) int {
	db, err := xorm.NewEngine("mysql", user+":"+password+"@("+host+":"+port+")/"+dbname+"?charset="+char)
	defer db.Close()
	if err != nil {
		return -1
	}
	_, err1 := db.ImportFile(dbName)
	if err1 != nil {
		return -3
	}
	return 0
}
func (c *DbOptController) DbRestore() {

	type RestoreStu struct {
		DbFilePath string `json:"DbFilePath"`
	}
	var code int = 0
	var getParams RestoreStu
	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &getParams)
	if err != nil {
		code = errmsg.NOTJSON
	} else {
		DbType, _ := config.Int("dbtype")
		protocolCommonFunc.CloseChanel()
		protocolCommon.IsRestoreDb = 1
		if DbType == 1 {
			err1 := SqliteSQLImport("data/db/ism.db", getParams.DbFilePath)
			if err1 != 0 {
				code = err1
			}
		} else if DbType == 4 {
			oceabaseuser, _ := config.String("oceanbaseuser")
			oceabasepwd, _ := config.String("oceanbasepwd")
			oceabasehost, _ := config.String("oceanbasehost")
			oceabaseport, _ := config.String("oceanbaseport")
			oceabasedbname, _ := config.String("oceanbasedbname")
			err1 := MysqlSQLImport(oceabasehost, oceabaseport, oceabasedbname, oceabaseuser, oceabasepwd, "utf8mb4", getParams.DbFilePath)
			if err1 != 0 {
				code = err1
			}
		} else {
			mysqluser, _ := config.String("mysqluser")
			mysqlpwd, _ := config.String("mysqlpwd")
			mysqlhost, _ := config.String("mysqlhost")
			mysqlport, _ := config.String("mysqlport")
			mysqldbname, _ := config.String("mysqldbname")

			err1 := MysqlSQLImport(mysqlhost, mysqlport, mysqldbname, mysqluser, mysqlpwd, "utf8", getParams.DbFilePath)
			if err1 != 0 {
				code = err1
			}
		}
		ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
		WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "还原了数据库"+getParams.DbFilePath, errmsg.JournalLevelInfo, c.Ctx.Input)
		models.CheckAllTables()
	}

	result := map[string]interface{}{
		"code": code,
	}
	protocolCommon.IsRestoreDb = 0
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}

func (c *DbOptController) GetDbConfig() {

	type DbConfigStu struct {
		DbType         int    `json:"DbType"`
		Mysqluser      string `json:"Mysqluser"`
		Mysqlpwd       string `json:"Mysqlpwd"`
		Mysqlhost      string `json:"Mysqlhost"`
		Mysqlport      string `json:"Mysqlport"`
		Mysqldbname    string `json:"Mysqldbname"`
		Oceanbaseuser  string `json:"Oceanbaseuser"`
		Oceanbasepwd   string `json:"Oceanbasepwd"`
		Oceanbasehost  string `json:"Oceanbasehost"`
		Oceanbaseport  string `json:"Oceanbaseport"`
		Oceanbasedbname string `json:"Oceanbasedbname"`
	}
	var Dbconfig DbConfigStu

	Dbconfig.DbType, _ = config.Int("dbtype")
	Dbconfig.Mysqluser, _ = config.String("mysqluser")
	Dbconfig.Mysqlpwd, _ = config.String("mysqlpwd")
	Dbconfig.Mysqlhost, _ = config.String("mysqlhost")
	Dbconfig.Mysqlport, _ = config.String("mysqlport")
	Dbconfig.Mysqldbname, _ = config.String("mysqldbname")
	Dbconfig.Oceanbaseuser, _ = config.String("oceanbaseuser")
	Dbconfig.Oceanbasepwd, _ = config.String("oceanbasepwd")
	Dbconfig.Oceanbasehost, _ = config.String("oceanbasehost")
	Dbconfig.Oceanbaseport, _ = config.String("oceanbaseport")
	Dbconfig.Oceanbasedbname, _ = config.String("oceanbasedbname")

	result := map[string]interface{}{
		"code":   0,
		"config": Dbconfig,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}

func (c *DbOptController) SetDbConfig() {

	type DbConfigStu struct {
		DbType         int    `json:"DbType"`
		Mysqluser      string `json:"Mysqluser"`
		Mysqlpwd       string `json:"Mysqlpwd"`
		Mysqlhost      string `json:"Mysqlhost"`
		Mysqlport      string `json:"Mysqlport"`
		Mysqldbname    string `json:"Mysqldbname"`
		Oceanbaseuser  string `json:"Oceanbaseuser"`
		Oceanbasepwd   string `json:"Oceanbasepwd"`
		Oceanbasehost  string `json:"Oceanbasehost"`
		Oceanbaseport  string `json:"Oceanbaseport"`
		Oceanbasedbname string `json:"Oceanbasedbname"`
	}
	var code int = 0
	var getParams DbConfigStu
	data := c.Ctx.Input.RequestBody
	//json数据封装到对象中
	err := json.Unmarshal(data, &getParams)
	if err != nil {
		code = errmsg.NOTJSON
	} else {
		config.Set("dbtype", fmt.Sprintf("%d", getParams.DbType))
		if getParams.DbType == 0 {
			connstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
				getParams.Mysqluser,
				getParams.Mysqlpwd,
				getParams.Mysqlhost,
				getParams.Mysqlport,
				getParams.Mysqldbname)

			_, err = gorm.Open(mysql.Open(connstr))
			if err == nil {
				config.Set("mysqluser", getParams.Mysqluser)
				config.Set("mysqlpwd", getParams.Mysqlpwd)
				config.Set("mysqlhost", getParams.Mysqlhost)
				config.Set("mysqlport", getParams.Mysqlport)
				config.Set("mysqldbname", getParams.Mysqldbname)
				config.SaveConfigFile("conf/app.conf")
				sqldb, _ := models.Db.DB()
				sqldb.Close()
				models.ReconnectDbServer()
			} else {
				code = -3
			}

		} else if getParams.DbType == 4 {
			connstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				getParams.Oceanbaseuser,
				getParams.Oceanbasepwd,
				getParams.Oceanbasehost,
				getParams.Oceanbaseport,
				getParams.Oceanbasedbname)

			_, err = gorm.Open(mysql.Open(connstr))
			if err == nil {
				config.Set("oceanbaseuser", getParams.Oceanbaseuser)
				config.Set("oceanbasepwd", getParams.Oceanbasepwd)
				config.Set("oceanbasehost", getParams.Oceanbasehost)
				config.Set("oceanbaseport", getParams.Oceanbaseport)
				config.Set("oceanbasedbname", getParams.Oceanbasedbname)
				config.SaveConfigFile("conf/app.conf")
				sqldb, _ := models.Db.DB()
				sqldb.Close()
				models.ReconnectDbServer()
			} else {
				code = -3
			}

		} else if getParams.DbType == 3 {
			// 达梦数据库在 macOS 上不可用
			code = -3
		} else {
			sqldb, _ := models.Db.DB()
			sqldb.Close()
			models.ReconnectDbServer()
			config.SaveConfigFile("conf/app.conf")
		}
		if code == 0 {
			protocolCommonFunc.CloseChanel()
			ISMScript.ScriptCloseChan()
			staticDataTask.PushStaticCloseChan()
		}
	}

	result := map[string]interface{}{
		"code": code,
	}

	c.Data["json"] = result

	c.ServeJSON() //返回json格式
}
func (c *DbOptController) DbDown() {

	var fileslist []string
	var saveFilePath string = "static/"
	var filePath string
	type RestoreStu struct {
		DbFilePath string `json:"DbFilePath"`
	}
	var code int = 0
	var getParams RestoreStu
	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &getParams)
	if err != nil {
		code = errmsg.NOTJSON
	} else {
		downName := filepath.Base(getParams.DbFilePath[:len(getParams.DbFilePath)-len(filepath.Ext(getParams.DbFilePath))])
		fileslist = append(fileslist, getParams.DbFilePath)
		filePath = saveFilePath + downName + ".zip"
		ZipFiles(filePath, fileslist, "data\\dbbackup\\", "")
	}

	result := map[string]interface{}{
		"code": code,
		"path": filePath,
	}
	protocolCommon.IsRestoreDb = 0
	c.Data["json"] = result
	c.ServeJSON() //返回json格式
}
