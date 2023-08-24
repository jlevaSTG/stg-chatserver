import ReactDOM from 'react-dom/client'
import './index.css'
import {
    QueryClient,
    QueryClientProvider,
} from '@tanstack/react-query'
import {MantineProvider} from '@mantine/core';
import {
    createBrowserRouter,
    RouterProvider,
} from "react-router-dom";
import "./index.css";
import Root from "./root/Root.tsx";
import Dashboard from "./pages/Dashboard.tsx";

const queryClient = new QueryClient()

const router = createBrowserRouter([
    {
        path: "/",
        element: <Root/>,
        children: [
            {
                path: "/",
                element: <Dashboard/>
            }
        ]
    },
]);

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
    <MantineProvider withGlobalStyles withNormalizeCSS>
        <QueryClientProvider client={queryClient}>

            <RouterProvider router={router}/>
        </QueryClientProvider>
    </MantineProvider>
)
