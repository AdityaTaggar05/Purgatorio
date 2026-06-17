import { useState, useMemo } from "react";
import { useGame } from "../../hooks/useGame";
import BattleGrid from "../battle/BattleGrid";
import TroopStackOverlay, { type TroopStack } from "../battle/TroopStackOverlay";
import type { TroopDeployment } from "../../types/battle";
import type { PlacedBuilding } from "../../types/building";
import type { ActiveBattle } from "../../app/providers/GameContext";
import type { UseBattleSocketResult } from "../../hooks/useBattleSocket";

interface DeploymentScreenProps {
  battle: ActiveBattle;
  socket: UseBattleSocketResult;
}

const CELL_SIZE = 14;
const GRID_SIZE = 30;

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

export default function DeploymentScreen({ battle, socket }: DeploymentScreenProps) {
  const { state, dispatch } = useGame();
  const { sendDeploy, state: socketState, error: socketError, deployCountdown } = socket;

  const catalog = state.troopCatalog ?? [];
  const army = state.army?.troops ?? {};

  const [deployCounts, setDeployCounts] = useState<Record<string, number>>({});
  const [deployments, setDeployments] = useState<TroopDeployment[]>([]);
  const [selectedTroop, setSelectedTroop] = useState<string | null>(null);
  const [deployError, setDeployError] = useState<string | null>(null);

  const defenderBuildings: PlacedBuilding[] = useMemo(() => [], []);

  const deployedMap = useMemo(() => {
    const cells = new Map<string, TroopDeployment>();
    for (const d of deployments) cells.set(`${d.position.x},${d.position.y}`, d);
    return cells;
  }, [deployments]);

  const deploymentZone = useMemo(() => {
    const cells: { x: number; y: number }[] = [];
    for (let y = 0; y < GRID_SIZE; y++) {
      for (let x = 0; x < 3; x++) cells.push({ x, y });
    }
    return cells;
  }, []);

  const totalDeployed = deployments.reduce((sum, d) => sum + d.count, 0);

  const getAvailableCount = (troopId: string): number => {
    const total = army[troopId] ?? 0;
    const deployed = deployments
      .filter((d) => d.troop_type === troopId)
      .reduce((s, d) => s + d.count, 0);
    return Math.max(0, total - deployed);
  };

  const handleCellClick = (x: number, y: number) => {
    if (!selectedTroop) return;
    const key = `${x},${y}`;
    const existing = deployedMap.get(key);

    if (existing) {
      if (existing.troop_type === selectedTroop) {
        setDeployments((prev) => prev.filter((d) => d !== existing));
        setDeployError(null);
      } else {
        setDeployError(`That cell already has ${existing.troop_type} deployed`);
      }
      return;
    }

    const count = deployCounts[selectedTroop] ?? 1;
    const available = getAvailableCount(selectedTroop);
    if (count > available) {
      setDeployError(`Only ${available} ${selectedTroop} available`);
      return;
    }

    setDeployments((prev) => [...prev, { troop_type: selectedTroop, position: { x, y }, count }]);
    setDeployError(null);
  };

  const incrementCount = (troopId: string) => {
    setDeployCounts((prev) => {
      const current = prev[troopId] ?? 1;
      const available = getAvailableCount(troopId);
      return { ...prev, [troopId]: Math.min(current + 1, Math.max(1, available)) };
    });
  };

  const decrementCount = (troopId: string) => {
    setDeployCounts((prev) => ({ ...prev, [troopId]: Math.max(1, (prev[troopId] ?? 1) - 1) }));
  };

  const handleDeploy = () => {
    if (deployments.length === 0) {
      setDeployError("Deploy at least one troop");
      return;
    }
    if (socketState !== "open") {
      setDeployError("Battle server not connected");
      return;
    }
    sendDeploy(deployments);
    dispatch({
      type: "SET_ACTIVE_BATTLE",
      payload: { ...battle, phase: "viewing", deployment: deployments },
    });
  };

  const handleCancel = () => {
    dispatch({ type: "SET_ACTIVE_BATTLE", payload: null });
  };

  const selectedCells = deployments.map((d) => d.position);

  const previewStacks: TroopStack[] = useMemo(
    () =>
      deployments.map((d, i) => {
        const troopDef = catalog.find((t) => t.id === d.troop_type);
        const unitHp = troopDef?.hp ?? 1;
        return {
          id: `preview_${i}`,
          troopType: d.troop_type,
          x: d.position.x,
          y: d.position.y,
          unitMaxHp: unitHp,
          totalUnits: d.count,
          hp: unitHp * d.count,
          alive: true,
        };
      }),
    [deployments, catalog]
  );

  const canDeploy = socketState === "open";
  const isDeployEnabled = canDeploy && deployments.length > 0;

  return (
    <div className="absolute inset-0 z-50 bg-black/95 flex flex-col">
      <div className="flex-1 flex flex-col p-4 min-h-0">
        <div className="flex items-center justify-between mb-3 shrink-0">
          <div>
            <h2 className="font-serif text-lg font-bold tracking-wider text-red-400">
              Deploy Forces
            </h2>
            <div className="text-[10px] uppercase tracking-widest text-gray-500">
              vs {battle.defenderName}
            </div>
          </div>
          <div className="text-right">
            <div className="text-[10px] uppercase tracking-wider text-gray-500">Deploy Window</div>
            <div className={`text-lg font-bold font-sans ${deployCountdown <= 5 ? "text-red-500" : "text-gray-300"}`}>
              {deployCountdown}s
            </div>
          </div>
        </div>

        <div className="flex-1 flex gap-4 min-h-0">
          <div className="w-72 shrink-0 flex flex-col gap-2 overflow-y-auto">
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
                    onClick={() => setSelectedTroop(isSelected ? null : troop.id)}
                    disabled={available === 0}
                    className="w-full text-left flex items-center gap-3 disabled:cursor-not-allowed"
                  >
                    <TroopThumb troopType={troop.id} />
                    <div className="flex-1">
                      <div className="flex items-center justify-between">
                        <div className="text-gray-200 font-bold text-xs tracking-wide">{troop.name}</div>
                        <div className="text-[10px] text-gray-500">{available} left</div>
                      </div>
                      <div className="text-[10px] text-gray-500 mt-0.5">
                        HP {troop.hp} · DPS {troop.dps}
                      </div>
                    </div>
                  </button>

                  {isSelected && available > 0 && (
                    <div className="flex items-center gap-1.5 mt-2 pl-11">
                      <button
                        type="button"
                        onClick={() => decrementCount(troop.id)}
                        className="w-5 h-5 flex items-center justify-center text-xs border border-purgatory-border rounded text-gray-500 hover:text-gray-200"
                      >
                        −
                      </button>
                      <span className="text-xs text-gray-200 w-5 text-center font-bold">{count}</span>
                      <button
                        type="button"
                        onClick={() => incrementCount(troop.id)}
                        className="w-5 h-5 flex items-center justify-center text-xs border border-purgatory-border rounded text-gray-500 hover:text-gray-200"
                      >
                        +
                      </button>
                    </div>
                  )}
                </div>
              );
            })}

            {catalog.length === 0 && (
              <div className="text-xs text-gray-600 italic">No troops trained yet.</div>
            )}
          </div>

          <div className="flex-1 flex items-center justify-center overflow-auto">
            <div className="relative" style={{ width: GRID_SIZE * CELL_SIZE, height: GRID_SIZE * CELL_SIZE }}>
              <BattleGrid
                buildings={defenderBuildings}
                gridW={GRID_SIZE}
                gridH={GRID_SIZE}
                cellSize={CELL_SIZE}
                deploymentZone={deploymentZone}
                selectedCells={selectedCells}
                onCellClick={canDeploy ? handleCellClick : undefined}
              />
              <TroopStackOverlay stacks={previewStacks} cellSize={CELL_SIZE} staticPreview />
            </div>
          </div>

          <div className="w-48 shrink-0 flex flex-col gap-2 overflow-y-auto">
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
      </div>

      <div className="flex items-center justify-between p-4 bg-purgatory-card/80 border-t border-purgatory-border shrink-0">
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
            onClick={handleDeploy}
            disabled={!isDeployEnabled}
            className="text-xs uppercase tracking-widest font-bold px-8 py-2 rounded border transition-all
              enabled:border-red-600/60 enabled:text-red-400 enabled:hover:bg-red-900/30 enabled:hover:border-red-500
              disabled:border-gray-700 disabled:text-gray-600 disabled:cursor-not-allowed"
          >
            {socketState === "connecting" ? "Connecting..." : "Deploy"}
          </button>
        </div>
      </div>
    </div>
  );
}
