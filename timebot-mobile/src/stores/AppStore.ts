// App store for global application state
import { makeAutoObservable } from 'mobx';
import { makePersistable } from 'mobx-persist-store';
import AsyncStorage from '@react-native-async-storage/async-storage';

class AppStore {
  // Observables
  isInitialized = false;
  isOnline = true;
  theme: 'light' | 'dark' = 'dark';
  lastSyncDate: string | null = null;
  notifications: Array<{
    id: string;
    type: 'success' | 'error' | 'warning' | 'info';
    title: string;
    message: string;
    timestamp: string;
    read: boolean;
  }> = [];

  // App settings
  settings = {
    autoSubmitTimecards: false,
    reminderTime: '17:00', // 5 PM
    reminderEnabled: true,
    biometricEnabled: false,
    keepLoggedIn: true,
  };

  constructor() {
    makeAutoObservable(this);

    // Persist app settings and preferences
    makePersistable(this, {
      name: 'AppStore',
      properties: ['theme', 'settings', 'notifications'],
      storage: AsyncStorage,
    });
  }

  // Computed
  get unreadNotificationCount(): number {
    return this.notifications.filter(n => !n.read).length;
  }

  get recentNotifications(): typeof this.notifications {
    return this.notifications
      .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
      .slice(0, 10);
  }

  // Actions
  setInitialized = (initialized: boolean) => {
    this.isInitialized = initialized;
  };

  setOnlineStatus = (isOnline: boolean) => {
    this.isOnline = isOnline;
  };

  setTheme = (theme: 'light' | 'dark') => {
    this.theme = theme;
  };

  updateSetting = <K extends keyof typeof this.settings>(
    key: K,
    value: typeof this.settings[K]
  ) => {
    this.settings[key] = value;
  };

  addNotification = (notification: {
    type: 'success' | 'error' | 'warning' | 'info';
    title: string;
    message: string;
  }) => {
    const newNotification = {
      id: Date.now().toString(),
      ...notification,
      timestamp: new Date().toISOString(),
      read: false,
    };

    this.notifications.unshift(newNotification);

    // Keep only last 50 notifications
    if (this.notifications.length > 50) {
      this.notifications = this.notifications.slice(0, 50);
    }
  };

  markNotificationAsRead = (id: string) => {
    const notification = this.notifications.find(n => n.id === id);
    if (notification) {
      notification.read = true;
    }
  };

  markAllNotificationsAsRead = () => {
    this.notifications.forEach(notification => {
      notification.read = true;
    });
  };

  removeNotification = (id: string) => {
    this.notifications = this.notifications.filter(n => n.id !== id);
  };

  clearAllNotifications = () => {
    this.notifications = [];
  };

  setLastSyncDate = (date: string) => {
    this.lastSyncDate = date;
  };

  // Helper methods for showing notifications
  showSuccess = (title: string, message: string) => {
    this.addNotification({ type: 'success', title, message });
  };

  showError = (title: string, message: string) => {
    this.addNotification({ type: 'error', title, message });
  };

  showWarning = (title: string, message: string) => {
    this.addNotification({ type: 'warning', title, message });
  };

  showInfo = (title: string, message: string) => {
    this.addNotification({ type: 'info', title, message });
  };

  // App lifecycle methods
  initialize = async (): Promise<void> => {
    try {
      // Any app initialization logic here
      console.log('Initializing App Store...');

      this.setInitialized(true);
    } catch (error) {
      console.error('Failed to initialize app store:', error);
    }
  };

  reset = async (): Promise<void> => {
    // Reset all state (useful for logout)
    this.isInitialized = false;
    this.isOnline = true;
    this.lastSyncDate = null;
    this.notifications = [];

    // Keep theme and settings on reset
  };
}

export default new AppStore();
