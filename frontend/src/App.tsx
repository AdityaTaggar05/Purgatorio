import { Routes, Route } from 'react-router-dom'

import LoginPage from "./app/routes/Login"
import RegisterPage from './app/routes/Register'
import GamePage from './app/routes/Game'
import SplashScreen from './app/routes/Splash'

function App() {
  return (
    <Routes>
      <Route path="/" element={<SplashScreen />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/game" element={<GamePage />} />
    </Routes>
  )
}

export default App
