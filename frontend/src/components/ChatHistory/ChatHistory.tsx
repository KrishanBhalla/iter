import { FC, useState, useEffect } from 'react';
import './ChatHistory.scss';
import Message from "../Message/Message"

interface ChatHistoryProps {
  lastMessage?: MessageEvent<string>
}

interface MessageHistory {
  data: string,
  origin: string
}

const ChatHistory: FC<ChatHistoryProps> = ({lastMessage}) => {

  const [messageHistory, setMessageHistory] = useState<MessageHistory[]>([])

  useEffect(() => {
      if (lastMessage === undefined) {
        return
      }
      setMessageHistory((prevHistory) => {
      let msgData = lastMessage.data.replaceAll(RegExp("\"", "g"), "")
      let msg = {data: msgData, origin: lastMessage.origin}
      if (prevHistory.length === 0) {
        return [msg]
      }
      if (prevHistory[prevHistory.length-1].origin === lastMessage.origin) {
        let prevMessage = prevHistory.pop() || {data: "", origin: ""}
        prevMessage.data += msgData
        return [...prevHistory, prevMessage]
      }
      return [...prevHistory, msg]
    })
  }, [lastMessage])

  const messages = messageHistory.map((msg, index) => (
    <Message key={index} msg={msg.data}></Message>
  ))
  return <div className="ChatHistory">
    {messages}
  </div>
}

export default ChatHistory;
