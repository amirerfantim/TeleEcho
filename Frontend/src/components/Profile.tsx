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
    const [isEditing, setIsEditing] = useState(false);
    const [editingField, setEditingField] = useState<string | null>(null);
    const EditButton = () => (
        <button className="btn btn-secondary" onClick={() => setIsEditing(!isEditing)}>
            {/* Icon for edit, can be replaced with an actual icon */}
            {isEditing ? "Cancel" : "Edit"}
        </button>
    );

    interface CardContentProps {
        field: string;
        initialValue: string;
    }

    interface CardContentProps {
        field: string;
        initialValue: string;
        saveChanges: (updatedUserInfo: { [key: string]: string }) => Promise<void>;
        setEditingField: React.Dispatch<React.SetStateAction<string | null>>;
        isEditing: boolean;
    }

    const handleEditClick = () => {
        setIsEditing(!isEditing);
        if (isEditing) {
            setEditingField(null);
        }
    };

    const CardContent: React.FC<CardContentProps> = ({ field, initialValue, saveChanges, setEditingField, isEditing }) => {
        const [value, setValue] = useState(initialValue);
        const [isFieldEditing, setIsFieldEditing] = useState(false);

        const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
            setValue(e.target.value);
        };

        const handleSave = async () => {
            await saveChanges({ [field]: value });
            setEditingField(null);
        };




        const handleEditClick = () => {
            setIsFieldEditing(true);
            setEditingField(field);
        };

        const handleCancel = () => {
            setIsFieldEditing(false);
            setEditingField(null);
            setValue(initialValue); // Reset value to initial
        };

        return (
            <div className={`detail-card ${field}-card`}>
                {isFieldEditing ? (
                    <>
                        <input type="text" value={value} onChange={handleInputChange} />
                        <div>
                            <button onClick={handleSave}>Save</button>
                            <button onClick={handleCancel}>Cancel</button>
                        </div>
                    </>
                ) : (
                    <>
                        <p>{value}</p>
                        {isEditing && <button onClick={handleEditClick}>Edit</button>}
                    </>
                )}
            </div>
        );
    };


    const saveChanges = async (updatedUserInfo: { [key: string]: string }) => {
        try {
            const formData = new FormData();

            for (const key in updatedUserInfo) {
                formData.append(key, updatedUserInfo[key]);
            }

            const token = getToken();
            const response = await fetch('http://127.0.0.1:8020/users', {
                method: 'PATCH',
                headers: {
                    'Authorization': ` ${token}`,
                },
                body: formData
            });

            if (response.ok) {
                console.log("User information updated successfully");
                // Handle successful response
            } else {
                // Handle errors
                const errorData = await response.json();
                console.error("Error updating user information:", errorData);
            }
        } catch (error) {
            console.error("Error in API request:", error);
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
            <button className="btn btn-secondary" onClick={handleEditClick}>
                {isEditing ? "Cancel" : "Edit"}
            </button>
            <div className="profile-picture-container">
                {user.profilePicture && (
                    <img src={user.profilePicture} alt={`${user.firstname} ${user.lastname}`} className="profile-picture" />
                )}
            </div>
            <div className="name-card">
                <p>{user.firstname} {user.lastname}</p>
            </div>
            <div className="details-container">
                <CardContent
                    field="username"
                    initialValue={user.username}
                    saveChanges={saveChanges}
                    setEditingField={setEditingField}
                    isEditing={isEditing}
                />
                <CardContent
                    field="phone"
                    initialValue={user.phone}
                    saveChanges={saveChanges}
                    setEditingField={setEditingField}
                    isEditing={isEditing}
                />
            </div>
            <div className="bio-card">
                <p>{user.Bio}</p>
            </div>
        </div>
    );
};
export default Profile;
