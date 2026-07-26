package models

import (
	"ISMServer/utils/errmsg"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/beego/beego/v2/core/logs"
	"github.com/go-basic/uuid"
)

const (
	modbusBulkModelBatchSize   = 100
	modbusBulkRealDataBatch    = 200
	modbusBulkRealDataFlushAt  = 5000
	modbusBulkPrefetchChunk    = 500
	modbusBulkUpdateWorkers    = 8
)

// ModbusBulkUpsertResult 全量点位导入统计
type ModbusBulkUpsertResult struct {
	Added   int
	Updated int
	Skipped int
}

func modbusAddrKey(muid, groupUuid string, addr int) string {
	return fmt.Sprintf("%s\x00%s\x00%d", muid, groupUuid, addr)
}

func modbusAuthToInt(auth string) int {
	if auth == "ReadWrite" {
		return 2
	}
	return 1
}

func chunkStrings(src []string, size int) [][]string {
	if size <= 0 {
		size = 500
	}
	var out [][]string
	for i := 0; i < len(src); i += size {
		end := i + size
		if end > len(src) {
			end = len(src)
		}
		out = append(out, src[i:end])
	}
	return out
}

// ModbusBulkUpsertRegisterAddresses 缓存预取 + 分批写入 + 有限并发更新。
// 语义与逐行调用 ModbusRegisterAddressAdd/Update 一致，但避免 N+1。
func ModbusBulkUpsertRegisterAddresses(items []ModbusDevicesDataModel) ModbusBulkUpsertResult {
	var result ModbusBulkUpsertResult
	if len(items) == 0 {
		return result
	}

	existByUuid := make(map[string]ModbusDevicesDataModel, len(items))
	existByKey := make(map[string]ModbusDevicesDataModel, len(items))

	uuidList := make([]string, 0, len(items))
	muidSet := make(map[string]struct{}, 64)
	for _, it := range items {
		if it.Uuid != "" {
			uuidList = append(uuidList, it.Uuid)
		}
		if it.Muid != "" {
			muidSet[it.Muid] = struct{}{}
		}
	}

	for _, chunk := range chunkStrings(uuidList, modbusBulkPrefetchChunk) {
		var found []ModbusDevicesDataModel
		if err := Db.Model(&ModbusDevicesDataModel{}).Where("uuid IN ?", chunk).Find(&found).Error; err != nil {
			logs.Error("ModbusBulkUpsert prefetch by uuid failed: %v", err)
			continue
		}
		for _, row := range found {
			existByUuid[row.Uuid] = row
			existByKey[modbusAddrKey(row.Muid, row.RegisterGroupUuid, row.RegisterAddress)] = row
		}
	}

	muidList := make([]string, 0, len(muidSet))
	for m := range muidSet {
		muidList = append(muidList, m)
	}
	for _, chunk := range chunkStrings(muidList, modbusBulkPrefetchChunk) {
		var found []ModbusDevicesDataModel
		if err := Db.Model(&ModbusDevicesDataModel{}).
			Select("uuid", "muid", "register_group_uuid", "register_address").
			Where("muid IN ?", chunk).
			Find(&found).Error; err != nil {
			logs.Error("ModbusBulkUpsert prefetch by muid failed: %v", err)
			continue
		}
		for _, row := range found {
			existByKey[modbusAddrKey(row.Muid, row.RegisterGroupUuid, row.RegisterAddress)] = row
			if _, ok := existByUuid[row.Uuid]; !ok {
				existByUuid[row.Uuid] = row
			}
		}
	}

	toAdd := make([]ModbusDevicesDataModel, 0, len(items)/4)
	toUpdate := make([]ModbusDevicesDataModel, 0, len(items))
	seenInFile := make(map[string]struct{}, len(items))

	for _, it := range items {
		if it.Muid == "" || it.RegisterGroupUuid == "" {
			result.Skipped++
			continue
		}
		fileKey := modbusAddrKey(it.Muid, it.RegisterGroupUuid, it.RegisterAddress)
		if _, dup := seenInFile[fileKey]; dup {
			result.Skipped++
			continue
		}
		seenInFile[fileKey] = struct{}{}

		if it.Uuid != "" {
			if exist, ok := existByUuid[it.Uuid]; ok && exist.Muid == it.Muid {
				it.Uuid = exist.Uuid
				toUpdate = append(toUpdate, it)
				continue
			}
		}
		if exist, ok := existByKey[fileKey]; ok {
			it.Uuid = exist.Uuid
			toUpdate = append(toUpdate, it)
			continue
		}
		if it.Uuid == "" {
			it.Uuid = uuid.New()
		}
		toAdd = append(toAdd, it)
		existByKey[fileKey] = it
		existByUuid[it.Uuid] = it
	}

	logs.Info("ModbusBulkUpsert: total=%d add=%d update=%d skipped=%d",
		len(items), len(toAdd), len(toUpdate), result.Skipped)

	result.Updated = modbusBulkUpdateParallel(toUpdate)
	added, addSkipped := modbusBulkAddBatched(toAdd)
	result.Added = added
	result.Skipped += addSkipped
	return result
}

func modbusBulkUpdateParallel(rows []ModbusDevicesDataModel) int {
	if len(rows) == 0 {
		return 0
	}
	workers := modbusBulkUpdateWorkers
	if workers > len(rows) {
		workers = len(rows)
	}
	jobs := make(chan ModbusDevicesDataModel, workers*2)
	var okCount int64
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range jobs {
				if ModbusRegisterAddressUpdate(row) == errmsg.SUCCSE {
					atomic.AddInt64(&okCount, 1)
				}
			}
		}()
	}
	for _, row := range rows {
		jobs <- row
	}
	close(jobs)
	wg.Wait()
	return int(okCount)
}

func modbusBulkAddBatched(rows []ModbusDevicesDataModel) (added int, skipped int) {
	if len(rows) == 0 {
		return 0, 0
	}

	if err := Db.Model(&ModbusDevicesDataModel{}).CreateInBatches(&rows, modbusBulkModelBatchSize).Error; err != nil {
		logs.Error("ModbusBulkUpsert CreateInBatches data model failed, fallback per-row: %v", err)
		return modbusBulkAddFallback(rows)
	}
	added = len(rows)

	muidSet := make(map[string]struct{}, 64)
	for _, r := range rows {
		muidSet[r.Muid] = struct{}{}
	}
	muidList := make([]string, 0, len(muidSet))
	for m := range muidSet {
		muidList = append(muidList, m)
	}

	devicesByMuid := make(map[string][]MonitorList, len(muidList))
	for _, chunk := range chunkStrings(muidList, modbusBulkPrefetchChunk) {
		var devices []MonitorList
		if err := Db.Model(&MonitorList{}).
			Select("uuid", "name", "project_uuid", "muid").
			Where("muid IN ?", chunk).
			Find(&devices).Error; err != nil {
			logs.Error("ModbusBulkUpsert prefetch devices failed: %v", err)
			continue
		}
		for _, d := range devices {
			devicesByMuid[d.Muid] = append(devicesByMuid[d.Muid], d)
		}
	}

	realBuf := make([]DeviceRealData, 0, modbusBulkRealDataFlushAt)
	flushReal := func() {
		if len(realBuf) == 0 {
			return
		}
		if err := Db.Model(&DeviceRealData{}).CreateInBatches(&realBuf, modbusBulkRealDataBatch).Error; err != nil {
			logs.Error("ModbusBulkUpsert CreateInBatches device_real_data failed: %v", err)
		}
		realBuf = realBuf[:0]
	}

	for _, addData := range rows {
		devs := devicesByMuid[addData.Muid]
		if len(devs) == 0 {
			continue
		}
		auth := modbusAuthToInt(addData.Auth)
		for _, v := range devs {
			realBuf = append(realBuf, DeviceRealData{
				Auth:                 auth,
				ProjectUuid:          v.ProjectUuid,
				DeviceUuid:           v.Uuid,
				DeviceName:           v.Name,
				Name:                 addData.Name,
				Uuid:                 uuid.New(),
				ModelDataUuid:        addData.Uuid,
				Type:                 1,
				DataUnit:             addData.DataUnit,
				IsAlarm:              addData.IsAlarm,
				AlarmOnValue:         addData.AlarmOnValue,
				AlarmLevel:           addData.AlarmLevel,
				AlarmMessage:         addData.AlarmMessage,
				AlarmClearMessage:    addData.AlarmClearMessage,
				IsRecord:             addData.IsRecord,
				RecordType:           addData.RecordType,
				RecordInterval:       addData.RecordInterval,
				RecordDataCharge:     addData.RecordDataCharge,
				ConversionExpression: addData.ConversionExpression,
				Value:                "",
				Muid:                 addData.Muid,
				DeviceType:           2,
			})
			if len(realBuf) >= modbusBulkRealDataFlushAt {
				flushReal()
			}
		}
	}
	flushReal()
	return added, 0
}

func modbusBulkAddFallback(rows []ModbusDevicesDataModel) (added int, skipped int) {
	for _, row := range rows {
		// ModbusRegisterAddressAdd 会覆盖 Uuid；保留已分配 uuid 需直接 Create
		if row.Uuid == "" {
			row.Uuid = uuid.New()
		}
		if err := Db.Model(&ModbusDevicesDataModel{}).Create(&row).Error; err != nil {
			skipped++
			continue
		}
		added++
		var devices []MonitorList
		if err := Db.Model(&MonitorList{}).Where("muid = ?", row.Muid).Find(&devices).Error; err != nil || len(devices) == 0 {
			continue
		}
		auth := modbusAuthToInt(row.Auth)
		realRows := make([]DeviceRealData, 0, len(devices))
		for _, v := range devices {
			realRows = append(realRows, DeviceRealData{
				Auth:                 auth,
				ProjectUuid:          v.ProjectUuid,
				DeviceUuid:           v.Uuid,
				DeviceName:           v.Name,
				Name:                 row.Name,
				Uuid:                 uuid.New(),
				ModelDataUuid:        row.Uuid,
				Type:                 1,
				DataUnit:             row.DataUnit,
				IsAlarm:              row.IsAlarm,
				AlarmOnValue:         row.AlarmOnValue,
				AlarmLevel:           row.AlarmLevel,
				AlarmMessage:         row.AlarmMessage,
				AlarmClearMessage:    row.AlarmClearMessage,
				IsRecord:             row.IsRecord,
				RecordType:           row.RecordType,
				RecordInterval:       row.RecordInterval,
				RecordDataCharge:     row.RecordDataCharge,
				ConversionExpression: row.ConversionExpression,
				Value:                "",
				Muid:                 row.Muid,
				DeviceType:           2,
			})
		}
		_ = Db.Model(&DeviceRealData{}).CreateInBatches(&realRows, modbusBulkRealDataBatch).Error
	}
	return added, skipped
}
