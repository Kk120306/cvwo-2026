import { useEffect, useState } from 'react';
import { setUser, logout } from '../store/slices/authSlice';
import { useAppDispatch } from './reduxHooks';

interface ValidateAuthResponse {
    user?: {
        id: string;
        username: string;
        role: string;
        avatarUrl: string;
        isAdmin: boolean;
    };
}

// Function to determine if the user has been authenticated
// Through the help of cookies passed to backend
export default function useAuth() {

    // Even though dispatch never changes, to satisfy react hooks linting rules
    const dispatch = useAppDispatch();
    // Loading state. Used globally in App to ensure routes are rendered only when state is ready
    const [loading, setLoading] = useState(true);


    // Calls the backend API to validate the user session
    // Sends Cookie, validates and if valid, sets the user in the redux store
    // If invalid, logs out the user by clearing the store
    // Use effect used to make sure it runs on apps first load, we add dependency on dispatch
    useEffect(() => {
        fetch(`${import.meta.env.VITE_BACKEND_HOST}/auth/validate`, { credentials: 'include' })
            .then(async (res) => {
                if (!res.ok) throw new Error('Failed to validate'); // If non-200 response
                const text = await res.text();
                if (!text) throw new Error('Empty response'); // If res body is empty, session invalid 
                return JSON.parse(text);
            })
            .then((data: ValidateAuthResponse) => {
                if (data.user) { // if JSON contains user 
                    dispatch(setUser(data.user));
                } else { // No user ie. invalid session 
                    dispatch(logout());
                }
            }) // Any other errors that may go undetected like backend error
            .catch(() => {
                dispatch(logout());
            })
            .finally(() => setLoading(false));
    }, [dispatch]);

    return loading;
}
