import { Toaster } from 'react-hot-toast';

// Toaster Provider compoenent with config. 
// https://react-hot-toast.com/docs/styling - read documentation for more 
const ToastProvider = () => {
    return (
        <Toaster
            reverseOrder={false}
            toastOptions={{
                style: {
                    background: '#333',
                    color: '#fff',
                },
                success: {
                    duration: 2000,
                    style: {
                        background: '#4caf50',
                        color: '#fff',
                    },
                },
                error: {
                    duration: 2000,
                    style: {
                        background: '#f44336',
                        color: '#fff',
                    },
                },
            }}
            containerStyle={{
                top: '80px',
            }}
        />
    );
};

export default ToastProvider;
