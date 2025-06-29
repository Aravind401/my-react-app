import React from 'react';
import { Link, useNavigate } from 'react-router-dom';

export default function Header() {
  const navigate = useNavigate();

  const logout = () => {
    localStorage.removeItem('token');
    navigate('/');
  };

  return (
    <header className="header">
      <h1>Shop Manager</h1>
      <nav>
        <Link to="/products">Products</Link>
        <Link to="/orders">Orders</Link>
        <button onClick={logout}>Logout</button>
      </nav>
    </header>
  );
}
