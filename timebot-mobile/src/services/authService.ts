// Authentication API service
import apiClient from './apiClient';
import type {
  LoginResponse,
  RefreshTokenResponse,
  ApiResponse,
} from '../types/api';
import type { LoginCredentials, SignupData, ForgotPasswordData } from '../types/user';

class AuthService {
  /**
   * Login user with email and password
   */
  async login(credentials: LoginCredentials): Promise<ApiResponse<LoginResponse>> {
    const response = await apiClient.post<LoginResponse>('/auth/login', credentials);

    if (response.success && response.data?.token) {
      // Store the auth token
      await apiClient.setAuthToken(response.data.token);
    }

    return response;
  }

  /**
   * Register new user
   */
  async signup(userData: SignupData): Promise<ApiResponse<LoginResponse>> {
    const response = await apiClient.post<LoginResponse>('/auth/signup', userData);

    if (response.success && response.data?.token) {
      // Store the auth token
      await apiClient.setAuthToken(response.data.token);
    }

    return response;
  }

  /**
   * Logout user
   */
  async logout(): Promise<void> {
    try {
      // Call logout endpoint to invalidate token on server
      await apiClient.post('/auth/logout');
    } catch (error) {
      // Even if server call fails, clear local token
      console.warn('Logout API call failed:', error);
    } finally {
      // Always clear local auth token
      await apiClient.clearAuthToken();
    }
  }

  /**
   * Refresh authentication token
   */
  async refreshToken(): Promise<ApiResponse<RefreshTokenResponse>> {
    return apiClient.post<RefreshTokenResponse>('/auth/refresh');
  }

  /**
   * Send forgot password email
   */
  async forgotPassword(data: ForgotPasswordData): Promise<ApiResponse<{ message: string }>> {
    return apiClient.post<{ message: string }>('/auth/forgot-password', data);
  }

  /**
   * Reset password with token
   */
  async resetPassword(token: string, newPassword: string): Promise<ApiResponse<{ message: string }>> {
    return apiClient.post<{ message: string }>('/auth/reset-password', {
      token,
      password: newPassword,
    });
  }

  /**
   * Verify email address
   */
  async verifyEmail(token: string): Promise<ApiResponse<{ message: string }>> {
    return apiClient.post<{ message: string }>('/auth/verify-email', { token });
  }

  /**
   * Resend verification email
   */
  async resendVerification(email: string): Promise<ApiResponse<{ message: string }>> {
    return apiClient.post<{ message: string }>('/auth/resend-verification', { email });
  }

  /**
   * Check if user is authenticated
   */
  async isAuthenticated(): Promise<boolean> {
    const token = await apiClient.getAuthToken();
    return !!token;
  }

  /**
   * Get current user profile
   */
  async getCurrentUser(): Promise<ApiResponse<LoginResponse['user']>> {
    return apiClient.get<LoginResponse['user']>('/auth/me');
  }
}

export default new AuthService();
