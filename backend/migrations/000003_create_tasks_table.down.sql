DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks;
DROP FUNCTION IF EXISTS update_updated_at();
DROP TABLE IF EXISTS tasks;
DROP TYPE IF EXISTS task_status;
DROP TYPE IF EXISTS task_priority;