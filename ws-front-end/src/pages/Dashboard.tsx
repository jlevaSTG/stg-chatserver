import {Grid, Modal} from "@mantine/core";


import 'ag-grid-community/styles/ag-grid.css'; // Core grid CSS, always needed
import 'ag-grid-community/styles/ag-theme-alpine.css'; // Optional theme CSS

import DashBoardHeader from "../components/DashboardHeader.tsx";
import {useQuery} from "@tanstack/react-query";

import {useDisclosure} from "@mantine/hooks";
import {Chats, useClientStore} from "../store/clientStore.ts";

import {

     ChatBubbleLeftEllipsisIcon,


} from '@heroicons/react/20/solid'
import ChatAccordion from "../components/ChatAccordion.tsx";





interface ManagerDetail {
    serverDetail: { clients: CLientData[], numberOfClients: number };
}

interface CLientData {
    client_id: string,
    login_in_at: string
}

function Dashboard() {
    const [opened, { close}] = useDisclosure(false);
    const {activeClient, setClients, setChats, activeSession} = useClientStore()


    useQuery({
            queryKey: ['clients'], queryFn: async () => {
                try {
                    const response = await fetch("/admin/stats");
                    if (response.ok) {
                        const data: ManagerDetail = await response.json();
                        console.log("Data fetched:", data);
                        setClients(data.serverDetail.clients)
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
    useQuery({
            queryKey: [activeClient], queryFn: async () => {
                try {
                    const response = await fetch(`/api/retrieve/${activeClient?.client_id}`);
                    if (response.ok) {
                        const data: Chats = await response.json();
                        console.log("Sessions Fetched:", data);
                        setChats(data)
                        return data
                    } else {
                        console.log("Fetch Error:", response.statusText);
                    }
                } catch (error) {
                    console.error("Error fetching data:", error);
                }
                return []
            },
        }
    );
    function formateDate(timestamp: any): string {
        const date = new Date(timestamp);
        return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
    }


    return (
        <Grid gutter={5} gutterXl={0} className={''}>
            <Grid.Col span={12}>
                <Modal opened={opened} onClose={close} title="CLient Connection Control" size={700}>
                    {JSON.stringify(activeClient?.client_id)}
                </Modal>
                <DashBoardHeader/>
                <main>
                    <div className="mx-auto px-4 sm:px-6 lg:px-8">
                        <div
                            className="mx-auto grid max-w-2xl grid-cols-1 grid-rows-1 items-start gap-x-8 gap-y-8 lg:mx-0 lg:max-w-none lg:grid-cols-3 ">

                            <div
                                className="bg-white -mx-4 px-4 py-8 shadow-sm ring-1 ring-gray-900/5 sm:mx-0 sm:rounded-lg sm:px-8 sm:pb-14 lg:col-span-2 lg:row-span-2 lg:row-end-2 xl:px-16 xl:pb-20 xl:pt-16">
                                <h2 className="text-base font-semibold leading-6 text-gray-600">User Detail</h2>
                                <dl className="mt-4 grid grid-cols-1 text-sm leading-6 sm:grid-cols-2 ">
                                    <div className="sm:pr-4 ">
                                        <dt className="inline text-gray-500">Logged In on</dt>
                                        {' '}
                                        <dd className="inline text-gray-700">
                                            <time dateTime="2023-23-01">{formateDate(activeClient?.login_in_at)}</time>
                                        </dd>
                                    </div>
                                    <div className="mt-2 sm:mt-0 sm:pl-4">
                                        <dt className="inline text-gray-500">User ID:</dt>
                                        {' '}
                                        <dd className="inline text-gray-700 ml-4 text-blue-400">
                                            <time>{activeClient?.client_id}</time>
                                        </dd>
                                    </div>
                                </dl>
                                <div className="mt-8 flow-root border-t">
                                    <div className={'mt-10'}>
                                        <ChatAccordion  />
                                    </div>
                                </div>

                            </div>

                            <div className="lg:col-start-3">
                                {/* Activity feed */}
                                <h2 className="text-sm font-semibold leading-6 text-gray-900">Chat Messages</h2>

                                <div className="flow-root mt-10">
                                    <ul role="list" className="-mb-8">
                                        {activeSession?.messages.map((m, activityItemIdx) => (
                                            <li key={m.created_at}>
                                                <div className="relative pb-8">
                                                    {activityItemIdx !== activeSession.messages.length - 1 ? (
                                                        <span className="absolute left-5 top-5 -ml-px h-full w-0.5 bg-gray-200" aria-hidden="true" />
                                                    ) : null}
                                                    <div className="relative flex items-start space-x-3">
                                                        {m.message_type === 'text_message' ? (
                                                            <>
                                                                <div className="relative">
                                                                    <img
                                                                        className="flex h-10 w-10 items-center justify-center rounded-full bg-gray-400 ring-8 ring-white"
                                                                        src={""}
                                                                        alt=""
                                                                    />
                                                                    <span className="absolute -bottom-0.5 -right-1 rounded-tl bg-white px-0.5 py-px">
                                                                        <ChatBubbleLeftEllipsisIcon className="h-5 w-5 text-gray-400" aria-hidden="true" />
                                                                      </span>
                                                                </div>
                                                                <div className="min-w-0 flex-1">
                                                                    <div>
                                                                        <div className="text-sm">
                                                                            <div  className="font-medium text-gray-900">
                                                                                {m.created_by}
                                                                            </div>
                                                                        </div>
                                                                        <p className="mt-0.5 text-sm text-gray-500">Commented {formateDate(m.created_at)}</p>
                                                                    </div>
                                                                    <div className="mt-2 text-sm text-gray-700">
                                                                        <p>{m.message}</p>
                                                                    </div>
                                                                </div>
                                                            </>
                                                        )  : null}
                                                    </div>
                                                </div>
                                            </li>
                                        ))}
                                    </ul>
                                </div>





                                {/* New comment form */}

                            </div>
                        </div>
                    </div>
                </main>
            </Grid.Col>
        </Grid>

    );
}

export default Dashboard;
