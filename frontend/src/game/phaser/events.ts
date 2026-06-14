import type { PlacedBuilding } from "../../../types/building";

export type PlacementMode = "none" | "place" | "move";

export const phaserEvents = {
  mode: "none" as PlacementMode,
  placementBuilding: null as { id: string; size: number } | null,
  onGridClick: null as ((x: number, y: number) => void) | null,
  onBuildingClick: null as ((building: PlacedBuilding) => void) | null,
  ghostPosition: null as { x: number; y: number } | null,
};
