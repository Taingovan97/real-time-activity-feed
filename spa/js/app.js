class App {
    constructor() {
        this.init();
    }

    init() {
        this.setupRouting();
        this.bindEvents();
        const limitSelector = document.getElementById('feed-limit');
        if (limitSelector) {
            limitSelector.value = String(this.getFeedLimit());
        }
        feedManager.connect(this.getFeedLimit());
        this.updateAuthUI();
    }

    setupRouting() {
        this.handleRoute();
        window.addEventListener('popstate', () => this.handleRoute());
    }

    bindEvents() {
        document.getElementById('feed-form')?.addEventListener('submit', (e) => this.handleEventSubmit(e, false));
        document.getElementById('profile-feed-form')?.addEventListener('submit', (e) => this.handleEventSubmit(e, true));
        document.getElementById('profile-back-btn')?.addEventListener('click', () => this.navigate('/'));
        document.getElementById('auth-back-to-feed')?.addEventListener('click', () => this.navigate('/'));
        document.getElementById('auth-back-to-feed-register')?.addEventListener('click', () => this.navigate('/'));
        document.getElementById('publish-event-btn')?.addEventListener('click', () => this.navigate('/login'));
        document.getElementById('go-home-btn')?.addEventListener('click', () => this.navigate('/'));
        document.getElementById('profile-link')?.addEventListener('click', (e) => {
            e.preventDefault();
            this.navigate('/profile');
        });
        document.getElementById('feed-limit')?.addEventListener('change', async (e) => {
            const limit = parseInt(e.target.value, 10);
            this.setFeedLimit(limit);
            await feedManager.updateLimit(limit);
        });
    }

    navigate(path) {
        window.history.pushState({}, '', path);
        this.handleRoute();
    }

    handleRoute() {
        const path = window.location.pathname;
        if (path === '/login') return this.showLogin();
        if (path === '/register') return this.showRegister();
        if (path === '/profile') return this.showProfile();
        if (path !== '/') return this.show404();
        this.showFeed();
    }

    showLogin() {
        this.toggleView('auth-container');
        document.getElementById('login-form')?.classList.remove('hidden');
        document.getElementById('register-form')?.classList.add('hidden');
    }

    showRegister() {
        this.toggleView('auth-container');
        document.getElementById('login-form')?.classList.add('hidden');
        document.getElementById('register-form')?.classList.remove('hidden');
    }

    showFeed() {
        this.toggleView('feed-container');
        this.updateAuthUI();
    }

    async showProfile() {
        if (!localStorage.getItem('accessToken')) {
            return this.navigate('/login');
        }

        this.toggleView('profile-container');
        await this.loadProfileData();
    }

    show404() {
        this.toggleView('not-found-container');
    }

    toggleView(activeID) {
        ['auth-container', 'feed-container', 'profile-container', 'not-found-container'].forEach((id) => {
            document.getElementById(id)?.classList.toggle('hidden', id !== activeID);
        });
    }

    async loadProfileData() {
        const currentUser = authManager.getCurrentUser();
        if (!currentUser) {
            return;
        }

        document.getElementById('profile-avatar-text').textContent = currentUser.username?.charAt(0).toUpperCase() || '';
        document.getElementById('profile-username').textContent = currentUser.username || '-';
        document.getElementById('profile-email').textContent = currentUser.email || '-';

        try {
            const response = await api.getFeed(100, 0);
            const entries = Array.isArray(response.data) ? response.data.filter((entry) => entry.user_id === currentUser.id) : [];
            document.getElementById('profile-latest-event-type').textContent = entries[0]?.event_type || '-';
            document.getElementById('profile-latest-event-time').textContent = entries[0] ? new Date(entries[0].created_at).toLocaleString() : '-';
            document.getElementById('profile-total-events').textContent = entries.length;
        } catch {
            document.getElementById('profile-latest-event-type').textContent = '-';
            document.getElementById('profile-latest-event-time').textContent = '-';
            document.getElementById('profile-total-events').textContent = '0';
        }
    }

    async handleEventSubmit(e, profileMode) {
        e.preventDefault();

        if (!localStorage.getItem('accessToken')) {
            return this.navigate('/login');
        }

        const typeInput = document.getElementById(profileMode ? 'profile-event-type-input' : 'event-type-input');
        const contentInput = document.getElementById(profileMode ? 'profile-event-content-input' : 'event-content-input');
        const errorDiv = document.getElementById(profileMode ? 'profile-event-error' : 'event-error');
        const successDiv = document.getElementById(profileMode ? 'profile-event-success' : 'event-success');

        errorDiv.classList.add('hidden');
        successDiv.classList.add('hidden');

        try {
            const response = await api.publishEvent(typeInput.value.trim(), contentInput.value.trim());
            if (!response.success) {
                throw new Error(response.message || 'Failed to publish event');
            }

            successDiv.textContent = 'Event published successfully.';
            successDiv.classList.remove('hidden');
            typeInput.value = '';
            contentInput.value = '';
            await feedManager.loadInitialFeed(this.getFeedLimit());
            if (profileMode) {
                await this.loadProfileData();
            }
        } catch (error) {
            errorDiv.textContent = error.message || 'Failed to publish event.';
            errorDiv.classList.remove('hidden');
        }
    }

    updateAuthUI() {
        const isAuthenticated = !!localStorage.getItem('accessToken');
        const currentUser = authManager.getCurrentUser();

        document.getElementById('event-publisher-card')?.classList.toggle('hidden', !isAuthenticated);
        document.getElementById('event-publisher-prompt')?.classList.toggle('hidden', isAuthenticated);
        document.getElementById('user-dropdown')?.classList.toggle('hidden', !isAuthenticated);

        if (currentUser) {
            document.getElementById('user-name').textContent = currentUser.username || '';
            document.getElementById('user-avatar-text').textContent = currentUser.username?.charAt(0).toUpperCase() || '';
            document.getElementById('dropdown-user-name').textContent = currentUser.username || '';
            document.getElementById('dropdown-user-email').textContent = currentUser.email || '';
        }
    }

    getFeedLimit() {
        return parseInt(localStorage.getItem('feedLimit') || '10', 10);
    }

    setFeedLimit(limit) {
        localStorage.setItem('feedLimit', String(limit));
    }
}

let app;
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        app = new App();
    });
} else {
    app = new App();
}
