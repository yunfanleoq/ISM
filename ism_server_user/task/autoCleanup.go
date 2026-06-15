/**
 * @ Description: 自动清理历史数据，防止数据库膨胀影响性能
 * @ Author: ISM Web组态软件
 * @ Create Time: 2024-06-15
 */

package tasks

import (
	"ISMServer/models"
	protocol_common "ISMServer/protocol/common"
	"time"

	"github.com/beego/beego/v2/core/config"
	gormlog "github.com/beego/beego/v2/core/logs"
)

func StartAutoCleanup() {
	go func() {
		// 启动后先等30秒，让数据库完全就绪
		time.Sleep(30 * time.Second)

		// 立即执行一次清理
		cleanHistoryData()
		cleanAlarmData()

		// 之后每小时执行一次
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			cleanHistoryData()
			cleanAlarmData()
		}
	}()
}

func cleanHistoryData() {
	days := -protocol_common.HistoryKeepDays
	cutoff := time.Now().AddDate(0, 0, days)
	gormlog.Info("自动清理: 删除 %s 之前的历史数据", cutoff.Format("2006-01-02 15:04:05"))

	result := models.Db.Unscoped().
		Where("record_time < ?", cutoff).
		Delete(&models.DevicesHistoryDataList{})

	if result.Error != nil {
		gormlog.Error("自动清理历史数据失败: %s", result.Error)
	} else {
		gormlog.Info("自动清理完成: 删除历史数据 %d 条", result.RowsAffected)

		// SQLite 需要 VACUUM 回收空间，OceanBase/MySQL 不需要
		dbtype, _ := config.Int("dbtype")
		if dbtype == 1 {
			models.Db.Exec("VACUUM;")
		}
	}
}

func cleanAlarmData() {
	days := -protocol_common.HistoryKeepDays
	cutoff := time.Now().AddDate(0, 0, days)

	result := models.Db.Unscoped().
		Where("clear_time < ? AND clear_time IS NOT NULL AND clear_time > '2000-01-01'", cutoff).
		Delete(&models.DevicesAlarmList{})

	if result.Error != nil {
		gormlog.Error("自动清理告警数据失败: %s", result.Error)
	} else if result.RowsAffected > 0 {
		gormlog.Info("自动清理完成: 删除告警数据 %d 条", result.RowsAffected)
	}
}
