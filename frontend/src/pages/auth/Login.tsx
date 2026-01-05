import { login } from '../../api/handleAuth';
import AuthForm from '../../components/authentication/AuthForm';
import { useAppDispatch, useAppSelector } from '../../hooks/reduxHooks';
import { setUser } from '../../store/slices/authSlice';
import { useNavigate, Link } from 'react-router-dom';
import { useEffect } from 'react';


// Login component, Uses Auth Form 
export default function Login() {
    const user = useAppSelector(state => state.auth.user);
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    // Checking if user is already logged in - if is redirect to posts
    useEffect(() => {
        if (user) navigate('/posts');
    }, [user, navigate]);

    if (user) return <p>Loading...</p>;

    // function to login the user with credentials 
    const handleLogin = async (username: string) => {
        const user = await login({ username });
        dispatch(setUser(user));
        navigate('/dashtest');
    };

    return (
        <AuthForm
            heading="Login"
            buttonLabel="Login"
            onSubmit={handleLogin}
            extraLink={<Link to="/signup">Dont have an account? Sign up</Link>}
        />
    );
}
