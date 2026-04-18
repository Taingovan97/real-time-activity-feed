# Continental Review Luxury Editorial Theme Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Redirect the SPA redesign from the stale light-dashboard pass to the approved Continental Review luxury-editorial theme while preserving all existing routes, feed behavior, auth behavior, and websocket updates.

**Architecture:** Keep the existing SPA structure centered in `spa/index.html` with behavior preserved in `spa/js/app.js` and `spa/js/feed.js`. Treat the current worktree branch as an intermediate state: replace the over-simplified light treatment with a navy-paper-ink system, stronger issue-cover framing, sharper typography, and classic table-first presentation without adding features or changing flow logic.

**Tech Stack:** Static HTML, inline CSS, Tailwind utility classes via CDN, vanilla JavaScript, WebSocket feed updates

---

## File Structure

- Modify: `spa/index.html`
  - Replace the current light-dashboard tokens and component styling with warm paper, deep navy, and ink tokens.
  - Retune header, publisher, feed, auth, profile, and 404 markup classes to match the approved luxury-editorial direction.
- Modify: `spa/js/feed.js`
  - Update runtime-generated feed rows and connection-state classes so live updates match the new table styling.
- Modify: `spa/js/app.js`
  - Only if required for class hooks or route-state compatibility after the HTML refactor.
- Reference: `docs/superpowers/specs/2026-04-18-swiss-editorial-light-theme-design.md`
- Execute in worktree: `.worktrees/swiss-editorial-light-theme`

## Task 1: Replace Base Tokens With Navy, Paper, And Ink

**Files:**
- Modify: `spa/index.html`
- Reference: `docs/superpowers/specs/2026-04-18-swiss-editorial-light-theme-design.md`

- [ ] **Step 1: Confirm the current worktree state before replacing theme primitives**

Run:

```bash
git status --short
```

Expected:

```text
Only intended frontend files are modified in the worktree.
```

- [ ] **Step 2: Replace the current simplified light tokens with the Continental Review token system**

In `spa/index.html`, replace the current theme primitives in the inline `<style>` block with:

```html
<style>
    :root {
        --paper: #f3ede2;
        --paper-soft: #ede5d8;
        --paper-strong: #fffaf0;
        --ink: #142235;
        --ink-soft: #4d5d72;
        --ink-muted: #718198;
        --navy: #152338;
        --navy-soft: #1d2d45;
        --navy-strong: #0f1b2b;
        --blue: #2d5baf;
        --blue-strong: #1d417f;
        --blue-soft: #d8e5f7;
        --rule: rgba(20, 34, 53, 0.14);
        --rule-strong: rgba(20, 34, 53, 0.24);
        --success: #14805e;
        --error: #b6534d;
        --shadow-soft: 0 18px 45px rgba(15, 27, 43, 0.08);
        --shadow-lift: 0 24px 60px rgba(15, 27, 43, 0.14);
    }

    * { font-family: 'Manrope', 'Segoe UI', sans-serif; }
    h1, h2, h3, .editorial-display {
        font-family: 'Cormorant Garamond', Georgia, serif;
        letter-spacing: 0.01em;
    }

    body {
        background:
            radial-gradient(circle at top left, rgba(29, 65, 127, 0.08), transparent 24%),
            linear-gradient(180deg, #f7f1e8 0%, var(--paper) 100%);
        color: var(--ink);
    }

    .continental-shell {
        min-height: 100vh;
        background: transparent;
    }

    .continental-frame {
        background: linear-gradient(180deg, var(--navy) 0%, var(--navy-soft) 100%);
        color: var(--paper-strong);
        box-shadow: var(--shadow-lift);
    }

    .continental-panel {
        background: rgba(255, 250, 240, 0.92);
        border: 1px solid rgba(20, 34, 53, 0.1);
        box-shadow: var(--shadow-soft);
    }

    .continental-panel-inset {
        background: linear-gradient(180deg, rgba(255,250,240,0.96), rgba(243,237,226,0.92));
        border: 1px solid rgba(20, 34, 53, 0.12);
    }

    .continental-rule {
        border-color: var(--rule);
    }

    .continental-input {
        background: rgba(255, 250, 240, 0.96);
        color: var(--ink);
        border: 1px solid var(--rule-strong);
    }

    .continental-input:focus {
        outline: none;
        border-color: rgba(45, 91, 175, 0.48);
        box-shadow: 0 0 0 4px rgba(45, 91, 175, 0.12);
    }

    .continental-btn {
        background: linear-gradient(180deg, var(--blue) 0%, var(--blue-strong) 100%);
        color: #fefaf3;
        border: 1px solid rgba(20, 34, 53, 0.12);
        transition: transform 180ms ease, box-shadow 180ms ease, background 180ms ease;
        box-shadow: 0 12px 28px rgba(29, 65, 127, 0.24);
    }

    .continental-btn:hover {
        transform: translateY(-1px);
        background: linear-gradient(180deg, #244f9d 0%, #173668 100%);
        box-shadow: 0 16px 36px rgba(29, 65, 127, 0.28);
    }

    .continental-btn-secondary {
        background: transparent;
        color: var(--ink);
        border: 1px solid var(--rule-strong);
        transition: background 180ms ease, border-color 180ms ease, color 180ms ease;
    }

    .continental-btn-secondary:hover {
        background: rgba(216, 229, 247, 0.36);
        border-color: rgba(45, 91, 175, 0.26);
        color: var(--blue-strong);
    }

    .continental-kicker {
        font-size: 0.72rem;
        letter-spacing: 0.18em;
        text-transform: uppercase;
        color: var(--ink-muted);
    }

    .feed-row {
        transition: background 180ms ease;
    }

    .feed-row:hover {
        background: rgba(216, 229, 247, 0.22);
    }

    .type-chip {
        display: inline-flex;
        align-items: center;
        padding: 0.22rem 0.72rem;
        border-radius: 999px;
        background: rgba(216, 229, 247, 0.58);
        border: 1px solid rgba(29, 65, 127, 0.12);
        color: var(--blue-strong);
        font-size: 0.78rem;
        font-weight: 600;
    }
 </style>
```

- [ ] **Step 3: Keep the editorial font pairing and remove leftover generic dashboard utility assumptions**

In `spa/index.html`, ensure the Google Fonts import remains:

```html
<link href="https://fonts.googleapis.com/css2?family=Manrope:wght@400;500;600;700;800&family=Cormorant+Garamond:wght@500;600;700&display=swap" rel="stylesheet">
```

And replace stale custom class names such as `editorial-panel`, `editorial-input`, and `editorial-btn` with the new `continental-*` names throughout the file.

- [ ] **Step 4: Run a smoke build after the token refactor**

Run:

```bash
go test ./cmd/server/... ./internal/...
```

Expected:

```text
PASS
```

- [ ] **Step 5: Commit the token and base-class replacement**

```bash
git add spa/index.html
git commit -m "feat: establish continental review theme tokens"
```

## Task 2: Rebuild The Header And Feed Screen As An Issue Cover

**Files:**
- Modify: `spa/index.html`

- [ ] **Step 1: Create a failing checklist for the feed route**

Use this acceptance check before edits:

```text
FAIL if the feed screen still looks like:
- a thin app bar instead of a strong masthead
- a generic white dashboard with evenly rounded cards
- soft product pills dominating the controls
- a feed table that lacks navy framing and editorial hierarchy
```

- [ ] **Step 2: Replace the current header with an issue-cover masthead**

In `spa/index.html`, refactor the top shell and header markup toward:

```html
<body class="continental-shell text-slate-900 min-h-screen">
    <div id="feed-container" class="min-h-screen">
        <header class="continental-frame relative overflow-hidden border-b border-[rgba(255,250,240,0.12)]">
            <div class="max-w-7xl mx-auto px-5 lg:px-8 py-8 md:py-12">
                <div class="flex flex-col gap-8 md:flex-row md:items-start md:justify-between">
                    <div class="max-w-3xl">
                        <p class="continental-kicker text-[rgba(255,250,240,0.62)] mb-3">Shared activity stream</p>
                        <h1 class="editorial-display text-[3rem] md:text-[4.4rem] leading-[0.9] text-[var(--paper-strong)]">Real-Time Activity Feed</h1>
                        <p class="mt-4 max-w-2xl text-[rgba(255,250,240,0.72)] text-base md:text-lg leading-7">
                            A live operations briefing presented as a continuous editorial ledger.
                        </p>
                    </div>
                    <div class="flex items-start gap-4">
                        <div id="connection-status" class="border border-[rgba(255,250,240,0.18)] bg-[rgba(255,250,240,0.06)] px-4 py-2 rounded-full">
                            <div class="flex items-center gap-2">
                                <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
                                <span class="text-sm text-[rgba(255,250,240,0.82)]">Live</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </header>
```

- [ ] **Step 3: Refactor the publisher and feed section into stronger editorial framing**

In `spa/index.html`, change the feed route body to:

```html
<main class="max-w-7xl mx-auto px-5 lg:px-8 py-8 md:py-12">
    <section id="event-publisher-prompt" class="continental-frame rounded-[34px] p-8 md:p-10 mb-10">
        <div class="flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
            <div class="max-w-2xl">
                <p class="continental-kicker text-[rgba(255,250,240,0.62)] mb-3">Contribute</p>
                <h2 class="editorial-display text-[2.4rem] md:text-[3rem] leading-[0.94] text-[var(--paper-strong)]">Publish a new event</h2>
                <p class="mt-3 text-[rgba(255,250,240,0.74)] leading-7">
                    Sign in to add decisions, approvals, incidents, and release notes to the live issue.
                </p>
            </div>
            <button id="publish-event-btn" class="continental-btn font-semibold py-3.5 px-7 rounded-full">Sign In</button>
        </div>
    </section>

    <section class="continental-panel rounded-[34px] p-7 md:p-9 lg:p-10">
        <div class="grid gap-8 lg:grid-cols-[minmax(0,1.3fr)_auto] lg:items-end mb-8">
            <div class="max-w-2xl">
                <p class="continental-kicker mb-3">Feature spread</p>
                <h2 class="editorial-display text-[2.6rem] md:text-[3.4rem] leading-[0.92] text-[var(--ink)]">Latest Feed</h2>
                <p class="mt-3 text-[var(--ink-soft)] text-base md:text-lg leading-7">
                    Recent actions and notifications across the system, presented as a classic briefing table.
                </p>
            </div>
            <div class="flex flex-wrap items-center gap-3 lg:justify-end">
                <select id="feed-event-type-filter" class="continental-input text-sm rounded-full px-4 py-2.5 min-w-[160px]">
                    <option value="">All Types</option>
                </select>
                <select id="feed-limit" class="continental-input text-sm rounded-full px-4 py-2.5 min-w-[130px]">
                    <option value="5">Last 5</option>
                    <option value="10" selected>Last 10</option>
                    <option value="50">Last 50</option>
                    <option value="100">Last 100</option>
                </select>
            </div>
        </div>
```

- [ ] **Step 4: Refine the search row and table wrapper without switching away from table-first**

In `spa/index.html`, update the search form and table wrapper to:

```html
<form id="feed-search-form" class="mb-8 border-t border-b continental-rule py-5 flex flex-col gap-4 lg:flex-row lg:items-center">
    <input id="feed-search-query" type="text" class="continental-input flex-1 px-5 py-3.5 rounded-full" placeholder="Search by user or message">
    <div class="flex gap-3">
        <button type="submit" class="continental-btn font-semibold py-3 px-6 rounded-full">Search</button>
        <button type="button" id="feed-search-clear" class="continental-btn-secondary font-semibold py-3 px-6 rounded-full">Clear</button>
    </div>
</form>

<div id="feed-content" class="hidden overflow-x-auto rounded-[26px] border border-[rgba(20,34,53,0.12)] bg-[rgba(255,250,240,0.78)]">
    <table class="w-full min-w-[760px]">
        <thead>
            <tr class="border-b border-[rgba(20,34,53,0.12)]">
                <th class="text-left py-4 px-4 text-xs font-semibold tracking-[0.18em] uppercase text-[var(--ink-muted)]">Time</th>
                <th class="text-left py-4 px-4 text-xs font-semibold tracking-[0.18em] uppercase text-[var(--ink-muted)]">User</th>
                <th class="text-left py-4 px-4 text-xs font-semibold tracking-[0.18em] uppercase text-[var(--ink-muted)]">Type</th>
                <th class="text-left py-4 px-4 text-xs font-semibold tracking-[0.18em] uppercase text-[var(--ink-muted)]">Message</th>
            </tr>
        </thead>
        <tbody id="feed-body"></tbody>
    </table>
</div>
```

- [ ] **Step 5: Verify and commit the feed-page refactor**

Run:

```bash
make start-dev
```

Expected manual result:

```text
- the header reads as a deep navy issue-cover band
- the page feels warmer and more dramatic than the stale light-dashboard pass
- the feed remains a classic table, not cards
- controls feel tailored and no longer look like generic SaaS pills
```

Then commit:

```bash
git add spa/index.html
git commit -m "feat: rebuild feed page as continental review layout"
```

## Task 3: Restyle Auth, Profile, And 404 Into The Same Editorial System

**Files:**
- Modify: `spa/index.html`

- [ ] **Step 1: Create a failing checklist for non-feed routes**

Use this acceptance check:

```text
FAIL if login, register, profile, or 404 still look like:
- the stale light-dashboard pass
- soft generic white cards without navy framing
- old purple accents or old editorial class names
- profile actions with weak CTA emphasis
```

- [ ] **Step 2: Convert login and register into navy-and-paper editorial forms**

In `spa/index.html`, refactor auth route markup toward:

```html
<div id="auth-container" class="hidden min-h-screen flex items-center justify-center px-4 py-12">
    <div id="login-form" class="w-full max-w-2xl">
        <div class="continental-frame rounded-[34px] p-8 md:p-10 lg:p-12">
            <p class="continental-kicker text-[rgba(255,250,240,0.62)] text-center mb-3">Editorial access</p>
            <h1 class="editorial-display text-[3.4rem] md:text-[4rem] text-center leading-[0.92] text-[var(--paper-strong)]">Welcome Back</h1>
            <p class="mt-4 text-center text-[rgba(255,250,240,0.72)] leading-7">
                Sign in to publish updates and contribute to the live operations issue.
            </p>
            <form id="login-form-element" class="mt-8 space-y-5">
                <input type="text" id="login-username" required class="continental-input w-full px-5 py-3.5 rounded-full" placeholder="Username">
                <input type="password" id="login-password" required class="continental-input w-full px-5 py-3.5 rounded-full" placeholder="Password">
                <div id="login-error" class="hidden text-sm text-rose-200"></div>
                <button type="submit" class="continental-btn w-full font-semibold py-3.5 px-4 rounded-full">Sign In</button>
            </form>
            <button id="auth-back-to-feed" class="continental-btn-secondary w-full font-semibold py-3.5 px-4 rounded-full mt-5 text-[var(--paper-strong)] border-[rgba(255,250,240,0.26)]">Back to Feed</button>
        </div>
    </div>
</div>
```

- [ ] **Step 3: Bring profile and 404 into the same high-contrast system**

In `spa/index.html`, restyle the profile route and 404 surface so:

```html
<div id="not-found-container" class="hidden min-h-screen flex items-center justify-center px-4 py-12">
    <div class="continental-frame rounded-[34px] p-10 md:p-12 text-center max-w-2xl">
        <p class="continental-kicker text-[rgba(255,250,240,0.62)] mb-3">Page state</p>
        <h1 class="editorial-display text-[5rem] md:text-[6rem] leading-none text-[var(--paper-strong)]">404</h1>
        <p class="mt-4 text-[rgba(255,250,240,0.74)] leading-7">
            The page you requested is not part of this current issue.
        </p>
        <button id="go-home-btn" class="continental-btn font-semibold py-3 px-6 rounded-full mt-8">Go to Home</button>
    </div>
</div>
```

And convert the profile stats, avatar area, and publish form to use `continental-frame`, `continental-panel`, `continental-panel-inset`, `continental-input`, `continental-btn`, and `continental-btn-secondary` rather than the stale simplified classes.

- [ ] **Step 4: Verify route consistency**

Run:

```bash
make start-dev
```

Expected manual result:

```text
- /login and /register feel like premium editorial forms
- /profile matches the navy-paper-ink system and keeps strong CTA emphasis
- /404 shares the same issue-cover tone
- no old purple or stale light-dashboard styling remains
```

- [ ] **Step 5: Commit the route refactor**

```bash
git add spa/index.html
git commit -m "feat: restyle secondary routes for continental review"
```

## Task 4: Align Runtime Feed Rendering With The New Table Style

**Files:**
- Modify: `spa/js/feed.js`
- Modify: `spa/js/app.js` only if route-state hooks need a class-name adjustment

- [ ] **Step 1: Create the failing runtime checklist**

Use this acceptance check:

```text
FAIL if live or initial feed rows still render with:
- stale class names from the previous pass
- weak contrast against the warm-paper table background
- chips or status text that break the blue-and-ink-only direction
```

- [ ] **Step 2: Replace runtime row classes with Continental Review table styling**

In `spa/js/feed.js`, update `renderFeed()` to use:

```javascript
visibleEntries.forEach((entry) => {
    const row = document.createElement('tr');
    row.className = 'feed-row border-b border-[rgba(20,34,53,0.12)] align-top';
    row.innerHTML = `
        <td class="py-4 px-4 text-sm text-[var(--ink-muted)] whitespace-nowrap">${new Date(entry.created_at).toLocaleTimeString()}</td>
        <td class="py-4 px-4 text-sm text-[var(--ink)] font-semibold">${entry.username || 'Unknown'}</td>
        <td class="py-4 px-4 text-sm"><span class="type-chip">${entry.event_type}</span></td>
        <td class="py-4 px-4 text-sm text-[var(--ink-soft)] leading-6">${entry.content}</td>
    `;
    bodyEl.appendChild(row);
});
```

- [ ] **Step 3: Verify connection-state styling still fits the navy masthead**

If needed, normalize `updateConnectionStatus()` in `spa/js/feed.js` to:

```javascript
indicator.className = connected ? 'w-2 h-2 rounded-full bg-emerald-400 animate-pulse' : 'w-2 h-2 rounded-full bg-rose-400';
if (text) {
    text.className = connected ? 'text-sm text-[rgba(255,250,240,0.82)]' : 'text-sm text-[rgba(255,250,240,0.62)]';
    text.textContent = connected ? 'Live' : 'Disconnected';
}
```

- [ ] **Step 4: Verify live updates manually**

Run:

```bash
make start-dev
```

Expected manual result:

```text
- initial rows and websocket-inserted rows match exactly
- the type chip remains restrained and readable
- connection status remains legible in the navy header
```

- [ ] **Step 5: Commit the runtime alignment**

```bash
git add spa/js/feed.js spa/js/app.js
git commit -m "feat: align live feed rendering with continental review theme"
```

## Task 5: Final Responsive Pass, Cleanup, And Verification

**Files:**
- Modify: `spa/index.html`
- Modify: `spa/js/feed.js` only if small visual-hook cleanup is needed
- Modify: `spa/js/app.js` only if small route styling cleanup is needed

- [ ] **Step 1: Create the last failing acceptance checklist**

Use this checklist:

```text
FAIL if:
- desktop still feels too flat or too even after the redesign
- mobile spacing causes the issue-cover header or controls to break awkwardly
- any route still leaks stale class names from the previous pass
- empty, loading, disconnected, success, or error states lose readability
```

- [ ] **Step 2: Apply final targeted polish only where needed**

Use focused adjustments such as:

```html
<main class="max-w-7xl mx-auto px-4 sm:px-5 lg:px-8 py-8 md:py-12">
```

```html
<div id="feed-loading" class="text-center py-14">
    <p class="mt-4 text-[var(--ink-soft)]">Loading activity feed...</p>
</div>
```

```html
<div id="feed-empty" class="hidden text-center py-14">
    <p class="text-[var(--ink-soft)] text-base">No events yet. Publish the first event.</p>
</div>
```

Only keep polish that reinforces the approved direction. Do not add features, toggles, or new flows.

- [ ] **Step 3: Run verification commands**

Run:

```bash
go test ./cmd/server/... ./internal/...
make lint
```

Expected:

```text
All Go tests pass.
Linter completes successfully.
```

- [ ] **Step 4: Run the full manual route verification**

Run:

```bash
make start-dev
```

Expected manual result:

```text
- feed, login, register, profile, and 404 all match the Continental Review direction
- the feed remains classic table-first
- the header reads like an issue cover
- search, filters, publish forms, auth flow, and websocket updates still behave correctly
- desktop and mobile both stay presentation-friendly
```

- [ ] **Step 5: Commit the final luxury-editorial pass**

```bash
git add spa/index.html spa/js/feed.js spa/js/app.js
git commit -m "feat: finalize continental review luxury theme"
```

## Self-Review

### Spec Coverage

- Warm paper, deep navy, ink-only color system: covered in Task 1
- Issue-cover header and stronger page framing: covered in Task 2
- Classic table-first feed treatment: covered in Tasks 2 and 4
- Luxury-editorial auth, profile, and 404 routes: covered in Task 3
- Reduced motion, responsive layout, and behavior preservation: covered in Task 5

No spec gaps remain.

### Placeholder Scan

- No `TODO`, `TBD`, or deferred implementation placeholders remain.
- Each task names exact files.
- Code-changing steps include concrete snippets.
- Verification steps include exact commands and expected outcomes.

### Type Consistency

- Theme classes are consistently named `continental-*`.
- Display typography hook remains `editorial-display`.
- Runtime chip class remains `type-chip`.
- Feed row class remains `feed-row`.
