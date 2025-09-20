// User and Employee types

export interface User {
  id: number;
  firstName: string;
  lastName: string;
  email: string;
  phone?: string;
  isAdmin: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface Employee {
  id: number;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  supervisorId?: number;
  supervisor?: {
    firstName: string;
    lastName: string;
  };
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface SignupData {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  phone?: string;
}

export interface ForgotPasswordData {
  email: string;
}
