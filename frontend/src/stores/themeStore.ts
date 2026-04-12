import { create } from 'zustand';

type ThemeMode = 'light' | 'dark';

interface ThemeState {
  mode: ThemeMode;
  toggleMode: () => void;
  setMode: (mode: ThemeMode) => void;
}

const getInitialMode = (): ThemeMode => {
  const stored = localStorage.getItem('themeMode');
  if (stored === 'light' || stored === 'dark') return stored;
  return 'light';
};

export const useThemeStore = create<ThemeState>((set) => ({
  mode: getInitialMode(),
  toggleMode: () => {
    set((state) => {
      const newMode = state.mode === 'light' ? 'dark' : 'light';
      localStorage.setItem('themeMode', newMode);
      return { mode: newMode };
    });
  },
  setMode: (mode) => {
    localStorage.setItem('themeMode', mode);
    set({ mode });
  },
}));