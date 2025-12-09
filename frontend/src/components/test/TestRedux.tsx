// Fully generated using Chatgpt Used purpley as a means to test if Redux is working
// And global state is being maintained. This falls under code review purposes
// Will be removed for final production build

import type { RootState } from '../../store';
import { setUser, logout } from '../../store/slices/authSlice';
import { useAppDispatch, useAppSelector } from '../../hooks/reduxHooks';
import { useEffect } from 'react';

export default function TestRedux() {
    const user = useAppSelector(state => state.auth.user);
    const dispatch = useAppDispatch();
    const isAuthenticated = useAppSelector((state: RootState) => state.auth.isAuthenticated);


    useEffect(() => {
        console.log('Current user:', user);
    }, [user]);



    // Simulate a user login
    const handleLogin = () => {
        dispatch(setUser({
            id: '1',
            username: 'Kai',
            email: 'kai@example.com',
            avatarURL: '',
            isAdmin: true,
        }));
    };

    // Simulate logout
    const handleLogout = () => dispatch(logout());

    return (
        <div style={{ padding: '2rem', border: '2px solid #ccc', maxWidth: '400px', margin: '2rem auto' }}>
            <h2>Redux Test Component</h2>
            <p><strong>Authenticated:</strong> {isAuthenticated ? 'Yes' : 'No'}</p>
            <p><strong>User:</strong> {user ? user.username : 'None'}</p>
            <button onClick={handleLogin} style={{ marginRight: '1rem' }}>Login (set user)</button>
            <button onClick={handleLogout}>Logout</button>
        </div>
    );
}
