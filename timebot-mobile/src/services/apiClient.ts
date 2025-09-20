// API client configuration
import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { API_BASE_URL } from '../constants';
import type { ApiResponse, ErrorResponse } from '../types/api';

// Storage keys
const TOKEN_KEY = 'auth_token';
const REFRESH_TOKEN_KEY = 'refresh_token';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Request interceptor to add auth token
    this.client.interceptors.request.use(
      async (config) => {
        const token = await AsyncStorage.getItem(TOKEN_KEY);
        if (token && config.headers) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      async (error) => {
        if (error.response?.status === 401) {
          // Token expired, try to refresh or logout
          await this.handleUnauthorized();
        }
        return Promise.reject(this.handleError(error));
      }
    );
  }

  private async handleUnauthorized(): Promise<void> {
    // Clear stored tokens
    await AsyncStorage.multiRemove([TOKEN_KEY, REFRESH_TOKEN_KEY]);
    // TODO: Navigate to login screen or trigger logout action
  }

  private handleError(error: any): ErrorResponse {
    if (error.response) {
      // Server responded with error status
      return {
        success: false,
        error: error.response.data?.error || 'Server error',
        message: error.response.data?.message,
        details: error.response.data?.details,
      };
    } else if (error.request) {
      // Network error
      return {
        success: false,
        error: 'Network error',
        message: 'Please check your internet connection',
      };
    } else {
      // Other error
      return {
        success: false,
        error: 'Unknown error',
        message: error.message || 'An unexpected error occurred',
      };
    }
  }

  // Generic request methods
  async get<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<T> = await this.client.get(url, config);
      return {
        success: true,
        data: response.data,
      };
    } catch (error) {
      throw error;
    }
  }

  async post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<T> = await this.client.post(url, data, config);
      return {
        success: true,
        data: response.data,
      };
    } catch (error) {
      throw error;
    }
  }

  async put<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<T> = await this.client.put(url, data, config);
      return {
        success: true,
        data: response.data,
      };
    } catch (error) {
      throw error;
    }
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<T> = await this.client.delete(url, config);
      return {
        success: true,
        data: response.data,
      };
    } catch (error) {
      throw error;
    }
  }

  // Token management
  async setAuthToken(token: string): Promise<void> {
    await AsyncStorage.setItem(TOKEN_KEY, token);
  }

  async getAuthToken(): Promise<string | null> {
    return await AsyncStorage.getItem(TOKEN_KEY);
  }

  async clearAuthToken(): Promise<void> {
    await AsyncStorage.multiRemove([TOKEN_KEY, REFRESH_TOKEN_KEY]);
  }

  // Health check
  async healthCheck(): Promise<ApiResponse<{ status: string; service: string }>> {
    return this.get('/health');
  }
}

export default new ApiClient();
