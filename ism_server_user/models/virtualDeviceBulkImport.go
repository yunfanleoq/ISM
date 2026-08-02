package models

import (
	"ISMServer/utils/errmsg"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/beego/beego/v2/core/logs"
)

const (
	virtualDeviceBulkPrefetchChunk = 500
	virtualDeviceBulkUpdateWorkers = 8
	virtualDeviceModelType         = 480
)

// VirtualDeviceBulkUpsertResult 虚拟设备全量点位导入统计
type VirtualDeviceBulkUpsertResult struct {
	Added   int
	Updated int
	Skipped int
}

func virtualDeviceNameKey(muid, name string) string {
	return fmt.Sprintf("%s\x00%s", muid, name)
}

// VirtualDeviceBulkUpsert 按 模型ID+数据ID / 模型ID+数据名称 upsert。
// 复用 VirtualDeviceDataAdd / VirtualDeviceDataEdit，保证 device_real_data 同步。
func VirtualDeviceBulkUpsert(items []VirtualDeviceDataModel) VirtualDeviceBulkUpsertResult {
	var result VirtualDeviceBulkUpsertResult
	if len(items) == 0 {
		return result
	}

	existByUuid := make(map[string]VirtualDeviceDataModel, len(items))
	existByName := make(map[string]VirtualDeviceDataModel, len(items))

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

	for _, chunk := range chunkStrings(uuidList, virtualDeviceBulkPrefetchChunk) {
		var found []VirtualDeviceDataModel
		if err := Db.Model(&VirtualDeviceDataModel{}).Where("uuid IN ?", chunk).Find(&found).Error; err != nil {
			logs.Error("VirtualDeviceBulkUpsert prefetch by uuid failed: %v", err)
			continue
		}
		for _, row := range found {
			existByUuid[row.Uuid] = row
			existByName[virtualDeviceNameKey(row.Muid, row.Name)] = row
		}
	}

	muidList := make([]string, 0, len(muidSet))
	for m := range muidSet {
		muidList = append(muidList, m)
	}
	for _, chunk := range chunkStrings(muidList, virtualDeviceBulkPrefetchChunk) {
		var found []VirtualDeviceDataModel
		if err := Db.Model(&VirtualDeviceDataModel{}).
			Select("uuid", "muid", "name").
			Where("muid IN ?", chunk).
			Find(&found).Error; err != nil {
			logs.Error("VirtualDeviceBulkUpsert prefetch by muid failed: %v", err)
			continue
		}
		for _, row := range found {
			existByName[virtualDeviceNameKey(row.Muid, row.Name)] = row
			if _, ok := existByUuid[row.Uuid]; !ok {
				existByUuid[row.Uuid] = row
			}
		}
	}

	toAdd := make([]VirtualDeviceDataModel, 0, len(items)/4)
	toUpdate := make([]VirtualDeviceDataModel, 0, len(items))
	seenInFile := make(map[string]struct{}, len(items))

	for _, it := range items {
		if it.Muid == "" || it.Name == "" {
			result.Skipped++
			continue
		}
		if it.ModelType == 0 {
			it.ModelType = virtualDeviceModelType
		}
		fileKey := virtualDeviceNameKey(it.Muid, it.Name)
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
		if exist, ok := existByName[fileKey]; ok {
			it.Uuid = exist.Uuid
			toUpdate = append(toUpdate, it)
			continue
		}
		toAdd = append(toAdd, it)
	}

	logs.Info("VirtualDeviceBulkUpsert: total=%d add=%d update=%d skipped=%d",
		len(items), len(toAdd), len(toUpdate), result.Skipped)

	result.Updated = virtualDeviceBulkUpdateParallel(toUpdate)
	added, addSkipped := virtualDeviceBulkAdd(toAdd)
	result.Added = added
	result.Skipped += addSkipped
	return result
}

func virtualDeviceBulkUpdateParallel(rows []VirtualDeviceDataModel) int {
	if len(rows) == 0 {
		return 0
	}
	workers := virtualDeviceBulkUpdateWorkers
	if workers > len(rows) {
		workers = len(rows)
	}
	jobs := make(chan VirtualDeviceDataModel, workers*2)
	var okCount int64
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range jobs {
				code := VirtualDeviceDataEdit(row.Muid, row.Uuid, row)
				if code == errmsg.SNMP_MODEL_ADD_SUCCSE || code == errmsg.SUCCSE {
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

func virtualDeviceBulkAdd(rows []VirtualDeviceDataModel) (added int, skipped int) {
	for _, row := range rows {
		code := VirtualDeviceDataAdd(row)
		if code == errmsg.SNMP_MODEL_ADD_SUCCSE {
			added++
		} else {
			skipped++
		}
	}
	return added, skipped
}
