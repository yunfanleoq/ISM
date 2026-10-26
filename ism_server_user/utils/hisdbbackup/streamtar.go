package hisdbbackup

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

// StreamTarArgs 拼 taosdump 结束后那一次拷贝的两条命令。
//
// 数据文件已经是 snappy avro，这里只打 tar，不再 gzip。再压几乎不减小体积，还占 CPU。
// 参数里因此是 -cf，不是 -czf，也不要再接 gzip。
//
// 容器侧是 tar -C /tmp -cf - <目录名>，包内带一层目录名。
// 宿主必须解到备份目录的父目录：tar -xf - -C <父目录>。
// 解压结果是 <父目录>/<目录名>/dbs.sql，也就是本次备份目录根上的 dbs.sql，HasDbsSQL 能找到。
// 如果把这一层解进备份目录自己里面，会变成 <备份目录>/<目录名>/dbs.sql，现有检查认不出。
//
// 不用 tar -C <目录> . 。成员「.」在部分 tar 上会让解压以非 0 退出，文件已经落上也会被当成失败。
//
// 这是流式参数：- 表示容器 tar 写 stdout、宿主 tar 读 stdin。调用方用管道接上，容器里不另存一份 tar。
// DockerArgs 不含 docker 本身，HostArgs 不含 tar 本身。
func StreamTarArgs(container, containerDir, hostOut string) (dockerArgs []string, hostArgs []string, err error) {
	container = strings.TrimSpace(container)
	if !validContainerName(container) {
		return nil, nil, fmt.Errorf("invalid container name")
	}
	dir := path.Clean(containerDir)
	if dir == "/tmp" || !strings.HasPrefix(dir, "/tmp/") {
		return nil, nil, fmt.Errorf("container dump dir must be a single directory under /tmp")
	}
	name := strings.TrimPrefix(dir, "/tmp/")
	if !isSinglePathElement(name) {
		return nil, nil, fmt.Errorf("container dump dir must be a single directory under /tmp")
	}
	host := filepath.Clean(strings.TrimSpace(hostOut))
	if !validHostBackupDir(host) {
		return nil, nil, fmt.Errorf("invalid host backup dir")
	}
	if filepath.Base(host) != name {
		return nil, nil, fmt.Errorf("host backup dir name must match container dump dir")
	}
	parent := filepath.Dir(host)
	dockerArgs = []string{"exec", container, "tar", "-C", "/tmp", "-cf", "-", name}
	hostArgs = []string{"-xf", "-", "-C", parent}
	return dockerArgs, hostArgs, nil
}

// ContainerTarMissing 是容器里没有 tar 时写给用户的失败原因。
// 缺 tar 就失败，不改用逐文件 docker cp。
func ContainerTarMissing(container string) string {
	return fmt.Sprintf("容器 %s 里没有 tar，无法把备份打成一个包拷出", container)
}

func validContainerName(name string) bool {
	if name == "" || strings.HasPrefix(name, "-") {
		return false
	}
	for _, r := range name {
		if r <= ' ' || r == '/' || r == '\\' || r == ':' {
			return false
		}
	}
	return true
}

func validHostBackupDir(host string) bool {
	if host == "" || host == "." || host == string(filepath.Separator) {
		return false
	}
	if host == ".." || strings.HasPrefix(host, ".."+string(filepath.Separator)) {
		return false
	}
	parent := filepath.Dir(host)
	if parent == "" || parent == "." || parent == string(filepath.Separator) {
		return false
	}
	return true
}
