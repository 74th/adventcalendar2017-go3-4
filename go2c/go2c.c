#include <stdio.h>
#include <unistd.h>
#include "libgo2c.h"

int main(int nArgs, char **args)
{
	// GOのポインタをunsafe.Pointer受け取る
	// "panic: runtime error: cgo result has Go pointer"のエラーになる
	// void *unsafeBox = createGoInstanceUnsafe();

	// GOのポインタを受け取る
	GoUintptr uintptrBox = createGoInstanceUintptr(42);
	printf("go pointer:%d\n", (int)uintptrBox);

	// GOのポインタ使って、GOで処理を実行する
	printf("go struct contents:%d\n", getBoxContents(uintptrBox));

	// 開放する
	freeBox(uintptrBox);

	// 1M個生成する
	int num = 1024 * 1024;
	GoUintptr *list = malloc(sizeof(GoUintptr) * num);
	for (int i = 0; i < num; i++)
	{
		list[i] = createGoInstanceUintptr(i);
	}

	// 63/64を開放する
	for (int i = 0; i < num; i++)
	{
		if(i%4!=0){
			freeBox(list[i]);
		}
	}

	// 1秒停止し、ガベージコレクションが走るのを待つ
	sleep(1);

	int missMatch = 0;
	// ガベージコレクションが走って違う値を返すか確認する
	for (int i = 0; i < num; i++)
	{
		if(i%64==0){
			int content = getBoxContents(list[i]);
			if (i != content)
			{
				missMatch = 1;
				printf("miss match:%d %d\n", i, content);
			}
		}
	}

	if(missMatch){
		// 同じ値を返したため、
		// 特にガベージコレクションで移動させられてはいない様子だった
	}

	// GOの文字列を読む
	char *gobyte = (char *)getGoBytes();
	printf("go string:%s\n", gobyte);

	return 0;
}