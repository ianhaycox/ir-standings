import { useAppDispatch, useAppSelector } from "../../app/hooks"
import {
    selectStandings,
    selectStatus,
    set,
} from "./standingsSlice"

interface StandingsProps {
    carClassID: number;
}

export const Standings = ({ carClassID }: StandingsProps) => {
    const dispatch = useAppDispatch()
    const standings = useAppSelector(selectStandings)
    const status = useAppSelector(selectStatus)


    const rows = standings.map(row => {
        let change = "---"

        if (row.change < 0) {
            change = "v" + row.change*-1
        }

        if (row.change > 0) {
            change = "^" + row.change
        }

        return (
        <div key={row.cust_id} className="row irc-row text-start p-0">
            <div className="col-1 p-0">{row.predicted_position}</div>
            <div className="col-1 p-0">{row.car_number}</div>
            <div className="col p-0">{row.driver_name}</div>
            <div className="col-1 p-0">{row.current_position}</div>
            <div className="col-1 p-0 text-end">{row.predicted_points}</div>
            <div className="col-1 px-1 text-end">
                <div className={"float-end " + (row.change > 0 ? "irc-change-up" : row.change < 0 ? "irc-change-down" : "irc-change-none")}>{change}</div>
            </div>
        </div>
        )
    }
)

    return (
        <div id={`car-class-id-${carClassID}`} className="irc-standings">
            <div className="container py-2">
                <div key={-1} className="row irc-header text-start">
                    <div className="col-1 p-0">Pos</div>
                    <div className="col-1 p-0">No.</div>
                    <div className="col p-0">Driver</div>
                    <div className="col-1 p-0">Prev</div>
                    <div className="col-1 p-0 text-end">Pts</div>
                    <div className="col-1 p-0 text-end">+/-</div>
                </div>
                {rows}
                <div key={-2} className="row irc-footer text-center">
                    <div className="col p-0">Connected: SOF 2034 Laps:23</div>
                </div>
            </div>

            <button
                className="button"
                aria-label="Add"
                onClick={() => {
                    dispatch(set())
                }}
            >
                Add
            </button>


        </div>
    )
}