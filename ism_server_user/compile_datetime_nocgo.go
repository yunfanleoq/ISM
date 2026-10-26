//go:build !cgo

package main

import "time"

func fetchCompileDateTime() (date string, tim string) {
	now := time.Now()
	return now.Format("Jan 02 2006"), now.Format("15:04:05")
}
