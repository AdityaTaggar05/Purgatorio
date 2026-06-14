import { useState } from "react";
import GameCanvas from '../../game/phaser/GameCanvas';
import { useAuth } from '../../hooks/useAuth';
import { useGame } from '../../hooks/useGame';
import GameHud from '../../ui/layout/GameHud';
import ShopPanel from '../../ui/panels/ShopPanel';

export default function GameDashboard() {
  const { user, logout } = useAuth();
  const { state } = useGame();
  const [shopOpen, setShopOpen] = useState(false);

  const handleAttack = () => console.log('Initiating Combat Encounter Instance...');
  const handleArmy = () => console.log('Opening Legion Management Array...');

  return (
    <div className="relative w-screen h-screen overflow-hidden bg-[#111111]">
      <GameCanvas layout={state.layout} />

      <GameHud
        username={user?.username || "Unknown Penitent"}
        level={user?.level || 1}
        economy={state.economy}
        sinMeter={state.sinMeter}
        onAscensionClick={() => setShopOpen(true)}
        onLogoutClick={logout}
        onAttackClick={handleAttack}
        onArmyClick={handleArmy}
      />

      <ShopPanel open={shopOpen} onClose={() => setShopOpen(false)} />
    </div>
  );
}
