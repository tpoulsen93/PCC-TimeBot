// Time calculation utilities

/**
 * Calculate hours between start and end time
 */
export const calculateHours = (startTime: string, endTime: string): number => {
  if (!startTime || !endTime) return 0;

  const [startHours, startMinutes] = startTime.split(':').map(Number);
  const [endHours, endMinutes] = endTime.split(':').map(Number);

  const startTotalMinutes = startHours * 60 + startMinutes;
  let endTotalMinutes = endHours * 60 + endMinutes;

  // Handle overnight shifts (end time is next day)
  if (endTotalMinutes < startTotalMinutes) {
    endTotalMinutes += 24 * 60; // Add 24 hours
  }

  const totalMinutes = endTotalMinutes - startTotalMinutes;
  return Math.round((totalMinutes / 60) * 100) / 100; // Round to 2 decimal places
};

/**
 * Calculate total hours from an array of time entries
 */
export const calculateTotalHours = (timeEntries: Array<{ hours: number }>): number => {
  const total = timeEntries.reduce((sum, entry) => sum + entry.hours, 0);
  return Math.round(total * 100) / 100; // Round to 2 decimal places
};

/**
 * Calculate overtime hours (anything over 40 hours per week)
 */
export const calculateOvertime = (totalHours: number, regularHoursLimit: number = 40): { regular: number; overtime: number } => {
  if (totalHours <= regularHoursLimit) {
    return {
      regular: Math.round(totalHours * 100) / 100,
      overtime: 0,
    };
  }

  return {
    regular: regularHoursLimit,
    overtime: Math.round((totalHours - regularHoursLimit) * 100) / 100,
  };
};

/**
 * Format hours to display string (e.g., "8.5 hrs")
 */
export const formatHours = (hours: number): string => {
  if (hours === 0) return '0 hrs';
  if (hours === 1) return '1 hr';
  return `${hours} hrs`;
};

/**
 * Validate time format (HH:mm)
 */
export const isValidTimeFormat = (time: string): boolean => {
  const timeRegex = /^([01]?[0-9]|2[0-3]):[0-5][0-9]$/;
  return timeRegex.test(time);
};

/**
 * Convert decimal hours to hours and minutes string
 */
export const hoursToHoursMinutes = (decimalHours: number): string => {
  const hours = Math.floor(decimalHours);
  const minutes = Math.round((decimalHours - hours) * 60);

  if (minutes === 0) {
    return `${hours}h`;
  }

  return `${hours}h ${minutes}m`;
};
