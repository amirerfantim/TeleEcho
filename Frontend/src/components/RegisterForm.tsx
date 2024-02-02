import React, { useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import './css/styles.css';
//
const RegisterForm = () => {
    const [formData, setFormData] = useState({
        username: '',
        firstname: '',
        lastname: '',
        phone: '',
        password: '',
        bio: ''
    });
    const [profilePicture, setProfilePicture] = useState<File | null>(null);
    const [successMessage, setSuccessMessage] = useState('');
    const [errorMessage, setErrorMessage] = useState('');
    const navigate = useNavigate();
    const formRef = useRef<HTMLFormElement>(null);
    const validatePassword = (password: string) => {
        return password.length >= 8;
    };

    const validateBio = (bio: string) => bio.length <= 100;
    const validateProfile = (file: File) => file.size <= 1048576; // Size must be less than or equal to 1MB


    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;

        // Reset error message on each change
        setErrorMessage('');

        if (name === 'profile') {
            const file = (e.target as HTMLInputElement).files?.[0];
            if (file && validateProfile(file)) {
                setProfilePicture(file);
            } else {
                setProfilePicture(null);
                if (file) {
                    setErrorMessage('Profile picture must be less than 1MB');
                }
            }
        } else if (name === 'bio') {
            if (validateBio(value)) {
                setFormData({ ...formData, [name]: value });
            } else {
                setErrorMessage('Biography must be less than 100 characters');
            }
        } else {
            setFormData({ ...formData, [name]: value });
        }
    };
    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setErrorMessage('');
        setSuccessMessage('');

        if (!validatePassword(formData.password)) {
            setErrorMessage('Password must be at least 8 characters long');
            return;
        }
        if (!validateBio(formData.bio)) {
            setErrorMessage('Biography must be less than 100 characters');
            return;
        }
        if (profilePicture && !validateProfile(profilePicture)) {
            setErrorMessage('Profile picture must be less than 1MB');
            return;
        }

        const formDataObj = new FormData();
        for (const key in formData) {
            formDataObj.append(key, formData[key as keyof typeof formData]);
        }
        if (profilePicture) {
            formDataObj.append('profile', profilePicture);
        }

        try {
            const response = await fetch('http://127.0.0.1:8020/register', {
                method: 'POST',
                body: formDataObj
            });
            const data = await response.json();

            if (response.ok) {
                setSuccessMessage('User created successfully. Redirecting to login...');
                setTimeout(() => navigate('/login'), 2000);
            } else {
                setErrorMessage(data || 'Registration failed');
            }
        } catch (error) {
            console.error(error);
            setErrorMessage('An error occurred while registering the user.');
        }
    };

    return (
        <div className="container">
            <div className="card">
                <div className="logo"></div>
                <h1 className="title">Register</h1>
                <form ref={formRef} onSubmit={handleSubmit}>
                    <input type="text" className="form-control" name="username" value={formData.username} onChange={handleChange} placeholder="Username" required />
                    <input type="text" className="form-control" name="firstname" value={formData.firstname} onChange={handleChange} placeholder="First Name" required />
                    <input type="text" className="form-control" name="lastname" value={formData.lastname} onChange={handleChange} placeholder="Last Name" required />
                    <input type="tel" className="form-control" name="phone" value={formData.phone} onChange={handleChange} placeholder="Phone" required />
                    <input type="password" className="form-control" name="password" value={formData.password} onChange={handleChange} placeholder="Password" required />
                    <input type="file" className="form-control" name="profile" onChange={handleChange}
                           placeholder="Profile Picture"/>
                    <textarea className="form-control" name="bio" value={formData.bio} onChange={handleChange} placeholder="Bio" rows={3}></textarea>
                    {successMessage && (
                        <div className="alert alert-success" role="alert">
                            {successMessage}
                        </div>
                    )}
                    {errorMessage && (
                        <div className="alert alert-danger" role="alert">
                            {errorMessage}
                        </div>
                    )}
                    <div className="button-container">
                        <button type="submit" className="btn btn-primary">Submit</button>
                        <button type="button" className="btn btn-secondary" onClick={() => navigate('/login')}>Login</button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default RegisterForm;