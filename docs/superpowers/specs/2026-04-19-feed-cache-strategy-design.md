# Feed Cache Strategy Design

## Goal

Introduce a practical cache strategy that improves scalability for feed reads without weakening correctness, complicating auth, or turning Redis into an unbounded query cache.

## Current State

- PostgreSQL is the source of truth for persisted data.
- Redis is currently used only for pub/sub fan-out of live feed deltas.
- `GET /api/v1/feed` reads directly from PostgreSQL for every request.
- `GET /api/v1/feed/ws` is pub/sub only and does not read cached feed history.
- Auth endpoints, including `GET /api/v1/auth/me`, read directly from the database-backed auth path.

This is a sound baseline, but it means repeated hot feed reads scale primarily with PostgreSQL.

## Design Principles

- PostgreSQL remains authoritative for feed history.
- Redis caching is additive and optional, not required for correctness.
- Cache only high-reuse views with predictable invalidation.
- Avoid caching arbitrary search queries in the first iteration.
- Keep cache responsibilities separate from Redis pub/sub broadcasting.

## Recommended Approach

Add a dedicated feed read cache for a narrow set of hot queries behind the feed read path.

### Scope

Cache only:

- Unfiltered first-page feed reads
- First-page feed reads filtered by `event_type`
- Optionally a very small number of additional early pages if metrics justify it

Do not cache initially:

- Arbitrary text search queries
- Deep pagination
- Auth endpoints
- Token validation or refresh flows
- WebSocket broadcast payload delivery

### Cache Pattern

Use cache-aside semantics on the feed read path:

1. The feed use case receives a read request.
2. If the request is cacheable, read Redis first.
3. On cache hit, return the cached payload.
4. On cache miss, query PostgreSQL.
5. Store the serialized result in Redis with a short TTL.
6. Return the database result.

If Redis is unavailable, reads fall back to PostgreSQL and still succeed.

## Cacheability Rules

Initial cacheability rules should be intentionally narrow:

- `query == ""`
- `offset == 0`
- `limit` must match an approved small set, such as the SPA default
- `event_type` may be empty or one valid fixed event type

This keeps key cardinality bounded and makes invalidation cheap.

## Cache Keys

Use explicit versioned keys so the format can evolve safely:

- `feed:v1:page:all:limit=<n>:offset=<n>`
- `feed:v1:page:type=<eventType>:limit=<n>:offset=<n>`

The cached value should contain:

- `entries`
- `total`

This mirrors the existing `GetFeed` result and avoids hidden reconstruction logic.

## TTL Strategy

Use short TTLs as a safety backstop rather than the primary consistency mechanism.

Recommended starting TTL:

- 15 to 60 seconds for feed pages

Short TTLs provide:

- resilience if invalidation is missed
- bounded staleness
- simple recovery from deploy-time key mismatches

## Invalidation Strategy

Invalidate hot feed keys after a successful event publish.

When `PublishEvent` persists a new event successfully:

- invalidate the unfiltered first-page key
- invalidate the first-page key for the published `event_type`

If later pages are cached in a future iteration, invalidate only the small, explicitly supported cached range rather than trying to update every cached page in place.

The design should prefer delete-and-rebuild over partial in-place mutation of cached pages. That keeps correctness simple and avoids ordering bugs.

## Architecture Changes

Add a feed cache abstraction owned by the feed module application layer.

Suggested interface shape:

```go
type FeedCache interface {
    GetFeed(ctx context.Context, eventType string, limit, offset int64) ([]domain.FeedEvent, int64, bool, error)
    SetFeed(ctx context.Context, eventType string, limit, offset int64, entries []domain.FeedEvent, total int64, ttl time.Duration) error
    InvalidateAfterPublish(ctx context.Context, eventType string) error
}
```

Notes:

- Keep cache concerns separate from `BroadcastService`.
- Inject the cache dependency into the feed read and publish use cases.
- Allow a no-op implementation so caching can be disabled without changing behavior.

## Module Responsibilities

### Application Layer

- decide whether a request is cacheable
- orchestrate cache read, DB fallback, cache write
- trigger invalidation after successful publish

### Infrastructure Layer

- implement Redis-backed `FeedCache`
- serialize and deserialize cached feed payloads
- generate Redis keys
- apply TTL and deletion semantics

### Adapters Layer

- unchanged API contract
- no cache-specific behavior in handlers

## Failure Behavior

Caching must not make the system less available.

- Redis read failure: log and fall back to PostgreSQL
- Redis write failure: log and return the DB result
- Redis invalidation failure: log, rely on TTL, do not fail the publish request

This matches the current broadcast philosophy, where a broadcast failure does not invalidate the successful persistence write.

## Why Not Cache Auth

`GET /api/v1/auth/me` is not the right first target:

- lower fan-out than feed reads
- tighter correctness expectations
- relatively cheap database path
- little payoff compared to the complexity of user-data invalidation

Auth should continue to rely on DB-backed reads and SPA-side local storage for session convenience.

## Why Not Cache WebSocket Deltas

Redis pub/sub already serves the live delivery requirement.

Caching the deltas would only be useful if the product needed replay, gap recovery, or short-term event history for reconnecting clients. The current design explicitly separates:

- feed history via `GET /feed`
- live deltas via `/feed/ws`

That separation should remain intact unless product requirements change.

## Rollout Plan

### Phase 1

- Add `FeedCache` interface
- Add no-op implementation
- Add Redis-backed cache implementation
- Cache only unfiltered first-page reads
- Invalidate unfiltered first-page key on publish

### Phase 2

- Add first-page cache entries per `event_type`
- Invalidate matching filtered first-page key on publish

### Phase 3

- Add metrics for hit rate, miss rate, DB fallback rate, and invalidation failures
- Re-evaluate whether `offset > 0` should be cached

### Deferred

- search query caching
- deep page caching
- replay buffers for WebSocket reconnects
- full feed materialization in Redis

These should remain deferred until traffic data proves they are necessary.

## Testing Strategy

Add tests for:

- cache hit on cacheable unfiltered first-page request
- cache miss with DB fallback and cache fill
- non-cacheable request bypasses cache
- publish invalidates expected keys
- Redis failures degrade gracefully without breaking reads or writes

Keep the use case tests focused on orchestration and infrastructure tests focused on key format, TTL behavior, and serialization.

## Success Criteria

- reduced PostgreSQL read load for hot feed views
- unchanged API contract for existing clients
- no correctness dependency on Redis availability
- bounded cache key cardinality
- straightforward extension path for future feed-view patterns

## Decision

Adopt a narrow Redis-backed cache-aside strategy for hot feed reads only, with short TTLs and publish-triggered invalidation, while preserving PostgreSQL as the source of truth and Redis pub/sub as the real-time transport mechanism.
