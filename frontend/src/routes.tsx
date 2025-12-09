import App from "./App";
import ErrorPage from "./components/ErrorPage";
import TestRedux from "./components/test/TestRedux";
import AnotherComponent from "./components/test/AnotherComponent";
import Home from "./components/Home";

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
            }
        ]
    },
]

export default routes;
