'use client';

import React, { createContext, useContext, useState, useEffect } from 'react';
import { User } from '@/models/yuh/user';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/auth_context';
import config from '@/config';

// define context type
type UserContextType = {
    user: User | null; // the user object
    setUser: (user: User | null) => void; // function to set the user
    clearUser: () => void; // function to clear the user
};

// create a default context
const UserContext = createContext<UserContextType>({
    user: null,
    setUser: () => {},
    clearUser: () => {},
});

// provider component
export const UserProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { token, isAuthenticated, fetchWithAuth } = useAuth();
    const [user, setUser] = useState<User | null>(null);
    
    // clear the user
    const clearUser = () => {
        setUser(null);
    }

    // fetch user data
    useEffect(() => {
        // if not authenticated, clear the user and return
        if (!isAuthenticated || !token) {
            clearUser();
            return;
        }

        // fetch user data from the API
        const fetchUser = async () => {
            try {
                const res = await fetchWithAuth(`${config.servicePath}/me`, {
                    method: 'GET',
                });

                if (!res.ok) {
                    throw new Error('Failed to fetch user data');
                }

                // parse the response
                const data: User = await res.json();
                setUser(data);
            } catch (err) {
                console.error(err);
                clearUser();
            }
        };

        fetchUser();
    }, [isAuthenticated, token]);

    return (
        <UserContext.Provider value={{ user, setUser, clearUser }}>
            {children}
        </UserContext.Provider>
    );
}


// custom hook
export const useUser = () => {
    const context = useContext(UserContext);
    if (!context) {
        throw new Error('useUser must be used within a UserProvider');
    }
    return context;
};