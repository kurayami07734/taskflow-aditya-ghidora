import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ThemeProvider } from '@mui/material/styles'
import CssBaseline from '@mui/material/CssBaseline'
import { LocalizationProvider } from '@mui/x-date-pickers'
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs'
import App from './App'
import theme from './theme/theme'
import { useAuthStore } from './stores/authStore'
import { SnackbarProvider } from './components/common/SnackbarProvider'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

const InitializeAuth = ({ children }: { children: React.ReactNode }) => {
  useAuthStore.getState().initialize()
  return <>{children}</>
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <LocalizationProvider dateAdapter={AdapterDayjs}>
          <SnackbarProvider>
            <InitializeAuth>
              <BrowserRouter>
                <App />
              </BrowserRouter>
            </InitializeAuth>
          </SnackbarProvider>
        </LocalizationProvider>
      </ThemeProvider>
    </QueryClientProvider>
  </StrictMode>,
)