import React, { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  Alert,
  ScrollView,
  Platform,
} from 'react-native';
import { LinearGradient } from 'expo-linear-gradient';
import { Ionicons } from '@expo/vector-icons';
import { observer } from 'mobx-react-lite';
import DateTimePicker from '@react-native-community/datetimepicker';
import { format, differenceInMinutes } from 'date-fns';

import { useThemeStyles } from '../../theme/ThemeContext';
import TimecardStore from '../../stores/TimecardStore';

export const SubmitTimeScreen: React.FC = observer(() => {
  const [selectedDate, setSelectedDate] = useState(new Date());
  const [startTime, setStartTime] = useState(new Date());
  const [endTime, setEndTime] = useState(new Date());
  const [showDatePicker, setShowDatePicker] = useState(false);
  const [showStartTimePicker, setShowStartTimePicker] = useState(false);
  const [showEndTimePicker, setShowEndTimePicker] = useState(false);
  const [notes, setNotes] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const theme = useThemeStyles();

  // Calculate total hours
  const calculateHours = (): string => {
    const diffMinutes = differenceInMinutes(endTime, startTime);
    if (diffMinutes <= 0) return '0.00';
    const hours = diffMinutes / 60;
    return hours.toFixed(2);
  };

  const handleSubmit = async () => {
    if (!selectedDate || !startTime || !endTime) {
      Alert.alert('Error', 'Please fill in all required fields');
      return;
    }

    if (endTime <= startTime) {
      Alert.alert('Error', 'End time must be after start time');
      return;
    }

    setIsSubmitting(true);
    try {
      await TimecardStore.createTimeEntry({
        date: format(selectedDate, 'yyyy-MM-dd'),
        startTime: format(startTime, 'HH:mm'),
        endTime: format(endTime, 'HH:mm'),
        description: notes.trim(),
      });

      Alert.alert(
        'Success',
        `Time entry submitted successfully!\nHours: ${calculateHours()}`,
        [{ text: 'OK', onPress: () => resetForm() }]
      );
    } catch (error) {
      Alert.alert('Error', 'Failed to submit time entry. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const resetForm = () => {
    setSelectedDate(new Date());
    setStartTime(new Date());
    setEndTime(new Date());
    setNotes('');
  };

  const onDateChange = (event: any, date?: Date) => {
    setShowDatePicker(false);
    if (date) {
      setSelectedDate(date);
    }
  };

  const onStartTimeChange = (event: any, time?: Date) => {
    setShowStartTimePicker(false);
    if (time) {
      setStartTime(time);
    }
  };

  const onEndTimeChange = (event: any, time?: Date) => {
    setShowEndTimePicker(false);
    if (time) {
      setEndTime(time);
    }
  };

  return (
    <LinearGradient
      colors={theme.colors.gradients.primary as [string, string]}
      style={styles.container}
    >
      <ScrollView
        style={styles.scrollView}
        contentContainerStyle={styles.scrollContent}
        keyboardShouldPersistTaps="handled"
      >
        {/* Header */}
        <View style={styles.header}>
          <Ionicons
            name="time-outline"
            size={32}
            color={theme.colors.purple.primary}
          />
          <Text style={[theme.styles.text.h2, styles.title]}>
            Submit Hours
          </Text>
          <Text style={[theme.styles.text.bodySecondary, styles.subtitle]}>
            Log your work hours for today
          </Text>
        </View>

        {/* Form Card */}
        <View style={[theme.styles.card.base, styles.formCard]}>
          {/* Date Picker */}
          <View style={styles.inputSection}>
            <Text style={theme.styles.text.label}>Date</Text>
            <TouchableOpacity
              style={[theme.styles.input.base, styles.dateInput]}
              onPress={() => setShowDatePicker(true)}
            >
              <Ionicons
                name="calendar-outline"
                size={20}
                color={theme.colors.text.tertiary}
                style={styles.inputIcon}
              />
              <Text style={[theme.styles.text.body, styles.dateText]}>
                {format(selectedDate, 'EEEE, MMMM d, yyyy')}
              </Text>
              <Ionicons
                name="chevron-down"
                size={20}
                color={theme.colors.text.tertiary}
              />
            </TouchableOpacity>
          </View>

          {/* Time Inputs */}
          <View style={styles.timeRow}>
            {/* Start Time */}
            <View style={[styles.inputSection, styles.timeInput]}>
              <Text style={theme.styles.text.label}>Start Time</Text>
              <TouchableOpacity
                style={[theme.styles.input.base, styles.timeButton]}
                onPress={() => setShowStartTimePicker(true)}
              >
                <Ionicons
                  name="play-outline"
                  size={20}
                  color={theme.colors.status.success}
                  style={styles.inputIcon}
                />
                <Text style={[theme.styles.text.body, styles.timeText]}>
                  {format(startTime, 'h:mm a')}
                </Text>
              </TouchableOpacity>
            </View>

            {/* End Time */}
            <View style={[styles.inputSection, styles.timeInput]}>
              <Text style={theme.styles.text.label}>End Time</Text>
              <TouchableOpacity
                style={[theme.styles.input.base, styles.timeButton]}
                onPress={() => setShowEndTimePicker(true)}
              >
                <Ionicons
                  name="stop-outline"
                  size={20}
                  color={theme.colors.status.error}
                  style={styles.inputIcon}
                />
                <Text style={[theme.styles.text.body, styles.timeText]}>
                  {format(endTime, 'h:mm a')}
                </Text>
              </TouchableOpacity>
            </View>
          </View>

          {/* Total Hours Display */}
          <View style={[styles.totalHoursCard, { backgroundColor: theme.colors.purple.primary + '20' }]}>
            <Ionicons
              name="time"
              size={24}
              color={theme.colors.purple.primary}
            />
            <View style={styles.totalHoursText}>
              <Text style={[theme.styles.text.caption, { color: theme.colors.purple.primary }]}>
                Total Hours
              </Text>
              <Text style={[theme.styles.text.h3, { color: theme.colors.purple.primary }]}>
                {calculateHours()}
              </Text>
            </View>
          </View>

          {/* Notes Input */}
          <View style={styles.inputSection}>
            <Text style={theme.styles.text.label}>Notes (Optional)</Text>
            <View style={[theme.styles.input.base, styles.notesInputWrapper]}>
              <Ionicons
                name="document-text-outline"
                size={20}
                color={theme.colors.text.tertiary}
                style={[styles.inputIcon, styles.notesIcon]}
              />
              <TextInput
                style={[theme.styles.input.base, styles.notesInput]}
                value={notes}
                onChangeText={setNotes}
                placeholder="Add notes about your work..."
                placeholderTextColor={theme.colors.input.placeholder}
                multiline
                numberOfLines={3}
                textAlignVertical="top"
              />
            </View>
          </View>

          {/* Submit Button */}
          <TouchableOpacity
            style={[theme.styles.button.primary, styles.submitButton]}
            onPress={handleSubmit}
            disabled={isSubmitting}
          >
            <LinearGradient
              colors={theme.colors.gradients.purple as [string, string]}
              style={styles.buttonGradient}
            >
              <Ionicons
                name={isSubmitting ? "hourglass-outline" : "checkmark"}
                size={20}
                color="white"
                style={styles.buttonIcon}
              />
              <Text style={theme.styles.text.buttonText}>
                {isSubmitting ? 'Submitting...' : 'Submit Hours'}
              </Text>
            </LinearGradient>
          </TouchableOpacity>
        </View>

        {/* Date Picker Modals */}
        {showDatePicker && (
          <DateTimePicker
            value={selectedDate}
            mode="date"
            display={Platform.OS === 'ios' ? 'spinner' : 'default'}
            onChange={onDateChange}
            maximumDate={new Date()}
          />
        )}

        {showStartTimePicker && (
          <DateTimePicker
            value={startTime}
            mode="time"
            display={Platform.OS === 'ios' ? 'spinner' : 'default'}
            onChange={onStartTimeChange}
          />
        )}

        {showEndTimePicker && (
          <DateTimePicker
            value={endTime}
            mode="time"
            display={Platform.OS === 'ios' ? 'spinner' : 'default'}
            onChange={onEndTimeChange}
          />
        )}
      </ScrollView>
    </LinearGradient>
  );
});

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  scrollView: {
    flex: 1,
  },
  scrollContent: {
    flexGrow: 1,
    padding: 20,
  },
  header: {
    alignItems: 'center',
    marginBottom: 32,
    marginTop: 20,
  },
  title: {
    marginTop: 12,
    marginBottom: 8,
    textAlign: 'center',
  },
  subtitle: {
    textAlign: 'center',
  },
  formCard: {
    padding: 24,
  },
  inputSection: {
    marginBottom: 24,
  },
  dateInput: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 12,
  },
  inputIcon: {
    marginRight: 12,
  },
  dateText: {
    flex: 1,
  },
  timeRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: 16,
    marginBottom: 24,
  },
  timeInput: {
    flex: 1,
    marginBottom: 0,
  },
  timeButton: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 12,
  },
  timeText: {
    flex: 1,
  },
  totalHoursCard: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 16,
    borderRadius: 12,
    marginBottom: 24,
  },
  totalHoursText: {
    marginLeft: 12,
    flex: 1,
  },
  notesInputWrapper: {
    minHeight: 80,
    alignItems: 'flex-start',
    paddingHorizontal: 16,
    paddingVertical: 12,
  },
  notesIcon: {
    marginTop: 2,
  },
  notesInput: {
    flex: 1,
    backgroundColor: 'transparent',
    borderWidth: 0,
    paddingHorizontal: 0,
    paddingVertical: 0,
    minHeight: 60,
    width: '100%',
  },
  submitButton: {
    marginTop: 8,
    padding: 0,
    overflow: 'hidden',
  },
  buttonGradient: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 16,
    paddingHorizontal: 24,
  },
  buttonIcon: {
    marginRight: 8,
  },
});
