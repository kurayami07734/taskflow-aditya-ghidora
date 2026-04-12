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

interface CreateProjectModalProps {
  open: boolean;
  onClose: () => void;
}

const CreateProjectModal = ({ open, onClose }: CreateProjectModalProps) => {
  const queryClient = useQueryClient();
  const [formData, setFormData] = useState({ name: '', description: '' });

  const mutation = useMutation({
    mutationFn: projectApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      setFormData({ name: '', description: '' });
      onClose();
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    mutation.mutate(formData);
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
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
          <Button onClick={onClose}>Cancel</Button>
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