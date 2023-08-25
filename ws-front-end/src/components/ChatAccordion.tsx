import { useClientStore} from "../store/clientStore.ts";

import {Accordion, createStyles, rem} from '@mantine/core';

const useStyles = createStyles((theme) => ({
    root: {
        backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[0],
        borderRadius: theme.radius.sm,
    },

    item: {
        backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[0],
        border: `${rem(1)} solid transparent`,
        position: 'relative',
        zIndex: 0,
        transition: 'transform 150ms ease',

        '&[data-active]': {
            transform: 'scale(1.03)',
            backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.white,
            boxShadow: theme.shadows.xl,
            borderColor: theme.colorScheme === 'dark' ? theme.colors.dark[4] : theme.colors.gray[2],
            borderRadius: theme.radius.md,
            zIndex: 1,
        },
    },

    chevron: {
        '&[data-rotate]': {
            transform: 'rotate(-90deg)',
        },
    },
}));


function formateDate(timestamp: any): string {
    const date = new Date(timestamp);
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
}
export default function ChatAccordion() {
    const {classes} = useStyles();
    const {chatSession, setActiveSession} = useClientStore()
    return (
        <Accordion
            // maw={420}
            mx="auto"
            variant="filled"
            classNames={classes}
            onChange={() => setActiveSession("")}

            className={classes.root}
        >
            {
                chatSession?.chatSessions?.map(s => (
                    <Accordion.Item value={s.chat_id} key={s.chat_id} onClick={() => setActiveSession(s.chat_id)}>
                        <Accordion.Control>Chat ID: {s.chat_id}</Accordion.Control>
                        <Accordion.Panel>
                            <div className="mt-8 flow-root">
                                <div className="">
                                    <div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                                        <table className="min-w-full divide-y divide-gray-300">
                                            <thead>
                                            <tr>
                                                <th scope="col" className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-0">
                                                    Client ID
                                                </th>
                                                <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                                                    Added By
                                                </th>
                                                <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                                                    Status
                                                </th>
                                                <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                                                    Joined At
                                                </th>
                                                <th scope="col" className="relative py-3.5 pl-3 pr-4 sm:pr-0">
                                                    <span className="sr-only">Edit</span>
                                                </th>
                                            </tr>
                                            </thead>
                                            <tbody className="divide-y divide-gray-200 bg-white">
                                            {s.participants.map((p) => (
                                                <tr key={p.id}>
                                                    <td className="whitespace-nowrap py-5 pl-4 pr-3 text-sm sm:pl-0">
                                                        <div className="flex items-center">
                                                            <div className="ml-4">
                                                                <div className="font-medium text-gray-900">{p.id}</div>
                                                            </div>
                                                        </div>
                                                    </td>
                                                    <td className="whitespace-nowrap px-3 py-5 text-sm text-gray-500">
                                                        <div className="text-gray-900">{p.added_by}</div>
                                                    </td>
                                                    <td className="whitespace-nowrap px-3 py-5 text-sm text-gray-500">
                                                      <span className="inline-flex items-center rounded-md bg-green-50 px-2 py-1 text-xs font-medium text-green-700 ring-1 ring-inset ring-green-600/20">
                                                       {p.active ? "Active" : "Disconnected"}
                                                      </span>
                                                    </td>
                                                    <td className="whitespace-nowrap px-3 py-5 text-sm text-gray-500">{formateDate(p.joinedAt)}</td>
                                                    <td className="relative whitespace-nowrap py-5 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                                                        <a href="#" className="text-indigo-600 hover:text-indigo-900">
                                                            Edit<span className="sr-only">, {p.id}</span>
                                                        </a>
                                                    </td>
                                                </tr>
                                            ))}
                                            </tbody>
                                        </table>
                                    </div>
                                </div>
                            </div>
                        </Accordion.Panel>
                    </Accordion.Item>
                ))
            }
        </Accordion>
    );
}