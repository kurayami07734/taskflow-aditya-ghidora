import { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  MenuItem,
  Select,
  FormControl,
  InputLabel,
  Box,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import dayjs, { Dayjs } from 'dayjs';
import { taskApi, userApi } from '../../api';
import type { Task, User } from '../../api/types';
import { useSnackbar } from '../../components/common/SnackbarProvider';
import UserSearch from '../../components/common/UserSearch';

interface TaskFormData {
  title: string;
  description: string;
  status: Task['status'];
  priority: Task['priority'];
  due_date: string | null;
  assignee_id: string | null;
}

interface TaskFormModalProps {
  open: boolean;
  onClose: () => void;
  projectId: string;
  task?: Task | null;
}

const TaskFormModal = ({ open, onClose, projectId, task }: TaskFormModalProps) => {
  const queryClient = useQueryClient();
  const { showError, showSuccess } = useSnackbar();
  const [assignee, setAssignee] = useState<User | null>(null);
  const [dueDate, setDueDate] = useState<Dayjs | null>(null);
  const [formData, setFormData] = useState<TaskFormData>({
    title: '',
    description: '',
    status: 'todo',
    priority: 'medium',
    due_date: null,
    assignee_id: null,
  });

  useEffect(() => {
    if (task) {
      setFormData({
        title: task.title,
        description: task.description || '',
        status: task.status,
        priority: task.priority,
        due_date: task.due_date,
        assignee_id: task.assignee_id,
      });
      setDueDate(task.due_date ? dayjs(task.due_date) : null);
      
      if (task.assignee_id) {
        userApi.get(task.assignee_id).then(({ data }) => setAssignee(data)).catch(() => setAssignee(null));
      } else {
        setAssignee(null);
      }
    } else {
      setFormData({
        title: '',
        description: '',
        status: 'todo',
        priority: 'medium',
        due_date: null,
        assignee_id: null,
      });
      setDueDate(null);
      setAssignee(null);
    }
  }, [task, open]);

  const isEditing = !!task;

  const mutation = useMutation({
    mutationFn: (data: TaskFormData) => 
      isEditing 
        ? taskApi.update(task.id, data)
        : taskApi.create(projectId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks', projectId] });
      onClose();
      showSuccess(isEditing ? 'Task updated successfully' : 'Task created successfully');
    },
    onError: (err: unknown) => {
      const axiosError = err as { response?: { data?: { message?: string } } };
      showError(axiosError.response?.data?.message || 'Failed to save task');
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const submitData = {
      ...formData,
      due_date: dueDate ? dueDate.toISOString() : null,
    };
    mutation.mutate(submitData);
  };

  const handleAssigneeSelect = (user: User | null) => {
    setAssignee(user);
    setFormData({ ...formData, assignee_id: user?.id ?? null });
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>{task ? 'Edit Task' : 'Create New Task'}</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            fullWidth
            label="Title"
            margin="normal"
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            required
          />
          <TextField
            fullWidth
            label="Description"
            margin="normal"
            multiline
            rows={3}
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          />
          <Box sx={{ display: 'flex', gap: 2, mt: 1 }}>
            <FormControl fullWidth margin="normal">
              <InputLabel>Status</InputLabel>
              <Select
                value={formData.status}
                label="Status"
                onChange={(e) => setFormData({ ...formData, status: e.target.value as Task['status'] })}
              >
                <MenuItem value="todo">Todo</MenuItem>
                <MenuItem value="in_progress">In Progress</MenuItem>
                <MenuItem value="done">Done</MenuItem>
              </Select>
            </FormControl>
            <FormControl fullWidth margin="normal">
              <InputLabel>Priority</InputLabel>
              <Select
                value={formData.priority}
                label="Priority"
                onChange={(e) => setFormData({ ...formData, priority: e.target.value as Task['priority'] })}
              >
                <MenuItem value="low">Low</MenuItem>
                <MenuItem value="medium">Medium</MenuItem>
                <MenuItem value="high">High</MenuItem>
              </Select>
            </FormControl>
          </Box>
          <Box sx={{ mt: 1 }}>
            <UserSearch
              selectedUser={assignee}
              onSelect={handleAssigneeSelect}
            />
          </Box>
          <Box sx={{ mt: 1 }}>
            <DatePicker
              label="Due Date"
              value={dueDate}
              onChange={(newValue) => setDueDate(newValue)}
              format="DD/MM/YYYY"
              slotProps={{
                textField: {
                  fullWidth: true,
                  margin: 'normal',
                },
              }}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Cancel</Button>
          <Button 
            type="submit" 
            variant="contained"
            disabled={mutation.isPending}
          >
            {mutation.isPending ? (task ? 'Saving...' : 'Creating...') : (task ? 'Save' : 'Create')}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default TaskFormModal;