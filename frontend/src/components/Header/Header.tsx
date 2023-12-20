import { FC } from 'react';
import './Header.scss';

interface HeaderProps {}

const Header: FC<HeaderProps> = () => (
  <div className="header">
      <h2>Iter</h2>
  </div>
);

export default Header;