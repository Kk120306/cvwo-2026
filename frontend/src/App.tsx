import './App.css'
import useAuth from './hooks/useAuth'
import { Outlet } from 'react-router-dom'


function App() {
    useAuth();
    return (
        <div>
            <Outlet />
        </div>
    )
}

export default App
