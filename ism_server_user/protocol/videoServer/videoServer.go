/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:01:09
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package videoserver

import (
	protocolCommon "ISMServer/protocol/common"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/format/mp4"
	"github.com/shirou/gopsutil/process"
)

var ConfigName = "conf/videoConfig.json"

var VideoWg sync.WaitGroup
var GVideoChan chan bool
var MonibucaServerProcess *exec.Cmd

func isChanClose() bool {
	select {
	case _, received := <-GVideoChan:
		return !received
	default:
	}
	return false
}
func testmp4() {
	filePath := filepath.Join(protocolCommon.RecordPath, "测试/2024-07-22/2024_07_22_16_10_23.mp4")
	RecordFw := &FileWriter{filePath: filePath}

	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err == nil {
		RecordFw.Reader = file
		RecordFw.Writer = file
		RecordFw.Seeker = file
		RecordFw.Closer = file
		fmt.Println(filePath)
		Recorder := mp4.NewMuxer(RecordFw)
		Recorder.WriteTrailer()
		RecordFw.Closer.Close()
	}
}
func VideoCloseChan() {

	isOpen := isChanClose()
	if !isOpen && GVideoChan != nil {
		close(GVideoChan)
	}
}
func Server() {
	var is_starting = 0
	// testmp4()
	for {

		if is_starting == 1 {
			VideoWg.Wait()
		}
		//等待数据库还原
		if protocolCommon.IsRestoreDb == 1 {
			time.Sleep(time.Second * 5)
			continue
		}
		GVideoChan = make(chan bool)
		Config = LoadConfig()
		is_starting = 0
		if len(Config.Streams) > 0 {
			for _, v := range Config.Streams {
				if v.IsUsed == 0 {
					d := &VideoCtl{waitGroup: &VideoWg, RecordInter: v.RecordInter, Key: v.Uuid, name: v.Name, url: v.URL, project_uuid: v.ProjectUuid, RecordEnable: v.RecordEnable, IsUsed: v.IsUsed}
					d.InitVideoDevice()
					go d.VideoGatherPthread()
					VideoWg.Add(1)
					is_starting = 1
				}
			}
		}
		if is_starting == 0 {
			time.Sleep(time.Second * 5)
		}
	}

}

func CheckProcessRunning() bool {

	pids, _ := process.Pids()
	for _, pid := range pids {
		pn, _ := process.NewProcess(pid)
		pName, _ := pn.Name()
		if strings.Contains(pName, "monibucaserverarm") || strings.Contains(pName, "MonibucaServerArm") || strings.Contains(pName, "MonibucaServer") || strings.Contains(pName, "monibucaserver") {
			return true
		}
	}
	return false
}
func ExecCommand(name string, args ...string) {
	MonibucaServerProcess = exec.Command(name, args...) // 拼接参数与命令
	var err error

	if err = MonibucaServerProcess.Start(); err != nil {
		logs.Error("国标视频服务启动失败")
		log.Println(err)
	} else {
		logs.Info("国标视频服务启动成功")
	}
}
func StartGB28281Server() {

	for {
		isExist := CheckProcessRunning()
		if !isExist {
			if runtime.GOARCH == "amd64" || runtime.GOARCH == "386" {
				if runtime.GOOS != "windows" {
					os.Remove("./vendorBin/MonibucaServer.exe")
					if _, err := os.Stat("./vendorBin/MonibucaServer"); err == nil {
						logs.Info("正在启动国标视频服务...")
						ExecCommand("./vendorBin/MonibucaServer", "-c", "./conf/MonibucaServer.yaml")
					}
				} else {
					os.Remove("./vendorBin/MonibucaServer")
					if _, err := os.Stat("./vendorBin/MonibucaServer.exe"); err == nil {
						logs.Info("正在启动国标视频服务...")
						ExecCommand("./vendorBin/MonibucaServer.exe", "-c", "./conf/MonibucaServer.yaml")
					}
				}
			} else if runtime.GOARCH == "arm64" {
				os.Remove("./vendorBin/MonibucaServer")
				os.Remove("./vendorBin/MonibucaServer.exe")
				if _, err := os.Stat("./vendorBin/MonibucaServerArm"); err == nil {
					logs.Info("正在启动国标视频服务...")
					ExecCommand("./vendorBin/MonibucaServerArm", "-c", "./conf/MonibucaServer.yaml")
				}
			} else {
				logs.Error("MonibucaServer 不支持当前系统架构", runtime.GOARCH)
			}
			if MonibucaServerProcess != nil {
				MonibucaServerProcess.Wait()
			}
		}
		time.Sleep(time.Second * 4)
	}
}
