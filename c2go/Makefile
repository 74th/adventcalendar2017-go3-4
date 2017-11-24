do:c2go.go libc2go.so
	go run c2go.go

c2go:c2go.go libc2go.so libc2go.h
	go build c2go.go

libc2go.so:libc2go.c libc2go.h
	gcc -shared -fPIC -o libc2go.so libc2go.c

clean:
	rm -rf c2go libc2go.so