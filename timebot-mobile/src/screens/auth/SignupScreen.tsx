import React, { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  Alert,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
} from 'react-native';
import { LinearGradient } from 'expo-linear-gradient';
import { Ionicons } from '@expo/vector-icons';
import { observer } from 'mobx-react-lite';
import type { StackNavigationProp } from '@react-navigation/stack';
import type { AuthStackParamList } from '../../types/navigation';
import AuthStore from '../../stores/AuthStore';
import { useThemeStyles } from '../../theme/ThemeContext';

type SignupScreenNavigationProp = StackNavigationProp<AuthStackParamList, 'Signup'>;

interface Props {
  navigation: SignupScreenNavigationProp;
}

export const SignupScreen: React.FC<Props> = observer(({ navigation }) => {
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const theme = useThemeStyles();

  const handleSignup = async () => {
    if (!firstName.trim() || !lastName.trim() || !email.trim() || !password.trim() || !confirmPassword.trim()) {
      Alert.alert('Error', 'Please fill in all fields');
      return;
    }

    if (password !== confirmPassword) {
      Alert.alert('Error', 'Passwords do not match');
      return;
    }

    if (password.length < 6) {
      Alert.alert('Error', 'Password must be at least 6 characters long');
      return;
    }

    setIsLoading(true);
    try {
      await AuthStore.signup({
        firstName: firstName.trim(),
        lastName: lastName.trim(),
        email: email.trim(),
        password,
      });
    } catch (error) {
      Alert.alert('Signup Failed', 'Unable to create account. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <LinearGradient
      colors={theme.colors.gradients.primary as [string, string]}
      style={styles.container}
    >
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.keyboardAvoid}
      >
        <ScrollView
          contentContainerStyle={styles.scrollContent}
          keyboardShouldPersistTaps="handled"
        >
          {/* Header */}
          <View style={styles.header}>
            <Text style={[theme.styles.text.h1, styles.title]}>
              Create Account
            </Text>
            <Text style={[theme.styles.text.bodySecondary, styles.subtitle]}>
              Join PCC TimeBot today
            </Text>
          </View>

          {/* Signup Form */}
          <View style={styles.form}>
            {/* Name Inputs */}
            <View style={styles.nameRow}>
              <View style={[styles.inputContainer, styles.nameInput]}>
                <Text style={theme.styles.text.label}>First Name</Text>
                <View style={[theme.styles.input.base, styles.inputWrapper]}>
                  <Ionicons
                    name="person-outline"
                    size={20}
                    color={theme.colors.text.tertiary}
                    style={styles.inputIcon}
                  />
                  <TextInput
                    style={[theme.styles.input.base, styles.textInput]}
                    value={firstName}
                    onChangeText={setFirstName}
                    placeholder="First name"
                    placeholderTextColor={theme.colors.input.placeholder}
                    autoCapitalize="words"
                    autoCorrect={false}
                  />
                </View>
              </View>

              <View style={[styles.inputContainer, styles.nameInput]}>
                <Text style={theme.styles.text.label}>Last Name</Text>
                <View style={[theme.styles.input.base, styles.inputWrapper]}>
                  <Ionicons
                    name="person-outline"
                    size={20}
                    color={theme.colors.text.tertiary}
                    style={styles.inputIcon}
                  />
                  <TextInput
                    style={[theme.styles.input.base, styles.textInput]}
                    value={lastName}
                    onChangeText={setLastName}
                    placeholder="Last name"
                    placeholderTextColor={theme.colors.input.placeholder}
                    autoCapitalize="words"
                    autoCorrect={false}
                  />
                </View>
              </View>
            </View>

            {/* Email Input */}
            <View style={styles.inputContainer}>
              <Text style={theme.styles.text.label}>Email</Text>
              <View style={[theme.styles.input.base, styles.inputWrapper]}>
                <Ionicons
                  name="mail-outline"
                  size={20}
                  color={theme.colors.text.tertiary}
                  style={styles.inputIcon}
                />
                <TextInput
                  style={[theme.styles.input.base, styles.textInput]}
                  value={email}
                  onChangeText={setEmail}
                  placeholder="Enter your email"
                  placeholderTextColor={theme.colors.input.placeholder}
                  keyboardType="email-address"
                  autoCapitalize="none"
                  autoCorrect={false}
                />
              </View>
            </View>

            {/* Password Input */}
            <View style={styles.inputContainer}>
              <Text style={theme.styles.text.label}>Password</Text>
              <View style={[theme.styles.input.base, styles.inputWrapper]}>
                <Ionicons
                  name="lock-closed-outline"
                  size={20}
                  color={theme.colors.text.tertiary}
                  style={styles.inputIcon}
                />
                <TextInput
                  style={[theme.styles.input.base, styles.textInput]}
                  value={password}
                  onChangeText={setPassword}
                  placeholder="Create a password"
                  placeholderTextColor={theme.colors.input.placeholder}
                  secureTextEntry={!showPassword}
                  autoCapitalize="none"
                  autoCorrect={false}
                />
                <TouchableOpacity
                  onPress={() => setShowPassword(!showPassword)}
                  style={styles.eyeIcon}
                >
                  <Ionicons
                    name={showPassword ? 'eye-off-outline' : 'eye-outline'}
                    size={20}
                    color={theme.colors.text.tertiary}
                  />
                </TouchableOpacity>
              </View>
            </View>

            {/* Confirm Password Input */}
            <View style={styles.inputContainer}>
              <Text style={theme.styles.text.label}>Confirm Password</Text>
              <View style={[theme.styles.input.base, styles.inputWrapper]}>
                <Ionicons
                  name="lock-closed-outline"
                  size={20}
                  color={theme.colors.text.tertiary}
                  style={styles.inputIcon}
                />
                <TextInput
                  style={[theme.styles.input.base, styles.textInput]}
                  value={confirmPassword}
                  onChangeText={setConfirmPassword}
                  placeholder="Confirm your password"
                  placeholderTextColor={theme.colors.input.placeholder}
                  secureTextEntry={!showConfirmPassword}
                  autoCapitalize="none"
                  autoCorrect={false}
                />
                <TouchableOpacity
                  onPress={() => setShowConfirmPassword(!showConfirmPassword)}
                  style={styles.eyeIcon}
                >
                  <Ionicons
                    name={showConfirmPassword ? 'eye-off-outline' : 'eye-outline'}
                    size={20}
                    color={theme.colors.text.tertiary}
                  />
                </TouchableOpacity>
              </View>
            </View>

            {/* Signup Button */}
            <TouchableOpacity
              style={[theme.styles.button.primary, styles.signupButton]}
              onPress={handleSignup}
              disabled={isLoading}
            >
              <LinearGradient
                colors={theme.colors.gradients.purple as [string, string]}
                style={styles.buttonGradient}
              >
                {isLoading ? (
                  <Text style={theme.styles.text.buttonText}>Creating Account...</Text>
                ) : (
                  <Text style={theme.styles.text.buttonText}>Create Account</Text>
                )}
              </LinearGradient>
            </TouchableOpacity>

            {/* Login Link */}
            <View style={styles.loginContainer}>
              <Text style={theme.styles.text.bodySecondary}>
                Already have an account?{' '}
              </Text>
              <TouchableOpacity onPress={() => navigation.navigate('Login')}>
                <Text style={[theme.styles.text.body, { color: theme.colors.purple.primary }]}>
                  Sign In
                </Text>
              </TouchableOpacity>
            </View>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    </LinearGradient>
  );
});

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  keyboardAvoid: {
    flex: 1,
  },
  scrollContent: {
    flexGrow: 1,
    justifyContent: 'center',
    padding: 24,
  },
  header: {
    alignItems: 'center',
    marginBottom: 32,
  },
  title: {
    textAlign: 'center',
    marginBottom: 8,
  },
  subtitle: {
    textAlign: 'center',
  },
  form: {
    width: '100%',
  },
  nameRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    gap: 12,
  },
  nameInput: {
    flex: 1,
  },
  inputContainer: {
    marginBottom: 20,
  },
  inputWrapper: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 0,
  },
  inputIcon: {
    marginLeft: 16,
    marginRight: 12,
  },
  textInput: {
    flex: 1,
    backgroundColor: 'transparent',
    borderWidth: 0,
    paddingHorizontal: 0,
  },
  eyeIcon: {
    padding: 8,
    marginRight: 8,
  },
  signupButton: {
    marginBottom: 24,
    marginTop: 12,
    padding: 0,
    overflow: 'hidden',
  },
  buttonGradient: {
    paddingVertical: 16,
    paddingHorizontal: 24,
    alignItems: 'center',
    justifyContent: 'center',
  },
  loginContainer: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
  },
});
