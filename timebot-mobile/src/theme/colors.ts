export const colors = {
  // Primary Dark Theme Colors
  background: {
    primary: '#0a0a0a',      // Deep black background
    secondary: '#1a1a1a',    // Slightly lighter black for cards
    tertiary: '#2a2a2a',     // Medium gray for elevated surfaces
    overlay: 'rgba(0, 0, 0, 0.8)', // Semi-transparent overlay
  },

  // Purple Accent Colors (from mockup)
  purple: {
    primary: '#8b5cf6',      // Main purple accent
    light: '#a78bfa',        // Lighter purple for highlights
    dark: '#7c3aed',         // Darker purple for pressed states
    gradient: ['#8b5cf6', '#a855f7'], // Purple gradient
  },

  // Text Colors
  text: {
    primary: '#ffffff',      // White primary text
    secondary: '#d1d5db',    // Light gray secondary text
    tertiary: '#9ca3af',     // Medium gray tertiary text
    disabled: '#6b7280',     // Disabled text
    accent: '#8b5cf6',       // Purple accent text
  },

  // Status Colors
  status: {
    success: '#10b981',      // Green for success/approved
    warning: '#f59e0b',      // Orange for pending
    error: '#ef4444',        // Red for errors/rejected
    info: '#3b82f6',         // Blue for info
  },

  // Border Colors
  border: {
    primary: '#374151',      // Main border color
    secondary: '#4b5563',    // Lighter border
    accent: '#8b5cf6',       // Purple accent border
    light: '#6b7280',        // Light border
  },

  // Button Colors
  button: {
    primary: '#8b5cf6',      // Purple primary button
    primaryPressed: '#7c3aed', // Pressed state
    secondary: '#374151',    // Secondary button
    secondaryPressed: '#4b5563', // Secondary pressed
    ghost: 'transparent',    // Ghost button
    disabled: '#4b5563',     // Disabled button
  },

  // Input Colors
  input: {
    background: '#1f2937',   // Dark input background
    border: '#374151',       // Input border
    borderFocus: '#8b5cf6',  // Focused input border
    placeholder: '#9ca3af',  // Placeholder text
  },

  // Card Colors
  card: {
    background: '#1f2937',   // Card background
    elevated: '#374151',     // Elevated card
    border: '#374151',       // Card border
  },

  // Tab Bar Colors
  tabBar: {
    background: '#111827',   // Tab bar background
    active: '#8b5cf6',       // Active tab
    inactive: '#6b7280',     // Inactive tab
    border: '#374151',       // Tab bar border
  },

  // Gradient Definitions
  gradients: {
    primary: ['#0a0a0a', '#1a1a1a'],     // Main background gradient
    purple: ['#8b5cf6', '#a855f7'],       // Purple accent gradient
    card: ['#1f2937', '#374151'],         // Card gradient
    header: ['#111827', '#1f2937'],       // Header gradient
    success: ['#10b981', '#059669'],      // Success gradient
    warning: ['#f59e0b', '#d97706'],      // Warning gradient
    error: ['#ef4444', '#dc2626'],        // Error gradient
  },

  // Opacity Variants
  opacity: {
    light: 0.1,
    medium: 0.3,
    heavy: 0.6,
    overlay: 0.8,
  },
};

// Color utility functions
export const getColorWithOpacity = (color: string, opacity: number): string => {
  return `${color}${Math.round(opacity * 255).toString(16).padStart(2, '0')}`;
};

export const rgba = (color: string, alpha: number): string => {
  // Convert hex to rgba if needed
  if (color.startsWith('#')) {
    const hex = color.replace('#', '');
    const r = parseInt(hex.substr(0, 2), 16);
    const g = parseInt(hex.substr(2, 2), 16);
    const b = parseInt(hex.substr(4, 2), 16);
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
  }
  return color;
};
