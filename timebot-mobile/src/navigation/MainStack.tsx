// Main app stack navigator with bottom tabs
import React from 'react';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { Ionicons } from '@expo/vector-icons';
import { observer } from 'mobx-react-lite';

// Types
import type { MainTabParamList } from '../types/navigation';

// Theme
import { colors, navigationStyles } from '../theme';

// Stores
import { AuthStore } from '../stores';

// Screens
import { SubmitTimeScreen } from '../screens';

// Temporary placeholder screens
import { View, Text, StyleSheet } from 'react-native';

const PlaceholderScreen: React.FC<{ title: string }> = ({ title }) => (
  <View style={styles.container}>
    <Text style={styles.text}>{title} Screen</Text>
    <Text style={styles.subtitle}>Coming soon...</Text>
  </View>
);

const Tab = createBottomTabNavigator<MainTabParamList>();

// Icon mapping for tabs
const getTabIcon = (routeName: keyof MainTabParamList, focused: boolean) => {
  let iconName: keyof typeof Ionicons.glyphMap;

  switch (routeName) {
    case 'Dashboard':
      iconName = focused ? 'stats-chart' : 'stats-chart-outline';
      break;
    case 'SubmitTime':
      iconName = focused ? 'time' : 'time-outline';
      break;
    case 'History':
      iconName = focused ? 'list' : 'list-outline';
      break;
    case 'Timecards':
      iconName = focused ? 'document-text' : 'document-text-outline';
      break;
    case 'Admin':
      iconName = focused ? 'person' : 'person-outline';
      break;
    default:
      iconName = 'help-outline';
  }

  return iconName;
};

const MainStack: React.FC = observer(() => {
  const isAdmin = AuthStore.isAdmin;

  return (
    <Tab.Navigator
      initialRouteName="Dashboard"
      screenOptions={({ route }) => ({
        tabBarIcon: ({ focused, color, size }) => {
          const iconName = getTabIcon(route.name, focused);
          return <Ionicons name={iconName} size={size} color={color} />;
        },
        tabBarActiveTintColor: colors.tabBar.active,
        tabBarInactiveTintColor: colors.tabBar.inactive,
        tabBarStyle: {
          backgroundColor: colors.tabBar.background,
          borderTopColor: colors.tabBar.border,
          borderTopWidth: 1,
          paddingBottom: 8,
          paddingTop: 8,
          height: 70,
        },
        tabBarLabelStyle: {
          fontSize: 11,
          fontWeight: '500',
          marginTop: 4,
        },
        headerShown: false,
      })}
    >
      <Tab.Screen
        name="Dashboard"
        component={() => <PlaceholderScreen title="Dashboard" />}
        options={{
          tabBarLabel: 'Dashboard',
        }}
      />
      <Tab.Screen
        name="SubmitTime"
        component={SubmitTimeScreen}
        options={{
          tabBarLabel: 'Submit',
        }}
      />
      <Tab.Screen
        name="History"
        component={() => <PlaceholderScreen title="History" />}
        options={{
          tabBarLabel: 'History',
        }}
      />
      <Tab.Screen
        name="Timecards"
        component={() => <PlaceholderScreen title="Timecards" />}
        options={{
          tabBarLabel: 'Timecards',
        }}
      />
      {isAdmin && (
        <Tab.Screen
          name="Admin"
          component={() => <PlaceholderScreen title="Admin" />}
          options={{
            tabBarLabel: 'Admin',
          }}
        />
      )}
    </Tab.Navigator>
  );
});

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: colors.background.primary,
    padding: 20,
  },
  text: {
    color: colors.text.primary,
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  subtitle: {
    color: colors.text.secondary,
    fontSize: 16,
  },
});

export default MainStack;
