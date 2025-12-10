// Fully generated using Chatgpt Used purpley as a means to test if Redux is working
// And global state is being maintained. This falls under code review purposes
// Will be removed for final production build

import { useAppSelector } from '../../hooks/reduxHooks';
import { useNavigate } from 'react-router-dom';
import { useEffect } from 'react';

export default function DashTest() {
    const user = useAppSelector(state => state.auth.user);
    const navigate = useNavigate();

    useEffect(() => {
        if (!user) navigate('/signup'); // only redirects if validate finished and no user
    }, [user, navigate]);

    if (!user) return <p>Loading...</p>;

    return (
        <div style={{ padding: '2rem', maxWidth: '600px', margin: '2rem auto' }}>
            <h1>Welcome to DashTest!</h1>
            <p>Hello, {user.username}! You are authenticated.</p>
        </div>
    );
}
