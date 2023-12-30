import { FC, KeyboardEvent } from 'react';
import './ChatInput.scss';
import {Message as MessageType} from "../../types"
import { MESSAGE_TYPE_CHAT, MESSAGE_TYPE_CONTEXT, VALID_MESSAGE_TYPE } from '../../api/websocket';


interface ChatInputProps {
  send: (event: KeyboardEvent<HTMLInputElement>, role: string, messageType: VALID_MESSAGE_TYPE) => void
  isVisible: boolean
  lastMessage?: MessageType
}

const ChatInput: FC<ChatInputProps> = ({send, isVisible, lastMessage}) => {
  if (lastMessage !== undefined && lastMessage.origin === "SET_COUNTRY") {
    return (
      <div className={"ChatInput"+(isVisible ? "Default":"Hidden")}>
        {isVisible ? <input onKeyDown={(e) => send(e, "user", MESSAGE_TYPE_CONTEXT)} placeholder="Hit enter to send"/> : undefined}
      </div>
    )
  }
  return (
    <div className={"ChatInput"+(isVisible ? "Default":"Hidden")}>
      {isVisible ? <input onKeyDown={(e) => send(e, "user", MESSAGE_TYPE_CHAT)} placeholder="Hit enter to send" /> : undefined}
    </div>
  )
};

export default ChatInput;
