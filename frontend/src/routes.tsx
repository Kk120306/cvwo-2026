import App from "./App"
import ErrorPage from "./components/ErrorPage"
import Home from "./pages/Home"
import Dashboard from "./pages/Dashboard"
import PostPage from "./pages/posts/PostPage"
import CreatePost from "./pages/posts/CreatePostPage"
import AddTopic from "./pages/topics/CreateTopic"
import Login from "./pages/auth/Login"
import SignUp from "./pages/auth/SignUp"

// Test Compoenents - Will be removed after development 
import TestRedux from "./components/test/TestRedux"
import AnotherComponent from "./components/test/AnotherComponent"
import DashTest from "./components/test/DashTest"

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

            // Testing 
            {
                path: "test",
                children: [
                    { index: true, element: <AnotherComponent /> },
                    { path: "redux", element: <TestRedux /> },
                    { path: "dash", element: <DashTest /> },
                ],
            },
        ],
    },
]

export default routes
