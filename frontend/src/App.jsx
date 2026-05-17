import { Link, Route, Routes, Navigate } from 'react-router-dom'
import Register from './pages/Register'
import Login from './pages/Login'
import Restaurants from './pages/Restaurants'
import Menu from './pages/Menu'
import CreateOrder from './pages/CreateOrder'
import OrderStatus from './pages/OrderStatus'

export default function App() {
  const token = localStorage.getItem('token')
  const user = JSON.parse(localStorage.getItem('user') || 'null')

  const logout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    window.location.href = '/login'
  }

  return (
    <>
      <nav>
        <strong>Food Delivery</strong>
        <Link to="/restaurants">Restaurants</Link>
        {token ? (
          <>
            <Link to="/orders/new">New Order</Link>
            {user?.userId && <Link to={`/orders/track`}>My Orders</Link>}
            <button type="button" className="secondary" style={{ width: 'auto', marginLeft: 'auto' }} onClick={logout}>Logout</button>
          </>
        ) : (
          <>
            <Link to="/login">Login</Link>
            <Link to="/register">Register</Link>
          </>
        )}
      </nav>
      <div className="container">
        <Routes>
          <Route path="/" element={<Navigate to="/restaurants" />} />
          <Route path="/register" element={<Register />} />
          <Route path="/login" element={<Login />} />
          <Route path="/restaurants" element={<Restaurants />} />
          <Route path="/restaurants/:id/menu" element={<Menu />} />
          <Route path="/orders/new" element={<CreateOrder />} />
          <Route path="/orders/:id" element={<OrderStatus />} />
          <Route path="/orders/track" element={<OrderStatus />} />
        </Routes>
      </div>
    </>
  )
}
