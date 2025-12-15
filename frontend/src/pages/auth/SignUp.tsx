import { signup, login } from '../../api/handleAuth';
import AuthForm from '../../components/authentication/AuthForm';
import { useAppDispatch } from '../../hooks/reduxHooks';
import { setUser } from '../../store/slices/authSlice';
import { Link, useNavigate } from 'react-router-dom';
import { normalizeUser } from '../../helpers/normalizer';
import { useAppSelector } from '../../hooks/reduxHooks';
import { useEffect } from 'react';

// SignUp component, Uses Auth Form
export default function SignUp() {
    const user = useAppSelector(state => state.auth.user);
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    // Checking if user is already logged in - if is redirect to posts
    useEffect(() => {
        if (user) navigate('/posts');
    }, [user, navigate]);

    if (user) return <p>Loading...</p>;

    // Calls both login and signup. Sign up only creates the user, Login is then used to generate cookie and get user data
    const handleSignup = async (username: string) => {
        await signup({ username });
        const user = await login({ username });
        dispatch(setUser(normalizeUser(user)));
        navigate('/dashtest');
    };

    return (
        <AuthForm
            heading="Sign Up"
            buttonLabel="Create Account"
            onSubmit={handleSignup}
            extraLink={<Link to="/login">Already have an account? Log In</Link>}
        />
    );
}
