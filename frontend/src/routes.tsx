import App from "./App";
import ErrorPage from "./components/ErrorPage";

const routes = [
    {
        path: "/",
        element: <App />,
        errorElement: <ErrorPage />,
        children: []
    }
]

export default routes;