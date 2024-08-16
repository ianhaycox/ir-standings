import { createSlice } from '@reduxjs/toolkit';
import type { RootState } from "../../app/store"
import type { PayloadAction } from "@reduxjs/toolkit"



export interface AlertState {
    type: string
    message: string
}

const initialState: AlertState = {
    type: "",
    message: "",
}

export const alertSlice = createSlice({
    name: "alert",
    initialState,
    reducers: {
        success: (state, action:PayloadAction<string>) => {
            state.type = 'alert-success';
            state.message =  action.payload;
          },
          error: (state, action:PayloadAction<string>) => {
            state.type = 'alert-danger';
            state.message = action.payload;
          },
          clear: state => {
            state.type = 'alert-success';
            state.message=  "";
          },
    },
})

export const { success, error, clear } = alertSlice.actions

export default alertSlice.reducer

export const getAlert = (state: RootState) => state.alert
