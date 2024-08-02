#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "libir.h"

int main() {
    printf("Using irStandings lib from C:\n");

    const char *fn = "../model/fixtures/2024-1-285-results-redacted.json";
    const char *json = "{\"series_id\":285,\"session_id\":1,\"subsession_id\":2,\"track\":\"Lime Rock\",\"count_best_of\":10,\"car_class_id\":84,\"top_n\":5,\"positions\":[{\"cust_id\":123,\"laps_complete\":10,\"car_id\":76}]}";

    GoString filename = {p:fn, n:strlen(fn)};
    GoString livePositions = {p:json, n:strlen(json)};
    struct LiveStandings_return ret;

    printf("name: %s %ld\n", filename.p, filename.n);

    ret = LiveStandings(filename, livePositions);
    printf("msg = %s, val = %lld\n", ret.r0, ret.r1);
    free(ret.r0);
}
