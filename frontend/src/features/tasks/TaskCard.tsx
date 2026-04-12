import { Card, CardContent, Typography, Chip, Box, IconButton } from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import type { Task } from '../../api/types';

interface TaskCardProps {
  task: Task;
  onEdit: (task: Task) => void;
  onDelete: (taskId: string) => void;
}

const priorityColors = {
  low: 'success',
  medium: 'warning',
  high: 'error',
} as const;

const statusLabels = {
  todo: 'To Do',
  in_progress: 'In Progress',
  done: 'Done',
};

const TaskCard = ({ task, onEdit, onDelete }: TaskCardProps) => {
  return (
    <Card 
      sx={{ 
        mb: 1.5, 
        cursor: 'pointer',
        '&:hover': { boxShadow: 3 }
      }}
      onClick={() => onEdit(task)}
    >
      <CardContent sx={{ pb: '12px !important', '&:last-child': { pb: '12px' } }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 1 }}>
          <Typography variant="subtitle1" sx={{ fontWeight: 500 }}>
            {task.title}
          </Typography>
          <IconButton 
            size="small" 
            onClick={(e) => {
              e.stopPropagation();
              onDelete(task.id);
            }}
          >
            <DeleteIcon fontSize="small" />
          </IconButton>
        </Box>
        
        {task.description && (
          <Typography 
            variant="body2" 
            color="text.secondary" 
            sx={{ 
              display: '-webkit-box',
              WebkitLineClamp: 2,
              WebkitBoxOrient: 'vertical',
              overflow: 'hidden',
              mb: 1
            }}
          >
            {task.description}
          </Typography>
        )}
        
        <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
          <Chip 
            label={statusLabels[task.status]} 
            size="small" 
            variant="outlined"
          />
          <Chip 
            label={task.priority} 
            size="small" 
            color={priorityColors[task.priority]}
          />
        </Box>
        
        {task.due_date && (
          <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
            Due: {new Date(task.due_date).toLocaleDateString()}
          </Typography>
        )}
      </CardContent>
    </Card>
  );
};

export default TaskCard;