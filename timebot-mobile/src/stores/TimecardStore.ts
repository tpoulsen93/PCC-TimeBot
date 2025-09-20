// Timecard store using MobX
import { makeAutoObservable, runInAction } from 'mobx';
import { makePersistable } from 'mobx-persist-store';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { timecardService } from '../services';
import { getCurrentDate, getCurrentWeekRange } from '../utils/dateHelpers';
import { calculateHours } from '../utils/timeCalculations';
import type { TimeEntry, TimeCard, CreateTimeEntryData } from '../types/timecard';
import type { TimeEntryResponse, TimecardResponse } from '../types/api';

class TimecardStore {
  // Observables
  timeEntries: TimeEntry[] = [];
  currentWeekTimecard: TimeCard | null = null;
  selectedTimeEntry: TimeEntry | null = null;
  isLoading = false;
  isSubmitting = false;
  error: string | null = null;
  lastFetchDate: string | null = null;

  // Form state for time submission
  submissionForm = {
    date: getCurrentDate(),
    startTime: '08:00',
    endTime: '17:00',
    project: '',
    description: '',
  };

  constructor() {
    makeAutoObservable(this);

    // Persist timecard data for offline access
    makePersistable(this, {
      name: 'TimecardStore',
      properties: ['timeEntries', 'currentWeekTimecard', 'lastFetchDate'],
      storage: AsyncStorage,
    });
  }

  // Computed
  get calculatedHours(): number {
    return calculateHours(this.submissionForm.startTime, this.submissionForm.endTime);
  }

  get currentWeekEntries(): TimeEntry[] {
    const { start, end } = getCurrentWeekRange();
    return this.timeEntries.filter(
      entry => entry.date >= start && entry.date <= end
    );
  }

  get currentWeekTotalHours(): number {
    return this.currentWeekEntries.reduce((total, entry) => total + entry.hours, 0);
  }

  get pendingEntries(): TimeEntry[] {
    return this.timeEntries.filter(entry => entry.status === 'pending');
  }

  get approvedEntries(): TimeEntry[] {
    return this.timeEntries.filter(entry => entry.status === 'approved');
  }

  get todayEntries(): TimeEntry[] {
    const today = getCurrentDate();
    return this.timeEntries.filter(entry => entry.date === today);
  }

  get hasUnsavedChanges(): boolean {
    // Check if form has been modified
    return (
      this.submissionForm.date !== getCurrentDate() ||
      this.submissionForm.startTime !== '08:00' ||
      this.submissionForm.endTime !== '17:00' ||
      this.submissionForm.project !== '' ||
      this.submissionForm.description !== ''
    );
  }

  // Actions
  setLoading = (loading: boolean) => {
    this.isLoading = loading;
    if (loading) {
      this.error = null;
    }
  };

  setSubmitting = (submitting: boolean) => {
    this.isSubmitting = submitting;
    if (submitting) {
      this.error = null;
    }
  };

  setError = (error: string | null) => {
    this.error = error;
  };

  setSelectedTimeEntry = (entry: TimeEntry | null) => {
    this.selectedTimeEntry = entry;
  };

  // Form actions
  updateSubmissionForm = (field: keyof typeof this.submissionForm, value: string) => {
    this.submissionForm[field] = value;
  };

  resetSubmissionForm = () => {
    this.submissionForm = {
      date: getCurrentDate(),
      startTime: '08:00',
      endTime: '17:00',
      project: '',
      description: '',
    };
  };

  // API methods
  fetchTimeEntries = async (params?: {
    startDate?: string;
    endDate?: string;
    forceRefresh?: boolean;
  }): Promise<void> => {
    if (!params?.forceRefresh && this.lastFetchDate === getCurrentDate()) {
      return; // Already fetched today
    }

    this.setLoading(true);

    try {
      const response = await timecardService.getTimeEntries(params);

      if (response.success && response.data) {
        runInAction(() => {
          this.timeEntries = response.data!.data.map(this.mapTimeEntryResponse);
          this.lastFetchDate = getCurrentDate();
          this.error = null;
        });
      }
    } catch (error: any) {
      runInAction(() => {
        this.error = error.error || error.message || 'Failed to fetch time entries';
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };

  fetchCurrentWeekTimecard = async (): Promise<void> => {
    this.setLoading(true);

    try {
      const response = await timecardService.getCurrentWeekTimecard();

      if (response.success && response.data) {
        runInAction(() => {
          this.currentWeekTimecard = this.mapTimecardResponse(response.data!);
          this.error = null;
        });
      }
    } catch (error: any) {
      runInAction(() => {
        this.error = error.error || error.message || 'Failed to fetch current week timecard';
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  };

  createTimeEntry = async (data?: CreateTimeEntryData): Promise<boolean> => {
    const entryData = data || {
      date: this.submissionForm.date,
      startTime: this.submissionForm.startTime,
      endTime: this.submissionForm.endTime,
      project: this.submissionForm.project || undefined,
      description: this.submissionForm.description || undefined,
    };

    this.setSubmitting(true);

    try {
      const response = await timecardService.createTimeEntry(entryData);

      if (response.success && response.data) {
        const newEntry = this.mapTimeEntryResponse(response.data);

        runInAction(() => {
          this.timeEntries.unshift(newEntry);
          this.error = null;
        });

        if (!data) {
          this.resetSubmissionForm();
        }

        return true;
      }
      return false;
    } catch (error: any) {
      runInAction(() => {
        this.error = error.error || error.message || 'Failed to create time entry';
      });
      return false;
    } finally {
      runInAction(() => {
        this.isSubmitting = false;
      });
    }
  };

  updateTimeEntry = async (id: number, updates: Partial<CreateTimeEntryData>): Promise<boolean> => {
    this.setSubmitting(true);

    try {
      const response = await timecardService.updateTimeEntry(id, updates);

      if (response.success && response.data) {
        const updatedEntry = this.mapTimeEntryResponse(response.data);

        runInAction(() => {
          const index = this.timeEntries.findIndex(entry => entry.id === id);
          if (index !== -1) {
            this.timeEntries[index] = updatedEntry;
          }
          this.error = null;
        });

        return true;
      }
      return false;
    } catch (error: any) {
      runInAction(() => {
        this.error = error.error || error.message || 'Failed to update time entry';
      });
      return false;
    } finally {
      runInAction(() => {
        this.isSubmitting = false;
      });
    }
  };

  deleteTimeEntry = async (id: number): Promise<boolean> => {
    this.setSubmitting(true);

    try {
      const response = await timecardService.deleteTimeEntry(id);

      if (response.success) {
        runInAction(() => {
          this.timeEntries = this.timeEntries.filter(entry => entry.id !== id);
          this.error = null;
        });

        return true;
      }
      return false;
    } catch (error: any) {
      runInAction(() => {
        this.error = error.error || error.message || 'Failed to delete time entry';
      });
      return false;
    } finally {
      runInAction(() => {
        this.isSubmitting = false;
      });
    }
  };

  submitTimecard = async (timecardId: number): Promise<boolean> => {
    this.setSubmitting(true);

    try {
      const response = await timecardService.submitTimecard(timecardId);

      if (response.success) {
        runInAction(() => {
          // Update the status of related time entries
          this.timeEntries = this.timeEntries.map(entry => {
            if (entry.date >= this.currentWeekTimecard?.startDate! &&
                entry.date <= this.currentWeekTimecard?.endDate!) {
              return { ...entry, status: 'pending' as const };
            }
            return entry;
          });
          this.error = null;
        });

        return true;
      }
      return false;
    } catch (error: any) {
      runInAction(() => {
        this.error = error.error || error.message || 'Failed to submit timecard';
      });
      return false;
    } finally {
      runInAction(() => {
        this.isSubmitting = false;
      });
    }
  };

  // Helper methods
  private mapTimeEntryResponse = (response: TimeEntryResponse): TimeEntry => {
    return {
      id: response.id,
      employeeId: response.employeeId,
      date: response.date,
      startTime: response.startTime,
      endTime: response.endTime,
      hours: response.hours,
      project: response.project,
      description: response.description,
      status: response.status as 'pending' | 'approved' | 'rejected',
      createdAt: response.createdAt,
      updatedAt: response.updatedAt,
    };
  };

  private mapTimecardResponse = (response: TimecardResponse): TimeCard => {
    return {
      id: response.id,
      employeeId: response.employeeId,
      employee: {
        firstName: '', // Will be filled from auth store or separate call
        lastName: '',
        email: '',
      },
      startDate: response.startDate,
      endDate: response.endDate,
      entries: response.entries.map(this.mapTimeEntryResponse),
      totalHours: response.totalHours,
      regularHours: Math.min(response.totalHours, 40),
      overtimeHours: Math.max(response.totalHours - 40, 0),
      status: response.status as 'draft' | 'submitted' | 'approved' | 'paid',
      createdAt: response.createdAt,
      updatedAt: response.updatedAt,
    };
  };

  clearError = () => {
    this.error = null;
  };

  // Initialize store
  initialize = async (): Promise<void> => {
    const weekRange = getCurrentWeekRange();
    await Promise.all([
      this.fetchCurrentWeekTimecard(),
      this.fetchTimeEntries({
        startDate: weekRange.start,
        endDate: weekRange.end
      }),
    ]);
  };
}

export default new TimecardStore();
