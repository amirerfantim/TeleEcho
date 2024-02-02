import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
// import './css/ChatsPage.css';
import {getToken} from "./AuthHelper"; // Ensure you have the CSS similar to ContactsPage.css

// Assuming the chat model includes these fields
interface Chat {
    chatID: number;
    user: {
        id: number;
        username: string;
        firstname: string;
        lastname: string;
    };
    unreadMessage: number;
}

const ChatsPage: React.FC = () => {
    const [chats, setChats] = useState<Chat[]>([]);
    const ws = useRef<WebSocket | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        ws.current = new WebSocket('ws://127.0.0.1:8020/all-chat?token=' + getToken());

        ws.current.onmessage = (event) => {
            const data = JSON.parse(event.data);
            setChats(data);
        };

        ws.current.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        return () => {
            ws.current?.close();
        };
    }, []);

    const navigateToChat = (chatID: number) => {
        const adjustedChatID = chatID % 2 === 1 ? chatID + 1 : chatID;
        navigate(`/chat/${adjustedChatID}`);
    };

    return (
        <div className="centered-container">
            <ul className="contacts-container">
                {chats.map((chat) => (
                    <li key={chat.chatID} className="contact-card" onClick={() => navigateToChat(chat.chatID)}>
                        <div className="contact-info">
                            <div className="name">{chat.user.firstname} {chat.user.lastname}</div>
                            <div className="details">@{chat.user.username} | Unread Messages: {chat.unreadMessage}</div>
                        </div>
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default ChatsPage;
