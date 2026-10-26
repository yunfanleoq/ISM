package controllers

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"ISMServer/utils/hisdbbackup"

	"github.com/beego/beego/v2/core/logs"
)

// 同时只允许一次历史库备份。第二次再点会和第一次抢同一套库，两边都停住。
var (
	hisBackupMu       sync.Mutex
	hisBackupRunState hisBackupRun
)

type hisBackupRun struct {
	Name       string
	Running    bool
	Prog       hisdbbackup.DumpProgress
	LastChange time.Time
}

func beginHisBackup(name string) bool {
	hisBackupMu.Lock()
	defer hisBackupMu.Unlock()
	if hisBackupRunState.Running {
		return false
	}
	hisBackupRunState = hisBackupRun{Name: name, Running: true, LastChange: time.Now()}
	return true
}

func finishHisBackup() {
	hisBackupMu.Lock()
	hisBackupRunState = hisBackupRun{}
	hisBackupMu.Unlock()
}

func noteHisBackupLine(line string) {
	hisBackupMu.Lock()
	defer hisBackupMu.Unlock()
	if !hisBackupRunState.Running {
		return
	}
	prev := hisBackupRunState.Prog
	next := hisdbbackup.ApplyDumpLine(prev, line)
	hisBackupRunState.Prog = next
	if next != prev {
		hisBackupRunState.LastChange = time.Now()
	}
}

func noteHisBackupPhase(phase string) {
	hisBackupMu.Lock()
	defer hisBackupMu.Unlock()
	if !hisBackupRunState.Running {
		return
	}
	hisBackupRunState.Prog.Phase = phase
	hisBackupRunState.LastChange = time.Now()
}

func hisBackupView(name string, valid bool) (status, text string, progress int) {
	hisBackupMu.Lock()
	run := hisBackupRunState
	hisBackupMu.Unlock()
	if run.Running && run.Name == name {
		text = run.Prog.Text()
		// 满 5 分钟张数没变只改提示。停进程看连续没有输出和 48 小时上限，不在这里杀。
		if run.Prog.Phase != "copy" && !run.LastChange.IsZero() && time.Since(run.LastChange) > 5*time.Minute {
			text += "（已超过 5 分钟没有新进度）"
		}
		return "running", text, run.Prog.Percent()
	}
	if valid {
		return "done", "备份完成", 100
	}
	return "incomplete", "不完整", 0
}

func currentDumpProgress() hisdbbackup.DumpProgress {
	hisBackupMu.Lock()
	defer hisBackupMu.Unlock()
	return hisBackupRunState.Prog
}

// runCmdTimeoutWatch 边跑边读输出。完整输出不进内存，只留最后一段给失败提示。
// 有新输出就继续等。连续 DumpStallLimit 没有输出，或墙上时钟到 DumpHardLimit，才停。
func runCmdTimeoutWatch(onLine func(string), name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}

	var stateMu sync.Mutex
	started := time.Now()
	lastOutput := started
	var stopReason string
	watchDone := make(chan struct{})
	go func() {
		defer close(watchDone)
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				stateMu.Lock()
				elapsed := time.Since(started)
				silent := time.Since(lastOutput)
				stateMu.Unlock()
				stop, reason := hisdbbackup.ShouldStopDump(elapsed, silent)
				if !stop {
					continue
				}
				stateMu.Lock()
				if stopReason == "" {
					stopReason = reason
				}
				stateMu.Unlock()
				cancel()
				return
			}
		}
	}()

	var tailMu sync.Mutex
	var tail strings.Builder
	consume := func(line string) {
		stateMu.Lock()
		lastOutput = time.Now()
		stateMu.Unlock()
		if onLine != nil {
			onLine(line)
		}
		tailMu.Lock()
		tail.WriteString(line)
		tail.WriteByte('\n')
		if tail.Len() > 4000 {
			s := tail.String()
			tail.Reset()
			if len(s) > 800 {
				s = s[len(s)-800:]
			}
			tail.WriteString(s)
		}
		tailMu.Unlock()
	}
	scan := func(r io.Reader) {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for sc.Scan() {
			consume(sc.Text())
		}
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		scan(stdout)
	}()
	go func() {
		defer wg.Done()
		scan(stderr)
	}()
	waitErr := cmd.Wait()
	cancel()
	wg.Wait()
	<-watchDone

	tailMu.Lock()
	out := []byte(tail.String())
	tailMu.Unlock()
	stateMu.Lock()
	reason := stopReason
	stateMu.Unlock()
	if waitErr != nil && reason != "" {
		return out, errors.New(reason)
	}
	return out, waitErr
}

var errSyncBusy = errors.New("sync busy")

// startDockerDirSync 在 taosdump 结束后把容器目录整份拷到宿主。
// 导出过程中默认不拷：每 30 秒对十几万小文件做 docker cp 会和 taosdump 抢同一批文件。
// 返回的函数做这一次最终拷贝，用流式 tar，不再逐文件 docker cp。
// 拷贝前进程被杀时，宿主还没有 dbs.sql，容器目录留给 settleDump，不在这里删。
func startDockerDirSync(container, containerDir, hostOut string) func() error {
	var mu sync.Mutex
	copying := false
	syncOnce := func() error {
		mu.Lock()
		if copying {
			mu.Unlock()
			return errSyncBusy
		}
		copying = true
		mu.Unlock()
		defer func() {
			mu.Lock()
			copying = false
			mu.Unlock()
		}()
		err := streamTarToHost(container, containerDir, hostOut)
		if err != nil {
			logs.Warn("HisDbBackUp sync to host failed: %v", err)
		}
		return err
	}
	stop := func() error {
		for {
			err := syncOnce()
			if err != errSyncBusy {
				return err
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
	interval := hisdbbackup.DumpHostSyncInterval
	if interval <= 0 {
		return stop
	}
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				_ = syncOnce()
			}
		}
	}()
	return func() error {
		close(done)
		return stop()
	}
}

// streamTarToHost 把容器 tar 的 stdout 接到宿主 tar 的 stdin。
// 不经过 shell，容器里也不落一份同等大小的 tar 文件。上限用 FinalHostCopyLimit。
// 容器没有 tar 时直接失败，不退回逐文件 docker cp。
func streamTarToHost(container, containerDir, hostOut string) error {
	dockerArgs, hostArgs, err := hisdbbackup.StreamTarArgs(container, containerDir, hostOut)
	if err != nil {
		return err
	}
	if err = ensureContainerTar(container); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), hisdbbackup.FinalHostCopyLimit)
	defer cancel()
	dockerCmd := exec.CommandContext(ctx, "docker", dockerArgs...)
	hostCmd := exec.CommandContext(ctx, "tar", hostArgs...)
	pr, pw := io.Pipe()
	dockerCmd.Stdout = pw
	hostCmd.Stdin = pr
	var dockerErrBuf, hostErrBuf bytes.Buffer
	dockerCmd.Stderr = &dockerErrBuf
	hostCmd.Stderr = &hostErrBuf

	// 先起容器 tar。它写完后关掉管道，宿主 tar 才能读到 EOF。
	// 宿主先退出时关掉读端，容器 tar 才不会堵在管道上。
	dockerDone := make(chan error, 1)
	go func() {
		err := dockerCmd.Run()
		_ = pw.Close()
		dockerDone <- err
	}()
	hostErr := hostCmd.Run()
	_ = pr.Close()
	dockerErr := <-dockerDone
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("拷到备份目录超时 %s", hisdbbackup.FinalHostCopyLimit)
	}
	if dockerErr == nil && hostErr == nil {
		return nil
	}
	return fmt.Errorf("流式 tar 拷贝失败: 容器 %v %s；宿主 %v %s",
		dockerErr, trimCmdDetail(dockerErrBuf.Bytes()), hostErr, trimCmdDetail(hostErrBuf.Bytes()))
}

func ensureContainerTar(container string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", "exec", container, "tar", "--help")
	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}
	detail := trimCmdDetail(out)
	blob := detail + " " + err.Error()
	if tarBinaryMissing(blob) {
		return errors.New(hisdbbackup.ContainerTarMissing(container))
	}
	if detail == "" {
		detail = err.Error()
	}
	return fmt.Errorf("容器 %s 无法执行 tar: %s", container, detail)
}

func tarBinaryMissing(detail string) bool {
	d := strings.ToLower(detail)
	if strings.Contains(d, "executable file not found") {
		return true
	}
	return strings.Contains(d, "exec: \"tar\"") && strings.Contains(d, "not found")
}

func trimCmdDetail(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 500 {
		s = s[len(s)-500:]
	}
	return s
}

func removeContainerDump(container, containerDir string) {
	if container == "" || containerDir == "" {
		return
	}
	_ = exec.Command("docker", "exec", container, "rm", "-rf", containerDir).Run()
}

// settleDump 决定这次能不能叫备份完成。数据表没导完、或文件没落到宿主机，都不算。
// prog 用拷贝前的进度。宿主已有 dbs.sql 就标不完整并保留；没有就留下容器目录。
func settleDump(hostOut, container, containerDir string, runErr error, failDetail string, prog hisdbbackup.DumpProgress) (keep bool, msg string, ok bool) {
	ready, why := hisdbbackup.BackupReady(hostOut, prog)
	if runErr != nil {
		if failDetail == "" {
			failDetail = runErr.Error()
		}
		if hisdbbackup.HasDbsSQL(hostOut) {
			reason := failDetail
			if prog.Total > 0 && !prog.DataComplete() {
				reason = hisdbbackup.InterruptReason(prog, failDetail)
			}
			_ = hisdbbackup.MarkIncomplete(hostOut, reason)
			removeContainerDump(container, containerDir)
			return true, "备份没有完成：" + reason, false
		}
		if containerDir != "" {
			return false, "备份没有完成：" + failDetail + "。容器里的目录还在 " + containerDir + "，没有删。", false
		}
		return false, "备份没有完成：" + failDetail, false
	}
	if !ready {
		_ = hisdbbackup.MarkIncomplete(hostOut, why)
		if hisdbbackup.HasDbsSQL(hostOut) {
			removeContainerDump(container, containerDir)
			return true, "备份没有完成：" + why, false
		}
		if containerDir != "" {
			return false, "备份没有完成：" + why + "。容器里的目录还在 " + containerDir + "，没有删。", false
		}
		return false, "备份没有完成：" + why, false
	}
	removeContainerDump(container, containerDir)
	if prog.Total > 0 {
		return true, fmt.Sprintf("备份完成，数据 %d/%d 已落到备份目录", prog.Done, prog.Total), true
	}
	return true, "备份完成，已落到备份目录", true
}
