import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import type { RootState } from "../../app/store"
import { fetchPastResults } from "./standingsAPI"

export interface StandingsState {
  gotResults: boolean
  status: "idle" | "loading" | "failed"
}

const initialState: StandingsState = {
  gotResults: false,
  status: "idle",
}

export const standingsSlice = createSlice({
  name: "standings",
  initialState,
  reducers: {
  },
  extraReducers: builder => {
    builder
      .addCase(getPastResults.pending, state => {
        state.status = "loading"
      })
      .addCase(getPastResults.fulfilled, (state, action) => {
        state.status = "idle"
        state.gotResults = action.payload
      })
      .addCase(getPastResults.rejected, state => {
        state.status = "failed"
      })
  },
})

export default standingsSlice.reducer

export const selectGotResults = (state: RootState) => state.standings.gotResults
export const selectStatus = (state: RootState) => state.standings.status

export const getPastResults = createAsyncThunk(
  "standings/getPastResults",
  async (isLoggedIn: boolean) => {


    if (!isLoggedIn) {
      return false
    }

    const response = await fetchPastResults()

    return response
  },
)
