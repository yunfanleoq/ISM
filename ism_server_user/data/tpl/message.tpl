告警主题:{{if eq  .DataName "device.DeviceStatus"}}设备状态{{else}}{{.DataName}}{{end}}
告警设备:{{ .DeviceName }}
告警详情:{{if eq  .Value "1"}}{{if eq  .AlarmMessage "device.DeviceStatusOffline"}}设备离线{{else if eq  .AlarmMessage "VideoManager.VideoOffline"}}视频设备离线{{else}}{{ .AlarmMessage }}{{end}}{{else}}{{if eq  .AlarmClearMessage "device.DeviceStatusOnline"}}设备在线{{else if eq  .AlarmClearMessage "VideoManager.VideoOnline"}}视频设备在线{{else}}{{ .AlarmClearMessage }}{{end}}{{end}}
