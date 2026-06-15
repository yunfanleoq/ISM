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
	protocolCommonFunc "ISMServer/protocol/commFunc"
	"ISMServer/utils/errmsg"
	"encoding/base64"
	"io/ioutil"
	"sync"
	"time"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/mattn/anko/vm"
)

var GScriptChan chan bool
var scriptWg sync.WaitGroup

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
func ISMScriptMailPthread() {
	var is_starting = 0
	go StartSysScript()
	for {
		if is_starting == 1 {
			scriptWg.Wait()
		}
		ScriptCloseChan()
		GScriptChan = make(chan bool)
		code, scriptList := getAllScript()
		if code == errmsg.SUCCSE && (len(scriptList) > 0) {
			for _, script := range scriptList {
				if script.Delay < 100 {
					script.Delay = 100
				}
				tempComponents, deErr := base64.StdEncoding.DecodeString(script.ScriptContent)
				if deErr == nil {
					script.ScriptContent = string(tempComponents)
				}
				d := &ISMScriptPthread{Script: script}
				go d.Run()
				scriptWg.Add(1)
			}
		} else {
			time.Sleep(time.Millisecond * 1000)
		}
		is_starting = 1
		time.Sleep(1 * time.Second)
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
						logs.Error("%s,Execute error: %v", filepath, scripterr)
					}
				}(systemScriptPath + file.Name())

			}
		}
	}

}
