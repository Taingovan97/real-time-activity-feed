const API_BASE_URL = '/api/v1';

class API {
    constructor() {
        this.accessToken = localStorage.getItem('accessToken');
        this.refreshToken = localStorage.getItem('refreshToken');
        this.refreshPromise = null;
        this.refreshBufferMinutes = 5;
    }

    updateTokens() {
        this.accessToken = localStorage.getItem('accessToken');
        this.refreshToken = localStorage.getItem('refreshToken');
    }

    async request(endpoint, options = {}) {
        this.updateTokens();

        if (this.accessToken && shouldRefreshToken(this.accessToken, this.refreshBufferMinutes)) {
            const refreshed = await this.ensureValidToken();
            if (!refreshed) {
                this.clearTokens();
                throw new Error('Authentication failed. Please log in again.');
            }
        }

        if (this.refreshPromise) {
            await this.refreshPromise;
        }

        const headers = {
            'Content-Type': 'application/json',
            ...options.headers,
        };

        if (this.accessToken) {
            headers.Authorization = `Bearer ${this.accessToken}`;
        }

        const response = await fetch(`${API_BASE_URL}${endpoint}`, { ...options, headers });
        const data = await response.json();

        if (!response.ok) {
            if (response.status === 401 && this.refreshToken) {
                const refreshed = await this.ensureValidToken();
                if (refreshed) {
                    return this.request(endpoint, options);
                }
                this.clearTokens();
            }
            throw new Error(data.message || data.error?.message || 'Request failed');
        }

        return data;
    }

    async ensureValidToken() {
        if (this.refreshPromise) {
            try {
                await this.refreshPromise;
                return true;
            } catch {
                return false;
            }
        }

        if (!this.refreshToken || isTokenExpired(this.refreshToken)) {
            return false;
        }

        if (this.accessToken && !isTokenExpired(this.accessToken) && !shouldRefreshToken(this.accessToken, this.refreshBufferMinutes)) {
            return true;
        }

        this.refreshPromise = this.refreshAccessToken();
        try {
            return await this.refreshPromise;
        } finally {
            this.refreshPromise = null;
        }
    }

    async refreshAccessToken() {
        const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refresh_token: this.refreshToken }),
        });

        const data = await response.json();
        if (!response.ok || !data.data?.token) {
            return false;
        }

        this.setTokens(data.data.token.access_token, data.data.token.refresh_token);
        if (typeof authManager !== 'undefined' && authManager.onTokenRefreshed) {
            authManager.onTokenRefreshed().catch(() => {});
        }

        return true;
    }

    setTokens(accessToken, refreshToken) {
        this.accessToken = accessToken;
        this.refreshToken = refreshToken;
        localStorage.setItem('accessToken', accessToken);
        localStorage.setItem('refreshToken', refreshToken);
    }

    clearTokens() {
        this.accessToken = null;
        this.refreshToken = null;
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
    }

    async register(username, email, password) {
        return this.request('/auth/register', {
            method: 'POST',
            body: JSON.stringify({ username, email, password }),
        });
    }

    async login(username, password) {
        return this.request('/auth/login', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
        });
    }

    async getFeed(limit = 10, offset = 0, eventType = '', query = '') {
        const params = new URLSearchParams({
            limit: String(limit),
            offset: String(offset),
        });
        if (eventType) {
            params.set('event_type', eventType);
        }
        if (query) {
            params.set('query', query);
        }
        return this.request(`/feed?${params.toString()}`, { method: 'GET' });
    }

    async publishEvent(eventType, content) {
        return this.request('/events', {
            method: 'POST',
            body: JSON.stringify({ event_type: eventType, content }),
        });
    }

    async getEventTypes() {
        return this.request('/feed/event-types', { method: 'GET' });
    }
}

const api = new API();
