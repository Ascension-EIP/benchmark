import 'dart:async';
import 'dart:io';
import 'package:flutter/material.dart';
import 'package:camera/camera.dart';
import 'package:go_router/go_router.dart';
import 'package:permission_handler/permission_handler.dart';
import '../models/video_info.dart';

class CameraScreen extends StatefulWidget {
  const CameraScreen({super.key});

  @override
  State<CameraScreen> createState() => _CameraScreenState();
}

class _CameraScreenState extends State<CameraScreen> {
  CameraController? _controller;
  List<CameraDescription>? _cameras;
  bool _isRecording = false;
  int _recordingTime = 0;
  Timer? _timer;
  bool _isInitialized = false;
  String? _errorMessage;

  @override
  void initState() {
    super.initState();
    _initializeCamera();
  }

  Future<void> _initializeCamera() async {
    // Demande les permissions
    final cameraStatus = await Permission.camera.request();
    final micStatus = await Permission.microphone.request();

    if (!cameraStatus.isGranted || !micStatus.isGranted) {
      setState(() {
        _errorMessage = 'Permissions caméra/micro refusées';
      });
      return;
    }

    try {
      // Récupère les caméras disponibles
      _cameras = await availableCameras();
      
      if (_cameras == null || _cameras!.isEmpty) {
        setState(() {
          _errorMessage = 'Aucune caméra disponible';
        });
        return;
      }

      // Initialise la caméra arrière
      _controller = CameraController(
        _cameras![0],
        ResolutionPreset.high,
        enableAudio: true,
      );

      await _controller!.initialize();
      
      if (!mounted) return;
      
      setState(() {
        _isInitialized = true;
      });
    } catch (e) {
      setState(() {
        _errorMessage = 'Erreur initialisation: $e';
      });
    }
  }

  Future<void> _startRecording() async {
    if (_controller == null || !_controller!.value.isInitialized) return;

    try {
      await _controller!.startVideoRecording();
      
      setState(() {
        _isRecording = true;
        _recordingTime = 0;
      });

      // Timer
      _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
        setState(() {
          _recordingTime++;
        });

        // Stop automatique après 30 secondes
        if (_recordingTime >= 30) {
          _stopRecording();
        }
      });
    } catch (e) {
      debugPrint('Erreur enregistrement: $e');
    }
  }

  Future<void> _stopRecording() async {
    if (_controller == null || !_controller!.value.isRecordingVideo) return;

    try {
      _timer?.cancel();
      
      final video = await _controller!.stopVideoRecording();
      
      setState(() {
        _isRecording = false;
      });

      // Récupère infos fichier
      final file = File(video.path);
      final fileSize = await file.length();
      final sizeInMB = fileSize / (1024 * 1024);

      debugPrint('📹 Vidéo enregistrée:');
      debugPrint('Path: ${video.path}');
      debugPrint('Taille: ${sizeInMB.toStringAsFixed(2)} MB');
      debugPrint('Durée: $_recordingTime secondes');

      // Navigation vers PreviewScreen
      if (!mounted) return;
      
      context.push(
        '/preview',
        extra: VideoInfo(
          path: video.path,
          sizeInMB: sizeInMB,
          duration: _recordingTime,
        ),
      );
    } catch (e) {
      debugPrint('Erreur arrêt: $e');
    }
  }

  Future<void> _toggleCamera() async {
    if (_cameras == null || _cameras!.length < 2) return;

    try {
      final currentIndex = _cameras!.indexOf(_controller!.description);
      final newIndex = (currentIndex + 1) % _cameras!.length;

      await _controller?.dispose();

      _controller = CameraController(
        _cameras![newIndex],
        ResolutionPreset.high,
        enableAudio: true,
      );

      await _controller!.initialize();
      
      if (!mounted) return;
      setState(() {});
    } catch (e) {
      debugPrint('Erreur flip caméra: $e');
    }
  }

  String _formatTime(int seconds) {
    final mins = seconds ~/ 60;
    final secs = seconds % 60;
    return '${mins.toString().padLeft(2, '0')}:${secs.toString().padLeft(2, '0')}';
  }

  @override
  void dispose() {
    _timer?.cancel();
    _controller?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (_errorMessage != null) {
      return Scaffold(
        backgroundColor: Colors.black,
        body: Center(
          child: Text(
            _errorMessage!,
            style: const TextStyle(color: Colors.white, fontSize: 16),
            textAlign: TextAlign.center,
          ),
        ),
      );
    }

    if (!_isInitialized) {
      return const Scaffold(
        backgroundColor: Colors.black,
        body: Center(
          child: CircularProgressIndicator(color: Colors.white),
        ),
      );
    }

    return Scaffold(
      backgroundColor: Colors.black,
      body: Stack(
        children: [
          // Caméra
          Positioned.fill(
            child: CameraPreview(_controller!),
          ),

          // Timer overlay
          if (_isRecording)
            Positioned(
              top: 60,
              left: 0,
              right: 0,
              child: Center(
                child: Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 16,
                    vertical: 8,
                  ),
                  decoration: BoxDecoration(
                    color: Colors.red.withOpacity(0.8),
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Container(
                        width: 12,
                        height: 12,
                        decoration: const BoxDecoration(
                          color: Colors.white,
                          shape: BoxShape.circle,
                        ),
                      ),
                      const SizedBox(width: 8),
                      Text(
                        _formatTime(_recordingTime),
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),

          // Controls
          Positioned(
            bottom: 40,
            left: 0,
            right: 0,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                // Flip camera button
                IconButton(
                  onPressed: _isRecording ? null : _toggleCamera,
                  icon: const Icon(Icons.flip_camera_android),
                  color: Colors.white,
                  iconSize: 40,
                ),

                // Record button
                GestureDetector(
                  onTap: _isRecording ? _stopRecording : _startRecording,
                  child: Container(
                    width: 80,
                    height: 80,
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      color: Colors.white.withOpacity(0.3),
                      border: Border.all(color: Colors.white, width: 4),
                    ),
                    child: Center(
                      child: Container(
                        width: 60,
                        height: 60,
                        decoration: BoxDecoration(
                          color: Colors.red,
                          borderRadius: BorderRadius.circular(
                            _isRecording ? 8 : 30,
                          ),
                        ),
                      ),
                    ),
                  ),
                ),

                // Spacer
                const SizedBox(width: 40),
              ],
            ),
          ),

          // Instructions
          if (!_isRecording)
            Positioned(
              bottom: 140,
              left: 0,
              right: 0,
              child: Center(
                child: Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 20,
                    vertical: 10,
                  ),
                  decoration: BoxDecoration(
                    color: Colors.black.withOpacity(0.6),
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: const Text(
                    'Appuyez pour enregistrer (max 30s)',
                    style: TextStyle(color: Colors.white, fontSize: 14),
                  ),
                ),
              ),
            ),
        ],
      ),
    );
  }
}