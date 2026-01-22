'use client';

import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { urlAPI } from '../lib/api';
import { MinifyResponse } from '../types';
import { toast, ToastContainer } from 'react-toastify';

import 'react-toastify/dist/ReactToastify.css';

const HomePage: React.FC = () => {
    const { user, isAuthenticated } = useAuth();
    const [url, setUrl] = useState('');
    const [minifyUrl, setMinifyUrl] = useState<MinifyResponse | null>(null);
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!url.trim()) {
            toast.error('Please enter a URL');
            return;
        }

        if (!url.match(/^https?:\/\/.+/)) {
            toast.error('Please enter a valid URL (must start with http:// or https://)');
            return;
        }
        
        setLoading(true);

        try {
            const result = await urlAPI.minify({
                url: url.trim(),
                user_id: user?.id,
            });
            setMinifyUrl(result);
            setUrl('');
            toast.success('URL successfully Minified!');
        } catch (err: any) {
            toast.error(err.response?.data?.error || 'Failed to Minify URL');
        } finally {
            setLoading(false);
        }
    };

    const copyToClipboard = (text: string) => navigator.clipboard.writeText(text);

    return (
        <div className="home-wrapper">
            <h1 className="home-title">Minify: Shorten your URLs</h1>

            <form onSubmit={handleSubmit} className="minify-form">
                <div className="form-row">
                    <label htmlFor="url">Enter a URL to Minify</label>
                    <input
                        type="url"
                        id="url"
                        value={url}
                        onChange={(e) => setUrl(e.target.value)}
                        placeholder="https://example.com/url"
                        className="input-field"
                        disabled={loading}
                    />
                </div>
                <button
                    type="submit"
                    disabled={loading}
                    className="button button-primary minify-btn"
                >
                    {loading ? 'Minifying...' : 'Minify URL'}
                </button>

                {!isAuthenticated && (
                    <div className="info-box">
                        💡 <strong>Tip:</strong> Create an account to save and manage your Minified URLs!
                    </div>
                )}
            </form>

            {minifyUrl && (
                <div className="output-box">
                    <div className="form-row">
                        <label>Short URL:</label>
                        <input
                            type="text"
                            value={minifyUrl.short_url}
                            readOnly
                            className="input-field"
                        />
                    </div>
                    <button
                        type="button"
                        onClick={() => copyToClipboard(minifyUrl.short_url)}
                        className="button button-secondary copy-btn"
                    >
                        Copy
                    </button>

                    <div className="form-row">
                        <label>Original URL:</label>
                        <input
                            type="text"
                            value={minifyUrl.original_url}
                            readOnly
                            className="input-field"
                        />
                    </div>
                </div>
            )}

            <ToastContainer
                position="top-right"
                autoClose={3000}
                hideProgressBar={false}
                newestOnTop
                closeOnClick
                rtl={false}
                pauseOnFocusLoss
                draggable
                pauseOnHover
            />
        </div>
    );
};

export default HomePage;
