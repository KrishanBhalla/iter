
import React, { useEffect, useState } from "react"
import Header from './components/Header';
import ChatInput from './components/ChatInput';
import ChatHistory from './components/ChatHistory';
import { connect, sendMsg } from "./api"
import logo from './logo.svg';
import './App.css';

function App() {
  const [chatHistory, setChatHistory] = useState<string[]>([])

  useEffect(() => {
    const handleNewMessage = (msg: MessageEvent<string>) => {
      setChatHistory((prevChatHistory) => [...prevChatHistory, msg.data]);
    };

    connect(handleNewMessage);

    // Clean up the WebSocket connection on component unmount
    return () => {
      // Close the WebSocket connection or perform cleanup if needed
    };
  }, []); // Empty dependency array ensures the effect runs only once on mount

  const send = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter') {
      const target = event.target as HTMLInputElement
      sendMsg(target.value)
      target.value = ""
      event.target = target
    }
  };

  return (
    <div className="App">
    <Header />
    <ChatHistory chatHistory={chatHistory}/>
    <ChatInput send={send}/>
  </div>
  );
}

export default App;
