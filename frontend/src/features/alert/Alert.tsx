import { useAppSelector, useAppDispatch } from "../../app/hooks"
import { getAlert, clear } from "./alertSlice"

export const Alert = () => {
    const dispatch = useAppDispatch();
    const alert = useAppSelector(getAlert);

    if (alert.message == "") return null;

    return (
        <div>
            <div className="m-3">
                <div className={`alert alert-dismissible ${alert.type}`}>
                    {alert.message}
                    <button type="button" className="btn-close" onClick={() => dispatch(clear())}></button>
                </div>
            </div>
        </div>
    );
}
