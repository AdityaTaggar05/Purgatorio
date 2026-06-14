export interface Troop {
  id: string;
  name: string;
  training_cost: number;
  space: number;
  hp: number;
  dps: number;
  speed: number;
  attack_range: number;
  preferred_target: string;
}

export interface ArmyResponse {
  troops: Record<string, number>;
  used_capacity: number;
  max_capacity: number;
}
