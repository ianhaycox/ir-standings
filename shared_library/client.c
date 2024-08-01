#include <stdio.h>
#include <stdlib.h>
#include "libir.h"

int main() {
    printf("Using irStandings lib from C:\n");

    GoString filename = {"foo.json"};
    GoString livePositions = {"{}"};
    struct LiveStandings_return ret;

    ret = LiveStandings(filename, livePositions);
    printf("msg = %s, val = %lld\n", ret.r0, ret.r1);
    free(ret.r0);
}
