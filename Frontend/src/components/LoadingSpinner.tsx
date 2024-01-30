import React from 'react';
import './css/LoadingSpinner.css'; // Make sure the path to your CSS file is correct

const LoadingSpinner: React.FC = () => {
    return (
        <div className="loading-spinner-container">
            <div className="loading-spinner"></div>
        </div>
    );
};

export default LoadingSpinner;
