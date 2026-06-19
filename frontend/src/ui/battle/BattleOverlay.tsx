import { useEffect, useMemo, useRef, useState } from "react";
import { useGame } from "../../hooks/useGame";
import { useBattleSocket } from "../../hooks/useBattleSocket";
import type { ActiveBattle } from "../../app/providers/GameContext";
import type { TroopStackData } from "../../game/phaser/entities/TroopStackSprite";
import type { HpChange, PositionChange, TickResult, TroopDeployment } from "../../types/battle";
import BattleCanvas from "../../game/phaser/BattleCanvas";
import DeploymentScreen from "../panels/DeploymentScreen";
import BattleResultScreen from "../panels/BattleResultScreen";
import { useAuth } from "../../hooks/useAuth";

interface BattleOverlayProps {
  battle: ActiveBattle;
}

const PAD = 6;
const TICK_INTERVAL_MS = 100;

function parseTroopSeq(entityId: string): number | null {
  const parts = entityId.split("_");
  if (parts.length < 3 || parts[0] !== "troop") return null;
  const seq = parseInt(parts[parts.length - 1], 10);
  return isNaN(seq) ? null : seq;
}

function applyHPChanges(stacks: TroopStackData[], changes: HpChange[]): TroopStackData[] {
  const next = [...stacks];
  for (const change of changes) {
    if (change.entity_type !== "troop") continue;
    const seq = parseTroopSeq(change.entity_id);
    if (seq === null) continue;
    const idx = seq - 1;
    if (idx < 0 || idx >= next.length) continue;
    const newHp = Math.max(0, next[idx].hp + change.delta);
    next[idx] = { ...next[idx], hp: newHp, alive: newHp > 0 };
  }
  return next;
}

export default function BattleOverlay({ battle }: BattleOverlayProps) {
  const { state, dispatch } = useGame();
  const { accessToken } = useAuth();
  const socket = useBattleSocket(battle.battleId, accessToken);
  const catalog = state.troopCatalog ?? [];
  const army = state.army?.troops ?? {};

  // Deployment state
  const [deployments, setDeployments] = useState<TroopDeployment[]>([]);
  const [selectedTroop, setSelectedTroop] = useState<string | null>(null);
  const [deployCounts, setDeployCounts] = useState<Record<string, number>>({});
  const [deployError, setDeployError] = useState<string | null>(null);

  // Playback state
  const [viewingStacks, setViewingStacks] = useState<TroopStackData[]>([]);
  const [currentTick, setCurrentTick] = useState(0);
  const [resultLoading, setResultLoading] = useState(false);
  const stacksRef = useRef<TroopStackData[]>([]);
  const cursorRef = useRef(0);
  const ticksRef = useRef<TickResult[]>(socket.ticks);
  const [buildingHpByKey, setBuildingHpByKey] = useState<Record<string, number>>({});
  const buildingHpRef = useRef<Record<string, number>>({});

  const gridW = battle.defenderLayout?.grid_w ?? 30;
  const gridH = battle.defenderLayout?.grid_h ?? 30;

  // Compute deployment zone
  const deploymentZone = useMemo(() => {
    const cells: { x: number; y: number }[] = [];
    const maxX = gridW + PAD * 2;
    const maxY = gridH + PAD * 2;
    for (let y = -PAD; y < gridH + PAD; y++) {
      for (let x = -PAD; x < gridW + PAD; x++) {
        const inBase = x >= 0 && x < gridW && y >= 0 && y < gridH;
        if (!inBase) cells.push({ x, y });
      }
    }
    return cells;
  }, [gridW, gridH]);

  const getAvailableCount = (troopId: string): number => {
    const total = army[troopId] ?? 0;
    const deployed = deployments
      .filter((d) => d.troop_type === troopId)
      .reduce((s, d) => s + d.count, 0);
    return Math.max(0, total - deployed);
  };

  // Keep ticksRef in sync
  useEffect(() => { ticksRef.current = socket.ticks; }, [socket.ticks]);

  // Initialise playback stacks — one per individual unit
  useEffect(() => {
    if (battle.phase !== "viewing") return;
    const initial: TroopStackData[] = [];
    deployments.forEach((d, di) => {
      const troopDef = catalog.find((t) => t.id === d.troop_type);
      const unitHp = troopDef?.hp ?? 1;
      for (let j = 0; j < d.count; j++) {
        const offset = (j - (d.count - 1) / 2) * 0.35;
        initial.push({
          id: `stack_${di}_${j}`,
          troopType: d.troop_type,
          x: d.position.x + offset,
          y: d.position.y,
          unitMaxHp: unitHp,
          totalUnits: 1,
          hp: unitHp,
          alive: true,
        });
      }
    });
    setViewingStacks(initial);
    stacksRef.current = initial;
    cursorRef.current = 0;
    buildingHpRef.current = {};
    setBuildingHpByKey({});
  }, [battle.phase, deployments, catalog]);

  // Map each backend entity position to its matching stack by sequence number
  const applyPositions = (stacks: TroopStackData[], positions: PositionChange[]) => {
    const next = [...stacks];
    for (const p of positions) {
      const seq = parseTroopSeq(p.entity_id);
      if (seq === null) continue;
      const idx = seq - 1;
      if (idx >= 0 && idx < next.length) {
        next[idx] = { ...next[idx], x: p.x, y: p.y };
      }
    }
    return next;
  };

  // Playback tick loop
  useEffect(() => {
    if (battle.phase !== "viewing") return;
    const interval = setInterval(() => {
      const ticks = ticksRef.current;
      if (cursorRef.current >= ticks.length) return;
      const tick = ticks[cursorRef.current];
      cursorRef.current += 1;

      // Apply troop HP changes
      let updated = applyHPChanges(stacksRef.current, tick.hp_changes);

      // Apply troop position changes
      if (tick.positions && tick.positions.length > 0) {
        updated = applyPositions(updated, tick.positions);
      }
      stacksRef.current = updated;
      setViewingStacks([...updated]);

      // Apply building HP changes
      let buildingChanged = false;
      for (const change of tick.hp_changes) {
        if (change.entity_type !== "building") continue;
        // Entity ID format: buildingtype_x_y -> extract "x_y" key
        const parts = change.entity_id.split("_");
        if (parts.length < 3) continue;
        const key = parts.slice(-2).join("_");
        const prev = buildingHpRef.current[key];
        if (prev !== change.new_hp) {
          buildingHpRef.current[key] = change.new_hp;
          buildingChanged = true;
        }
      }
      if (buildingChanged) {
        setBuildingHpByKey({ ...buildingHpRef.current });
      }

      setCurrentTick(tick.tick);
    }, TICK_INTERVAL_MS);
    return () => clearInterval(interval);
  }, [battle.phase]);

  // Auto-transition to result when battle ends
  useEffect(() => {
    if (socket.battleResult && battle.phase === "viewing") {
      setResultLoading(false);
      dispatch({
        type: "SET_ACTIVE_BATTLE",
        payload: { ...battle, phase: "result", outcome: socket.battleResult.outcome, destruction: socket.battleResult.destruction, loot: socket.battleResult.loot, newSinMeter: socket.battleResult.sin_meter, duration: socket.battleResult.duration_ticks },
      });
    }
  }, [socket.battleResult, battle, dispatch]);

  const handleCellClick = (x: number, y: number) => {
    if (!selectedTroop) return;
    const key = `${x},${y}`;

    const existing = deployments.find(
      (d) => `${d.position.x},${d.position.y}` === key
    );
    if (existing) {
      if (existing.troop_type === selectedTroop) {
        setDeployments((prev) => prev.filter((d) => `${d.position.x},${d.position.y}` !== key));
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

    const deployment = { troop_type: selectedTroop, position: { x, y }, count };
    setDeployments((prev) => [...prev, deployment]);
    socket.sendDeploy([deployment]);
    setDeployError(null);

    const remaining = getAvailableCount(selectedTroop) - count;
    if (remaining <= 0) {
      const nextTroop = catalog.find(
        (t) => getAvailableCount(t.id) > 0 && t.id !== selectedTroop
      );
      setSelectedTroop(nextTroop?.id ?? null);
    }
  };

  const handleEndBattle = () => {
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
    } else {
      socket.sendSkip();
      setResultLoading(true);
      cursorRef.current = Infinity;
    }
  };

  const handleStartBattle = () => {
    if (deployments.length === 0) {
      setDeployError("Deploy at least one troop");
      return;
    }
    if (socket.state !== "open" && socket.state !== "deployed") {
      setDeployError("Battle server not connected");
      return;
    }
    socket.sendDone();
    dispatch({
      type: "SET_ACTIVE_BATTLE",
      payload: { ...battle, phase: "viewing", deployment: deployments },
    });
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

  // Preview stacks during deployment — one per individual unit
  const previewStacks = useMemo<TroopStackData[]>(
    () => {
      const stacks: TroopStackData[] = [];
      deployments.forEach((d, di) => {
        const troopDef = catalog.find((t) => t.id === d.troop_type);
        const unitHp = troopDef?.hp ?? 1;
        for (let j = 0; j < d.count; j++) {
          const offset = (j - (d.count - 1) / 2) * 0.35;
          stacks.push({
            id: `preview_${di}_${j}`,
            troopType: d.troop_type,
            x: d.position.x + offset,
            y: d.position.y,
            unitMaxHp: unitHp,
            totalUnits: 1,
            hp: unitHp,
            alive: true,
          });
        }
      });
      return stacks;
    },
    [deployments, catalog]
  );

  const selectedCells = useMemo(() => deployments.map((d) => d.position), [deployments]);

  const nowDeploying = battle.phase === "deploying";
  const nowViewing = battle.phase === "viewing";
  const catalogIds = useMemo(() => catalog.map((t) => t.id), [catalog]);

  return (
    <div className="absolute inset-0 z-50 bg-black/95">
      {/* Phaser canvas — mounted for the entire battle lifecycle */}
      <BattleCanvas
        layout={battle.defenderLayout}
        buildingHpByKey={buildingHpByKey}
        stacks={nowViewing ? viewingStacks : previewStacks}
        deploymentZone={nowDeploying ? deploymentZone : null}
        selectedCells={nowDeploying ? selectedCells : []}
        troopCatalogIds={catalogIds}
        interactiveDeployment={nowDeploying}
        staticPreview
        selectedTroop={nowDeploying ? selectedTroop : null}
        onCellClick={nowDeploying ? handleCellClick : undefined}
      />

      {/* Deployment UI overlay */}
      {nowDeploying && (
        <DeploymentScreen
          battle={battle}
          socket={socket}
          deployments={deployments}
          selectedTroop={selectedTroop}
          onSelectTroop={setSelectedTroop}
          deployCounts={deployCounts}
          deployError={deployError}
          onIncrementCount={incrementCount}
          onDecrementCount={decrementCount}
          onStartBattle={handleStartBattle}
          onSetDeployError={setDeployError}
        />
      )}

      {/* Playback HUD overlay */}
      {nowViewing && (
        <div className="absolute inset-0 flex flex-col pointer-events-none">
          <div className="flex items-center justify-between p-4 shrink-0 border-b border-purgatory-border bg-black/80 pointer-events-auto">
            <div>
              <h2 className="font-serif text-lg font-bold tracking-wider text-red-400">Battle in Progress</h2>
              <div className="text-[10px] uppercase tracking-widest text-gray-500">vs {battle.defenderName}</div>
            </div>
            <div className="flex items-center gap-4">
              <div className="text-right">
                <div className="text-[10px] uppercase tracking-wider text-gray-500">Elapsed</div>
                <div className="text-lg font-bold font-sans text-gray-300">{(currentTick / 10).toFixed(1)}s</div>
              </div>
              <div className="flex items-center gap-2 text-[10px] uppercase tracking-widest text-red-400">
                <span className="w-2 h-2 rounded-full bg-red-500 animate-pulse" />Live
              </div>
            </div>
          </div>
          <div className="mt-auto flex items-center justify-between p-4 bg-purgatory-card/80 border-t border-purgatory-border shrink-0 pointer-events-auto">
            <div className="text-sm text-gray-500">
              {socket.error ? (
                <span className="text-red-400">{socket.error}</span>
              ) : resultLoading ? (
                <span className="text-amber-400 animate-pulse">Ending battle...</span>
              ) : (
                `${viewingStacks.filter((s) => s.alive).length} of ${viewingStacks.length} troops still standing`
              )}
            </div>
            <button
              type="button"
              onClick={handleEndBattle}
              className="text-xs uppercase tracking-widest font-bold px-6 py-2 rounded border border-red-600/60 text-red-400 hover:bg-red-900/30 hover:border-red-500 transition-all"
            >
              End Battle
            </button>
          </div>
        </div>
      )}

      {battle.phase === "result" && <BattleResultScreen battle={battle} />}
    </div>
  );
}
