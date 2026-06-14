import type { ApiClient } from "../client";
import type { ApiResponse } from "../../types/api";
import type { UserEconomy } from "../../types/economy";

export function getEconomy(api: ApiClient): Promise<ApiResponse<UserEconomy>> {
  return api.get<UserEconomy>("/user/economy");
}

export function collectResources(api: ApiClient): Promise<ApiResponse<UserEconomy>> {
  return api.post<UserEconomy>("/user/economy/collect");
}
