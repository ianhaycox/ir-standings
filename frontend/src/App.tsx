import logo from './assets/images/logo-universal.png';
import './App.css';
import { Login } from "./features/login/Login";
import { Standings } from './features/standings/Standings';
import { useAppSelector } from "./app/hooks"
import { isLoggedIn } from "./features/login/loginSlice"
import { Alert } from './features/alert/Alert';

function App() {
    const loggedIn = useAppSelector(isLoggedIn)

    return (
        <div id="App" className="container-sm my-2">
            <div className="py-3">
                {loggedIn ? (
                    <Standings carClassID={83} />
                ) : (
                    <div>
                        <Alert />
                        <Login />
                    </div>
                )}
            </div>
        </div>
    )
}

export default App
