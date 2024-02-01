import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './css/GroupsPage.css';
import profileIcon from "../assets/profile.png";
import groupIcon from "../assets/group.png";
import contactIcon from "../assets/contact.png";
import deleteIcon from "../assets/delete_icon.png";
import addUserIcon from "../assets/add.png";


interface Group {
    id: number;
    name: string;
    description: string;
    profilePicture: string;
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

    const fetchGroups = async () => {
        try {
            const response = await fetch('http://127.0.0.1:8020/group', {
                method: 'GET',
                headers: {
                    'Authorization': '' + localStorage.getItem('token'), // Replace with your auth token
                },
            });
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const data = await response.json();
            setGroups(data);
        } catch (error) {
            console.error('There has been a problem with your fetch operation:', error);
        }
    };

    useEffect(() => {
        fetchGroups();
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
                    'Authorization': '' + localStorage.getItem('token'), // Replace with your auth token
                },
                body: formData,
            });

            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            await response.json();
            setShowAddGroupModal(false);
            setNewGroupName('');
            setNewGroupDescription('');
            setGroupProfilePicture(null);
            fetchGroups(); // Refresh the groups list
        } catch (error) {
            console.error('Failed to create the group:', error);
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
                    // 'Content-Type': 'multipart/form-data' is not needed, browser sets it along with the correct boundary
                },
                body: formData,
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Failed to add user to group');
                    }
                    return response.json();
                })
                .then(() => {
                    setShowAddUsersModal(false);
                    fetchGroups();
                })
                .catch(error => setError(error.message));
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
            <button className="add-button-first" onClick={() => setShowAddGroupModal(true)}>+</button>

            <div className="groups-container">
                {groups.map(group => (
                    <div key={group.id} className="group-card" onClick={() => navigateToGroup(group.id)}>
                        <h3>{group.name}</h3>
                        <p>{group.description}</p>
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
                                    <button onClick={() => handleAddUsersToGroup(contact.username)}>Add to Group</button>
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
