import { useMemo } from 'react';
import { Box, Typography, Paper, useTheme } from '@mui/material';
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
  const theme = useTheme();

  const tasksByStatus = useMemo(() => ({
    todo: tasks.filter((task) => task.status === 'todo'),
    in_progress: tasks.filter((task) => task.status === 'in_progress'),
    done: tasks.filter((task) => task.status === 'done'),
  }), [tasks]);

  return (
    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 3, pb: 2 }}>
      {columns.map((column) => {
        const columnTasks = tasksByStatus[column.id];
        return (
          <Paper 
            key={column.id}
            sx={{ 
              flex: '1 1 300px', 
              maxWidth: { xs: '100%', md: 'calc(33.333% - 16px)' },
              p: 2,
              bgcolor: theme.palette.mode === 'dark' ? 'grey.900' : 'grey.50'
            }}
          >
            <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 2 }}>
              {column.title} ({columnTasks.length})
            </Typography>
            <Box>
              {columnTasks.map((task) => (
                <TaskCard
                  key={task.id}
                  task={task}
                  onEdit={onEditTask}
                  onDelete={onDeleteTask}
                />
              ))}
              {columnTasks.length === 0 && (
                <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', py: 2 }}>
                  No tasks
                </Typography>
              )}
            </Box>
          </Paper>
        );
      })}
    </Box>
  );
};

export default TaskBoard;