import type { ApiClient } from "../client";
import type { ApiResponse } from "../../types/api";
import type { ArmyResponse, Troop } from "../../types/army";

export function getCatalog(api: ApiClient): Promise<ApiResponse<{ catalog: Troop[]; army: ArmyResponse }>> {
  return api.get<{ catalog: Troop[]; army: ArmyResponse }>("/army");
}

export function trainTroop(api: ApiClient, troopId: string, count: number): Promise<ApiResponse<ArmyResponse>> {
  return api.post<ArmyResponse>("/army/train", { troop_id: troopId, count });
}
