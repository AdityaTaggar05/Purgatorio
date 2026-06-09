import GameCanvas from '../../game/phaser/GameCanvas';

export default function GameDashboard() {
  return (
    <div className="relative w-screen h-screen overflow-hidden bg-[#111111]">

      <GameCanvas />

      <div className="absolute inset-0 pointer-events-none z-10 p-6 flex flex-col justify-between">
        <header className="pointer-events-auto flex justify-between items-center bg-[#181716]/80 backdrop-blur-md border border-[#2c2927] p-4 rounded-lg shadow-lg max-w-2xl mx-auto w-full">
          <div>
            <h1 className="font-serif text-xl font-bold tracking-widest text-gray-200">
              FIRST TERRACE
            </h1>
            <p className="text-xs text-amber-500/80 uppercase tracking-[0.2em]">
              The Proud
            </p>
          </div>
          <button className="bg-[#1e1d1c] hover:bg-amber-950/50 border border-[#2c2927] text-gray-300 px-4 py-2 rounded transition-all">
            Inventory
          </button>
        </header>
      </div>
    </div>
  );
}
