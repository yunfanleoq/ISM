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
	"ISMServer/utils/backupverify"
	"ISMServer/utils/errmsg"
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/beego/beego/v2/adapter/logs"
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

// 备份完整性校验关注的关键表（见 utils/backupverify.CriticalTables）
var _ = backupverify.CriticalTables

// 客户端/会话侧抬高包上限，避免大 INSERT（组态/能源快照）Error 1153
// 现场 315MB+ 备份中大 INSERT 需 ≥256MB；取 512MB 留余量（会话 SET + DSN maxAllowedPacket）
const mysqlMaxAllowedPacket = 512 * 1024 * 1024

// OceanBase 专有：整表 dump 会撞默认 ob_query_timeout=10s（Error 4012）
// 单位微秒；1 小时。非 OB 上 SET 失败仅 Warn。
const obQueryTimeoutUs = int64(3600 * 1000 * 1000)

type backupTableCount struct {
	Table    string `json:"table"`
	DbCount  int64  `json:"dbCount"`
	SqlCount int64  `json:"sqlCount"`
	Ok       bool   `json:"ok"`
}

func rejectEmptyBackupFile(fullPath string) error {
	info, err := os.Stat(fullPath)
	if err != nil {
		return err
	}
	if info.Size() <= 0 {
		_ = os.Remove(fullPath)
		return fmt.Errorf("备份文件为空(0B)，已删除: %s", fullPath)
	}
	return nil
}

// verifyBackupSQLCounts 对比关键表库内 COUNT 与 SQL INSERT 行数；SQL 少于库内则判失败并删除文件。
func verifyBackupSQLCounts(sqlPath string, tables []string, projectID string) ([]backupTableCount, error) {
	raw, err := backupverify.VerifySQLCounts(sqlPath, tables, projectID)
	results := make([]backupTableCount, 0, len(raw))
	for _, c := range raw {
		results = append(results, backupTableCount{Table: c.Table, DbCount: c.DbCount, SqlCount: c.SqlCount, Ok: c.Ok})
	}
	return results, err
}

func mysqlDumpDSN(user, password, host, port, dbname, char string) string {
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&maxAllowedPacket=%d&timeout=120s&readTimeout=3600s&writeTimeout=3600s",
		user, password, host, port, dbname, char, mysqlMaxAllowedPacket)
}

// applyOceanBaseDumpSessionTimeouts 抬高 OB 查询/事务超时，避免 dump 大表 Error 4012。
// 纯 MySQL / SQLite 上变量不存在，失败仅打 Warn。
func applyOceanBaseDumpSessionTimeouts(engine *xorm.Engine) {
	if engine == nil {
		return
	}
	if _, err := engine.Exec(fmt.Sprintf("SET SESSION ob_query_timeout=%d", obQueryTimeoutUs)); err != nil {
		logs.Warn("SET SESSION ob_query_timeout failed (continue, non-OB?): %v", err)
	}
	if _, err := engine.Exec(fmt.Sprintf("SET SESSION ob_trx_timeout=%d", obQueryTimeoutUs)); err != nil {
		logs.Warn("SET SESSION ob_trx_timeout failed (continue, non-OB?): %v", err)
	}
}

func isObQueryTimeoutErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "4012") ||
		strings.Contains(msg, "ob_query_timeout") ||
		strings.Contains(msg, "maximum query timeout")
}

// MysqlSQLDump ...
func MysqlSQLDump(host, port, dbname, user, password, char, backupfilePath, zippwd string, tables []string, ProjectID string) (string, error) {
	db, err := xorm.NewEngine("mysql", mysqlDumpDSN(user, password, host, port, dbname, char))
	if err != nil {
		return "", err
	}
	defer db.Close()

	applyOceanBaseDumpSessionTimeouts(db)

	nowtime := time.Now().Format("2006-01-02_15-04-05")
	dist := "Mysql_Backup_" + nowtime + ".sql"
	fullPath := filepath.Join(backupfilePath, dist)
	db.ProjectUUid = ProjectID
	err = db.DumpAllToFile(fullPath, tables)
	if err != nil {
		_ = os.Remove(fullPath)
		return "", err
	}
	if err := rejectEmptyBackupFile(fullPath); err != nil {
		return "", err
	}
	return dist, nil
}

func SqliteSQLDump(dbpath, backupfilePath, zippwd string, tables []string, ProjectID string) (string, error) {
	db, err := xorm.NewEngine("sqlite3", dbpath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	nowtime := time.Now().Format("2006-01-02_15-04-05")
	dist := "Sqlite3_Backup_" + nowtime + ".sql"
	fullPath := filepath.Join(backupfilePath, dist)
	db.ProjectUUid = ProjectID
	err = db.DumpAllToFile(fullPath, tables)
	if err != nil {
		_ = os.Remove(fullPath)
		return "", err
	}
	if err := rejectEmptyBackupFile(fullPath); err != nil {
		return "", err
	}
	return dist, nil
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
	var fileName string
	var err error
	if DbType == 1 {
		fileName, err = SqliteSQLDump("data/db/ism.db", SavePath, "", Tables, ProjectID)
	} else if DbType == 4 {
		oceabaseuser, _ := config.String("oceanbaseuser")
		oceabasepwd, _ := config.String("oceanbasepwd")
		oceabasehost, _ := config.String("oceanbasehost")
		oceabaseport, _ := config.String("oceanbaseport")
		oceabasedbname, _ := config.String("oceanbasedbname")
		fileName, err = MysqlSQLDump(oceabasehost, oceabaseport, oceabasedbname, oceabaseuser, oceabasepwd, "utf8mb4", SavePath, "", Tables, ProjectID)
	} else {
		mysqluser, _ := config.String("mysqluser")
		mysqlpwd, _ := config.String("mysqlpwd")
		mysqlhost, _ := config.String("mysqlhost")
		mysqlport, _ := config.String("mysqlport")
		mysqldbname, _ := config.String("mysqldbname")
		fileName, err = MysqlSQLDump(mysqlhost, mysqlport, mysqldbname, mysqluser, mysqlpwd, "utf8", SavePath, "", Tables, ProjectID)
	}
	if err != nil {
		logs.Error("BackProjectData dump failed: %v", err)
		return -2
	}
	if _, vErr := verifyBackupSQLCounts(filepath.Join(SavePath, fileName), Tables, ProjectID); vErr != nil {
		logs.Error("BackProjectData verify failed: %v", vErr)
		return -6
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
	var fileName string
	var fileSize int64
	var tableCounts []backupTableCount
	var errMsg string

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
		var dumpErr error
		if DbType == 1 {
			fileName, dumpErr = SqliteSQLDump("data/db/ism.db", SavePath, "", getParams.Tables, "")
		} else if DbType == 4 {
			oceabaseuser, _ := config.String("oceanbaseuser")
			oceabasepwd, _ := config.String("oceanbasepwd")
			oceabasehost, _ := config.String("oceanbasehost")
			oceabaseport, _ := config.String("oceanbaseport")
			oceabasedbname, _ := config.String("oceanbasedbname")
			fileName, dumpErr = MysqlSQLDump(oceabasehost, oceabaseport, oceabasedbname, oceabaseuser, oceabasepwd, "utf8mb4", SavePath, "", getParams.Tables, "")
		} else {
			mysqluser, _ := config.String("mysqluser")
			mysqlpwd, _ := config.String("mysqlpwd")
			mysqlhost, _ := config.String("mysqlhost")
			mysqlport, _ := config.String("mysqlport")
			mysqldbname, _ := config.String("mysqldbname")
			fileName, dumpErr = MysqlSQLDump(mysqlhost, mysqlport, mysqldbname, mysqluser, mysqlpwd, "utf8", SavePath, "", getParams.Tables, "")
		}
		if dumpErr != nil {
			if isObQueryTimeoutErr(dumpErr) {
				code = -7
				errMsg = "备份查询超时(OceanBase ob_query_timeout)：大表 dump 超过会话时限。请确认 ism_server 已含 SET SESSION ob_query_timeout，或执行 scripts/tune_ob_max_allowed_packet.sh 抬高全局超时后重试。详情: " + dumpErr.Error()
			} else {
				code = -2
				errMsg = dumpErr.Error()
			}
			logs.Error("DbBackUp dump failed: %v", dumpErr)
		} else {
			fullPath := filepath.Join(SavePath, fileName)
			if info, stErr := os.Stat(fullPath); stErr == nil {
				fileSize = info.Size()
			}
			counts, vErr := verifyBackupSQLCounts(fullPath, getParams.Tables, "")
			tableCounts = counts
			if vErr != nil {
				code = -6
				errMsg = vErr.Error()
				fileName = ""
				fileSize = 0
				logs.Error("DbBackUp verify failed: %v", vErr)
			}
		}
	}

	result := map[string]interface{}{
		"code":        code,
		"fileName":    fileName,
		"fileSize":    fileSize,
		"fileSizeStr": formatFileSize(fileSize),
		"tableCounts": tableCounts,
		"msg":         errMsg,
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
	if err != nil {
		return -1
	}
	defer db.Close()
	_, err1 := db.ImportFile(dbName)
	if err1 != nil {
		fmt.Printf("SqliteSQLImport failed: %v\n", err1)
		logs.Error("SqliteSQLImport failed: %v", err1)
		if strings.Contains(err1.Error(), "missing xorm magic") {
			return -5
		}
		if strings.Contains(err1.Error(), "partial failure") {
			return -4
		}
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

// lastMysqlImportErr 供 DbRestore 区分主键冲突 / packet / 其它部分失败文案。
var lastMysqlImportErr error

// clearMysqlTablesBeforeImport 还原前清空目标库表，避免 INSERT 撞已有 PRIMARY（Error 1062）。
func clearMysqlTablesBeforeImport(db *xorm.Engine) error {
	if db == nil {
		return fmt.Errorf("nil engine")
	}
	sqlDB := db.DB()
	if sqlDB == nil {
		return fmt.Errorf("nil sql.DB")
	}
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS=0")
	rows, err := sqlDB.Query("SHOW TABLES")
	if err != nil {
		return err
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var name string
		if scanErr := rows.Scan(&name); scanErr != nil {
			continue
		}
		if name == "" {
			continue
		}
		tables = append(tables, name)
	}
	for _, t := range tables {
		// TRUNCATE 失败（权限/引擎）则 DELETE
		if _, terr := sqlDB.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", t)); terr != nil {
			if _, derr := sqlDB.Exec(fmt.Sprintf("DELETE FROM `%s`", t)); derr != nil {
				logs.Warn("restore clear table %s failed: truncate=%v delete=%v", t, terr, derr)
			}
		}
	}
	_, _ = sqlDB.Exec("SET FOREIGN_KEY_CHECKS=1")
	logs.Info("restore: cleared %d tables before import", len(tables))
	return nil
}

func classifyMysqlImportCode(err error) int {
	if err == nil {
		return 0
	}
	msg := err.Error()
	if strings.Contains(msg, "missing xorm magic") {
		return -5
	}
	if strings.Contains(msg, "Duplicate entry") || strings.Contains(msg, "1062") {
		return -4
	}
	if strings.Contains(msg, "partial failure") {
		return -4
	}
	if strings.Contains(msg, "max_allowed_packet") || strings.Contains(msg, "1153") {
		return -4
	}
	return -3
}

func restoreFailureMsg(code int, importErr error) string {
	if code != -4 {
		return fmt.Sprintf("restore failed code=%d", code)
	}
	msg := ""
	if importErr != nil {
		msg = importErr.Error()
	}
	if strings.Contains(msg, "Duplicate entry") || strings.Contains(msg, "1062") {
		return "还原失败：主键冲突(Duplicate entry)。已尝试清空表后仍冲突，请检查备份 SQL 是否自带重复主键，或联系支持。"
	}
	if strings.Contains(msg, "max_allowed_packet") || strings.Contains(msg, "1153") {
		return "还原失败：语句过大(max_allowed_packet)。请执行 scripts/tune_ob_max_allowed_packet.sh 抬高 GLOBAL 后重试（OceanBase 上 SESSION 该变量只读）。"
	}
	return "还原数据库失败（含语句错误），未当作成功。请查看 ism_server 日志中 MysqlSQLImport 详情后重试。"
}

func MysqlSQLImport(host, port, dbname, user, password, char string, dbName string) int {
	lastMysqlImportErr = nil
	db, err := xorm.NewEngine("mysql", mysqlDumpDSN(user, password, host, port, dbname, char))
	if err != nil {
		lastMysqlImportErr = err
		return -1
	}
	defer db.Close()
	// OceanBase：SESSION max_allowed_packet 常为只读，优先抬 GLOBAL（需权限；失败仅 Warn）
	if _, setErr := db.DB().Exec(fmt.Sprintf("SET GLOBAL max_allowed_packet=%d", mysqlMaxAllowedPacket)); setErr != nil {
		logs.Warn("SET GLOBAL max_allowed_packet failed (continue; 请执行 scripts/tune_ob_max_allowed_packet.sh): %v", setErr)
	}
	if _, setErr := db.DB().Exec(fmt.Sprintf("SET SESSION max_allowed_packet=%d", mysqlMaxAllowedPacket)); setErr != nil {
		logs.Warn("SET SESSION max_allowed_packet failed (OB 上多为只读，依赖 GLOBAL): %v", setErr)
	}
	applyOceanBaseDumpSessionTimeouts(db)
	if clearErr := clearMysqlTablesBeforeImport(db); clearErr != nil {
		logs.Warn("restore clear tables failed (continue import): %v", clearErr)
	}
	_, err1 := db.ImportFile(dbName)
	if err1 != nil {
		lastMysqlImportErr = err1
		fmt.Printf("MysqlSQLImport failed: %v\n", err1)
		logs.Error("MysqlSQLImport failed: %v", err1)
		return classifyMysqlImportCode(err1)
	}
	return 0
}
func (c *DbOptController) DbRestore() {

	type RestoreStu struct {
		DbFilePath string `json:"DbFilePath"`
	}
	var code int = 0
	var syncCreated, syncSkipped int
	var getParams RestoreStu
	var errMsg string
	data := c.Ctx.Input.RequestBody

	//json数据封装到对象中
	err := json.Unmarshal(data, &getParams)
	if err != nil {
		code = errmsg.NOTJSON
	} else {
		DbType, _ := config.Int("dbtype")
		// 先置还原标志，再停采集/脚本，避免 DROP 表时仍有并发读 → Error 1146
		protocolCommon.IsRestoreDb = 1
		protocolCommonFunc.CloseChanel()
		ISMScript.ScriptCloseChan()
		time.Sleep(2 * time.Second)
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
		if code == 0 {
			models.CheckAllTables()
			// 重建 GORM 连接，避免还原用独立 xorm 引擎写库后运行时连接池读到旧状态
			models.ReconnectDbServer()
			syncCreated, syncSkipped = models.SyncAllProjectsDeviceRealData()
			protocolCommonFunc.CloseChanel()
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid,
				fmt.Sprintf("还原后补建实时点位：创建%d条, 跳过%d条", syncCreated, syncSkipped),
				errmsg.JournalLevelInfo, c.Ctx.Input)
		} else {
			errMsg = restoreFailureMsg(code, lastMysqlImportErr)
			logs.Error("DbRestore failed: code=%d path=%s msg=%s", code, getParams.DbFilePath, errMsg)
			WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid,
				fmt.Sprintf("还原数据库失败 code=%d file=%s", code, getParams.DbFilePath),
				errmsg.JournalLevelInfo, c.Ctx.Input)
		}
	}

	result := map[string]interface{}{
		"code":        code,
		"syncCreated": syncCreated,
		"syncSkipped": syncSkipped,
		"msg":         errMsg,
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
func (c *DbOptController) DbDeleteBackup() {
	type DelStu struct {
		DbFilePath string `json:"DbFilePath"`
	}
	result := map[string]interface{}{"code": errmsg.ERROR}
	var getParams DelStu
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &getParams); err != nil {
		result["code"] = errmsg.NOTJSON
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	if getParams.DbFilePath == "" {
		result["msg"] = "empty path"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	absSave, _ := filepath.Abs(SavePath)
	absTarget, err := filepath.Abs(getParams.DbFilePath)
	if err != nil {
		result["msg"] = "invalid path"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	rel, err := filepath.Rel(absSave, absTarget)
	if err != nil || strings.HasPrefix(rel, "..") {
		result["code"] = -9
		result["msg"] = "path outside backup dir"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	info, err := os.Stat(absTarget)
	if err != nil || info.IsDir() {
		result["msg"] = "file not found"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	if err := os.Remove(absTarget); err != nil {
		result["msg"] = err.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	ProjectUuid := c.Ctx.Request.Header.Get("ProjectUuid")
	WriteOperationJournal(c.Ctx.Request.Header.Get("Authorization"), ProjectUuid, "删除备份 "+filepath.Base(absTarget), errmsg.JournalLevelInfo, c.Ctx.Input)
	result["code"] = errmsg.SUCCSECODE
	c.Data["json"] = result
	c.ServeJSON()
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
		result := map[string]interface{}{
			"code": code,
			"path": "",
		}
		protocolCommon.IsRestoreDb = 0
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	if getParams.DbFilePath == "" {
		result := map[string]interface{}{
			"code": errmsg.ERROR,
			"path": "",
		}
		protocolCommon.IsRestoreDb = 0
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	downName := filepath.Base(getParams.DbFilePath[:len(getParams.DbFilePath)-len(filepath.Ext(getParams.DbFilePath))])
	fileslist = append(fileslist, getParams.DbFilePath)
	filePath = saveFilePath + downName + ".zip"
	if zipErr := ZipFiles(filePath, fileslist, "data\\dbbackup\\", ""); zipErr != nil {
		result := map[string]interface{}{
			"code": errmsg.ERROR,
			"path": "",
			"msg":  zipErr.Error(),
		}
		protocolCommon.IsRestoreDb = 0
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	protocolCommon.IsRestoreDb = 0
	// 直接流式返回 zip，避免前端再请求 /static 导致 7080 上 404
	c.Ctx.Output.Download(filePath, downName+".zip")
}
