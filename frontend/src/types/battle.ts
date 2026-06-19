import type { BaseLayout } from "./building";

export interface MatchListEntry {
  user_id: string;
  username: string;
  terrace_level: number;
}

export interface InitiateResponse {
  battle_id: string;
  defender_name: string;
  defender_layout?: BaseLayout;
}

export interface BattleResultResponse {
  battle_id: string;
  outcome: "victory" | "defeat" | "threshold_failed";
  destruction: number;
  loot: number;
  duration: number;
  sin_meter: number;
}

export interface TroopDeployment {
  troop_type: string;
  position: { x: number; y: number };
  count: number;
}

export interface PositionChange {
  entity_id: string;
  x: number;
  y: number;
}

export interface TickResult {
  tick: number;
  hp_changes: HpChange[];
  positions: PositionChange[];
  done?: boolean;
}

export interface HpChange {
  entity_id: string;
  entity_type: "troop" | "building";
  new_hp: number;
  delta: number;
}

export interface ReplayData {
  battle_id: string;
  attacker_id: string;
  defender_id: string;
  outcome: string;
  data: {
    deployment: TroopDeployment[];
    seed: number;
    base_snapshot_id: string;
  };
}
