// Root store that combines all other stores
import AuthStore from './AuthStore';
import TimecardStore from './TimecardStore';
import AppStore from './AppStore';

export class RootStore {
  authStore = AuthStore;
  timecardStore = TimecardStore;
  appStore = AppStore;

  // Initialize all stores
  initialize = async (): Promise<void> => {
    try {
      console.log('Initializing stores...');

      // Initialize app store first
      await this.appStore.initialize();

      // Initialize auth store (will restore persisted auth state)
      await this.authStore.initialize();

      // Initialize timecard store if user is authenticated
      if (this.authStore.isAuthenticated) {
        await this.timecardStore.initialize();
      }

      console.log('All stores initialized successfully');
    } catch (error) {
      console.error('Store initialization failed:', error);
      this.appStore.showError(
        'Initialization Error',
        'Failed to initialize app. Please restart the application.'
      );
    }
  };

  // Reset all stores (useful for logout)
  reset = async (): Promise<void> => {
    await Promise.all([
      this.authStore.logout(),
      this.timecardStore.clearError(),
      this.appStore.reset(),
    ]);
  };
}

// Create root store instance
export const rootStore = new RootStore();

// Export individual stores for easier access
export { AuthStore, TimecardStore, AppStore };

// Export default for convenience
export default rootStore;
