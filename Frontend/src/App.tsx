import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import RegisterForm from './components/RegisterForm';
import LoginForm from './components/LoginForm';
import Profile from './components/Profile';
import React from "react";
import ContactsPage from "./components/ContactsPage";
import GroupsPage from "./components/GroupsPage";
import DirectChatPage from "./components/DirectChatPage";
import ChatsPage from "./components/ChatsPage";

const App: React.FC = () => {
    return (
        <Router>
            <Routes>
                <Route path="/register" element={<RegisterForm />} />
                <Route path="/login" element={<LoginForm />} />
                <Route path="/profile" element={<Profile />} />
                <Route path="/contacts" element={<ContactsPage />} />
                <Route path="/group" element={<GroupsPage />} />
                <Route path="/chat/:chatId" element={<DirectChatPage />} />
                <Route path="/chats" element={<ChatsPage />} />


            </Routes>
        </Router>
    );
};

export default App;
