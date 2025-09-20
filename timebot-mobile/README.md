# PCC TimeBot Mobile App

A React Native mobile application for time tracking and payroll management.

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── ui/             # Basic UI components (buttons, inputs, etc.)
│   └── forms/          # Form-specific components
├── screens/            # Screen components
│   ├── auth/           # Authentication screens
│   ├── main/           # Main app screens
│   └── admin/          # Admin-only screens
├── navigation/         # Navigation configuration
├── stores/             # MobX state management stores
├── services/           # API services and external integrations
├── types/              # TypeScript type definitions
├── utils/              # Utility functions
├── styles/             # Global styles and themes
├── constants/          # App constants and configuration
└── index.ts           # Main export file
```

## Coding Standards

### File Naming Conventions
- **Components**: PascalCase (e.g., `TimeEntryCard.tsx`)
- **Screens**: PascalCase (e.g., `LoginScreen.tsx`)
- **Utilities**: camelCase (e.g., `dateHelpers.ts`)
- **Types**: camelCase (e.g., `userTypes.ts`)
- **Stores**: PascalCase with "Store" suffix (e.g., `AuthStore.ts`)

### Import/Export Patterns
- Use named exports for utilities and types
- Use default exports for components and screens
- Create index files for easy imports
- Group imports: external packages, internal modules, relative imports

### Component Structure
```typescript
// External imports
import React from 'react';
import { View, Text, StyleSheet } from 'react-native';

// Internal imports
import { COLORS, SPACING } from '../constants';
import { formatDate } from '../utils';

// Types
interface Props {
  title: string;
  onPress: () => void;
}

// Component
const MyComponent: React.FC<Props> = ({ title, onPress }) => {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>{title}</Text>
    </View>
  );
};

// Styles
const styles = StyleSheet.create({
  container: {
    padding: SPACING.md,
    backgroundColor: COLORS.surface,
  },
  title: {
    color: COLORS.text,
    fontSize: 16,
  },
});

export default MyComponent;
```

### MobX Store Pattern
```typescript
import { makeAutoObservable } from 'mobx';
import { makePersistable } from 'mobx-persist-store';

class ExampleStore {
  // Observables
  loading = false;
  data: any[] = [];

  constructor() {
    makeAutoObservable(this);
    makePersistable(this, {
      name: 'ExampleStore',
      properties: ['data'],
      storage: AsyncStorage,
    });
  }

  // Actions
  setLoading = (loading: boolean) => {
    this.loading = loading;
  };

  // Computed
  get hasData() {
    return this.data.length > 0;
  }
}

export default new ExampleStore();
```

## Design System

The app follows the dark theme design from the mockup:

### Colors
- **Primary**: #4f46e5 (Indigo)
- **Secondary**: #7c3aed (Purple)
- **Background**: #111827 (Dark gray)
- **Surface**: #1f2937 (Gray)
- **Text**: #f9fafb (Light gray)

### Typography
- Uses system fonts: SF Pro (iOS), Roboto (Android)
- Font sizes: 12px (xs) to 32px (xxxl)

### Spacing
- Consistent spacing scale: 4px, 8px, 16px, 24px, 32px, 48px

### Components
- All components follow the design system
- Consistent border radius, shadows, and animations
- Dark theme optimized for readability

## State Management

### MobX Stores
- **AuthStore**: Authentication state and user data
- **TimecardStore**: Time entries and timecard data
- **UserStore**: User preferences and settings
- **AppStore**: Global app state (loading, errors, etc.)

### Data Flow
1. UI components observe store data
2. User actions trigger store methods
3. Store methods call API services
4. API responses update store state
5. UI automatically re-renders

## Development Guidelines

### Testing
- Write unit tests for utilities and stores
- Write integration tests for complex flows
- Test components with React Native Testing Library

### Performance
- Use FlatList for large data sets
- Implement proper memoization
- Optimize images and assets
- Monitor bundle size

### Accessibility
- Add accessibility labels and hints
- Support screen readers
- Ensure proper color contrast
- Test with accessibility tools

### Error Handling
- Graceful error handling in all API calls
- User-friendly error messages
- Proper loading states
- Offline support where possible
