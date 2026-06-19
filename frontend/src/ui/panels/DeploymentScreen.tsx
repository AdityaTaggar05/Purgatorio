import { useState, useMemo } from "react";
import { useGame } from "../../hooks/useGame";
import type { TroopDeployment } from "../../types/battle";
import type { ActiveBattle } from "../../app/providers/GameContext";
import type { UseBattleSocketResult } from "../../hooks/useBattleSocket";

interface DeploymentScreenProps {
  battle: ActiveBattle;
  socket: UseBattleSocketResult;
  deployments: TroopDeployment[];
  selectedTroop: string | null;
  onSelectTroop: (id: string | null) => void;
  deployCounts: Record<string, number>;
  deployError: string | null;
  onIncrementCount: (troopId: string) => void;
  onDecrementCount: (troopId: string) => void;
  onStartBattle: () => void;
  onSetDeployError: (err: string | null) => void;
}

function TroopThumb({ troopType }: { troopType: string }) {
  const [failed, setFailed] = useState(false);

  if (failed) {
    return (
      <div className="w-8 h-8 shrink-0 rounded bg-red-900/40 border border-red-700/40 flex items-center justify-center text-[9px] font-bold text-red-300 uppercase">
        {troopType.slice(0, 2)}
      </div>
    );
  }
  return (
    <img
      src={`/assets/${troopType}.png`}
      alt={troopType}
      onError={() => setFailed(true)}
      className="w-8 h-8 shrink-0 object-contain"
      draggable={false}
    />
  );
}

export default function DeploymentScreen({
  battle,
  socket,
  deployments,
  selectedTroop,
  onSelectTroop,
  deployCounts,
  deployError,
  onIncrementCount,
  onDecrementCount,
  onStartBattle,
  onSetDeployError,
}: DeploymentScreenProps) {
  const { state, dispatch } = useGame();
  const { state: socketState, error: socketError, deployCountdown } = socket;
  const catalog = state.troopCatalog ?? [];
  const army = state.army?.troops ?? {};

  const getAvailableCount = (troopId: string): number => {
    const total = army[troopId] ?? 0;
    const deployed = deployments
      .filter((d) => d.troop_type === troopId)
      .reduce((s, d) => s + d.count, 0);
    return Math.max(0, total - deployed);
  };

  const totalDeployed = deployments.reduce((sum, d) => sum + d.count, 0);
  const canDeploy = socketState === "open" || socketState === "deployed";
  const isStartEnabled = canDeploy && deployments.length > 0;

  const handleCancel = () => dispatch({ type: "SET_ACTIVE_BATTLE", payload: null });

  return (
    <div className="absolute inset-0 flex flex-col z-10 pointer-events-none">
      {/* Top bar */}
      <div className="flex items-center justify-between p-4 shrink-0 pointer-events-auto">
        <div>
          <h2 className="font-serif text-lg font-bold tracking-wider text-red-400">Deploy Forces</h2>
          <div className="text-[10px] uppercase tracking-widest text-gray-500">vs {battle.defenderName}</div>
        </div>
        <div className="text-right">
          <div className="text-[10px] uppercase tracking-wider text-gray-500">Deploy Window</div>
          <div className={`text-lg font-bold font-sans ${deployCountdown <= 5 ? "text-red-500" : "text-gray-300"}`}>
            {deployCountdown}s
          </div>
        </div>
      </div>

      {/* Main area */}
      <div className="flex-1 flex gap-4 px-4 min-h-0">
        {/* Troop selector */}
        <div className="w-72 shrink-0 flex flex-col gap-2 overflow-y-auto pointer-events-auto">
          <div className="text-xs uppercase tracking-wider text-gray-500 font-bold mb-1">Your Legion</div>
          {catalog.map((troop) => {
            const available = getAvailableCount(troop.id);
            const count = deployCounts[troop.id] ?? 1;
            const isSelected = selectedTroop === troop.id;

            return (
              <div
                key={troop.id}
                className={`p-3 rounded border transition-all ${
                  isSelected
                    ? "bg-red-900/20 border-red-600/50"
                    : "bg-purgatory-input border-purgatory-border hover:border-gray-600"
                } ${available === 0 ? "opacity-30" : ""}`}
              >
                <button
                  type="button"
                  onClick={() => onSelectTroop(isSelected ? null : troop.id)}
                  disabled={available === 0}
                  className="w-full text-left flex items-center gap-3 disabled:cursor-not-allowed"
                >
                  <TroopThumb troopType={troop.id} />
                  <div className="flex-1">
                    <div className="flex items-center justify-between">
                      <div className="text-gray-200 font-bold text-xs tracking-wide">{troop.name}</div>
                      <div className="text-[10px] text-gray-500">{available} left</div>
                    </div>
                    <div className="text-[10px] text-gray-500 mt-0.5">HP {troop.hp} · DPS {troop.dps}</div>
                  </div>
                </button>

                {isSelected && available > 0 && (
                  <div className="flex items-center gap-1.5 mt-2 pl-11">
                    <button
                      type="button"
                      onClick={() => onDecrementCount(troop.id)}
                      className="w-5 h-5 flex items-center justify-center text-xs border border-purgatory-border rounded text-gray-500 hover:text-gray-200"
                    >
                      −
                    </button>
                    <span className="text-xs text-gray-200 w-5 text-center font-bold">{count}</span>
                    <button
                      type="button"
                      onClick={() => onIncrementCount(troop.id)}
                      className="w-5 h-5 flex items-center justify-center text-xs border border-purgatory-border rounded text-gray-500 hover:text-gray-200"
                    >
                      +
                    </button>
                  </div>
                )}
              </div>
            );
          })}
          {catalog.length === 0 && <div className="text-xs text-gray-600 italic">No troops trained yet.</div>}
        </div>

        <div className="flex-1" /> {/* spacer for Phaser canvas */}

        {/* Deployed list */}
        <div className="w-48 shrink-0 flex flex-col gap-2 overflow-y-auto pointer-events-auto">
          <div className="text-xs uppercase tracking-wider text-gray-500 font-bold mb-1">Deployed</div>
          <div className="text-sm text-gray-400 mb-2">{totalDeployed} troops deployed</div>
          {deployments.map((d, i) => (
            <div key={i} className="bg-purgatory-input border border-purgatory-border rounded p-2 text-xs">
              <div className="text-gray-300 font-bold">{d.count}× {d.troop_type}</div>
              <div className="text-gray-500">at ({d.position.x}, {d.position.y})</div>
            </div>
          ))}
          {deployments.length === 0 && (
            <div className="text-xs text-gray-600 italic">
              Select a troop, then click a cell in the green deployment zone.
            </div>
          )}
        </div>
      </div>

      {/* Bottom bar */}
      <div className="flex items-center justify-between p-4 bg-purgatory-card/80 border-t border-purgatory-border shrink-0 pointer-events-auto">
        <div className="flex-1">
          {(deployError || socketError) && (
            <div className="text-sm text-red-400">{deployError || socketError}</div>
          )}
          {socketState === "connecting" && (
            <div className="text-sm text-amber-400 flex items-center gap-2">
              <span className="w-3 h-3 border border-amber-500/30 border-t-amber-400 rounded-full animate-spin inline-block" />
              Connecting to battle server...
            </div>
          )}
        </div>
        <div className="flex gap-3">
          <button
            type="button"
            onClick={handleCancel}
            className="text-xs uppercase tracking-widest font-bold px-6 py-2 rounded border border-purgatory-border text-gray-500 hover:text-gray-300 hover:border-gray-600 transition-all"
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={onStartBattle}
            disabled={!isStartEnabled}
            className="text-xs uppercase tracking-widest font-bold px-8 py-2 rounded border transition-all
              enabled:border-red-600/60 enabled:text-red-400 enabled:hover:bg-red-900/30 enabled:hover:border-red-500
              disabled:border-gray-700 disabled:text-gray-600 disabled:cursor-not-allowed"
          >
            {socketState === "connecting" ? "Connecting..." : "Start Battle"}
          </button>
        </div>
      </div>
    </div>
  );
}
