import { AppBar, Toolbar, Typography, Button, Box } from "@mui/material";
import { Link } from "react-router-dom";

// Simple navbar component - just for MVP
const Navbar = () => {
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
            </Toolbar>
        </AppBar>
    );
};

export default Navbar;
