import iracingLogo from './assets/images/iRacing-Inline-Color-White.svg';
import seriesLogo from './assets/images/IMSA Vintage Series.png';
import './App.css';
import { Login } from "./components/login/Login";
import { Standings } from './components/standings/Standings';
import { Telemetry } from './components/telemetry/telemetry';
import { useAppSelector, useAppDispatch } from "./app/hooks"
import { isLoggedIn } from "./components/login/loginSlice"
import { Alert } from './components/alert/Alert';
import { fetchConfiguration } from './components/config/configAPI';
import { selectConfiguration } from './components/config/configSlice';

function App() {
    const loggedIn = useAppSelector(isLoggedIn)
    const config = useAppSelector(selectConfiguration)
    const dispatch = useAppDispatch();

    dispatch(fetchConfiguration)

    return (
        <div id="App" className="container-sm my-0">
            <div className="py-1">
                {loggedIn ? (
                    <div>
                        <Telemetry />
                        <Standings config={config} />
                    </div>
                ) : (
                    <div>
                        <div className="container banner">
                                <img className="float-start" src={iracingLogo} id="iracing-logo" alt="iRacing" />
                                <img className="float-end" src={seriesLogo} id="series-logo" alt="Kamel" />
                        </div>
                        <Alert />
                        <Login />
                    </div>
                )}
            </div>
        </div>
    )
}

export default App
