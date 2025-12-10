import { login } from '../../api/handleAuth';
import AuthForm from './AuthForm';
import { useAppDispatch } from '../../hooks/reduxHooks';
import { setUser } from '../../store/slices/authSlice';
import { useNavigate } from 'react-router-dom';
import { normalizeUser } from '../../helpers/normalizeUser';
import { Link } from 'react-router-dom';

// Login component, Uses Auth Form 
export default function Login() {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    const handleLogin = async (username: string) => {
        const user = await login({ username });
        dispatch(setUser(normalizeUser(user)));
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
