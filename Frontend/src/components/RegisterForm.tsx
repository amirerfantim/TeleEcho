
import React, { useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import './RegisterForm.css'; // Make sure this CSS file is in the same directory

interface FormData {
    username: string;
    firstname: string;
    lastname: string;
    phone: string;
    password: string;
    profile: string;
    bio: string;
}

const RegisterForm = () => {
    const [formData, setFormData] = useState<FormData>({
        username: '',
        firstname: '',
        lastname: '',
        phone: '',
        password: '',
        profile: '',
        bio: ''
    });
    const [successMessage, setSuccessMessage] = useState('');
    const [errorMessage, setErrorMessage] = useState('');
    const navigate = useNavigate();
    const formRef = useRef<HTMLFormElement>(null);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setErrorMessage('');
        setSuccessMessage('');

        if (formRef.current) {
            const formDataObj = new FormData(formRef.current);
            for (const key in formData) {
                formDataObj.append(key, formData[key as keyof FormData]);
            }

            try {
                const response = await fetch('http://127.0.0.1:8020/user/register', {
                    method: 'POST',
                    body: formDataObj
                });
                const data = await response.json();

                if (response.ok) {
                    setSuccessMessage('User created successfully. Redirecting to login...');
                    setTimeout(() => navigate('/login'), 2000); // Redirect after 2 seconds
                } else {
                    setErrorMessage(data);
                }
            } catch (error) {
                console.error(error);
                setErrorMessage('An error occurred while registering the user.');
            }
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
                    <input type="text" className="form-control" name="profile" value={formData.profile} onChange={handleChange} placeholder="Profile Picture URL" />
                    <textarea className="form-control" name="bio" value={formData.bio} onChange={handleChange} placeholder="Bio" rows={3}></textarea>
                    <div className="button-container">
                        <button type="submit" className="btn btn-primary">Submit</button>
                        <button type="button" className="btn btn-secondary" onClick={() => navigate('/login')}>Login</button>
                    </div>
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
                </form>
            </div>
        </div>
    );
};

export default RegisterForm;