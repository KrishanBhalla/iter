import { FC } from 'react';
import './Message.scss';


interface MessageProps {
  msg: string
  role: string
}

const Message: FC<MessageProps> = ({msg, role}: MessageProps) => {
  if (role === "system" || role === "user") {
    return <div className={"Message-" + role}>{msg}</div>
  }
  return <div></div>
};

export default Message;
