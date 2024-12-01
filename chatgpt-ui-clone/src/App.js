import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false); // To show a loading indicator
  const [mood, setMood] = useState('');

  const moods = [
    'mood-happy',
    'mood-calm',
    'mood-energetic',
    'mood-sad',
    'mood-sassy',
    'mood-sarcastic',
    'mood-funny'
  ];

  const moodEmojis = {
    'mood-happy': '😄', 
    'mood-calm': '🌿', 
    'mood-energetic': '⚡', 
    'mood-sad': '😢', 
    'mood-sassy': '💅', 
    'mood-sarcastic': '🙃', 
    'mood-funny': '😂',
  };
  
  useEffect(() => {
    // Select a random mood on component mount
    const randomMood = moods[Math.floor(Math.random() * moods.length)];
    setMood(randomMood);
  }, []); // Empty dependency array ensures this runs only once

  const handleSend = async (event) => {
    event.preventDefault();
    const input = event.target.elements.messageInput;
    const userMessage = input.value.trim();
    if (userMessage) {
      setMessages([...messages, { sender: 'user', text: userMessage }]);
      input.value = '';
      setLoading(true);

      try {
        const response = await axios.post('https://emogpt.onrender.com/api/sendMessage', {
          message: userMessage,
          mood: mood,
        });

        setMessages((prevMessages) => [
          ...prevMessages,
          { sender: 'bot', text: response.data.message },
        ]);
      } catch (error) {
        console.error('Error fetching the response:', error);
        setMessages((prevMessages) => [
          ...prevMessages,
          { sender: 'bot', text: 'Oops! Something went wrong. 🚧' },
        ]);
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <div className={`chat-container ${mood}`}>
      <div className="chat-header">
        EmoGPT: Your Emotional Ally 
        <span className="mood-emoji">{moodEmojis[mood]}</span> {/* Dynamic Emoji */}

      </div>
      <div className="chat-messages">
        {messages.map((msg, index) => (
          <div className={`message ${msg.sender}`} key={index}>
            <span>{msg.text}</span>
          </div>
        ))}
        {loading && (
          <div className="message bot">
            <span>Typing...</span>
          </div>
        )}
      </div>
      <form className="chat-input" onSubmit={handleSend}>
        <input
          name="messageInput"
          type="text"
          placeholder="Type your message here..."
          autoComplete="off"
        />
        <button type="submit">Send 🚀</button>
      </form>
      <div className="chat-footer">
        <p>
          Built by{' '}
          <a
            href="https://x.com/ShivamSouravJha"
            target="_blank"
            rel="noopener noreferrer"
          >
            @ShivamSouravJha
          </a>
          |   Funded by{' '}
          <a
            href="https://x.com/sonichigo1219"
            target="_blank"
            rel="noopener noreferrer"
          >
            @sonichigo1219
          </a>{' '}
        </p>
      </div>
    </div>
  );
}

export default App;
