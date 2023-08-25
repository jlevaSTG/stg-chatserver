import {create} from "zustand";

export interface ClientData {
    client_id: string,
    login_in_at: string
}

export interface ClientStore {
    clients: ClientData[];
    chatSession: Chats | undefined
    activeSession: Session | undefined
    activeClient: ClientData | undefined;
    setClients: (clients: ClientData[]) => void
    setActiveClient: (clientId: string) => void
    setChats: (chats: Chats) => void
    setActiveSession: (chatId: string) => void

}

export interface Chats {
    chatSessions: Session[]
}

export interface Session {
    chat_id: string
    created_at: string
    created_by: string
    active: boolean
    participants: Participant[]
    messages: Messages[]
}

export interface Participant {
    id: string
    active: boolean
    added_by: string
    joinedAt: string
    disconnectedAt: string
}

export interface Messages {
    created_at: string
    message_type: string
    chat_id: string
    created_by: string
    message: string
}

export const useClientStore = create<ClientStore>((set) => ({
    clients: [],
    chatSession: undefined,
    activeClient: undefined,
    activeSession: undefined,
    setClients: (clients: ClientData[]) => {
        set(() => ({
            clients,
        }));
    },
    setActiveClient: (clientId: string) => {
        set((state) => {
            const filteredClients = state.clients.filter(
                (c) => c.client_id === clientId
            );
            return {
                activeClient: filteredClients.length > 0 ? filteredClients[0] : undefined,
            };
        });
    },
    setChats: (chats: Chats) => {
        set(()=> ({
            chatSession: chats,
        }));
    },
    setActiveSession: (chatId: string) => {
        set((state: ClientStore) => {
            const filteredChats = state.chatSession?.chatSessions?.filter(
                (c) => c.chat_id === chatId
            );
            const activeSession: Session | undefined = filteredChats && filteredChats.length > 0 ? filteredChats[0] : undefined;
            return { activeSession };  // Make sure to return an object
        });
    },
}));