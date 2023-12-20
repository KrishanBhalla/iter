import { FC } from 'react';
import './Message.scss';


interface MessageProps {
  msg: string
  role: string
}

const Message: FC<MessageProps> = ({msg, role}: MessageProps) => (
  
  <div className={"Message-" + role}>{msg}</div>
);

export default Message;
