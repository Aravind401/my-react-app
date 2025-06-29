import React, { useState, useEffect } from 'react';
import api from './api';

export default function Orders() {
  const [orders, setOrders] = useState([]);
  const [orderForm, setOrderForm] = useState({ product_id: '', quantity: '' });

  const fetchOrders = async () => {
    const res = await api.get('/order');
    setOrders(res.data);
  };

  const createOrder = async () => {
    await api.post('/order', {
      product_id: parseInt(orderForm.product_id),
      quantity: parseInt(orderForm.quantity),
    });
    fetchOrders();
  };

  useEffect(() => { fetchOrders(); }, []);

  return (
    <div>
      <h2>Orders</h2>
      <input placeholder="Product ID" onChange={(e) => setOrderForm({ ...orderForm, product_id: e.target.value })} />
      <input placeholder="Quantity" onChange={(e) => setOrderForm({ ...orderForm, quantity: e.target.value })} />
      <button onClick={createOrder}>Place Order</button>
      <ul>
        {orders.map(o => (
          <li key={o.order_id}>
            Order #{o.order_id} - Product: {o.product_id}, Quantity: {o.quantity}, Date: {new Date(o.order_date).toLocaleString()}
          </li>
        ))}
      </ul>
    </div>
  );
}
