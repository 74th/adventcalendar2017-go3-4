package main

import (
	"C"
	"unsafe"
)

func main() {}

type box struct {
	content int
}

// //export createGoInstancePointer
func createGoInstancePointer(n C.int) *box {
	b := new(box)
	b.content = int(n)

	// そのままポインタを返す関数を
	// cgoとしてexportするとビルドエラーになる
	return b
}

//export createGoInstanceUnsafe
func createGoInstanceUnsafe(n C.int) unsafe.Pointer {
	b := new(box)
	b.content = int(n)

	// unsafe.Pointerにすると、ビルドエラーにはならないが、
	// "panic: runtime error: cgo result has Go pointer"の実行時エラーになる
	return unsafe.Pointer(b)
}

// インスタンス保管用
var boxInstances map[uintptr]*box

//export createGoInstanceUintptr
func createGoInstanceUintptr(n C.int) uintptr {
	b := new(box)
	b.content = int(n)

	// ガベージコレクションで削除されないように、
	// グローバル変数とつなげておくためのマップ
	if boxInstances == nil {
		boxInstances = make(map[uintptr]*box)
	}

	// unsafeにてポインタを、uintptrに変換する
	p := uintptr(unsafe.Pointer(b))

	// ガベコレを避けるためにmapに入れる
	boxInstances[p] = b

	return p
}

//export getBoxContents
func getBoxContents(p uintptr) C.int {
	b := (*box)(unsafe.Pointer(p))
	return C.int(b.content)
}

//export freeBox
func freeBox(p uintptr) {
	// インスタンス保管用マップから削除
	delete(boxInstances, p)
}

// Cの文字列を読む
//export readCBytes
func readCBytes(size C.int, p *C.char) {
	// 一旦固定長のbyteの配列にする
	// array := byte[1024](p)
	// それをスライスにする
	// b := array[:int(size)]
	// fmt.Printf("C:%s", b)
}

var buffer []byte

// getGoBytes GOの文字列を返す
//export getGoBytes
func getGoBytes() uintptr {

	// 最後にヌル文字を入れる
	buffer = []byte("GoBytes\x00")

	// Byte列の1文字目のポインタを返すことで、
	// 文字列とする
	return uintptr(unsafe.Pointer(&buffer[0]))
}
