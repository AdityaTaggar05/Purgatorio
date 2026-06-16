import { useState, useEffect } from "react";
import { useGame } from "../../hooks/useGame";
import * as battleApi from "../../api/endpoints/battle";
import type { MatchListEntry } from "../../types/battle";

interface MatchmakingPanelProps {
  open: boolean;
  onClose: () => void;
}

export default function MatchmakingPanel({ open, onClose }: MatchmakingPanelProps) {
  const { state, api, dispatch } = useGame();
  const [players, setPlayers] = useState<MatchListEntry[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [initiatingId, setInitiatingId] = useState<string | null>(null);

  useEffect(() => {
    if (!open) return;
    let cancelled = false;
    setError(null);
    setLoading(true);

    battleApi.getMatchList(api).then((res) => {
      if (cancelled) return;
      if (res.success) {
        setPlayers(res.data);
      } else {
        setError(res.error?.message ?? "Failed to load match list");
      }
      setLoading(false);
    });

    return () => { cancelled = true; };
  }, [open, api]);

  const handleAttack = async (player: MatchListEntry) => {
    setError(null);
    setInitiatingId(player.user_id);
    const res = await battleApi.initiateBattle(api, player.user_id);
    setInitiatingId(null);

    if (res.success) {
      dispatch({
        type: "SET_ACTIVE_BATTLE",
        payload: {
          battleId: res.data.battle_id,
          defenderName: player.username,
          defenderTerraceLevel: player.terrace_level,
          phase: "deploying",
          deployment: [],
        },
      });
      onClose();
    } else {
      setError(res.error?.message ?? "Failed to initiate battle");
    }
  };

  const level = state.layout ? Math.ceil(state.layout.grid_w / 2) : 1;

  return (
    <div className={`absolute inset-0 z-30 transition-all duration-300 ${open ? "pointer-events-auto" : "pointer-events-none"}`}>
      <div
        className={`absolute inset-0 transition-opacity duration-300 ${open ? "opacity-100 pointer-events-auto" : "opacity-0 pointer-events-none"}`}
        onClick={onClose}
      />

      <div className={`absolute top-0 right-0 h-full w-96 bg-purgatory-card border-l border-purgatory-border shadow-2xl overflow-y-auto transition-transform duration-300 ease-out ${open ? "translate-x-0" : "translate-x-full"}`}>
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="font-serif text-xl font-bold tracking-wider text-gray-200">
                Purgatorial Matchmaking
              </h2>
              <div className="text-[10px] uppercase tracking-widest text-red-500/70 mt-1">
                Terrace {level}
              </div>
            </div>
            <button
              onClick={onClose}
              className="text-gray-500 hover:text-gray-300 transition-colors p-1"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          {error && (
            <div className="mb-4 bg-red-900/20 border border-red-900/40 rounded p-3 text-sm text-red-300">
              {error}
            </div>
          )}

          {loading ? (
            <div className="flex justify-center py-12">
              <div className="w-6 h-6 border-2 border-red-500/20 border-t-red-500 rounded-full animate-spin" />
            </div>
          ) : players.length === 0 ? (
            <div className="text-center py-12">
              <div className="text-gray-500 text-sm">No penitents found at your terrace level.</div>
              <div className="text-gray-600 text-xs mt-2">Ascend higher to find worthy adversaries.</div>
            </div>
          ) : (
            <div className="space-y-3">
              {players.map((player) => (
                <div
                  key={player.user_id}
                  className="bg-purgatory-input border border-purgatory-border rounded p-4"
                >
                  <div className="flex items-start justify-between">
                    <div>
                      <div className="text-gray-200 font-bold text-sm tracking-wide">
                        {player.username}
                      </div>
                      <div className="text-[10px] uppercase tracking-widest text-gray-500 mt-0.5">
                        Terrace {player.terrace_level}
                      </div>
                    </div>
                    <button
                      onClick={() => handleAttack(player)}
                      disabled={initiatingId !== null}
                      className="text-xs uppercase tracking-widest font-bold px-4 py-1.5 rounded border transition-all
                        enabled:border-red-900/50 enabled:text-red-400 enabled:hover:bg-red-900/20 enabled:hover:border-red-600
                        disabled:border-gray-700 disabled:text-gray-600 disabled:cursor-not-allowed"
                    >
                      {initiatingId === player.user_id ? (
                        <span className="inline-block w-3 h-3 border border-red-500/40 border-t-red-500 rounded-full animate-spin" />
                      ) : (
                        "Attack"
                      )}
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
