import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import AnimatedText from './AnimatedText';
import Navbar from './Navbar';
import { Message, useInterviewStore } from '../store';
import { AnswerChat, EndChat } from '../js/wailsjs/go/main/App';

interface ChatScreenProps {
  setError: (error: string | null) => void;
}

const ChatScreen: React.FC<ChatScreenProps> = ({ setError }) => {
  const { messages, initialText, initialAudio, isIntroDone, interviewId, interviewSecret, hasEnded, addMessage, setIsIntroDone, setHasEnded, resetStore } = useInterviewStore();
  const [isRecording, setIsRecording] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const [hasStarted, setHasStarted] = useState(false);

  const navigate = useNavigate();

  const audioRef = useRef<HTMLAudioElement | null>(null);
  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const chatContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!initialText) {
      navigate('/');
    }
  }, [navigate, initialText]);

  useEffect(() => {
    if (!isIntroDone && initialAudio && initialAudio !== 'undefined') {
      playAudio(initialAudio);
      setIsIntroDone(true);
    }
  }, [isIntroDone, initialAudio, setIsIntroDone]);

  useEffect(() => {
    if (chatContainerRef.current) {
      chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
    }
  }, [messages]);

  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      mediaRecorderRef.current = new MediaRecorder(stream);

      const audioChunks: BlobPart[] = [];
      mediaRecorderRef.current.ondataavailable = (event) => {
        audioChunks.push(event.data);
      };

      mediaRecorderRef.current.onstop = () => {
        const audioBlob = new Blob(audioChunks, { type: 'audio/wav' });
        sendAudioToServer(audioBlob);
      };

      mediaRecorderRef.current.start();
      setIsRecording(true);
    } catch (error) {
      console.error('Error accessing microphone:', error);
      setError('Failed to access microphone. Please check your permissions and try again.');
    }
  };

  const stopRecording = () => {
    if (mediaRecorderRef.current && isRecording) {
      mediaRecorderRef.current.stop();
      setIsRecording(false);
    }
  };

  const sendAudioToServer = async (audioBlob: Blob) => {
    setIsProcessing(true);

    try {
      const audioArray = new Uint8Array(await audioBlob.arrayBuffer());
      const response = await AnswerChat(interviewId, interviewSecret, Array.from(audioArray));

      const userMessage: Message = { text: response?.prompt?.text ?? '', isUser: true };
      const botMessage: Message = { text: response?.answer?.text ?? '', isUser: false, isAnimated: true };

      addMessage(userMessage);
      addMessage(botMessage);

      if (response?.answer?.audio) {
        playAudio(response.answer.audio);
      }

      setHasStarted(true);
    } catch (error) {
      console.error('Error sending audio:', error);
      setError('Failed to send your response. Please check your connection and try again.');
    } finally {
      setIsProcessing(false);
    }
  };

  const playAudio = (base64Audio: string | null) => {
    if (!base64Audio) {
      return
    }

    // Stop any currently playing audio
    if (audioRef.current) {
      audioRef.current.pause();
      audioRef.current.currentTime = 0;
    }

    // Create a new audio element and play the audio
    audioRef.current = new Audio(`data:audio/mp3;base64,${base64Audio}`);
    audioRef.current.play();
  };

  const endInterview = async () => {
    setIsProcessing(true);

    try {
      const response = await EndChat(interviewId, interviewSecret)

      const botMessage: Message = { text: response?.answer?.text ?? '', isUser: false, isAnimated: true };
      addMessage(botMessage);

      if (response?.answer?.audio) {
        playAudio(response.answer.audio);
      }

      setHasEnded(true);

    } catch (error) {
      console.error('Error ending interview:', error);
      setError('Failed to end the interview. Please check your connection and try again.');
    } finally {
      setIsProcessing(false);
    }
  };

  const handleStartOver = () => {
    resetStore();
    navigate('/');
  };

  const handleBack = () => {
    navigate('/');
  };

  return (
    <div className="flex flex-col h-screen bg-[#1E1E2E] text-white">
      <Navbar
        showBackIcon
        showForwardIcon
        showStartOver
        onBack={handleBack}
        onStartOver={handleStartOver}
        disableForward={true}
      />
      <div ref={chatContainerRef} className="flex-grow overflow-y-auto px-4 py-2">
        {messages.map((message, index) => (
          <div key={index} className={`mb-4 ${message.isUser ? 'text-right' : 'text-left'}`}>
            <span className={`inline-block p-3 rounded-2xl ${message.isUser
              ? 'bg-[#3E64FF] text-white'
              : 'bg-[#2B2B3B] text-white'
              }`}>
              {message.isAnimated
                ? <AnimatedText message={message} />
                : message.text
              }
            </span>
          </div>
        ))}
      </div>

      <div className="flex justify-between items-center space-x-4 p-4 bg-[#1E1E2E]">
        <button
          onClick={isRecording ? stopRecording : startRecording}
          disabled={isProcessing || hasEnded}
          className={`w-full p-4 rounded-xl font-bold text-lg transition-all duration-300 ${isProcessing || hasEnded
            ? 'bg-[#2B2B3B] text-gray-400 cursor-not-allowed'
            : isRecording
              ? 'bg-[#FF3E3E] text-white animate-pulse'
              : 'bg-[#3E64FF] text-white hover:bg-opacity-90'
            }`}
        >
          {isProcessing
            ? 'Processing...'
            : isRecording
              ? 'Stop Recording'
              : hasEnded
                ? 'This interview has ended'
                : 'Start Recording'
          }
        </button>
        {hasStarted && !hasEnded && (
          <button
            onClick={endInterview}
            disabled={isProcessing || isRecording || hasEnded}
            className={`w-3/12 p-4 rounded-xl font-bold text-lg hover:bg-opacity-90 transition-all duration-300 ${isProcessing || isRecording || hasEnded
              ? 'bg-[#2B2B3B] text-gray-400 cursor-not-allowed'
              : 'bg-[#FF3E3E] text-white hover:bg-opacity-90'
              }`}
          >
            End
          </button>
        )}
      </div>
    </div>
  );
};

export default ChatScreen;