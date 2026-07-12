/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-09-06 16:08:57
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package syncDataTask

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
)

// FullPrewarmEnabled is a temporary rollback switch. The default live mode
// avoids loading the complete device_real_data table during process startup.
func FullPrewarmEnabled() bool {
	mode, err := config.String("realdata_prewarm")
	return err == nil && strings.EqualFold(strings.TrimSpace(mode), "full")
}

// SyncDevicesDataToMemory keeps the legacy full prewarm available only as an
// explicit rollback mode. Normal startup relies on DB snapshots plus live data.
func SyncDevicesDataToMemory() {
	logs.Info("正在同步设备数据到内存...")
	r := time.Now()
	var lastID uint
	for {
		var batch []models.DeviceRealData
		if err := models.Db.Model(&models.DeviceRealData{}).
			Select("id, device_uuid, uuid, value, device_name, name").
			Where("id > ?", lastID).
			Order("id ASC").
			Limit(5000).
			Find(&batch).Error; err != nil {
			logs.Error("同步设备数据到内存失败: %v", err)
			return
		}
		if len(batch) == 0 {
			break
		}
		for _, realData := range batch {
			protocol_common.DeviceRealDataMapByUUID.Store(realData.DeviceUuid+realData.Uuid, realData.Value)
			protocol_common.DeviceRealDataMapByUUID.Store(realData.Uuid, realData.Value)
			protocol_common.DeviceRealDataMap.Store(realData.DeviceName+"->"+realData.Name, realData.Value)
		}
		lastID = batch[len(batch)-1].ID
		// SQLite 单连接环境下主动让出查询窗口，避免启动阶段登录/API 被连续批次饿死。
		time.Sleep(25 * time.Millisecond)
	}
	d := time.Since(r)
	logs.Info("数据同步完成,耗时:%s", d)
}
