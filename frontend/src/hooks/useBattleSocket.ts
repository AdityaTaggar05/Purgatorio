import { useState, useEffect, useCallback, useRef } from "react";
import { BattleSocket, type BattleEndPayload } from "../api/ws";
import type { TroopDeployment, TickResult } from "../types/battle";

export interface UseBattleSocketResult {
  sendDeploy: (deployment: TroopDeployment[]) => void;
  sendDone: () => void;
  ticks: TickResult[];
  battleResult: BattleEndPayload | null;
  error: string | null;
  state: "connecting" | "open" | "deployed" | "closed" | "error";
  deployCountdown: number;
}

export function useBattleSocket(battleId: string, token: string | null): UseBattleSocketResult {
  const [ticks, setTicks] = useState<TickResult[]>([]);
  const [battleResult, setBattleResult] = useState<BattleEndPayload | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [state, setState] = useState<UseBattleSocketResult["state"]>("connecting");
  const [deployCountdown, setDeployCountdown] = useState(30);
  const socketRef = useRef<BattleSocket | null>(null);
  const deployedRef = useRef(false);

  useEffect(() => {
    deployedRef.current = false;
    setTicks([]);
    setBattleResult(null);
    setError(null);
    setState("connecting");
    setDeployCountdown(30);

    const socket = new BattleSocket(battleId, token);
    socketRef.current = socket;

    socket.onOpen(() => {
      if (!deployedRef.current) setState("open");
    });

    socket.onDeployCountdown((seconds) => {
      setDeployCountdown(seconds);
    });

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

    socket.connect();

    return () => {
      socket.disconnect();
    };
  }, [battleId]);

  const sendDeploy = useCallback(
    (deployment: TroopDeployment[]) => {
      deployedRef.current = true;
      socketRef.current?.sendDeploy(deployment);
      setState("deployed");
    },
    []
  );

  const sendDone = useCallback(() => {
    socketRef.current?.sendDone();
  }, []);

  return { sendDeploy, sendDone, ticks, battleResult, error, state, deployCountdown };
}
