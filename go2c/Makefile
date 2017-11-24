do:go2c
	./go2c

libgo2c.so:libgo2c.go
	go build -buildmode=c-shared -o libgo2c.so libgo2c.go

libgo2c.h:libgo2c.go
	go tool cgo -exportheader libgo2c.h libgo2c.go

go2c:go2c.c libgo2c.so libgo2c.h
	gcc -o go2c go2c.c -I./ -L./ -lgo2c

clean:
	rm -rf go2c libgo2c.so libgo2c.h