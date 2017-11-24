# cgoでポインタ渡しのテストコード

アドベントカレンダーGo3 11/4のコードです。

https://qiita.com/74th/private/0362bea2012ef253c539

CでGOのポインタを扱う方法と、GoでCのポインタを扱う方法を紹介する。

## GoのポインタをCで扱う場合

そのままGoのポインタを渡すコードを作成すると、ビルドエラーになる

```go:libgo2c.go
import "C"

//export createGoInstancePointer
func createGoInstancePointer(n C.int) *box {
	b := new(box)
	b.content = int(n)
	return b
}
```

```
go build -buildmode=c-shared -o libgo2c.so go2c.go
Go type not supported in export: struct
```

unsafe.Pointerにしてから渡すコードを作成すると、ビルドできるが、実行時エラーになる

```go:libgo2c.go
import "C"

//export createGoInstanceUnsafe
func createGoInstanceUnsafe(n C.int) unsafe.Pointer {
	b := new(box)
	b.content = int(n)
	return unsafe.Pointer(b)
}
```

```c:main.c
void *unsafeBox = createGoInstanceUnsafe();
```

```
go build -buildmode=c-shared -o libgo2c.so go2c.go
gcc -o main main.c -I./ -L./ -lgo2c
panic: runtime error: cgo result has Go pointer
```

unsafe.Pointer->uintptrに変換すると、渡すことができる

```go:libgo2c.go
import "C"

// ガベージコレクションで削除されないように、
// グローバル変数とつなげておくための補完先
var boxInstances map[uintptr]*box

//export createGoInstanceUintptr
func createGoInstanceUintptr(n C.int) uintptr {
	b := new(box)
	b.content = int(n)
	p := uintptr(unsafe.Pointer(b))

	if boxInstances == nil {
		boxInstances = make(map[uintptr]*box)
	}
	boxInstances[p] = b

	return p
}
```

もちろんuintptrを、unsafe.Pointer経由で元の方に復元できる

```go:libgo2c.go
//export getBoxContents
func getBoxContents(p uintptr) C.int {
	b := (*box)(unsafe.Pointer(p))
	return C.int(b.content)
}
```

## Goのバイト列をCで扱う場合

Goのバイト列の先頭のポインタをCに渡すことで、CでGoのバイト列を扱うことができる

```go:libc2go.go
import "C"

var buffer []byte

//export getGoBytes
func getGoBytes() uintptr {
	buffer = []byte("GoBytes\x00")
	return uintptr(unsafe.Pointer(&buffer[0]))
}
```

```c:main.c
char *gobyte = (char *)getGoBytes();
printf("go string:%s\n", gobyte);
```

## CのポインタをGOで扱う場合

簡単なCのライブラリのコードを示す

```c:libc2go.h
typedef struct{
	int content;
} box;

box *createCInstance(int n);
```

```c:libc2go.c
#include <stdio.h>
#include "libc2go.h"

box *createCInstance(int n)
{
	box *b = malloc(sizeof(box));
	b->content = n;
	return b;
}
```

Goのコード中にCのライブラリを参照するのLDFLAGSを追記する
GoはCのポインタを扱うことができる

```go:c2go.go
// #cgo LDFLAGS: -L./ -lc2go
// #include "libc2go.h"
import "C"

func main() {
	b := C.createCInstance(42)
	fmt.Printf("c pointer:%p\n", b)
	fmt.Printf("c struct contents:%d\n", int(b.content))
}
```

```sh
gcc -shared -fPIC -o libc2go.so libc2go.c
go run c2go.go
c pointer:0x4202900
c struct contents:42
```

## Cのバイト列をGoで扱う場合

C.GoString()を使うことで変換できるが、コピーが行われる。
直接参照するためには、unsafe.Pointer->固定長配列のポインタ->スライスの順に変換を行う

```go:c2go.go
//  char *getCByte()
cbyte := C.getCByte()

fmt.Printf("GoString():%s\n", C.GoString(cbyte))

// 1.unsafe.Pointerに変換
cbytePointer := unsafe.Pointer(cbyte)
// 2.固定長の配列のポインタに変換
cbyteArray := (*[1024]byte)(cbytePointer)
// 3.スライスに変換
cbyteSlice := cbyteArray[:7]
// cのバイト列をそのままスライスとして利用
fmt.Printf("cbyteSlice():%s\n", cbyteSlice)
```

```
gcc -shared -fPIC -o libc2go.so libc2go.c
go run c2go.go
GoString():CByte
cbyteSlice():CByte
```
