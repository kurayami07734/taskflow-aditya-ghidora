import { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
} from '@mui/material';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { projectApi } from '../../api';
import { useSnackbar } from '../../components/common/SnackbarProvider';

interface CreateProjectModalProps {
  open: boolean;
  onClose: () => void;
}

const CreateProjectModal = ({ open, onClose }: CreateProjectModalProps) => {
  const queryClient = useQueryClient();
  const { showError, showSuccess } = useSnackbar();
  const [formData, setFormData] = useState({ name: '', description: '' });

  const mutation = useMutation({
    mutationFn: projectApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      setFormData({ name: '', description: '' });
      onClose();
      showSuccess('Project created successfully');
    },
    onError: (err: unknown) => {
      const axiosError = err as { response?: { data?: { message?: string } } };
      showError(axiosError.response?.data?.message || 'Failed to create project');
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    mutation.mutate(formData);
  };

  const handleClose = () => {
    setFormData({ name: '', description: '' });
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>Create New Project</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            fullWidth
            label="Project Name"
            margin="normal"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
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
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button 
            type="submit" 
            variant="contained"
            disabled={mutation.isPending}
          >
            {mutation.isPending ? 'Creating...' : 'Create'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default CreateProjectModal;