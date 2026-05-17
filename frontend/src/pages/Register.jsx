import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api } from '../api/client'

export default function Register() {
  const navigate = useNavigate()
  const [form, setForm] = useState({ name: '', email: '', password: '', phone: '', address: '' })
  const [error, setError] = useState('')

  const submit = async (e) => {
    e.preventDefault()
    setError('')
    try {
      const res = await api.register(form)
      localStorage.setItem('token', res.token)
      localStorage.setItem('user', JSON.stringify(res.user))
      navigate('/restaurants')
    } catch (err) {
      setError(err.message)
    }
  }

  return (
    <div className="card">
      <h1>Register</h1>
      <form onSubmit={submit}>
        {['name', 'email', 'password', 'phone', 'address'].map((field) => (
          <label key={field}>
            {field}
            <input
              type={field === 'password' ? 'password' : 'text'}
              value={form[field]}
              onChange={(e) => setForm({ ...form, [field]: e.target.value })}
              required={field !== 'phone'}
            />
          </label>
        ))}
        {error && <p className="error">{error}</p>}
        <button type="submit">Create account</button>
      </form>
      <p>Already have an account? <Link to="/login">Login</Link></p>
    </div>
  )
}
