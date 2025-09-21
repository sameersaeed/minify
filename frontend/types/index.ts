export interface User {
  id: number;
  username: string;
  email: string;
  created_at: string;
}

export interface URL {
  id: number;
  short_code: string;
  original_url: string;
  user_id?: number;
  clicks: number;
  created_at: string;
  updated_at: string;
}

export interface MinifyRequest {
  url: string;
  user_id?: number;
}

export interface MinifyResponse {
  short_url: string;
  original_url: string;
  short_code: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
}

export interface OverviewStats {
  total_users: number;
  total_urls: number;
  total_clicks: number;
  recent_users: string[];
  timeframe_data: {
    hour?: TimeframeStats;
    day?: TimeframeStats;
    month?: TimeframeStats;
    year?: TimeframeStats;
  };
}

export interface PopularURL {
  short_code: string;
  original_url: string;
  clicks: number;
  username?: string;
}

export interface TimeframeStats {
  period: string;
  click_count: number;
  url_count: number;
  unique_users: number;
}
