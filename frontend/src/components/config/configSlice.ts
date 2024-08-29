import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import type { RootState } from "../../app/store"
import { fetchConfiguration } from "./configAPI"
import { main } from "../../../wailsjs/go/models"

export interface ConfigState {
    config: main.Config,
    status: "idle" | "loading" | "failed"
}

const initialState: ConfigState = {
    config: {show_topn: 10},
    status: "idle",
}

export const configSlice = createSlice({
    name: "config",
    initialState,
    reducers: {
    },
    extraReducers: builder => {
        builder
            .addCase(getConfiguration.pending, state => {
                state.status = "loading"
            })
            .addCase(getConfiguration.fulfilled, (state, action) => {
                state.status = "idle"
                state.config = action.payload
            })
            .addCase(getConfiguration.rejected, state => {
                state.status = "failed"
            })
    },
})

export default configSlice.reducer

export const selectConfiguration = (state: RootState) => state.config.config

export const getConfiguration = createAsyncThunk(
    "config/fetchConfiguration",
    async () => {
        return await fetchConfiguration()
    },
)
