import { FC } from 'react';
import './Message.scss';


interface MessageProps {
  msg: string
  role: string
  isVisible: boolean
}

const Message: FC<MessageProps> = ({msg, role, isVisible}: MessageProps) => {
  if (role === "system" || role === "user") {
    return <div className={(isVisible ? ("Message-" + role) : "Hidden")}>{msg}</div>
  }
  return <div></div>
};

export default Message;
