import { toast } from 'react-hot-toast';

// Typescript prop for Signup args
interface AuthProps {
    username: string;
}

// Function to handle user signup with the backend API 
// Calls post /auth/signup endpoint
export async function signup({ username }: AuthProps) {
    const res = await fetch(`${import.meta.env.VITE_BACKEND_HOST}/auth/signup`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ username }),
    })

    if (!res.ok) {
        const err = await res.json();
        toast.error(err.message || 'Signup failed');
        throw new Error(err.message || 'Signup failed');

    }

    return true;
}

// Function to handle user login with the backend API
// Calls post /auth/login endpoint
// Since the login returns the user data, we return data in order to set it in redux store 
export async function login({ username }: AuthProps) {
    const res = await fetch(`${import.meta.env.VITE_BACKEND_HOST}/auth/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ username }),
    });

    if (!res.ok) {
        const err = await res.json();
        toast.error(err.message || 'Login failed');
        throw new Error(err.message || 'Login failed, Please try again.');
    }

    const data = await res.json();
    toast.success('Login successful!');

    return data.user;
}


// Function to handle user logout with the backend API
// Calls post /auth/logout endpoint
// Calls to clear the cookie session
export async function logoutAccount() {
    await fetch(`${import.meta.env.VITE_BACKEND_HOST}/auth/logout`, {
        method: 'POST',
        credentials: 'include',
    });
    return true;
}