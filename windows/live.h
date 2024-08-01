#pragma once

#include <string>
#include <vector>

struct CurrentPosition {
	int custID;                       // "cust_id"
	int finishPositionInClass;        // "finish_position_in_class"
	int lapsComplete;                 // "laps_complete"
	int carID;                        // "car_id"
};

struct LiveResults {
	int seriesID;                              // "series_id"
	int sessionID;                             // "session_id"
	int subsessionID;                          // "subsession_id"
	std::string track;                         // "track"
	int countBestOf;                           // "count_best_of"
	int carClassID;                            // "car_class_id"
	int topN;                                  // "top_n"
	std::vector<CurrentPosition> positions;   // "positions"
};

struct PredictedStanding {
	std::string driverName;           // "driver_name"
	std::string carNumber;            // "car_number"
	int         currentPosition;      // "current_position"
	int         predictedPosition;    // "predicted_position"
	int         currentPoints;        // "current_points"
	int         predictedPoints;      // "predicted_points"
	int         change;               // "change"
};


class Live
{
public:
    Live(const int selectedClassID) {
        m_selectedClassID = selectedClassID;
    }

    std::vector<struct PredictedStanding> LatestStandings(LiveResults lr);

private:
    int m_selectedClassID;

};