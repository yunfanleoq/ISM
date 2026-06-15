/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:02:14
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package ISMScript

import (
	"ISMServer/models"
	protocolCommonFunc "ISMServer/protocol/commFunc"
	"time"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/mattn/anko/vm"
)

type ISMScriptPthread struct {
	Script models.ISMScript
}

func ExecScript(sct models.ISMScript) {
	GoSctVm := protocolCommonFunc.ScriptDefine()
	_, err := vm.Execute(GoSctVm, nil, sct.ScriptContent)
	if err != nil {
		logs.Error("%s,Execute error: %v", sct.ScriptName, err)
	}
	time.Sleep(time.Millisecond * time.Duration(sct.Delay))
}

func (t *ISMScriptPthread) Run() {

	for {
		//检测协程是否主动退出
		select {
		case <-GScriptChan:
			scriptWg.Done()
			logs.Info("脚本主动退出", t.Script.ScriptName)
			return
		default:
		}
		ExecScript(t.Script)
	}
}
