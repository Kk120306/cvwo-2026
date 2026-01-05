import { AppBar, Toolbar, Typography, Button, Box, Avatar } from "@mui/material";
import { Link } from "react-router-dom";
import { useAppSelector, useAppDispatch } from "../hooks/reduxHooks";
import { logout } from '../store/slices/authSlice';
import { logoutAccount } from '../api/handleAuth';
import { useNavigate } from "react-router-dom";


// Simple navbar component - just for MVP
const Navbar = () => {
    const user = useAppSelector((state) => state.auth.user);
    const dispatch = useAppDispatch();
    const navigate = useNavigate();

    const handleLogout = async () => {
        dispatch(logout());
        await logoutAccount();
        navigate(0); // Refresh the page to update state
    }

    return (
        <AppBar
            position="static"
            elevation={0}
            sx={{
                background: "white",
                color: "black",
                borderBottom: "1px solid #eee",
            }}
        >
            <Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
                <Typography
                    variant="h6"
                    component={Link}
                    to="/"
                    sx={{ textDecoration: "none", fontWeight: 700 }}
                >
                    My Application
                </Typography>

                {!user ? (
                    <Box sx={{ display: "flex", gap: 2 }}>
                        <Button component={Link} to="/login">
                            Login
                        </Button>
                        <Button
                            component={Link}
                            to="/signup"
                            variant="contained"
                            sx={{ boxShadow: "none" }}
                        >
                            Sign Up
                        </Button>
                    </Box>
                ) : (
                    <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
                        <Button component={Link} to="/posts/create">
                            Create Post
                        </Button>
                        <Button onClick={handleLogout}>
                            Logout
                        </Button>
                        <Avatar
                            component={Link}
                            to={`/profile/${user.username}`}
                            src={user.avatarURL}
                            alt={user.username}
                            sx={{ width: 32, height: 32 }}
                        />
                    </Box>
                )}
            </Toolbar>
        </AppBar>
    );
};

export default Navbar;
