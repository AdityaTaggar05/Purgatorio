import type { PlacedBuilding } from "../../types/building";

interface BattleGridProps {
  buildings: PlacedBuilding[];
  gridW?: number;
  gridH?: number;
  deploymentZone?: { x: number; y: number }[] | null;
  selectedCells?: { x: number; y: number }[];
  onCellClick?: (x: number, y: number) => void;
  cellSize?: number;
}

export default function BattleGrid({
  buildings,
  gridW = 30,
  gridH = 30,
  deploymentZone,
  selectedCells,
  onCellClick,
  cellSize = 14,
}: BattleGridProps) {
  const occupancy = new Map<string, PlacedBuilding>();
  for (const b of buildings) {
    for (let dx = 0; dx < b.size; dx++) {
      for (let dy = 0; dy < b.size; dy++) {
        occupancy.set(`${b.x + dx},${b.y + dy}`, b);
      }
    }
  }

  const selectedSet = new Set(selectedCells?.map((c) => `${c.x},${c.y}`) ?? []);
  const deploySet = new Set(deploymentZone?.map((c) => `${c.x},${c.y}`) ?? []);

  const cellBlocks: { building: PlacedBuilding; topX: number; topY: number; w: number; h: number }[] = [];
  const seen = new Set<string>();
  for (const b of buildings) {
    const key = `${b.x},${b.y}`;
    if (seen.has(key)) continue;
    seen.add(key);
    let w = 1;
    let h = 1;
    while (occupancy.has(`${b.x + w},${b.y}`) && occupancy.get(`${b.x + w},${b.y}`)?.building_id === b.building_id) w++;
    while (occupancy.has(`${b.x},${b.y + h}`) && occupancy.get(`${b.x},${b.y + h}`)?.building_id === b.building_id) h++;
    cellBlocks.push({ building: b, topX: b.x, topY: b.y, w, h });
  }

  const categoryColor = (cat: string) => {
    switch (cat) {
      case "defense": return "bg-red-800/60 border-red-600/40";
      case "resource": return "bg-amber-800/60 border-amber-600/40";
      case "army": return "bg-teal-800/60 border-teal-600/40";
      default: return "bg-gray-700/60 border-gray-600/40";
    }
  };

  return (
    <div 
      className="relative inline-block bg-black/60 border border-purgatory-border rounded overflow-hidden"
      style={{ width: gridW * cellSize, height: gridH * cellSize }}
    >
      <div className="absolute inset-0 grid" style={{
        gridTemplateColumns: `repeat(${gridW}, ${cellSize}px)`,
        gridTemplateRows: `repeat(${gridH}, ${cellSize}px)`,
      }}>
        {Array.from({ length: gridW * gridH }).map((_, i) => {
          const x = i % gridW;
          const y = Math.floor(i / gridW);
          const key = `${x},${y}`;
          const isDeployZone = deploySet.has(key);
          const isSelected = selectedSet.has(key);
          const bldg = occupancy.get(key);

          let cellClass = "border border-gray-800/30";
          if (isDeployZone && !isSelected) cellClass = "border border-green-800/40 bg-green-900/10";
          if (isDeployZone && isSelected) cellClass = "border border-green-500/60 bg-green-700/40";
          if (bldg) cellClass += " cursor-not-allowed";

          return (
            <div
              key={i}
              className={cellClass + (onCellClick && !bldg ? " cursor-pointer hover:bg-white/10" : "")}
              style={{ width: cellSize, height: cellSize }}
              onClick={() => { if (!bldg) onCellClick?.(x, y); }}
            />
          );
        })}
      </div>

      {cellBlocks.map((block, i) => (
        <div
          key={i}
          className={`absolute border ${categoryColor(block.building.category)} flex items-center justify-center`}
          style={{
            left: block.topX * cellSize,
            top: block.topY * cellSize,
            width: block.w * cellSize,
            height: block.h * cellSize,
          }}
        >
          <span className="text-white/70 text-[6px] uppercase tracking-tight font-bold truncate px-0.5 leading-tight">
            {block.building.name.split(" ")[0]}
          </span>
        </div>
      ))}
    </div>
  );
}
