import type { HudProps } from '../../types/hud';

export default function GameHud({
  username,
  level,
  penitence,
  grace,
  sinMeter,
  onAscensionClick,
  onLogoutClick,
  onAttackClick,
  onArmyClick
}: HudProps) {

  const getPercent = (current: number, max: number) => Math.min(100, Math.max(0, (current / max) * 100));

  return (
    <div className="absolute inset-0 pointer-events-none z-20 font-serif select-none p-6 flex flex-col justify-between">

      {/* ================= TOP ANCHOR ROW ================= */}
      <div className="flex justify-between items-start w-full">

        {/* TOP LEFT */}
        <div className="pointer-events-auto flex items-center gap-3 bg-purgatory-card/90 backdrop-blur-md border border-purgatory-border px-4 py-3 rounded shadow-xl">
          <div className="w-10 h-10 bg-purgatory-input border border-amber-600/30 flex items-center justify-center rounded">
            <span className="text-amber-500 font-bold text-sm tracking-wider">LV.{level}</span>
          </div>
          <div>
            <div className="text-gray-200 font-sans font-semibold tracking-wide text-sm">{username}</div>
            <div className="text-[10px] text-gray-500 tracking-[0.15em] uppercase">Penitent Soul</div>
          </div>
        </div>

        {/* TOP CENTER */}
        <div className="pointer-events-auto flex flex-col items-center w-80 bg-purgatory-card/90 backdrop-blur-md border border-purgatory-border p-3 rounded shadow-xl mx-4">
          <div className="flex justify-between w-full text-[10px] tracking-[0.2em] text-red-500/90 uppercase font-bold mb-1.5">
            <span>Sin Meter</span>
            <span>{sinMeter}%</span>
          </div>
          <div className="w-full h-2.5 bg-black/50 border border-purgatory-border rounded-sm overflow-hidden p-[1px]">
            <div
              className="h-full bg-gradient-to-r from-red-900 via-red-700 to-red-600 shadow-[0_0_8px_rgba(220,38,38,0.3)] transition-all duration-500 ease-out"
              style={{ width: `${sinMeter}%` }}
            />
          </div>
        </div>

        {/* TOP RIGHT */}
        <div className="pointer-events-auto flex flex-col gap-2 w-72 bg-purgatory-card/90 backdrop-blur-md border border-purgatory-border p-3 rounded shadow-xl">
          {/* Penitence Meter */}
          <div className="flex flex-col w-full">
            <div className="flex justify-between text-[10px] tracking-wider mb-0.5 font-sans">
              <span className="text-purple-400 font-serif uppercase tracking-[0.15em] text-[9px] font-bold">{penitence.label}</span>
              <span className="text-gray-400 text-xs font-medium">{penitence.current} / {penitence.max}</span>
            </div>
            <div className="w-full h-2 bg-black/50 border border-purgatory-border rounded-sm overflow-hidden">
              <div
                className="h-full bg-gradient-to-r from-purple-900 to-purple-500 transition-all duration-300"
                style={{ width: `${getPercent(penitence.current, penitence.max)}%` }}
              />
            </div>
          </div>

          {/* Grace Meter */}
          <div className="flex flex-col w-full">
            <div className="flex justify-between text-[10px] tracking-wider mb-0.5 font-sans">
              <span className="text-teal-400 font-serif uppercase tracking-[0.15em] text-[9px] font-bold">{grace.label}</span>
              <span className="text-gray-400 text-xs font-medium">{grace.current} / {grace.max}</span>
            </div>
            <div className="w-full h-2 bg-black/50 border border-purgatory-border rounded-sm overflow-hidden">
              <div
                className="h-full bg-gradient-to-r from-teal-900 to-teal-500 transition-all duration-300"
                style={{ width: `${getPercent(grace.current, grace.max)}%` }}
              />
            </div>
          </div>
        </div>

      </div>

      {/* ================= BOTTOM ANCHOR ROW ================= */}
      <div className="flex justify-between items-end w-full">

        {/* BOTTOM LEFT */}
        <div className="pointer-events-auto flex gap-3 items-center">
          {/* Attack Button */}
          <button
            onClick={onAttackClick}
            className="group flex flex-col items-center justify-center w-20 h-20 bg-purgatory-input hover:bg-red-950/40 border border-red-900/40 hover:border-red-600/60 rounded text-gray-300 hover:text-red-400 cursor-pointer transition-all shadow-lg shadow-black/50"
          >
            <svg className="w-6 h-6 transition-transform group-hover:scale-110" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M15.362 5.214A8.252 8.252 0 0 1 12 21 8.25 8.25 0 0 1 6.038 7.047 8.287 8.287 0 0 0 9 9.601a8.983 8.983 0 0 1 3.361-6.867 8.21 8.21 0 0 0 3 2.48Z" />
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 18a3.75 3.75 0 0 0 .495-7.467 5.99 5.99 0 0 0-1.925 3.546 5.974 5.974 0 0 1-2.133-1A3.75 3.75 0 0 0 12 18Z" />
            </svg>
            <span className="text-[9px] uppercase tracking-widest font-sans font-bold mt-1">Attack</span>
          </button>

          {/* Army Button */}
          <button
            onClick={onArmyClick}
            className="group flex flex-col items-center justify-center w-14 h-14 bg-purgatory-input hover:bg-teal-950/30 border border-teal-900/30 hover:border-teal-600/50 rounded text-gray-300 hover:text-teal-400 cursor-pointer transition-all shadow-lg shadow-black/50"
          >
            <svg className="w-6 h-6 transition-transform group-hover:scale-110" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M18 18.72a9.094 9.094 0 0 0 3.741-.479 3 3 0 0 0-4.65-1.12 3.614 3.614 0 0 0-5.182 0 3 3 0 0 0-4.65 1.12 9.088 9.088 0 0 0 3.931.485M15 10a3 3 0 1 1-6 0 3 3 0 0 1 6 0Zm6 3a2.25 2.25 0 1 1-4.5 0 2.25 2.25 0 0 1 4.5 0Zm-13.5 0a2.25 2.25 0 1 1-4.5 0 2.25 2.25 0 0 1 4.5 0Z" />
            </svg>
            <span className="text-[9px] uppercase tracking-widest font-sans font-bold mt-1">Legion</span>
          </button>
        </div>

        {/* BOTTOM RIGHT */}
        <div className="pointer-events-auto flex items-center gap-4">
          {/* Small Logout Button */}
          <button
            onClick={onLogoutClick}
            title="Abandon Ascent"
            className="flex items-center justify-center w-10 h-10 bg-purgatory-card border border-purgatory-border hover:border-red-900/60 text-gray-500 hover:text-red-500 rounded cursor-pointer transition-colors shadow-lg"
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15M12 9l-3 3m0 0 3 3m-3-3h12.75" />
            </svg>
          </button>

          {/* Big Shop Button (Altar of Exchange) */}
          <button
            onClick={onAscensionClick}
            className="group relative flex items-center gap-3 bg-gradient-to-b from-purgatory-input to-purgatory-card hover:from-purgatory-card hover:to-purgatory-input border border-amber-600/40 hover:border-amber-500/80 px-6 py-3.5 rounded shadow-2xl cursor-pointer transition-all duration-300"
          >
            <div className="absolute inset-0 bg-amber-500/0 group-hover:bg-amber-500/[0.02] transition-colors rounded" />

            <svg className="w-5 h-5 text-amber-500/80 group-hover:text-amber-400 transition-transform group-hover:scale-110" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 21v-7.5a.75.75 0 0 1 .75-.75h3a.75.75 0 0 1 .75.75V21m-4.5 0H2.36m11.14 0H18m0 0h3.64m-1.39 0V9.349M3.75 21V9.349m0 0a3.001 3.001 0 0 0 3.75-.615A2.993 2.993 0 0 0 9.75 9.75c.896 0 1.7-.393 2.25-1.016a2.993 2.993 0 0 0 2.25 1.016c.896 0 1.7-.393 2.25-1.015a3.001 3.001 0 0 0 3.75.614m-16.5 0a3.004 3.004 0 0 1-.621-4.72l1.189-1.19A1.5 1.5 0 0 1 5.378 3h13.243a1.5 1.5 0 0 1 1.06.44l1.19 1.189a3 3 0 0 1-.621 4.72M6.75 18h3.75a.75.75 0 0 0 .75-.75V13.5a.75.75 0 0 0-.75-.75H6.75a.75.75 0 0 0-.75.75v3.75c0 .414.336.75.75.75Z" />
            </svg>

            <div className="text-left">
              <div className="text-xs font-bold tracking-[0.2em] text-gray-200 uppercase group-hover:text-white transition-colors">Altar of Exchange</div>
              <div className="text-[9px] font-sans text-amber-600/70 uppercase tracking-widest font-semibold mt-0.5">Acquire Relics</div>
            </div>
          </button>
        </div>

      </div>
    </div>
  );
}
