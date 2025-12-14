import useAuth from './hooks/useAuth'
import { Outlet } from 'react-router-dom'
import Navbar from './components/Navbar'
import Footer from './components/Footer'
import ToastProvider from './components/ToastProvider'


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
            <Outlet />
            <Footer />
        </div>
    )
}

export default App
