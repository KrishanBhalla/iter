import { FC, KeyboardEvent } from 'react';
import './ChatInput.scss';

interface ChatInputProps {
  send: (event: KeyboardEvent<HTMLInputElement>) => void
}

const ChatInput: FC<ChatInputProps> = ({send}) => {
  return (
    <div className="ChatInput">
      <input onKeyDown={send} placeholder="Hit enter to send" />
    </div>
  )
};

export default ChatInput;
