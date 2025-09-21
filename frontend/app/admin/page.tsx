'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../../context/AuthContext';
import { analyticsAPI } from '../../lib/api';
import { OverviewStats, PopularURL } from '../../types';
import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    PieChart,
    Pie,
    Cell,
} from 'recharts';

const AdminPage: React.FC = () => {
    const { isAdmin, loading: authLoading } = useAuth();
    const router = useRouter();
    const [stats, setStats] = useState<OverviewStats | null>(null);
    const [popularUrls, setPopularUrls] = useState<PopularURL[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    // redirect user to homepage if they're not an admin
    useEffect(() => {
        if (!authLoading && !isAdmin) {
            router.push('/');
            return;
        }

        if (isAdmin) fetchAnalytics();
    }, [isAdmin, authLoading, router]);

    const fetchAnalytics = async () => {
        try {
            setLoading(true);
            const [overviewData, popularData] = await Promise.all([
                analyticsAPI.getOverview(),
                analyticsAPI.getPopularURLs(10),
            ]);
            setStats(overviewData);
            setPopularUrls(popularData);
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to load analytics');
        } finally {
            setLoading(false);
        }
    };

    if (authLoading || loading) {
        return (
            <div className="dashboard-wrapper">
                <div className="dashboard-loading">
                    <p>Loading admin dashboard...</p>
                </div>
            </div>
        );
    }

    if (!isAdmin) return null;

    const timeframeData = stats?.timeframe_data
        ? Object.entries(stats.timeframe_data).map(([period, data]: [string, any]) => ({
              period: period.charAt(0).toUpperCase() + period.slice(1),
              clicks: data?.click_count || 0,
              urls: data?.url_count || 0,
              users: data?.unique_users || 0,
          }))
        : [];

    const pieData = [
        { name: 'Total Users', value: stats?.total_users || 0, color: '#3B82F6' },
        { name: 'Total URLs', value: stats?.total_urls || 0, color: '#10B981' },
        { name: 'Total Clicks', value: stats?.total_clicks || 0, color: '#F59E0B' },
    ];

    return (
        <div className="dashboard-wrapper">
            <div className="dashboard-header">
                <h1>Admin Dashboard</h1>
                <p>System overview and analytics for Minify URL shortener</p>
            </div>

            {error && <div className="dashboard-info">{error}</div>}

            <div className="stats-container">
                <div className="stats-card">
                    <h3>Total Users</h3>
                    <p>{stats?.total_users || 0}</p>
                </div>
                <div className="stats-card">
                    <h3>Total URLs</h3>
                    <p>{stats?.total_urls || 0}</p>
                </div>
                <div className="stats-card">
                    <h3>Total Clicks</h3>
                    <p>{stats?.total_clicks || 0}</p>
                </div>
                <div className="stats-card">
                    <h3>Avg Clicks/URL</h3>
                    <p>{stats?.total_urls ? Math.round((stats.total_clicks || 0) / stats.total_urls) : 0}</p>
                </div>
            </div>

            <div className="urls-section">
                <div className="chart-box">
                    <h3>Activity by Timeframe</h3>
                    <ResponsiveContainer width="100%" height={300}>
                        <BarChart data={timeframeData}>
                            <CartesianGrid strokeDasharray="3 3" />
                            <XAxis dataKey="period" />
                            <YAxis />
                            <Tooltip />
                            <Bar dataKey="clicks" fill="#3B82F6" name="Clicks" />
                            <Bar dataKey="urls" fill="#10B981" name="URLs" />
                            <Bar dataKey="users" fill="#F59E0B" name="Users" />
                        </BarChart>
                    </ResponsiveContainer>
                </div>

                <div className="chart-box">
                    <h3>System Overview</h3>
                    <ResponsiveContainer width="100%" height={300}>
                        <PieChart>
                            <Pie
                                data={pieData}
                                cx="50%"
                                cy="50%"
                                labelLine={false}
                                label={({ name, value }) => `${name}: ${value}`}
                                outerRadius={80}
                                dataKey="value"
                            >
                                {pieData.map((entry, index) => (
                                    <Cell key={`cell-${index}`} fill={entry.color} />
                                ))}
                            </Pie>
                            <Tooltip />
                        </PieChart>
                    </ResponsiveContainer>
                </div>
            </div>

            <div className="urls-section">
                <div className="output-box">
                    <h3>Popular URLs</h3>
                    {popularUrls.length === 0 ? (
                        <p>No URLs found</p>
                    ) : (
                        popularUrls.map((url, index) => (
                            <div key={url.short_code} className="popular-url-row">
                                <div>
                                    <p>#{index + 1} /{url.short_code}</p>
                                    <p title={url.original_url}>{url.original_url}</p>
                                    {url.username && <p>by {url.username}</p>}
                                </div>
                                <span>{url.clicks} clicks</span>
                            </div>
                        ))
                    )}
                </div>

                <div className="output-box">
                    <h3>Recent Users</h3>
                    {!stats?.recent_users || stats.recent_users.length === 0 ? (
                        <p>No recent users</p>
                    ) : (
                        stats.recent_users.map((username, index) => (
                            <div key={index} className="recent-user-row">
                                <div className="avatar">{username.charAt(0).toUpperCase()}</div>
                                <span>{username}</span>
                            </div>
                        ))
                    )}
                </div>
            </div>

            <button className="button-primary" onClick={fetchAnalytics}>
                Refresh Analytics
            </button>
        </div>
    );
};

export default AdminPage;
