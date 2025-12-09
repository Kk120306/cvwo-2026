import { configureStore } from '@reduxjs/toolkit';
import authReducer from '../store/slices/authSlice';

// Configuring the slices 
// TODO : Add more slices - Could do posts etc. 
export const store = configureStore({
  reducer: {
    auth: authReducer,
  },
});

// Types for TypeScript purposes
// https://redux.js.org/tutorials/typescript-quick-start - Read Define Root state and dispatch types
// Infer the `RootState`,  `AppDispatch`, and `AppStore` types from the store itself
export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
export type AppStore = typeof store



