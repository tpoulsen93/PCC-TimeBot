// Authentication stack navigator
import React from 'react';
import { createStackNavigator } from '@react-navigation/stack';

// Screens
import { LoginScreen, SignupScreen } from '../screens/auth';

// Types
import type { AuthStackParamList } from '../types/navigation';

// Theme
import { colors } from '../theme';

const Stack = createStackNavigator<AuthStackParamList>();

const AuthStack: React.FC = () => {
  return (
    <Stack.Navigator
      initialRouteName="Login"
      screenOptions={{
        headerShown: false,
        gestureEnabled: false,
        cardStyle: { backgroundColor: colors.background.primary },
      }}
    >
      <Stack.Screen name="Login" component={LoginScreen} />
      <Stack.Screen name="Signup" component={SignupScreen} />
    </Stack.Navigator>
  );
};

export default AuthStack;
