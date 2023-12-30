
import React, { useEffect, useState } from "react"
import Header from './components/Header';
import ChatInput from './components/ChatInput';
import ChatHistory from './components/ChatHistory';
import Message from './components/Message';
import { Websocket as WS, Countries } from "./api"
import { Message as MessageType } from './types';
import './App.css';
import CountryDropdown from "./components/CountryDropdown";


function App() {
  const [lastMessage, setLastMessage] = useState<MessageType>()
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [countries, setCountries] = useState<string[]>([])
  const [isChatVisible, setIsChatVisible] = useState<boolean>(false)

  useEffect(() => {
    async function getCountries() {
      setIsLoading(true);
      const countryList = await Countries.getCountries()
      setCountries(countryList)
      setIsLoading(false)
    }
    getCountries()
    const handleNewMessage = (msg: MessageType) => {
      setLastMessage(() => {
        return {data: JSON.parse(msg.data), origin: "system"}
      });
    };
    WS.connect(handleNewMessage);
    return () => {};
  }, []); // Empty dependency array ensures the effect runs only once on mount

  const send = (event: React.KeyboardEvent<HTMLInputElement>, role: string = "user") => {
    if (event.key === 'Enter') {
      const target = event.target as HTMLInputElement
      WS.sendMsg(target.value)
      setLastMessage(() => ({data: target.value, origin: role}))
      target.value = ""
      event.target = target
    }
    if (!isChatVisible) {
      setIsChatVisible(true)
    }
  };

  if (isLoading) {
    return <p>Loading countries...</p>;
  }
  return (
    <div className="App">
    <Header />
    <Message msg={"Where would you like to go?"} role={"system"}/>
    <CountryDropdown countries={countries} send={send} isDisabled={isChatVisible}/>
    <ChatHistory lastMessage={lastMessage}/>
    <ChatInput send={send} isVisible={isChatVisible}/>
  </div>
  );
}

export default App;
