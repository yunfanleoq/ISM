package controllers

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestApplyHxRoomPatchKeepsRest(t *testing.T) {
	src := []byte("{\n  \"texts\": [\n    {\n      \"s\": \"5B2-T1\\n备用\",\n      \"x\": 3164.29,\n      \"y\": 139.42,\n      \"w\": 68.09\n    },\n    {\n      \"s\": \"keep\",\n      \"x\": 10,\n      \"y\": 20,\n      \"w\": 1\n    }\n  ],\n  \"knives\": [\n    {\n      \"cad\": \"T1\",\n      \"x\": 1.5,\n      \"y\": 2.5\n    }\n  ],\n  \"wires\": [\n    {\n      \"net\": \"w\",\n      \"pts\": [\n        [\n          10.5,\n          20.25\n        ],\n        [\n          10.5,\n          40\n        ]\n      ],\n      \"sw\": 2\n    }\n  ],\n  \"boxes\": []\n}\n")
	x := 3204.29
	y := 139.42
	out, err := applyHxRoomPatch(src, hxRoomPatch{
		Room:  "5B2",
		Texts: []hxNumPatch{{I: 0, X: &x, Y: &y}},
		Wires: []hxWirePatch{{I: 0, Pts: [][]float64{{50.5, 20.25}, {50.5, 40}}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(out, []byte("3164.29")) {
		t.Fatalf("old x still there: %s", out)
	}
	if !bytes.Contains(out, []byte("3204.29")) || !bytes.Contains(out, []byte("\"y\": 139.42")) {
		t.Fatalf("text not updated: %s", out)
	}
	if !bytes.Contains(out, []byte("\"s\": \"keep\"")) || !bytes.Contains(out, []byte("\"x\": 1.5")) {
		t.Fatalf("other objects changed: %s", out)
	}
	if !bytes.Contains(out, []byte("50.5")) || !bytes.Contains(out, []byte("\"sw\": 2")) {
		t.Fatalf("wire pts: %s", out)
	}
	back := 3164.29
	out2, err := applyHxRoomPatch(out, hxRoomPatch{
		Texts: []hxNumPatch{{I: 0, X: &back, Y: &y}},
		Wires: []hxWirePatch{{I: 0, Pts: [][]float64{{10.5, 20.25}, {10.5, 40}}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out2, src) {
		t.Fatalf("round trip\n got %s\nwant %s", out2, src)
	}
	if !json.Valid(out2) {
		t.Fatal("invalid json")
	}
}

func TestApplyHxRoomPatchReal5B2(t *testing.T) {
	path := "../ism-front-end-v2/public/hx-dcim/data/hx-power-rooms/5B2.json"
	src, err := os.ReadFile(path)
	if err != nil {
		path = "../../ism-front-end-v2/public/hx-dcim/data/hx-power-rooms/5B2.json"
		src, err = os.ReadFile(path)
	}
	if err != nil {
		t.Skip(err)
	}
	x := 3204.29
	y := 179.42
	out, err := applyHxRoomPatch(src, hxRoomPatch{Texts: []hxNumPatch{{I: 0, X: &x, Y: &y}}})
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(out, src) {
		t.Fatal("file unchanged")
	}
	backX := 3164.29
	backY := 139.42
	out2, err := applyHxRoomPatch(out, hxRoomPatch{Texts: []hxNumPatch{{I: 0, X: &backX, Y: &backY}}})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out2, src) {
		t.Fatalf("5B2 round trip changed %d bytes", len(out2)-len(src))
	}
}

func TestWriteHxRoomAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "room.json")
	orig := []byte("{\"room\":\"5A1\",\"texts\":[]}\n")
	if err := os.WriteFile(path, orig, 0644); err != nil {
		t.Fatal(err)
	}
	next := []byte("{\"room\":\"5A1\",\"texts\":[{\"s\":\"a\"}]}\n")
	if err := writeHxRoomAtomic(path, next); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, next) || !json.Valid(got) {
		t.Fatalf("written %s", got)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("temp left behind")
	}
	if err := writeHxRoomAtomic(path, []byte("{")); err == nil {
		t.Fatal("invalid json accepted")
	}
	got, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, next) {
		t.Fatalf("invalid json changed file: %s", got)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("temp left after invalid json")
	}

	ro := t.TempDir()
	roFile := filepath.Join(ro, "room.json")
	if err := os.WriteFile(roFile, orig, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(ro, 0555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(ro, 0755) })
	if err := writeHxRoomAtomic(roFile, next); err == nil {
		t.Fatal("unwritable dir accepted")
	}
	got, err = os.ReadFile(roFile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, orig) {
		t.Fatalf("unwritable dir changed file: %s", got)
	}
	if _, err := os.Stat(roFile + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("temp left after unwritable dir")
	}
}
