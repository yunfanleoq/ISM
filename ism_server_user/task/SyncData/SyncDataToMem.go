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
	"time"

	"github.com/beego/beego/v2/core/logs"
)

// UpgradeStaticData 升级静态数据函数
func SyncDevicesDataToMemory() {
	logs.Info("正在同步设备数据到内存...")
	r := time.Now()
	var DeviceRealData []models.DeviceRealData
	models.Db.Model(&models.DeviceRealData{}).Where("id > 0").Find(&DeviceRealData)
	if len(DeviceRealData) != 0 {
		for _, RealData := range DeviceRealData {
			protocol_common.DeviceRealDataMapByUUID.Store(RealData.DeviceUuid+RealData.Uuid, RealData.Value)
			protocol_common.DeviceRealDataMapByUUID.Store(RealData.Uuid, RealData.Value)
			protocol_common.DeviceRealDataMap.Store(RealData.DeviceName+"->"+RealData.Name, RealData.Value)
		}
	}
	d := time.Since(r)
	logs.Info("数据同步完成,耗时:%s", d)
}
