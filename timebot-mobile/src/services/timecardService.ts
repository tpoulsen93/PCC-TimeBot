// Timecard API service
import apiClient from './apiClient';
import type {
  ApiResponse,
  PaginatedResponse,
  TimeEntryResponse,
  TimecardResponse,
  CreateTimeEntryRequest,
  UpdateTimeEntryRequest,
} from '../types/api';
import type { TimeEntry, TimeCard } from '../types/timecard';

class TimecardService {
  /**
   * Get all time entries for current user
   */
  async getTimeEntries(params?: {
    startDate?: string;
    endDate?: string;
    page?: number;
    pageSize?: number;
  }): Promise<ApiResponse<PaginatedResponse<TimeEntryResponse>>> {
    const queryParams = new URLSearchParams();

    if (params?.startDate) queryParams.append('startDate', params.startDate);
    if (params?.endDate) queryParams.append('endDate', params.endDate);
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.pageSize) queryParams.append('pageSize', params.pageSize.toString());

    const url = `/timecards${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return apiClient.get<PaginatedResponse<TimeEntryResponse>>(url);
  }

  /**
   * Get specific time entry by ID
   */
  async getTimeEntry(id: number): Promise<ApiResponse<TimeEntryResponse>> {
    return apiClient.get<TimeEntryResponse>(`/timecards/${id}`);
  }

  /**
   * Create new time entry
   */
  async createTimeEntry(data: CreateTimeEntryRequest): Promise<ApiResponse<TimeEntryResponse>> {
    return apiClient.post<TimeEntryResponse>('/timecards', data);
  }

  /**
   * Update existing time entry
   */
  async updateTimeEntry(id: number, data: UpdateTimeEntryRequest): Promise<ApiResponse<TimeEntryResponse>> {
    return apiClient.put<TimeEntryResponse>(`/timecards/${id}`, data);
  }

  /**
   * Delete time entry
   */
  async deleteTimeEntry(id: number): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(`/timecards/${id}`);
  }

  /**
   * Get timecard summary for a specific period
   */
  async getTimecardSummary(startDate: string, endDate: string): Promise<ApiResponse<TimecardResponse>> {
    const params = new URLSearchParams({
      startDate,
      endDate,
    });
    return apiClient.get<TimecardResponse>(`/timecards/summary?${params.toString()}`);
  }

  /**
   * Get current week timecard
   */
  async getCurrentWeekTimecard(): Promise<ApiResponse<TimecardResponse>> {
    return apiClient.get<TimecardResponse>('/timecards/current-week');
  }

  /**
   * Submit timecard for approval
   */
  async submitTimecard(timecardId: number): Promise<ApiResponse<TimecardResponse>> {
    return apiClient.post<TimecardResponse>(`/timecards/${timecardId}/submit`);
  }

  /**
   * Get payroll report for date range
   */
  async getPayrollReport(startDate: string, endDate: string): Promise<ApiResponse<any>> {
    const params = new URLSearchParams({
      startDate,
      endDate,
    });
    return apiClient.get<any>(`/reports/payroll?${params.toString()}`);
  }

  /**
   * Email timecard to user
   */
  async emailTimecard(timecardId: number): Promise<ApiResponse<{ message: string }>> {
    return apiClient.post<{ message: string }>(`/timecards/${timecardId}/email`);
  }

  /**
   * Get available projects for time entry
   */
  async getProjects(): Promise<ApiResponse<Array<{ id: string; name: string; description?: string }>>> {
    return apiClient.get<Array<{ id: string; name: string; description?: string }>>('/projects');
  }

  /**
   * Bulk create time entries
   */
  async bulkCreateTimeEntries(entries: CreateTimeEntryRequest[]): Promise<ApiResponse<TimeEntryResponse[]>> {
    return apiClient.post<TimeEntryResponse[]>('/timecards/bulk', { entries });
  }

  /**
   * Get time entry statistics
   */
  async getTimeEntryStats(params?: {
    startDate?: string;
    endDate?: string;
  }): Promise<ApiResponse<{
    totalHours: number;
    regularHours: number;
    overtimeHours: number;
    totalEntries: number;
    averageHoursPerDay: number;
  }>> {
    const queryParams = new URLSearchParams();

    if (params?.startDate) queryParams.append('startDate', params.startDate);
    if (params?.endDate) queryParams.append('endDate', params.endDate);

    const url = `/timecards/stats${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return apiClient.get(url);
  }
}

export default new TimecardService();
