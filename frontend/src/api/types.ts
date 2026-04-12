export interface User {
  id: string;
  name: string;
  email: string;
  created_at: string;
}

export interface Project {
  id: string;
  name: string;
  description: string | null;
  created_at: string;
}

export interface Pagination {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface ProjectListResponse {
  projects: Project[];
  pagination: Pagination;
}

export interface Task {
  id: string;
  project_id: string;
  title: string;
  description: string | null;
  status: 'todo' | 'in_progress' | 'done';
  priority: 'low' | 'medium' | 'high';
  due_date: string | null;
  assignee: string | null;
  created_at: string;
  updated_at: string;
}

export interface TaskListResponse {
  tasks: Task[];
  pagination: Pagination;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface CreateProjectRequest {
  name: string;
  description: string;
}

export interface UpdateProjectRequest {
  name?: string;
  description?: string;
}

export interface CreateTaskRequest {
  title: string;
  description: string;
  status?: 'todo' | 'in_progress' | 'done';
  priority?: 'low' | 'medium' | 'high';
  due_date?: string | null;
  assignee?: string | null;
}

export interface UpdateTaskRequest {
  title?: string;
  description?: string;
  status?: 'todo' | 'in_progress' | 'done';
  priority?: 'low' | 'medium' | 'high';
  due_date?: string | null;
  assignee?: string | null;
}