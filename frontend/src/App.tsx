
import React, { useEffect, useState } from "react"
import Header from './components/Header';
import ChatInput from './components/ChatInput';
import ChatHistory from './components/ChatHistory';
import { connect, sendMsg } from "./api"
import { Message } from './types';
import './App.css';

function App() {
  const [lastMessage, setLastMessage] = useState<Message>()

  useEffect(() => {
    const handleNewMessage = (msg: Message) => {
      setLastMessage(() => {
        return {data: JSON.parse(msg.data), origin: "system"}
      });
    };
    connect(handleNewMessage);
    return () => {};
  }, []); // Empty dependency array ensures the effect runs only once on mount

  const send = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter') {
      const target = event.target as HTMLInputElement
      sendMsg(target.value)
      setLastMessage(() => ({data: target.value, origin: "user"}))
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
