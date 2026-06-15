import { createContext } from "react";
import type { ApiClient } from "../../api/client";
import type { UserEconomy } from "../../types/economy";
import type { BaseLayout } from "../../types/building";
import type { Troop, ArmyResponse } from "../../types/army";

export interface GameState {
  economy: UserEconomy | null;
  layout: BaseLayout | null;
  troopCatalog: Troop[] | null;
  army: ArmyResponse | null;
  sinMeter: number;
  isLoading: boolean;
  error: string | null;
}

export interface GameContextType {
  state: GameState;
  api: ApiClient;
  dispatch: React.Dispatch<GameAction>;
}

export type GameAction =
  | { type: "SET_ECONOMY"; payload: UserEconomy }
  | { type: "SET_LAYOUT"; payload: BaseLayout }
  | { type: "SET_TROOP_CATALOG"; payload: Troop[] }
  | { type: "SET_ARMY"; payload: ArmyResponse }
  | { type: "SET_SIN_METER"; payload: number }
  | { type: "SET_LOADING"; payload: boolean }
  | { type: "SET_ERROR"; payload: string | null };

export function gameReducer(state: GameState, action: GameAction): GameState {
  switch (action.type) {
    case "SET_ECONOMY":
      return { ...state, economy: action.payload };
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
    default:
      return state;
  }
}

export const GameContext = createContext<GameContextType | undefined>(undefined);
