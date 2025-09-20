// Root navigation container
import React from 'react';
import { NavigationContainer, DefaultTheme } from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';
import { observer } from 'mobx-react-lite';

// Navigation stacks
import AuthStack from './AuthStack';
import MainStack from './MainStack';

// Types
import type { RootStackParamList } from '../types/navigation';

// Stores
import { AuthStore } from '../stores';

// Theme
import { colors } from '../theme';

// Dark theme for React Navigation
const DarkNavigationTheme = {
  ...DefaultTheme,
  dark: true,
  colors: {
    ...DefaultTheme.colors,
    primary: colors.purple.primary,
    background: colors.background.primary,
    card: colors.background.secondary,
    text: colors.text.primary,
    border: colors.border.primary,
    notification: colors.purple.primary,
  },
};

const Stack = createStackNavigator<RootStackParamList>();

const RootNavigator: React.FC = observer(() => {
  const isAuthenticated = AuthStore.isAuthenticated;

  return (
    <NavigationContainer theme={DarkNavigationTheme}>
      <Stack.Navigator
        screenOptions={{
          headerShown: false,
          gestureEnabled: false,
        }}
      >
        {isAuthenticated ? (
          <Stack.Screen name="Main" component={MainStack} />
        ) : (
          <Stack.Screen name="Auth" component={AuthStack} />
        )}
      </Stack.Navigator>
    </NavigationContainer>
  );
});

export default RootNavigator;
