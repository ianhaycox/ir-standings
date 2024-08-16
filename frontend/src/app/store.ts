import type { Action, ThunkAction } from "@reduxjs/toolkit"
import { configureStore, combineReducers } from "@reduxjs/toolkit"
import standingsReducer from "../features/standings/standingsSlice"
import loginReducer from "../features/login/loginSlice"
import alertReducer from "../features/alert/alertSlice"

export const rootReducer = combineReducers({
    standings: standingsReducer,
    login: loginReducer,
    alert: alertReducer,
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
