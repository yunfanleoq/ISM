/**
 * 定时历史存档：到点从实时库取最新值写入历史（不依赖采集成功事件）。
 * RecordType==1（定时存储）由本任务独占；采集路径的 HistoryDataWrite 会忽略该类型。
 */
package alarmTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"sync"
	"sync/atomic"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type timedHistoryPoint struct {
	DataUuid       string
	DeviceUuid     string
	ProjectUuid    string
	DeviceName     string
	DataName       string
	ModelDataUuid  string
	DataUnit       string
	RecordInterval int // 秒
}

var (
	timedPointsMu     sync.RWMutex
	timedPoints       []timedHistoryPoint
	timedLastCycle    sync.Map // deviceUuid+dataUuid -> last archived cycle unix
	timedPointsLoaded int32
)

// reloadTimedHistoryPoints 从业务库加载「开启历史 + 定时存储」点位元数据（不含实时值）。
func reloadTimedHistoryPoints() {
	var rows []models.DeviceRealData
	err := models.Db.Model(&models.DeviceRealData{}).
		Select("uuid, device_uuid, project_uuid, device_name, name, model_data_uuid, data_unit, record_interval, record_type").
		Where("is_record = ? AND record_type = ?", 1, 1).
		Find(&rows).Error
	if err != nil {
		logs.Error("reload timed history points failed: %v", err)
		return
	}
	next := make([]timedHistoryPoint, 0, len(rows))
	for _, row := range rows {
		interval := row.RecordInterval
		if interval <= 0 {
			interval = 1
		}
		next = append(next, timedHistoryPoint{
			DataUuid:       row.Uuid,
			DeviceUuid:     row.DeviceUuid,
			ProjectUuid:    row.ProjectUuid,
			DeviceName:     row.DeviceName,
			DataName:       row.Name,
			ModelDataUuid:  row.ModelDataUuid,
			DataUnit:       row.DataUnit,
			RecordInterval: interval,
		})
	}
	timedPointsMu.Lock()
	timedPoints = next
	timedPointsMu.Unlock()
	atomic.StoreInt32(&timedPointsLoaded, 1)
	logs.Error("timed history snapshot points loaded: %d", len(next))
}

func snapshotDueTimedPoints(now time.Time) (wrote, skippedNoValue, skippedNotDue int) {
	timedPointsMu.RLock()
	points := timedPoints
	timedPointsMu.RUnlock()
	if len(points) == 0 {
		return
	}

	for _, p := range points {
		interval := p.RecordInterval
		if interval <= 0 {
			interval = 1
		}
		cycleDur := time.Duration(interval) * time.Second
		cycleTime := now.Truncate(cycleDur)
		cycleUnix := cycleTime.Unix()
		key := p.DeviceUuid + p.DataUuid

		if last, ok := timedLastCycle.Load(key); ok {
			if lastUnix, ok2 := last.(int64); ok2 && lastUnix >= cycleUnix {
				skippedNotDue++
				continue
			}
		}

		val, ok := protocol_common.LoadDeviceRealValue(p.DataUuid, p.DeviceName, p.DataName)
		if !ok {
			skippedNoValue++
			protocol_common.ErrorThrottled(
				"history:snapshot:novalue:"+key,
				"timed history snapshot: no realtime value yet, skip cycle device=%s data=%s uuid=%s cycle=%s",
				p.DeviceName, p.DataName, p.DataUuid, cycleTime.Format("2006-01-02 15:04:05"),
			)
			continue
		}

		sample := models.DevicesHistoryDataList{
			DataName:         p.DataName,
			DeviceUuid:       p.DeviceUuid,
			ProjectUuid:      p.ProjectUuid,
			DeviceName:       p.DeviceName,
			DataUuid:         p.DataUuid,
			ModelDataUuid:    p.ModelDataUuid,
			RecordTime:       cycleTime,
			DataUnit:         p.DataUnit,
			DataValue:        val,
			RecordInterval:   interval,
			RecordType:       1,
			RecordDataCharge: "",
			RecordDataTimely: "",
		}
		protocol_common.HistoryDataWriteSnapshot(sample)
		timedLastCycle.Store(key, cycleUnix)
		wrote++
	}
	return
}

// DealWithTimedRealtimeHistorySnapshot 定时从实时库快照写历史。
// 到点即取最新实时值；采集周期慢或某轮采失败不影响——只要实时库里有值就会落档。
func DealWithTimedRealtimeHistorySnapshot() {
	// 启动稍晚，等实时值预热/采集线程起来
	time.Sleep(15 * time.Second)
	reloadTimedHistoryPoints()

	reloadTicker := time.NewTicker(5 * time.Minute)
	defer reloadTicker.Stop()
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	var lastLogUnix int64
	for {
		select {
		case <-reloadTicker.C:
			reloadTimedHistoryPoints()
		case now := <-tick.C:
			wrote, noVal, notDue := snapshotDueTimedPoints(now)
			if wrote > 0 || noVal > 0 {
				u := now.Unix()
				if u-lastLogUnix >= 60 {
					lastLogUnix = u
					logs.Error(
						"timed history snapshot tick: wrote=%d noRealtimeValue=%d notDue=%d points=%d",
						wrote, noVal, notDue, len(timedPoints),
					)
				}
			}
		}
	}
}
