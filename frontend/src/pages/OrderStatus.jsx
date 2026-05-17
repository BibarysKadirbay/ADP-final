import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { api } from '../api/client'

export default function OrderStatus() {
  const { id } = useParams()
  const user = JSON.parse(localStorage.getItem('user') || '{}')
  const [order, setOrder] = useState(null)
  const [orders, setOrders] = useState([])
  const [error, setError] = useState('')

  useEffect(() => {
    if (id && id !== 'track') {
      api.getOrder(id).then(setOrder).catch((err) => setError(err.message))
    } else if (user.userId) {
      api.getUserOrders(user.userId).then((res) => setOrders(res.orders || [])).catch((err) => setError(err.message))
    }
  }, [id, user.userId])

  const refresh = () => {
    if (order?.orderId) api.getOrder(order.orderId).then(setOrder)
  }

  const cancel = async () => {
    if (!order?.orderId) return
    await api.cancelOrder(order.orderId)
    refresh()
  }

  if (id === 'track' || (!id && user.userId)) {
    return (
      <div>
        <h1>My Orders</h1>
        {error && <p className="error">{error}</p>}
        {orders.map((o) => (
          <div key={o.orderId} className="card">
            <p><strong>{o.orderId}</strong> — {o.status} / {o.paymentStatus}</p>
            <p>Total: {o.totalPrice}</p>
          </div>
        ))}
      </div>
    )
  }

  if (!order) return <p>Loading...</p>

  return (
    <div className="card">
      <h1>Order Status</h1>
      {error && <p className="error">{error}</p>}
      <pre>{JSON.stringify(order, null, 2)}</pre>
      <button type="button" onClick={refresh}>Refresh</button>
      <button type="button" className="secondary" onClick={cancel}>Cancel order</button>
    </div>
  )
}
