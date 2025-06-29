import React, { useState } from 'react';
import api from '../api';
import { useNavigate } from 'react-router-dom';

export default function Login() {
  const [credentials, setCredentials] = useState({ username: '', password: '' });
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const login = async () => {
    try {
      const res = await api.post('/login', credentials);
      localStorage.setItem('token', res.data.token);
      navigate('/products');
    } catch {
      setError('Invalid username or password');
    }
  };

  return (
    <div className="form-container">
      <h2>Login</h2>
      <input
        placeholder="Username"
        onChange={(e) => setCredentials({ ...credentials, username: e.target.value })}
      />
      <input
        placeholder="Password"
        type="password"
        onChange={(e) => setCredentials({ ...credentials, password: e.target.value })}
      />
      <button onClick={login}>Login</button>
      {error && <p className="error">{error}</p>}
    </div>
  );
}
