import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import type { NativeStackNavigationProp } from '@react-navigation/native-stack';

type RootStackParamList = {
  Home: undefined;
  Camera: undefined;
  Preview: {
    videoUri: string;
    videoSize: string;
    duration: number;
  };
};

type HomeScreenNavigationProp = NativeStackNavigationProp<RootStackParamList, 'Home'>;

export default function HomeScreen() {
  const navigation = useNavigation<HomeScreenNavigationProp>();

  return (
    <View style={styles.container}>
      <Text style={styles.title}>🧗 Escalade Benchmark</Text>
      <Text style={styles.subtitle}>React Native - Expo</Text>
      
      <TouchableOpacity
        style={styles.button}
        onPress={() => navigation.navigate('Camera')}
      >
        <Text style={styles.buttonText}>📹 Enregistrer une vidéo</Text>
      </TouchableOpacity>

      <View style={styles.info}>
        <Text style={styles.infoText}>Test de capture vidéo</Text>
        <Text style={styles.infoText}>Max 30 secondes</Text>
        <Text style={styles.infoText}>Qualité automatique</Text>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#1a1a1a',
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  title: {
    fontSize: 32,
    fontWeight: 'bold',
    color: '#fff',
    marginBottom: 10,
  },
  subtitle: {
    fontSize: 18,
    color: '#888',
    marginBottom: 60,
  },
  button: {
    backgroundColor: '#ff4444',
    paddingHorizontal: 40,
    paddingVertical: 20,
    borderRadius: 12,
    marginBottom: 40,
  },
  buttonText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: '600',
  },
  info: {
    alignItems: 'center',
  },
  infoText: {
    color: '#666',
    fontSize: 14,
    marginVertical: 4,
  },
});