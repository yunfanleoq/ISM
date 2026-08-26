/**
 * 定时历史存档：到点从实时库取最新值写入历史（不依赖采集成功事件）。
 * RecordType==1（定时存储）由本任务独占；采集路径的 HistoryDataWrite 会忽略该类型。
 *
 * 取值顺序：内存 cache（UUID / deviceName->pointName）→ 回退读 device_real_data.value
 * → 仍无值则复用上一档已写入值。禁止按点 Error 刷屏；写入走异步缓冲，不堵 1s ticker。
 */
package alarmTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type timedHistoryPoint struct {
	DataUuid         string
	DeviceUuid       string
	ProjectUuid      string
	DeviceName       string
	DataName         string
	ModelDataUuid    string
	DataUnit         string
	RecordInterval   int // 秒
	RecordType       int
	RecordDataTimely string
}

func hourlyIntervalSeconds(timely string, recordInterval int) int {
	switch strings.TrimSpace(timely) {
	case "1":
		return 5 * 60
	case "2":
		return 10 * 60
	case "3":
		return 15 * 60
	case "4":
		return 30 * 60
	case "5":
		return 60 * 60
	}
	if recordInterval > 0 {
		return recordInterval
	}
	return 3600
}

const timedDbFallbackChunk = 500

var (
	timedPointsMu     sync.RWMutex
	timedPoints       []timedHistoryPoint
	timedLastCycle    sync.Map // deviceUuid+dataUuid -> last archived cycle unix
	timedLastValue    sync.Map // dataUuid -> last written value
	timedPointsLoaded int32
)

func reloadTimedHistoryPoints() {
	var rows []models.DeviceRealData
	err := models.Db.Model(&models.DeviceRealData{}).
		Select("uuid, device_uuid, project_uuid, device_name, name, model_data_uuid, data_unit, record_interval, record_type, record_data_timely").
		Where("is_record = ? AND record_type IN ?", 1, []int{1, 4}).
		Find(&rows).Error
	if err != nil {
		logs.Error("reload timed history points failed: %v", err)
		return
	}
	timed1, timed4 := 0, 0
	next := make([]timedHistoryPoint, 0, len(rows))
	for _, row := range rows {
		if row.RecordType == 4 {
			timed4++
		} else {
			timed1++
		}
		interval := row.RecordInterval
		if row.RecordType == 4 {
			interval = hourlyIntervalSeconds(row.RecordDataTimely, row.RecordInterval)
		} else if interval <= 0 {
			interval = 1
		}
		next = append(next, timedHistoryPoint{
			DataUuid:         row.Uuid,
			DeviceUuid:       row.DeviceUuid,
			ProjectUuid:      row.ProjectUuid,
			DeviceName:       row.DeviceName,
			DataName:         row.Name,
			ModelDataUuid:    row.ModelDataUuid,
			DataUnit:         row.DataUnit,
			RecordInterval:   interval,
			RecordType:       row.RecordType,
			RecordDataTimely: row.RecordDataTimely,
		})
	}
	timedPointsMu.Lock()
	timedPoints = next
	timedPointsMu.Unlock()
	atomic.StoreInt32(&timedPointsLoaded, 1)
	logs.Info("timed history snapshot points loaded: total=%d type1=%d type4=%d", len(next), timed1, timed4)
}

func usableRealtimeValue(val string, ok bool) (string, bool) {
	if !ok {
		return "", false
	}
	val = strings.TrimSpace(val)
	if val == "" {
		return "", false
	}
	return val, true
}

func loadTimedPointValuesFromDB(uuids []string) map[string]string {
	out := make(map[string]string, len(uuids))
	if len(uuids) == 0 {
		return out
	}
	for start := 0; start < len(uuids); start += timedDbFallbackChunk {
		end := start + timedDbFallbackChunk
		if end > len(uuids) {
			end = len(uuids)
		}
		chunk := uuids[start:end]
		var rows []models.DeviceRealData
		err := models.Db.Model(&models.DeviceRealData{}).
			Select("uuid, value").
			Where("uuid IN ?", chunk).
			Find(&rows).Error
		if err != nil {
			logs.Warn("timed history snapshot: load device_real_data.value failed: %v", err)
			continue
		}
		for _, row := range rows {
			if v, ok := usableRealtimeValue(row.Value, true); ok {
				out[row.Uuid] = v
			}
		}
	}
	return out
}

func writeTimedSnapshot(p timedHistoryPoint, val string, interval int, cycleTime time.Time, cycleUnix int64, key string) {
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
		RecordType:       p.RecordType,
		RecordDataCharge: "",
		RecordDataTimely: p.RecordDataTimely,
	}
	protocol_common.HistoryDataWriteSnapshot(sample)
	timedLastCycle.Store(key, cycleUnix)
	if p.DataUuid != "" {
		timedLastValue.Store(p.DataUuid, val)
	}
}

func snapshotDueTimedPoints(now time.Time) (wrote, fromDB, reusedLast, skippedNoValue, skippedNotDue int) {
	timedPointsMu.RLock()
	points := timedPoints
	timedPointsMu.RUnlock()
	if len(points) == 0 {
		return
	}

	type dueItem struct {
		p         timedHistoryPoint
		interval  int
		cycleTime time.Time
		cycleUnix int64
		key       string
		val       string
		memOK     bool
	}

	due := make([]dueItem, 0, len(points))
	missUuids := make([]string, 0)
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

		item := dueItem{p: p, interval: interval, cycleTime: cycleTime, cycleUnix: cycleUnix, key: key}
		if val, ok := usableRealtimeValue(protocol_common.LoadDeviceRealValue(p.DataUuid, p.DeviceName, p.DataName)); ok {
			item.val = val
			item.memOK = true
		} else if p.DataUuid != "" {
			missUuids = append(missUuids, p.DataUuid)
		}
		due = append(due, item)
	}

	dbVals := loadTimedPointValuesFromDB(missUuids)
	for _, item := range due {
		val := item.val
		ok := item.memOK
		src := "mem"
		if !ok {
			dbRaw, dbHit := dbVals[item.p.DataUuid]
			if dbVal, dbOK := usableRealtimeValue(dbRaw, dbHit); dbOK {
				val = dbVal
				ok = true
				src = "db"
				protocol_common.StoreDeviceRealValue(item.p.DataUuid, item.p.DeviceName, item.p.DataName, val)
				fromDB++
			}
		}
		if !ok {
			if last, hit := timedLastValue.Load(item.p.DataUuid); hit {
				if lastVal, lastOK := usableRealtimeValue(last.(string), true); lastOK {
					val = lastVal
					ok = true
					src = "last"
					reusedLast++
				}
			}
		}
		if !ok {
			skippedNoValue++
			continue
		}
		writeTimedSnapshot(item.p, val, item.interval, item.cycleTime, item.cycleUnix, item.key)
		wrote++
		_ = src
	}
	return
}

// DealWithTimedRealtimeHistorySnapshot 定时从实时库快照写历史。
func DealWithTimedRealtimeHistorySnapshot() {
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
			wrote, fromDB, reused, noVal, notDue := snapshotDueTimedPoints(now)
			if wrote > 0 || noVal > 0 {
				u := now.Unix()
				if u-lastLogUnix >= 60 {
					lastLogUnix = u
					logs.Error(
						"timed history snapshot tick: wrote=%d fromDB=%d reusedLast=%d noRealtimeValue=%d notDue=%d points=%d",
						wrote, fromDB, reused, noVal, notDue, len(timedPoints),
					)
				}
			}
		}
	}
}
