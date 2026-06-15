package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/fatih/color"
	"github.com/getlantern/systray"
	"github.com/shiena/ansicolor"
	"github.com/shirou/gopsutil/process"
)

var upgradeDir string = "data/upgrade/"
var ISMProcess *exec.Cmd

func printLogo() {

	file, err := os.Open("conf/logo.txt")
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n') //注意是字符
		if err == io.EOF {
			if len(line) != 0 {
				color.Red(line)
			}
			break
		}
		if err != nil {
			fmt.Println("read file failed, err:", err)
			return
		}
		color.Red(line)
	}
}

func UnzipAndUpgrade(src string, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		filePath := path.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := file.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()
			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func ExecCommand(name string, args ...string) {
	w := ansicolor.NewAnsiColorWriter(os.Stdout)
	ISMProcess = exec.Command(name, args...) // 拼接参数与命令
	stdout, err3 := ISMProcess.StdoutPipe()
	if err3 != nil {
		fmt.Println("cmd.StdoutPipe: ", err3)
		return
	}
	var err error

	if err = ISMProcess.Start(); err != nil {
		log.Println(err)
	}

	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Fprint(w, line)
	}
}

func CheckProcessRunning() bool {

	pids, _ := process.Pids()
	for _, pid := range pids {
		pn, _ := process.NewProcess(pid)
		pName, _ := pn.Name()
		if strings.Contains(pName, "ISMServer") || strings.Contains(pName, "ismserver") {
			return true
		}
	}
	return false
}
func UpgradeFileExists(path string) bool {
	_, err := os.Stat(path)
	//当为空文件或文件夹存在
	if err == nil {
		return true
	}
	//os.IsNotExist(err)为true，文件或文件夹不存在
	if os.IsNotExist(err) {
		return false
	}
	//其它类型，不确定是否存在
	return false
}

// 托盘icon 图标
var Data []byte

func onReady() {
	systray.SetIcon(Data)
	systray.SetTitle("Server is running")
	systray.SetTooltip("服务正在运行, 右键点击打开菜单！")
	mShow := systray.AddMenuItem("显示", "显示窗口")
	systray.AddSeparator()
	mHide := systray.AddMenuItem("隐藏", "隐藏窗口")

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	user32 := syscall.NewLazyDLL("user32.dll")
	// https://docs.microsoft.com/en-us/windows/console/getconsolewindow
	getConsoleWindows := kernel32.NewProc("GetConsoleWindow")
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-showwindowasync
	showWindowAsync := user32.NewProc("ShowWindowAsync")
	consoleHandle, r2, err := getConsoleWindows.Call()
	if consoleHandle == 0 {
		fmt.Println("Error call GetConsoleWindow: ", consoleHandle, r2, err)
	}

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				mShow.Disable()
				mHide.Enable()
				r1, r2, err := showWindowAsync.Call(consoleHandle, 5)
				if r1 != 1 {
					fmt.Println("Error call ShowWindow @SW_SHOW: ", r1, r2, err)
				}
			case <-mHide.ClickedCh:
				mHide.Disable()
				mShow.Enable()
				r1, r2, err := showWindowAsync.Call(consoleHandle, 0)
				if r1 != 1 {
					fmt.Println("Error call ShowWindow @SW_HIDE: ", r1, r2, err)
				}
			}
		}
	}()
}
func onExit() {

}
func CheckAndRestart() {
	for {
		isExist := CheckProcessRunning()
		if !isExist {

			if UpgradeFileExists(upgradeDir + "lastversion.zip") {
				UnzipAndUpgrade(upgradeDir+"lastversion.zip", "")
				os.Remove(upgradeDir + "lastversion.zip")
				logs.Error("ISM 升级完成准备重启")
			}
			if runtime.GOARCH == "amd64" || runtime.GOARCH == "386" {
				if runtime.GOOS != "windows" {
					os.Remove("./ISMServer.exe")
					ExecCommand("./ISMServer", "")
				} else {
					os.Remove("./ISMServer")
					ExecCommand("./ISMServer.exe", "")
				}
			} else if (runtime.GOARCH == "arm64") || (runtime.GOARCH == "arm") {
				os.Remove("./ISMServer.exe")
				os.Remove("./ISMServer")
				ExecCommand("./ISMServerForArm", "")
			} else {
				logs.Error("ISM 不支持当前系统架构", runtime.GOARCH)
			}

			logs.Error("ISM 退出，准备重启启动")
			if ISMProcess != nil {
				// 托盘程序逻辑
				ISMProcess.Wait()
			}
		}
		time.Sleep(time.Second * 4)
	}
}
func main() {

	printLogo()

	iconfile, readErr := os.Open("static/favicon.ico")
	if readErr != nil {
		fmt.Println("读取托盘图标失败")
		fmt.Println(readErr)
	} else {
		Data, _ = io.ReadAll(iconfile)
		iconfile.Close()
	}
	//日志初始化
	err := logs.SetLogger(logs.AdapterFile, `{"filename":"logs/ismwatcher.log","level":7,"daily":true,"maxdays":30,"color":true}`)
	if err != nil {
		panic(err)
	}
	logs.EnableFuncCallDepth(false)
	logs.Async(1e3)
	color.Red(" \r\n  服务正在启动,请稍等......")
	for {
		go CheckAndRestart()

		if runtime.GOOS == "windows" {
			systray.Run(onReady, onExit)
		}
		time.Sleep(time.Second * 4)
	}
}
