
import React, { useEffect, useState } from "react"
import Header from './components/Header';
import ChatInput from './components/ChatInput';
import ChatHistory from './components/ChatHistory';
import { connect, sendMsg } from "./api"
import logo from './logo.svg';
import './App.css';

function App() {
  const [lastMessage, setLastMessage] = useState<MessageEvent<string>>()

  useEffect(() => {
    const handleNewMessage = (msg: MessageEvent<string>) => {
      setLastMessage(() => msg);
    };
    connect(handleNewMessage);
    return () => {};
  }, []); // Empty dependency array ensures the effect runs only once on mount

  const send = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter') {
      const target = event.target as HTMLInputElement
      sendMsg(target.value)
      setLastMessage(() => new MessageEvent("message", {data: target.value, origin: "user"}))
      target.value = ""
      event.target = target
    }
  };

  return (
    <div className="App">
    <Header />
    <ChatHistory lastMessage={lastMessage}/>
    <ChatInput send={send}/>
  </div>
  );
}

export default App;
