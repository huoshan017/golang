package main

/*
#cgo LDFLAGS: -L ./ libuv.lib
#include <stdio.h>
#include <stdlib.h>
#include "libuv/uv.h"
void print(char* str) {
	printf("%s\n", str);
}
*/
import "C"

import "unsafe"

func main() {
	s := "hello"
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.print(cs)
	var buffer *C.char
	var size *C.size_t
	C.uv_cwd(buffer, size)
}
