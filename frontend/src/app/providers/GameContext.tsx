import { createContext } from "react";
import type { ApiClient } from "../../api/client";
import type { UserEconomy } from "../../types/economy";
import type { BaseLayout } from "../../types/building";
import type { Troop, ArmyResponse } from "../../types/army";
import type { TroopDeployment } from "../../types/battle";

export type BattlePhase = "matching" | "deploying" | "viewing" | "result";

export interface ActiveBattle {
  battleId: string;
  defenderName: string;
  defenderTerraceLevel: number;
  phase: BattlePhase;
  deployment: TroopDeployment[];
  defenderLayout: BaseLayout | null;
  outcome?: "victory" | "defeat" | "threshold_failed";
  destruction?: number;
  loot?: number;
  newSinMeter?: number;
  duration?: number;
}

export interface GameState {
  economy: UserEconomy | null;
  layout: BaseLayout | null;
  troopCatalog: Troop[] | null;
  army: ArmyResponse | null;
  sinMeter: number;
  isLoading: boolean;
  error: string | null;
  checkInResult: string | null;
  activeBattle: ActiveBattle | null;
}

export interface GameContextType {
  state: GameState;
  api: ApiClient;
  dispatch: React.Dispatch<GameAction>;
}

export type GameAction =
  | { type: "SET_ECONOMY"; payload: UserEconomy }
  | { type: "SET_LAYOUT"; payload: BaseLayout | null }
  | { type: "SET_TROOP_CATALOG"; payload: Troop[] | null }
  | { type: "SET_ARMY"; payload: ArmyResponse | null }
  | { type: "SET_SIN_METER"; payload: number }
  | { type: "SET_LOADING"; payload: boolean }
  | { type: "SET_ERROR"; payload: string | null }
  | { type: "SET_CHECK_IN_RESULT"; payload: string | null }
  | { type: "SET_ACTIVE_BATTLE"; payload: ActiveBattle | null };

export function gameReducer(state: GameState, action: GameAction): GameState {
  switch (action.type) {
    case "SET_ECONOMY":
      return {
        ...state,
        economy: {
          ...action.payload,
          overflow_penitence: action.payload.overflow_penitence ?? state.economy?.overflow_penitence ?? 0,
        },
      };
    case "SET_LAYOUT":
      return { ...state, layout: action.payload };
    case "SET_TROOP_CATALOG":
      return { ...state, troopCatalog: action.payload };
    case "SET_ARMY":
      return { ...state, army: action.payload };
    case "SET_SIN_METER":
      return { ...state, sinMeter: action.payload };
    case "SET_LOADING":
      return { ...state, isLoading: action.payload };
    case "SET_ERROR":
      return { ...state, error: action.payload };
    case "SET_CHECK_IN_RESULT":
      return { ...state, checkInResult: action.payload };
    case "SET_ACTIVE_BATTLE":
      return { ...state, activeBattle: action.payload };
    default:
      return state;
  }
}

export const GameContext = createContext<GameContextType | undefined>(undefined);
