import ThemedButton from "@/components/ThemedButton";
import { useWizard } from "@/components/wizard";
import React from "react";
import { View } from "react-native";

interface FormButtonsProps {
  styles: {
    startButton: object;
    startButtonText: object;
  };
  onComplete?: () => void;
}

export const FormButtons = (props: FormButtonsProps) => {
  const { styles, onComplete } = props;
  const { nextStep, currentStep } = useWizard();

  const handlePress = () => {
    if (currentStep === 0) {
      nextStep();
    } else {
      onComplete?.();
    }
  };

  const renderText = (step: number) => {
    switch (step) {
      case 0:
        return "Start Chatting";
      default:
        return "I'm Ready!";
    }
  };

  return (
    <View style={{ paddingHorizontal: 20 }}>
      <ThemedButton
        text={renderText(currentStep)}
        onPress={handlePress}
        styles={{
          touchableOpacity: styles.startButton,
          text: styles.startButtonText,
        }}
      />
    </View>
  );
};
