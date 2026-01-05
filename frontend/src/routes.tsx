import App from "./App"
import ErrorPage from "./components/ErrorPage"
import Home from "./pages/Home"
import Dashboard from "./pages/dashboard/Dashboard"
import PostPage from "./pages/posts/Post"
import CreatePost from "./pages/posts/CreatePost"
import AddTopic from "./pages/topics/CreateTopic"
import Login from "./pages/auth/Login"
import SignUp from "./pages/auth/SignUp"
import ProfilePage from "./pages/profile/Profile"
import { Navigate } from "react-router-dom"

const routes = [
    {
        path: "/",
        element: <App />,
        errorElement: <ErrorPage />,
        children: [
            // Home
            { index: true, element: <Home /> },

            // Auth
            { path: "login", element: <Login /> },
            { path: "signup", element: <SignUp /> },

            // Posts - Shows all posts with topics and serach
            {
                path: "posts",
                children: [
                    { index: true, element: <Dashboard /> },
                    { path: "create", element: <CreatePost /> },
                    { path: ":id", element: <PostPage /> },
                ],
            },

            // Topics
            {
                path: "topics",
                children: [
                    { path: "create", element: <AddTopic /> },
                ],
            },

            // Profile 
            {
                path: "profile",
                children: [
                    {
                        index: true,
                        element: <Navigate to="/posts" replace />
                    },
                    {
                        path: ":username",
                        element: <ProfilePage />
                    }
                ]
            },
        ],
    },
]

export default routes
