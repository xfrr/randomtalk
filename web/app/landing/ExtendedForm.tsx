import React, { useMemo } from "react";

import { CountryPicker } from "@/components/CountryPicker";
import { ThemedTextInput } from "@/components/ThemedTextInput";
import "@expo/match-media";
import { View } from "react-native";
import { baseStyles } from "./styles";
import { ThemedPicker } from "@/components/ThemedPicker";
import Genders, { Gender } from "@/utils/genders";
import { Picker } from "@react-native-picker/picker";
import { GenderPicker } from "@/components/GenderPicker";
import { BirthdatePicker } from "@/components/BirthdatePicker";

export interface ExtendedFormState {
  birthdate: Date;
  gender: Gender;
}

export interface ExtendedFormActions {
  setBirthdate: (birthdate: Date) => void;
  setGender: (gender: Gender) => void;
}

export interface ExtendedFormProps {
  state: ExtendedFormState;
  actions: ExtendedFormActions;
  styles: {
    input: object;
    formContainer: object;
  };
}

export default function ExtendedForm(props: ExtendedFormProps) {
  const { state, actions, styles } = props;
  const { birthdate, gender } = state;
  const { setBirthdate, setGender } = actions;

  return (
    <View style={styles.formContainer}>
      <BirthdatePicker
        required
        selectedDate={birthdate}
        onChange={(birthdate, age) => {
          setBirthdate(birthdate);
          console.log("User age:", age);
        }}
        minimumDate={new Date(1920, 0, 1)}
        maximumDate={
          new Date(new Date().setFullYear(new Date().getFullYear() - 18))
        }
        minAge={18}
        maxAge={100}
        label="birthdate.label"
        placeholder="birthdate.placeholder"
        styles={{
          input: styles.input,
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
          },
          ageText: {
            color: "#4caf50",
          },
        }}
      />

      <GenderPicker
        type="rounded"
        selectedValue={gender}
        onValueChange={(val) => setGender(val as Gender)}
        containerStyle={{
          marginVertical: 0,
          shadowColor: "#000",
          shadowOffset: { width: 0, height: 2 },
          shadowOpacity: 0.15,
          shadowRadius: 3,
          elevation: 3,
        }}
      />
    </View>
  );
}
