import { FC, KeyboardEvent } from 'react';
import './ChatInput.scss';


interface ChatInputProps {
  send: (event: KeyboardEvent<HTMLInputElement>) => void
  isVisible: boolean
}

const ChatInput: FC<ChatInputProps> = ({send, isVisible}) => {
  return (
    <div className={"ChatInput"+(isVisible ? "Default":"Hidden")}>
      {isVisible ? <input onKeyDown={send} placeholder="Hit enter to send" /> : undefined}
    </div>
  )
};

export default ChatInput;
