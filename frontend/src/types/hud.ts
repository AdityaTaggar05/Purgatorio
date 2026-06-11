export interface ResourceMeter {
  current: number;
  max: number;
  label: string;
}

export interface HudProps {
  username: string;
  level: number;
  penitence: ResourceMeter;
  grace: ResourceMeter;
  sinMeter: number;
  onAscensionClick: () => void;
  onLogoutClick: () => void;
  onAttackClick: () => void;
  onArmyClick: () => void;
}
