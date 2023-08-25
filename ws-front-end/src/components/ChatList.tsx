import {Avatar} from "@mantine/core";
import {ChatSession, ChatSessionsResponse} from "../types/modals.ts";


interface ChatListProps {
    chatSessionsResponse: ChatSessionsResponse | undefined;
    selectedChat: ChatSession | undefined
    setSelectedChat: (session: ChatSession) => void
}

function formateDate(timestamp: any): string {
    const date = new Date(timestamp);
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
}


const ChatList: React.FC<ChatListProps> = ({selectedChat, setSelectedChat, chatSessionsResponse}) => {
    return (
        <div>
            <ul role="list" className="divide-y divide-gray-100  border rounded-md bg-white overflow-y-auto h-[40rem] ">
                {
                    chatSessionsResponse?.chatSessions.map((s) => (
                        <li
                            onClick={() => setSelectedChat(s)}
                            key={s.chat_id}
                            className={`flex flex-wrap items-center justify-between gap-x-6 gap-y-4 py-5 sm:flex-nowrap hover:shadow-xl p-8 hover:bg-gray-200 hover:border-gray-200 hover:border hover:rounded-md ${selectedChat?.chat_id === s.chat_id ? "bg-gray-200 shadow-xl" : ""}`}
                        >
                            <div>
                                <p className="text-sm font-semibold leading-6 text-gray-900 ">
                                    <a className="hover:underline">
                                        {s.chat_id}
                                    </a>
                                </p>
                                <div className="mt-1 flex items-center gap-x-2 text-xs leading-5 text-gray-500">
                                    <p>
                                        <a className="hover:underline">
                                            {s.created_by}
                                        </a>
                                    </p>
                                    <svg viewBox="0 0 2 2" className="h-0.5 w-0.5 fill-current">
                                        <circle cx={1} cy={1} r={1}/>
                                    </svg>
                                    <p>
                                        <time dateTime={s.created_at}>{formateDate(s.created_at)}</time>
                                    </p>
                                </div>
                            </div>
                            <dl className="flex w-full flex-none justify-between gap-x-8 sm:w-auto">
                                <div className="flex -space-x-0.5">
                                    <dt className="sr-only">Commenters</dt>
                                    <Avatar.Group spacing="sm">
                                        {s.participants.slice(0, 5).map((p, index) => (
                                            <Avatar key={index} color="cyan" radius="xl">{p.id.slice(0, 2)}</Avatar>
                                        ))}
                                        {s.participants.length > 5 && (
                                            <Avatar radius="xl">+{s.participants.length - 5}</Avatar>
                                        )}
                                    </Avatar.Group>
                                </div>
                            </dl>
                        </li>
                    ))
                }
            </ul>
        </div>
    )
}

export default ChatList;