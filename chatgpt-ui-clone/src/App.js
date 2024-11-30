import React, { useState } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false); // To show a loading indicator

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
        });

        setMessages(prevMessages => [
          ...prevMessages,
          { sender: 'bot', text: response.data.message },
        ]);
      } catch (error) {
        console.error('Error fetching the response:', error);
        setMessages(prevMessages => [
          ...prevMessages,
          { sender: 'bot', text: 'Oops! Something went wrong. 🚧' },
        ]);
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <div className="chat-container">
      <div className="chat-header">EmoGPT: Your Emotional Ally 🧡</div>
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
    </div>
  );
}

export default App;
