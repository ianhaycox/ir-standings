import { Login } from "../../../wailsjs/go/main/App";

export const login = (username:string, password:string) => {
  return Login(username, password)
}