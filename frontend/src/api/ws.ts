import type { TroopDeployment, TickResult, BattleResultResponse } from "../types/battle";
import { WS_BASE_URL } from "../config";

export interface BattleEndPayload {
  outcome: "victory" | "defeat" | "threshold_failed";
  destruction: number;
  loot: number;
  sin_meter: number;
  duration_ticks: number;
}

type SocketState = "connecting" | "open" | "closed" | "error";

type ServerMessage =
  | { type: "tick_batch"; ticks: TickResult[]; batch_start: number }
  | { type: "battle_end" } & BattleEndPayload
  | { type: "error"; message: string };

export class BattleSocket {
  private ws: WebSocket | null = null;
  private battleId: string;
  private _state: SocketState = "connecting";
  private tickBatchCallbacks: Array<(ticks: TickResult[], batchStart: number) => void> = [];
  private battleEndCallbacks: Array<(result: BattleEndPayload) => void> = [];
  private errorCallbacks: Array<(message: string) => void> = [];
  private deployTimer: ReturnType<typeof setTimeout> | null = null;
  private _deployCountdown = 30;

  constructor(battleId: string) {
    this.battleId = battleId;
  }

  get state(): SocketState {
    return this._state;
  }

  get deployCountdown(): number {
    return this._deployCountdown;
  }

  connect(): void {
    this._state = "connecting";
    const url = `${WS_BASE_URL}/battle/${this.battleId}/ws`;
    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      this._state = "open";
      this.startDeployCountdown();
    };

    this.ws.onmessage = (event) => {
      try {
        const msg: ServerMessage = JSON.parse(event.data);
        switch (msg.type) {
          case "tick_batch":
            this.tickBatchCallbacks.forEach((cb) => cb(msg.ticks, msg.batch_start));
            break;
          case "battle_end":
            this.stopDeployCountdown();
            this.battleEndCallbacks.forEach((cb) =>
              cb({
                outcome: msg.outcome,
                destruction: msg.destruction,
                loot: msg.loot,
                sin_meter: msg.sin_meter,
                duration_ticks: msg.duration_ticks,
              })
            );
            this.disconnect();
            break;
          case "error":
            this.stopDeployCountdown();
            this.errorCallbacks.forEach((cb) => cb(msg.message));
            this._state = "error";
            break;
        }
      } catch {
        this.errorCallbacks.forEach((cb) => cb("Failed to parse server message"));
        this._state = "error";
      }
    };

    this.ws.onclose = () => {
      this._state = "closed";
      this.stopDeployCountdown();
    };

    this.ws.onerror = () => {
      if (this._state === "connecting") {
        this._state = "error";
        this.errorCallbacks.forEach((cb) => cb("Failed to connect to battle server"));
      }
    };
  }

  sendDeploy(deployment: TroopDeployment[]): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return;
    this.stopDeployCountdown();
    this.ws.send(JSON.stringify({ type: "deploy", troops: deployment }));
  }

  disconnect(): void {
    this.stopDeployCountdown();
    if (this.ws) {
      this.ws.close(1000);
    }
  }

  onTickBatch(callback: (ticks: TickResult[], batchStart: number) => void): void {
    this.tickBatchCallbacks.push(callback);
  }

  onBattleEnd(callback: (result: BattleEndPayload) => void): void {
    this.battleEndCallbacks.push(callback);
  }

  onError(callback: (message: string) => void): void {
    this.errorCallbacks.push(callback);
  }

  private startDeployCountdown(): void {
    this._deployCountdown = 30;
    this.deployTimer = setInterval(() => {
      this._deployCountdown--;
      if (this._deployCountdown <= 0) {
        this.stopDeployCountdown();
        this.errorCallbacks.forEach((cb) => cb("Deployment timed out"));
        this.disconnect();
      }
    }, 1000);
  }

  private stopDeployCountdown(): void {
    if (this.deployTimer) {
      clearInterval(this.deployTimer);
      this.deployTimer = null;
    }
  }
}
