import { useState } from "react";
import { useGame } from "../../hooks/useGame";
import * as economyApi from "../../api/endpoints/economy";
import * as armyApi from "../../api/endpoints/army";
import type { ActiveBattle } from "../../app/providers/GameContext";

interface BattleResultScreenProps {
  battle: ActiveBattle;
}

const OUTCOME_LABEL: Record<string, string> = {
  victory: "VICTORY",
  defeat: "DEFEAT",
  threshold_failed: "PURGATION FAILED",
};

const OUTCOME_COLOR: Record<string, string> = {
  victory: "text-amber-400",
  defeat: "text-red-500",
  threshold_failed: "text-gray-400",
};

export default function BattleResultScreen({ battle }: BattleResultScreenProps) {
  const { api, dispatch } = useGame();
  const [closing, setClosing] = useState(false);

  const outcome = battle.outcome ?? "defeat";
  const destruction = battle.destruction ?? 0;
  const loot = battle.loot ?? 0;
  const duration = battle.duration ?? 0;
  const newSinMeter = battle.newSinMeter;

  const handleClose = async () => {
    setClosing(true);

    if (typeof newSinMeter === "number") {
      dispatch({ type: "SET_SIN_METER", payload: newSinMeter });
    }

    const [econRes, armyRes] = await Promise.all([
      economyApi.getEconomy(api),
      armyApi.getMyTroops(api),
    ]);
    if (econRes.success) dispatch({ type: "SET_ECONOMY", payload: econRes.data });
    if (armyRes.success) dispatch({ type: "SET_ARMY", payload: armyRes.data });

    dispatch({ type: "SET_ACTIVE_BATTLE", payload: null });
  };

  return (
    <div className="absolute inset-0 z-50 bg-black/95 flex flex-col items-center justify-center gap-8 p-8">
      <div className="text-center">
        <div className="text-[10px] uppercase tracking-[0.3em] text-gray-500 mb-2">
          vs {battle.defenderName}
        </div>
        <div className={`font-serif text-5xl font-bold tracking-widest ${OUTCOME_COLOR[outcome] ?? "text-gray-300"}`}>
          {OUTCOME_LABEL[outcome] ?? outcome.toUpperCase()}
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4 w-full max-w-md">
        <div className="bg-purgatory-card border border-purgatory-border rounded p-4">
          <div className="text-[10px] uppercase tracking-widest text-gray-500 mb-2">Destruction</div>
          <div className="text-2xl font-bold text-gray-200 mb-2">{destruction}%</div>
          <div className="h-2 bg-black/50 border border-purgatory-border rounded-sm overflow-hidden">
            <div
              className="h-full bg-gradient-to-r from-red-900 to-red-500"
              style={{ width: `${Math.min(100, Math.max(0, destruction))}%` }}
            />
          </div>
        </div>

        <div className="bg-purgatory-card border border-purgatory-border rounded p-4">
          <div className="text-[10px] uppercase tracking-widest text-gray-500 mb-2">Loot</div>
          <div className="text-2xl font-bold text-purple-400">{loot.toLocaleString()}</div>
          <div className="text-[10px] uppercase tracking-wider text-purple-500/70 mt-1">Penitence</div>
        </div>

        <div className="bg-purgatory-card border border-purgatory-border rounded p-4">
          <div className="text-[10px] uppercase tracking-widest text-gray-500 mb-2">Sin Meter</div>
          <div className="text-2xl font-bold text-red-400">
            {typeof newSinMeter === "number" ? `${newSinMeter}%` : "—"}
          </div>
        </div>

        <div className="bg-purgatory-card border border-purgatory-border rounded p-4">
          <div className="text-[10px] uppercase tracking-widest text-gray-500 mb-2">Duration</div>
          <div className="text-2xl font-bold text-gray-200">{(duration / 20).toFixed(1)}s</div>
        </div>
      </div>

      <button
        type="button"
        onClick={handleClose}
        disabled={closing}
        className="text-xs uppercase tracking-widest font-bold px-10 py-3 rounded border transition-all
          enabled:border-amber-500/50 enabled:text-amber-400 enabled:hover:bg-amber-500/10 enabled:hover:border-amber-400
          disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {closing ? "Returning..." : "Return to Terrace"}
      </button>
    </div>
  );
}
