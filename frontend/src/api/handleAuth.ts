// Typescript prop for Signup args
interface SignupProps {
    username: string;
    email: string;
    password: string;
}

// Typescript prop for Login args
interface LoginProps {
    email: string;
    password: string;
}


// Function to handle user signup with the backend API 
// Calls post /auth/signup endpoint
export async function signup({ username, email, password }: SignupProps) {
    const res = await fetch(`${import.meta.env.VITE_BACKEND_HOST}/auth/signup`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ username, email, password }),
    })

    if (!res.ok) {
        const err = await res.json();
        throw new Error(err.message || 'Signup failed');
    }

    return true;
}

// Function to handle user login with the backend API
// Calls post /auth/login endpoint
// Since the login returns the user data, we return data in order to set it in redux store 
export async function login({ email, password }: LoginProps) {
    const res = await fetch(`${import.meta.env.VITE_BACKEND_HOST}/auth/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ email, password }),
    });

    if (!res.ok) {
        const err = await res.json();
        throw new Error(err.message || 'Login failed');
    }

    const data = await res.json();

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