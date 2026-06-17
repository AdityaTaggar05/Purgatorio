import { useEffect, useReducer, useRef, type ReactNode } from "react";
import { GameContext, gameReducer, type GameState, type GameContextType } from "./GameContext";
import { ApiClient } from "../../api/client";
import { useAuth } from "../../hooks/useAuth";
import { API_BASE_URL } from "../../config";
import type { Troop } from "../../types/army";

const initialState: GameState = {
  economy: null,
  layout: null,
  troopCatalog: null,
  army: null,
  sinMeter: 0,
  isLoading: true,
  error: null,
  checkInResult: null,
  activeBattle: null,
};

let cachedTroopCatalog: Troop[] | null = null;

export function GameProvider({ children }: { children: ReactNode }) {
  const { accessToken, getFreshToken, logout } = useAuth();
  const [state, dispatch] = useReducer(gameReducer, initialState);

  const apiRef = useRef<GameContextType["api"] | null>(null);
  if (!apiRef.current) {
    apiRef.current = new ApiClient({
      baseUrl: API_BASE_URL,
      getToken: () => accessToken,
      onTokenRefresh: getFreshToken,
      onAuthFailure: logout,
    });
  }
  const api = apiRef.current;

  useEffect(() => {
    if (!accessToken) return;
    let cancelled = false;

    async function hydrate() {
      dispatch({ type: "SET_CHECK_IN_RESULT", payload: null });

      const checkInRes = await api.post<{ completed_upgrades: { building_id: string; x: number; y: number; from_level: number; to_level: number }[] }>("/base/check-in");

      const [economyRes, layoutRes, armyRes] = await Promise.all([
        api.get<{ penitence: number; grace: number; max_penitence: number; overflow_penitence?: number }>("/user/economy"),
        api.get<{ buildings: unknown[]; grid_w: number; grid_h: number }>("/base/layout"),
        api.get<{ troops: Record<string, number>; used_capacity: number; max_capacity: number }>("/army/my-troops"),
      ]);

      if (cancelled) return;

      if (economyRes.success) {
        dispatch({ type: "SET_ECONOMY", payload: economyRes.data });
      }
      if (layoutRes.success) {
        dispatch({ type: "SET_LAYOUT", payload: layoutRes.data as GameState["layout"] });
      }
      if (armyRes.success) {
        dispatch({ type: "SET_ARMY", payload: armyRes.data as GameState["army"] });
      }

      if (checkInRes.success && checkInRes.data.completed_upgrades.length > 0) {
        const names = checkInRes.data.completed_upgrades
          .map(u => `${u.building_id} Lv.${u.from_level} → ${u.to_level}`)
          .join(", ");
        dispatch({ type: "SET_CHECK_IN_RESULT", payload: `Upgrades completed: ${names}` });
      }

      if (!cachedTroopCatalog) {
        const catalogRes = await api.get<{ troops: Troop[] }>("/army/troops");
        if (!cancelled && catalogRes.success) {
          cachedTroopCatalog = catalogRes.data.troops;
        }
      }
      dispatch({ type: "SET_TROOP_CATALOG", payload: cachedTroopCatalog });

      dispatch({ type: "SET_LOADING", payload: false });
    }

    hydrate();
    return () => { cancelled = true; };
  }, [accessToken]);

  return (
    <GameContext.Provider value={{ state, api, dispatch }}>
      {children}
    </GameContext.Provider>
  );
}
