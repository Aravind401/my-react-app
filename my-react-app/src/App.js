import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Login from './components/Login';
import Products from './components/Products';
import Orders from './components/Orders';
import Header from './components/Header';
import Footer from './components/Footer';
import ProtectedRoute from './components/ProtectedRoute'; // 👈 Import this
import './styles.css';

export default function App() {
  return (
    <Router>
      <div id="app-root">
        <Header />
        <main className="main-content">
          <Routes>
            <Route path="/" element={<Login />} />
            <Route path="/products" element={
              <ProtectedRoute>
                <Products />
              </ProtectedRoute>
            } />
            <Route path="/orders" element={
              <ProtectedRoute>
                <Orders />
              </ProtectedRoute>
            } />
          </Routes>
        </main>
        <Footer />
      </div>
    </Router>
  );
}
