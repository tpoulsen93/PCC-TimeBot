export interface Employee {
  id: number;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  supervisorId: number | null;
  isAdmin: boolean;
}

export interface HistoryEntry {
  date: string;
  hours: number;
  location: string;
  message: string;
}

export interface HistoryResponse {
  entries: HistoryEntry[];
}

export interface SummaryDay {
  date: string;
  hours: number;
  location: string;
}

export interface SummaryResponse {
  name: string;
  weekStart: string;
  weekEnd: string;
  days: SummaryDay[];
  totalHours: number;
  payday: string;
}

export interface SubmitResponse {
  message: string;
  hours: number;
  date: string;
}

export interface AdminTimecard {
  employeeId: number;
  name: string;
  totalHours: number;
}

export interface AdminTimecardsResponse {
  start: string;
  end: string;
  timecards: AdminTimecard[];
  totalHours: number;
  cost: number;
  payday?: string;
}
