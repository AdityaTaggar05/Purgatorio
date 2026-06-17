import { useState } from "react";
import EntityBar from "./EntityBar";

export interface TroopStack {
  id: string;
  troopType: string;
  x: number;
  y: number;
  unitMaxHp: number;
  totalUnits: number;
  hp: number;
  alive: boolean;
}

interface TroopStackOverlayProps {
  stacks: TroopStack[];
  cellSize: number;
  staticPreview?: boolean;
}

function TroopIcon({ troopType, size }: { troopType: string; size: number }) {
  const [failed, setFailed] = useState(false);

  if (failed) {
    return (
      <div
        className="flex items-center justify-center rounded-full bg-red-900/70 border border-red-500/50 text-[8px] font-bold text-red-200 uppercase"
        style={{ width: size, height: size }}
      >
        {troopType.slice(0, 2)}
      </div>
    );
  }

  return (
    <img
      src={`/assets/${troopType}.png`}
      alt={troopType}
      onError={() => setFailed(true)}
      className="object-contain drop-shadow-md"
      style={{ width: size, height: size }}
      draggable={false}
    />
  );
}

export default function TroopStackOverlay({ stacks, cellSize, staticPreview = false }: TroopStackOverlayProps) {
  const iconSize = cellSize * 1.8;

  return (
    <div className="absolute inset-0 pointer-events-none">
      {stacks.map((s) => {
        const maxHp = s.unitMaxHp * s.totalUnits;
        const displayCount = staticPreview
          ? s.totalUnits
          : Math.max(0, Math.min(s.totalUnits, Math.ceil((s.hp / Math.max(1, s.unitMaxHp)))));

        return (
          <div
            key={s.id}
            className={`absolute flex flex-col items-center transition-all duration-500 ${
              s.alive ? "opacity-100 scale-100" : "opacity-0 scale-50"
            }`}
            style={{
              left: s.x * cellSize + cellSize / 2 - iconSize / 2,
              top: s.y * cellSize - iconSize * 0.55,
            }}
          >
            {!staticPreview && (
              <div className="mb-0.5">
                <EntityBar hp={s.hp} maxHp={maxHp} width={iconSize} />
              </div>
            )}
            <div className="relative">
              <TroopIcon troopType={s.troopType} size={iconSize} />
              {displayCount > 1 && (
                <div className="absolute -bottom-1 -right-1 bg-black/80 border border-purgatory-border text-[8px] text-gray-200 font-bold rounded-full w-4 h-4 flex items-center justify-center">
                  {displayCount}
                </div>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}
