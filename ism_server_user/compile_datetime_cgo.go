//go:build cgo

package main

/*
#include<stdint.h>
#include<string.h>
void getCompileDateTime(uint8_t dt[12], uint8_t tm[9]) {
  strcpy((char*)dt, __DATE__);
  strcpy((char*)tm, __TIME__);
}
*/
import "C"

import "unsafe"

func fetchCompileDateTime() (date string, tim string) {
	dt := make([]byte, 12)
	tm := make([]byte, 10)
	C.getCompileDateTime(
		(*C.uint8_t)(unsafe.Pointer(&dt[0])),
		(*C.uint8_t)(unsafe.Pointer(&tm[0])),
	)
	return string(dt), string(tm)
}
