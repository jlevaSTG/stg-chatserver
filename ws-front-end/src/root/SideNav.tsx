import {useState} from 'react';
import {createStyles, Navbar, UnstyledButton, Tooltip, rem} from '@mantine/core';

import {
    IconGauge,
    IconDeviceDesktopAnalytics,
} from '@tabler/icons-react';
import {useQuery} from "@tanstack/react-query";

const useStyles = createStyles((theme) => ({
    wrapper: {
        display: 'flex',
    },

    aside: {
        flex: `0 0 ${rem(60)}`,
        backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.white,
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        borderRight: `${rem(1)} solid ${
            theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.colors.gray[3]
        }`,
    },

    main: {
        flex: 1,
        backgroundColor: 'white',
    },

    mainLink: {
        width: rem(44),
        height: rem(44),
        borderRadius: theme.radius.md,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.colors.gray[7],

        '&:hover': {
            backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[5] : theme.colors.gray[0],
        },
    },

    mainLinkActive: {
        '&, &:hover': {
            backgroundColor: theme.fn.variant({variant: 'light', color: theme.primaryColor}).background,
            color: theme.fn.variant({variant: 'light', color: theme.primaryColor}).color,
        },
    },

    title: {
        boxSizing: 'border-box',
        fontFamily: `Greycliff CF, ${theme.fontFamily}`,
        marginBottom: theme.spacing.xl,
        backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.white,
        padding: theme.spacing.md,
        paddingTop: rem(18),
        height: rem(60),
        borderBottom: `${rem(1)} solid ${
            theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.colors.gray[3]
        }`,
    },

    logo: {
        boxSizing: 'border-box',
        width: '100%',
        display: 'flex',
        justifyContent: 'center',
        height: rem(60),
        paddingTop: theme.spacing.md,
        borderBottom: `${rem(1)} solid ${
            theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.colors.gray[3]
        }`,
        marginBottom: theme.spacing.xl,
    },

    link: {
        boxSizing: 'border-box',
        display: 'block',
        textDecoration: 'none',
        borderTopRightRadius: theme.radius.md,
        borderBottomRightRadius: theme.radius.md,
        color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.colors.gray[7],
        padding: `0 ${theme.spacing.md}`,
        fontSize: theme.fontSizes.sm,
        marginRight: theme.spacing.md,
        fontWeight: 500,
        height: rem(44),
        lineHeight: rem(44),

        '&:hover': {
            backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[5] : theme.colors.gray[1],
            color: theme.colorScheme === 'dark' ? theme.white : theme.black,
        },
    },

    linkActive: {
        '&, &:hover': {
            borderLeftColor: theme.fn.variant({variant: 'filled', color: theme.primaryColor})
                .background,
            backgroundColor: theme.fn.variant({variant: 'filled', color: theme.primaryColor})
                .background,
            color: theme.white,
        },
    },
}));

const mainLinksMockdata = [
    {icon: IconGauge, label: 'Admin Dashboard'},
    {icon: IconDeviceDesktopAnalytics, label: 'Analytics'},
];

interface ManagerDetail {
    serverDetail: { clients: CLientData[], numberOfClients: number };
}

interface CLientData {
    client_id: string,
    login_in_at: string
}


function SideNav() {
    const {classes, cx} = useStyles();
    const [active, setActive] = useState('Admin Dashboard');

    const mainLinks = mainLinksMockdata.map((link) => (
        <Tooltip
            label={link.label}
            position="right"
            withArrow
            transitionProps={{duration: 0}}
            key={link.label}
        >
            <UnstyledButton
                onClick={() => setActive(link.label)}
                className={cx(classes.mainLink, {[classes.mainLinkActive]: link.label === active})}
            >
                <link.icon size="1.4rem" stroke={1.5}/>
            </UnstyledButton>
        </Tooltip>
    ));


    const {data} = useQuery({
            queryKey: ['clients'], queryFn: async () => {
                try {
                    const response = await fetch("/admin/stats");
                    if (response.ok) {
                        const data: ManagerDetail = await response.json();
                        console.log("Data fetched:", data);
                        return data.serverDetail.clients
                    } else {
                        console.log("Fetch Error:", response.statusText);
                    }
                } catch (error) {
                    console.error("Error fetching data:", error);
                }
                return []
            },
        }
    )


    return (
        <Navbar className={''} width={{sm: 400}}>
            <Navbar.Section grow className={classes.wrapper}>
                <div className={classes.aside + " mt-4"}>
                    <div className={classes.logo}>

                    </div>
                    {mainLinks}
                </div>
                <div className={classes.main}>
                    <div className={classes.title + " mt-4 "}>
                        <p className={'text-gray-600'}>{active}</p>
                    </div>

                    {<ul role="list" className="divide-y divide-gray-100">
                        {data?.map((client) => (
                            <li key={client.client_id} className="relative flex justify-between gap-x-6 py-4">
                                <div className="flex min-w-0 gap-x-4 ml-8">
                                    <div className="min-w-0 flex-auto">
                                        <p className="text-sm font-semibold leading-6 text-gray-900 cursor-pointer">
                                            <span className="absolute inset-x-0 -top-px bottom-0"/>
                                            {client.client_id}
                                        </p>
                                        <p className="mt-1 flex text-xs leading-5 text-gray-500">
                                            leslie.alexander@example.com
                                        </p>
                                    </div>
                                </div>

                            </li>
                        ))}
                    </ul>}
                </div>
            </Navbar.Section>
        </Navbar>

    )
}

export default SideNav;