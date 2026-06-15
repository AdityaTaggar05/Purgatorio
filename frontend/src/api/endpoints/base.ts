import type { ApiClient } from "../client";
import type { ApiResponse } from "../../types/api";
import type { BaseLayout } from "../../types/building";

export function getLayout(api: ApiClient): Promise<ApiResponse<BaseLayout>> {
  return api.get<BaseLayout>("/base/layout");
}

export function placeBuilding(
  api: ApiClient,
  buildingId: string,
  x: number,
  y: number
): Promise<ApiResponse<BaseLayout>> {
  return api.post<BaseLayout>("/base/layout", { building_id: buildingId, x, y });
}

export function moveBuilding(
  api: ApiClient,
  buildingId: string,
  fromX: number,
  fromY: number,
  toX: number,
  toY: number
): Promise<ApiResponse<BaseLayout>> {
  return api.put<BaseLayout>("/base/layout", {
    building_id: buildingId,
    from_x: fromX,
    from_y: fromY,
    to_x: toX,
    to_y: toY,
  });
}

export function removeBuilding(
  api: ApiClient,
  buildingId: string,
  x: number,
  y: number
): Promise<ApiResponse<void>> {
  return api.del<void>("/base/layout", { building_id: buildingId, x, y });
}

export function upgradeBuilding(
  api: ApiClient,
  buildingId: string,
  x: number,
  y: number
): Promise<ApiResponse<void>> {
  return api.post<void>("/base/upgrade", { building_id: buildingId, x, y });
}

interface CheckInResponse {
  completed_upgrades: {
    building_id: string;
    x: number;
    y: number;
    from_level: number;
    to_level: number;
  }[];
}

export function checkIn(api: ApiClient): Promise<ApiResponse<CheckInResponse>> {
  return api.post<CheckInResponse>("/base/check-in");
}
