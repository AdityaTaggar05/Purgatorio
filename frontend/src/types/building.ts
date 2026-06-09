export interface BuildingData {
  id: string;
  x: number;
  y: number;
  size: number;
  metadata?: object;
}

export interface BaseLayout {
  userID: string;
  tiles: number;
  subgridFactor: number;
  buildings: BuildingData[];
}
