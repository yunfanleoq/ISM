/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-07-07 14:59:26
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package main

import (
	ISMVendor "ISMServer/ISMVendor"
	"ISMServer/controllers"
	license "ISMServer/license"
	protocols "ISMServer/protocol"
	protocol_common "ISMServer/protocol/common"
	_ "ISMServer/routers"
	tasks "ISMServer/task"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	ISMServer "github.com/beego/beego/v2/server/web"
	"github.com/fsnotify/fsnotify"
)

/*
#include<stdint.h>
#include<string.h>
void getCompileDateTime(uint8_t  dt[12],uint8_t tm[9]){
  strcpy(dt, __DATE__); //Mmm dd yyyy
  strcpy(tm,__TIME__);  //hh:mm:ss
}
*/
import "C"

// 不同平台启动指令不同
var commands = map[string]string{
	"windows": "explorer",
	"darwin":  "open",
	"linux":   "xdg-open",
}

const VERSION string = "V3.01.RC07"
const VERSION_DATE string = "2026年03月17日"

// 验证软件是否过期
// createDateStr: 创建时间字符串（如 "2023-10-01"）
// authDays: 授权天数
// 返回值：true=已过期，false=未过期，error=解析失败
func SoftAuthIsExpired() (bool, error) {

	beijingLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return false, fmt.Errorf("加载北京时间时区失败: %v", err)
	}

	layout := "2006-01-02_15-04-05"
	createDate, err := time.ParseInLocation(layout, protocol_common.AuthorizationCreateDate, beijingLoc)
	if err != nil {
		return false, fmt.Errorf("解析创建时间失败：%v", err)
	}
	for {
		if !protocol_common.IsAuthLimit {
			protocol_common.IsAuthTimeLimit = false
			time.Sleep(10 * time.Second)
		}
		// 2. 计算过期时间（创建时间 + 授权天数）
		expireDate := createDate.AddDate(0, 0, protocol_common.AuthorizationDays)
		// 3. 对比当前时间（使用本地时间或UTC，根据需求调整）
		now := time.Now().In(beijingLoc)
		remaining := expireDate.Sub(now)
		if remaining <= 0 {
			protocol_common.IsAuthTimeLimit = true
		} else {
			protocol_common.IsAuthTimeLimit = false
			// 计算剩余的天数、小时和分钟
			days := int(remaining.Hours() / 24)
			remainingAfterDays := remaining - time.Duration(days)*24*time.Hour

			hours := int(remainingAfterDays.Hours())
			remainingAfterHours := remainingAfterDays - time.Duration(hours)*time.Hour

			minutes := int(remainingAfterHours.Minutes())

			protocol_common.AuthRemainingTimeDays = days
			protocol_common.AuthRemainingTimeHours = hours
			protocol_common.AuthRemainingTimeMinutes = minutes
		}
		time.Sleep(1 * time.Minute)
	}
}
func Open(uri string) error {
	// runtime.GOOS获取当前平台
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	cmd := exec.Command(run, uri)
	return cmd.Run()
}
func licenseConfigCheck() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logs.Error("new watcher failed: ", err)
	}
	defer watcher.Close()
	for {
		err = watcher.Add("data/auth/active.dat")
		if err != nil {
			logs.Error("add failed:", err)
			time.Sleep(time.Second * 10)
			continue
		}
		break
	}
	for {
		select {
		case event := <-watcher.Events:
			{
				switch event.Op {
				case fsnotify.Create:
					logs.Info("授权文件已经更新")
					license.CheckLicense()
					break
				case fsnotify.Write:
					logs.Info("授权文件已经更新")
					license.CheckLicense()
					break
				case fsnotify.Remove:
					logs.Info("授权文件已经删除")
					protocol_common.IsLicense = false
				}
			}
		case err := <-watcher.Errors:
			logs.Error(err.Error())
		}
	}

}
func debugServer() {
	http.ListenAndServe(":28081", nil)
}
func freeMemory() {
	for {
		time.Sleep(time.Second * 600)
		logs.Info("开始清理内存和缓存")
		runtime.GC()
	}
}
func main() {

	var RecordPathErr error

	runtime.GOMAXPROCS(runtime.NumCPU())

	runtime.SetBlockProfileRate(1)

	dt := make([]byte, 12)
	tm := make([]byte, 10)
	C.getCompileDateTime((*C.uint8_t)(unsafe.Pointer(&dt[0])), (*C.uint8_t)(unsafe.Pointer(&tm[0])))
	dts, tms := string(dt), string(tm)

	logFilesSavaDays, logFilesSavaDaysErr_ := config.Int("logFilesSavaDays")
	if logFilesSavaDaysErr_ != nil {
		logFilesSavaDays = 3
	}
	//是否打开调试界面
	IsDebug, IsDebugerr := config.Bool("IsDebug")
	if IsDebugerr != nil {
		IsDebug = false
	}
	if IsDebug {
		go debugServer()
		logs.Info("debug server is running at 28081")
		logs.EnableFuncCallDepth(true)
		logs.SetLogFuncCallDepth(4)
	} else {
		logs.EnableFuncCallDepth(false)
		logs.SetLogFuncCallDepth(2)
	}
	//日志初始化
	err := logs.SetLogger(logs.AdapterFile, "{\"rotate\": true,\"filename\":\"logs/ism.log\",\"level\":3,\"daily\":true,\"maxdays\":"+fmt.Sprintf("%d", logFilesSavaDays)+",\"color\":true}")
	if err != nil {
		panic(err)
	}

	logs.Async(1e3)
	logs.Info("System is starting....")
	logs.Info("系统运行版本:%s,发布时期:%s", VERSION, dts+" "+tms)
	if _, err := os.Stat(controllers.SavePath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(controllers.SavePath, os.ModePerm)
		}
	}
	//导出的历史记录目录
	if _, err := os.Stat("static/HistoryData/"); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll("static/HistoryData/", os.ModePerm)
		}
	}
	//导出的历史记录模板
	if _, err := os.Stat("static/reportTemplete/"); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll("static/reportTemplete/", os.ModePerm)
		}
	}

	//授权目录
	if _, err := os.Stat("data/auth/"); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll("data/auth/", os.ModePerm)
		}
	}

	//录像目录
	protocol_common.RecordPath, RecordPathErr = config.String("RecordPath")
	if RecordPathErr != nil || protocol_common.RecordPath == "" {
		protocol_common.RecordPath = "static/RecordVideo/"
	}
	if _, err := os.Stat(protocol_common.RecordPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(protocol_common.RecordPath, os.ModePerm)
		}
	}
	logs.Info("录像目录", protocol_common.RecordPath)
	ISMServer.SetStaticPath("/record", protocol_common.RecordPath)

	//检查是否授权
	license.CheckLicense()

	//检测系统关闭事件
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdownChan
		fmt.Println("程序退出......")
		if protocol_common.GSaveHistoryDataLevelDb != nil {
			protocol_common.GSaveHistoryDataLevelDb.Close()
		}
		os.Exit(0)
	}()
	if protocol_common.IsAuthTimeLimit {
		logs.Error("软件授权已过期，请联系管理员获取正式授权")
	} else {
		//授权验证
		go SoftAuthIsExpired()
	}
	if !protocol_common.IsOem {
		logs.Error("目前使用的是个人免费版本,购买企业版本，请访问 www.ismctl.com 咨询购买。")
		Open("http://www.ismctl.com")
	}
	//定时清理内存
	// go freeMemory()
	//go lisceseConfigCheck()
	//初始化任务
	tasks.TasksServer()
	//协议解析初始化
	protocols.ProtocolsServer()
	//初始化第三方库
	ISMVendor.VendorServer()
	//事件总线启动
	// EventBusServer.EventBusStart()
	ISMPort, _ := config.Int("httpport")

	AutoOpenUrl, AutoOpenUrlErr := config.Bool("AutoOpenUrl")
	if AutoOpenUrlErr == nil && AutoOpenUrl {
		//使用系统默认浏览器打开
		Open(fmt.Sprintf("http://127.0.0.1:%d", ISMPort))
	}

	//web服务启动
	ISMServer.Run()
}
