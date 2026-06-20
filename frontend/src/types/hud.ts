import type { UserEconomy } from "./economy";

export interface HudProps {
  username: string;
  level: number;
  economy: UserEconomy | null;
  sinMeter: number;
  onAscensionClick: () => void;
  onLogoutClick: () => void;
  onAttackClick: () => void;
  onArmyClick: () => void;
  onCollectClick: () => void;
}
