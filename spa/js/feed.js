class FeedManager {
    constructor() {
        this.socket = null;
        this.reconnectTimer = null;
        this.shouldReconnect = true;
        this.isConnected = false;
        this.entries = [];
        this.limit = 10;
        this.eventTypeFilter = '';
    }

    async loadInitialFeed(limit = 10, eventType = this.eventTypeFilter) {
        try {
            const response = await api.getFeed(limit, 0, eventType);
            if (response.success) {
                this.limit = limit;
                this.eventTypeFilter = eventType;
                this.entries = Array.isArray(response.data) ? response.data : [];
                this.total = response.meta?.total || this.entries.length;
                this.renderFeed();
            }
        } catch (error) {
            console.error('Error loading activity feed:', error);
        }
    }

    async connect(limit = 10, eventType = '') {
        await this.loadInitialFeed(limit, eventType);
        this.shouldReconnect = true;
        this.openSocket();
    }

    disconnect() {
        this.shouldReconnect = false;
        if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;
        }
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }

    async updateLimit(limit) {
        this.limit = limit;
        await this.loadInitialFeed(limit, this.eventTypeFilter);
    }

    async updateEventTypeFilter(eventType) {
        this.eventTypeFilter = eventType;
        await this.loadInitialFeed(this.limit, eventType);
    }

    dedupeByID(entries) {
        const seen = new Set();
        return entries.filter((entry) => {
            if (seen.has(entry.id)) {
                return false;
            }
            seen.add(entry.id);
            return true;
        });
    }

    openSocket() {
        if (this.socket) {
            this.socket.close();
        }

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const url = `${protocol}//${window.location.host}/api/v1/feed/ws`;
        this.socket = new WebSocket(url);

        this.socket.onopen = () => {
            this.isConnected = true;
            this.updateConnectionStatus(true);
        };

        this.socket.onmessage = (event) => {
            try {
                const payload = JSON.parse(event.data);
                if (payload.success && payload.data) {
                    if (this.eventTypeFilter && payload.data.event_type !== this.eventTypeFilter) {
                        return;
                    }
                    this.entries.unshift(payload.data);
                    this.entries = this.dedupeByID(this.entries).slice(0, Math.max(this.limit, 100));
                    this.total = Math.max(this.total || 0, this.entries.length);
                    this.renderFeed();
                }
            } catch (error) {
                console.error('Error parsing feed update:', error);
            }
        };

        this.socket.onerror = () => {
            this.isConnected = false;
            this.updateConnectionStatus(false);
        };

        this.socket.onclose = () => {
            this.socket = null;
            this.isConnected = false;
            this.updateConnectionStatus(false);

            if (!this.shouldReconnect) {
                return;
            }

            if (this.reconnectTimer) {
                clearTimeout(this.reconnectTimer);
            }
            this.reconnectTimer = setTimeout(() => this.openSocket(), 3000);
        };
    }

    updateConnectionStatus(connected) {
        const statusEl = document.getElementById('connection-status');
        if (!statusEl) {
            return;
        }
        const indicator = statusEl.querySelector('span');
        const text = statusEl.querySelectorAll('span')[1];
        indicator.className = connected ? 'w-2 h-2 rounded-full bg-green-500 animate-pulse' : 'w-2 h-2 rounded-full bg-red-500';
        if (text) {
            text.textContent = connected ? 'Live' : 'Disconnected';
        }
    }

    renderFeed() {
        const bodyEl = document.getElementById('feed-body');
        const loadingEl = document.getElementById('feed-loading');
        const emptyEl = document.getElementById('feed-empty');
        const contentEl = document.getElementById('feed-content');
        const totalEl = document.getElementById('total-events');

        if (loadingEl) loadingEl.classList.add('hidden');
        if (totalEl) totalEl.textContent = this.total || this.entries.length;

        if (!bodyEl) {
            return;
        }

        const visibleEntries = this.entries.slice(0, this.limit);
        bodyEl.innerHTML = '';

        if (visibleEntries.length === 0) {
            if (emptyEl) emptyEl.classList.remove('hidden');
            if (contentEl) contentEl.classList.add('hidden');
            return;
        }

        if (emptyEl) emptyEl.classList.add('hidden');
        if (contentEl) contentEl.classList.remove('hidden');

        visibleEntries.forEach((entry) => {
            const row = document.createElement('tr');
            row.className = 'feed-row border-b border-purple-500/20';
            row.innerHTML = `
                <td class="py-4 px-4 text-sm text-slate-300 whitespace-nowrap">${new Date(entry.created_at).toLocaleTimeString()}</td>
                <td class="py-4 px-4 text-sm text-slate-200 font-medium">${entry.username || 'Unknown'}</td>
                <td class="py-4 px-4 text-sm text-purple-300">${entry.event_type}</td>
                <td class="py-4 px-4 text-sm text-slate-300">${entry.content}</td>
            `;
            bodyEl.appendChild(row);
        });
    }
}

const feedManager = new FeedManager();
