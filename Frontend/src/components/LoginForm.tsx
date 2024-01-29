import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './css/LoginForm.css';


interface Credentials {
    username: string;
    password: string;
}

const LoginForm: React.FC = () => {
    const [credentials, setCredentials] = useState<Credentials>({ username: '', password: '' });
    const [errorMessage, setErrorMessage] = useState<string>('');
    const navigate = useNavigate();

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setCredentials({ ...credentials, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setErrorMessage('');

        try {
            const response = await fetch('http://127.0.0.1:8020/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: `username=${encodeURIComponent(credentials.username)}&password=${encodeURIComponent(credentials.password)}`,
            });
            const data = await response.json();

            if (response.ok) {
                localStorage.setItem('token', data.token as string);
                navigate('/dashboard');
            } else {
                setErrorMessage(data.message || 'An error occurred');
            }
        } catch (error) {
            console.error('Login error:', error);
            setErrorMessage('An error occurred while attempting to log in.');
        }
    };

    return (
        <div className="container">
            <div className="card">
                <div className="logo"></div>
                <h1 className="title">Login</h1>
                <form onSubmit={handleSubmit}>
                    <input type="text" className="form-control" name="username" value={credentials.username} onChange={handleChange} placeholder="Username" required />
                    <input type="password" className="form-control" name="password" value={credentials.password} onChange={handleChange} placeholder="Password" required />
                    <div className="button-container">
                        <button type="submit" className="btn btn-primary">Login</button>
                        <button type="button" className="btn btn-secondary" onClick={() => navigate('/register')}>Register</button>
                    </div>
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

export default LoginForm;
