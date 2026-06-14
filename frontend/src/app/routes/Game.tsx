import { useState, useCallback, useEffect } from "react";
import GameCanvas from '../../game/phaser/GameCanvas';
import { useAuth } from '../../hooks/useAuth';
import { useGame } from '../../hooks/useGame';
import GameHud from '../../ui/layout/GameHud';
import ShopPanel from '../../ui/panels/ShopPanel';
import PlacementToolbar from '../../ui/panels/PlacementToolbar';
import { phaserEvents } from '../../game/phaser/events';
import * as baseApi from '../../api/endpoints/base';
import type { PlacedBuilding } from '../../types/building';

export default function GameDashboard() {
  const { user, logout } = useAuth();
  const { state, api, dispatch } = useGame();
  const [shopOpen, setShopOpen] = useState(false);
  const [buildingMenu, setBuildingMenu] = useState<PlacedBuilding | null>(null);

  const refetchLayout = useCallback(async () => {
    const res = await baseApi.getLayout(api);
    if (res.success) {
      dispatch({ type: "SET_LAYOUT", payload: res.data });
    }

    const econRes = await api.get<{ penitence: number; grace: number; max_penitence: number }>("/user/economy");
    if (econRes.success) {
      dispatch({ type: "SET_ECONOMY", payload: econRes.data });
    }
  }, [api, dispatch]);

  useEffect(() => {
    phaserEvents.onBuildingClick = (b) => setBuildingMenu(b);
    return () => { phaserEvents.onBuildingClick = null; };
  }, []);

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

      <PlacementToolbar
        buildingMenu={buildingMenu}
        onCloseMenu={() => setBuildingMenu(null)}
        onLayoutChanged={refetchLayout}
      />
    </div>
  );
}
