import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import type { RootState } from "../../app/store"
import { fetchLatestStandings } from "./telemetryAPI"

export interface SerializablePredictedStanding {
    driving: boolean;
    cust_id: number;
    driver_name: string;
    car_number?: string;
    current_position?: number;
    predicted_position: number;
    current_points: number;
    predicted_points: number;
    change: number;
    car_names: string[];
}

export interface SerializableStanding {
    sof_by_car_class: number;
    car_class_id: number;
    car_class_name: string;
    class_leader_laps_complete: number;
    items: SerializablePredictedStanding[];
}

export interface SerializablePredictedStandings {
    status: string;
    track_name: string;
    count_best_of: number;
    standings: { [key: number]: SerializableStanding }; // By Car Class ID
}

export interface TelemetryState {
    standings: SerializablePredictedStandings
    status: "idle" | "loading" | "failed"
}

const initialState: TelemetryState = {
    standings: <SerializablePredictedStandings>{},
    status: "idle",
}

export const telemetrySlice = createSlice({
    name: "telemetry",
    initialState,
    reducers: {
    },
    extraReducers: builder => {
        builder
            .addCase(getLatestStandings.pending, state => {
                state.status = "loading"
            })
            .addCase(getLatestStandings.fulfilled, (state, action) => {
                state.status = "idle"
                state.standings = action.payload
            })
            .addCase(getLatestStandings.rejected, state => {
                state.status = "failed"
            })
    },
})

export default telemetrySlice.reducer

export const selectLatestStandings = (state: RootState) => state.telemetry.standings
export const selectTelemetryStatus = (state: RootState) => state.telemetry.status

export const getLatestStandings = createAsyncThunk(
    "telemetry/fetchLatestStandings",
    async () => {

        const response = await fetchLatestStandings()

        return response
    },
)
