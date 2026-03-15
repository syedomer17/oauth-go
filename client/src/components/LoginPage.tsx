import { API_BASE_URL } from '../constants/api'
import { GithubIcon, GoogleIcon } from './OAuthIcons'

const LoginPage = () => (
  <div className="auth-page">
    <div className="auth-card">
      <div className="auth-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z"
          />
        </svg>
      </div>

      <h1 className="auth-title">Welcome back</h1>
      <p className="auth-subtitle">Sign in to your account to continue</p>

      <div className="auth-buttons">
        <a href={`${API_BASE_URL}/auth/google/login`} className="auth-btn google">
          <GoogleIcon />
          Continue with Google
        </a>

        <a href={`${API_BASE_URL}/auth/github/login`} className="auth-btn github">
          <GithubIcon />
          Continue with GitHub
        </a>
      </div>

      <p className="auth-legal">
        By signing in you agree to our <a href="#">Terms</a> and <a href="#">Privacy Policy</a>.
      </p>
    </div>
  </div>
)

export default LoginPage
