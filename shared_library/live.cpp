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
#include "goir.h"
#include "picojson.h"
#include "live.h"
//#include "cgo.cpp"

std::vector<struct PredictedStanding> Live::LatestStandings(std::string fn, LiveResults lr) {
    picojson::object liveResults;

    liveResults["series_id"] = picojson::value(static_cast<double>(lr.seriesID));
    liveResults["season_id"] = picojson::value(static_cast<double>(lr.sessionID));
    liveResults["subsession_id"] = picojson::value(static_cast<double>(lr.subsessionID));
    liveResults["track"] = picojson::value(static_cast<std::string>(lr.track));
    liveResults["count_best_of"] = picojson::value(static_cast<double>(lr.countBestOf));
    liveResults["car_class_id"] = picojson::value(static_cast<double>(lr.carClassID));
    liveResults["top_n"] = picojson::value(static_cast<double>(lr.topN));

    picojson::array positions(lr.positions.size());
    for (int i = 0; i < lr.positions.size(); ++i) {
        picojson::object p;

        p["cust_id"] = picojson::value(static_cast<double>(lr.positions[i].custID));
        p["finish_position_in_class"] = picojson::value(static_cast<double>(lr.positions[i].finishPositionInClass));
        p["laps_complete"] = picojson::value(static_cast<double>(lr.positions[i].lapsComplete));
        p["car_id"] = picojson::value(static_cast<double>(lr.positions[i].carID));

        positions[i].set(p);
    }

    picojson::value v;
    v.set(picojson::array(positions));
    liveResults["results"] = v;

    const picojson::value value = picojson::value(liveResults);
    const std::string json = value.serialize(true);

    GoString goJSON = {json.c_str()};
    GoString filename = { fn.c_str() };
    struct LiveStandings_return ret;

    ret = LiveStandings(filename, goJSON);

    picojson::value result;

    std::string err = picojson::parse(result, ret.r0);
    free(ret.r0);
    if (!err.empty()) {
        printf("Live response is not valid JSON!\n%s\n", err.c_str());

        std::vector<struct PredictedStanding> r;
        return r;
    }

    std::vector<struct PredictedStanding> predictedStandings;
    picojson::array a = result.get<picojson::array>();

    for (int i = 0; i < (int)a.size(); ++i)
    {
        PredictedStanding ls = {
            a.at(i).get("driver_name").to_str(),
            a.at(i).get("car_number").to_str(),
            (int)a.at(i).get("current_position").get<double>(),
            (int)a.at(i).get("predicted_position").get<double>(),
            (int)a.at(i).get("current_points").get<double>(),
            (int)a.at(i).get("predicted_points").get<double>(),
            (int)a.at(i).get("change").get<double>(),
        };

        predictedStandings.push_back(ls);
    }

    return predictedStandings;
};
