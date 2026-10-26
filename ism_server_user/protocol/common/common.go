/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-05 15:34:25
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package protocolCommon

import (
	"ISMServer/middleware"
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gorilla/websocket"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"gorm.io/gorm"
)

const (
	MAXTRIGGERITEM = 100000 //支持最大的触发器数量
)

var HistoryRecordDbType int
var HistoryRecordTsDb *sql.DB
var HistoryRecordInfluxdbDb api.WriteAPIBlocking
var HistoryRecordInfluxdbQuery api.QueryAPI
var HistoryRecordInfluxdbBucket string
var HistoryRecordClickHouseDb *gorm.DB
var InsideDbType int
var ClearAlarmType int
var HistoryKeepDays int = 7 // 默认保留7天历史测点数据
// AlarmKeepDays 已消除告警硬删保留天数；0 表示不自动硬删告警。默认 90。
var AlarmKeepDays int = 90
// ActiveAlarmClearThreshold：clear_time 严格小于该值视为「未消除」实时告警（协议写入哨兵为 2006-01-02）。
// 自动硬删只允许 clear_time >= 该阈值（已人工/恢复消除）且早于保留截止日。
const ActiveAlarmClearThreshold = "2007-01-02 15:04:05"

var HistoryRecordPG *gorm.DB

var GGatherDataQueue = (*middleware.Queue)(nil)

var GAlarmQueue = (*middleware.Queue)(nil)
var GHistoryDataQueue = (*middleware.Queue)(nil)
var GRealDataQueue = (*middleware.Queue)(nil)

var GTriggerDataQueue sync.Map

var PushGAlarmQueue = (*middleware.Queue)(nil)

var GGatherSystemDataQueue = (*middleware.Queue)(nil)

var GCustomDataQueue sync.Map

var GSaveHistoryDataQueue = (*middleware.Queue)(nil)
var GSaveHistoryDataLevelDb *leveldb.DB

// 历史数据批量写入相关变量
var (
	historyDataBuffer      []interface{}
	historyDataBufferSize  int
	historyDataMutex       sync.Mutex
	historyDataFlushTicker *time.Ticker
)

var RecordPath string

var ModbusDebug bool = false
var IEC104Debug bool = false
var IsLicense = true
var IsOem = true
var ISMProtectedID = ""
var ConfigPageCount int = 100000
var ConfigAppCount int = 100000
var IsAuthLimit bool = false
var IsAuthTimeLimit bool = false
var AuthorizationDays int = 30
var AuthorizationCreateDate string = ""
var AuthRemainingTimeDays int = 0
var AuthRemainingTimeHours int = 0
var AuthRemainingTimeMinutes int = 0

var IsAuthGenuine bool = false
var WsRWMux *sync.Mutex

type UpdateStu struct {
	Uuid          string
	Value         string
	Oid           string
	ModelDataUuid string
	DataName      string
}

type PushRealDataWebData struct {
	DeviceUuid  string
	ProjectUuid string
	Cmd         string
	DeviceName  string
	Data        []UpdateStu
}
type PushSystemDataWebData struct {
	ProjectUuid string
	Cmd         string
	Data        []UpdateStu
}

type PushRealAlarmWebData struct {
	DeviceUuid string
	Cmd        string
	AlarmID    string
	Value      int
}

type PushAlarm struct {
	ID                uint
	DeviceName        string
	DataName          string
	DeviceUuid        string
	ProjectUuid       string
	DataUuid          string
	ModelDataUuid     string
	Value             string
	Cmd               string
	RealValue         string
	HappenTime        time.Time
	AlarmLevel        int
	AlarmMessage      string
	AlarmClearMessage string
	AlarmOnValue      int // 0 或 1：实时值等于该字符串时触发告警，默认 1
	SuppressNotice    bool
}

// ResolveAlarmOnValue 未配置时默认 1 报警
func ResolveAlarmOnValue(v int) int {
	if v == 0 || v == 1 {
		return v
	}
	return 1
}

// IsAlarmValueActive 判断当前值是否处于告警态
func IsAlarmValueActive(value string, alarmOnValue int) bool {
	trigger := ResolveAlarmOnValue(alarmOnValue)
	if trigger == 0 {
		return value == "0"
	}
	return value == "1"
}

type TriggerRealData struct {
	DeviceName        string
	DataName          string
	DeviceUuid        string
	ProjectUuid       string
	DataUuid          string
	ModelDataUuid     string
	Value             string
	DataType          int
	AlarmShield       int
	GatherTime        time.Time
	AlarmLevel        int
	IsAlarm           int
	AlarmMessage      string
	DataUnit          string
	AlarmClearMessage string
	IsRecord          int
	RecordInterval    int
	RecordType        int
	RecordDataCharge  string
}
type ISMNodeProjectConnStu struct {
	NodeName    string
	ProjectUuid string
	Uuid        string
	ChanelConn  *websocket.Conn
	ConnRwMutex *sync.Mutex
}

type ISMNodePushDataStu struct {
	ProjectUuid string
	MsgType     int
	Message     interface{}
}

var IsRestoreDb int = 0
var DeviceAlarmTriggerMap sync.Map
var DeviceCustomDataMap sync.Map
var DeviceRealDataMap sync.Map
var DeviceRealDataMapByUUID sync.Map
var ISMNodeDeviceRealDataMapByUUID sync.Map
var ISMSystemDataMapByUUID sync.Map
var DeviceProjectRealDataMap sync.Map
var NetworkNodePushDataChanel chan interface{}
var SendDataDelay int
var AlarmCacheCount int
var RealDataChanelCache int
var HistoryCacheCount int
var DTUConnCache int
var HistoryPartitionType int = 0
var DataFloatNumberFunc map[string]govaluate.ExpressionFunction

// RealDataFrontendPush 由 websocket 包注册，供 models 等避免循环依赖地推送实时值到前端。
var RealDataFrontendPush func(msg PushRealDataWebData)

// NotifyRealDataFrontend 推送 RealData 到前端（经合并窗）；未注册时为 no-op。
func NotifyRealDataFrontend(msg PushRealDataWebData) {
	if RealDataFrontendPush != nil {
		RealDataFrontendPush(msg)
	}
}

var ISMNodeProjectConn sync.Map
var ISMNodeServerProjectConn sync.Map
var ISMNodeProjectDataCollects sync.Map

var DeviceProjectAllRealDataMap sync.Map

// initHistoryDataBuffer 初始化历史数据缓冲区
func initHistoryDataBuffer(bufferSize int, flushIntervalMs int) {
	historyDataBufferSize = bufferSize
	historyDataBuffer = make([]interface{}, 0, bufferSize)

	// 启动定期刷新定时器
	historyDataFlushTicker = time.NewTicker(time.Duration(flushIntervalMs) * time.Millisecond)
	go func() {
		for range historyDataFlushTicker.C {
			flushHistoryDataBuffer()
		}
	}()
}

// flushHistoryDataBuffer 刷新历史数据缓冲区，批量写入LevelDB
func flushHistoryDataBuffer() {
	historyDataMutex.Lock()
	defer historyDataMutex.Unlock()

	if len(historyDataBuffer) == 0 {
		return
	}

	if GSaveHistoryDataLevelDb != nil {
		// 创建批量写入对象
		wb := new(leveldb.Batch)

		// 遍历缓冲区数据，添加到批量写入中
		// 遍历缓冲区数据，添加到批量写入中
		for _, data := range historyDataBuffer {
			// 使用 gob 序列化，比 JSON 更高效
			var buf bytes.Buffer
			err := gob.NewEncoder(&buf).Encode(data)
			if err == nil {
				// 使用时间戳+随机数作为键
				timestamp := time.Now().UnixNano()
				random := rand.Intn(1000000) // 增加随机数范围，减少冲突概率
				key := fmt.Sprintf("%d_%d", timestamp, random)
				wb.Put([]byte(key), buf.Bytes())
			}
		}

		// 执行批量写入
		if err := GSaveHistoryDataLevelDb.Write(wb, nil); err != nil {
			// 写入失败，记录错误
			logs.Error("批量写入历史数据失败: ", err)
		} else {
			// 写入成功，清空缓冲区
			historyDataBuffer = make([]interface{}, 0, historyDataBufferSize)
		}
	}
}

// HistoryDataWrite 写入历史数据（批量写入版本）
// 数据会先进入内存缓冲区，达到一定数量或时间后批量写入LevelDB。
// RecordType==1（定时存储）由实时库快照任务独占，采集路径调用本函数会被忽略。
func HistoryDataWrite(HistoryData any) {
	if isTimedRecordOwnedBySnapshot(HistoryData) {
		return
	}
	historyDataWriteInternal(HistoryData)
}

// HistoryDataWriteSnapshot 仅供定时快照任务写入 RecordType==1。
func HistoryDataWriteSnapshot(HistoryData any) {
	historyDataWriteInternal(HistoryData)
}

func isTimedRecordOwnedBySnapshot(sample interface{}) bool {
	if sample == nil {
		return false
	}
	v := reflect.ValueOf(sample)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}
	f := v.FieldByName("RecordType")
	if !f.IsValid() {
		return false
	}
	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.Int() == 1
	}
	return false
}

func historyDataWriteInternal(HistoryData any) {
	historyDataMutex.Lock()
	defer historyDataMutex.Unlock()

	historyDataBuffer = append(historyDataBuffer, HistoryData)
	if len(historyDataBuffer) >= historyDataBufferSize {
		go flushHistoryDataBuffer()
	}
}

// HistoryDataFlush 手动刷新历史数据缓冲区
func HistoryDataFlush() {
	flushHistoryDataBuffer()
}
func ProtocolCommonInit() {
	var err error
	GGatherDataQueue = middleware.NewQueue(0)
	GAlarmQueue = middleware.NewQueue(0)
	PushGAlarmQueue = middleware.NewQueue(0)
	GHistoryDataQueue = middleware.NewQueue(1000000)
	GRealDataQueue = middleware.NewQueue(0)
	GGatherSystemDataQueue = middleware.NewQueue(0)
	SendDataDelay, err = config.Int("sendalarmspeed")
	if err != nil {
		SendDataDelay = 100
	}
	SendDataDelay, err = config.Int("sendalarmspeed")
	if err != nil {
		SendDataDelay = 100
	}
	//每个设备的告警数量缓存
	AlarmCacheCount, err = config.Int("AlarmCacheCount")
	if err != nil {
		AlarmCacheCount = 500
	}

	//每个设备的历史数据对比数量缓存
	HistoryCacheCount, err = config.Int("HistoryCacheCount")
	if err != nil {
		HistoryCacheCount = 500
	}

	//每个DTU连接的数量缓存
	DTUConnCache, err = config.Int("DTUConnCache")
	if err != nil {
		DTUConnCache = 500
	}

	GSaveHistoryDataQueueLength, err1 := config.Int("HistoryDataQueueLength")
	if err1 != nil {
		GSaveHistoryDataQueueLength = 100000 //默认1百万条
	}

	GSaveHistoryDataQueue = middleware.NewQueue(GSaveHistoryDataQueueLength)

	// 初始化历史数据批量写入缓冲区（正式环境默认更大批量、更长间隔，降低刷盘）
	historyDataBufferSize, err := config.Int("HistoryDataBufferSize")
	if err != nil {
		historyDataBufferSize = 10000
	}
	historyDataFlushInterval, err := config.Int("HistoryDataFlushInterval")
	if err != nil {
		historyDataFlushInterval = 2000
	}
	initHistoryDataBuffer(historyDataBufferSize, historyDataFlushInterval)

	RealDataChanelCache, err = config.Int("RealDataChanelCache")
	if err != nil || RealDataChanelCache <= 0 {
		// 默认加大缓冲，降低采集风暴时 RealDataChanel full 丢包概率
		RealDataChanelCache = 10000
	}
	NetworkNodePushDataChanel = make(chan interface{}, RealDataChanelCache)

	opts := &opt.Options{
		BlockCacheCapacity: 128 * 1024 * 1024, // 128MB 缓存
		WriteBuffer:        64 * 1024 * 1024,  // 64MB 写入缓冲区
		Compression:        opt.DefaultCompression,
		Filter:             filter.NewBloomFilter(10),
	}
	// LevelDB 作为历史数据写入缓冲区，在 macOS 上 flock() 可能因锁残留返回 EAGAIN
	// 通过重试机制处理，失败后降级为仅内存缓冲（不持久化），不终止服务
	GSaveHistoryDataLevelDb, err = leveldb.OpenFile("data/historyData.db", opts)
	if err != nil {
		logs.Warning("打开历史数据缓冲区失败（第1次）: %s，将在1秒后重试...", err)
		time.Sleep(1 * time.Second)
		GSaveHistoryDataLevelDb, err = leveldb.OpenFile("data/historyData.db", opts)
	}
	if err != nil {
		logs.Warning("打开历史数据缓冲区失败（第2次）: %s，将在3秒后重试...", err)
		time.Sleep(3 * time.Second)
		GSaveHistoryDataLevelDb, err = leveldb.OpenFile("data/historyData.db", opts)
	}
	if err != nil {
		logs.Warning("打开历史数据缓冲区失败（第3次）: %s，尝试修复...", err)
		os.RemoveAll("data/historyData.db")
		GSaveHistoryDataLevelDb, err = leveldb.OpenFile("data/historyData.db", opts)
	}
	if err != nil {
		logs.Error("LevelDB 历史数据缓冲区最终打开失败: %s", err)
		logs.Warning("降级运行：历史数据将直接写入数据库，不使用 LevelDB 本地缓冲")
		GSaveHistoryDataLevelDb = nil
	}

	DataFloatNumberFunc = map[string]govaluate.ExpressionFunction{
		"Round": func(args ...interface{}) (interface{}, error) {
			// 示例：返回参数的和
			var dataString string
			numberFloat, ok := args[1].(float64)
			if !ok {
				return nil, fmt.Errorf("unsupported type %T", numberFloat)
			}
			number := int(numberFloat)
			switch number {
			case 1:
				dataString = fmt.Sprintf("%.1f", args[0])
			case 2:
				dataString = fmt.Sprintf("%.2f", args[0])
			case 3:
				dataString = fmt.Sprintf("%.3f", args[0])
			case 4:
				dataString = fmt.Sprintf("%.4f", args[0])
			case 5:
				dataString = fmt.Sprintf("%.5f", args[0])
			case 6:
				dataString = fmt.Sprintf("%.6f", args[0])
			case 7:
				dataString = fmt.Sprintf("%.7f", args[0])
			case 8:
				dataString = fmt.Sprintf("%.8f", args[0])
			default:
				dataString = fmt.Sprintf("%.2f", args[0])
			}
			ff, err := strconv.ParseFloat(dataString, 64)
			if err != nil {
				return nil, err
			}
			return ff, nil
		},
	}

}
