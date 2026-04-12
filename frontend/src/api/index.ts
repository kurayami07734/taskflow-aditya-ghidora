import api from './axios';
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  Project,
  CreateProjectRequest,
  UpdateProjectRequest,
  Task,
  CreateTaskRequest,
  UpdateTaskRequest,
} from './types';

export const authApi = {
  login: (data: LoginRequest) =>
    api.post<AuthResponse>('/auth/login', data),
  
  register: (data: RegisterRequest) =>
    api.post<AuthResponse>('/auth/register', data),
};

export const projectApi = {
  list: () => api.get<Project[]>('/projects'),
  
  get: (id: string) => api.get<Project>(`/projects/${id}`),
  
  create: (data: CreateProjectRequest) =>
    api.post<Project>('/projects', data),
  
  update: (id: string, data: UpdateProjectRequest) =>
    api.patch<Project>(`/projects/${id}`, data),
  
  delete: (id: string) => api.delete(`/projects/${id}`),
};

export const taskApi = {
  listByProject: (projectId: string) =>
    api.get<Task[]>(`/projects/${projectId}/tasks`),
  
  create: (projectId: string, data: CreateTaskRequest) =>
    api.post<Task>(`/projects/${projectId}/tasks`, data),
  
  update: (taskId: string, data: UpdateTaskRequest) =>
    api.patch<Task>(`/tasks/${taskId}`, data),
  
  delete: (taskId: string) => api.delete(`/tasks/${taskId}`),
};