import React, { useMemo } from "react";

import { CountryPicker } from "@/components/CountryPicker";
import { ThemedTextInput } from "@/components/ThemedTextInput";
import "@expo/match-media";
import { View } from "react-native";
import { baseStyles } from "./styles";

export interface MainFormState {
  username: string;
  interests: string;
  country: string;
}

export interface MainFormActions {
  setUsername: (username: string) => void;
  setInterests: (interests: string) => void;
  setCountry: (country: string) => void;
}

export interface MainFormProps {
  state: MainFormState;
  actions: MainFormActions;
  styles: {
    input: object;
  };
}

export default function MainForm(props: MainFormProps) {
  const { state, actions } = props;
  const { username, interests, country } = state;
  const { setUsername, setInterests, setCountry } = actions;

  // merge the base styles with the responsive overrides
  const styles = useMemo(() => {
    return {
      ...baseStyles,
      input: [baseStyles.input, props.styles.input],
    };
  }, [props.styles.input]);

  return (
    <View style={styles.formContainer}>
      <ThemedTextInput
        required
        minLength={4}
        helperText="Username must be at least 4 characters"
        successMessage="Looks good!"
        style={styles.input}
        placeholder="Insert your username"
        placeholderTextColor="#AAA"
        onChangeText={setUsername}
        value={username}
      />

      <ThemedTextInput
        style={styles.input}
        placeholder="Insert your interests"
        placeholderTextColor="#AAA"
        helperText="Separate interests with commas,(e.g., music, coding)"
        successMessage="Looks good!"
        onChangeText={setInterests}
        value={interests}
      />

      <CountryPicker
        selectedValue={country}
        onValueChange={(val) => setCountry(val)}
        type="rounded"
        containerStyle={{
          marginTop: 10,
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
