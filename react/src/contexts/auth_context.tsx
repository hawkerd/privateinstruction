'use client'
import React, { createContext, useContext, useEffect, useState } from 'react';
import config from '@/config';
import { useMemo } from 'react';

// define context type
type AuthContextType = {
    token: string | null; // the JWT token
    login: (token: string) => void; // function to set the token
    logout: () => void; // function to clear the token
    fetchWithAuth: (input: RequestInfo, init?: RequestInit) => Promise<Response>; // function to fetch with auth
    isAuthenticated: boolean; // boolean indicating if the user is authenticated (convenience)
}

// set up default context
const AuthContext = createContext<AuthContextType>({
    token: null,
    login: () => {},
    logout: () => {},
    fetchWithAuth: async (input: RequestInfo, init?: RequestInit) => {
        throw new Error('fetchWithAuth not implemented');
    },
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
    // state variable for the token
    const [token, setToken] = useState<string | null>(null);

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

    // function to refresh the JWT
    const refreshToken = async (oldToken: string) => {
        const response = await fetch(`${config.servicePath}/auth/refresh`, {
            method: 'POST',
            credentials: 'include',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${oldToken}`,
            },
        });
        console.log('refreshToken response', response);
        if (!response.ok) {
            logout();
            return null;
        }
        const data = await response.json();
        const newToken = data.access_token;
        return newToken;
    };

    // load the token from local storage on mount
    useEffect(() => {
        const init = async () => {
            const storedToken = localStorage.getItem('authToken');
            if (storedToken) {
                setToken(storedToken);
                if (isTokenExpired(storedToken)) {
                    const t = await refreshToken(storedToken);
                    if (t) {
                        login(t);
                    }
                }
            }
        }
        init();
    }, []);

    const isAuthenticated = useMemo(() => {
        return !!token && !isTokenExpired(token);
    }, [token]);


    // fetch function to handle requests with auth
    const fetchWithAuth = async (input: RequestInfo, init?: RequestInit): Promise<Response> => {
        // set up the request headers
        const headers = new Headers(init?.headers);
        headers.set('Content-Type', 'application/json');
        if (token) {
        headers.set('Authorization', `Bearer ${token}`);
        }

        // make the request
        let response = await fetch(input, {
        ...init,
        credentials: 'include',
        headers,
        });

        if (response.status !== 401) {
        return response;
        }

        if (!token) {
            return response;
        }
        const newToken = await refreshToken(token);
        if (!newToken) {
            return response;
        }
        login(newToken);

        // Retry original request with new token
        headers.set('Authorization', `Bearer ${newToken}`);

        response = await fetch(input, {
        ...init,
        credentials: 'include',
        headers,
        });

        return response;
    };

    // return the provider
    return (
        <AuthContext.Provider value={{ token, login, logout, fetchWithAuth, isAuthenticated }}>
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