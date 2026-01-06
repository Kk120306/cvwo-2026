import { toast } from 'react-hot-toast';

// Typescript prop for Signup args
interface AuthProps {
    username: string;
}

const baseUrl = import.meta.env.VITE_BACKEND_HOST;

// Function to handle user signup with the backend API 
// Calls post /auth/signup endpoint
export async function signup({ username }: AuthProps) {
    const endPoint = `${baseUrl}/auth/signup`;

    const res = await fetch(endPoint, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ username }),
    })

    if (!res.ok) {
        const err = await res.json();
        toast.error(err.message || 'Signup failed, this username may already be taken.');
        throw new Error(err.message || 'Signup failed');

    }

    return true;
}

// Function to handle user login with the backend API
// Calls post /auth/login endpoint
// Since the login returns the user data, we return data in order to set it in redux store 
export async function login({ username }: AuthProps) {
    const endPoint = `${baseUrl}/auth/login`;
    const res = await fetch(endPoint, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ username }),
    });

    if (!res.ok) {
        const err = await res.json();
        toast.error(err.message || 'Login failed, Please make sure the username is correct.');
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
    const endPoint = `${baseUrl}/auth/logout`;
    await fetch(endPoint, {
        method: 'POST',
        credentials: 'include',
    });
    return true;
}