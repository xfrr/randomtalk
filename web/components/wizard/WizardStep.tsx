import { ReactNode } from "react";
import { View } from "react-native";

interface WizardStepProps {
  children: ReactNode;
}

export const WizardStep = (props: WizardStepProps) => {
  const { children } = props;
  return <View>{children}</View>;
};

export default WizardStep;
