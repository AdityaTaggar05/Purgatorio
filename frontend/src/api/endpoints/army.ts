import type { ApiClient } from "../client";
import type { ApiResponse } from "../../types/api";
import type { Troop, ArmyResponse } from "../../types/army";

export function getCatalog(api: ApiClient): Promise<ApiResponse<{ troops: Troop[] }>> {
  return api.get<{ troops: Troop[] }>("/army/troops");
}

export function getMyTroops(api: ApiClient): Promise<ApiResponse<ArmyResponse>> {
  return api.get<ArmyResponse>("/army/my-troops");
}

export function trainTroop(api: ApiClient, troopId: string, count: number): Promise<ApiResponse<ArmyResponse>> {
  return api.post<ArmyResponse>("/army/train", { troop_id: troopId, count });
}

export function detrainTroop(api: ApiClient, troopId: string, count: number): Promise<ApiResponse<ArmyResponse>> {
  return api.post<ArmyResponse>("/army/detrain", { troop_id: troopId, count });
}
