import { useAppSelector, useAppDispatch } from "../../app/hooks"
import { selectLatestStandings, selectCarClassId, setCarClassId } from "../telemetry/telemetrySlice"
import { Config } from "../config/configSlice"

type Props = {
    config: Config;
}

export const Standings = (props: Props) => {
    const standings = useAppSelector(selectLatestStandings)
    const dispatch = useAppDispatch();

    const { show_topn } = props.config
    
    document.body.classList.add('irc-transparent')

    let carClassID = useAppSelector(selectCarClassId)

    if (standings.standings === undefined || !(carClassID in standings.standings)) return (
        <div className="irc-standings">
            <div className="spinner-border text-primary" role="status">
                <span className="visually-hidden">Loading...</span>
            </div>
                <div key={carClassID} id={`car-class-id-${carClassID}`} className="irc-standings small">
                    <div className='container py-2'>
                        {header(0)}
                        {dummyRows(10)}
                        {footer(null, 0, "", 0, "Waiting for iRacing", 0, 0)}
                    </div>
                </div>
        </div>
    );

    const tables: JSX.Element[] = []

    const rows: JSX.Element[] = []

    for (var i = 0; i < standings.standings[carClassID].items.length; i++) {
        const row = standings.standings[carClassID].items[i]

        if (i == show_topn) {
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

    if (rows.length !== 0) {
        const h = header(carClassID);
        const f = footer(dispatch, carClassID, standings.standings[carClassID].car_class_name, standings.standings[carClassID].sof_by_car_class,
            standings.track_name, standings.count_best_of, standings.standings[carClassID].class_leader_laps_complete);

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
        <div className="irc-standings">{tables}</div>
    )
}

const header = (carClassID: number) => {
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

const footer = (dispatch: any, carClassID: number, carClassName: string, sof: number, trackName: string, bestOf: number, lapsComplete: number) => {
    return (
        <div key={`footer-${carClassID}`} className="row irc-footer text-center pt-2">
            <div className="col-2 ps-0 text-center">
                <div className={`irc-toggle-class irc-car-number irc-car-number-${carClassID}`}
                    onClick={() => dispatch(setCarClassId(carClassID))}>
                    {carClassName}
                </div>
            </div>
            <div className="col-2 p-0 text-start">SOF :{sof}</div>
            <div className="col p-0 text-start">{trackName}</div>
            <div className="col-1 p-0 text-end">Best: {bestOf}</div>
            <div className="col-2 p-0 text-end">Laps: {lapsComplete}</div>
        </div>
    )
}

const dummyRows = (num:number) => {
    const rows: JSX.Element[] = []

    for (var i = 0; i < num; i++) {
        rows.push(
            <div key={i} className="row irc-row text-start p-0 irc-absent">
                <div className="col-1 p-0">{i+1}</div>
                <div className="col-1 p-0">
                    <div className={`float-start px-0 irc-box irc-car-number irc-car-number`}>
                        <i className="bi-dash-lg"></i>
                    </div>
                </div>
                <div className="col-4 p-0"></div>
                <div className="col-3 p-0 text-start"></div>
                <div className="col-1 p-0"></div>
                <div className="col-1 p-0 text-end"></div>
                <div className="col-1 px-1 text-end">

                    <div className="float-end irc-box ps-0 irc-change-none">
                        <i className="bi-dash-lg"></i>
                    </div>
                </div>
            </div>
        )
    }

    return rows
}