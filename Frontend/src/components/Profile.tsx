import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { isAuthenticated, getToken } from './AuthHelper';
import './css/Profile.css';

interface User {
    username: string;
    firstname: string;
    lastname: string;
    phone: string;
    bio: string;
    profilePicture?: string;
}

const Profile: React.FC = () => {
    const [user, setUser] = useState<User | null>(null);
    const [errorMessage, setErrorMessage] = useState<string>('');
    const navigate = useNavigate();

    useEffect(() => {
        if (!isAuthenticated()) {
            navigate('/login');
            return;
        }

        const fetchUserInfo = async () => {
            try {
                const token = getToken();
                const response = await fetch('http://127.0.0.1:8020/users', {
                    method: 'GET',
                    headers: {
                        'Authorization': ` ${token}`,
                        'Content-Type': 'application/json'
                    },
                });

                const data = await response.json();
                if (response.ok) {
                    // Check if the profile picture is in Base64 format
                    if (data.profilePicture && !data.profilePicture.startsWith('http')) {
                        data.profilePicture = `data:image/png;base64,${data.profilePicture}`;
                    }
                    setUser(data);
                } else {
                    handleServerError(data);
                }
            } catch (error) {
                console.error('Error fetching user info:', error);
                setErrorMessage('An error occurred while fetching user information.');
            }
        };

        fetchUserInfo();
    }, [navigate]);

    const handleServerError = (data: any) => {
        setErrorMessage(data.error || data || 'Failed to retrieve user information.');
    };

    if (errorMessage) {
        return <div className="profile-container">
            <div className="alert alert-danger">
                {errorMessage}
            </div>
        </div>;
    }

    if (!user) {
        return <div className="profile-container">Loading...</div>;
    }

    return (
        <div className="profile-container">
            <h1>User Profile</h1>
            <div className="user-info">
                <p><strong>Username:</strong> {user.username}</p>
                <p><strong>First Name:</strong> {user.firstname}</p>
                <p><strong>Last Name:</strong> {user.lastname}</p>
                <p><strong>Phone:</strong> {user.phone}</p>
                <p><strong>Bio:</strong> {user.bio}</p>
                {user.profilePicture && (
                    <img src={user.profilePicture} alt="Profile" className="profile-picture" />
                )}
            </div>
        </div>
    );
};

export default Profile;
