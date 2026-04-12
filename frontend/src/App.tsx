import { Routes, Route, Navigate, useNavigate } from 'react-router-dom'
import { useEffect } from 'react'
import AuthGuard from './components/auth/AuthGuard'
import AppLayout from './components/layout/AppLayout'
import Login from './pages/Login'
import Register from './pages/Register'
import ProjectList from './pages/ProjectList'
import ProjectDetail from './pages/ProjectDetail'

function App() {
  const navigate = useNavigate()

  useEffect(() => {
    const handleUnauthorized = () => {
      navigate('/login', { replace: true })
    }
    window.addEventListener('auth:unauthorized', handleUnauthorized)
    return () => window.removeEventListener('auth:unauthorized', handleUnauthorized)
  }, [navigate])

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