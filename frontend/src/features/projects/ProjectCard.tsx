import { Link as RouterLink } from 'react-router-dom';
import { Card, CardContent, Typography, CardActionArea, Box } from '@mui/material';
import type { Project } from '../../api/types';

interface ProjectCardProps {
  project: Project;
}

const ProjectCard = ({ project }: ProjectCardProps) => {
  return (
    <Card sx={{ height: '100%' }}>
      <CardActionArea component={RouterLink} to={`/projects/${project.id}`}>
        <CardContent>
          <Typography variant="h6" component="div" gutterBottom>
            {project.name}
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{
            display: '-webkit-box',
            WebkitLineClamp: 2,
            WebkitBoxOrient: 'vertical',
            overflow: 'hidden',
            mb: 1
          }}>
            {project.description || 'No description'}
          </Typography>
          <Box sx={{ mt: 1 }}>
            <Typography variant="caption" color="text.secondary">
              Created: {new Date(project.created_at).toLocaleDateString()}
            </Typography>
          </Box>
        </CardContent>
      </CardActionArea>
    </Card>
  );
};

export default ProjectCard;