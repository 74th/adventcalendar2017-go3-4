package main

// #cgo LDFLAGS: -L./ -lc2go
// #include "libc2go.h"
import "C"

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/dustin/go-humanize"
)

func main() {

	// Cのポインタを受け取る
	b := C.createCInstance(42)
	fmt.Printf("c pointer:%p\n", b)

	// Cの構造体の中を直接参照できる
	fmt.Printf("c struct contents:%d\n", int(b.content))

	// Cのバイト列を得る
	cbyte := C.getCByte()

	// 通常はC.GoString()で変換する
	// この方法ではコピーが走る
	fmt.Printf("GoString():%s\n", C.GoString(cbyte))

	// 1.unsafe.Pointerに変換
	cbytePointer := unsafe.Pointer(cbyte)
	// 2.固定長の配列のポインタに変換
	cbyteArray := (*[1024]byte)(cbytePointer)
	// 3.固定長の配列からスライスを作る
	cbyteSlice := cbyteArray[:7]
	// cのバイト列をそのままスライスとして利用
	fmt.Printf("cbyteSlice():%s\n", cbyteSlice)

	count := int64(0)
	timer := time.NewTimer(time.Second * 10)
	fmt.Printf("CallEmtpy Test started. Wait 10 seconds.\n")

LOOP:
	for {
		select {
		case <-timer.C:
			break LOOP
		default:
		}
		C.callEmpty()
		count++
	}
	fmt.Printf("%s calls / sec\n", humanize.Comma(count/10))

}
