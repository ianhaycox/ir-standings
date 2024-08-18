import { Login } from "../../../wailsjs/go/main/App";

// A mock function to mimic making an async request for data
export const login = (username:string, password:string) => {
  return Login(username, password)
}