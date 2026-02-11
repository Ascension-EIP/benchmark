import React, { useState } from 'react';
import { View, Text, TouchableOpacity, StyleSheet, Alert } from 'react-native';
import { VideoView, useVideoPlayer } from 'expo-video';
import { useRoute, useNavigation } from '@react-navigation/native';
import type { RouteProp } from '@react-navigation/native';

type RootStackParamList = {
  Home: undefined;
  Camera: undefined;
  Preview: {
    videoUri: string;
    videoSize: string;
    duration: number;
  };
};

type PreviewScreenRouteProp = RouteProp<RootStackParamList, 'Preview'>;

export default function PreviewScreen() {
  const route = useRoute<PreviewScreenRouteProp>();
  const navigation = useNavigation();
  const { videoUri, videoSize, duration } = route.params;

  const [isUploading, setIsUploading] = useState(false);

  // Nouveau player expo-video
  const player = useVideoPlayer(videoUri, (player) => {
    player.loop = true;
    player.play();
  });

  const simulateUpload = async () => {
    setIsUploading(true);
    
    try {
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      console.log('📤 Upload simulé:');
      console.log('URI:', videoUri);
      console.log('Taille:', videoSize, 'MB');
      
      Alert.alert(
        'Upload réussi ✅',
        `Vidéo de ${videoSize} MB uploadée avec succès!`,
        [
          {
            text: 'OK',
            onPress: () => navigation.navigate('Home' as never),
          },
        ]
      );
    } catch (error) {
      console.error('Erreur upload:', error);
      Alert.alert('Erreur', 'Échec de l\'upload');
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <View style={styles.container}>
      <VideoView
        style={styles.video}
        player={player}
        allowsFullscreen
        allowsPictureInPicture
      />

      <View style={styles.infoContainer}>
        <View style={styles.infoRow}>
          <Text style={styles.infoLabel}>Durée :</Text>
          <Text style={styles.infoValue}>{duration}s</Text>
        </View>
        <View style={styles.infoRow}>
          <Text style={styles.infoLabel}>Taille :</Text>
          <Text style={styles.infoValue}>{videoSize} MB</Text>
        </View>
        <View style={styles.infoRow}>
          <Text style={styles.infoLabel}>URI :</Text>
          <Text style={styles.infoValue} numberOfLines={1}>
            {videoUri.split('/').pop()}
          </Text>
        </View>
      </View>

      <View style={styles.buttonContainer}>
        <TouchableOpacity
          style={[styles.button, styles.uploadButton]}
          onPress={simulateUpload}
          disabled={isUploading}
        >
          <Text style={styles.buttonText}>
            {isUploading ? '⏳ Upload en cours...' : '📤 Simuler upload'}
          </Text>
        </TouchableOpacity>

        <TouchableOpacity
          style={[styles.button, styles.retryButton]}
          onPress={() => navigation.goBack()}
        >
          <Text style={styles.buttonText}>🔄 Réenregistrer</Text>
        </TouchableOpacity>

        <TouchableOpacity
          style={[styles.button, styles.homeButton]}
          onPress={() => navigation.navigate('Home' as never)}
        >
          <Text style={styles.buttonText}>🏠 Accueil</Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000',
  },
  video: {
    width: '100%',
    height: 300,
    backgroundColor: '#1a1a1a',
  },
  infoContainer: {
    backgroundColor: '#1a1a1a',
    padding: 20,
    margin: 20,
    borderRadius: 12,
  },
  infoRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginVertical: 8,
  },
  infoLabel: {
    color: '#888',
    fontSize: 16,
  },
  infoValue: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
    maxWidth: '60%',
  },
  buttonContainer: {
    padding: 20,
  },
  button: {
    paddingVertical: 16,
    borderRadius: 12,
    marginVertical: 8,
    alignItems: 'center',
  },
  uploadButton: {
    backgroundColor: '#4CAF50',
  },
  retryButton: {
    backgroundColor: '#ff4444',
  },
  homeButton: {
    backgroundColor: '#333',
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});