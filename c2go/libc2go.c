#include <stdio.h>
#include "libc2go.h"

box *createCInstance(int n)
{
	box *b = malloc(sizeof(box));
	b->content = n;
	return b;
}

char *getCByte()
{
	return "CByte";
}