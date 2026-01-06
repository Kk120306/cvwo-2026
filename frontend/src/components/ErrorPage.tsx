import { Box, Typography, Button, Container } from '@mui/material';
import { Link } from 'react-router-dom';

const ErrorPage = () => {
    return (
        <Container maxWidth="sm">
            <Box
                sx={{
                    minHeight: '70vh',
                    display: 'flex',
                    flexDirection: 'column',
                    justifyContent: 'center',
                    alignItems: 'center',
                    textAlign: 'center',
                }}
            >
                <Typography variant="h3" fontWeight={700} gutterBottom>
                    Oops!
                </Typography>

                <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
                    Sorry, an unexpected error has occurred.
                </Typography>

                <Button
                    component={Link}
                    to="/"
                    variant="contained"
                    size="large"
                >
                    Go back home
                </Button>
            </Box>
        </Container>
    );
};

export default ErrorPage;
