#include <string.h>
#include <stdio.h>
#include "windows.h"
#include "libgoir.h"

typedef struct LiveStandings_return LSR;
typedef LSR(CALLBACK* LPFNDLLFUNC1)(GoString, GoString);

static HINSTANCE hDLL = NULL;               // Handle to DLL
static LPFNDLLFUNC1 lpfnDllFunc1 = NULL;    // Function pointer

char* GoLatestStandings(const char* fn, const char* json) {

    GoString filename = { fn, (ptrdiff_t)strlen(fn) };
    GoString livePositions = {json, (ptrdiff_t)strlen(json)};

    LSR hrReturnVal;

    if (hDLL == NULL) {
        hDLL = LoadLibrary("libgoir");
    }

    if (NULL != hDLL)
    {
        if (NULL == lpfnDllFunc1) {
            lpfnDllFunc1 = (LPFNDLLFUNC1)GetProcAddress(hDLL, "LiveStandings");
        }

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
    }
    else
    {
        return "DLL not found";
    }

    return hrReturnVal.r0;
} 