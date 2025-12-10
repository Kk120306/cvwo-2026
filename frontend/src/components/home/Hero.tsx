import { Box, Container, Stack, Typography, Button } from '@mui/material';
import { Link } from 'react-router-dom';

const Hero = () => {
    return (
        <Box
            sx={{
                minHeight: "100vh",
                display: "flex",
                alignItems: "center",
                background: "linear-gradient(135deg, #f5f5f5, #e9e9ff)",
            }}
        >
            <Container maxWidth="lg">
                <Stack spacing={3}>
                    <Typography
                        variant="h2"
                        fontWeight={700}
                        sx={{
                            lineHeight: 1.15,
                            letterSpacing: "-1px",
                        }}
                    >
                        Expand your knowledge
                    </Typography>

                    <Typography variant="h6" color="text.secondary" maxWidth="sm">
                        jfkdjsalfjkdsjaklfjdsajlfdjsakljfioewfjewipfjeoiwnfpkewnfnekwnf
                        fdnslkajfkoejnweakonfkeansklfnwepijfipoew
                    </Typography>
                    
                    <Stack direction="row" spacing={2}>
                        <Button
                            variant="contained"
                            size="large"
                            component={Link}
                            to="/signup"
                        >
                            Get Started
                        </Button>
                    </Stack>
                </Stack>
            </Container>
        </Box>
    );
};

export default Hero;
