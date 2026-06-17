import { useEffect, useMemo, useRef, useState } from "react";
import { useGame } from "../../hooks/useGame";
import * as battleApi from "../../api/endpoints/battle";
import BattleGrid from "../battle/BattleGrid";
import TroopStackOverlay, { type TroopStack } from "../battle/TroopStackOverlay";
import type { ActiveBattle } from "../../app/providers/GameContext";
import type { UseBattleSocketResult } from "../../hooks/useBattleSocket";
import type { HpChange, TickResult } from "../../types/battle";

interface BattleViewerProps {
  battle: ActiveBattle;
  socket: UseBattleSocketResult;
}

const CELL_SIZE = 14;
const GRID_SIZE = 30;
const TICK_INTERVAL_MS = 100;

// Convention from the battle backend: "troop_{troopType}_{index}". None of
function parseTroopType(entityId: string): string | null {
  const parts = entityId.split("_");
  if (parts.length < 3 || parts[0] !== "troop") return null;
  return parts.slice(1, -1).join("_");
}

function applyDamageToType(stacks: TroopStack[], troopType: string, delta: number): TroopStack[] {
  let targetIndex = -1;
  let lowestFraction = Infinity;
  stacks.forEach((s, i) => {
    if (s.troopType !== troopType || !s.alive) return;
    const fraction = s.hp / Math.max(1, s.unitMaxHp * s.totalUnits);
    if (fraction < lowestFraction) {
      lowestFraction = fraction;
      targetIndex = i;
    }
  });
  if (targetIndex === -1) return stacks;

  return stacks.map((s, i) => {
    if (i !== targetIndex) return s;
    const newHp = Math.max(0, s.hp + delta);
    return { ...s, hp: newHp, alive: newHp > 0 };
  });
}

function applyHpChanges(stacks: TroopStack[], changes: HpChange[]): TroopStack[] {
  let next = stacks;
  for (const change of changes) {
    if (change.entity_type !== "troop") continue;
    const troopType = parseTroopType(change.entity_id);
    if (!troopType) continue;
    next = applyDamageToType(next, troopType, change.delta);
  }
  return next;
}

export default function BattleViewer({ battle, socket }: BattleViewerProps) {
  const { state, api, dispatch } = useGame();
  const catalog = state.troopCatalog ?? [];

  // Computed once, from the deployment that kicked off this battle.
  const initialStacks = useMemo<TroopStack[]>(() => {
    return battle.deployment.map((d, i) => {
      const troopDef = catalog.find((t) => t.id === d.troop_type);
      const unitHp = troopDef?.hp ?? 1;
      return {
        id: `stack_${i}`,
        troopType: d.troop_type,
        x: d.position.x,
        y: d.position.y,
        unitMaxHp: unitHp,
        totalUnits: d.count,
        hp: unitHp * d.count,
        alive: true,
      };
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const [stacks, setStacks] = useState<TroopStack[]>(initialStacks);
  const [currentTick, setCurrentTick] = useState(0);
  const [skipping, setSkipping] = useState(false);
  const [skipError, setSkipError] = useState<string | null>(null);

  const stacksRef = useRef(stacks);
  const cursorRef = useRef(0);
  const ticksRef = useRef(socket.ticks);

  useEffect(() => {
    ticksRef.current = socket.ticks;
  }, [socket.ticks]);

  // Plays back whatever ticks have arrived
  useEffect(() => {
    const interval = setInterval(() => {
      const ticks = ticksRef.current;
      if (cursorRef.current >= ticks.length) return;
      const tick: TickResult = ticks[cursorRef.current];
      cursorRef.current += 1;
      const updated = applyHpChanges(stacksRef.current, tick.hp_changes);
      stacksRef.current = updated;
      setStacks(updated);
      setCurrentTick(tick.tick);
    }, TICK_INTERVAL_MS);
    return () => clearInterval(interval);
  }, []);

  const handleSkip = async () => {
    setSkipping(true);
    setSkipError(null);
    const res = await battleApi.getBattleResult(api, battle.battleId);
    setSkipping(false);
    if (res.success) {
      dispatch({
        type: "SET_ACTIVE_BATTLE",
        payload: {
          ...battle,
          phase: "result",
          outcome: res.data.outcome,
          destruction: res.data.destruction,
          loot: res.data.loot,
          newSinMeter: res.data.sin_meter,
          duration: res.data.duration,
        },
      });
    } else {
      setSkipError(res.error?.message ?? "Battle still resolving — try again in a moment.");
    }
  };

  const elapsedSeconds = (currentTick / 10).toFixed(1);
  const liveStacksCount = stacks.filter((s) => s.alive).length;

  return (
    <div className="absolute inset-0 z-50 bg-black/95 flex flex-col">
      <div className="flex items-center justify-between p-4 shrink-0 border-b border-purgatory-border">
        <div>
          <h2 className="font-serif text-lg font-bold tracking-wider text-red-400">Battle in Progress</h2>
          <div className="text-[10px] uppercase tracking-widest text-gray-500">vs {battle.defenderName}</div>
        </div>
        <div className="flex items-center gap-4">
          <div className="text-right">
            <div className="text-[10px] uppercase tracking-wider text-gray-500">Elapsed</div>
            <div className="text-lg font-bold font-sans text-gray-300">{elapsedSeconds}s</div>
          </div>
          <div className="flex items-center gap-2 text-[10px] uppercase tracking-widest text-red-400">
            <span className="w-2 h-2 rounded-full bg-red-500 animate-pulse" />
            Live
          </div>
        </div>
      </div>

      <div className="flex-1 flex items-center justify-center overflow-auto p-4">
        <div className="relative" style={{ width: GRID_SIZE * CELL_SIZE, height: GRID_SIZE * CELL_SIZE }}>
          <div className="absolute -top-6 left-0 text-[9px] uppercase tracking-widest text-gray-600">
            Defender base layout unavailable — showing troop positions only
          </div>
          <BattleGrid buildings={[]} gridW={GRID_SIZE} gridH={GRID_SIZE} cellSize={CELL_SIZE} />
          <TroopStackOverlay stacks={stacks} cellSize={CELL_SIZE} />
        </div>
      </div>

      <div className="flex items-center justify-between p-4 bg-purgatory-card/80 border-t border-purgatory-border shrink-0">
        <div className="text-sm text-gray-500">
          {socket.error ? (
            <span className="text-red-400">{socket.error} — Skip will try to fetch the final result directly.</span>
          ) : skipError ? (
            <span className="text-red-400">{skipError}</span>
          ) : (
            `${liveStacksCount} of ${stacks.length} troop group${stacks.length === 1 ? "" : "s"} still standing`
          )}
        </div>
        <button
          type="button"
          onClick={handleSkip}
          disabled={skipping}
          className="text-xs uppercase tracking-widest font-bold px-6 py-2 rounded border transition-all
            enabled:border-gray-600 enabled:text-gray-300 enabled:hover:border-gray-400
            disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {skipping ? "Fetching result..." : "Skip to Result"}
        </button>
      </div>
    </div>
  );
}
