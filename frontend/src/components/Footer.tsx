import { Box, Typography, Container } from "@mui/material";

// Simple footer component - just for MVP
const Footer = () => {
    return (
        <Box
            component="footer"
            sx={{
                py: 3,
                background: "white",
                mt: "auto",
            }}
        >
            <Container maxWidth="lg">
                <Typography
                    variant="body2"
                    align="center"
                >
                    My Application. All rights reserved.
                </Typography>
            </Container>
        </Box>
    );
};

export default Footer;
