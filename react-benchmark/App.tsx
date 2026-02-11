import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import HomeScreen from './src/screens/HomeScreen';
import CameraScreen from './src/screens/CameraScreen';
import PreviewScreen from './src/screens/PreviewScreen';

export type RootStackParamList = {
  Home: undefined;
  Camera: undefined;
  Preview: {
    videoUri: string;
    videoSize: string;
    duration: number;
  };
};

const Stack = createNativeStackNavigator<RootStackParamList>();

export default function App() {
  console.log('🚀 APP LOADED - NAVIGATION ACTIVE');
  return (
    <NavigationContainer>
      <Stack.Navigator
        initialRouteName="Home"
        screenOptions={{
          headerStyle: { backgroundColor: '#1a1a1a' },
          headerTintColor: '#fff',
          headerTitleStyle: { fontWeight: 'bold' },
        }}
      >
        <Stack.Screen 
          name="Home" 
          component={HomeScreen}
          options={{ title: 'Escalade Benchmark' }}
        />
        <Stack.Screen 
          name="Camera" 
          component={CameraScreen}
          options={{ 
            headerShown: false,
            animation: 'fade',
          }}
        />
        <Stack.Screen 
          name="Preview" 
          component={PreviewScreen}
          options={{ title: 'Aperçu vidéo' }}
        />
      </Stack.Navigator>
    </NavigationContainer>
  );
}