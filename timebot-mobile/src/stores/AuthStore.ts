// Authentication store using MobX
import { makeAutoObservable, runInAction } from 'mobx';
import { makePersistable } from 'mobx-persist-store';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { authService } from '../services';
import type { User, LoginCredentials, SignupData } from '../types/user';
import type { LoginResponse } from '../types/api';

class AuthStore {
  // Observables
  user: User | null = null;
  token: string | null = null;
  isLoading = false;
  error: string | null = null;

  constructor() {
    makeAutoObservable(this);

    // Make authentication data persistent
    makePersistable(this, {
      name: 'AuthStore',
      properties: ['user', 'token'],
      storage: AsyncStorage,
    });
  }

  // Computed
  get isAuthenticated(): boolean {
    return !!this.token && !!this.user;
  }

  get isAdmin(): boolean {
    return this.user?.isAdmin ?? false;
  }

  get fullName(): string {
    if (!this.user) return '';
    return `${this.user.firstName} ${this.user.lastName}`;
  }

  // Actions
  setLoading = (loading: boolean) => {
    this.isLoading = loading;
    if (loading) {
      this.error = null;
    }
  };

  setError = (error: string | null) => {
    this.error = error;
  };

  setUser = (user: User | null) => {
    this.user = user;
  };

  setToken = (token: string | null) => {
    this.token = token;
  };

  // Auth methods
  login = async (credentials: LoginCredentials): Promise<boolean> => {
    this.setLoading(true);

    try {
      const response = await authService.login(credentials);

      if (response.success && response.data) {
        runInAction(() => {
          this.user = response.data!.user;
          this.token = response.data!.token;
          this.error = null;
        });
        return true;
      } else {
        throw new Error('Login failed');
      }
    } catch (error: any) {
      runInAction(() => {
        this.error = error.error || error.message || 'Login failed';
      });
      return false;
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };

  signup = async (userData: SignupData): Promise<boolean> => {
    this.setLoading(true);

    try {
      const response = await authService.signup(userData);

      if (response.success && response.data) {
        runInAction(() => {
          this.user = response.data!.user;
          this.token = response.data!.token;
          this.error = null;
        });
        return true;
      } else {
        throw new Error('Signup failed');
      }
    } catch (error: any) {
      runInAction(() => {
        this.error = error.error || error.message || 'Signup failed';
      });
      return false;
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };

  logout = async (): Promise<void> => {
    this.setLoading(true);

    try {
      await authService.logout();
    } catch (error) {
      console.warn('Logout API call failed:', error);
    } finally {
      runInAction(() => {
        this.user = null;
        this.token = null;
        this.error = null;
        this.isLoading = false;
      });
    }
  };

  forgotPassword = async (email: string): Promise<boolean> => {
    this.setLoading(true);

    try {
      const response = await authService.forgotPassword({ email });

      if (response.success) {
        runInAction(() => {
          this.error = null;
        });
        return true;
      } else {
        throw new Error('Failed to send reset email');
      }
    } catch (error: any) {
      runInAction(() => {
        this.error = error.error || error.message || 'Failed to send reset email';
      });
      return false;
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };

  refreshToken = async (): Promise<boolean> => {
    try {
      const response = await authService.refreshToken();

      if (response.success && response.data?.token) {
        runInAction(() => {
          this.token = response.data!.token;
        });
        return true;
      }
      return false;
    } catch (error) {
      console.warn('Token refresh failed:', error);
      return false;
    }
  };

  getCurrentUser = async (): Promise<void> => {
    if (!this.token) return;

    try {
      const response = await authService.getCurrentUser();

      if (response.success && response.data) {
        runInAction(() => {
          this.user = response.data!;
        });
      }
    } catch (error) {
      console.warn('Failed to get current user:', error);
      // Don't logout on this error, might be network issue
    }
  };

  clearError = () => {
    this.error = null;
  };

  // Initialize store (call on app startup)
  initialize = async (): Promise<void> => {
    if (this.token && !this.user) {
      await this.getCurrentUser();
    }
  };
}

export default new AuthStore();
