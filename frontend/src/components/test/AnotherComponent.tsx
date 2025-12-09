// Fully generated using Chatgpt Used purpley as a means to test if Redux is working
// And global state is being maintained. This falls under code review purposes
// Will be removed for final production build

import { useAppDispatch, useAppSelector } from '../../hooks/reduxHooks';
import { useEffect } from 'react';

export default function AnotherComponent() {
    const user = useAppSelector(state => state.auth.user);
    const dispatch = useAppDispatch();

    useEffect(() => {
        console.log('Current user:', user);
    }, [user]);

    return (
        <div>
            <h3>Another Component</h3>
            <p>User from Redux: {user ? user.username : 'None'}</p>
        </div>
    );
}
