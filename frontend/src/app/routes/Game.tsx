import GameCanvas from '../../game/phaser/GameCanvas';
import { useAuth } from '../../hooks/useAuth';
import { useGame } from '../../hooks/useGame';
import GameHud from '../../ui/layout/GameHud';

export default function GameDashboard() {
  const { user, logout } = useAuth();
  const { state } = useGame();

  const handleAscension = () => console.log('Opening Altar Store Interface...');
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
        onAscensionClick={handleAscension}
        onLogoutClick={logout}
        onAttackClick={handleAttack}
        onArmyClick={handleArmy}
      />
    </div>
  );
}
