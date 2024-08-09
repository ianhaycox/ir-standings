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
#include "live.h"
#include "latest-standings.h"
#include "json.hpp"
using json = nlohmann::json;

std::vector<struct PredictedStanding> LatestStandings(std::string fn, LiveResults lr) {

    json liveResults;

    liveResults["series_id"] = lr.seriesID;
    liveResults["season_id"] = lr.sessionID;
    liveResults["subsession_id"] = lr.subsessionID;
    liveResults["track"] = lr.track;
    liveResults["count_best_of"] = lr.countBestOf;
    liveResults["car_class_id"] = lr.carClassID;
    liveResults["top_n"] = lr.topN;

    json positions = json::array();
    for (int i = 0; i < lr.positions.size(); ++i) {
        json p;

        p["cust_id"] = lr.positions[i].custID;
        p["finish_position_in_class"] = lr.positions[i].finishPositionInClass;
        p["laps_complete"] = lr.positions[i].lapsComplete;
        p["car_id"] = lr.positions[i].carID;

        liveResults["positions"].push_back(p);
    }

    std::string jsonStr = liveResults.dump();
    char* ret = GoLatestStandings(fn.c_str(), jsonStr.c_str());

    json response = json::parse(ret);
//    free(ret);


    std::vector<struct PredictedStanding> predictedStandings;

    for (int i = 0; i < (int)response.size(); ++i) {
        PredictedStanding ls = {
            response.at(i)["cust_id"],
            response.at(i)["driver_name"],
            response.at(i)["current_position"],
            response.at(i)["predicted_position"],
            response.at(i)["current_points"],
            response.at(i)["predicted_points"],
            response.at(i)["change"],
        };

        predictedStandings.push_back(ls);
    }

    return predictedStandings;
};
