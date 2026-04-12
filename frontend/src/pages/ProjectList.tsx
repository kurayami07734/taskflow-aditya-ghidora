import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  Typography,
  Box,
  Button,
  CircularProgress,
  Alert,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import { projectApi } from '../api';
import ProjectCard from '../features/projects/ProjectCard';
import CreateProjectModal from '../features/projects/CreateProjectModal';

const ProjectList = () => {
  const [modalOpen, setModalOpen] = useState(false);

  const { data: projects, isLoading, isError, error } = useQuery({
    queryKey: ['projects'],
    queryFn: async () => {
      const { data } = await projectApi.list();
      return data;
    },
  });

  if (isLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (isError) {
    return (
      <Alert severity="error" sx={{ mt: 2 }}>
        Failed to load projects: {(error as Error).message}
      </Alert>
    );
  }

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4">Projects</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setModalOpen(true)}
        >
          New Project
        </Button>
      </Box>

      {!projects || projects.length === 0 ? (
        <Box sx={{ textAlign: 'center', mt: 4 }}>
          <Typography variant="h6" color="text.secondary">
            No projects yet
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            Create your first project to get started!
          </Typography>
        </Box>
      ) : (
        <Box sx={{ 
          display: 'grid', 
          gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' }, 
          gap: 3 
        }}>
          {projects.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))}
        </Box>
      )}

      <CreateProjectModal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
      />
    </Box>
  );
};

export default ProjectList;