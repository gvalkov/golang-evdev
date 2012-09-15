// +build ignore

#include "linux/input.h"
#include "stdlib.h"
#include "stdio.h"

#define MAX_NAME_SIZE 256

int main (int argc, char** argv) {
    char arr[MAX_NAME_SIZE];

    printf("EVIOCGID   = 0x%x\n", EVIOCGID);
    printf("EVIOCGNAME = 0x%x\n", EVIOCGNAME(sizeof(arr)));
    printf("EVIOCGPHYS = 0x%x\n", EVIOCGPHYS(sizeof(arr)));
    printf("EVIOCGBIT  = 0x%x\n", EVIOCGBIT(0, EV_MAX));
    /* printf("EVIOCGABS  = 0x%x\n", EVIOCGABS(0, now what?!)); */
}
