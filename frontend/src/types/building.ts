export interface PlacedBuilding {
  building_id: string;
  name: string;
  category: "defense" | "resource" | "army" | "other";
  level: number;
  x: number;
  y: number;
  size: number;
  hp?: number;
  dps?: number;
  attack_range?: number;
  production_rate?: number;
  storage_capacity?: number;
  metadata?: BuildingMetadata;
  upgrade_cost?: number;
  upgrade_time?: number;
}

export interface BuildingMetadata {
  upgrade_ends_at?: string;
}

export interface BaseLayout {
  buildings: PlacedBuilding[];
  grid_w: number;
  grid_h: number;
}

export interface BuildingInfo {
  id: string;
  name: string;
  size: number;
  price: number;
  currency: "penitence" | "grace";
  category: "defense" | "resource" | "army" | "other";
}

export interface ShopItem {
  building: BuildingInfo;
  current_owned: number;
  max_allowed: number;
  can_buy: boolean;
}
