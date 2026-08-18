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
	protocolCommon "ISMServer/protocol/common"
	protocolCommonFunc "ISMServer/protocol/commFunc"
	"time"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/mattn/anko/ast"
	"github.com/mattn/anko/env"
	"github.com/mattn/anko/parser"
	"github.com/mattn/anko/vm"
)

type ISMScriptPthread struct {
	Script models.ISMScript
	env    *env.Env
	stmt   ast.Stmt
	wakeCh chan struct{}
	deps   []string
}

// ExecScript runs a script once (used by manual/task paths). Still parses each call.
func ExecScript(sct models.ISMScript) {
	if protocolCommon.IsRestoreDb == 1 {
		return
	}
	GoSctVm := protocolCommonFunc.ScriptDefine()
	_, err := vm.Execute(GoSctVm, nil, sct.ScriptContent)
	if err != nil {
		protocolCommon.ErrorThrottled("script:"+sct.ScriptName, "%s,Execute error: %v", sct.ScriptName, err)
	}
	time.Sleep(time.Millisecond * time.Duration(sct.Delay))
}

func (t *ISMScriptPthread) prepare() error {
	t.env = protocolCommonFunc.ScriptDefine()
	stmt, err := parser.ParseSrc(t.Script.ScriptContent)
	if err != nil {
		return err
	}
	t.stmt = stmt
	t.deps = ExtractScriptDeps(t.Script.ScriptContent)
	t.wakeCh = make(chan struct{}, 1)
	registerScriptWake(t.Script.ScriptUuid, t.deps, t.wakeCh)
	return nil
}

func (t *ISMScriptPthread) runOnce() {
	if protocolCommon.IsRestoreDb == 1 {
		return
	}
	_, err := vm.Run(t.env, nil, t.stmt)
	if err != nil {
		protocolCommon.ErrorThrottled("script:"+t.Script.ScriptName, "%s,Execute error: %v", t.Script.ScriptName, err)
	}
}

func (t *ISMScriptPthread) Run() {
	defer scriptWg.Done()

	if err := t.prepare(); err != nil {
		protocolCommon.ErrorThrottled("script:"+t.Script.ScriptName, "%s,Parse error: %v", t.Script.ScriptName, err)
		logs.Info("anko-onchange: %s parse failed, idle until reload", t.Script.ScriptName)
		for {
			select {
			case <-GScriptChan:
				logs.Info("脚本主动退出", t.Script.ScriptName)
				return
			case <-time.After(time.Second):
			}
		}
	}

	delay := t.Script.Delay
	if delay < 1000 {
		delay = 1000
	}
	logs.Info("anko-onchange: %s deps=%d delay=%d", t.Script.ScriptName, len(t.deps), delay)

	// Run once on start so outputs are consistent before first change/tick.
	t.runOnce()

	ticker := time.NewTicker(time.Millisecond * time.Duration(delay))
	defer ticker.Stop()

	for {
		select {
		case <-GScriptChan:
			logs.Info("脚本主动退出", t.Script.ScriptName)
			return
		case <-t.wakeCh:
			// drain coalesced wakes
			for {
				select {
				case <-t.wakeCh:
					continue
				default:
				}
				break
			}
			t.runOnce()
		case <-ticker.C:
			t.runOnce()
		}
	}
}
