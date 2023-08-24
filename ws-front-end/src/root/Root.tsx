import { AppShell} from '@mantine/core';
import SideNav from "./SideNav.tsx";
import {Outlet} from "react-router-dom";


function Root() {
    return (
        <AppShell
            className="h-screen"
            padding=""
            navbar={<SideNav />}
            // header={<Header height={60} p="xs">{/* Header content */}</Header>}
            styles={(theme) => ({
                main: { backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0] },
            })}
        >
           <Outlet />
        </AppShell>
    );
}

export default Root;