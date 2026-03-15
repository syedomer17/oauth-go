import { useState } from 'react'
import { GithubIcon, GoogleIcon } from './OAuthIcons'
import type { User } from '../types/user'

interface DashboardProps {
  user: User
  onLogout: () => void
}

const Dashboard = ({ user, onLogout }: DashboardProps) => {
  const [avatarFailed, setAvatarFailed] = useState(false)
  const showAvatar = user.avatar && !avatarFailed

  return (
    <div className="auth-page">
      <div className="auth-card dashboard">
        {showAvatar ? (
          <img
            src={user.avatar}
            alt={user.name}
            className="user-avatar"
            onError={() => setAvatarFailed(true)}
          />
        ) : (
          <div className="user-avatar user-avatar-fallback" aria-label={user.name}>
            {user.name?.charAt(0).toUpperCase() || 'U'}
          </div>
        )}
        <h1 className="auth-title">{user.name}</h1>
        {user.provider === 'github' && user.username ? (
          <p className="user-handle">@{user.username}</p>
        ) : null}
        <p className="user-email">{user.email}</p>
        <span className="user-badge">
          {user.provider === 'google' ? <GoogleIcon /> : <GithubIcon />}
          {user.provider}
        </span>
        <button onClick={onLogout} className="auth-btn logout" type="button">
          Sign out
        </button>
      </div>
    </div>
  )
}

export default Dashboard
