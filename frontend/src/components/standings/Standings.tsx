import { useAppSelector, useAppDispatch } from "../../app/hooks"
import { selectCarClassId, selectStatus, setCarClassId } from "./standingsSlice"
import { selectLatestStandings } from "../telemetry/telemetrySlice"

type Props = {
    topN: number;
}

interface IProps_Square {
    message: string,
    onClick: (e: React.MouseEvent<HTMLElement>) => void,
}

export const Standings = (props: Props) => {
    const standings = useAppSelector(selectLatestStandings)
    const status = useAppSelector(selectStatus)
    const dispatch = useAppDispatch();

    let { topN } = props

    let carClassID = useAppSelector(selectCarClassId)
    if (carClassID == 0) {
        carClassID = 83
   }


    if (status == "loading" || status == "failed" || standings.standings === undefined || !(carClassID in standings.standings)) return (
        <div className="text-center">
            <div className="spinner-border text-primary" role="status">
                <span className="visually-hidden">Loading...</span>
            </div>
        </div>
    );

    const tables: JSX.Element[] = []

    const header = () => {
        return (
            <div key={`header-${carClassID}`} className="row irc-header text-start pb-1">
                <div className="col-1 p-0">Pos</div>
                <div className="col-1 ps-1">No.</div>
                <div className="col p-0">Driver</div>
                <div className="col-1 p-0">Prev</div>
                <div className="col-1 p-0 text-end">Pts</div>
                <div className="col-1 pe-2 text-end">+/-</div>
            </div>
        )
    }

    const rows: JSX.Element[] = []

    for (var i = 0; i < standings.standings[carClassID].items.length; i++) {
        const row = standings.standings[carClassID].items[i]

        if (i == topN) {
            break;
        }

        let change = ""
        let changeIcon = "bi-dash-lg"

        if (row.change < 0) {
            change = "" + row.change * -1
            changeIcon = "float-start bi-caret-down-fill"
        }

        if (row.change > 0) {
            change = "" + row.change
            changeIcon = "float-start bi-caret-up-fill"
        }

        let rowClass = "row irc-row text-start p-0"
        if (!row.driving) {
            rowClass += " irc-absent"
        }

        let car_number = row.car_number
        let carNumberIcon = "bi-hash"
        if (row.car_number == "") {
            carNumberIcon = "bi-dash-lg"
        }

        rows.push(
            <div key={`${row.cust_id}-${carClassID}`} className={rowClass}>
                <div className="col-1 p-0">{row.predicted_position}</div>
                <div className="col-1 p-0">
                    <div className={`float-start px-0 irc-box irc-car-number irc-car-number-${carClassID}`}>
                        <i className={carNumberIcon}></i>
                        {car_number}
                    </div>
                </div>
                <div className="col-4 p-0">{row.driver_name}</div>
                <div className="col-3 p-0 text-start">{row.car_names}</div>
                <div className="col-1 p-0">{row.current_position}</div>
                <div className="col-1 p-0 text-end">{row.predicted_points}</div>
                <div className="col-1 px-1 text-end">

                    <div className={"float-end irc-box ps-0 " + (row.change > 0 ? "irc-change-up" : row.change < 0 ? "irc-change-down" : "irc-change-none")}>
                        <i className={changeIcon}></i>
                        {change}
                    </div>
                </div>
            </div>
        )
    }

    const footer = () => {
        return (
            <div key={`footer-${carClassID}`} className="row irc-footer text-center pt-2">
                <div className="col-2 ps-0 text-center">
                    <div className={`irc-toggle-class irc-car-number irc-car-number-${carClassID}`} onClick={() => dispatch(setCarClassId(carClassID))}>
                        {standings.standings[carClassID].car_class_name}
                    </div>
                </div>
                <div className="col-2 p-0 text-start">SOF :{standings.standings[carClassID].sof_by_car_class}</div>
                <div className="col p-0 text-center">{standings.track_name}</div>
                <div className="col-1 p-0 text-end">Best: {standings.count_best_of}</div>
                <div className="col-2 p-0 text-end">Laps: {standings.standings[carClassID].class_leader_laps_complete}</div>
            </div>
        )
    }

    if (rows.length !== 0) {
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
    }

    return (
        <div>{tables}</div>
    )
}
