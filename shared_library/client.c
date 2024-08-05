#include <stdio.h>
#include "latest-standings.h"

int main(int argc, char *argv[]) {
    printf("Using irStandings lib from C:\n");

    const char* ret = GoLatestStandings("file", "{\"test\":1}");

    printf("%s", ret);

    // free(ret);
}
