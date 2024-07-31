#pragma once

#include <string>
#include <vector>

struct LivePositions {
	int SeriesID;
	int SessionID;
	int SubsessionID;
	std::string Track;
	int CountBestOf;
	int CarClassID;
	int TopN;
	std::vector<LiveResults> Results;
};

struct LiveResults {
	int CustID;
	int FinishPositionInClass;
	int LapsComplete;
	int CarID;
};

struct LiveStandings {
	std::string DriverName;
	int         CurrentPosition;
	int         PredictedPosition;
	int         CurrentPoints;
	int         PredictedPoints;
	int         Change;
};


class Live
{
public:
    Live(const int selectedClassID) {
        m_selectedClassID = selectedClassID;
    }

    std::vector<LiveStandings> LatestStandings(LivePositions lp);

private:
    int m_selectedClassID;

};