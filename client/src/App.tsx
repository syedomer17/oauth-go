import { useEffect, useState } from 'react'
import { Navigate, Route, Routes, useNavigate } from 'react-router-dom'
import './App.css'
import Dashboard from './components/Dashboard'
import LoadingScreen from './components/LoadingScreen'
import LoginPage from './components/LoginPage'
import { API_BASE_URL } from './constants/api'
import type { User } from './types/user'

const App = () => {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  useEffect(() => {
    const fetchCurrentUser = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/api/user/me`, {
          credentials: 'include',
        })

        if (!response.ok) {
          setUser(null)
          return
        }

        const data = await response.json()
        setUser(data.user)
      } catch {
        setUser(null)
      } finally {
        setLoading(false)
      }
    }

    fetchCurrentUser()
  }, [])

  const handleLogout = async () => {
    await fetch(`${API_BASE_URL}/auth/logout`, {
      method: 'POST',
      credentials: 'include',
    })
    setUser(null)
    navigate('/', { replace: true })
  }

  if (loading) {
    return <LoadingScreen />
  }

  return (
    <Routes>
      <Route
        path="/"
        element={user ? <Navigate to="/me" replace /> : <LoginPage />}
      />

      <Route
        path="/me"
        element={
          user ? (
            <Dashboard user={user} onLogout={handleLogout} />
          ) : (
            <Navigate to="/" replace />
          )
        }
      />

      <Route path="*" element={<Navigate to={user ? '/me' : '/'} replace />} />
    </Routes>
  )
}

export default App
