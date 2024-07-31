/*
MIT License

Copyright (c) 2024 Ian Haycox

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

#pragma once

#include <string>
#include "ir-standings.h"
#include "picojson.h"

std::vector<LiveStandings> Live::LatestStandings(LivePositions lp) {
    std::vector<LiveStandings> result;

    std::string json = picojson::value(lp).serialize();

    GoString goJSON = {json};
    struct LiveStandings_return ret;

    ret = LiveStandings(goJSON);

    picojson::value result;

    std::string err = picojson::parse(result, ret.r0);
    if (!err.empty()) {
        std::vector<LiveStandings>{};
    }

//     printf("msg = %s, val = %lld\n", ret.r0, ret.r1);
    free(ret.r0);

    return result;
};
