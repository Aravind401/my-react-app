import React, { useEffect, useState } from 'react';
import api from '../api';

export default function Orders() {
  const [orders, setOrders] = useState([]);
  const [form, setForm] = useState({
    product_id: '',
    quantity: '',
    order_date: ''
  });

  const load = async () => {
    try {
      const res = await api.get('/order');
      // Ensure it's always an array
      setOrders(Array.isArray(res.data) ? res.data : []);
    } catch (error) {
      console.error("Failed to fetch orders:", error);
      setOrders([]); // fallback to empty list on failure
    }
  };

  const create = async () => {
    if (!form.product_id || !form.quantity) {
      return alert("Product ID and Quantity are required");
    }

    const payload = {
      product_id: parseInt(form.product_id),
      quantity: parseInt(form.quantity),
    };

    if (form.order_date) {
      payload.order_date = new Date(form.order_date).toISOString();
    }

    try {
      await api.post('/order', payload);
      setForm({ product_id: '', quantity: '', order_date: '' });
      load();
    } catch (err) {
      console.error("Order creation failed:", err);
      alert("Failed to create order");
    }
  };

  useEffect(() => {
    load();
  }, []);

  return (
    <div className="form-container">
      <h2>Orders</h2>
      <input
        placeholder="Product ID"
        value={form.product_id}
        onChange={(e) => setForm({ ...form, product_id: e.target.value })}
      />
      <input
        placeholder="Quantity"
        value={form.quantity}
        onChange={(e) => setForm({ ...form, quantity: e.target.value })}
      />
      <input
        type="datetime-local"
        value={form.order_date}
        onChange={(e) => setForm({ ...form, order_date: e.target.value })}
      />
      <button onClick={create}>Place Order</button>

      <ul>
        {orders && orders.length > 0 ? (
          orders.map((o) => (
            <li key={o.order_id}>
              Order #{o.order_id} — Product: {o.product_id} — Qty: {o.quantity} — Date:{" "}
              {o.order_date ? new Date(o.order_date).toLocaleString() : 'N/A'}
            </li>
          ))
        ) : (
          <li>No orders found.</li>
        )}
      </ul>
    </div>
  );
}
