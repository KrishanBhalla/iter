import { FC } from 'react';
import './Message.scss';

interface MessageProps {
  msg: string
}

const Message: FC<MessageProps> = ({msg}: MessageProps) => (
  
  <div className="Message">
    {msg}
  </div>
);

export default Message;
