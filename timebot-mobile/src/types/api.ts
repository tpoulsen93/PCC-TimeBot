// API response types

export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
  };
}

export interface ErrorResponse {
  success: false;
  error: string;
  message?: string;
  details?: Record<string, string[]>;
}

// Health check response
export interface HealthCheckResponse {
  status: string;
  service: string;
}

// Authentication responses
export interface LoginResponse {
  user: {
    id: number;
    firstName: string;
    lastName: string;
    email: string;
    isAdmin: boolean;
  };
  token: string;
  refreshToken?: string;
}

export interface RefreshTokenResponse {
  token: string;
  refreshToken?: string;
}

// Time entry responses
export interface TimeEntryResponse {
  id: number;
  employeeId: number;
  date: string;
  startTime: string;
  endTime: string;
  hours: number;
  project?: string;
  description?: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface TimecardResponse {
  id: number;
  employeeId: number;
  startDate: string;
  endDate: string;
  totalHours: number;
  status: string;
  entries: TimeEntryResponse[];
  createdAt: string;
  updatedAt: string;
}

// Employee responses
export interface EmployeeResponse {
  id: number;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  isActive: boolean;
  supervisor?: {
    id: number;
    firstName: string;
    lastName: string;
  };
  createdAt: string;
  updatedAt: string;
}

// Request types
export interface CreateTimeEntryRequest {
  date: string;
  startTime: string;
  endTime: string;
  project?: string;
  description?: string;
}

export interface UpdateTimeEntryRequest extends Partial<CreateTimeEntryRequest> {
  // All fields are optional for updates
}

export interface CreateEmployeeRequest {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  supervisorId?: number;
}

export interface UpdateEmployeeRequest extends Partial<CreateEmployeeRequest> {
  isActive?: boolean;
}
