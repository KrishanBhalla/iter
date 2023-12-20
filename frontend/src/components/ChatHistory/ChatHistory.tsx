import { FC, useState, useEffect } from 'react';
import './ChatHistory.scss';
import Message from "../Message/Message"

interface ChatHistoryProps {
  lastMessage?: MessageEvent<string>
}

const ChatHistory: FC<ChatHistoryProps> = ({lastMessage}) => {

  const [messageHistory, setMessageHistory] = useState<string[]>([])
  const [messageOrigin, setMessageOrigin] = useState<string>("")

  useEffect(() => {
      if (lastMessage === undefined) {
        return
      }
      setMessageOrigin(() => lastMessage.origin)
      setMessageHistory((prevHistory) => {
      let msgData = lastMessage.data.replaceAll(RegExp("\"\\n$", "g"), "").replaceAll(RegExp("^\"", "g"), "")
      if (messageOrigin === lastMessage.origin) {
        let prevMessage = prevHistory.pop() || ""
        prevMessage += msgData
        return [...prevHistory, prevMessage]
      }
      return [...prevHistory, msgData]
    })
  }, [lastMessage, messageOrigin])

  const messages = messageHistory.map((msg) => (
    <Message msg={msg}></Message>
  ))
  console.log(messages)
  console.log(messageOrigin)
  return <div className="ChatHistory">
    {messages}
  </div>
}

export default ChatHistory;
