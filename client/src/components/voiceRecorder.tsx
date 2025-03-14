import React, { useState, useRef, useEffect } from 'react';
import { Box, Button, Typography, Card } from '@mui/material';

type AudioRecorder = {
  mediaRecorder: MediaRecorder | null;
  audioChunks: Blob[];
};

const VoiceRecorder: React.FC = () => {
  const [recording, setRecording] = useState<boolean>(false);
  const [time, setTime] = useState<number>(0);
  const [warning, setWarning] = useState<boolean>(false);
  const [audioRecorder, setAudioRecorder] = useState<AudioRecorder>({ mediaRecorder: null, audioChunks: [] });
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  const MAX_TIME = 300; // 5 minutes in seconds

  useEffect(() => {
    if (recording && time >= 240 && !warning) {
      setWarning(true); // Trigger the warning at 4 minutes
    }

    if (time >= MAX_TIME) {
      handleStop(); // Stop recording when max time is reached
    }
  }, [time, recording, warning]);

  useEffect(() => {
    if (recording) {
      intervalRef.current = setInterval(() => {
        setTime((prev) => prev + 1);
      }, 1000);

      return () => {
        if (intervalRef.current) clearInterval(intervalRef.current);
      };
    }
  }, [recording]);

  const drawFrequency = () => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    ctx.clearRect(0, 0, canvas.width, canvas.height);

    const randomBars = Array.from({ length: 30 }, () => Math.random() * canvas.height);
    randomBars.forEach((barHeight, index) => {
      const color = warning ? 'red' : '#1976d2'; // Change color to red when warning
      ctx.fillStyle = color;
      ctx.fillRect(index * 10, canvas.height - barHeight, 8, barHeight);
    });
  };

  useEffect(() => {
    if (recording) {
      const interval = setInterval(drawFrequency, 100);
      return () => clearInterval(interval);
    }
  }, [recording, warning]);

  const handleStart = async () => {
    setRecording(true);
    setTime(0);
    setWarning(false);

    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    const recorder = new MediaRecorder(stream);

    recorder.ondataavailable = (event) => {
      setAudioRecorder((prev) => ({ ...prev, audioChunks: [...prev.audioChunks, event.data] }));
    };

    recorder.start();
    setAudioRecorder({ mediaRecorder: recorder, audioChunks: [] });
  };

  const handleStop = () => {
    setRecording(false);
    if (intervalRef.current) clearInterval(intervalRef.current);

    if (audioRecorder.mediaRecorder) {
      audioRecorder.mediaRecorder.stop();
      audioRecorder.mediaRecorder.stream.getTracks().forEach((track) => track.stop());

      const audioBlob = new Blob(audioRecorder.audioChunks, { type: 'audio/mp3' });
      const audioUrl = URL.createObjectURL(audioBlob);

      const link = document.createElement('a');
      link.href = audioUrl;
      link.download = 'recording.mp3';
      link.click();

      // Clean up the URL
      URL.revokeObjectURL(audioUrl);
    }
  };

  const formatTime = (time: number): string => {
    const minutes = Math.floor(time / 60);
    const seconds = time % 60;
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  return (
    <Box
      display="flex"
      justifyContent="center"
      alignItems="center"
      height="100vh"
      bgcolor="#f5f5f5"
    >
      <Card sx={{ padding: 4, width: 400, textAlign: 'center' }}>
        <Typography variant="h5" gutterBottom>
          Voice Recorder
        </Typography>

        <Box
          display="flex"
          justifyContent="space-between"
          alignItems="center"
          marginBottom={2}
        >
          <Typography variant="h6">{formatTime(time)}</Typography>
          <canvas
            ref={canvasRef}
            width={300}
            height={50}
            style={{ border: '1px solid #ccc' }}
          ></canvas>
          <Button
            variant="contained"
            color="error"
            onClick={handleStop}
            disabled={!recording}
          >
            Stop
          </Button>
        </Box>

        {warning && (
          <Typography color="error" variant="body2">
            Only 1 minute left!
          </Typography>
        )}

        {!recording && (
          <Button
            variant="contained"
            onClick={handleStart}
            disabled={recording}
          >
            Start Recording
          </Button>
        )}
      </Card>
    </Box>
  );
};

export default VoiceRecorder;
