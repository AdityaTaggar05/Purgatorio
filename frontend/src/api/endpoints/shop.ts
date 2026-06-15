import type { ApiClient } from "../client";
import type { ApiResponse } from "../../types/api";
import type { ShopItem } from "../../types/building";

export function getShop(api: ApiClient): Promise<ApiResponse<{ items: ShopItem[] }>> {
  return api.get<{ items: ShopItem[] }>("/shop");
}

export function buyBuilding(api: ApiClient, buildingId: string): Promise<ApiResponse<void>> {
  return api.post<void>("/shop/buy", { building_id: buildingId });
}
