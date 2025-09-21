'use client';

import React from 'react';
import Link from 'next/link';
import { useAuth } from '../context/AuthContext';

const Header: React.FC = () => {
    const { user, isAuthenticated, isAdmin, logout } = useAuth();

    return (
        <header className="header">
            <div className="header-container">
                <Link href="/" className="brand-title">
                    Minify
                </Link>

                <nav className="header-nav">
                    {isAuthenticated ? (
                        <>
                            <Link href="/dashboard" className="nav-link">
                                Dashboard
                            </Link>

                            {isAdmin && (
                                <Link href="/admin" className="nav-link">
                                    Admin
                                </Link>
                            )}

                            <div className="nav-item-group">
                                <span className="nav-user">
                                    Welcome, {user?.username}
                                </span>
                                <button
                                    onClick={logout}
                                    className="button"
                                >
                                    Logout
                                </button>
                            </div>
                        </>
                    ) : (
                        <div className="nav-item-group">
                            <Link href="/login" className="nav-link">
                                Login
                            </Link>
                            <Link href="/register" className="button button-primary">
                                Sign Up
                            </Link>
                        </div>
                    )}
                </nav>
            </div>
        </header>
    );
};

export default Header;
