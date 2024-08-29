import type { Action, ThunkAction } from "@reduxjs/toolkit"
import { configureStore, combineReducers } from "@reduxjs/toolkit"
import loginReducer from "../components/login/loginSlice"
import alertReducer from "../components/alert/alertSlice"
import telemetryReducer from "../components/telemetry/telemetrySlice"
import configSlice from "../components/config/configSlice"

export const rootReducer = combineReducers({
  config: configSlice,
  login: loginReducer,
  alert: alertReducer,
  telemetry: telemetryReducer,
})

export const store = configureStore({
  reducer: rootReducer,
  middleware: getDefaultMiddleware =>
    getDefaultMiddleware()
})

// Infer the type of `store`
export type AppStore = typeof store
export type RootState = ReturnType<typeof rootReducer>
// Infer the `AppDispatch` type from the store itself
export type AppDispatch = AppStore["dispatch"]
// Define a reusable type describing thunk functions
export type AppThunk<ThunkReturnType = void> = ThunkAction<
  ThunkReturnType,
  RootState,
  unknown,
  Action
>
