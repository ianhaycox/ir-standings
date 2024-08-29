import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import type { RootState } from "../../app/store"
import { fetchLatestStandings } from "./telemetryAPI"
import type { PayloadAction } from "@reduxjs/toolkit"
import {live} from "../../../wailsjs/go/models"

export interface TelemetryState {
    standings: live.PredictedStandings
    status: "idle" | "loading" | "failed"
    selectedCarClassID: number
}

const initialState: TelemetryState = {
    standings: <live.PredictedStandings>{},
    status: "idle",
    selectedCarClassID: 0,
}

export const telemetrySlice = createSlice({
    name: "telemetry",
    initialState,
    reducers: {
        setCarClassId: (state, action: PayloadAction<number>) => {

            let selectedCarClassID: number = 0;
            const len = state.standings.car_class_ids.length;

            // Cycle through car class IDs
            for (let i = 0; i < len; i++) {
              if (action.payload == state.standings.car_class_ids[i]) {
                let next = i

                if (next < len - 1) {
                  next += 1
                } else {
                  next = 0
                }

                selectedCarClassID = state.standings.car_class_ids[next]

                break
              }
            }

            if (selectedCarClassID == 0) {
                selectedCarClassID = state.standings.self_car_class_id
            }

            // Just choose one
            if (selectedCarClassID == 0 && len > 0) {
              selectedCarClassID = state.standings.car_class_ids[0]
            }

            state.selectedCarClassID = selectedCarClassID
          },
    },
    extraReducers: builder => {
        builder
            .addCase(getLatestStandings.pending, state => {
                state.status = "loading"
            })
            .addCase(getLatestStandings.fulfilled, (state, action) => {
                state.status = "idle"
                state.standings = action.payload
                if (state.selectedCarClassID == 0) {
                    state.selectedCarClassID = action.payload.self_car_class_id
                }
            })
            .addCase(getLatestStandings.rejected, state => {
                state.status = "failed"
            })
    },
})

export default telemetrySlice.reducer

export const selectLatestStandings = (state: RootState) => state.telemetry.standings
export const selectTelemetryStatus = (state: RootState) => state.telemetry.status
export const { setCarClassId } = telemetrySlice.actions
export const selectCarClassId = (state: RootState) => state.telemetry.selectedCarClassID

export const getLatestStandings = createAsyncThunk(
    "telemetry/fetchLatestStandings",
    async () => {

        const response = await fetchLatestStandings()

        return response
    },
)
