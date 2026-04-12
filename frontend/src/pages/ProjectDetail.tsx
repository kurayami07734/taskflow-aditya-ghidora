import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Typography,
  Box,
  Button,
  CircularProgress,
  Alert,
  Breadcrumbs,
  Link,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import DeleteIcon from '@mui/icons-material/Delete';
import { projectApi, taskApi } from '../api';
import TaskBoard from '../features/tasks/TaskBoard';
import TaskFormModal from '../features/tasks/TaskFormModal';
import ConfirmDialog from '../components/common/ConfirmDialog';
import type { Task } from '../api/types';
import { useSnackbar } from '../components/common/SnackbarProvider';

const ProjectDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { showError, showSuccess } = useSnackbar();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [deleteTask, setDeleteTask] = useState<Task | null>(null);

  const { data: project, isLoading: projectLoading, isError: projectError } = useQuery({
    queryKey: ['project', id],
    queryFn: async () => {
      const { data } = await projectApi.get(id!);
      return data;
    },
    enabled: !!id,
  });

  if (projectError) {
    showError('Failed to load project');
  }

  const { data: tasksResponse, isLoading: tasksLoading, isError: tasksError } = useQuery({
    queryKey: ['tasks', id],
    queryFn: async () => {
      const { data } = await taskApi.listByProject(id!);
      return data;
    },
    enabled: !!id,
  });

  const tasks = tasksResponse?.tasks || [];

  const deleteTaskMutation = useMutation({
    mutationFn: taskApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks', id] });
      showSuccess('Task deleted successfully');
      setDeleteTask(null);
    },
    onError: (err: unknown) => {
      const axiosError = err as { response?: { data?: { message?: string } } };
      showError(axiosError.response?.data?.message || 'Failed to delete task');
    },
  });

  if (projectLoading || tasksLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (projectError) {
    return (
      <Alert severity="error" sx={{ mt: 2 }}>
        Failed to load project
      </Alert>
    );
  }

  if (!project) {
    return (
      <Alert severity="warning" sx={{ mt: 2 }}>
        Project not found
      </Alert>
    );
  }

  return (
    <Box>
      <Breadcrumbs sx={{ mb: 2 }}>
        <Link 
          underline="hover" 
          color="inherit" 
          onClick={() => navigate('/projects')}
          sx={{ cursor: 'pointer' }}
        >
          Projects
        </Link>
        <Typography color="text.primary">{project.name}</Typography>
      </Breadcrumbs>

      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Box>
          <Typography variant="h4">{project.name}</Typography>
          {project.description && (
            <Typography variant="body1" color="text.secondary" sx={{ mt: 1 }}>
              {project.description}
            </Typography>
          )}
        </Box>
        <Box sx={{ display: 'flex', gap: 1 }}>
          <Button
            variant="outlined"
            color="error"
            startIcon={<DeleteIcon />}
            onClick={() => setDeleteTask({ id: '', title: '', description: null, status: 'todo', priority: 'medium', project_id: id!, assignee: null, due_date: null, created_at: '', updated_at: '' } as Task)}
          >
            Delete Project
          </Button>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => {
              setEditingTask(null);
              setModalOpen(true);
            }}
          >
            Add Task
          </Button>
        </Box>
      </Box>

      {tasksError ? (
        <Alert severity="error" sx={{ mt: 2 }}>
          Failed to load tasks
        </Alert>
      ) : tasks && tasks.length > 0 ? (
        <TaskBoard
          tasks={tasks}
          onEditTask={(task) => {
            setEditingTask(task);
            setModalOpen(true);
          }}
          onDeleteTask={(task) => setDeleteTask(task)}
        />
      ) : (
        <Box sx={{ textAlign: 'center', mt: 4 }}>
          <Typography variant="h6" color="text.secondary">
            No tasks yet
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            Create your first task to get started!
          </Typography>
        </Box>
      )}

      <TaskFormModal
        open={modalOpen}
        onClose={() => {
          setModalOpen(false);
          setEditingTask(null);
        }}
        projectId={id!}
        task={editingTask}
      />

      <ConfirmDialog
        open={!!deleteTask}
        title={deleteTask?.title ? 'Delete Task' : 'Delete Project'}
        message={
          deleteTask?.title
            ? `Are you sure you want to delete the task "${deleteTask.title}"?`
            : `Are you sure you want to delete the project "${project.name}"? This will also delete all tasks in this project.`
        }
        onConfirm={() => {
          if (deleteTask?.title) {
            deleteTaskMutation.mutate(deleteTask.id);
          } else {
            projectApi.delete(id!).then(() => {
              showSuccess('Project deleted successfully');
              navigate('/projects');
            }).catch((err) => {
              const axiosError = err as { response?: { data?: { message?: string } } };
              showError(axiosError.response?.data?.message || 'Failed to delete project');
            });
          }
        }}
        onCancel={() => setDeleteTask(null)}
      />
    </Box>
  );
};

export default ProjectDetail;