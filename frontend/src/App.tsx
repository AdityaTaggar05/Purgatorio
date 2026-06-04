import { Routes, Route } from 'react-router-dom'

import LoginPage from "./app/routes/Login"
import RegisterPage from './app/routes/Register'
import GamePage from './app/routes/Game'

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/game" element={<GamePage />} />
    </Routes>
  )
}

export default App
