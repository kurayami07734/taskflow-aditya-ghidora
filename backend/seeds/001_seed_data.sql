-- Seed data for testing
-- Email: test@example.com
-- Password: password123 (bcrypt hash with cost 12)
-- Hash: $2a$12$HkHVdkjAFjJWaHEgevQl4OC6d3BpyXJZQbygFnzpWYFJZL3YjjTnK

INSERT INTO users (id, name, email, password, created_at) 
VALUES (
    '028de2bc-6c51-4eb4-994a-0b305d507b65',
    'Test User',
    'test@example.com',
    '$2a$12$HkHVdkjAFjJWaHEgevQl4OC6d3BpyXJZQbygFnzpWYFJZL3YjjTnK',
    NOW()
) ON CONFLICT (email) DO NOTHING;

INSERT INTO projects (id, name, description, owner_id, created_at)
VALUES (
    '4ed6f48f-eb63-4412-9941-5f2880c34488',
    'Test Project',
    'This is a sample project for testing',
    '028de2bc-6c51-4eb4-994a-0b305d507b65',
    NOW()
) ON CONFLICT DO NOTHING;

INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at)
VALUES 
(
    '307f75c6-4a60-4dc1-a8f2-54830d78981a',
    'Setup project structure',
    'Initialize the project with proper folder structure',
    'done',
    'high',
    '4ed6f48f-eb63-4412-9941-5f2880c34488',
    '028de2bc-6c51-4eb4-994a-0b305d507b65',
    '2026-04-01',
    NOW(),
    NOW()
),
(
    '3dea209d-bdd4-49b8-8b3c-3ce59e38e1a0',
    'Implement authentication',
    'Add login and registration functionality',
    'in_progress',
    'high',
    '4ed6f48f-eb63-4412-9941-5f2880c34488',
    '028de2bc-6c51-4eb4-994a-0b305d507b65',
    '2026-04-15',
    NOW(),
    NOW()
),
(
    '6a401351-1475-48e7-b7a7-23252aadd195',
    'Write documentation',
    'Create README and API documentation',
    'todo',
    'medium',
    '4ed6f48f-eb63-4412-9941-5f2880c34488',
    '028de2bc-6c51-4eb4-994a-0b305d507b65',
    '2026-04-20',
    NOW(),
    NOW()
) ON CONFLICT DO NOTHING;