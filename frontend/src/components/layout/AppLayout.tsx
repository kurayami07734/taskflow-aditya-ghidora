import { Box, Container } from '@mui/material';
import { Outlet } from 'react-router-dom';
import Navbar from './Navbar';

const AppLayout = () => {
  return (
    <Box sx={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Navbar />
      <Container maxWidth="lg" sx={{ flex: 1, py: 3 }}>
        <Outlet />
      </Container>
    </Box>
  );
};

export default AppLayout;