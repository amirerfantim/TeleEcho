import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './css/GroupsPage.css';
import profileIcon from "../assets/profile.png";
import groupIcon from "../assets/group.png";
import contactIcon from "../assets/contact.png";
import deleteIcon from "../assets/delete_icon.png";
import addUserIcon from "../assets/add.png";
import chatIcon from "../assets/chat.png";


interface Group {
    id: number;
    name: string;
    description: string;
    profilePicture: string;
    users: Contact[]; // Add this line to include users in each group
}


interface Contact {
    id: number;
    username: string;
    firstname: string;
    lastname: string;
    phone: string;
}

const GroupsPage: React.FC = () => {
    const [groups, setGroups] = useState<Group[]>([]);
    const [showAddGroupModal, setShowAddGroupModal] = useState(false);
    const [newGroupName, setNewGroupName] = useState('');
    const [newGroupDescription, setNewGroupDescription] = useState('');
    const [groupProfilePicture, setGroupProfilePicture] = useState<File | null>(null);
    const navigate = useNavigate();
    const [contacts, setContacts] = useState<Contact[]>([]);
    const [showAddUsersModal, setShowAddUsersModal] = useState(false);
    const [selectedGroupId, setSelectedGroupId] = useState<number | null>(null);
    const [error, setError] = useState('');
    const [responseMessage, setResponseMessage] = useState('');

    const fetchGroupsAndUsers = async () => {
        try {
            const response = await fetch('http://127.0.0.1:8020/group', {
                method: 'GET',
                headers: {
                    'Authorization': '' + localStorage.getItem('token'),
                },
            });
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const groupsData = await response.json();

            const groupsWithUsers = await Promise.all(groupsData.map(async (group: Group) => {
                const usersResponse = await fetch(`http://127.0.0.1:8020/group/all?groupID=${group.id}`, {
                    method: 'GET',
                    headers: {
                        'Authorization': '' + localStorage.getItem('token'),
                    },
                });
                if (!usersResponse.ok) {
                    throw new Error('Failed to fetch users for group ' + group.id);
                }
                const users = await usersResponse.json();
                return { ...group, users }; // Add the users to the group object
            }));

            setGroups(groupsWithUsers);
        } catch (error) {
            console.error('There has been a problem with your fetch operation:', error);
        }
    };


    useEffect(() => {
        fetchGroupsAndUsers();
    }, []);

    const handleNewGroupNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setNewGroupName(e.target.value);
    };

    const handleNewGroupDescriptionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setNewGroupDescription(e.target.value);
    };

    const handleGroupProfilePictureChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files[0]) {
            setGroupProfilePicture(e.target.files[0]);
        }
    };

    const fetchContacts = () => {
        fetch('http://127.0.0.1:8020/contacts', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': '' + localStorage.getItem('token'),
            },
        })
            .then(response => {
                if (response.ok) {
                    return response.json();
                }
                throw new Error('Network response was not ok.');
            })
            .then(data => {
                if (Array.isArray(data)) {
                    setContacts(data);
                } else {
                    setError('The contact list is empty');
                }
            })
            .catch(error => setError(error.message));
    };

    const openAddUsersModal = (groupId: number) => {
        setSelectedGroupId(groupId);
        setShowAddUsersModal(true);
        fetchContacts();
    };

    const handleCreateGroup = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        if (!groupProfilePicture) {
            alert('Please select a profile picture for the group.');
            return;
        }

        const formData = new FormData();
        formData.append('name', newGroupName);
        formData.append('description', newGroupDescription);
        formData.append('profile', groupProfilePicture);

        try {
            const response = await fetch('http://127.0.0.1:8020/group', {
                method: 'POST',
                headers: {
                    'Authorization': '' + localStorage.getItem('token'),
                },
                body: formData,
            });

            setShowAddGroupModal(false);


            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || 'Failed to create the group');
            }

            await response.json();
            setResponseMessage('Group created successfully.');
            setNewGroupName('');
            setNewGroupDescription('');
            setGroupProfilePicture(null);
            fetchGroupsAndUsers(); // Refresh the groups list
        } catch (error) {
            if (error instanceof Error) {
                setResponseMessage(error.message);
            } else {
                setResponseMessage('An unexpected error occurred');
            }
        }
    };



    const handleDeleteGroup = async (groupId: number) => {
        // Show confirmation dialog
        const confirmDelete = window.confirm("Are you sure you want to delete this group?");
        if (!confirmDelete) {
            return; // Stop the function if user does not confirm
        }

        try {
            const response = await fetch(`http://127.0.0.1:8020/group?groupID=${groupId}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': '' + localStorage.getItem('token'), // Replace with your auth token
                    'Content-Type': 'application/json',
                },
            });

            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            // Remove the group from the state
            setGroups(groups.filter(group => group.id !== groupId));
        } catch (error) {
            console.error('There has been a problem with your delete operation:', error);
        }
    };

    const handleAddUsersToGroup = (username: string) => {
        if (selectedGroupId) {
            const formData = new FormData();
            formData.append('groupID', selectedGroupId.toString());
            formData.append('username', username);

            fetch(`http://127.0.0.1:8020/group`, {
                method: 'PATCH',
                headers: {
                    'Authorization': '' + localStorage.getItem('token'),
                },
                body: formData,
            })
                .then(response => {
                    setShowAddUsersModal(false);
                    if (!response.ok) {
                        // Extracting error message from response
                        return response.text().then(text => {
                            throw new Error(text || 'Failed to add user to group');
                        });
                    }
                    return response.json();
                })
                .then(() => {
                    setResponseMessage('User added to the group successfully.');
                    fetchGroupsAndUsers();
                })
                .catch(error => {
                    setResponseMessage(error.message);
                });
        }
    };




    const navigateToGroup = (groupId: number) => {
        navigate(`/group/${groupId}`);
    };

    return (
        <div className="groups-page">
            <button className="navigate-first" onClick={() => navigate('/profile')}>
                <img src={profileIcon} alt="Profile"/>
            </button>
            <button className="navigate-second" onClick={() => navigate('/contacts')}>
                <img src={contactIcon} alt="Groups"/>
            </button>
            <button className="navigate-third" onClick={() => navigate('/chats')}>
                <img src={chatIcon} alt="Chat"/>
            </button>
            <button className="add-button-first" onClick={() => setShowAddGroupModal(true)}>+</button>

            {responseMessage && (
                <div className="response-message">{responseMessage}</div>
            )}

            <div className="groups-container">
                {groups.map(group => (
                    <div key={group.id} className="group-card" onClick={() => navigateToGroup(group.id)}>
                        <h3>{group.name}</h3>
                        <p>{group.description}</p>
                        {/* Display users for each group */}
                        <div className="group-users">
                            {group.users && group.users.map(user => (
                                <div key={user.id} className="group-user">
                                    {user.firstname} {user.lastname}
                                </div>
                            ))}
                        </div>
                        <div className="group-card-actions">
                            <button className="delete-button" onClick={(e) => {
                                e.stopPropagation(); // Prevents navigating to the group detail page
                                handleDeleteGroup(group.id);
                            }}>
                                <img src={deleteIcon} alt="Delete"/>
                            </button>
                            <button className="add-user-group-button" onClick={(e) => {
                                e.stopPropagation(); // Prevent navigating to group detail
                                openAddUsersModal(group.id);
                            }}>
                                <img src={addUserIcon} alt="Add User"/>
                            </button>
                        </div>
                    </div>
                ))}
            </div>
            {showAddGroupModal && (
                <div className="modal">
                    <div className="modal-content">
                        <span className="close-button" onClick={() => setShowAddGroupModal(false)}>&times;</span>
                        <h2>Add New Group</h2>
                        <form onSubmit={handleCreateGroup}>
                            <label htmlFor="group-name">Group Name:</label>
                            <input
                                id="group-name"
                                type="text"
                                value={newGroupName}
                                onChange={handleNewGroupNameChange}
                                placeholder="Group Name"
                                required
                            />
                            <label htmlFor="group-description">Description:</label>
                            <input
                                id="group-description"
                                type="text"
                                value={newGroupDescription}
                                onChange={handleNewGroupDescriptionChange}
                                placeholder="Group Description"
                                required
                            />
                            <label htmlFor="group-profile-picture">Profile Picture:</label>
                            <input
                                id="group-profile-picture"
                                type="file"
                                onChange={handleGroupProfilePictureChange}
                                accept="image/*"
                                required
                            />
                            <button type="submit">Create Group</button>
                        </form>
                    </div>
                </div>
            )}

            {showAddUsersModal && (
                <div className="modal">
                    <div className="modal-content">
                        <span className="close-button" onClick={() => setShowAddUsersModal(false)}>&times;</span>
                        <h2>Add Users to Group</h2>
                        <div>
                            {contacts.map(contact => (
                                <div key={contact.id} className="contact-card">
                                    <div>{contact.firstname} {contact.lastname}</div>
                                    <button onClick={() => handleAddUsersToGroup(contact.username)}>Add to Group
                                    </button>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );


};

export default GroupsPage;
