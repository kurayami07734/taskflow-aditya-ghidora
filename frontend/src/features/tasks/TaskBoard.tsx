import { Box, Typography, Paper } from '@mui/material';
import type { Task } from '../../api/types';
import TaskCard from './TaskCard';

interface TaskBoardProps {
  tasks: Task[];
  onEditTask: (task: Task) => void;
  onDeleteTask: (task: Task) => void;
}

interface Column {
  id: Task['status'];
  title: string;
}

const columns: Column[] = [
  { id: 'todo', title: 'Todo' },
  { id: 'in_progress', title: 'In Progress' },
  { id: 'done', title: 'Done' },
];

const TaskBoard = ({ tasks, onEditTask, onDeleteTask }: TaskBoardProps) => {
  const getTasksByStatus = (status: Task['status']) => 
    tasks.filter((task) => task.status === status);

  return (
    <Box sx={{ display: 'flex', gap: 3, overflowX: 'auto', pb: 2 }}>
      {columns.map((column) => (
        <Paper 
          key={column.id}
          sx={{ 
            flex: '1 1 300px', 
            minWidth: 280,
            p: 2,
            bgcolor: 'grey.50'
          }}
        >
          <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 2 }}>
            {column.title} ({getTasksByStatus(column.id).length})
          </Typography>
          <Box>
            {getTasksByStatus(column.id).map((task) => (
              <TaskCard
                key={task.id}
                task={task}
                onEdit={onEditTask}
                onDelete={onDeleteTask}
              />
            ))}
            {getTasksByStatus(column.id).length === 0 && (
              <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', py: 2 }}>
                No tasks
              </Typography>
            )}
          </Box>
        </Paper>
      ))}
    </Box>
  );
};

export default TaskBoard;