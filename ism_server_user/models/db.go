/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-21 10:38:34
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	protocol_common "ISMServer/protocol/common"
	ISMNetwork "ISMServer/task/network"
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/config"
	gormlog "github.com/beego/beego/v2/core/logs"
	"github.com/go-basic/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	Db *gorm.DB
)

type Writer struct {
}

func (w Writer) Printf(format string, args ...interface{}) {
	// log.Infof(format, args...)
	logString := fmt.Sprintf(format, args...)

	gormlog.Info(logString)
}

func init() {
	//初始化网络
	ISMNetwork.InitNetwork()
	ReconnectDbServer()
}

func CheckAllTables() {
	// 自动迁移
	r := time.Now()
	gormlog.Info("正在检查系统表,请稍等......")
	err := Db.AutoMigrate(&BacnetDevicesDataModel{}, &SQLReportTemplete{}, &CJT188DevicesDataModel{}, &VirtualDeviceDataModel{}, &ModbusTcpDataPushModel{}, &IEC104DataPushModel{}, &SystemDataInterface{}, &SystemDataTemplete{}, &DisplayModelsUserList{}, &HJ212DevicesDataModel{}, &OutConnectList{}, &IEC61850DevicesDataModel{}, &IEC104DevicesDataModel{}, &Dlt645DevicesDataModel{}, &MqttDevicesDataModel{}, &ISMScript{}, &ReportTemplete{}, &TaskPlanList{}, &SimS7DataModel{}, &UserApiAccessToken{}, &RESTFulDataModel{}, &CustomData{}, &SystemDataModel{}, &OpcuaDevicesDataModel{}, &RolesList{}, &AlarmNotice{}, &SystemJournal{}, &StaticData{}, &ProjectVideoList{}, &ProjectUser{}, &ProjectLists{}, &AlarmTrigger{}, &DevicesHistoryDataList{}, &DevicesAlarmList{}, &SystemImge{}, &ModbusDevicesDataModel{}, &ModbusDevicesRegisterGroup{}, &User{}, &DisplayModelLayer{}, &DevicesModel{}, &SnmpDevicesDataModel{}, &MonitorList{}, &DeviceRealData{}, &DevicesSupportList{}, &DisplayModels{})
	if err != nil {
		gormlog.Info(err)
	}
	d := time.Since(r)
	gormlog.Info("系统表检查完成,耗时:%s", d)
}
func ReconnectDbServer() {
	var err error
	var adminUser User
	var adminRole RolesList
	var dbConnect *gorm.DB

	dbConnect = nil
	sqlLogger := logger.New(
		Writer{},
		logger.Config{
			SlowThreshold:             120000 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Error,              // Log level
			IgnoreRecordNotFoundError: true,                      // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                      // Disable color
		},
	)
	var dbConfig = gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 外键约束
		SkipDefaultTransaction:                   true, // 禁用默认事务（提高运行速度）
		PrepareStmt:                              false,
		Logger:                                   sqlLogger,
		CreateBatchSize:                          3000,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	}
	var (
		dbtype, _      = config.Int("dbtype")
		mysqluser, _   = config.String("mysqluser")
		mysqlpwd, _    = config.String("mysqlpwd")
		mysqlhost, _   = config.String("mysqlhost")
		mysqlport, _   = config.String("mysqlport")
		mysqldbname, _ = config.String("mysqldbname")

		postgresuser, _   = config.String("postgresuser")
		postgrespwd, _    = config.String("postgrespwd")
		postgreshost, _   = config.String("postgreshost")
		postgresport, _   = config.String("postgresport")
		postgresdbname, _ = config.String("postgresdbname")

		oceanbaseuser, _   = config.String("oceanbaseuser")
		oceanbasepwd, _    = config.String("oceanbasepwd")
		oceanbasehost, _   = config.String("oceanbasehost")
		oceanbaseport, _   = config.String("oceanbaseport")
		oceanbasedbname, _ = config.String("oceanbasedbname")

		historyKeepDays, _ = config.Int("history_keep_days")
	)
	if historyKeepDays > 0 {
		protocol_common.HistoryKeepDays = historyKeepDays
	}

	gormlog.Info("正在连接数据库,请稍等......")

	if dbtype == 1 {
		dbConnect, err = gorm.Open(sqlite.Open("data/db/ism.db"), &dbConfig)
		if err == nil {
			// 启用 WAL 模式：允许读写并发，大幅提升稳定性
			dbConnect.Exec("PRAGMA journal_mode=WAL")
			dbConnect.Exec("PRAGMA synchronous=NORMAL")
			dbConnect.Exec("PRAGMA busy_timeout=5000")
			dbConnect.Exec("PRAGMA cache_size=-8000") // 8MB cache
		}
	} else if dbtype == 0 {
		// 数据库连接字符串
		connstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			mysqluser,
			mysqlpwd,
			mysqlhost,
			mysqlport,
			mysqldbname)
		dbConnect, err = gorm.Open(mysql.Open(connstr), &dbConfig)
	} else if dbtype == 2 {
		// 数据库连接字符串

		connstr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			postgreshost,
			postgresuser,
			postgrespwd,
			postgresdbname,
			postgresport,
		)

		dbConnect, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  connstr,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &dbConfig)
	} else if dbtype == 3 {
		// 达梦数据库在 macOS 上不可用，跳过
		panic("达梦数据库不支持当前平台，请使用 SQLite(dbtype=1)、MySQL(dbtype=0) 或 PostgreSQL(dbtype=2)")
	} else if dbtype == 4 {
		// OceanBase MySQL 兼容模式
		connstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s",
			oceanbaseuser,
			oceanbasepwd,
			oceanbasehost,
			oceanbaseport,
			oceanbasedbname)
		dbConnect, err = gorm.Open(mysql.Open(connstr), &dbConfig)
	}
	//把系统使用的数据库类型作为全局去使用
	protocol_common.InsideDbType = dbtype
	if err != nil {
		panic("数据库连接失败！")
	}
	gormlog.Info("数据库连接成功")
	Db = dbConnect
	if dbtype != 3 {
		CheckAllTables()
	}
	sqldb, _ := dbConnect.DB()

	if dbtype == 1 {
		// SetMaxIdleCons 设置连接池中的最大闲置连接数。
		sqldb.SetMaxIdleConns(1)

		// SetMaxOpenCons 设置数据库的最大连接数量。
		sqldb.SetMaxOpenConns(1)

	} else if dbtype == 4 {
		// OceanBase 连接池配置：支持较高并发
		sqldb.SetMaxIdleConns(20)
		sqldb.SetMaxOpenConns(100)
		sqldb.SetConnMaxLifetime(60 * time.Second)
	} else {
		// SetMaxIdleCons 设置连接池中的最大闲置连接数。
		sqldb.SetMaxIdleConns(10)

		// SetMaxOpenCons 设置数据库的最大连接数量。
		sqldb.SetMaxOpenConns(80)

	}

	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	if dbtype != 4 {
		sqldb.SetConnMaxLifetime(30 * time.Second)
	}
	protocol_common.ClearAlarmType, err = config.Int("clearalarmtype")
	if err != nil {
		protocol_common.ClearAlarmType = 0 //默认清除告警
	}
	if protocol_common.ClearAlarmType == 0 {
		//把告警没有结束的全部结束,不再更新keeptime
		var clearAlarm []DevicesAlarmList
		Db.Model(&DevicesAlarmList{}).Where("clear_time < ?", "2007-01-02 15:04:05").Find(&clearAlarm)
		for _, device := range clearAlarm {
			var upateAlarm DevicesAlarmList

			upateAlarm.ClearTime = time.Now()
			upateAlarm.KeepTime = (float64)((upateAlarm.ClearTime.UnixMilli() - device.HappenTime.UnixMilli()) / 1000.0)
			Db.Model(&DevicesAlarmList{}).Where("clear_time < ?", "2007-01-02 15:04:05").Updates(&upateAlarm)
		}
	}
	//单机版初始化超级管理员数据
	AdminUserError := Db.Model(&User{}).Where("username = 'admin'").First(&adminUser)

	if errors.Is(AdminUserError.Error, gorm.ErrRecordNotFound) {
		adminUser.Username = "admin"
		adminUser.Password = "$2a$10$h9swLjbTTcSVUCqQDt6nAetw.FVRLPE0WPDzqloprYRO7PDtLC5Ii" // bcrypt(MD5("123456"))
		adminUser.Role = "Admin"
		adminUser.Name = "超级管理员"
		adminUser.Uuid = uuid.New()

		Db.Model(&User{}).Create(&adminUser) // 通过数据的指针来创建
	}

	AdminUserError = Db.Model(&RolesList{}).Where("role_id = ?", "Admin").First(&adminRole)

	if errors.Is(AdminUserError.Error, gorm.ErrRecordNotFound) {
		var insertrole = make([]RolesList, 3)
		insertrole[0].RoleName = "管理员"
		insertrole[0].RoleId = "Admin"
		insertrole[0].Permission = "add,edit,delete,view"

		insertrole[1].RoleName = "操作员"
		insertrole[1].RoleId = "Operator"
		insertrole[1].Permission = "add,edit,delete,view"

		insertrole[2].RoleName = "普通用户"
		insertrole[2].RoleId = "User"
		insertrole[2].Permission = "view"

		Db.Model(&RolesList{}).CreateInBatches(&insertrole, 3) // 通过数据的指针来创建
	}

	var supportDeviceListCount int64
	Db.Model(&DevicesSupportList{}).Where("ID > 0").Count(&supportDeviceListCount)
	if supportDeviceListCount != 15 {
		var delDevice DevicesSupportList
		Db.Unscoped().Model(&DevicesSupportList{}).Where("ID > 0").Delete(&delDevice)

		var insertSupportDevice = make([]DevicesSupportList, 15)
		insertSupportDevice[0].Name = "device.SnmpDevice"
		insertSupportDevice[0].Described = "SNMP设备"
		insertSupportDevice[0].Type = 1

		insertSupportDevice[1].Name = "device.ModbusDevice"
		insertSupportDevice[1].Described = "Modbus设备"
		insertSupportDevice[1].Type = 2

		insertSupportDevice[2].Name = "device.OPCUADevice"
		insertSupportDevice[2].Described = "OPCUA设备"
		insertSupportDevice[2].Type = 3

		insertSupportDevice[3].Name = "device.RESTFulDevice"
		insertSupportDevice[3].Described = "RESTFul设备"
		insertSupportDevice[3].Type = 5

		insertSupportDevice[4].Name = "device.StaticDevice"
		insertSupportDevice[4].Described = "静态数据"
		insertSupportDevice[4].Type = 6

		insertSupportDevice[5].Name = "device.CustomDevice"
		insertSupportDevice[5].Described = "自定义数据"
		insertSupportDevice[5].Type = 7

		insertSupportDevice[6].Name = "device.SimS7Device"
		insertSupportDevice[6].Described = "西门子PLC S7设备"
		insertSupportDevice[6].Type = 15

		insertSupportDevice[7].Name = "device.MqttDevice"
		insertSupportDevice[7].Described = "MQTT设备"
		insertSupportDevice[7].Type = 20

		insertSupportDevice[8].Name = "device.DLT645Device"
		insertSupportDevice[8].Described = "DLT645电表"
		insertSupportDevice[8].Type = 30

		insertSupportDevice[9].Name = "device.IEC104Device"
		insertSupportDevice[9].Described = "IEC104电力规约"
		insertSupportDevice[9].Type = 40

		insertSupportDevice[10].Name = "device.IEC61850Device"
		insertSupportDevice[10].Described = "IEC61850电力规约"
		insertSupportDevice[10].Type = 350

		insertSupportDevice[11].Name = "device.HJ212Device"
		insertSupportDevice[11].Described = "环保2017"
		insertSupportDevice[11].Type = 470

		insertSupportDevice[12].Name = "device.VisDevice"
		insertSupportDevice[12].Described = "虚拟设备"
		insertSupportDevice[12].Type = 480

		insertSupportDevice[13].Name = "device.CJT188Device"
		insertSupportDevice[13].Described = "CJT188协议"
		insertSupportDevice[13].Type = 490

		insertSupportDevice[14].Name = "device.BacnetDevice"
		insertSupportDevice[14].Described = "BACnet协议"
		insertSupportDevice[14].Type = 500

		Db.Model(&DevicesSupportList{}).CreateInBatches(&insertSupportDevice, len(insertSupportDevice)) // 通过数据的指针来创建
	}
	//启动后把全部设备的状态全部更新为离线，然后重新判断,暂时不启用

	//Db.Model(&MonitorList{}).Where("type = 1 and status !=3").Update("status", 0)
}
