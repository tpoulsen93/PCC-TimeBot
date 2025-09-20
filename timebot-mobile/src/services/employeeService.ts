// Employee API service
import apiClient from './apiClient';
import type {
  ApiResponse,
  PaginatedResponse,
  EmployeeResponse,
  CreateEmployeeRequest,
  UpdateEmployeeRequest,
} from '../types/api';

class EmployeeService {
  /**
   * Get all employees (admin only)
   */
  async getEmployees(params?: {
    page?: number;
    pageSize?: number;
    search?: string;
    isActive?: boolean;
  }): Promise<ApiResponse<PaginatedResponse<EmployeeResponse>>> {
    const queryParams = new URLSearchParams();

    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.pageSize) queryParams.append('pageSize', params.pageSize.toString());
    if (params?.search) queryParams.append('search', params.search);
    if (params?.isActive !== undefined) queryParams.append('isActive', params.isActive.toString());

    const url = `/employees${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return apiClient.get<PaginatedResponse<EmployeeResponse>>(url);
  }

  /**
   * Get specific employee by ID
   */
  async getEmployee(id: number): Promise<ApiResponse<EmployeeResponse>> {
    return apiClient.get<EmployeeResponse>(`/employees/${id}`);
  }

  /**
   * Create new employee (admin only)
   */
  async createEmployee(data: CreateEmployeeRequest): Promise<ApiResponse<EmployeeResponse>> {
    return apiClient.post<EmployeeResponse>('/employees', data);
  }

  /**
   * Update existing employee (admin only)
   */
  async updateEmployee(id: number, data: UpdateEmployeeRequest): Promise<ApiResponse<EmployeeResponse>> {
    return apiClient.put<EmployeeResponse>(`/employees/${id}`, data);
  }

  /**
   * Deactivate employee (admin only)
   */
  async deactivateEmployee(id: number): Promise<ApiResponse<EmployeeResponse>> {
    return apiClient.put<EmployeeResponse>(`/employees/${id}`, { isActive: false });
  }

  /**
   * Reactivate employee (admin only)
   */
  async reactivateEmployee(id: number): Promise<ApiResponse<EmployeeResponse>> {
    return apiClient.put<EmployeeResponse>(`/employees/${id}`, { isActive: true });
  }

  /**
   * Get employee timecard history (admin only)
   */
  async getEmployeeTimecards(employeeId: number, params?: {
    startDate?: string;
    endDate?: string;
    page?: number;
    pageSize?: number;
  }): Promise<ApiResponse<PaginatedResponse<any>>> {
    const queryParams = new URLSearchParams();

    if (params?.startDate) queryParams.append('startDate', params.startDate);
    if (params?.endDate) queryParams.append('endDate', params.endDate);
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.pageSize) queryParams.append('pageSize', params.pageSize.toString());

    const url = `/employees/${employeeId}/timecards${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return apiClient.get<PaginatedResponse<any>>(url);
  }

  /**
   * Approve employee timecard (admin only)
   */
  async approveTimecard(timecardId: number): Promise<ApiResponse<{ message: string }>> {
    return apiClient.post<{ message: string }>(`/timecards/${timecardId}/approve`);
  }

  /**
   * Reject employee timecard (admin only)
   */
  async rejectTimecard(timecardId: number, reason?: string): Promise<ApiResponse<{ message: string }>> {
    return apiClient.post<{ message: string }>(`/timecards/${timecardId}/reject`, { reason });
  }

  /**
   * Get employee dashboard stats (admin only)
   */
  async getEmployeeStats(params?: {
    startDate?: string;
    endDate?: string;
  }): Promise<ApiResponse<{
    totalEmployees: number;
    activeEmployees: number;
    totalHoursThisWeek: number;
    pendingTimecards: number;
    approvedTimecards: number;
  }>> {
    const queryParams = new URLSearchParams();

    if (params?.startDate) queryParams.append('startDate', params.startDate);
    if (params?.endDate) queryParams.append('endDate', params.endDate);

    const url = `/employees/stats${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return apiClient.get(url);
  }

  /**
   * Get pending submissions for admin dashboard
   */
  async getPendingSubmissions(params?: {
    page?: number;
    pageSize?: number;
  }): Promise<ApiResponse<PaginatedResponse<any>>> {
    const queryParams = new URLSearchParams();

    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.pageSize) queryParams.append('pageSize', params.pageSize.toString());

    const url = `/timecards/pending${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return apiClient.get<PaginatedResponse<any>>(url);
  }

  /**
   * Send reminder emails to employees
   */
  async sendReminders(employeeIds?: number[]): Promise<ApiResponse<{ message: string; sent: number }>> {
    return apiClient.post<{ message: string; sent: number }>('/employees/send-reminders', {
      employeeIds,
    });
  }

  /**
   * Export employee data (admin only)
   */
  async exportEmployeeData(format: 'csv' | 'xlsx' = 'csv'): Promise<ApiResponse<{ downloadUrl: string }>> {
    return apiClient.get<{ downloadUrl: string }>(`/employees/export?format=${format}`);
  }

  /**
   * Get employee work summary
   */
  async getEmployeeWorkSummary(employeeId: number, params: {
    startDate: string;
    endDate: string;
  }): Promise<ApiResponse<{
    employeeId: number;
    employeeName: string;
    totalHours: number;
    regularHours: number;
    overtimeHours: number;
    totalDays: number;
    averageHoursPerDay: number;
    projects: Array<{
      name: string;
      hours: number;
    }>;
  }>> {
    const queryParams = new URLSearchParams(params);
    return apiClient.get(`/employees/${employeeId}/summary?${queryParams.toString()}`);
  }
}

export default new EmployeeService();
