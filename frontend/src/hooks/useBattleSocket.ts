import { useState, useEffect, useCallback, useRef } from "react";
import { BattleSocket, type BattleEndPayload } from "../api/ws";
import type { TroopDeployment, TickResult } from "../types/battle";

interface UseBattleSocketResult {
  sendDeploy: (deployment: TroopDeployment[]) => void;
  ticks: TickResult[];
  battleResult: BattleEndPayload | null;
  error: string | null;
  state: "connecting" | "open" | "deployed" | "closed" | "error";
  deployCountdown: number;
}

export function useBattleSocket(battleId: string): UseBattleSocketResult {
  const [ticks, setTicks] = useState<TickResult[]>([]);
  const [battleResult, setBattleResult] = useState<BattleEndPayload | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [state, setState] = useState<UseBattleSocketResult["state"]>("connecting");
  const [deployCountdown, setDeployCountdown] = useState(30);
  const socketRef = useRef<BattleSocket | null>(null);

  useEffect(() => {
    const socket = new BattleSocket(battleId);
    socketRef.current = socket;

    socket.onTickBatch((incomingTicks) => {
      setTicks((prev) => [...prev, ...incomingTicks]);
    });

    socket.onBattleEnd((result) => {
      setBattleResult(result);
      setState("closed");
    });

    socket.onError((msg) => {
      setError(msg);
      setState("error");
    });

    const countdownInterval = setInterval(() => {
      if (socket.state === "open") {
        setDeployCountdown(socket.deployCountdown);
      }
      if (socket.state === "open" && !battleResult) {
        setState("open");
      }
    }, 500);

    socket.connect();

    return () => {
      clearInterval(countdownInterval);
      socket.disconnect();
    };
  }, [battleId]);

  const sendDeploy = useCallback(
    (deployment: TroopDeployment[]) => {
      socketRef.current?.sendDeploy(deployment);
      setState("deployed");
    },
    []
  );

  return { sendDeploy, ticks, battleResult, error, state, deployCountdown };
}
