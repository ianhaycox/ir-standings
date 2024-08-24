import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"
import type { RootState } from "../../app/store"
import { fetchPastResults } from "./standingsAPI"
import type { PayloadAction } from "@reduxjs/toolkit"


export interface StandingsState {
  gotResults: boolean
  status: "idle" | "loading" | "failed"
  selectedCarClassID: number
}

const initialState: StandingsState = {
  gotResults: false,
  status: "idle",
  selectedCarClassID: 0,
}

export const standingsSlice = createSlice({
  name: "standings",
  initialState,
  reducers: {
    setCarClassId: (state, action:PayloadAction<number>) => {
      if (action.payload === 83) {
        state.selectedCarClassID =  84;
      } else {
        state.selectedCarClassID =  83;
      }
    },
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
export const selectCarClassId = (state: RootState) => state.standings.selectedCarClassID
export const { setCarClassId } = standingsSlice.actions

export const getPastResults = createAsyncThunk(
  "standings/getPastResults",
  async (isLoggedIn:boolean) => {


    if (!isLoggedIn) {
      return false
    }

    const response = await fetchPastResults()

    return response
  },
)
