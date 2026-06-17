import { useEffect } from "react";
import { useGame } from "../../hooks/useGame";
import { useBattleSocket } from "../../hooks/useBattleSocket";
import type { ActiveBattle } from "../../app/providers/GameContext";
import DeploymentScreen from "../panels/DeploymentScreen";
import BattleViewer from "../panels/BattleViewer";
import BattleResultScreen from "../panels/BattleResultScreen";
import { useAuth } from "../../hooks/useAuth";

interface BattleOverlayProps {
  battle: ActiveBattle;
}

export default function BattleOverlay({ battle }: BattleOverlayProps) {
  const { dispatch } = useGame();
  const { accessToken } = useAuth()
  const socket = useBattleSocket(battle.battleId, accessToken);

  useEffect(() => {
    if (socket.battleResult && battle.phase === "viewing") {
      dispatch({
        type: "SET_ACTIVE_BATTLE",
        payload: {
          ...battle,
          phase: "result",
          outcome: socket.battleResult.outcome,
          destruction: socket.battleResult.destruction,
          loot: socket.battleResult.loot,
          newSinMeter: socket.battleResult.sin_meter,
          duration: socket.battleResult.duration_ticks,
        },
      });
    }
  }, [socket.battleResult, battle, dispatch]);

  if (battle.phase === "deploying") {
    return <DeploymentScreen battle={battle} socket={socket} />;
  }
  if (battle.phase === "viewing") {
    return <BattleViewer battle={battle} socket={socket} />;
  }
  if (battle.phase === "result") {
    return <BattleResultScreen battle={battle} />;
  }
  return null;
}
