// Validation utilities

/**
 * Email validation
 */
export const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

/**
 * Password validation
 */
export const isValidPassword = (password: string): boolean => {
  return password.length >= 8;
};

/**
 * Phone number validation
 */
export const isValidPhoneNumber = (phone: string): boolean => {
  const phoneRegex = /^\(?([0-9]{3})\)?[-. ]?([0-9]{3})[-. ]?([0-9]{4})$/;
  return phoneRegex.test(phone);
};

/**
 * Hours validation (0-24)
 */
export const isValidHours = (hours: number): boolean => {
  return hours >= 0 && hours <= 24;
};

/**
 * Date validation (YYYY-MM-DD format)
 */
export const isValidDate = (dateString: string): boolean => {
  const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
  if (!dateRegex.test(dateString)) return false;

  const date = new Date(dateString);
  return date instanceof Date && !isNaN(date.getTime());
};

/**
 * Required field validation
 */
export const isRequired = (value: string | number): boolean => {
  if (typeof value === 'string') {
    return value.trim().length > 0;
  }
  return value !== null && value !== undefined;
};

/**
 * Get validation error message
 */
export const getValidationError = (field: string, value: string): string | null => {
  if (!isRequired(value)) {
    return `${field} is required`;
  }

  switch (field.toLowerCase()) {
    case 'email':
      return isValidEmail(value) ? null : 'Please enter a valid email address';
    case 'password':
      return isValidPassword(value) ? null : 'Password must be at least 8 characters long';
    case 'phone':
      return isValidPhoneNumber(value) ? null : 'Please enter a valid phone number';
    default:
      return null;
  }
};
