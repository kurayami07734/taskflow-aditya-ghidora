import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ThemeProvider } from '@mui/material/styles'
import CssBaseline from '@mui/material/CssBaseline'
import { LocalizationProvider } from '@mui/x-date-pickers'
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs'
import App from './App'
import { getTheme } from './theme/theme'
import { useThemeStore } from './stores/themeStore'
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

const ThemedApp = ({ children }: { children: React.ReactNode }) => {
  const mode = useThemeStore((state) => state.mode)
  return <ThemeProvider theme={getTheme(mode)}>{children}</ThemeProvider>
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <ThemedApp>
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
      </ThemedApp>
    </QueryClientProvider>
  </StrictMode>,
)