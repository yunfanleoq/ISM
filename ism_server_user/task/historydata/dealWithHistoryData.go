/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-07 14:41:49
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package alarmTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/gob"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"gorm.io/driver/postgres"

	"github.com/beego/beego/v2/core/config"
	_ "github.com/taosdata/driver-go/v3/taosRestful"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DeviceHistoryDataTemp = make(map[string]models.DevicesHistoryDataList, protocol_common.HistoryCacheCount)

var OnceWriteHistoryNumber int = 100
var HistoryWriteWorkerCount int = defaultHistoryWriteWorkerCount()

type Writer struct {
}

type historyWriteBatch struct {
	keys      [][]byte
	data      []models.DevicesHistoryDataList
	waitGroup *sync.WaitGroup
}

func (w Writer) Printf(format string, args ...interface{}) {
	// log.Infof(format, args...)
	logString := fmt.Sprintf(format, args...)

	logs.Info(logString)
}

func defaultHistoryWriteWorkerCount() int {
	workerCount := runtime.NumCPU()
	if workerCount < 2 {
		return 2
	}
	if workerCount > 8 {
		return 8
	}
	return workerCount
}

func newHistoryWriteBatch(batchSize int, waitGroup *sync.WaitGroup) historyWriteBatch {
	if batchSize <= 0 {
		batchSize = 100
	}
	return historyWriteBatch{
		keys:      make([][]byte, 0, batchSize),
		data:      make([]models.DevicesHistoryDataList, 0, batchSize),
		waitGroup: waitGroup,
	}
}

func cloneHistoryLevelDbKey(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

func deleteHistoryBatchKeys(keys [][]byte) {
	if protocol_common.GSaveHistoryDataLevelDb == nil {
		return
	}
	for _, key := range keys {
		if err := protocol_common.GSaveHistoryDataLevelDb.Delete(key, nil); err != nil {
			logs.Error(fmt.Sprintf("delete history leveldb key failed: %v", err))
		}
	}
}

func writeInsideHistoryData(writeDeviceHistoryData []models.DevicesHistoryDataList) error {
	if len(writeDeviceHistoryData) == 0 {
		return nil
	}
	if protocol_common.InsideDbType == 0 {
		return insertHistoryData(writeDeviceHistoryData, OnceWriteHistoryNumber)
	}
	return models.Db.Model(&models.DevicesHistoryDataList{}).CreateInBatches(&writeDeviceHistoryData, OnceWriteHistoryNumber).Error
}

func writeTDengineHistoryData(writeDeviceHistoryData []models.DevicesHistoryDataList) error {
	if len(writeDeviceHistoryData) == 0 {
		return nil
	}
	if protocol_common.HistoryRecordTsDb == nil {
		return fmt.Errorf("history record tdengine db is nil")
	}

	var insertSQL strings.Builder
	insertSQL.Grow(len(writeDeviceHistoryData) * 256)
	tagNo := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(0xFFFFFF)
	insertSQL.WriteString(fmt.Sprintf("INSERT INTO ISMHistoryDb.HistoryDatas USING ISMHistoryDb.TempleteHistoryDatas  TAGS(%d) VALUES", tagNo))
	for _, historyData := range writeDeviceHistoryData {
		insertSQL.WriteString(fmt.Sprintf(" ('%s', '%s','%s','%s','%s','%s','%s','%s','%s') ",
			historyData.RecordTime.Format("2006-01-02 15:04:05.0000"),
			historyData.DataName,
			historyData.DeviceUuid,
			historyData.ProjectUuid,
			historyData.DeviceName,
			historyData.DataUuid,
			historyData.ModelDataUuid,
			historyData.DataUnit,
			historyData.DataValue,
		))
	}
	_, err := protocol_common.HistoryRecordTsDb.Exec(insertSQL.String())
	return err
}

func writeInfluxHistoryData(writeDeviceHistoryData []models.DevicesHistoryDataList) error {
	if len(writeDeviceHistoryData) == 0 {
		return nil
	}
	if protocol_common.HistoryRecordInfluxdbDb == nil {
		return fmt.Errorf("history record influxdb db is nil")
	}

	writeInfluxDeviceHistoryData := make([]*write.Point, 0, len(writeDeviceHistoryData))
	for _, historyData := range writeDeviceHistoryData {
		tags := map[string]string{
			"DataName":      historyData.DataName,
			"DeviceUuid":    historyData.DeviceUuid,
			"ProjectUuid":   historyData.ProjectUuid,
			"DeviceName":    historyData.DeviceName,
			"DataUuid":      historyData.DataUuid,
			"ModelDataUuid": historyData.ModelDataUuid,
			"DataUnit":      historyData.DataUnit,
		}
		fields := map[string]interface{}{
			"DataValue": historyData.DataValue,
		}
		writeInfluxDeviceHistoryData = append(writeInfluxDeviceHistoryData, write.NewPoint("ISMHistoryData", tags, fields, historyData.RecordTime))
	}
	return protocol_common.HistoryRecordInfluxdbDb.WritePoint(context.Background(), writeInfluxDeviceHistoryData...)
}

func writePGHistoryData(writeDeviceHistoryData []models.DevicesHistoryDataList) error {
	if len(writeDeviceHistoryData) == 0 {
		return nil
	}
	if protocol_common.HistoryRecordPG == nil {
		return fmt.Errorf("history record pg db is nil")
	}

	writePgDeviceHistoryData := make([]models.DevicesPgHistoryData, 0, len(writeDeviceHistoryData))
	for _, historyData := range writeDeviceHistoryData {
		writePgDeviceHistoryData = append(writePgDeviceHistoryData, models.DevicesPgHistoryData{
			DataName:      historyData.DataName,
			DeviceUuid:    historyData.DeviceUuid,
			ProjectUuid:   historyData.ProjectUuid,
			DeviceName:    historyData.DeviceName,
			DataUuid:      historyData.DataUuid,
			ModelDataUuid: historyData.ModelDataUuid,
			RecordTime:    historyData.RecordTime,
			DataUnit:      historyData.DataUnit,
			DataValue:     historyData.DataValue,
		})
	}
	return insertHistoryPgData(writePgDeviceHistoryData, OnceWriteHistoryNumber)
}

func writeClickHouseHistoryData(writeDeviceHistoryData []models.DevicesHistoryDataList) error {
	if len(writeDeviceHistoryData) == 0 {
		return nil
	}
	if protocol_common.HistoryRecordClickHouseDb == nil {
		return fmt.Errorf("history record clickhouse db is nil")
	}

	writeCHDeviceHistoryData := make([]models.DevicesCHHistoryData, 0, len(writeDeviceHistoryData))
	for _, historyData := range writeDeviceHistoryData {
		writeCHDeviceHistoryData = append(writeCHDeviceHistoryData, models.DevicesCHHistoryData{
			DataName:      historyData.DataName,
			DeviceUuid:    historyData.DeviceUuid,
			ProjectUuid:   historyData.ProjectUuid,
			DeviceName:    historyData.DeviceName,
			DataUuid:      historyData.DataUuid,
			ModelDataUuid: historyData.ModelDataUuid,
			RecordTime:    historyData.RecordTime,
			DataUnit:      historyData.DataUnit,
			DataValue:     historyData.DataValue,
		})
	}
	return protocol_common.HistoryRecordClickHouseDb.Model(&models.DevicesCHHistoryData{}).CreateInBatches(&writeCHDeviceHistoryData, OnceWriteHistoryNumber).Error
}

func writeHistoryBatch(batch historyWriteBatch) error {
	if len(batch.data) == 0 {
		return nil
	}

	if protocol_common.HistoryRecordDbType == 2 {
		return writeTDengineHistoryData(batch.data)
	} else if protocol_common.HistoryRecordDbType == 4 {
		return writeInfluxHistoryData(batch.data)
	} else if protocol_common.HistoryRecordDbType == 5 {
		return writePGHistoryData(batch.data)
	} else if protocol_common.HistoryRecordDbType == 3 {
		return writeClickHouseHistoryData(batch.data)
	}
	return writeInsideHistoryData(batch.data)
}

func historyWriteWorker(jobCh <-chan historyWriteBatch) {
	for batch := range jobCh {
		if err := writeHistoryBatch(batch); err != nil {
			logs.Error(fmt.Sprintf("write history data batch failed: %v", err))
		} else {
			deleteHistoryBatchKeys(batch.keys)
		}
		if batch.waitGroup != nil {
			batch.waitGroup.Done()
		}
	}
}

// 根据日期动态创建表（如果表不存在）并插入订单数据
func insertHistoryData(writeDeviceHistoryData []models.DevicesHistoryDataList, OnceWriteHistoryNumber int) error {
	t := time.Now()
	var tableName string
	if protocol_common.HistoryPartitionType == 1 {
		tableName = fmt.Sprintf("devices_history_data_%s", t.Format("2006"))
	} else if protocol_common.HistoryPartitionType == 2 {
		tableName = fmt.Sprintf("devices_history_data_%s", t.Format("200601"))
	} else if protocol_common.HistoryPartitionType == 3 {
		tableName = fmt.Sprintf("devices_history_data_%s", t.Format("20060102"))
	} else if protocol_common.HistoryPartitionType == 4 {
		tableName = fmt.Sprintf("devices_history_data_%s", t.Format("20060102_15"))
	} else {
		tableName = "devices_history_data_list"
	}
	if protocol_common.HistoryPartitionType != 0 {
		// 要检查的表名
		// 使用 raw SQL 查询来检查表名是否存在
		var count int64
		models.Db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name =?", tableName).Scan(&count)
		if count <= 0 {
			logs.Info("创建历史数据表格", tableName)
			// 自动迁移，创建表
			if err := models.Db.Table(tableName).AutoMigrate(&models.DevicesHistoryDataList{}); err != nil {
				return fmt.Errorf("failed to auto migrate table: %v", err)
			}
		}
	}

	ee := models.Db.Table(tableName).CreateInBatches(&writeDeviceHistoryData, OnceWriteHistoryNumber).Error
	return ee
}
func insertHistoryPgData(writeDeviceHistoryData []models.DevicesPgHistoryData, OnceWriteHistoryNumber int) error {
	t := time.Now()
	var tableName string
	if protocol_common.HistoryPartitionType == 1 {
		tableName = fmt.Sprintf("devices_history_data_%s", t.Format("2006"))
	} else if protocol_common.HistoryPartitionType == 2 {
		tableName = fmt.Sprintf("devices_history_data_%s", t.Format("200601"))
	} else if protocol_common.HistoryPartitionType == 3 {
		tableName = fmt.Sprintf("devices_history_data_%s", t.Format("20060102"))
	} else if protocol_common.HistoryPartitionType == 4 {
		tableName = fmt.Sprintf("devices_history_data_%s", t.Format("20060102_15"))
	} else {
		tableName = "devices_history_data_list"
	}
	if protocol_common.HistoryPartitionType != 0 {
		// 要检查的表名
		// 使用 raw SQL 查询来检查表名是否存在
		var exists bool
		query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", tableName)
		err := protocol_common.HistoryRecordPG.Raw(query).Scan(&exists).Error
		if !exists || err != nil {
			logs.Info("创建历史数据表格", tableName)
			// 自动迁移，创建表
			if err := protocol_common.HistoryRecordPG.Table(tableName).AutoMigrate(&models.DevicesPgHistoryData{}); err != nil {
				return fmt.Errorf("failed to auto migrate table: %v", err)
			}
		}
	}

	ee := protocol_common.HistoryRecordPG.Table(tableName).CreateInBatches(&writeDeviceHistoryData, OnceWriteHistoryNumber).Error
	return ee
}
func DealWithHistoryData() {
	for {
		data, code := protocol_common.GHistoryDataQueue.QueuePull()
		if data == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		if code != -1 {
			var build strings.Builder
			HistoryData := data.(models.DevicesHistoryDataList)
			build.WriteString(HistoryData.DeviceUuid)
			build.WriteString(HistoryData.DataUuid)
			key := build.String()
			dataTemp, isExist := DeviceHistoryDataTemp[key]
			if !isExist {
				if HistoryData.RecordType == 2 {
					protocol_common.HistoryDataWrite(HistoryData)
				} else {
					DeviceHistoryDataTemp[key] = HistoryData
				}
			} else {
				if HistoryData.RecordType == 1 {
					if HistoryData.RecordInterval == 0 {
						HistoryData.RecordInterval = 1
					}
					if (HistoryData.RecordTime.Unix() - dataTemp.RecordTime.Unix()) >= int64(HistoryData.RecordInterval) {
						//models.Db.Model(&models.DevicesHistoryDataList{}).Create(&HistoryData)
						//protocol_common.GSaveHistoryDataQueue.QueuePush(HistoryData)
						protocol_common.HistoryDataWrite(HistoryData)
						DeviceHistoryDataTemp[key] = HistoryData
					}
				} else if HistoryData.RecordType == 0 {
					ChargeValue, err3 := strconv.ParseFloat(HistoryData.RecordDataCharge, 32)
					if err3 != nil {
						time.Sleep(time.Millisecond * 100)
						continue
					}
					currentValue, err := strconv.ParseFloat(HistoryData.DataValue, 32)
					if err != nil {
						time.Sleep(time.Millisecond * 100)
						continue
					}
					oldValue, err1 := strconv.ParseFloat(dataTemp.DataValue, 32)
					if err1 != nil {
						time.Sleep(time.Millisecond * 100)
						continue
					}
					if math.Abs(currentValue-oldValue) >= ChargeValue {
						//models.Db.Model(&models.DevicesHistoryDataList{}).Create(&HistoryData)
						//protocol_common.GSaveHistoryDataQueue.QueuePush(HistoryData)
						protocol_common.HistoryDataWrite(HistoryData)
						DeviceHistoryDataTemp[key] = HistoryData
					}
				} else if HistoryData.RecordType == 2 {
					protocol_common.HistoryDataWrite(HistoryData)
				} else if HistoryData.RecordType == 3 {
					ChargeValue, err3 := strconv.ParseFloat(HistoryData.RecordDataCharge, 32)
					if err3 != nil {
						return
					}
					currentValue, err := strconv.ParseFloat(HistoryData.DataValue, 32)
					if err != nil {
						return
					}
					oldValue, err1 := strconv.ParseFloat(dataTemp.DataValue, 32)
					if err1 != nil {
						return
					}
					if oldValue == 0 {
						DeviceHistoryDataTemp[key] = HistoryData
						return
					}
					DiffValue := math.Abs(currentValue - oldValue)
					cr := (DiffValue / oldValue) * 100
					if cr >= ChargeValue {
						protocol_common.HistoryDataWrite(HistoryData)
						DeviceHistoryDataTemp[key] = HistoryData
					}
				}
			}
		}

		time.Sleep(time.Millisecond * 1)
	}
}
func HistoryRecordDb() {
	HistoryConf, err := config.NewConfig("ini", "conf/historyData.conf")
	if err != nil {
		protocol_common.HistoryRecordDbType = 1
		return
	}
	protocol_common.HistoryRecordDbType, err = HistoryConf.Int("HistoryRecordDbType")
	if err != nil {
		protocol_common.HistoryRecordDbType = 1
		return
	}

	OnceWriteHistoryNumber, err = HistoryConf.Int("OnceWriteHistoryCounts")
	if err != nil {
		OnceWriteHistoryNumber = 2000
	}
	HistoryWriteWorkerCount, err = HistoryConf.Int("WriteWorkerCount")
	if err != nil || HistoryWriteWorkerCount <= 0 {
		HistoryWriteWorkerCount = defaultHistoryWriteWorkerCount()
	}
	protocol_common.HistoryPartitionType, err = HistoryConf.Int("partitiontype")
	if err != nil {
		protocol_common.HistoryPartitionType = 0
	}
	if protocol_common.HistoryRecordDbType == 2 {

		TDenginePort, _ := HistoryConf.String("TDengine::TDenginePort")
		TDengineHost, _ := HistoryConf.String("TDengine::TDengineHost")
		UserName, _ := HistoryConf.String("TDengine::UserName")
		PassWord, _ := HistoryConf.String("TDengine::PassWord")
		logs.Info("正在连接涛思数据,TDengine 数据库连接信息", TDengineHost, TDenginePort, UserName, PassWord)
		var taosDSN = UserName + ":" + PassWord + "@" + "http(" + TDengineHost + ":" + TDenginePort + ")/" //"root:taosdata@tcp(localhost:6030)/"
		taos, err := sql.Open("taosRestful", taosDSN)
		if err != nil {
			logs.Error("failed to connect TDengine")
			return
		}
		_, err = taos.Exec("CREATE DATABASE if not exists ISMHistoryDb")
		if err != nil {
			fmt.Println("failed to create database, err:", err)
		}
		_, err = taos.Exec("CREATE STABLE if not exists ISMHistoryDb.TempleteHistoryDatas (record_time TIMESTAMP, data_name NCHAR(255), device_uuid NCHAR(255), project_uuid NCHAR(255),device_name NCHAR(255),data_uuid NCHAR(255),model_data_uuid NCHAR(255),data_unit NCHAR(255),data_value NCHAR(255)) TAGS (groupId int);")
		if err != nil {
			fmt.Println("failed to create stable, err:", err)
		}
		logs.Info("正在连接涛思数据成功")
		protocol_common.HistoryRecordTsDb = taos
	} else if protocol_common.HistoryRecordDbType == 3 {

		sqlLogger := logger.New(
			Writer{},
			logger.Config{
				SlowThreshold:             1000 * time.Millisecond, // Slow SQL threshold
				LogLevel:                  logger.Error,            // Log level
				IgnoreRecordNotFoundError: true,                    // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,                    // Disable color
			},
		)

		dbConfig := gorm.Config{
			PrepareStmt:     false,
			Logger:          sqlLogger,
			CreateBatchSize: 100,
		}

		ChickHousePort, _ := HistoryConf.String("ChickHouse::ChickHousePort")
		ChickHouseHost, _ := HistoryConf.String("ChickHouse::ChickHouseHost")
		UserName, _ := HistoryConf.String("ChickHouse::UserName")
		DataBase, _ := HistoryConf.String("ChickHouse::DataBase")
		PassWord, _ := HistoryConf.String("ChickHouse::PassWord")
		ReadTimeout, _ := HistoryConf.String("ChickHouse::ReadTimeout")
		ConnectTimeout, _ := HistoryConf.String("ChickHouse::ConnectTimeout")
		logs.Info("正在连接ClickHouse 数据库连接信息", ChickHouseHost, ChickHousePort, UserName, PassWord, DataBase)
		dsn := "clickhouse://" + UserName + ":" + PassWord + "@" + ChickHouseHost + ":" + ChickHousePort + "/" + DataBase + "?dial_timeout=" + ConnectTimeout + "&read_timeout=" + ReadTimeout
		housedb, err := gorm.Open(clickhouse.Open(dsn), &dbConfig)
		if err != nil {
			logs.Error("ClickHouse 数据库连接失败 ", err)
			return
		} else {
			logs.Info("ClickHouse 数据库连接成功")
		}
		// Auto Migrate
		err = housedb.AutoMigrate(&models.DevicesCHHistoryData{})
		if err != nil {
			logs.Info(err)
		}
		protocol_common.HistoryRecordClickHouseDb = housedb

	} else if protocol_common.HistoryRecordDbType == 4 {
		Url, _ := HistoryConf.String("Influxdb::Url")
		Token, _ := HistoryConf.String("Influxdb::Token")

		Org, _ := HistoryConf.String("Influxdb::Org")
		Bucket, _ := HistoryConf.String("Influxdb::Bucket")

		// Create HTTP client
		httpClient := &http.Client{
			Timeout: time.Second * time.Duration(1200),
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: 1200 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 5 * time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		}
		logs.Info("正在连接Influxdb 数据库连接信息", Url, Token, Org, Bucket)
		influxdb2Default := influxdb2.DefaultOptions()
		influxdb2Default.SetHTTPClient(httpClient)
		influxdb2Default.SetBatchSize(50000)
		influxdb2Default.SetUseGZip(true)
		client := influxdb2.NewClientWithOptions(Url, Token, influxdb2Default)
		protocol_common.HistoryRecordInfluxdbDb = client.WriteAPIBlocking(Org, Bucket)
		protocol_common.HistoryRecordInfluxdbQuery = client.QueryAPI(Org)
		protocol_common.HistoryRecordInfluxdbBucket = Bucket
	} else if protocol_common.HistoryRecordDbType == 5 {
		postgresuser, _ := HistoryConf.String("PG::User")
		postgrespwd, _ := HistoryConf.String("PG::PassWord")
		postgreshost, _ := HistoryConf.String("PG::Host")
		postgresport, _ := HistoryConf.String("PG::Port")
		postgresdbname, _ := HistoryConf.String("PG::DbName")

		connstr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			postgreshost,
			postgresuser,
			postgrespwd,
			postgresdbname,
			postgresport,
		)
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
		logs.Info("正在连接PG 数据库连接信息", postgreshost, postgresport, postgresuser, postgresdbname)
		protocol_common.HistoryRecordPG, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  connstr,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &dbConfig)
		if err != nil {
			logs.Error("PG数据库连接失败", err)
			return
		}
		logs.Info("PG数据库连接成功")
		var exists bool
		query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = 'public')", "devices_history_data_list")
		protocol_common.HistoryRecordPG.Raw(query).Scan(&exists)
		if !exists {
			// 自动迁移，创建表
			if err := protocol_common.HistoryRecordPG.Table("devices_history_data_list").AutoMigrate(&models.DevicesPgHistoryData{}); err != nil {
				logs.Error("failed to auto migrate table: %v", err)
				return
			}
		}
	}
}
func TraversalHistoryDb() {
	var waitGroup sync.WaitGroup
	jobCh := make(chan historyWriteBatch, 1)
	done := make(chan struct{})
	go func() {
		historyWriteWorker(jobCh)
		close(done)
	}()
	dispatchCount := dispatchHistoryDbBatches(jobCh, &waitGroup)
	if dispatchCount > 0 {
		waitGroup.Wait()
	}
	close(jobCh)
	<-done
	/*
		var InsertSqlBuilder strings.Builder
		InsertSqlBuilder.Grow(10240)
		if protocol_common.GSaveHistoryDataLevelDb == nil {
		var writeDeviceHistoryData []models.DevicesHistoryDataList
		var writePgDeviceHistoryData []models.DevicesPgHistoryData
		var writeCHDeviceHistoryData []models.DevicesCHHistoryData
		var writeInfluxDeviceHistoryData []*write.Point
		var batchKeys [][]byte
		// 遍历数据库
		if protocol_common.GSaveHistoryDataLevelDb == nil {
			time.Sleep(time.Millisecond * 100)
			return
		}

		iter := protocol_common.GSaveHistoryDataLevelDb.NewIterator(nil, nil)

		rand.Seed(time.Now().UnixNano())
		InsertSqlBuilder.WriteString(fmt.Sprintf("INSERT INTO ISMHistoryDb.HistoryDatas USING ISMHistoryDb.TempleteHistoryDatas  TAGS(%d) VALUES", rand.Intn(0xFFFFFF)))

		var indexNo int = 0
		for iter.Next() {
			var HistoryData models.DevicesHistoryDataList
			var HistoryPgData models.DevicesPgHistoryData
			var HistoryCHData models.DevicesCHHistoryData
			key := iter.Key()
			err := json.Unmarshal(iter.Value(), &HistoryData)
			if err == nil {
				if protocol_common.HistoryRecordDbType == 2 {
					indexNo++
					InsertSqlBuilder.WriteString(fmt.Sprintf(" ('%s', '%s','%s','%s','%s','%s','%s','%s','%s') ", HistoryData.RecordTime.Format("2006-01-02 15:04:05.0000"), HistoryData.DataName, HistoryData.DeviceUuid, HistoryData.ProjectUuid, HistoryData.DeviceName, HistoryData.DataUuid, HistoryData.ModelDataUuid, HistoryData.DataUnit, HistoryData.DataValue))
					if indexNo >= OnceWriteHistoryNumber {
						_, err := protocol_common.HistoryRecordTsDb.Exec(InsertSqlBuilder.String())
						if err != nil {
							logs.Error(err)
						} else {
							protocol_common.GSaveHistoryDataLevelDb.Delete(key, nil)
						}
						InsertSqlBuilder.Reset()
						InsertSqlBuilder.WriteString(fmt.Sprintf("INSERT INTO ISMHistoryDb.HistoryDatas USING ISMHistoryDb.TempleteHistoryDatas  TAGS(%d) VALUES", rand.Intn(0xFFFFFF)))
						indexNo = 0
					}
				} else if protocol_common.HistoryRecordDbType == 4 {
					tags := map[string]string{
						"DataName":      HistoryData.DataName,
						"DeviceUuid":    HistoryData.DeviceUuid,
						"ProjectUuid":   HistoryData.ProjectUuid,
						"DeviceName":    HistoryData.DeviceName,
						"DataUuid":      HistoryData.DataUuid,
						"ModelDataUuid": HistoryData.ModelDataUuid,
						"DataUnit":      HistoryData.DataUnit,
					}
					fields := map[string]interface{}{
						"DataValue": HistoryData.DataValue,
					}
					point := write.NewPoint("ISMHistoryData", tags, fields, HistoryData.RecordTime)
					writeInfluxDeviceHistoryData = append(writeInfluxDeviceHistoryData, point)
					batchKeys = append(batchKeys, key)
					indexNo++
					if indexNo >= OnceWriteHistoryNumber {
						if err := protocol_common.HistoryRecordInfluxdbDb.WritePoint(context.Background(), writeInfluxDeviceHistoryData...); err != nil {
							logs.Error(err)
						} else {
							// 批量删除已写入的数据
							for _, k := range batchKeys {
								protocol_common.GSaveHistoryDataLevelDb.Delete(k, nil)
							}
						}
						writeInfluxDeviceHistoryData = []*write.Point{}
						batchKeys = [][]byte{}
						indexNo = 0
					}
				} else if protocol_common.HistoryRecordDbType == 5 {
					HistoryPgData.DataName = HistoryData.DataName
					HistoryPgData.DeviceUuid = HistoryData.DeviceUuid
					HistoryPgData.ProjectUuid = HistoryData.ProjectUuid
					HistoryPgData.DeviceName = HistoryData.DeviceName
					HistoryPgData.DataUuid = HistoryData.DataUuid
					HistoryPgData.ModelDataUuid = HistoryData.ModelDataUuid
					HistoryPgData.RecordTime = HistoryData.RecordTime
					HistoryPgData.DataUnit = HistoryData.DataUnit
					HistoryPgData.DataValue = HistoryData.DataValue
					writePgDeviceHistoryData = append(writePgDeviceHistoryData, HistoryPgData)
					protocol_common.GSaveHistoryDataLevelDb.Delete(key, nil)
				} else if protocol_common.HistoryRecordDbType == 3 {
					HistoryCHData.DataName = HistoryData.DataName
					HistoryCHData.DeviceUuid = HistoryData.DeviceUuid
					HistoryCHData.ProjectUuid = HistoryData.ProjectUuid
					HistoryCHData.DeviceName = HistoryData.DeviceName
					HistoryCHData.DataUuid = HistoryData.DataUuid
					HistoryCHData.ModelDataUuid = HistoryData.ModelDataUuid
					HistoryCHData.RecordTime = HistoryData.RecordTime
					HistoryCHData.DataUnit = HistoryData.DataUnit
					HistoryCHData.DataValue = HistoryData.DataValue
					writeCHDeviceHistoryData = append(writeCHDeviceHistoryData, HistoryCHData)
					protocol_common.GSaveHistoryDataLevelDb.Delete(key, nil)
				} else {
					writeDeviceHistoryData = append(writeDeviceHistoryData, HistoryData)
					protocol_common.GSaveHistoryDataLevelDb.Delete(key, nil)
				}
			}
		}
		if protocol_common.HistoryRecordDbType == 2 {
			if indexNo >= 1 && indexNo < OnceWriteHistoryNumber {
				_, err := protocol_common.HistoryRecordTsDb.Exec(InsertSqlBuilder.String())
				if err != nil {
					logs.Error(err)
				}
			}
			InsertSqlBuilder.Reset()
			indexNo = 0
		} else if protocol_common.HistoryRecordDbType == 1 && len(writeDeviceHistoryData) > 0 {
			if protocol_common.InsideDbType == 0 {
				insertHistoryData(writeDeviceHistoryData, OnceWriteHistoryNumber)
			} else {
				//sqlite3数据不使用分表
				models.Db.Model(&models.DevicesHistoryDataList{}).CreateInBatches(&writeDeviceHistoryData, OnceWriteHistoryNumber)
			}
		} else if protocol_common.HistoryRecordDbType == 5 && len(writePgDeviceHistoryData) > 0 {
			insertHistoryPgData(writePgDeviceHistoryData, OnceWriteHistoryNumber)
			// protocol_common.HistoryRecordPG.Model(&models.DevicesPgHistoryData{}).CreateInBatches(&writePgDeviceHistoryData, OnceWriteHistoryNumber)
		} else if protocol_common.HistoryRecordDbType == 4 && len(writeInfluxDeviceHistoryData) > 0 {
			if err := protocol_common.HistoryRecordInfluxdbDb.WritePoint(context.Background(), writeInfluxDeviceHistoryData...); err != nil {
				logs.Error(err)
			} else {
				// 批量删除剩余的已写入数据
				for _, k := range batchKeys {
					protocol_common.GSaveHistoryDataLevelDb.Delete(k, nil)
				}
			}
			writeInfluxDeviceHistoryData = []*write.Point{}
			batchKeys = [][]byte{}
			indexNo = 0
		} else if protocol_common.HistoryRecordDbType == 3 && len(writeCHDeviceHistoryData) > 0 {
			if protocol_common.HistoryRecordClickHouseDb != nil {
				protocol_common.HistoryRecordClickHouseDb.Model(&models.DevicesCHHistoryData{}).CreateInBatches(&writeCHDeviceHistoryData, OnceWriteHistoryNumber)
			}
		}
		InsertSqlBuilder.Reset()
		iter.Release()
	*/
}

func dispatchHistoryDbBatches(jobCh chan<- historyWriteBatch, waitGroup *sync.WaitGroup) int {
	if protocol_common.GSaveHistoryDataLevelDb == nil {
		return 0
	}

	batchSize := OnceWriteHistoryNumber
	if batchSize <= 0 {
		batchSize = 100
	}

	iter := protocol_common.GSaveHistoryDataLevelDb.NewIterator(nil, nil)
	defer iter.Release()

	batch := newHistoryWriteBatch(batchSize, waitGroup)
	dispatchCount := 0
	for iter.Next() {
		var historyData models.DevicesHistoryDataList
		// 使用 gob 反序列化，替代 JSON 反序列化
		var buf bytes.Buffer
		buf.Write(iter.Value())
		err := gob.NewDecoder(&buf).Decode(&historyData)
		if err != nil {
			logs.Error(fmt.Sprintf("unmarshal history data failed: %v", err))
			continue
		}

		batch.data = append(batch.data, historyData)
		batch.keys = append(batch.keys, cloneHistoryLevelDbKey(iter.Key()))
		if len(batch.data) >= batchSize {
			if waitGroup != nil {
				waitGroup.Add(1)
			}
			jobCh <- batch
			dispatchCount++
			batch = newHistoryWriteBatch(batchSize, waitGroup)
		}
	}

	if err := iter.Error(); err != nil {
		logs.Error(fmt.Sprintf("iterate history leveldb failed: %v", err))
	}

	if len(batch.data) > 0 {
		if waitGroup != nil {
			waitGroup.Add(1)
		}
		jobCh <- batch
		dispatchCount++
	}

	return dispatchCount
}

func DealWithSaveHistoryData() {
	workerCount := HistoryWriteWorkerCount
	if workerCount <= 0 {
		workerCount = defaultHistoryWriteWorkerCount()
	}
	// 从配置文件读取线程数
	configWorkerCount, err := config.Int("HistoryWriteWorkerCount")
	if err == nil && configWorkerCount > 0 {
		workerCount = configWorkerCount
	}

	jobCh := make(chan historyWriteBatch, workerCount*4)
	for i := 0; i < workerCount; i++ {
		go historyWriteWorker(jobCh)
	}

	for {
		var waitGroup sync.WaitGroup
		dispatchCount := dispatchHistoryDbBatches(jobCh, &waitGroup)
		if dispatchCount <= 0 {
			time.Sleep(time.Millisecond * 50)
			continue
		}
		waitGroup.Wait()
	}
}
