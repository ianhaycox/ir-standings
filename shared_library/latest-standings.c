#include <string.h>
#include <stdio.h>
#include "libgoir.h"

const char* GoLatestStandings(const char* fn, const char* json) {

    GoString filename = {p:fn, n:strlen(fn)};
    GoString livePositions = {p:json, n:strlen(json)};
    struct LiveStandings_return ret;

    printf("name: %s %ld\n", filename.p, filename.n);

    ret = LiveStandings(filename, livePositions);
    printf("msg = %s, val = %lld\n", ret.r0, ret.r1);

    return ret.r0;
}