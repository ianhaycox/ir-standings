import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import type { RootState } from "../../app/store"
import { error, success, clear } from "../../features/alert/alertSlice"

import { login } from "./loginAPI"

export interface LoginState {
    ok: boolean
    status: "idle" | "loading" | "failed"
}

const initialState: LoginState = {
    ok: false,
    status: "idle",
}

export interface Credentials {
    username: string;
    password: string;
}

export const loginSlice = createSlice({
    name: "login",
    initialState,
    reducers: {},
    extraReducers: builder => {
        builder
            // Handle the action types defined by the `incrementAsync` thunk defined below.
            // This lets the slice reducer update the state with request status and results.
            .addCase(loginAsync.pending, state => {
                state.status = "loading"
                state.ok = false
            })
            .addCase(loginAsync.fulfilled, (state, action) => {
                state.status = "idle"
                state.ok = action.payload
            })
            .addCase(loginAsync.rejected, state => {
                state.status = "failed"
                state.ok = false
            })
    }
})

export default loginSlice.reducer

export const isLoggedIn = (state: RootState) => state.login.ok
export const selectStatus = (state: RootState) => state.login.status


export const loginAsync = createAsyncThunk(
    "login/authenticate",
    async (credentials: Credentials, thunkAPI) => {
        thunkAPI.dispatch(clear())

        const { username, password } = credentials;
        const response = await login(username, password)

        if (response) {
            thunkAPI.dispatch(success("Logged in"))
        } else {
            thunkAPI.dispatch(error("Invalid username or password"))
        }

        // The value we return becomes the `fulfilled` action payload
        return response
    },
)

