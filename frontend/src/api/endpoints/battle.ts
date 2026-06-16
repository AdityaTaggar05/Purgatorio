import type { ApiClient } from "../client";
import type { ApiResponse } from "../../types/api";
import type {
  MatchListEntry,
  InitiateResponse,
  BattleResultResponse,
  ReplayData,
} from "../../types/battle";

export function getMatchList(api: ApiClient): Promise<ApiResponse<MatchListEntry[]>> {
  return api.get<MatchListEntry[]>("/battle/matchlist");
}

export function initiateBattle(
  api: ApiClient,
  defenderId: string
): Promise<ApiResponse<InitiateResponse>> {
  return api.post<InitiateResponse>("/battle/initiate", { defender_id: defenderId });
}

export function getBattleResult(
  api: ApiClient,
  battleId: string
): Promise<ApiResponse<BattleResultResponse>> {
  return api.get<BattleResultResponse>(`/battle/${battleId}/result`);
}

export function getAttackReplays(
  api: ApiClient
): Promise<ApiResponse<ReplayData[]>> {
  return api.get<ReplayData[]>("/battle/replays/attacks");
}

export function getDefenseReplays(
  api: ApiClient
): Promise<ApiResponse<ReplayData[]>> {
  return api.get<ReplayData[]>("/battle/replays/defenses");
}

export function getReplayData(
  api: ApiClient,
  battleId: string
): Promise<ApiResponse<ReplayData>> {
  return api.get<ReplayData>(`/battle/${battleId}/replay`);
}
