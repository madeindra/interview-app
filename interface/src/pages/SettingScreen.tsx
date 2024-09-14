import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Navbar from './Navbar';
import { useInterviewStore } from '../store';
import { AreKeyExist, UpdateAPIKeys } from '../js/wailsjs/go/main/App';

interface SettingScreenProps {
    backendHost: string;
    setError: (error: string | null) => void;
}

const SettingScreen: React.FC<SettingScreenProps> = ({ backendHost, setError }) => {
    const { messages } = useInterviewStore();

    const [openaiKeyInput, setOpenaiKeyInput] = useState('');
    const [elevenlabsKeyInput, setElevenlabsKeyInput] = useState('');

    const navigate = useNavigate();

    const handleSave = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            await UpdateAPIKeys(openaiKeyInput, elevenlabsKeyInput);
            navigate('/');
        } catch (error) {
            setError('Failed to save settings. Please check your connection and try again.');
        }
    };

    const checkAPIKeys = async () => {
        try {
            const response = await AreKeyExist();
            if (response) {
                console.log("API keys found, redirecting to home")
                navigate('/');
            }
        } catch (error) {
            console.error(error);
        }
    };

    useEffect(() => {
        checkAPIKeys();
    }, []);


    return (
        <div className="flex flex-col h-screen bg-[#1E1E2E] text-white">
            {messages.length > 0 && (
                <Navbar
                    backendHost={backendHost}
                    showBackIcon
                    showForwardIcon
                    disableBack={true}
                />
            )}
            <div className="container mx-auto mt-10 p-4 flex-grow">
                <div className="max-w-md mx-auto bg-[#2B2B3B] p-8 rounded-xl shadow-lg">
                    <h1 className="text-3xl font-bold mb-6 text-center text-white">Settings</h1>
                    <form onSubmit={handleSave} className="space-y-6">
                        <div>
                            <label htmlFor="openai-key" className="block mb-2 text-white font-semibold">
                                OpenAI Key
                            </label>
                            <input
                                type="password"
                                id="openai-key"
                                value={openaiKeyInput}
                                placeholder='sk-...'
                                onChange={(e) => setOpenaiKeyInput(e.target.value)}
                                className="w-full p-3 bg-[#3A3A4E] text-white border border-[#4A4A5E] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3E64FF]"
                            />
                        </div>
                        <div>
                            <label htmlFor="elevenlabs-key" className="block mb-2 text-white font-semibold">
                                ElevenLabs Key
                            </label>
                            <input
                                type="password"
                                id="elevenlabs-key"
                                value={elevenlabsKeyInput}
                                placeholder='sk-...'
                                onChange={(e) => setElevenlabsKeyInput(e.target.value)}
                                className="w-full p-3 bg-[#3A3A4E] text-white border border-[#4A4A5E] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3E64FF]"
                            />
                        </div>
                        <button
                            type="submit"
                            className="w-full p-4 bg-[#3E64FF] text-white font-bold rounded-xl hover:bg-opacity-90 transition-all duration-300"
                        >
                            Save Settings
                        </button>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default SettingScreen;
