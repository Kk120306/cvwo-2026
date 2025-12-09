import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

// User Info that we want to store in slice. 
interface User {
    id: string;
    username: string;
    email: string;
    avatarURL: string;
    isAdmin: boolean;
}

// Type for the auth slice state
interface AuthState {
    user: User | null;
    isAuthenticated: boolean;
}

// Initially have no user with authentication as false 
const initialState: AuthState = {
    user: null,
    isAuthenticated: false,
};

// Create the auth slice
// Set User takes a User object and updates the state
// Logout clears the user info and sets isAuthenticated to false 
// https://redux.js.org/tutorials/typescript-quick-start - Refer to Application Usage
const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: { 
        setUser: (state, action: PayloadAction<User>) => {
            state.user = action.payload;
            state.isAuthenticated = true;
        },
        logout: (state) => {
            state.user = null;
            state.isAuthenticated = false;
        },
    },
});


export const { setUser, logout } = authSlice.actions;
export default authSlice.reducer;
