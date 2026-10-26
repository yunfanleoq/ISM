package hisdbbackup

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStreamTarArgsExtractsOntoBackupRoot(t *testing.T) {
	const name = "ISM_TDengine_Backup_2026-10-22_09-00-00"
	hostOut := filepath.Join("data", "hisdbbackup", name)
	dockerArgs, hostArgs, err := StreamTarArgs("tdengine", "/tmp/"+name, hostOut)
	if err != nil {
		t.Fatal(err)
	}
	wantDocker := []string{"exec", "tdengine", "tar", "-C", "/tmp", "-cf", "-", name}
	if strings.Join(dockerArgs, "\x00") != strings.Join(wantDocker, "\x00") {
		t.Fatalf("docker args = %#v", dockerArgs)
	}
	wantHost := []string{"-xf", "-", "-C", filepath.Join("data", "hisdbbackup")}
	if strings.Join(hostArgs, "\x00") != strings.Join(wantHost, "\x00") {
		t.Fatalf("host args = %#v", hostArgs)
	}
	// 包内是 <目录名>/dbs.sql，解到父目录后 dbs.sql 在 hostOut 根上。
	if filepath.Join(hostArgs[len(hostArgs)-1], name, "dbs.sql") != filepath.Join(hostOut, "dbs.sql") {
		t.Fatal("dbs.sql would not land on the host backup dir root")
	}
	joined := strings.Join(append(append([]string{}, dockerArgs...), hostArgs...), " ")
	if strings.Contains(joined, "gzip") || strings.Contains(joined, "-czf") || strings.Contains(joined, " -z") {
		t.Fatalf("must not gzip an already snappy dump: %s", joined)
	}
	for _, a := range dockerArgs {
		if a == "-t" || a == "-T" {
			t.Fatalf("docker tar args must not carry taosdump thread flags: %#v", dockerArgs)
		}
	}
}

func TestStreamTarArgsUsesContainerName(t *testing.T) {
	const name = "ISM_TDengine_Backup_x"
	dockerArgs, _, err := StreamTarArgs("site_td", "/tmp/"+name+"/", filepath.Join("data", "hisdbbackup", name))
	if err != nil {
		t.Fatal(err)
	}
	if dockerArgs[1] != "site_td" {
		t.Fatalf("container = %s", dockerArgs[1])
	}
	if dockerArgs[len(dockerArgs)-1] != name {
		t.Fatalf("dump dir = %#v", dockerArgs)
	}
}

func TestStreamTarArgsRejectsBadPaths(t *testing.T) {
	const name = "ISM_TDengine_Backup_ok"
	host := filepath.Join("data", "hisdbbackup", name)
	cases := []struct {
		container, containerDir, hostOut string
	}{
		{"", "/tmp/" + name, host},
		{"-tdengine", "/tmp/" + name, host},
		{"td engine", "/tmp/" + name, host},
		{"td/engine", "/tmp/" + name, host},
		{"tdengine", "/tmp", host},
		{"tdengine", "/etc/passwd", host},
		{"tdengine", "/tmp/../etc", host},
		{"tdengine", "/tmp/a/b", host},
		{"tdengine", "/tmp/" + name, ""},
		{"tdengine", "/tmp/" + name, ".."},
		{"tdengine", "/tmp/" + name, filepath.Join("data", "hisdbbackup", "other")},
		{"tdengine", "/tmp/" + name, name},
	}
	for _, tc := range cases {
		if _, _, err := StreamTarArgs(tc.container, tc.containerDir, tc.hostOut); err == nil {
			t.Fatalf("accepted %#v", tc)
		}
	}
}

func TestContainerTarMissingNamesContainer(t *testing.T) {
	msg := ContainerTarMissing("tdengine")
	if !strings.Contains(msg, "tdengine") || !strings.Contains(msg, "没有 tar") {
		t.Fatal(msg)
	}
	if strings.Contains(msg, "docker cp") || strings.Contains(msg, "gzip") {
		t.Fatal(msg)
	}
}
