import React, { useState, useEffect, useRef } from 'react';
import { getToken } from "./AuthHelper";
import './css/DirectChatPage.css';
import {useParams} from "react-router-dom";
import { useNavigate } from 'react-router-dom';
import groupIcon from "../assets/group.png"; // Adjust path as per your project structure
import contactIcon from "../assets/contact.png";
import profileIcon from "../assets/profile.png"; // Adjust path as per your project structure
import chatIcon from "../assets/chat.png" ; // Adjust path as per your project structure


interface Message {
    id: number;
    content: string;
    senderId: number;
    createdAt: string;
}

const DirectChatPage: React.FC = () => {
    const {chatId} = useParams<{ chatId: string }>();
    // const chatId = "10";
    const [messages, setMessages] = useState<Message[]>([]);
    const [newMessage, setNewMessage] = useState('');
    const wsSend = useRef<WebSocket | null>(null);
    const wsReceive = useRef<WebSocket | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        // Connect to SendMessage WebSocket
        const token = '' + getToken();
        wsSend.current = new WebSocket(`ws://127.0.0.1:8020/write-message?chatID=${chatId}&token=${token}`);
        wsSend.current.onopen = () => console.log('Send WS Connection established');
        wsSend.current.onmessage = (e) => {
            const message = JSON.parse(e.data);
            addMessageSorted(message);
        };

        // Connect to GetMessage WebSocket
        wsReceive.current = new WebSocket(`ws://127.0.0.1:8020/read-message?chatID=${chatId}&token=${token}`);
        wsReceive.current.onopen = () => {
            console.log('Receive WS Connection established');
            fetchMessagesPeriodically({count: 10}); // Fetch last 10 messages periodically
        };
        wsReceive.current.onmessage = (e) => {
            const receivedMessages = JSON.parse(e.data);
            setMessagesSortById(receivedMessages);
        };

        return () => {
            wsSend.current?.close();
            wsReceive.current?.close();
        };
    }, [chatId]); // Reconnect if chatId changes

    const sendMessage = () => {
        if (newMessage.trim() !== '') {
            wsSend.current?.send(JSON.stringify({content: newMessage, stat: 'sending'}));
            setNewMessage('');
        }
    };

    const fetchMessagesPeriodically = ({count}: { count: any }) => {
        // Request messages periodically
        const requestMessages = () => {
            if(wsReceive.current && wsReceive.current.readyState === WebSocket.OPEN) {
                wsReceive.current.send(JSON.stringify({count: count, stat: 'request'}));
            }
        };

        requestMessages(); // Initial fetch
        const intervalId = setInterval(requestMessages, 5000); // Fetch every 5 seconds

        return () => clearInterval(intervalId); // Cleanup on component unmount or when chatId changes
    };

    const addMessageSorted = (newMessage: Message) => {
        setMessages((prevMessages) => {
            const updatedMessages = [...prevMessages, newMessage];
            return updatedMessages.sort((a, b) => a.id - b.id);
        });
    };

    const setMessagesSortById = (newMessages: Message[]) => {
        const sortedMessages = [...newMessages].sort((a, b) => a.id - b.id);
        setMessages(sortedMessages);
    };

    return (
        <div className="chat-page">
            <button className="navigate-first" onClick={() => navigate('/profile')}>
                <img src={profileIcon} alt="Profile"/>
            </button>
            <button className="navigate-second" onClick={() => navigate('/contacts')}>
                <img src={contactIcon} alt="Contacts"/>
            </button>
            <button className="navigate-third" onClick={() => navigate('/group')}>
                <img src={groupIcon} alt="Group"/>
            </button>
            <button className="navigate-fourth" onClick={() => navigate('/chats')}>
                <img src={chatIcon} alt="chats"/>
            </button>
            <div className="messages-container">
                {messages.map((message) => (
                    <div key={message.id} className={`message ${message.senderId === 22 ? 'sent' : 'received'}`}>
                        {message.content}
                    </div>
                ))}
            </div>
            <div className="message-input">
                <input type="text" value={newMessage} onChange={(e) => setNewMessage(e.target.value)}/>
                <button onClick={sendMessage}>Send</button>
            </div>
        </div>
    );
};

export default DirectChatPage;
