import React from 'react';
import { Platform, View } from 'react-native';
import { StatusBar as ExpoStatusBar } from 'expo-status-bar';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { colors } from '../../theme';

interface StatusBarProps {
  backgroundColor?: string;
  barStyle?: 'light-content' | 'dark-content';
}

export const StatusBar: React.FC<StatusBarProps> = ({
  backgroundColor = colors.background.primary,
  barStyle = 'light-content',
}) => {
  const insets = useSafeAreaInsets();

  return (
    <>
      <ExpoStatusBar
        style={barStyle === 'light-content' ? 'light' : 'dark'}
        backgroundColor={backgroundColor}
      />
      {Platform.OS === 'ios' && (
        <View
          style={{
            height: insets.top,
            backgroundColor,
          }}
        />
      )}
    </>
  );
};
