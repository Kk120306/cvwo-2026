import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import routes from './routes';
import { Provider } from 'react-redux';
import { store } from './store';
import { ThemeProvider } from '@mui/material';
import theme from './utils/theme';

const router = createBrowserRouter(routes);

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <Provider store={store}>
            {/* https://mui.com/material-ui/customization/theming/ */}
            <ThemeProvider theme={theme} >
                <RouterProvider router={router} />
            </ThemeProvider>
        </Provider>
    </StrictMode>,
)