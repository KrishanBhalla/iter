import { FC, useState, useEffect } from 'react';
import './ChatHistory.scss';
import Message from "../Message/Message"
import { Message as IMessage} from "../../types"

interface ChatHistoryProps {
  lastMessage?: IMessage
}


const ChatHistory: FC<ChatHistoryProps> = ({lastMessage}) => {

  const [messageHistory, setMessageHistory] = useState<IMessage[]>([])

  useEffect(() => {
      if (lastMessage === undefined) {
        return
      }
      setMessageHistory((prevHistory) => {
      let msgData = lastMessage.data
  
      let msg = {data: msgData, origin: lastMessage.origin}
      if (prevHistory.length === 0) {
        return [msg]
      }
      if (prevHistory[prevHistory.length-1].origin === lastMessage.origin) {
        let prevMessage = prevHistory.pop() || {data: "", origin: "system"}
        prevMessage.data += msgData
        return [...prevHistory, prevMessage]
      }
      return [...prevHistory, msg]
    })
  }, [lastMessage])
  const messages = messageHistory.map((msg, index) => (
    <Message key={index} msg={msg.data} role={msg.origin} isVisible={true}/>
  ))
  return <div className="ChatHistory">
    {messages}
  </div>
}

export default ChatHistory;
