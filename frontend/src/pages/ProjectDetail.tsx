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
import { projectApi, taskApi } from '../api';
import TaskBoard from '../features/tasks/TaskBoard';
import TaskFormModal from '../features/tasks/TaskFormModal';
import type { Task } from '../api/types';

const ProjectDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);

  const { data: project, isLoading: projectLoading, isError: projectError } = useQuery({
    queryKey: ['project', id],
    queryFn: async () => {
      const { data } = await projectApi.get(id!);
      return data;
    },
    enabled: !!id,
  });

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
          onDeleteTask={(taskId) => {
            if (confirm('Are you sure you want to delete this task?')) {
              deleteTaskMutation.mutate(taskId);
            }
          }}
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
    </Box>
  );
};

export default ProjectDetail;