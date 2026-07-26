/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-06-15 10:33:48
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package ISMScript

import (
	"ISMServer/models"
	protocolCommon "ISMServer/protocol/common"
	protocolCommonFunc "ISMServer/protocol/commFunc"
	"ISMServer/task/ISMScript/bitunpack"
	ISMScriptFunc "ISMServer/task/ISMScript/func"
	"ISMServer/utils/errmsg"
	"encoding/base64"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"time"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/mattn/anko/vm"
)

var GScriptChan chan bool
var scriptWg sync.WaitGroup
var valueChangeHookOnce sync.Once
var valueChangeDepth int32

func init() {
	valueChangeHookOnce.Do(func() {
		bitunpack.Configure(ISMScriptFunc.SetDeviceData, ISMScriptFunc.PeekDeviceDataValue)
		protocolCommon.SetDeviceValueChangeHandler(onDeviceValueChanged)
	})
}

func onDeviceValueChanged(deviceName, pointName, oldValue, newValue string) {
	// Prevent deep re-entrancy when BitUnpack SetDeviceData writes targets.
	if atomic.AddInt32(&valueChangeDepth, 1) > 16 {
		atomic.AddInt32(&valueChangeDepth, -1)
		return
	}
	defer atomic.AddInt32(&valueChangeDepth, -1)

	bitunpack.ApplySource(deviceName, pointName, newValue)
	WakeScriptsForKey(deviceName + "->" + pointName)
	_ = oldValue
}

func isChanClose() bool {
	select {
	case _, received := <-GScriptChan:
		return !received
	default:
	}
	return false
}

func ScriptCloseChan() {
	isOpen := isChanClose()
	if !isOpen && GScriptChan != nil {
		close(GScriptChan)
	}
}
func getAllScript() (int, []models.ISMScript) {
	var GetScrpt []models.ISMScript
	err := models.Db.Model(&models.ISMScript{}).Where("ID >= 0 and script_type=0 and is_disable=0").Find(&GetScrpt).Error
	if err != nil {
		return errmsg.ERROR, GetScrpt
	}
	return errmsg.SUCCSE, GetScrpt
}

func decodeScriptContent(script *models.ISMScript) {
	tempComponents, deErr := base64.StdEncoding.DecodeString(script.ScriptContent)
	if deErr == nil {
		script.ScriptContent = string(tempComponents)
	}
}

func ISMScriptMailPthread() {
	var is_starting = 0
	go StartSysScript()
	for {
		if is_starting == 1 {
			scriptWg.Wait()
		}
		ScriptCloseChan()
		GScriptChan = make(chan bool)

		bitunpack.Clear()
		clearScriptWakes()

		code, scriptList := getAllScript()
		nativeCount := 0
		ankoCount := 0
		if code == errmsg.SUCCSE && (len(scriptList) > 0) {
			for _, script := range scriptList {
				decodeScriptContent(&script)
				if rules, ok := bitunpack.Compile(script.ScriptUuid, script.ScriptName, script.ScriptContent); ok {
					bitunpack.Register(rules)
					logs.Info("native-bitunpack: %s rules=%d", script.ScriptName, len(rules))
					nativeCount++
					continue
				}
				if script.Delay < 1000 {
					script.Delay = 1000
				}
				d := &ISMScriptPthread{Script: script}
				scriptWg.Add(1)
				go d.Run()
				ankoCount++
			}
			if nativeCount > 0 {
				bitunpack.SettleAll()
			}
			logs.Info("script scheduler: native-bitunpack=%d anko-onchange=%d", nativeCount, ankoCount)
		}

		if ankoCount > 0 {
			// Wait on next loop until CRUD closes GScriptChan and anko workers exit.
			is_starting = 1
			time.Sleep(1 * time.Second)
		} else if nativeCount > 0 {
			// Pure native path: no WaitGroup workers; block until reload signal.
			is_starting = 0
			<-GScriptChan
		} else {
			is_starting = 0
			time.Sleep(time.Millisecond * 1000)
		}
	}
}
func StartSysScript() {

	var systemScriptPath string = "sys_script/"

	dirs, err := ioutil.ReadDir(systemScriptPath)
	if err != nil {
		logs.Error("不能读取系统脚本目录")
	} else {
		for _, file := range dirs {

			if !file.IsDir() {

				go func(filepath string) {
					content, err := ioutil.ReadFile(filepath)
					if err != nil {
						return
					}
					GoSysSctVm := protocolCommonFunc.ScriptDefine()
					_, scripterr := vm.Execute(GoSysSctVm, nil, string(content))

					if scripterr != nil {
						protocolCommon.ErrorThrottled("sysscript:"+filepath, "%s,Execute error: %v", filepath, scripterr)
					}
				}(systemScriptPath + file.Name())

			}
		}
	}

}
