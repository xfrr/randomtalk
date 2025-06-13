import { ReactNode } from "react";
import { View } from "react-native";

interface WizardFooterProps {
  children: ReactNode;
}

export const WizardFooter = (props: WizardFooterProps) => {
  const { children } = props;
  return <View>{children}</View>;
};

export default WizardFooter;
