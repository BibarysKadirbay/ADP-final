import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { api } from '../api/client'

export default function Menu() {
  const { id } = useParams()
  const [menu, setMenu] = useState(null)
  const [error, setError] = useState('')

  useEffect(() => {
    api.getMenu(id).then(setMenu).catch((err) => setError(err.message))
  }, [id])

  if (error) return <p className="error">{error}</p>
  if (!menu) return <p>Loading menu...</p>

  return (
    <div>
      <Link to="/restaurants">← Back</Link>
      <h1>Menu</h1>
      {(menu.categories || []).map((cat) => (
        <div key={cat.category?.id} className="card">
          <h3>{cat.category?.name}</h3>
          <ul>
            {(cat.items || []).map((item) => (
              <li key={item.id}>
                {item.name} — {item.price} ({item.isAvailable ? 'available' : 'unavailable'})
              </li>
            ))}
          </ul>
        </div>
      ))}
      <Link to="/orders/new" state={{ restaurantId: id }}>Order from this restaurant</Link>
    </div>
  )
}
