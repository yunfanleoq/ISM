package hisdbbackup

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// DumpHardLimit 是 taosdump 的墙上时钟上限。张数还在涨也不能超过它，避免真卡死不返回。
	// 11 万张子表的表结构加数据用 48 小时。固定 6 小时会把还在写的导出判成失败。
	DumpHardLimit = 48 * time.Hour
	// DumpStallLimit 连续这么久一条新输出都没有，才按中断结束。
	// 界面「超过 5 分钟没有新进度」只改提示，不杀进程。
	DumpStallLimit = 30 * time.Minute
	// DumpHostSyncInterval 是导出过程中把容器目录全量拷到宿主的间隔。
	// 0 表示 taosdump 结束前不拷。每 30 秒拷十几万小文件会和 taosdump 抢磁盘，
	// 表结构能从约 15 分钟拖到超过 6 小时。进度看 taosdump 输出，不看目录大小。
	DumpHostSyncInterval = time.Duration(0)
	// FinalHostCopyLimit 是结束后那一次流式 tar 拷到宿主的上限。这是唯一一次拷贝，要留够时间把已有文件拷全。
	FinalHostCopyLimit = 2 * time.Hour
	// DataDumpThreads 是 taosdump -T / --thread-num。只加速导数据，不加速表结构。3.3.6.13 默认 8。
	DataDumpThreads = 16
)

// historyDumpNotice 跟在超时原因后面，避免和业务库备份搞混。
func historyDumpNotice() string {
	return "（这是 TDengine 历史库备份，不是业务库备份；业务库请走数据库管理页）"
}

// ShouldStopDump 决定要不要停掉还在跑的 taosdump。
// elapsed 从启动算起，silent 从最后一行输出算起。有输出就继续等，直到墙上时钟上限。
func ShouldStopDump(elapsed, silent time.Duration) (bool, string) {
	if elapsed >= DumpHardLimit {
		return true, fmt.Sprintf("超时 %s%s", DumpHardLimit, historyDumpNotice())
	}
	if silent >= DumpStallLimit {
		return true, fmt.Sprintf("连续 %s 没有新输出，超时中断%s", DumpStallLimit, historyDumpNotice())
	}
	return false, ""
}

// DataDumpThreadFlag 返回 taosdump 的数据线程参数。必须是 -T，不能是 -t。
// 3.3.6.13 里 -t 是 WebSocket 超时，传了会直接 PROGRAM ERROR。
func DataDumpThreadFlag() (flag, value string) {
	return "-T", strconv.Itoa(DataDumpThreads)
}

// InterruptReason 是写进 backup.incomplete 的中断原因。
// 表结构阶段写明还没导完和已看到的张数。已经数到总表数时，不要写成「数据还没开始」。
// prog 必须是最终拷贝把阶段改成 copy 之前的进度，否则数据阶段的张数会丢。
func InterruptReason(p DumpProgress, failDetail string) string {
	failDetail = strings.TrimSpace(failDetail)
	if p.Phase == "data" && p.Total > 0 && !p.DataComplete() {
		return joinInterrupt(fmt.Sprintf("导出中断，数据 %d/%d", p.Done, p.Total), failDetail)
	}
	if p.Phase == "schema" || (p.Total > 0 && !p.DataComplete() && p.Phase != "data") {
		if p.Total > 0 {
			return joinInterrupt(fmt.Sprintf("导出中断，表结构还没导完（已看到 %d/%d 张）", p.Done, p.Total), failDetail)
		}
		return joinInterrupt("导出中断，表结构还没导完", failDetail)
	}
	return failDetail
}

func joinInterrupt(head, detail string) string {
	if detail == "" {
		return head
	}
	return head + "。" + detail
}

// taosdump 3.3 的两段输出。表结构行没有总数，总数来自 child tables。
// 数据行形如：ismhistorydb.templetehistorydatas [2800/119166] write avro 100% of hd_xxx
var (
	reDumpChildTables = regexp.MustCompile(`child tables \((\d+)\)`)
	reDumpSchema      = regexp.MustCompile(`dumping out schema:\s*(\d+)\s+from`)
	reDumpSchemaOK    = regexp.MustCompile(`total\s+(\d+)\s+table\(s\).*schema dumped`)
	reDumpData        = regexp.MustCompile(`\[(\d+)/(\d+)\]`)
)

// DumpProgress 是一次 taosdump 的进度。完成（100%）不在这里给出，
// 调用方必须等备份文件已经落到宿主机再标完成。
type DumpProgress struct {
	Phase      string
	Done       int
	Total      int
	ReachedEnd bool
}

// ApplyDumpLine 用 taosdump 的一行输出推进进度。认不出的行保持原样。
func ApplyDumpLine(p DumpProgress, line string) DumpProgress {
	if m := reDumpChildTables.FindStringSubmatch(line); m != nil {
		if n, err := strconv.Atoi(m[1]); err == nil && n > 0 {
			p.Total = n
		}
		if p.Phase == "" {
			p.Phase = "schema"
		}
		return p
	}
	if m := reDumpSchemaOK.FindStringSubmatch(line); m != nil {
		if n, err := strconv.Atoi(m[1]); err == nil && n > 0 {
			p.Total = n
		}
		if p.Phase == "" || p.Phase == "schema" {
			p.Phase = "schema"
		}
		return p
	}
	if m := reDumpSchema.FindStringSubmatch(line); m != nil {
		p.Phase = "schema"
		if n, err := strconv.Atoi(m[1]); err == nil {
			p.Done = n
		}
		return p
	}
	if m := reDumpData.FindStringSubmatch(line); m != nil {
		p.Phase = "data"
		if n, err := strconv.Atoi(m[1]); err == nil {
			p.Done = n
		}
		if n, err := strconv.Atoi(m[2]); err == nil && n > 0 {
			p.Total = n
		}
		if p.Total > 0 && p.Done >= p.Total {
			p.ReachedEnd = true
		}
		return p
	}
	return p
}

// DataComplete 数据阶段已经打到最后一张表。表结构导完不算。
func (p DumpProgress) DataComplete() bool {
	return p.ReachedEnd && p.Total > 0 && p.Done >= p.Total && p.Phase == "data"
}

// Percent 表结构最多 30，导数据最多 99。100 只留给已经落盘的备份。
func (p DumpProgress) Percent() int {
	switch p.Phase {
	case "copy":
		return 99
	case "done":
		return 100
	case "data":
		if p.Total <= 0 {
			return 30
		}
		v := 30 + p.Done*69/p.Total
		if v > 99 {
			return 99
		}
		return v
	case "schema":
		if p.Total <= 0 || p.Done <= 0 {
			return 0
		}
		v := p.Done * 30 / p.Total
		if v > 30 {
			return 30
		}
		if v == 0 && p.Done > 0 {
			return 1
		}
		return v
	default:
		return 0
	}
}

// Text 给界面直接显示。进行中一律以「正在备份中」开头。
func (p DumpProgress) Text() string {
	switch p.Phase {
	case "copy":
		return "正在备份中 正在拷到备份目录"
	case "schema":
		if p.Total > 0 && p.Done > 0 {
			return fmt.Sprintf("正在备份中 表结构 %d/%d", p.Done, p.Total)
		}
		return "正在备份中 表结构"
	case "data":
		if p.Total > 0 {
			return fmt.Sprintf("正在备份中 数据 %d/%d", p.Done, p.Total)
		}
		return "正在备份中 数据"
	default:
		return "正在备份中"
	}
}
