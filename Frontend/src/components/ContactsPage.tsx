// ContactsPage.tsx
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

    return (
        <div className="contacts-container">
            {error && <p>Error: {error}</p>}
            {contacts && contacts.length > 0 ? (
                <ul>
                    {contacts.map(contact => (
                        <li key={contact.id}>{contact.name} - {contact.phone}</li>
                    ))}
                </ul>
            ) : (
                <p>No contacts found.</p>
            )}
        </div>
    );
};


export default ContactsPage;
