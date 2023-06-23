import React from 'react';
import { Toolbar } from 'primereact/toolbar';
import { Link } from 'react-router-dom';

const NavBar = () => {
  const startElements = (
    <nav>
      <ul>
        <li>
          <Link to="/">Home</Link>
        </li>
        <li>
          <Link to="/applications">Applications</Link>
        </li>
      </ul>
    </nav>
  );
  const endElements = <></>;
  return (
    <div>
      <Toolbar start={startElements} end={endElements} />
    </div>
  );
};

export default NavBar;
