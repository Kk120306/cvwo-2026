import { useState } from 'react';
import {
    Box,
    Typography,
    TextField,
    Button,
    Paper,
    Container
} from '@mui/material';

// Props for the AuthForm component
interface AuthFormProps {
    heading: string;
    buttonLabel: string;
    onSubmit: (username: string) => Promise<void>;
    extraLink: React.ReactNode;
}

// Authentication form component that is used for both login and signup since both only require username
export default function AuthForm({ heading, buttonLabel, onSubmit, extraLink }: AuthFormProps) {
    const [username, setUsername] = useState('');

    // handles form submission 
    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            await onSubmit(username);
        } catch {
            console.error('Error during form submission');
        }
    };

    return (
        <Container>
            <Box
                sx={{
                    marginTop: 8,
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                }}
            >
                <Paper
                    elevation={3}
                    sx={{
                        padding: 4,
                        width: '100%',
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                    }}
                >
                    <Typography component="h1" variant="h4" fontWeight={600} mb={3}>
                        {heading}
                    </Typography>

                    <Box component="form" onSubmit={handleSubmit} sx={{ width: '100%' }}>
                        <TextField
                            margin="normal"
                            required
                            fullWidth
                            id="name"
                            label="Username"
                            name="name"
                            autoComplete="username"
                            autoFocus
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                        />

                        <Button
                            type="submit"
                            fullWidth
                            variant="contained"
                            sx={{ mt: 3, mb: 2 }}
                        >
                            {buttonLabel}
                        </Button>

                        <Box sx={{ textAlign: 'center' }}>
                            {extraLink}
                        </Box>
                    </Box>
                </Paper>
            </Box>
        </Container>
    );
}