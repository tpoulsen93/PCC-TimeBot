// Timecard and time entry types

export interface TimeEntry {
  id: number;
  employeeId: number;
  date: string; // YYYY-MM-DD format
  startTime: string; // HH:mm format
  endTime: string; // HH:mm format
  hours: number;
  project?: string;
  description?: string;
  status: 'pending' | 'approved' | 'rejected';
  createdAt: string;
  updatedAt: string;
}

export interface TimeCard {
  id: number;
  employeeId: number;
  employee: {
    firstName: string;
    lastName: string;
    email: string;
  };
  startDate: string; // YYYY-MM-DD format
  endDate: string; // YYYY-MM-DD format
  entries: TimeEntry[];
  totalHours: number;
  regularHours: number;
  overtimeHours: number;
  status: 'draft' | 'submitted' | 'approved' | 'paid';
  submittedAt?: string;
  approvedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateTimeEntryData {
  date: string;
  startTime: string;
  endTime: string;
  project?: string;
  description?: string;
}

export interface UpdateTimeEntryData extends Partial<CreateTimeEntryData> {
  id: number;
}

export interface WeeklyTimecard {
  weekOf: string; // Start date of the week
  entries: TimeEntry[];
  totalHours: number;
  regularHours: number;
  overtimeHours: number;
  days: {
    [date: string]: {
      entries: TimeEntry[];
      totalHours: number;
    };
  };
}

export interface PayrollPeriod {
  startDate: string;
  endDate: string;
  employees: {
    employeeId: number;
    timecard: TimeCard;
  }[];
  totalHours: number;
  totalRegularHours: number;
  totalOvertimeHours: number;
}
