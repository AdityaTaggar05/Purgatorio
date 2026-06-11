export interface BuildingData {
  id: string;
  x: number;
  y: number;
  size: number;
  metadata?: object;
}

export interface BaseLayout {
  user_id: string;
  tiles: number;
  subgrid_factor: number;
  buildings: BuildingData[];
}
