import { StyleSheet } from "react-native";

export default StyleSheet.create({
  /* The gradient behind everything */
  gradientContainer: {
    flex: 1,
  },

  container: {
    flex: 1,
  },

  // Header
  headerContainer: {
    paddingTop: 60,
    paddingBottom: 25,
    paddingHorizontal: 20,
    backgroundColor: "transparent",
    alignItems: "flex-start",
    // Shadow/elevation for depth
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 6,
    elevation: 6,
    marginBottom: 10,
  },
  title: {
    color: "#fff",
    fontSize: 28,
    fontWeight: "700",
    letterSpacing: 0.8,
    marginBottom: 5,
    textShadowColor: "rgba(0, 0, 0, 0.2)",
    textShadowOffset: { width: 1, height: 2 },
    textShadowRadius: 3,
  },
  subTitle: {
    color: "#BFBFBF",
    fontSize: 14,
    fontWeight: "400",
  },

  // Chat
  chatContainer: {
    flex: 1,
    paddingHorizontal: 14,
  },
  messageContainer: {
    marginVertical: 5,
    padding: 12,
    borderRadius: 14,
    maxWidth: "80%",
  },
  myMessage: {
    alignSelf: "flex-end",
    backgroundColor: "#543AB7",
    // Subtle shadow
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.25,
    shadowRadius: 3,
    elevation: 3,
  },
  otherMessage: {
    alignSelf: "flex-start",
    backgroundColor: "#332D45",
  },
  messageText: {
    color: "#fff",
    fontSize: 16,
    lineHeight: 20,
  },
  timestamp: {
    color: "#dcdcdc",
    fontSize: 10,
    marginTop: 4,
    textAlign: "right",
  },
  inputContainer: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: "#2B2540",
    borderTopWidth: 1,
    borderTopColor: "#444",
    paddingHorizontal: 10,
    paddingVertical: 10,
  },
  input: {
    flex: 1,
    color: "#fff",
    backgroundColor: "#3E3760",
    borderRadius: 10,
    paddingHorizontal: 12,
    paddingVertical: 8,
    marginRight: 10,
  },
  sendButton: {
    backgroundColor: "#543AB7",
    borderRadius: 10,
    paddingHorizontal: 16,
    paddingVertical: 10,
    // Another subtle shadow
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.25,
    shadowRadius: 3,
    elevation: 3,
  },
  sendButtonText: {
    color: "#fff",
    fontWeight: "600",
    fontSize: 16,
  },
});
