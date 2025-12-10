import { signup, login } from '../../api/handleAuth';
import AuthForm from './AuthForm';
import { useAppDispatch } from '../../hooks/reduxHooks';
import { setUser } from '../../store/slices/authSlice';
import { Link, useNavigate } from 'react-router-dom';
import { normalizeUser } from '../../helpers/normalizeUser';

// SignUp component, Uses Auth Form
export default function SignUp() {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

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
