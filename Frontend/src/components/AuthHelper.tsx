export const isAuthenticated = (): boolean => {
    return !!localStorage.getItem('token');
};

export const getToken = (): string | null => {
    return localStorage.getItem('token');
};

export const removeToken = () => {
    localStorage.removeItem('token');
};
