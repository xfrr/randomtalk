import React, { useEffect, useState } from "react";
import { View, Text, Pressable } from "react-native";
import { DatePickerModal } from "react-native-paper-dates";
import { en, registerTranslation } from "react-native-paper-dates";
import { useTranslation } from "react-i18next";
import { ThemedText } from "./ThemedText";

// Register default locale
registerTranslation("en", en);

export type BirthdatePickerProps = {
  selectedDate?: Date;
  onChange: (date: Date, age: number) => void;
  minimumDate?: Date;
  maximumDate?: Date;
  label?: string;
  placeholder?: string;
  required?: boolean;
  minAge?: number;
  maxAge?: number;
  styles?: {
    input?: object;
    inputText?: object;
    label?: object;
    errorText?: object;
    ageText?: object;
  };
};

export function BirthdatePicker({
  selectedDate,
  onChange,
  minimumDate,
  maximumDate = new Date(),
  label = "birthdate.label",
  placeholder = "birthdate.placeholder",
  required = false,
  minAge,
  maxAge,
  styles = {
    input: {},
    inputText: {
      fontSize: 16,
      color: "#fff",
    },
    label: {
      fontWeight: "500",
      marginBottom: 5,
    },
    errorText: {
      color: "#ff4d4f",
      marginTop: 4,
    },
    ageText: {
      marginTop: 8,
      color: "#4caf50",
    },
  },
}: BirthdatePickerProps) {
  const { t } = useTranslation();

  const [open, setOpen] = useState(false);
  const [date, setDate] = useState<Date | undefined>(selectedDate);
  const [error, setError] = useState<string | null>(null);

  // Run validation on initial selectedDate
  useEffect(() => {
    if (selectedDate) {
      validateDate(selectedDate);
    }
  }, [selectedDate]);

  const handleConfirm = (params: { date: Date | undefined }) => {
    setOpen(false);

    if (!params.date) return;

    const selected = params.date;
    setDate(selected);

    if (!validateDate(selected)) return;

    const age = calculateAge(selected);
    onChange(selected, age);
  };

  const validateDate = (selected: Date): boolean => {
    const age = calculateAge(selected);

    if (required && !selected) {
      setError(t("birthdate.error.required"));
      return false;
    }

    if (minAge && age < minAge) {
      setError(t("birthdate.error.tooYoung", { minAge }));
      return false;
    }

    if (maxAge && age > maxAge) {
      setError(t("birthdate.error.tooOld", { maxAge }));
      return false;
    }

    setError(null);
    return true;
  };

  const formattedDate = date ? date.toLocaleDateString() : t(placeholder);

  return (
    <View>
      {label && <ThemedText style={styles.label}>{t(label)}</ThemedText>}

      <Pressable
        onPress={() => setOpen(true)}
        accessibilityLabel={t("birthdate.accessibilityLabel")}
        style={[
          {
            padding: 12,
            borderWidth: 1,
            borderColor: error ? "#ff4d4f" : "#ccc",
            borderRadius: 8,
            backgroundColor: "#fff",
            marginBottom: 8,
          },
          styles.input,
        ]}
      >
        <Text style={styles.inputText}>{formattedDate}</Text>
      </Pressable>

      <DatePickerModal
        locale={t("birthdate.locale") || "en"}
        mode="single"
        visible={open}
        onDismiss={() => setOpen(false)}
        date={date || new Date(2000, 0, 1)}
        onConfirm={handleConfirm}
        validRange={{
          startDate: minimumDate,
          endDate: maximumDate,
        }}
      />

      {date && !error && (
        <Text style={styles.ageText}>
          🎉 {t("birthdate.ageDisplay", { count: calculateAge(date) })}
        </Text>
      )}

      {error && <Text style={styles.errorText}>{error}</Text>}
    </View>
  );
}

// Utility function to calculate age
function calculateAge(birthdate: Date): number {
  const today = new Date();
  let age = today.getFullYear() - birthdate.getFullYear();
  const monthDiff = today.getMonth() - birthdate.getMonth();

  if (
    monthDiff < 0 ||
    (monthDiff === 0 && today.getDate() < birthdate.getDate())
  ) {
    age--;
  }

  return age;
}
