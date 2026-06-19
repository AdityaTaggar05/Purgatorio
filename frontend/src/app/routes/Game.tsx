import { useState, useCallback, useEffect, useRef } from "react";
import GameCanvas from '../../game/phaser/GameCanvas';
import { useAuth } from '../../hooks/useAuth';
import { useGame } from '../../hooks/useGame';
import GameHud from '../../ui/layout/GameHud';
import ShopPanel from '../../ui/panels/ShopPanel';
import ArmyPanel from '../../ui/panels/ArmyPanel';
import MatchmakingPanel from '../../ui/panels/MatchmakingPanel';
import BattleOverlay from '../../ui/battle/BattleOverlay';
import PlacementToolbar from '../../ui/panels/PlacementToolbar';
import UpgradeSnackbar from '../../ui/panels/UpgradeSnackbar';
import { phaserEvents } from '../../game/phaser/events';
import * as baseApi from '../../api/endpoints/base';
import * as economyApi from '../../api/endpoints/economy';
import type { PlacedBuilding } from '../../types/building';

const SIN_DRAIN_PER_MINUTE = 1 / 6; // 10% per hour; matches backend
const DRAIN_INTERVAL_MS = 10_000;

export default function GameDashboard() {
  const { user, logout } = useAuth();
  const { state, api, dispatch } = useGame();
  const [shopOpen, setShopOpen] = useState(false);
  const [armyOpen, setArmyOpen] = useState(false);
  const [matchmakingOpen, setMatchmakingOpen] = useState(false);
  const [buildingMenu, setBuildingMenu] = useState<PlacedBuilding | null>(null);
  const [snackbarMsg, setSnackbarMsg] = useState<string | null>(null);

  const [displaySin, setDisplaySin] = useState(state.sinMeter);
  const sinInfoRef = useRef({ sin: state.sinMeter, syncedAt: Date.now() });

  useEffect(() => {
    sinInfoRef.current = { sin: state.sinMeter, syncedAt: Date.now() };
    setDisplaySin(state.sinMeter);
  }, [state.sinMeter]);

  useEffect(() => {
    const interval = setInterval(() => {
      const { sin, syncedAt } = sinInfoRef.current;
      const elapsedMin = (Date.now() - syncedAt) / 60000;
      const drained = Math.floor(elapsedMin * SIN_DRAIN_PER_MINUTE);
      setDisplaySin(Math.max(0, sin - drained));
    }, DRAIN_INTERVAL_MS);
    return () => clearInterval(interval);
  }, []);

  const drainTimeEstimate = displaySin > 0
    ? Math.ceil(displaySin / (SIN_DRAIN_PER_MINUTE * 60))
    : 0; // hours

  useEffect(() => {
    if (state.checkInResult) {
      setSnackbarMsg(state.checkInResult);
      dispatch({ type: "SET_CHECK_IN_RESULT", payload: null });
    }
  }, [state.checkInResult, dispatch]);

  const selectBuilding = useCallback((b: PlacedBuilding | null) => {
    setBuildingMenu(b);

    const buildings = phaserEvents.getActiveBuildings?.() ?? [];
    buildings.forEach(sprite => {
      sprite.selected = !!(b && sprite.buildingData.x === b.x && sprite.buildingData.y === b.y && sprite.buildingData.building_id === b.building_id);
    });
  }, []);

  const refetchLayout = useCallback(async () => {
    const res = await baseApi.getLayout(api);
    if (res.success) {
      dispatch({ type: "SET_LAYOUT", payload: res.data });
    }

    const econRes = await api.get<{ penitence: number; grace: number; max_penitence: number; overflow_penitence?: number }>("/user/economy");
    if (econRes.success) {
      dispatch({ type: "SET_ECONOMY", payload: econRes.data });
    }
  }, [api, dispatch]);

  useEffect(() => {
    phaserEvents.onBuildingClick = selectBuilding;
    return () => { phaserEvents.onBuildingClick = null; };
  }, [selectBuilding]);

  const handleAttack = () => setMatchmakingOpen(true);

  const handleCheckIn = useCallback(async () => {
    const res = await baseApi.checkIn(api);
    if (res.success) {
      await refetchLayout();
      if (res.data.completed_upgrades.length > 0) {
        const names = res.data.completed_upgrades
          .map(u => `${u.building_id} Lv.${u.from_level} → ${u.to_level}`)
          .join(", ");
        setSnackbarMsg(`Upgrades completed: ${names}`);
      }
    }
  }, [api, refetchLayout]);

  const handleCollect = useCallback(async () => {
    const res = await economyApi.collectResources(api);
    if (res.success) {
      dispatch({ type: "SET_ECONOMY", payload: res.data });
    }
  }, [api, dispatch]);

  return (
    <div className="relative w-screen h-screen overflow-hidden bg-[#111111]">
      <GameCanvas layout={state.layout} />

      <GameHud
        username={user?.username || "Unknown Penitent"}
        level={user?.level || 1}
        economy={state.economy}
        sinMeter={displaySin}
        onAscensionClick={() => setShopOpen(true)}
        onLogoutClick={logout}
        onAttackClick={handleAttack}
        onArmyClick={() => setArmyOpen(true)}
        onCheckInClick={handleCheckIn}
        onCollectClick={handleCollect}
      />

      <ShopPanel open={shopOpen} onClose={() => setShopOpen(false)} />
      <ArmyPanel open={armyOpen} onClose={() => setArmyOpen(false)} />
      <MatchmakingPanel open={matchmakingOpen} onClose={() => setMatchmakingOpen(false)} />

      {state.activeBattle && <BattleOverlay battle={state.activeBattle} />}

      <PlacementToolbar
        buildingMenu={buildingMenu}
        onCloseMenu={() => selectBuilding(null)}
        onLayoutChanged={refetchLayout}
      />

      {snackbarMsg && (
        <UpgradeSnackbar message={snackbarMsg} onDone={() => setSnackbarMsg(null)} />
      )}
    </div>
  );
}
