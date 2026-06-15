/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-08-02 15:37:42
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package snmpprotocols

import (
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	"sync"
	"time"
)

type deviceStu struct {
	Name                 string
	Uuid                 string
	ExtraData            string
	Muid                 string
	Port                 int
	Timeout              int
	Interval             int
	GatherNumber         int
	FailedTimes          int
	Version              int
	IsEnable             int
	Writecomm            string
	Readcomm             string
	SnmpUserName         string
	SnmpSecurityLevel    int
	SnmpAuthAlgorithm    int
	SnmpUserPassword     string
	SnmpPrivacyAlgorithm int
	SnmpPrivacyPassword  string
	ProjectUuid          string
}

var GSnmpChan chan bool

var wg sync.WaitGroup

func isChanClose() bool {
	select {
	case _, received := <-GSnmpChan:
		return !received
	default:
	}
	return false
}

func SnmpCloseChan() {

	isOpen := isChanClose()
	if !isOpen && GSnmpChan != nil {
		close(GSnmpChan)
	}
}
func waitForGather() {
	//do nothing
	for {
		select {
		case <-GSnmpChan:
			wg.Done()
			return
		default:
		}
		time.Sleep(time.Second * 5)
	}
}
func SnmpServer() {

	var is_starting = 0
	for {

		if is_starting == 1 {
			wg.Wait()
		}
		//等待数据库还原
		if protocolCommon.IsRestoreDb == 1 {
			time.Sleep(time.Second * 5)
			continue
		}
		SnmpCloseChan()
		GSnmpChan = make(chan bool)
		var results []deviceStu

		models.Db.Raw("SELECT monitor_list.project_uuid,monitor_list.uuid,monitor_list.name,monitor_list.is_enable, monitor_list.extra_data,monitor_list.muid ,devices_model.port,monitor_list.interval,monitor_list.timeout,devices_model.gather_number,monitor_list.failed_times,devices_model.version,devices_model.writecomm,devices_model.readcomm,devices_model.snmp_user_name,devices_model.snmp_security_level,devices_model.snmp_auth_algorithm,devices_model.snmp_user_password,devices_model.snmp_privacy_algorithm,devices_model.snmp_privacy_password FROM monitor_list, devices_model WHERE monitor_list.muid = devices_model.uuid and devices_model.type=1").Scan(&results)
		var gather_is_start int = 0
		if len(results) > 0 {
			for _, device := range results {
				if device.IsEnable == 0 {
					continue
				}
				var getOids []models.DeviceRealData
				models.Db.Model(&models.DeviceRealData{}).Where("device_uuid = ? and oid !=''", device.Uuid).Find(&getOids)
				d := &SnmpCtl{waitGroup: &wg, failedTimes: 0, deviceStatus: 0}
				d.InitDeviceSnmpInfo(device, getOids)
				go d.GatherSnmpOids()
				wg.Add(1)
				gather_is_start = 1
			}
		} else {
			time.Sleep(time.Second * 5)
		}

		if gather_is_start == 0 {
			go waitForGather()
			wg.Add(1)
		}

		is_starting = 1
	}
}
