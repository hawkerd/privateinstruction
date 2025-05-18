'use client'
import React, { createContext, useContext, useEffect, useState } from 'react';

// define context type
type AuthContextType = {
    token: string | null; // the JWT token
    login: (token: string) => void; // function to set the token
    logout: () => void; // function to clear the token
    isAuthenticated: boolean; // boolean indicating if the user is authenticated (convenience)
}

// set up default context
const AuthContext = createContext<AuthContextType>({
    token: null,
    login: () => {},
    logout: () => {},
    isAuthenticated: false,
});

// utility function to check if the token is expired
function isTokenExpired(token: string): boolean {
    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        const now = Math.floor(Date.now() / 1000);
        return payload.exp < now;
    } catch {
        return true; // assume expired if we can't decode
    }
}

// provider component
export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [token, setToken] = useState<string | null>(null);

    // load the token from local storage on mount
    useEffect(() => {
        const storedToken = localStorage.getItem('authToken');
        if (storedToken && !isTokenExpired(storedToken)) {
            setToken(storedToken);
        }
    }, []);

    // login function to set the token
    const login = (token: string) => {
        setToken(token);
        localStorage.setItem('authToken', token);
    };

    // logout function to clear the token
    const logout = () => {
        setToken(null);
        localStorage.removeItem('authToken');
    };

    // return the provider
    return (
        <AuthContext.Provider value={{ token, login, logout, isAuthenticated: !!token }}>
            {children}
        </AuthContext.Provider>
    );
}

// custom hook
export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}