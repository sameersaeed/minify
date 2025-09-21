'use client';

import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { urlAPI } from '../lib/api';
import { MinifyResponse } from '../types';

const HomePage: React.FC = () => {
    const { user, isAuthenticated } = useAuth();
    const [url, setUrl] = useState('');
    const [minifyUrl, setMinifyUrl] = useState<MinifyResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!url.trim()) {
            setError('Please enter a URL');
            return;
        }

        if (!url.match(/^https?:\/\/.+/)) {
            setError('Please enter a valid URL (must start with http:// or https://)');
            return;
        }
        
        setLoading(true);
        setError('');

        try {
            const result = await urlAPI.minify({
                url: url.trim(),
                user_id: user?.id,
            });
            setMinifyUrl(result);
            setUrl('');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to Minify URL');
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
                    <label htmlFor="url">Enter a URL</label>
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

                {error && <div className="alert alert-error">{error}</div>}

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
                    <h3 className="output-title">Success! 🎉</h3>

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
        </div>
    );
};

export default HomePage;
