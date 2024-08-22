import { useAppDispatch, useAppSelector } from "../../app/hooks"
import {
    selectGotResults,
    selectStatus,
    getPastResults,
} from "./standingsSlice"
import { selectLatestStandings, selectTelemetryStatus } from "../telemetry/telemetrySlice"

export const Standings = () => {
    const standings = useAppSelector(selectLatestStandings)
    const status = useAppSelector(selectStatus)
    const telemetryStatus = useAppSelector(selectTelemetryStatus)



    if (status == "loading" || status == "failed") return (<p>Loading...</p>);

    const tables: JSX.Element[] = []

    tables.push(<p>{telemetryStatus}</p>)

    for (const carClassID in standings.standings) {

        const header = () => {
            return (
                <div key={`header-${carClassID}`} className="row irc-header text-start">
                    <div className="col-1 p-0">Pos</div>
                    <div className="col-1 ps-1">No.</div>
                    <div className="col p-0">Driver</div>
                    <div className="col-1 p-0">Prev</div>
                    <div className="col-1 p-0 text-end">Pts</div>
                    <div className="col-1 pe-2 text-end">+/-</div>
                </div>
            )
        }

        const rows = standings.standings[carClassID].items.map(row => {
            let change = "---"

            if (row.change < 0) {
                change = "v" + row.change * -1
            }

            if (row.change > 0) {
                change = "^" + row.change
            }

            let rowClass = "row irc-row text-start p-0"
            if (!row.driving) {
                rowClass += " irc-absent"
            }

            let car_number = row.car_number
            if (row.car_number == "") {
                car_number = "---"
            }

            return (
                <div key={`${row.cust_id}-${carClassID}`} className={rowClass}>
                    <div className="col-1 p-0">{row.predicted_position}</div>
                    <div className="col-1 p-0">
                        <div className={`float-start irc-box irc-car-number-${carClassID}`}>#{car_number}</div>
                    </div>
                    <div className="col-4 p-0">{row.driver_name}</div>
                    <div className="col-3 p-0 text-start">{row.car_names}</div>
                    <div className="col-1 p-0">{row.current_position}</div>
                    <div className="col-1 p-0 text-end">{row.predicted_points}</div>
                    <div className="col-1 px-1 text-end">
                        <div className={"float-end irc-box " + (row.change > 0 ? "irc-change-up" : row.change < 0 ? "irc-change-down" : "irc-change-none")}>{change}</div>
                    </div>
                </div>
            )
        })

        const footer = () => {
            return (
                <div key={`footer-${carClassID}`} className="row irc-footer text-center">
                    <div className="col-2 p-0 text-start">{standings.standings[carClassID].car_class_name}</div>
                    <div className="col-2 p-0 text-start">SOF:{standings.standings[carClassID].sof_by_car_class}</div>
                    <div className="col p-0 text-center">{standings.track_name}</div>
                    <div className="col-2 p-0 text-end">Laps:{standings.standings[carClassID].class_leader_laps_complete}</div>
                </div>
            )
        }

        const h = header();
        const f = footer();

        const table = [
            <div key={carClassID} id={`car-class-id-${carClassID}`} className="irc-standings small">
                <div className='container py-2'>
                    {h}
                    {rows}
                    {f}
                </div>
            </div>
        ]

        tables.push(...table)
        table.push(<hr />)
    }

    return (
        <div>{tables}</div>
    )
}
