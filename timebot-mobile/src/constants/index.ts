// App constants
export const APP_NAME = 'PCC TimeBot';
export const VERSION = '1.0.0';

// API Configuration
export const API_BASE_URL = 'http://localhost:8080/api/v1';

// Colors from the mockup design
export const COLORS = {
  // Primary brand colors
  primary: '#4f46e5',
  primaryDark: '#7c3aed',
  accent: '#8b5cf6',

  // Background colors
  background: '#111827',
  surface: '#1f2937',
  card: '#1f2937',

  // Text colors
  text: '#f9fafb',
  textSecondary: '#9ca3af',
  textMuted: '#6b7280',

  // Status colors
  success: '#10b981',
  warning: '#f59e0b',
  error: '#ef4444',
  info: '#3b82f6',

  // Border colors
  border: '#374151',
  divider: '#374151',

  // Gradient colors
  gradientStart: '#1a1a2e',
  gradientEnd: '#16213e',
  buttonGradientStart: '#4f46e5',
  buttonGradientEnd: '#7c3aed',
} as const;

// Screen dimensions and spacing
export const SPACING = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  xxl: 48,
} as const;

export const BORDER_RADIUS = {
  sm: 6,
  md: 12,
  lg: 16,
  xl: 24,
} as const;

// Font sizes
export const FONT_SIZES = {
  xs: 12,
  sm: 14,
  md: 16,
  lg: 18,
  xl: 20,
  xxl: 24,
  xxxl: 32,
} as const;

// Navigation routes
export const ROUTES = {
  // Auth routes
  LOGIN: 'Login',
  SIGNUP: 'Signup',
  FORGOT_PASSWORD: 'ForgotPassword',

  // Main app routes
  TAB_NAVIGATOR: 'TabNavigator',
  DASHBOARD: 'Dashboard',
  SUBMIT_TIME: 'SubmitTime',
  HISTORY: 'History',
  TIMECARDS: 'Timecards',

  // Admin routes
  ADMIN_DASHBOARD: 'AdminDashboard',
  EMPLOYEE_MANAGEMENT: 'EmployeeManagement',
  REPORTS: 'Reports',
  SETTINGS: 'Settings',
} as const;

// Time format constants
export const TIME_FORMATS = {
  TIME_12H: 'h:mm A',
  TIME_24H: 'HH:mm',
  DATE: 'YYYY-MM-DD',
  DATE_DISPLAY: 'MMM DD, YYYY',
  DATETIME: 'YYYY-MM-DD HH:mm:ss',
} as const;

// Validation rules
export const VALIDATION = {
  MIN_PASSWORD_LENGTH: 8,
  MAX_HOURS_PER_DAY: 24,
  MIN_HOURS_PER_DAY: 0,
} as const;
