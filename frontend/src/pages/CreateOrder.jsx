import { useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { api } from '../api/client'

export default function CreateOrder() {
  const navigate = useNavigate()
  const location = useLocation()
  const user = JSON.parse(localStorage.getItem('user') || '{}')
  const [form, setForm] = useState({
    restaurant_id: location.state?.restaurantId || '',
    total_price: 2500,
    address: user.address || '',
    comment: '',
  })
  const [error, setError] = useState('')

  const submit = async (e) => {
    e.preventDefault()
    setError('')
    try {
      const res = await api.createOrder({
        user_id: user.userId,
        restaurant_id: form.restaurant_id,
        total_price: Number(form.total_price),
        address: form.address,
        comment: form.comment,
        user_email: user.email,
      })
      navigate(`/orders/${res.orderId}`)
    } catch (err) {
      setError(err.message)
    }
  }

  return (
    <div className="card">
      <h1>Create Order</h1>
      <form onSubmit={submit}>
        <label>Restaurant ID<input value={form.restaurant_id} onChange={(e) => setForm({ ...form, restaurant_id: e.target.value })} required /></label>
        <label>Total price<input type="number" value={form.total_price} onChange={(e) => setForm({ ...form, total_price: e.target.value })} required /></label>
        <label>Delivery address<input value={form.address} onChange={(e) => setForm({ ...form, address: e.target.value })} required /></label>
        <label>Comment<input value={form.comment} onChange={(e) => setForm({ ...form, comment: e.target.value })} /></label>
        {error && <p className="error">{error}</p>}
        <button type="submit">Place order</button>
      </form>
    </div>
  )
}
