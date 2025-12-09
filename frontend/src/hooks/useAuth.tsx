import { useEffect } from 'react';
import { setUser, logout } from '../store/slices/authSlice';
import { useAppDispatch } from './reduxHooks';

// Function to determine if the user has been authenticated
// Through the help of cookies passed to backend
export default function useAuth() {

    // Even though dispatch never changes, to satisfy react hooks linting rules
    // Using type safe useDispatch 
    const dispatch = useAppDispatch();

    // Calls the backend API to validate the user session
    // Sends Cookie, validates and if valid, sets the user in the redux store
    // If invalid, logs out the user by clearing the store
    // Use effect used to make sure it runs on apps first load, we add dependency on dispatch
    useEffect(() => {
        fetch(`${import.meta.env.BACKEND_HOST}/validate`, { credentials: 'include' })
            .then(res => res.json())
            .then(data => dispatch(setUser(data.message)))
            .catch(() => dispatch(logout()));
    }, [dispatch]);
}
