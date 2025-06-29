import React, { useEffect, useState } from 'react';
import api from '../api';

export default function Products() {
  const [products, setProducts] = useState([]);
  const [form, setForm] = useState({ name: '', description: '', price: '', in_stock: true });

  const load = async () => {
    const res = await api.get('/product');
    setProducts(res.data);
  };

  const add = async () => {
    if (!form.name || !form.price) {
      alert('Name and Price required');
      return;
    }
    await api.post('/product', {
      ...form,
      price: parseFloat(form.price),
      in_stock: form.in_stock === 'true'
    });
    setForm({ name: '', description: '', price: '', in_stock: true });
    load();
  };

  const del = async (id) => {
    await api.delete(`/product/${id}`);
    load();
  };

  useEffect(() => {
    load();
  }, []);

  return (
    <div className="form-container">
      <h2>Products</h2>
      <input
        placeholder="Name"
        value={form.name}
        onChange={(e) => setForm({ ...form, name: e.target.value })}
      />
      <input
        placeholder="Description"
        value={form.description}
        onChange={(e) => setForm({ ...form, description: e.target.value })}
      />
      <input
        placeholder="Price"
        type="number"
        value={form.price}
        onChange={(e) => setForm({ ...form, price: e.target.value })}
      />
      <select
        value={form.in_stock}
        onChange={(e) => setForm({ ...form, in_stock: e.target.value })}
      >
        <option value="true">In Stock</option>
        <option value="false">Out of Stock</option>
      </select>
      <button onClick={add}>Add Product</button>

      <ul>
        {products.map((p) => (
          <li key={p.id}>
            <strong>{p.name}</strong> – ₹{p.price} – {p.in_stock ? 'In Stock' : 'Out of Stock'}
            <button onClick={() => del(p.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
}
