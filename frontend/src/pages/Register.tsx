import { useState, useMemo } from 'react';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import {
  Box,
  Paper,
  TextField,
  Button,
  Typography,
  Link,
  InputAdornment,
  IconButton,
  LinearProgress,
  AppBar,
  Toolbar,
} from '@mui/material';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import LightMode from '@mui/icons-material/LightMode';
import DarkMode from '@mui/icons-material/DarkMode';
import { authApi } from '../api';
import { useAuthStore } from '../stores/authStore';
import { useThemeStore } from '../stores/themeStore';
import { useSnackbar } from '../components/common/SnackbarProvider';

const validateEmail = (email: string): boolean => {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
};

const getPasswordStrength = (password: string): { score: number; label: string; color: string } => {
  if (!password) return { score: 0, label: '', color: 'grey' };
  
  let score = 0;
  if (password.length >= 8) score += 1;
  if (password.length >= 12) score += 1;
  if (/[a-z]/.test(password) && /[A-Z]/.test(password)) score += 1;
  if (/\d/.test(password)) score += 1;
  if (/[!@#$%^&*(),.?":{}|<>]/.test(password)) score += 1;

  if (score <= 1) return { score: score / 5, label: 'Weak', color: 'error' };
  if (score <= 3) return { score: score / 5, label: 'Medium', color: 'warning' };
  return { score: score / 5, label: 'Strong', color: 'success' };
};

const Register = () => {
  const navigate = useNavigate();
  const setAuth = useAuthStore((state) => state.setAuth);
  const { showError } = useSnackbar();
  const { mode, toggleMode } = useThemeStore();
  const [formData, setFormData] = useState({ name: '', email: '', password: '' });
  const [errors, setErrors] = useState({ name: '', email: '', password: '' });
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const passwordStrength = useMemo(() => getPasswordStrength(formData.password), [formData.password]);

  const validateForm = (): boolean => {
    const newErrors = { name: '', email: '', password: '' };
    let isValid = true;

    if (!formData.name) {
      newErrors.name = 'Name is required';
      isValid = false;
    }

    if (!formData.email) {
      newErrors.email = 'Email is required';
      isValid = false;
    } else if (!validateEmail(formData.email)) {
      newErrors.email = 'Invalid email format';
      isValid = false;
    }

    if (!formData.password) {
      newErrors.password = 'Password is required';
      isValid = false;
    } else if (formData.password.length < 8) {
      newErrors.password = 'Password must be at least 8 characters';
      isValid = false;
    }

    setErrors(newErrors);
    return isValid;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setLoading(true);

    try {
      const { data } = await authApi.register(formData);
      setAuth(data.user, data.token);
      navigate('/projects');
    } catch (err: unknown) {
      const axiosError = err as { response?: { data?: { message?: string } } };
      showError(axiosError.response?.data?.message || 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ minHeight: '100vh', bgcolor: 'background.default', position: 'relative' }}>
      <AppBar position="absolute" color="transparent" elevation={0}>
        <Toolbar>
          <Typography variant="h5" component="div" sx={{ flexGrow: 1, fontWeight: 600 }}>
            TaskFlow
          </Typography>
          <IconButton
            onClick={toggleMode}
            title={`Switch to ${mode === 'light' ? 'dark' : 'light'} mode`}
            sx={{
              bgcolor: 'background.paper',
              borderRadius: 2,
              position: 'absolute',
              right: 16,
              '&:hover': { bgcolor: 'action.hover' },
            }}
          >
            {mode === 'light' ? <DarkMode /> : <LightMode />}
          </IconButton>
        </Toolbar>
      </AppBar>

      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          minHeight: '100vh',
        }}
      >
        <Paper sx={{ p: 4, width: 400 }}>
          <Typography variant="h4" component="h1" gutterBottom sx={{ textAlign: 'center', mb: 3 }}>
            Create Account
          </Typography>
          
          <form onSubmit={handleSubmit}>
            <TextField
              fullWidth
              label="Name"
              margin="normal"
              value={formData.name}
              onChange={(e) => {
                setFormData({ ...formData, name: e.target.value });
                if (errors.name) {
                  setErrors({ ...errors, name: '' });
                }
              }}
              error={!!errors.name}
              helperText={errors.name}
              required
            />
            <TextField
              fullWidth
              label="Email"
              type="email"
              margin="normal"
              value={formData.email}
              onChange={(e) => {
                setFormData({ ...formData, email: e.target.value });
                if (errors.email) {
                  setErrors({ ...errors, email: '' });
                }
              }}
              error={!!errors.email}
              helperText={errors.email}
              required
            />
            <TextField
              fullWidth
              label="Password"
              type={showPassword ? 'text' : 'password'}
              margin="normal"
              value={formData.password}
              onChange={(e) => {
                setFormData({ ...formData, password: e.target.value });
                if (errors.password) {
                  setErrors({ ...errors, password: '' });
                }
              }}
              error={!!errors.password}
              helperText={errors.password || "Minimum 8 characters"}
              required
              slotProps={{
                input: {
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        onClick={() => setShowPassword(!showPassword)}
                        edge="end"
                      >
                        {showPassword ? <VisibilityOff /> : <Visibility />}
                      </IconButton>
                    </InputAdornment>
                  ),
                },
              }}
            />
            {formData.password && (
              <Box sx={{ mt: 1, mb: 2 }}>
                <LinearProgress
                  variant="determinate"
                  value={passwordStrength.score * 100}
                  color={passwordStrength.color as 'error' | 'warning' | 'success'}
                  sx={{ height: 6, borderRadius: 3 }}
                />
                <Typography variant="caption" sx={{ mt: 0.5, display: 'block' }}>
                  Password strength: 
                  <Typography
                    component="span"
                    variant="caption"
                    sx={{ color: `${passwordStrength.color}.main`, fontWeight: 500, ml: 0.5 }}
                  >
                    {passwordStrength.label}
                  </Typography>
                </Typography>
              </Box>
            )}
            <Button
              fullWidth
              variant="contained"
              type="submit"
              disabled={loading}
              sx={{ mt: 1 }}
            >
              {loading ? 'Registering...' : 'Register'}
            </Button>
          </form>
          
          <Typography variant="body2" sx={{ mt: 2, textAlign: 'center' }}>
            Already have an account?{' '}
            <Link component={RouterLink} to="/login">
              Login
            </Link>
          </Typography>
        </Paper>
      </Box>
    </Box>
  );
};

export default Register;