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
	int seriesID = 0;                         // "series_id"
	int sessionID = 0;                        // "session_id"
	int subsessionID= 0;                      // "subsession_id"
	std::string track= "";                    // "track"
	int countBestOf = 10;                     // "count_best_of"
	int carClassID = 0;                       // "car_class_id"
	int topN = 5;                             // "top_n"
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

std::vector<struct PredictedStanding> LatestStandings(std::string fn, LiveResults lr);

