import { useEffect, useState } from "react";
import { useGame } from "../../hooks/useGame";
import { phaserEvents } from "../../game/phaser/events";
import * as baseApi from "../../api/endpoints/base";
import type { ShopItem, PlacedBuilding } from "../../types/building";

function formatTime(seconds: number): string {
  if (seconds < 60) return `${seconds}s`;
  const mins = Math.floor(seconds / 60);
  if (mins < 60) return `${mins}m`;
  const hours = Math.floor(mins / 60);
  const remMins = mins % 60;
  if (hours < 24) return remMins > 0 ? `${hours}h ${remMins}m` : `${hours}h`;
  const days = Math.floor(hours / 24);
  const remHours = hours % 24;
  return remHours > 0 ? `${days}d ${remHours}h` : `${days}d`;
}

function UpgradeTimer({ endsAt }: { endsAt: string }) {
  const endMs = new Date(endsAt).getTime();
  const [remaining, setRemaining] = useState(() => Date.now())

  useEffect(() => {
    const interval = setInterval(() => {
      setRemaining(Math.max(0, endMs - Date.now()))
    })
    return () => clearInterval(interval);
  }, [endMs]);

  const seconds = Math.ceil(remaining / 1000);
  const minutes = Math.floor(seconds / 60);
  const secs = seconds % 60;

  return (
    <div className="mb-3 text-[10px] text-teal-400/80">
      <div>Upgrading...</div>
      {seconds > 0 ? (
        <div className="font-bold text-teal-300">
          {minutes}:{String(secs).padStart(2, "0")} remaining
        </div>
      ) : (
        <div className="font-bold text-amber-400">Completing...</div>
      )}
    </div>
  );
}

interface PlacementToolbarProps {
  buildingMenu: PlacedBuilding | null;
  onCloseMenu: () => void;
  onLayoutChanged: () => void;
}

export default function PlacementToolbar({ buildingMenu, onCloseMenu, onLayoutChanged }: PlacementToolbarProps) {
  const { state, api } = useGame();
  const [inventoryOpen, setInventoryOpen] = useState(false);

  const layout = state.layout

  const inventory = state.inventory ?? new Map<ShopItem, number>();

  const cancelMode = () => {
    phaserEvents.mode = "none";
    phaserEvents.placementBuilding = null;
    phaserEvents.onGridClick = null;
    phaserEvents.ghostPosition = null;
  };

  const closeAndCancel = () => {
    cancelMode();
    onCloseMenu();
  };

  const enterPlaceMode = (item: { id: string; size: number }) => {
    cancelMode();
    phaserEvents.mode = "place";
    phaserEvents.placementBuilding = { id: item.id, size: item.size };
    setInventoryOpen(false);

    phaserEvents.onGridClick = async () => {
      const pos = phaserEvents.ghostPosition;
      if (!pos) return;

      const occupied = layout?.buildings.some(b => {
        return pos.x < b.x + b.size && pos.x + item.size > b.x &&
          pos.y < b.y + b.size && pos.y + item.size > b.y;
      });
      if (occupied) return;

      const res = await baseApi.placeBuilding(api, item.id, pos.x, pos.y);
      if (res.success) {
        if (state.inventory) for (const b of state.inventory.keys()) {
          if (b.building.id === item.id) {
            if (state.inventory.get(b) === 1) state.inventory.delete(b);
            else state.inventory.set(b, (state.inventory.get(b) ?? 1) - 1);
          }
        }
        await onLayoutChanged();
        cancelMode();
      }
    };
  };

  const enterMoveMode = () => {
    cancelMode();
    if (!buildingMenu) return;
    phaserEvents.mode = "place";
    phaserEvents.placementBuilding = { id: buildingMenu.building_id, size: buildingMenu.size };

    phaserEvents.onGridClick = async () => {
      const pos = phaserEvents.ghostPosition;
      if (!pos) return;

      const occupied = layout?.buildings.some(b => {
        if (b.x === buildingMenu.x && b.y === buildingMenu.y) return false;
        return pos.x < b.x + b.size && pos.x + buildingMenu.size > b.x &&
          pos.y < b.y + b.size && pos.y + buildingMenu.size > b.y;
      });
      if (occupied) return;

      const res = await baseApi.moveBuilding(api, buildingMenu.building_id, buildingMenu.x, buildingMenu.y, pos.x, pos.y);
      if (res.success) {
        await onLayoutChanged();
        cancelMode();
        onCloseMenu();
      }
    };
  };

  const handleRemove = async () => {
    if (!buildingMenu) return;
    const res = await baseApi.removeBuilding(api, buildingMenu.building_id, buildingMenu.x, buildingMenu.y);
    if (res.success) {
      await onLayoutChanged();
      onCloseMenu();
    }
  };

  const handleUpgrade = async () => {
    if (!buildingMenu) return;
    const res = await baseApi.upgradeBuilding(api, buildingMenu.building_id, buildingMenu.x, buildingMenu.y);
    if (res.success) {
      await onLayoutChanged();
      onCloseMenu();
    }
  };

  const activeMode = phaserEvents.mode !== "none";

  return (
    <>
      {/* Bottom toolbar */}
      <div className="absolute bottom-4 left-1/2 -translate-x-1/2 z-30 pointer-events-auto flex flex-col items-center gap-2">
        {/* Inventory tray */}
        {inventoryOpen && !activeMode && (
          <div className="bg-purgatory-card/95 backdrop-blur-md border border-purgatory-border rounded-lg p-3 shadow-2xl flex gap-2 flex-wrap max-w-lg">
            {inventory.size === 0 && (
              <div className="text-xs text-gray-500 px-2 py-1">No buildings to place</div>
            )}
            {[...inventory].map(([item, count]) => (
              <button
                key={item.building.id}
                onClick={() => enterPlaceMode(item.building)}
                className="flex items-center gap-2 bg-purgatory-input border border-purgatory-border hover:border-amber-500/40 rounded px-3 py-2 transition-all hover:bg-amber-500/5"
              >
                <div className="text-xs text-gray-300 font-bold">{item.building.name}</div>
                <div className="text-[9px] text-gray-500">
                  {count}×
                </div>
              </button>
            ))}
          </div>
        )}

        <div className="flex gap-2 bg-purgatory-card/95 backdrop-blur-md border border-purgatory-border rounded-lg p-2 shadow-2xl">
          <button
            onClick={() => { cancelMode(); setInventoryOpen(!inventoryOpen); }}
            className={`px-3 py-1.5 rounded text-xs uppercase tracking-widest font-bold transition-all ${inventoryOpen
              ? "bg-amber-500/10 border border-amber-500/40 text-amber-400"
              : "border border-purgatory-border text-gray-400 hover:text-gray-200 hover:border-gray-600"
              }`}
          >
            {inventory.size > 0 ? `Inventory (${inventory.size})` : "Inventory"}
          </button>

          {activeMode && (
            <button
              onClick={cancelMode}
              className="px-4 py-1.5 border border-red-900/50 text-red-400 rounded text-xs uppercase tracking-widest font-bold hover:bg-red-950/30 transition-all"
            >
              Cancel
            </button>
          )}
        </div>
      </div>

      {/* Building context menu */}
      {buildingMenu && phaserEvents.mode === "none" && !inventoryOpen && (
        <div className="absolute bottom-24 left-1/2 -translate-x-1/2 z-30 pointer-events-auto">
          <div className="bg-purgatory-card border border-purgatory-border rounded-lg p-4 shadow-2xl w-72">
            <div className="flex items-center justify-between mb-3">
              <div>
                <div className="text-gray-200 font-bold text-sm">{buildingMenu.name}</div>
                <div className="text-[10px] uppercase tracking-widest text-gray-500">
                  Lv.{buildingMenu.level} · ({buildingMenu.x}, {buildingMenu.y})
                </div>
              </div>
              <button onClick={closeAndCancel} className="text-gray-500 hover:text-gray-300 p-1">
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            {buildingMenu.hp != null && (
              <div className="mb-3 text-[10px] text-gray-400">
                HP: {buildingMenu.hp}
                {buildingMenu.dps != null && ` · DPS: ${buildingMenu.dps}`}
                {buildingMenu.attack_range != null && ` · Range: ${buildingMenu.attack_range}`}
              </div>
            )}

            {!buildingMenu.metadata?.upgrade_ends_at && buildingMenu.upgrade_cost != null && (
              <div className="mb-3 text-[10px] text-amber-500/70">
                Next upgrade: {buildingMenu.upgrade_cost} penitence · {formatTime(buildingMenu.upgrade_time ?? 0)}
              </div>
            )}

            {!buildingMenu.metadata?.upgrade_ends_at && buildingMenu.upgrade_cost == null && (
              <div className="mb-3 text-[10px] text-gray-500">MAX LEVEL — no further upgrades</div>
            )}

            {buildingMenu.metadata?.upgrade_ends_at && (
              <UpgradeTimer endsAt={buildingMenu.metadata.upgrade_ends_at} />
            )}

            <div className="flex gap-2">
              <button onClick={enterMoveMode} className="flex-1 px-3 py-1.5 border border-blue-900/40 text-blue-400 rounded text-xs uppercase tracking-wider font-bold hover:bg-blue-950/20 transition-all">
                Move
              </button>
              <button onClick={handleRemove} className="flex-1 px-3 py-1.5 border border-red-900/40 text-red-400 rounded text-xs uppercase tracking-wider font-bold hover:bg-red-950/20 transition-all">
                Remove
              </button>
              {!buildingMenu.metadata?.upgrade_ends_at && (
                <button
                  onClick={handleUpgrade}
                  className="flex-1 px-3 py-1.5 border border-amber-500/40 text-amber-400 rounded text-xs uppercase tracking-wider font-bold hover:bg-amber-500/10 transition-all disabled:opacity-30 disabled:cursor-not-allowed"
                  disabled={buildingMenu.upgrade_cost == null || (state.economy?.penitence ?? 0) < buildingMenu.upgrade_cost}
                >
                  Upgrade
                </button>
              )}
            </div>
          </div>
        </div>
      )}
    </>
  );
}
