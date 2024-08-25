import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import type { RootState } from "../../app/store"
import { fetchLatestStandings } from "./telemetryAPI"
import type { PayloadAction } from "@reduxjs/toolkit"

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
    car_class_ids: number[];
    self_car_class_id: number;
    standings: { [key: number]: SerializableStanding }; // By Car Class ID
}

export interface TelemetryState {
    standings: SerializablePredictedStandings
    status: "idle" | "loading" | "failed"
    selectedCarClassID: number
}

const initialState: TelemetryState = {
    standings: <SerializablePredictedStandings>{},
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
