import App from "./App";
import ErrorPage from "./components/ErrorPage";
import TestRedux from "./components/test/TestRedux";
import AnotherComponent from "./components/test/AnotherComponent";
import Home from "./pages/Home";
import SignUp from "./components/authentication/SignUp";
import DashTest from "./components/test/DashTest";
import Login from "./components/authentication/Login";


// Define the routes for the application 
const routes = [
    {
        path: "/",
        element: <App />,
        errorElement: <ErrorPage />,
        children: [
            {
                index: true,
                element: <Home />,
            },
            {
                path: "test",
                element: <AnotherComponent />,
            },
            {
                path: "test2",
                element: <TestRedux />,
            },
            {
                path: "signup",
                element: <SignUp />,
            },
            {
                path : "dashtest",
                element: <DashTest />,
            }, 
            {
                path : "login",
                element: <Login />,
            }
        ]
    },
]

export default routes;
