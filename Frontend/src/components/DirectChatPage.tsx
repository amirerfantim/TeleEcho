import React, { useState, useEffect, useRef } from 'react';
import { getToken } from "./AuthHelper";
import './css/DirectChatPage.css';
import {useParams} from "react-router-dom";
import { useNavigate } from 'react-router-dom';
import groupIcon from "../assets/group.png"; // Adjust path as per your project structure
import contactIcon from "../assets/contact.png";
import profileIcon from "../assets/profile.png"; // Adjust path as per your project structure
import chatIcon from "../assets/chat.png" ;
import deleteIcon from "../assets/delete_icon.png"; // Adjust path as per your project structure


interface Message {
    id: number;
    content: string;
    senderId: number;
    createdAt: string;
}

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


const DirectChatPage: React.FC = () => {
    const {chatId} = useParams<{ chatId: string }>();
    // const chatId = "10";
    const [messages, setMessages] = useState<Message[]>([]);
    const [newMessage, setNewMessage] = useState('');
    const wsSend = useRef<WebSocket | null>(null);
    const wsReceive = useRef<WebSocket | null>(null);
    const navigate = useNavigate();
    const [showModal, setShowModal] = useState(false);
    const [copiedMessage, setCopiedMessage] = useState('');
    const [chats, setChats] = useState<Chat[]>([]); // Assuming Chat interface is defined


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

    useEffect(() => {
        const token = '' + getToken();
        const wsSendMessageUrl = `ws://127.0.0.1:8020/write-message?chatID=${chatId}&token=${token}`;
        wsSend.current = new WebSocket(wsSendMessageUrl);

        wsSend.current.onopen = () => console.log('Send WS Connection established');

        // Cleanup function to close WebSocket connections
        return () => {
            wsSend.current?.close();
        };
    }, [chatId]);

    useEffect(() => {
        const token = getToken();
        const wsUrl = `ws://127.0.0.1:8020/all-chat?token=${token}`;
        const wsChats = new WebSocket(wsUrl);

        wsChats.onmessage = (e) => {
            const receivedChats = JSON.parse(e.data);
            console.log("Fetched chats:", receivedChats); // Log the fetched chats
            setChats(receivedChats);
        };

        return () => {
            wsChats.close();
        };
    }, []); // Empty d



    const handleCopyMessage = (content: string) => {
        setCopiedMessage(content);
        setShowModal(true);
    };


    const handleSendCopiedMessage = ({chatID}: { chatID: any }) => {
        if (copiedMessage.trim() !== '') {
            sendMessageToChat({chatID: chatID, messageContent: copiedMessage});
            setCopiedMessage('');
            setShowModal(false);
        }
    };




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

    const deleteChat = async () => {
        const confirmDelete = window.confirm("Are you sure you want to delete this chat?");
        if (confirmDelete) {
            try {
                const response = await fetch(`http://127.0.0.1:8020/chat/${chatId}`, {
                    method: 'DELETE',
                    headers: {
                        'Authorization': `${getToken()}`,
                        'Content-Type': 'application/json',
                    },
                });

                if (response.ok) {
                    alert("Chat deleted successfully.");
                    navigate('/chats'); // Redirect to a safe route after deletion
                } else {
                    alert("Failed to delete chat.");
                }
            } catch (error) {
                console.error("Error deleting chat:", error);
                alert("An error occurred while deleting the chat.");
            }
        }
    }

    class ChatSelectionModal extends React.Component<{ isOpen: boolean, onClose: () => void, onChatSelect: (chatID: number) => void }> {
        render() {
            const { isOpen, onClose, onChatSelect } = this.props;
            if (!isOpen) return null;

            return (
                <div className="modal">
                    <div className="modal-content">
                        <span className="close" onClick={onClose}>&times;</span>
                        <ul>
                            {chats.map((chat) => (
                                <li key={chat.chatID} onClick={() => onChatSelect(chat.chatID)}>
                                    {chat.user.firstname} {chat.user.lastname}
                                </li>
                            ))}
                        </ul>
                    </div>
                </div>
            );
        }
    }

    // Helper function to send a message using a new WebSocket connection
    const sendMessageToChat = ({chatID, messageContent}: { chatID: any, messageContent: any }) => {
        const token = '' + getToken(); // Ensure you have a function to get the current user's token
        const adjustedChatID = chatID % 2 === 1 ? chatID + 1 : chatID;
        const wsUrl = `ws://127.0.0.1:8020/write-message?chatID=${adjustedChatID}&token=${token}`;
        const tempWs = new WebSocket(wsUrl);

        tempWs.onopen = () => {
            console.log('Temporary WS Connection for sending message established');
            console.log(chatID)
            tempWs.send(JSON.stringify({
                chatID: chatID,
                content: messageContent,
                stat: 'sending',
            }));

            // Close the connection after sending the message
            tempWs.close();
        };

        tempWs.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        // Optionally handle the tempWs.onclose event
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
            <button className="add-button-second" onClick={deleteChat}>
                <img src={deleteIcon} alt="Profile"/>
            </button>
            <div className="messages-container">
                {messages.map((message) => (
                    <div key={message.id} className={`message ${message.senderId === 22 ? 'sent' : 'received'}`}>
                        {message.content}
                        <button onClick={() => handleCopyMessage(message.content)}>Copy</button>
                    </div>

                ))}
            </div>
            <div className="message-input">
                <input type="text" value={newMessage} onChange={(e) => setNewMessage(e.target.value)}/>
                <button onClick={sendMessage}>Send</button>
            </div>

            {showModal && (
                <ChatSelectionModal
                    isOpen={showModal}
                    onClose={() => setShowModal(false)}
                    onChatSelect={(chatID) => handleSendCopiedMessage({chatID: chatID})}
                />
            )}

        </div>
    );
};

export default DirectChatPage;
