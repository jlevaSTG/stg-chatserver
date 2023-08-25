import {Box, Button, Center,  Grid, Group,  Modal, TextInput, ThemeIcon} from "@mantine/core";
import DashBoardHeader from "../components/DashboardHeader.tsx";
import {Messages} from "../store/clientStore.ts";
import {useEffect, useRef, useState} from "react";
import {IconTrash} from "@tabler/icons-react";
import {useForm} from "@mantine/form";
import ChatList from "../components/ChatList.tsx";
import {useQuery} from "@tanstack/react-query";
import {ChatSession, ChatSessionsResponse, Message} from "../types/modals.ts";
import {useDisclosure} from "@mantine/hooks";

interface InitChat {
    id: string
    created_by: string
    participants: {
        id: string
    }[]
    message: string
}

interface SendMsg {
    created_by: string,
    chat_id: string,
    message: string
}

interface InitChatResponse {
    chat: Chat
}

interface Chat {
    chat_id: string
    created_at: string
    created_by: string
    active: boolean
    participants: Participant[]
    messages: Messages[]
}

interface Participant {
    id: string
    active: boolean
    added_by: string
    joinedAt: string
    disconnectedAt: string
}


interface TextChatMessage {
    message_type: string;
    created_at: string;
    payload: Payload;
}

interface Payload {
    chat_id: string;
    sent_by: string;
    message: string;
    resource_url: string;
}

interface ChatMessageResponse {
    chat: {
        created_at: string;
        created_by: string;
        chat_id: string;
        participants: string[] | null;
        message: string;
    };
}

function formateDate(timestamp: any): string {
    const date = new Date(timestamp);
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
}

export default function AdminChat() {
    // const [currentChat, setCurrentChat] = useState<InitChatResponse | undefined>()
    const messagesEndRef = useRef<any>(null);
    const [selectedChat, setSelectedChat] = useState<ChatSession>();
    const [opened, { open, close }] = useDisclosure(false);

    function addMessage(msg: Message) {
        setSelectedChat(prevChat => {
            if (!prevChat) return undefined;


            // Create a new array containing the existing messages plus the new message
            const updatedMessages = [...prevChat.messages, msg];
            // Return a new ChatSession object with the updated messages array
            return {
                ...prevChat,
                messages: updatedMessages,
            };
        });
    }

    const {data, refetch} = useQuery<ChatSessionsResponse, Error>(
        ['chatSessions'],
        async () => {
            const response = await fetch(`api/retrieve/admin`);
            if (!response.ok) {
                throw new Error(`Fetch Error: ${response.statusText}`);
            }
            const data: ChatSessionsResponse = await response.json();
            setSelectedChat(data.chatSessions[0])
            console.log("Sessions Fetched:", data);
            return data;
        },
    )


    useEffect(() => {
        if (messagesEndRef.current) {
            messagesEndRef.current.scrollIntoView({behavior: 'smooth'});
        }
    }, [selectedChat]);

    useEffect(() => {
        // Initialize WebSocket connection
        const ws = new WebSocket('ws://localhost:8080/ws?userId=admin');

        // Function to run when a message is received
        ws.onmessage = (event) => {
            const parsedData: TextChatMessage = JSON.parse(event.data);
            console.log("Parsed data as TextChatMessage:", parsedData);
            if (parsedData.message_type === "text_chat_message") {
                const payload: Payload = parsedData.payload
                const msg: Message = {
                    created_at: parsedData.created_at,
                    message_type: parsedData.message_type,
                    chat_id: payload.chat_id,
                    created_by: payload.sent_by,
                    message: payload.message
                }

                addMessage(msg);

            }
        };

        // Function to run when the connection opens
        ws.onopen = () => {
            console.log('WebSocket connection opened');
            ws.send('Hello Server');
        };

        // Function to run if an error occurs
        ws.onerror = (error) => {
            console.error(`WebSocket Error: ${error}`);
        };

        // Function to run when the connection closes
        ws.onclose = () => {
            console.log('WebSocket connection closed');
        };

        // Cleanup: Close the WebSocket connection when this component is unmounted
        return () => {
            ws.close();
        };
    }, []); // Emp


    const form = useForm({
        initialValues: {
            message: "",
            clients: [],
        },
    });

    const sendMessageForm = useForm(
        {
            initialValues: {
                message: "",
            }
        }
    )


    async function initChat(form: {
        message: string,
        clients: {
            id: string
        }[]
    }) {
        const initChat: InitChat = {
            id: "admin",
            created_by: "admin",
            participants: form.clients.map(c => ({id: c.id})),
            message: form.message
        };
        console.log("initchat", initChat)
        try {
            const response = await fetch("/api/initChat", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(initChat),
            });
            const chatResponse: InitChatResponse = await response.json();
            console.log("Success:", chatResponse);
            // setCurrentChat(data)
            await refetch()
            setSelectedChat(data?.chatSessions[data.chatSessions.length - 1])
        } catch (error) {
            console.error("Error:", error);
        }

    }


    const fields = form.values.clients.map((_, index) => (
            <Group mt="xs" key={index}>
                <Center>
                    <ThemeIcon variant="outline" color={'gray'} onClick={() => form.removeListItem('clients', index)}>
                        <IconTrash className={'bg-white'} size="1.2rem"/>
                    </ThemeIcon>
                </Center>
                <TextInput placeholder="John Doe" {...form.getInputProps(`clients.${index}.id`)} />
            </Group>
        )
    );

    async function sendChat(values: ReturnType<(values: {
        message: string
    }) => {
        message: string
    }>) {
        if (selectedChat) {
            const sendMsg: SendMsg = {
                chat_id: selectedChat.chat_id,
                message: values.message,
                created_by: "admin",
            }
            try {
                const response = await fetch("/api/sendMsg", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(sendMsg),
                });
                const data: ChatMessageResponse = await response.json();
                console.log("Success:", data);
                const msg: Message = {
                    created_at: data.chat.created_at,
                    message_type: "text_chat_message",
                    chat_id: data.chat.chat_id,
                    created_by: data.chat.created_by,
                    message: data.chat.message
                }

                addMessage(msg);
                sendMessageForm.reset();
            } catch (error) {
                console.error("Error:", error);
            }


        }
    }

    return (
        <>
            <Modal opened={opened} onClose={close} title="Authentication">
                <form onSubmit={form.onSubmit((values) => initChat(values))}>
                    <Box maw={500} mx="auto">
                        <TextInput type="text" placeholder="Message"
                                   className="block w-full  ml-2 focus:text-gray-700"
                                   name="message" required
                                   {...form.getInputProps('message')}
                        />
                        <div>
                            {fields}
                        </div>
                        <Group position="center" mt="md">
                            <Button className={'bg-white'} variant={'default'}
                                    onClick={() => form.insertListItem('clients', {id: ''})}>
                                Add Client
                            </Button>
                            <Button className={'bg-white'} variant={'default'} type={'submit'}>
                                Start Chat
                            </Button>
                        </Group>

                    </Box>
                </form>
            </Modal>
            <DashBoardHeader/>
            <Grid gutter={5} gutterXl={0} className={''}>
                <Grid.Col span={4}>
                    <ChatList chatSessionsResponse={data} selectedChat={selectedChat}
                              setSelectedChat={setSelectedChat}/>
                </Grid.Col>
                <Grid.Col span={8}>

                    <main>
                        {
                            data ? (
                                    <div className="container px-4 sm:px-6 lg:px-8">
                                        <div className=" border rounded-md border-gray-300">
                                            <div>
                                                <div className="w-full">
                                                    <div
                                                        className="relative flex items-center p-3 border-b border-gray-300">
                                                        {/*<span className="block ml-2 text-sm text-gray-600">Chat with <span className="font-bold">{activeClient?.client_id}</span></span>*/}
                                                    </div>

                                                    <div className="relative w-full p-6 overflow-y-auto h-[40rem] bg-white">


                                                        <ul className="space-y-2">

                                                            {
                                                                selectedChat?.messages
                                                                    .map((m, idx) => (
                                                                        <li className={"flex flex-col " + (m.created_by === "admin" ? "justify-start items-start" : "justify-end items-end")}
                                                                            key={idx}>
                                                                            <div
                                                                                className={
                                                                                    m.created_by === "admin" ?
                                                                                        "relative max-w-xl px-4 py-2 text-gray-700 rounded border border-gray-300 shadow" :
                                                                                        "relative max-w-xl px-4 py-2 text-gray-700 bg-gray-100 rounded border border-gray-300  shadow"
                                                                                }>
                                                                                <span className="block">{m.message}</span>
                                                                            </div>
                                                                            <span
                                                                                className="mt-2 text-xs text-gray-500 uppercase">{m.created_by} - {formateDate(m.created_at)}</span>
                                                                        </li>
                                                                    ))
                                                            }

                                                        </ul>
                                                        <div ref={messagesEndRef}></div>
                                                    </div>

                                                    <form
                                                        onSubmit={sendMessageForm.onSubmit((values) => sendChat(values))}
                                                        className="flex items-center justify-between w-full p-3 border-t border-gray-300 bg-white">
                                                        <button onClick={open}>
                                                            <svg xmlns="http://www.w3.org/2000/svg"
                                                                 className="w-5 h-5 text-gray-500" fill="none"
                                                                 viewBox="0 0 24 24"
                                                                 stroke="currentColor">
                                                                <path strokeLinecap="round" strokeLinejoin="round"
                                                                      strokeWidth="2"
                                                                      d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"/>
                                                            </svg>
                                                        </button>

                                                        <TextInput type="text" placeholder="Message"
                                                                   className="block w-full  ml-2 focus:text-gray-700"
                                                                   name="message"
                                                                   {...sendMessageForm.getInputProps('message')}
                                                                   required/>

                                                        <Button variant="outline"  type={'submit'} className={'br-white border-white'}

                                                        >
                                                            <svg
                                                                className="w-5 h-5 text-gray-500 origin-center transform rotate-90"
                                                                xmlns="http://www.w3.org/2000/svg"
                                                                viewBox="0 0 20 20" fill="currentColor">
                                                                <path
                                                                    d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z"/>
                                                            </svg>
                                                        </Button>
                                                    </form>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                )
                                :
                                <></>
                        }
                    </main>
                </Grid.Col>
            </Grid>
        </>
    )
}