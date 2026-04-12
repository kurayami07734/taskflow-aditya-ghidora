import { Routes, Route, Navigate } from 'react-router-dom'
import AuthGuard from './components/auth/AuthGuard'
import AppLayout from './components/layout/AppLayout'
import Login from './pages/Login'
import Register from './pages/Register'
import ProjectList from './pages/ProjectList'
import ProjectDetail from './pages/ProjectDetail'

function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      
      <Route element={<AuthGuard />}>
        <Route element={<AppLayout />}>
          <Route path="/projects" element={<ProjectList />} />
          <Route path="/projects/:id" element={<ProjectDetail />} />
        </Route>
      </Route>
      
      <Route path="/" element={<Navigate to="/projects" replace />} />
    </Routes>
  )
}

export default App