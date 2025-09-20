import React from 'react';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { ThemeProvider } from './src/theme/ThemeContext';
import { StatusBar } from './src/components/common/StatusBar';
import { RootNavigator } from './src/navigation';

export default function App() {
  return (
    <SafeAreaProvider>
      <ThemeProvider>
        <StatusBar />
        <RootNavigator />
      </ThemeProvider>
    </SafeAreaProvider>
  );
}
