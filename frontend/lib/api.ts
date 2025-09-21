import axios from 'axios';
import Cookies from 'js-cookie';
import {
    MinifyRequest,
    MinifyResponse,
    LoginRequest,
    LoginResponse,
    CreateUserRequest,
    User,
    URL,
    OverviewStats,
    PopularURL,
    TimeframeStats,
} from '../types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: true,
});

// add token to requests if available
api.interceptors.request.use((config) => {
    const token = Cookies.get('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// handle token expiration
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            Cookies.remove('token');
            Cookies.remove('user');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

// user auth
export const authAPI = {
    login: async (data: LoginRequest): Promise<LoginResponse> => {
        const response = await api.post('/api/v1/users/login', data);
        return response.data;
    },

    register: async (data: CreateUserRequest): Promise<User> => {
        const response = await api.post('/api/v1/users', data);
        return response.data;
    },

    logout: () => {
        Cookies.remove('token');
        Cookies.remove('user');
    },

    getCurrentUser: (): User | null => {
        const userCookie = Cookies.get('user');
        return userCookie ? JSON.parse(userCookie) : null;
    },

    isAuthenticated: (): boolean => {
        return !!Cookies.get('token');
    },

    // TODO: make this more strict, especially if admins get modify access for user data
    // not optimal, just set like this for easier testing
    isAdmin: (): boolean => {
        const user = authAPI.getCurrentUser();
        if (!user) return false;
        return user.username === 'admin' || user.email.includes('admin');
    }
};

// URL 
export const urlAPI = {
    minify: async (data: MinifyRequest): Promise<MinifyResponse> => {
        const response = await api.post('/api/v1/minify', data);
        return response.data;
    },

    getUserURLs: async (userId: number): Promise<URL[]> => {
        const response = await api.get(`/api/v1/urls?user_id=${userId}`);
        return response.data;
    },
};

// analytics 
export const analyticsAPI = {
    getOverview: async (): Promise<OverviewStats> => {
        const response = await api.get('/api/v1/analytics/overview');
        return response.data;
    },

    getPopularURLs: async (limit = 10): Promise<PopularURL[]> => {
        const response = await api.get(`/api/v1/analytics/popular?limit=${limit}`);
        return response.data;
    },

    getTimeframeStats: async (period: string): Promise<TimeframeStats> => {
        const response = await api.get(`/api/v1/analytics/timeframe/${period}`);
        return response.data;
    },
};

export default api;
