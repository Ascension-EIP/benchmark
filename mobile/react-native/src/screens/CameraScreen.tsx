import React, { useState, useRef, useEffect } from "react";
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  Alert,
  Dimensions,
} from "react-native";
import { Camera, CameraView, CameraType } from "expo-camera";
import * as Haptics from "expo-haptics";
import * as FileSystem from "expo-file-system/legacy";
import { useNavigation } from "@react-navigation/native";
import type { NativeStackNavigationProp } from "@react-navigation/native-stack";

const { width, height } = Dimensions.get("window");

type RootStackParamList = {
  Home: undefined;
  Camera: undefined;
  Preview: {
    videoUri: string;
    videoSize: string;
    duration: number;
  };
};

type CameraScreenNavigationProp = NativeStackNavigationProp<
  RootStackParamList,
  "Camera"
>;

export default function CameraScreen() {
  const [hasPermission, setHasPermission] = useState<boolean | null>(null);
  const [isRecording, setIsRecording] = useState(false);
  const [recordingTime, setRecordingTime] = useState(0);
  const [cameraType, setCameraType] = useState<CameraType>("back");

  const cameraRef = useRef<CameraView>(null);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const navigation = useNavigation<CameraScreenNavigationProp>();

  useEffect(() => {
    (async () => {
      const { status } = await Camera.requestCameraPermissionsAsync();
      setHasPermission(status === "granted");
    })();
  }, []);

  useEffect(() => {
    if (isRecording) {
      timerRef.current = setInterval(() => {
        setRecordingTime((prev) => prev + 1);
      }, 1000);
    } else {
      if (timerRef.current) {
        clearInterval(timerRef.current);
      }
      setRecordingTime(0);
    }

    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current);
      }
    };
  }, [isRecording]);

  const startRecording = async () => {
    if (!cameraRef.current || isRecording) return;

    try {
      await Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Medium);
      setIsRecording(true);

      const video = await cameraRef.current.recordAsync({
        maxDuration: 30,
      });

      if (video) {
        await handleVideoRecorded(video.uri);
      }
    } catch (error) {
      console.error("Erreur enregistrement:", error);
      Alert.alert("Erreur", "Impossible d'enregistrer la vidéo");
      setIsRecording(false);
    }
  };

  const stopRecording = async () => {
    if (!cameraRef.current || !isRecording) return;

    try {
      await Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Heavy);
      cameraRef.current.stopRecording();
      setIsRecording(false);
    } catch (error) {
      console.error("Erreur arrêt:", error);
      setIsRecording(false);
    }
  };

  const handleVideoRecorded = async (uri: string) => {
    try {
      const fileInfo = await FileSystem.getInfoAsync(uri);

      if (fileInfo.exists && "size" in fileInfo) {
        const sizeInMB = (fileInfo.size / (1024 * 1024)).toFixed(2);

        console.log("📹 Vidéo enregistrée:");
        console.log("URI:", uri);
        console.log("Taille:", sizeInMB, "MB");
        console.log("Durée:", recordingTime, "secondes");

        navigation.navigate("Preview", {
          videoUri: uri,
          videoSize: sizeInMB,
          duration: recordingTime,
        });
      }
    } catch (error) {
      console.error("Erreur récupération infos:", error);
      Alert.alert("Erreur", "Impossible de récupérer les infos vidéo");
    }
  };

  const toggleCameraType = async () => {
    await Haptics.selectionAsync();
    setCameraType((current) => (current === "back" ? "front" : "back"));
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
  };

  if (hasPermission === null) {
    return (
      <View style={styles.container}>
        <Text style={styles.message}>Vérification des permissions...</Text>
      </View>
    );
  }

  if (hasPermission === false) {
    return (
      <View style={styles.container}>
        <Text style={styles.message}>
          Accès à la caméra refusé.{"\n"}
          Activez-le dans les réglages.
        </Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <CameraView
        ref={cameraRef}
        style={styles.camera}
        facing={cameraType}
        mode="video"
      />

      <View style={styles.topOverlay}>
        {isRecording && (
          <View style={styles.recordingIndicator}>
            <View style={styles.recordingDot} />
            <Text style={styles.timerText}>{formatTime(recordingTime)}</Text>
          </View>
        )}
      </View>

      <View style={styles.bottomOverlay}>
        <TouchableOpacity
          style={styles.flipButton}
          onPress={toggleCameraType}
          disabled={isRecording}
        >
          <Text style={styles.flipText}>🔄</Text>
        </TouchableOpacity>

        <TouchableOpacity
          style={[
            styles.recordButton,
            isRecording && styles.recordButtonActive,
          ]}
          onPress={isRecording ? stopRecording : startRecording}
          activeOpacity={0.7}
        >
          <View
            style={[
              styles.recordButtonInner,
              isRecording && styles.recordButtonInnerActive,
            ]}
          />
        </TouchableOpacity>

        <View style={styles.flipButton} />
      </View>

      {!isRecording && (
        <View style={styles.instructionsOverlay}>
          <Text style={styles.instructionsText}>
            Appuyez pour enregistrer (max 30s)
          </Text>
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#000",
  },
  camera: {
    position: "absolute",
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
  },
  message: {
    flex: 1,
    textAlign: "center",
    color: "#fff",
    fontSize: 16,
    padding: 20,
    justifyContent: "center",
    alignItems: "center",
  },
  topOverlay: {
    position: "absolute",
    top: 60,
    left: 0,
    right: 0,
    alignItems: "center",
  },
  recordingIndicator: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: "rgba(255, 0, 0, 0.8)",
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 20,
  },
  recordingDot: {
    width: 12,
    height: 12,
    borderRadius: 6,
    backgroundColor: "#fff",
    marginRight: 8,
  },
  timerText: {
    color: "#fff",
    fontSize: 18,
    fontWeight: "bold",
  },
  bottomOverlay: {
    position: "absolute",
    bottom: 40,
    left: 0,
    right: 0,
    flexDirection: "row",
    justifyContent: "space-around",
    alignItems: "center",
    paddingHorizontal: 30,
  },
  flipButton: {
    width: 50,
    height: 50,
    justifyContent: "center",
    alignItems: "center",
  },
  flipText: {
    fontSize: 32,
  },
  recordButton: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: "rgba(255, 255, 255, 0.3)",
    justifyContent: "center",
    alignItems: "center",
    borderWidth: 4,
    borderColor: "#fff",
  },
  recordButtonActive: {
    backgroundColor: "rgba(255, 0, 0, 0.3)",
  },
  recordButtonInner: {
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: "#ff0000",
  },
  recordButtonInnerActive: {
    borderRadius: 8,
    width: 40,
    height: 40,
  },
  instructionsOverlay: {
    position: "absolute",
    bottom: 140,
    left: 0,
    right: 0,
    alignItems: "center",
  },
  instructionsText: {
    color: "#fff",
    fontSize: 14,
    backgroundColor: "rgba(0, 0, 0, 0.6)",
    paddingHorizontal: 20,
    paddingVertical: 10,
    borderRadius: 20,
  },
});
