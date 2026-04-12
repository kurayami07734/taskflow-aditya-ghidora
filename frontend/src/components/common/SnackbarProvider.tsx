import { createContext, useContext, useState, useCallback } from 'react';
import { Snackbar, Alert, AlertColor } from '@mui/material';

interface SnackbarState {
  open: boolean;
  message: string;
  severity: AlertColor;
}

interface SnackbarContextType {
  showError: (message: string) => void;
  showSuccess: (message: string) => void;
  showWarning: (message: string) => void;
  showInfo: (message: string) => void;
}

const SnackbarContext = createContext<SnackbarContextType | null>(null);

export const useSnackbar = () => {
  const context = useContext(SnackbarContext);
  if (!context) {
    throw new Error('useSnackbar must be used within SnackbarProvider');
  }
  return context;
};

export const SnackbarProvider = ({ children }: { children: React.ReactNode }) => {
  const [state, setState] = useState<SnackbarState>({
    open: false,
    message: '',
    severity: 'error',
  });

  const showMessage = useCallback((message: string, severity: AlertColor) => {
    setState({ open: true, message, severity });
  }, []);

  const showError = useCallback((message: string) => showMessage(message, 'error'), [showMessage]);
  const showSuccess = useCallback((message: string) => showMessage(message, 'success'), [showMessage]);
  const showWarning = useCallback((message: string) => showMessage(message, 'warning'), [showMessage]);
  const showInfo = useCallback((message: string) => showMessage(message, 'info'), [showMessage]);

  const handleClose = () => setState((prev) => ({ ...prev, open: false }));

  return (
    <SnackbarContext.Provider value={{ showError, showSuccess, showWarning, showInfo }}>
      {children}
      <Snackbar
        open={state.open}
        autoHideDuration={6000}
        onClose={handleClose}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert onClose={handleClose} severity={state.severity} sx={{ width: '100%' }}>
          {state.message}
        </Alert>
      </Snackbar>
    </SnackbarContext.Provider>
  );
};