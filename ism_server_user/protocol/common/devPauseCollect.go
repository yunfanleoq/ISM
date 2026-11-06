package protocolCommon

import "github.com/beego/beego/v2/core/config"

// DevPauseCollect 开发期暂停设备采集、重连和定时历史快照。
// 配置见 conf/app.conf 的 dev_pause_collect。缺省或读失败时保持原有采集行为。
func DevPauseCollect() bool {
	v, err := config.Bool("dev_pause_collect")
	if err != nil {
		return false
	}
	return v
}
