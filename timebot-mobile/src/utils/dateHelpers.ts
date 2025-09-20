// Date and time helper functions
import { format, parseISO, startOfWeek, endOfWeek, addDays, subDays } from 'date-fns';

/**
 * Format a date string to display format
 */
export const formatDate = (dateString: string, formatString: string = 'MMM dd, yyyy'): string => {
  try {
    const date = parseISO(dateString);
    return format(date, formatString);
  } catch (error) {
    console.error('Date formatting error:', error);
    return dateString;
  }
};

/**
 * Get current date in YYYY-MM-DD format
 */
export const getCurrentDate = (): string => {
  return format(new Date(), 'yyyy-MM-dd');
};

/**
 * Get current time in HH:mm format
 */
export const getCurrentTime = (): string => {
  return format(new Date(), 'HH:mm');
};

/**
 * Get the start and end dates of the current week
 */
export const getCurrentWeekRange = (): { start: string; end: string } => {
  const now = new Date();
  const start = startOfWeek(now, { weekStartsOn: 1 }); // Monday
  const end = endOfWeek(now, { weekStartsOn: 1 }); // Sunday

  return {
    start: format(start, 'yyyy-MM-dd'),
    end: format(end, 'yyyy-MM-dd'),
  };
};

/**
 * Check if a date is today
 */
export const isToday = (dateString: string): boolean => {
  const today = format(new Date(), 'yyyy-MM-dd');
  return dateString === today;
};

/**
 * Get days of the week for a given date range
 */
export const getDaysInRange = (startDate: string, endDate: string): string[] => {
  const start = parseISO(startDate);
  const end = parseISO(endDate);
  const days: string[] = [];

  let current = start;
  while (current <= end) {
    days.push(format(current, 'yyyy-MM-dd'));
    current = addDays(current, 1);
  }

  return days;
};
