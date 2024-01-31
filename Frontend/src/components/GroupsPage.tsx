import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './css/GroupsPage.css';
import profileIcon from "../assets/profile.png";
import groupIcon from "../assets/group.png";
import contactIcon from "../assets/contact.png";

interface Group {
    id: number;
    name: string;
    description: string;
    profilePicture: string;
}

const GroupsPage: React.FC = () => {
    const [groups, setGroups] = useState<Group[]>([]);
    const [showAddGroupModal, setShowAddGroupModal] = useState(false);
    const [newGroupName, setNewGroupName] = useState('');
    const [newGroupDescription, setNewGroupDescription] = useState('');
    const [groupProfilePicture, setGroupProfilePicture] = useState<File | null>(null);
    const navigate = useNavigate();

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
            <button className="add-contact" onClick={() => setShowAddGroupModal(true)}>+</button>
            <div className="groups-container">
                {groups.map(group => (
                    <div key={group.id} className="group-card" onClick={() => navigateToGroup(group.id)}>
                        <h3>{group.name}</h3>
                        <p>{group.description}</p>
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
        </div>
    );


};

export default GroupsPage;
