import useAuth from './hooks/useAuth'
import { Outlet } from 'react-router-dom'
import Navbar from './components/Navbar'
import Footer from './components/Footer'
import ToastProvider from './components/provider/ToastProvider'
import { Container } from '@mui/material'


function App() {
    // Use the useAuth hook to check authentication status
    // Ensures that the app waits for auth validation through redux before rendering routes
    // Maps to all routes through Outlet
    const loading = useAuth();

    if (loading) return <p>Loading...</p>;
    return (
        <div>
            <ToastProvider />
            <Navbar />
            <Container maxWidth="md" sx={{ mt: 4, mb: 4, minHeight: '80vh' }}>
                <Outlet />
            </Container>
            <Footer />
        </div>
    )
}

export default App
