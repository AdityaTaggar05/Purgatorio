import { useState } from 'react';
import GameCanvas from '../../game/phaser/GameCanvas';
import { useAuth } from '../../hooks/useAuth';
import GameHud from '../../ui/layout/GameHud';

export default function GameDashboard() {
  const { user, logout } = useAuth();

  // Mock variables to be finally fetched from the backend
  const [penitence, setPenitence] = useState({ current: 420, max: 1000, label: 'Penitence' });
  const [grace, setGrace] = useState({ current: 85, max: 120, label: 'Grace' });
  const [sinMeter, setSinMeter] = useState(68);

  const handleAscension = () => console.log('Opening Altar Store Interface...');
  const handleAttack = () => console.log('Initiating Combat Encounter Instance...');
  const handleArmy = () => console.log('Opening Legion Management Array...');

  return (
    <div className="relative w-screen h-screen overflow-hidden bg-[#111111]">

      <GameCanvas />

      <GameHud
        username={user?.username || "Unknown Penitent"}
        level={4}
        penitence={penitence}
        grace={grace}
        sinMeter={sinMeter}
        onAscensionClick={handleAscension}
        onLogoutClick={logout}
        onAttackClick={handleAttack}
        onArmyClick={handleArmy}
      />
    </div>
  );
}
