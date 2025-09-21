'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../../context/AuthContext';
import { urlAPI } from '../../lib/api';
import { URL } from '../../types';

const DashboardPage: React.FC = () => {
  const { user } = useAuth();
  const [urls, setUrls] = useState<URL[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const router = useRouter();

  // redirect user to homepage if they're not logged in
  useEffect(() => {
    if (!user?.id) {
      router.replace('/'); 
      return;
    }
    fetchUserURLs();
  }, [user]);

  useEffect(() => {
    if (user?.id) fetchUserURLs();
  }, [user]);

  const fetchUserURLs = async () => {
    if (!user?.id) return;
    try {
      setLoading(true);
      const data = await urlAPI.getUserURLs(user.id);
      setUrls(data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load URLs');
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading) {
    return (
      <div className="dashboard-wrapper">
        <div className="dashboard-loading">
          <p>Loading your Minified URLs...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="dashboard-wrapper">
      <div className="dashboard-header">
        <h1>Dashboard</h1>
        <p>Welcome back, {user?.username}! Here are your Minified URLs.</p>
      </div>

      <div className="stats-container">
        <div className="stats-card">
          <h3>Total URLs</h3>
          <p>{urls.length}</p>
        </div>
        <div className="stats-card">
          <h3>Total Clicks</h3>
          <p>{urls.reduce((sum, url) => sum + url.clicks, 0)}</p>
        </div>
        <div className="stats-card">
          <h3>Average Clicks</h3>
          <p>
            {urls.length
              ? Math.round(urls.reduce((sum, url) => sum + url.clicks, 0) / urls.length)
              : 0}
          </p>
        </div>
      </div>

      <div className="urls-section">
        <div className="urls-header">
          <h2>Your URLs</h2>
          <button onClick={fetchUserURLs}>Refresh</button>
        </div>

        {error && <div className="dashboard-info">{error}</div>}

        {urls.length === 0 ? (
          <div className="dashboard-info">
            <p>You haven't Minified any URLs yet</p>
            <a href="/">Minify your first URL</a>
          </div>
        ) : (
          <table className="urls-table">
            <thead>
              <tr>
                <th>Short URL</th>
                <th>Original URL</th>
                <th>Clicks</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {urls.map(url => (
                <tr key={url.id}>
                  <td>
                    <code>
                      {process.env.NEXT_PUBLIC_API_URL?.replace(/^https?:\/\//, '')}/{url.short_code}
                    </code>
                    <button
                      onClick={() =>
                        copyToClipboard(`${process.env.NEXT_PUBLIC_API_URL}/${url.short_code}`)
                      }
                      title="Copy to clipboard"
                    >
                      📋
                    </button>
                  </td>
                  <td title={url.original_url}>{url.original_url}</td>
                  <td>{url.clicks}</td>
                  <td>{formatDate(url.created_at)}</td>
                  <td>
                    <a
                      href={`${process.env.NEXT_PUBLIC_API_URL}/${url.short_code}`}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      Visit
                    </a>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};

export default DashboardPage;
