#include <string.h>
#include <stdio.h>
#include "windows.h"
#include "libgoir.h"

char* GoLatestStandings(const char* fn, const char* json) {

    GoString filename = { fn, (ptrdiff_t)strlen(fn) };
    GoString livePositions = {json, (ptrdiff_t)strlen(json)};
    typedef struct LiveStandings_return LSR;

    printf("name: %s %lld\n", filename.p, filename.n);

    typedef LSR(CALLBACK* LPFNDLLFUNC1)(GoString, GoString);

    HINSTANCE hDLL;               // Handle to DLL
    LPFNDLLFUNC1 lpfnDllFunc1;    // Function pointer
    LSR hrReturnVal;

    hDLL = LoadLibrary("libgoir");
    if (NULL != hDLL)
    {

        lpfnDllFunc1 = (LPFNDLLFUNC1)GetProcAddress(hDLL, "LiveStandings");
        if (NULL != lpfnDllFunc1)
        {
            // call the function
            hrReturnVal = lpfnDllFunc1(filename, livePositions);
        }
        else
        {
            // report the error
            return "Error";
        }
        FreeLibrary(hDLL);
    }
    else
    {
        return "DLL not found";
    }


//    printf("msg = %s, val = %lld\n", ret.r0, ret.r1);

    return hrReturnVal.r0;
} 