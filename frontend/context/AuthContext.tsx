'use client';

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useRouter } from 'next/navigation';
import Cookies from 'js-cookie';
import { User, LoginRequest, CreateUserRequest } from '../types';
import { authAPI } from '../lib/api';

interface AuthContextType {
    user: User | null;
    isAuthenticated: boolean;
    isAdmin: boolean;
    login: (data: LoginRequest) => Promise<void>;
    register: (data: CreateUserRequest) => Promise<void>;
    logout: () => void;
    loading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);

    const router = useRouter();

    useEffect(() => {
        // check if user is logged in already
        const token = Cookies.get('token');
        const userData = Cookies.get('user');

        if (token && userData) {
            try {
                setUser(JSON.parse(userData));
            } catch (error) {
                console.error('Error parsing user data:', error);
                Cookies.remove('token');
                Cookies.remove('user');
            }
        }
        setLoading(false);
    }, []);

    const login = async (data: LoginRequest) => {
        try {
            const response = await authAPI.login(data);
            const { token, user: userData } = response;

            // store token + user data
            Cookies.set('token', token, { expires: 7 }); // 7 days as defined in server (see handlers/user.go)
            Cookies.set('user', JSON.stringify(userData), { expires: 7 });

            setUser(userData);
        } catch (error) {
            console.error('Login error:', error);
            throw error;
        }
    };

    const register = async (data: CreateUserRequest) => {
        try {
            const userData = await authAPI.register(data);

            // auto-login after registration
            await login({ username: data.username, password: data.password });
        } catch (error) {
            console.error('Registration error:', error);
            throw error;
        }
    };

    // unset user, remove tokens, and redirect to homepage
    const logout = () => {
        authAPI.logout();
        setUser(null);
        Cookies.remove('token');
        Cookies.remove('user');
        router.push('/');
    };

    const value: AuthContextType = {
        user,
        isAuthenticated: !!user,
        isAdmin: user ? authAPI.isAdmin() : false,
        login,
        register,
        logout,
        loading,
    };

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
