import React, {
  useState,
  useEffect,
  useCallback,
  forwardRef,
  useRef,
} from "react";
import {
  TextInput,
  type TextInputProps,
  StyleSheet,
  ColorValue,
  View,
  Text,
  TouchableOpacity,
  Animated,
} from "react-native";
import { useThemeColor } from "@/hooks/useThemeColor";
import { useTranslation } from "react-i18next";
import debounce from "lodash.debounce";
import { CheckCircle, XCircle, Eye, EyeOff } from "lucide-react-native";
import { animationConfig } from "@/utils/animationConfig";

export type TextInputType = "default" | "underlined" | "rounded";

export interface ThemedTextInputProps extends TextInputProps {
  required?: boolean;
  lightColor?: string;
  darkColor?: string;
  type?: TextInputType;
  errorMessage?: string;
  successMessage?: string;
  minLength?: number;
  pattern?: RegExp;
  showPasswordToggle?: boolean;
  helperText?: string;
  asyncValidator?: (value: string) => Promise<string | null>;
}

// ForwardRef for Formik and other external refs
export const ThemedTextInput = forwardRef<TextInput, ThemedTextInputProps>(
  (
    {
      style,
      lightColor,
      darkColor,
      type = "default",
      required = false,
      errorMessage,
      successMessage,
      minLength,
      pattern,
      showPasswordToggle = false,
      helperText,
      asyncValidator,
      onBlur,
      onFocus,
      onChangeText,
      value = "",
      secureTextEntry,
      ...rest
    },
    ref
  ): JSX.Element => {
    const { t } = useTranslation();

    const textColor = useThemeColor(
      { light: lightColor, dark: darkColor },
      "text"
    );
    const placeholderColor = useThemeColor(
      { light: "#AAA", dark: "#666" },
      "placeholderText"
    );
    const backgroundColor = useThemeColor({}, "background");

    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<boolean>(false);
    const [isTouched, setIsTouched] = useState<boolean>(false);
    const [isPasswordVisible, setIsPasswordVisible] = useState<boolean>(false);
    const fadeAnim = useRef(new Animated.Value(0)).current;
    const shakeAnim = useRef(new Animated.Value(0)).current;

    // Debounced sync + async validation
    const validate = useCallback(
      debounce(async (text: string) => {
        let validationError: string | null = null;

        if (required && text.trim() === "") {
          validationError = errorMessage || t("validation.required");
        } else if (minLength && text.trim().length < minLength) {
          validationError = t("validation.minLength", { count: minLength });
        } else if (pattern && !pattern.test(text)) {
          validationError = t("validation.invalidFormat");
        }

        if (!validationError && asyncValidator) {
          try {
            const asyncError = await asyncValidator(text.trim());
            if (asyncError) validationError = asyncError;
          } catch (e) {
            validationError = t("validation.asyncError");
          }
        }

        if (validationError) {
          setError(validationError);
          setSuccess(false);
          triggerShake();
        } else {
          setError(null);
          setSuccess(text.trim() !== "");
        }
      }, 300),
      [required, errorMessage, minLength, pattern, asyncValidator, t]
    );

    useEffect(() => {
      if (isTouched) {
        validate(value);
      }
      return () => validate.cancel();
    }, [value, validate, isTouched]);

    useEffect(() => {
      Animated.timing(fadeAnim, {
        toValue: error || success ? 1 : 0,
        duration: 200,
        ...animationConfig,
      }).start();
    }, [error, success, fadeAnim]);

    const triggerShake = () => {
      Animated.sequence([
        Animated.timing(shakeAnim, {
          toValue: 5,
          duration: 50,
          ...animationConfig,
        }),
        Animated.timing(shakeAnim, {
          toValue: -5,
          duration: 50,
          ...animationConfig,
        }),
        Animated.timing(shakeAnim, {
          toValue: 3,
          duration: 50,
          ...animationConfig,
        }),
        Animated.timing(shakeAnim, {
          toValue: 0,
          duration: 50,
          ...animationConfig,
        }),
      ]).start();
    };

    const handleBlur = (e: any) => {
      validate.flush();
      onBlur?.(e);
    };

    const handleFocus = (e: any) => {
      setIsTouched(true);
      onFocus?.(e);
    };

    const handleChangeText = (text: string) => {
      onChangeText?.(text);
    };

    const togglePasswordVisibility = () => {
      setIsPasswordVisible(!isPasswordVisible);
    };

    const showError = isTouched && error;
    const showSuccess = isTouched && success;

    return (
      <View style={{ marginBottom: showError ? 8 : 15 }}>
        <Animated.View
          style={[
            styles.inputContainer,
            {
              transform: [{ translateX: shakeAnim }],
            },
          ]}
        >
          <TextInput
            ref={ref}
            style={[
              {
                flex: 1,
                color: textColor as string,
                backgroundColor: backgroundColor as string,
              },
              type === "default" && styles.default,
              type === "underlined" && styles.underlined,
              type === "rounded" && styles.rounded,
              showError && styles.errorBorder,
              showSuccess && styles.successBorder,
              style,
            ]}
            placeholderTextColor={placeholderColor as ColorValue}
            onBlur={handleBlur}
            onFocus={handleFocus}
            onChangeText={handleChangeText}
            value={value}
            secureTextEntry={
              showPasswordToggle ? !isPasswordVisible : secureTextEntry
            }
            {...rest}
          />

          {/* Error / Success Icon */}
          {showError && <XCircle size={20} stroke="red" style={styles.icon} />}
          {showSuccess && (
            <CheckCircle size={20} stroke="green" style={styles.icon} />
          )}

          {/* Password Visibility Toggle */}
          {showPasswordToggle && (
            <TouchableOpacity
              style={[
                styles.icon,
                { right: showError || showSuccess ? 30 : 10 },
              ]}
              onPress={togglePasswordVisibility}
              activeOpacity={0.7}
            >
              {isPasswordVisible ? (
                <EyeOff size={20} stroke="#666" />
              ) : (
                <Eye size={20} stroke="#666" />
              )}
            </TouchableOpacity>
          )}
        </Animated.View>

        {/* Helper text (always visible) */}
        {helperText && <Text style={styles.helperText}>{helperText}</Text>}

        {/* Animated error or success message */}
        <Animated.View style={{ opacity: fadeAnim }}>
          {showError && <Text style={styles.errorText}>{error}</Text>}
          {showSuccess && successMessage && (
            <Text style={styles.successText}>{successMessage}</Text>
          )}
        </Animated.View>
      </View>
    );
  }
);

ThemedTextInput.displayName = "ThemedTextInput";

const styles = StyleSheet.create({
  inputContainer: {
    position: "relative",
    flexDirection: "row",
    alignItems: "center",
  },
  default: {
    backgroundColor: "#3E3760",
    color: "#fff",
    paddingHorizontal: 15,
    borderRadius: 10,
    paddingRight: 40, // space for icon
  },
  underlined: {
    fontSize: 16,
    lineHeight: 24,
    borderBottomWidth: 1,
    borderBottomColor: "#999",
    paddingVertical: 6,
    paddingRight: 40,
  },
  rounded: {
    fontSize: 16,
    lineHeight: 24,
    borderRadius: 20,
    paddingVertical: 8,
    paddingHorizontal: 14,
    paddingRight: 40,
  },
  errorBorder: {
    borderColor: "red",
    borderWidth: 1,
  },
  successBorder: {
    borderColor: "green",
    borderWidth: 1,
  },
  icon: {
    position: "absolute",
    right: 10,
  },
  errorText: {
    color: "red",
    marginTop: 4,
    fontSize: 12,
  },
  successText: {
    color: "green",
    marginTop: 4,
    fontSize: 12,
  },
  helperText: {
    color: "#888",
    marginTop: 4,
    fontSize: 12,
  },
});
