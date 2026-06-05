import { useAuth } from "../../hooks/useAuth";
import * as api from "../../features/auth/authApi"

export default function GamePage() {
  const { user, logout } = useAuth()

  return <div className="flex items-center justify-center flex-col h-screen bg-[color-purgatory-bg]">
    <h1> Logged in as {user?.username} </h1>
    <span>
      <button onClick={async () => {
        const ok = await api.logout()

        if (ok) logout()
      }}
        type="submit"
        className="w-full font-serif bg-purgatory-border hover:bg-amber-950/40 border border-amber-500/30 hover:border-amber-500/60 text-gray-200 px-16 py-3 rounded text-sm tracking-widest font-semibold transition-all duration-300 mt-4 shadow-lg hover:cursor-pointer"
      >
        Log Out
      </button>
    </span>
  </div>
}
