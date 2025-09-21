'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '../../context/AuthContext';

const LoginPage: React.FC = () => {
  const [formData, setFormData] = useState({ username: '', password: '' });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const { login } = useAuth();
  const router = useRouter();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      await login(formData);
      router.push('/dashboard');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-wrapper">
      <h1 className="auth-title">Welcome Back</h1>
      <p className="auth-subtitle">Sign in to your account</p>

      <form onSubmit={handleSubmit} className="auth-form">
        <div className="form-row">
          <label htmlFor="username">Username</label>
          <input
            className="input-field"
            type="text"
            id="username"
            name="username"
            value={formData.username}
            onChange={handleChange}
            required
            disabled={loading}
          />
        </div>

        <div className="form-row">
          <label htmlFor="password">Password</label>
          <input
            className="input-field"
            type="password"
            id="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            required
            disabled={loading}
          />
        </div>

        {error && <div className="info-box">{error}</div>}

        <button
          type="submit"
          disabled={loading}
          className="button button-primary"
        >
          {loading ? 'Signing In...' : 'Sign In'}
        </button>
      </form>

      <div className="auth-footer">
        <p>
          Don&apos;t have an account?{' '}
          <Link href="/register" className="nav-link">Sign up</Link>
        </p>
      </div>
    </div>
  );
};

export default LoginPage;
