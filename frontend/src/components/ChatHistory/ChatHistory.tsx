import { FC } from 'react';
import './ChatHistory.scss';
import Message from "../Message/Message"

interface ChatHistoryProps {
  chatHistory: string[]
}

const ChatHistory: FC<ChatHistoryProps> = ({chatHistory}) => {
  const messages = chatHistory.map((msg) => (
    <Message msg={msg}></Message>
  ))
  return <div className="ChatHistory">
    {messages}
  </div>
}

export default ChatHistory;
