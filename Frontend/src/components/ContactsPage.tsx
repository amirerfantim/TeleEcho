import React, { useState, useEffect } from 'react';
import './css/ContactsPage.css';

interface Contact {
    id: number;
    name: string;
    phone: string;
}

const ContactsPage = () => {
    const [contacts, setContacts] = useState<Contact[]>([]);
    const [error, setError] = useState('');
    const [successMessage, setSuccessMessage] = useState('');
    const [showAddContactForm, setShowAddContactForm] = useState(false);
    const [newContactPhone, setNewContactPhone] = useState('');

    useEffect(() => {
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
                    setError('Data format is incorrect');
                }
            })
            .catch(error => setError(error.message));
    }, []);

    const toggleAddContactForm = () => {
        setShowAddContactForm(!showAddContactForm);
    };

    const handleNewContactPhoneChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setNewContactPhone(e.target.value);
    };

    const handleAddContact = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        setError('');
        setSuccessMessage('');

        fetch('http://127.0.0.1:8020/contacts', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'Authorization': '' + localStorage.getItem('token'),
            },
            body: `phone=${encodeURIComponent(newContactPhone)}`
        })
            .then(response => {
                // Parse the JSON response body
                return response.json().then(data => {
                    if (response.ok) {
                        // Handle success
                        setSuccessMessage('Contact added successfully');
                        setShowAddContactForm(false);
                        setNewContactPhone('');
                    } else {
                        // Handle errors, check if the response has an "error" key
                        if (data && data.error) {
                            throw new Error(data.error);
                        } else {
                            throw new Error('An error occurred while adding the contact.');
                        }
                    }
                });
            })
            .catch(error => {
                // Display errors
                setError(error.message);
            });
    };

    return (
        <div className="contacts-container">
            {successMessage && <p className="success-message">{successMessage}</p>}
            {error && <p className="error-message">{error}</p>}
            <button onClick={toggleAddContactForm}>Add Contact</button>
            {showAddContactForm && (
                <form onSubmit={handleAddContact}>
                    <input
                        type="text"
                        value={newContactPhone}
                        onChange={handleNewContactPhoneChange}
                        placeholder="Enter phone number"
                    />
                    <button type="submit">Create Contact</button>
                </form>
            )}
            <ul>
                {contacts.map(contact => (
                    <li key={contact.id}>{contact.name} - {contact.phone}</li>
                ))}
            </ul>
        </div>
    );
};

export default ContactsPage;
