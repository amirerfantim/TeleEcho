import React, {useState, useEffect} from 'react';
import {useNavigate} from 'react-router-dom';
import {isAuthenticated, getToken} from './AuthHelper';
import './css/Profile.css';
import LoadingSpinner from "./LoadingSpinner";

interface User {
    username: string;
    firstname: string;
    lastname: string;
    phone: string;
    Bio: string;
    profilePicture?: string;
}

const Profile: React.FC = () => {
    const [user, setUser] = useState<User | null>(null);
    const [errorMessage, setErrorMessage] = useState<string>('');
    const navigate = useNavigate();
    const handleProfilePicHover = (isHover: boolean) => {
        const profilePicElement = document.querySelector('.profile-picture');
        if (profilePicElement) {
            profilePicElement.classList.toggle('hover', isHover);
        }
    };

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

                if (response.ok) {
                    const data = await response.json();
                    if (data.profilePicture && !data.profilePicture.startsWith('http')) {
                        data.profilePicture = `data:image/png;base64,${data.profilePicture}`;
                    }
                    setUser(data);
                } else {
                    const contentType = response.headers.get('content-type');
                    if (contentType && contentType.includes('application/json')) {
                        const data = await response.json();
                        handleServerError(data);
                    } else {
                        const text = await response.text();
                        handleServerError({error: text});
                    }
                }
            } catch (error) {
                console.error('Error fetching user info:', error);
                setErrorMessage('An error occurred while fetching user information.');
            }
        };

        fetchUserInfo();
    }, [navigate]);

    const handleServerError = (data: any) => {
        setErrorMessage('Invalid token. Redirecting to login...');

        const timeoutId = setTimeout(() => {
            navigate('/login');
        }, 2000);

        return () => clearTimeout(timeoutId);
    };

    if (errorMessage) {
        return <div className="profile-container error-container">
            <div className="alert alert-danger">
                {errorMessage}
            </div>
        </div>;
    }

    if (!user) {
        return <LoadingSpinner />;
    }


    return (
        <div className="profile-container">
            <div className="profile-picture-container">
                {user?.profilePicture && (
                    <img src={user.profilePicture} alt={`${user?.firstname} ${user?.lastname}`} className="profile-picture" />
                )}
            </div>
            <div className="name-card">
                <p>{user?.firstname} {user?.lastname}</p>
            </div>
            <div className="details-container">
                <div className="detail-card username-card">
                    <p>{user?.username}</p>
                </div>
                <div className="detail-card phone-card">
                    <p>{user?.phone}</p>
                </div>
            </div>
            <div className="bio-card">
                <p>{user?.Bio}</p>
            </div>
        </div>
    );
};
export default Profile;
