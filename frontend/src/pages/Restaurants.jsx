import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../api/client'

export default function Restaurants() {
  const [restaurants, setRestaurants] = useState([])
  const [error, setError] = useState('')

  useEffect(() => {
    api.listRestaurants()
      .then((res) => setRestaurants(res.restaurants || []))
      .catch((err) => setError(err.message))
  }, [])

  return (
    <div>
      <h1>Restaurants</h1>
      {error && <p className="error">{error}</p>}
      <div className="grid">
        {restaurants.map((r) => (
          <div key={r.id} className="card">
            <h3>{r.name}</h3>
            <p>{r.cuisineType} · {r.city}</p>
            <p>{r.description}</p>
            <Link to={`/restaurants/${r.id}/menu`}>View menu →</Link>
          </div>
        ))}
      </div>
    </div>
  )
}
