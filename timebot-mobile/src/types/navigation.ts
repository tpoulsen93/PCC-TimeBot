// Navigation types
import { StackNavigationProp } from '@react-navigation/stack';
import { BottomTabNavigationProp } from '@react-navigation/bottom-tabs';

// Root Stack Navigator
export type RootStackParamList = {
  Auth: undefined;
  Main: undefined;
};

// Auth Stack Navigator
export type AuthStackParamList = {
  Login: undefined;
  Signup: undefined;
  ForgotPassword: undefined;
};

// Main Tab Navigator
export type MainTabParamList = {
  Dashboard: undefined;
  SubmitTime: undefined;
  History: undefined;
  Timecards: undefined;
  Admin?: undefined;
};

// Navigation prop types
export type RootNavigationProp = StackNavigationProp<RootStackParamList>;
export type AuthNavigationProp = StackNavigationProp<AuthStackParamList>;
export type MainTabNavigationProp = BottomTabNavigationProp<MainTabParamList>;
