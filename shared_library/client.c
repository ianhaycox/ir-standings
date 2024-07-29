#include <stdio.h>
#include <stdlib.h>
#include "ir-standings.h"

int main() {
    printf("Using irStandings lib from C:\n");

    GoString xx = {"and goodnight"};
    struct LiveStandings_return ret;

    ret = LiveStandings(xx);
    printf("msg = %s, val = %lld\n", ret.r0, ret.r1);
    free(ret.r0);
}
