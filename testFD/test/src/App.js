import logo from './logo.svg';
import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    const eventSource = new EventSource('http://localhost:8080/startSSE');

    eventSource.onmessage = function(event) {
      console.log('New message:', event.data);
      console.log("this is the custome json: ", JSON.parse(event.data));
      setMessages(prevMessages => [...prevMessages, event.data]);
    };

    eventSource.onerror = function(error) {
      console.error('EventSource failed:', error);
      eventSource.close();
    };

    return () => {
      eventSource.close();
    };
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <div className=''>{" < SSE pipeline > "}</div>
        <div>
          <h2>Messages</h2>
          <ul>
            {messages.map((msg, index) => (
              <li key={index}>{msg}</li>
            ))}
          </ul>
        </div>
      </header>
    </div>
  );
}

export default App;


