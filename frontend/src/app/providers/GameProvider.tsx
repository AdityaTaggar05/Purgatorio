import { useEffect, useReducer, useRef, type ReactNode } from "react";
import { GameContext, gameReducer, type GameState, type GameContextType } from "./GameContext";
import { ApiClient } from "../../api/client";
import { useAuth } from "../../hooks/useAuth";
import { API_BASE_URL } from "../../config";
import type { Troop } from "../../types/army";
import * as baseApi from "../../api/endpoints/base";
import * as shopApi from "../../api/endpoints/shop";
import type { ShopItem } from "../../types/shop";

const initialState: GameState = {
  economy: null,
  layout: null,
  inventory: null,
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

  const layoutRef = useRef<GameState["layout"]>(null);
  layoutRef.current = state.layout;

  useEffect(() => {
    if (!accessToken) return;
    let cancelled = false;
    let interval: number

    async function hydrate() {
      dispatch({ type: "SET_CHECK_IN_RESULT", payload: null });

      const checkInRes = await api.post<{ completed_upgrades: { building_id: string; x: number; y: number; from_level: number; to_level: number }[] }>("/base/check-in");

      const [economyRes, layoutRes, armyRes, combatRes] = await Promise.all([
        api.get<{ penitence: number; grace: number; max_penitence: number; overflow_penitence?: number }>("/user/economy"),
        api.get<{ buildings: unknown[]; grid_w: number; grid_h: number }>("/base/layout"),
        api.get<{ troops: Record<string, number>; used_capacity: number; max_capacity: number }>("/army/my-troops"),
        api.get<{ sin_meter: number }>("/user/combat"),
      ]);

      if (cancelled) return;

      if (economyRes.success) {
        dispatch({ type: "SET_ECONOMY", payload: economyRes.data });
      }
      if (layoutRes.success) {
        const layout = layoutRes.data as GameState["layout"]
        dispatch({ type: "SET_LAYOUT", payload: layout });

        interval = setInterval(async () => {
          const layout = layoutRef.current;
          if (!layout) return;

          const upgrades = layout.buildings
            .filter((b) => b.metadata?.upgrade_ends_at)
            .filter((b) => new Date(b.metadata!.upgrade_ends_at!).getTime() <= Date.now());

          if (upgrades.length === 0) return;

          const res = await baseApi.checkIn(api);
          if (!res.success || res.data.completed_upgrades.length === 0) return;

          const [layoutRes2, economyRes2] = await Promise.all([
            api.get<{ buildings: unknown[]; grid_w: number; grid_h: number }>("/base/layout"),
            api.get<{ penitence: number; grace: number; max_penitence: number; overflow_penitence?: number }>("/user/economy"),
          ]);

          if (layoutRes2.success) {
            dispatch({ type: "SET_LAYOUT", payload: layoutRes2.data as GameState["layout"] });
          }
          if (economyRes2.success) {
            dispatch({ type: "SET_ECONOMY", payload: economyRes2.data });
          }

          const names = res.data.completed_upgrades
            .map(u => `${u.building_id} Lv.${u.from_level} → ${u.to_level}`)
            .join(", ");
          dispatch({ type: "SET_CHECK_IN_RESULT", payload: `Upgrades completed: ${names}` });
        }, 1500)

        shopApi.getShop(api).then((res) => {
          const placedCounts = new Map<string, number>();

          layout?.buildings.forEach(b => {
            placedCounts.set(b.building_id, (placedCounts.get(b.building_id) ?? 0) + 1);
          });

          const ownedItems = res.data.items.filter(item => item.current_owned > 0);
          const placeableItems = ownedItems.filter(item => item.current_owned > (placedCounts.get(item.building.id) ?? 0));

          const inventory = new Map<ShopItem, number>();

          placeableItems.forEach((item) => {
            inventory.set(item, item.current_owned - (placedCounts.get(item.building.id) ?? 0))
          })

          dispatch({ type: "SET_INVENTORY", payload: inventory })
        })
      }
      if (armyRes.success) {
        dispatch({ type: "SET_ARMY", payload: armyRes.data as GameState["army"] });
      }
      if (combatRes.success) {
        dispatch({ type: "SET_SIN_METER", payload: combatRes.data.sin_meter });
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
    return () => { cancelled = true; clearInterval(interval) };
  }, [accessToken, api]);

  return (
    <GameContext.Provider value={{ state, api, dispatch }}>
      {children}
    </GameContext.Provider>
  );
}
