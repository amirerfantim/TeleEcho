import React, { useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import './RegisterForm.css';

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
                    setTimeout(() => navigate('/login'), 2000);
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
        <div className="container mt-5">
            <div className="row justify-content-center">
                <div className="col-md-6">
                    <div className="card">
                        <div className="card-header">
                            <h2>Register</h2>
                        </div>
                        <div className="card-body">
                            <form ref={formRef} onSubmit={handleSubmit}>
                                <div className="mb-3">
                                    <label htmlFor="username" className="form-label">Username</label>
                                    <input type="text" className="form-control" id="username" name="username"
                                           value={formData.username}
                                           onChange={handleChange} required/>
                                </div>
                                <div className="mb-3">
                                    <label htmlFor="firstname" className="form-label">First Name</label>
                                    <input type="text" className="form-control" id="firstname" name="firstname"
                                           value={formData.firstname} onChange={handleChange} required/>
                                </div>
                                <div className="mb-3">
                                    <label htmlFor="lastname" className="form-label">Last Name</label>
                                    <input type="text" className="form-control" id="lastname" name="lastname"
                                           value={formData.lastname}
                                           onChange={handleChange} required/>
                                </div>
                                <div className="mb-3">
                                    <label htmlFor="phone" className="form-label">Phone</label>
                                    <input type="tel" className="form-control" id="phone" name="phone"
                                           value={formData.phone}
                                           onChange={handleChange} required/>
                                </div>
                                <div className="mb-3">
                                    <label htmlFor="password" className="form-label">Password</label>
                                    <input type="password" className="form-control" id="password" name="password"
                                           value={formData.password} onChange={handleChange} required/>
                                </div>
                                <div className="mb-3">
                                    <label htmlFor="profile" className="form-label">Profile Picture URL</label>
                                    <input type="text" className="form-control" id="profile" name="profile"
                                           value={formData.profile}
                                           onChange={handleChange}/>
                                </div>
                                <div className="mb-3">
                                    <label htmlFor="bio" className="form-label">Bio</label>
                                    <textarea className="form-control" id="bio" name="bio" value={formData.bio}
                                              onChange={handleChange}
                                              rows={3}></textarea>
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

                                <div className="button-container">
                                    <button type="submit" className="btn btn-primary">Submit</button>
                                    <button type="button" className="btn btn-success" onClick={() => navigate('/login')}>Login</button>
                                </div>

                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );

};

export default RegisterForm;
