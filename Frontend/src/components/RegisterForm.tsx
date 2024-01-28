import React, { useState, useRef } from 'react';
import axios from 'axios';

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

    const formRef = useRef<HTMLFormElement>(null);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        if (formRef.current) {
            const formDataObj = new FormData(formRef.current);
            formDataObj.append('username', formData.username);
            formDataObj.append('firstname', formData.firstname);
            formDataObj.append('lastname', formData.lastname);
            formDataObj.append('phone', formData.phone);
            formDataObj.append('password', formData.password);
            formDataObj.append('profile', formData.profile);
            formDataObj.append('bio', formData.bio);

            try {
                const response = await axios.post('http://127.0.0.1:8020/user/register', formDataObj, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
                console.log(response.data);
                // Handle successful response
            } catch (error) {
                console.error(error);
                // Handle error
            }
        }
    };

    return (
        <div className="container mt-5">
            <h2>Register</h2>
            <form ref={formRef} onSubmit={handleSubmit} className="mt-4">
                <div className="mb-3">
                    <label htmlFor="username" className="form-label">Username</label>
                    <input type="text" className="form-control" id="username" name="username" value={formData.username}
                           onChange={handleChange} required/>
                </div>
                <div className="mb-3">
                    <label htmlFor="firstname" className="form-label">First Name</label>
                    <input type="text" className="form-control" id="firstname" name="firstname"
                           value={formData.firstname} onChange={handleChange} required/>
                </div>
                <div className="mb-3">
                    <label htmlFor="lastname" className="form-label">Last Name</label>
                    <input type="text" className="form-control" id="lastname" name="lastname" value={formData.lastname}
                           onChange={handleChange} required/>
                </div>
                <div className="mb-3">
                    <label htmlFor="phone" className="form-label">Phone</label>
                    <input type="tel" className="form-control" id="phone" name="phone" value={formData.phone}
                           onChange={handleChange} required/>
                </div>
                <div className="mb-3">
                    <label htmlFor="password" className="form-label">Password</label>
                    <input type="password" className="form-control" id="password" name="password"
                           value={formData.password} onChange={handleChange} required/>
                </div>
                <div className="mb-3">
                    <label htmlFor="profile" className="form-label">Profile Picture URL</label>
                    <input type="text" className="form-control" id="profile" name="profile" value={formData.profile}
                           onChange={handleChange}/>
                </div>
                <div className="mb-3">
                    <label htmlFor="bio" className="form-label">Bio</label>
                    <textarea className="form-control" id="bio" name="bio" value={formData.bio} onChange={handleChange}
                              rows={3}></textarea>
                </div>
                <button type="submit" className="btn btn-primary">Submit</button>
            </form>
        </div>
    );
};

export default RegisterForm;
