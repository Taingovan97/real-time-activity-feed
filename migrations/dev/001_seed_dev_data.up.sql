-- Dev seed data migration
-- This migration seeds test users and activity feed events for development only
-- All inserts are idempotent using ON CONFLICT clauses

-- Seed test users
-- Password for all users: password123
-- Bcrypt hash: $2a$10$9DSodToLn2m.h3i1uYQocu//OKlkvjHk3rdtB463cKQ0bSKDupJwC
INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
VALUES 
    ('00000000-0000-0000-0000-000000000001', 'alice', 'alice@example.com', '$2a$10$9DSodToLn2m.h3i1uYQocu//OKlkvjHk3rdtB463cKQ0bSKDupJwC', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000002', 'bob', 'bob@example.com', '$2a$10$9DSodToLn2m.h3i1uYQocu//OKlkvjHk3rdtB463cKQ0bSKDupJwC', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000003', 'charlie', 'charlie@example.com', '$2a$10$9DSodToLn2m.h3i1uYQocu//OKlkvjHk3rdtB463cKQ0bSKDupJwC', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000004', 'dave', 'dave@example.com', '$2a$10$9DSodToLn2m.h3i1uYQocu//OKlkvjHk3rdtB463cKQ0bSKDupJwC', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000005', 'eve', 'eve@example.com', '$2a$10$9DSodToLn2m.h3i1uYQocu//OKlkvjHk3rdtB463cKQ0bSKDupJwC', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000006', 'frank', 'frank@example.com', '$2a$10$9DSodToLn2m.h3i1uYQocu//OKlkvjHk3rdtB463cKQ0bSKDupJwC', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000007', 'grace', 'grace@example.com', '$2a$10$9DSodToLn2m.h3i1uYQocu//OKlkvjHk3rdtB463cKQ0bSKDupJwC', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000008', 'henry', 'henry@example.com', '$2a$10$9DSodToLn2m.h3i1uYQocu//OKlkvjHk3rdtB463cKQ0bSKDupJwC', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000009', 'ivy', 'ivy@example.com', '$2a$10$9DSodToLn2m.h3i1uYQocu//OKlkvjHk3rdtB463cKQ0bSKDupJwC', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000010', 'jack', 'jack@example.com', '$2a$10$9DSodToLn2m.h3i1uYQocu//OKlkvjHk3rdtB463cKQ0bSKDupJwC', NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- Seed activity feed events
INSERT INTO activity_feed_events (id, user_id, event_type, content, created_at)
VALUES 
    ('00000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000001', 'login', 'Alice signed in from the web dashboard', NOW() - INTERVAL '10 minutes'),
    ('00000000-0000-0000-0000-000000000102', '00000000-0000-0000-0000-000000000002', 'upload', 'Bob uploaded the quarterly metrics file', NOW() - INTERVAL '9 minutes'),
    ('00000000-0000-0000-0000-000000000103', '00000000-0000-0000-0000-000000000003', 'comment', 'Charlie commented on the release checklist', NOW() - INTERVAL '8 minutes'),
    ('00000000-0000-0000-0000-000000000104', '00000000-0000-0000-0000-000000000004', 'notification', 'Dave acknowledged the deployment alert', NOW() - INTERVAL '7 minutes'),
    ('00000000-0000-0000-0000-000000000105', '00000000-0000-0000-0000-000000000005', 'report', 'Eve generated the daily operations report', NOW() - INTERVAL '6 minutes'),
    ('00000000-0000-0000-0000-000000000106', '00000000-0000-0000-0000-000000000006', 'task', 'Frank completed the Redis cache refresh task', NOW() - INTERVAL '5 minutes'),
    ('00000000-0000-0000-0000-000000000107', '00000000-0000-0000-0000-000000000007', 'message', 'Grace posted a handoff note for support', NOW() - INTERVAL '4 minutes'),
    ('00000000-0000-0000-0000-000000000108', '00000000-0000-0000-0000-000000000008', 'sync', 'Henry synced user permissions', NOW() - INTERVAL '3 minutes'),
    ('00000000-0000-0000-0000-000000000109', '00000000-0000-0000-0000-000000000009', 'status', 'Ivy updated the incident channel status', NOW() - INTERVAL '2 minutes'),
    ('00000000-0000-0000-0000-000000000110', '00000000-0000-0000-0000-000000000010', 'approval', 'Jack approved the release checklist', NOW() - INTERVAL '1 minute')
ON CONFLICT (id) DO NOTHING;
