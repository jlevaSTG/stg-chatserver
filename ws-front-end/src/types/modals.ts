export interface Message {
    created_at: string;
    message_type: string;
    chat_id: string;
    created_by: string;
    message: string;
    resource_url?: string;
}

export interface Participant {
    id: string;
    active: boolean;
    added_by: string;
    joinedAt: string;
    disconnectedAt: string;
}

export interface ChatSession {
    chat_id: string;
    created_at: string;
    created_by: string;
    active: boolean;
    participants: Participant[];
    messages: Message[];
}

export interface ChatSessionsResponse {
    chatSessions: ChatSession[];
}
