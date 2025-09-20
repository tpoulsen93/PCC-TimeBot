import { colors, getColorWithOpacity, rgba } from './colors';
import {
  spacing,
  borderRadius,
  typography,
  shadows,
  animations,
  layout,
  breakpoints
} from './tokens';
import {
  buttonStyles,
  textStyles,
  cardStyles,
  inputStyles,
  containerStyles,
  badgeStyles,
  navigationStyles,
} from './styles';

// Re-export everything
export { colors, getColorWithOpacity, rgba };
export {
  spacing,
  borderRadius,
  typography,
  shadows,
  animations,
  layout,
  breakpoints
};
export {
  buttonStyles,
  textStyles,
  cardStyles,
  inputStyles,
  containerStyles,
  badgeStyles,
  navigationStyles,
};

// Main theme object
export const theme = {
  colors,
  spacing,
  borderRadius,
  typography,
  shadows,
  animations,
  layout,
  breakpoints,
  styles: {
    button: buttonStyles,
    text: textStyles,
    card: cardStyles,
    input: inputStyles,
    container: containerStyles,
    badge: badgeStyles,
    navigation: navigationStyles,
  },
};

export type Theme = typeof theme;
