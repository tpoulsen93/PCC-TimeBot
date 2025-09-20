import React, { createContext, useContext, ReactNode } from 'react';
import { theme, Theme } from './index';

interface ThemeContextType {
  theme: Theme;
  isDark: boolean;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

interface ThemeProviderProps {
  children: ReactNode;
}

export const ThemeProvider: React.FC<ThemeProviderProps> = ({ children }) => {
  const contextValue: ThemeContextType = {
    theme,
    isDark: true, // Always dark theme for now, could be made dynamic later
  };

  return (
    <ThemeContext.Provider value={contextValue}>
      {children}
    </ThemeContext.Provider>
  );
};

export const useTheme = (): ThemeContextType => {
  const context = useContext(ThemeContext);
  if (context === undefined) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
};

// Convenience hook for just the theme object
export const useThemeStyles = () => {
  const { theme } = useTheme();
  return theme;
};
