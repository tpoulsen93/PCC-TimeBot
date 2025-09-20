import { TextStyle, ViewStyle } from 'react-native';
import { colors } from './colors';
import { typography, spacing, borderRadius, shadows } from './tokens';

// Button Styles
export const buttonStyles = {
  primary: {
    backgroundColor: colors.button.primary,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    borderRadius: borderRadius.lg,
    alignItems: 'center' as const,
    justifyContent: 'center' as const,
    ...shadows.md,
  },
  secondary: {
    backgroundColor: colors.button.secondary,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    borderRadius: borderRadius.lg,
    alignItems: 'center' as const,
    justifyContent: 'center' as const,
    ...shadows.sm,
  },
  ghost: {
    backgroundColor: 'transparent',
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    borderRadius: borderRadius.lg,
    borderWidth: 1,
    borderColor: colors.border.accent,
    alignItems: 'center' as const,
    justifyContent: 'center' as const,
  },
  small: {
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    borderRadius: borderRadius.md,
  },
  large: {
    paddingHorizontal: spacing.xl,
    paddingVertical: spacing.lg,
    borderRadius: borderRadius.xl,
  },
} as const;

// Text Styles
export const textStyles = {
  h1: {
    fontSize: typography.fontSize['4xl'],
    fontWeight: typography.fontWeight.bold,
    color: colors.text.primary,
    lineHeight: typography.lineHeight.tight,
    marginBottom: spacing.md,
  } as TextStyle,
  h2: {
    fontSize: typography.fontSize['3xl'],
    fontWeight: typography.fontWeight.semibold,
    color: colors.text.primary,
    lineHeight: typography.lineHeight.tight,
    marginBottom: spacing.sm,
  } as TextStyle,
  h3: {
    fontSize: typography.fontSize['2xl'],
    fontWeight: typography.fontWeight.semibold,
    color: colors.text.primary,
    lineHeight: typography.lineHeight.tight,
    marginBottom: spacing.sm,
  } as TextStyle,
  h4: {
    fontSize: typography.fontSize.xl,
    fontWeight: typography.fontWeight.medium,
    color: colors.text.primary,
    lineHeight: typography.lineHeight.normal,
    marginBottom: spacing.sm,
  } as TextStyle,
  body: {
    fontSize: typography.fontSize.base,
    fontWeight: typography.fontWeight.normal,
    color: colors.text.primary,
    lineHeight: typography.lineHeight.normal,
  } as TextStyle,
  bodySecondary: {
    fontSize: typography.fontSize.base,
    fontWeight: typography.fontWeight.normal,
    color: colors.text.secondary,
    lineHeight: typography.lineHeight.normal,
  } as TextStyle,
  caption: {
    fontSize: typography.fontSize.sm,
    fontWeight: typography.fontWeight.normal,
    color: colors.text.tertiary,
    lineHeight: typography.lineHeight.normal,
  } as TextStyle,
  label: {
    fontSize: typography.fontSize.sm,
    fontWeight: typography.fontWeight.medium,
    color: colors.text.primary,
    marginBottom: spacing.xs,
  } as TextStyle,
  buttonText: {
    fontSize: typography.fontSize.base,
    fontWeight: typography.fontWeight.semibold,
    color: colors.text.primary,
  } as TextStyle,
  tabText: {
    fontSize: typography.fontSize.xs,
    fontWeight: typography.fontWeight.medium,
  } as TextStyle,
} as const;

// Card Styles
export const cardStyles = {
  base: {
    backgroundColor: colors.card.background,
    borderRadius: borderRadius.lg,
    padding: spacing.lg,
    ...shadows.md,
  } as ViewStyle,
  elevated: {
    backgroundColor: colors.card.elevated,
    borderRadius: borderRadius.lg,
    padding: spacing.lg,
    ...shadows.lg,
  } as ViewStyle,
  bordered: {
    backgroundColor: colors.card.background,
    borderRadius: borderRadius.lg,
    padding: spacing.lg,
    borderWidth: 1,
    borderColor: colors.border.primary,
  } as ViewStyle,
  compact: {
    backgroundColor: colors.card.background,
    borderRadius: borderRadius.md,
    padding: spacing.md,
    ...shadows.sm,
  } as ViewStyle,
} as const;

// Input Styles
export const inputStyles = {
  base: {
    backgroundColor: colors.input.background,
    borderWidth: 1,
    borderColor: colors.input.border,
    borderRadius: borderRadius.md,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    fontSize: typography.fontSize.base,
    color: colors.text.primary,
  } as ViewStyle & TextStyle,
  focused: {
    borderColor: colors.input.borderFocus,
    ...shadows.purple,
  } as ViewStyle,
  error: {
    borderColor: colors.status.error,
  } as ViewStyle,
  large: {
    paddingVertical: spacing.md,
    fontSize: typography.fontSize.lg,
  } as ViewStyle & TextStyle,
} as const;

// Container Styles
export const containerStyles = {
  screen: {
    flex: 1,
    backgroundColor: colors.background.primary,
  } as ViewStyle,
  safeArea: {
    flex: 1,
    backgroundColor: colors.background.primary,
  } as ViewStyle,
  scrollContainer: {
    flexGrow: 1,
    padding: spacing.lg,
  } as ViewStyle,
  centeredContainer: {
    flex: 1,
    justifyContent: 'center' as const,
    alignItems: 'center' as const,
    padding: spacing.lg,
    backgroundColor: colors.background.primary,
  } as ViewStyle,
  cardContainer: {
    padding: spacing.lg,
    gap: spacing.md,
  } as ViewStyle,
} as const;

// Status Badge Styles
export const badgeStyles = {
  success: {
    backgroundColor: colors.status.success,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
    borderRadius: borderRadius.full,
    alignSelf: 'flex-start' as const,
  } as ViewStyle,
  warning: {
    backgroundColor: colors.status.warning,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
    borderRadius: borderRadius.full,
    alignSelf: 'flex-start' as const,
  } as ViewStyle,
  error: {
    backgroundColor: colors.status.error,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
    borderRadius: borderRadius.full,
    alignSelf: 'flex-start' as const,
  } as ViewStyle,
  info: {
    backgroundColor: colors.status.info,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
    borderRadius: borderRadius.full,
    alignSelf: 'flex-start' as const,
  } as ViewStyle,
} as const;

// Navigation Styles
export const navigationStyles = {
  header: {
    backgroundColor: colors.background.secondary,
    elevation: 0,
    shadowOpacity: 0,
    borderBottomWidth: 1,
    borderBottomColor: colors.border.primary,
  } as ViewStyle,
  tabBar: {
    backgroundColor: colors.tabBar.background,
    borderTopWidth: 1,
    borderTopColor: colors.border.primary,
    height: 80,
    paddingBottom: spacing.sm,
    paddingTop: spacing.sm,
  } as ViewStyle,
  tabItem: {
    alignItems: 'center' as const,
    justifyContent: 'center' as const,
    flex: 1,
  } as ViewStyle,
} as const;
