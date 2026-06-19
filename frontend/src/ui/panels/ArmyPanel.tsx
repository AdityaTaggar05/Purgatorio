import { useState, useEffect } from "react";
import { useGame } from "../../hooks/useGame";
import * as armyApi from "../../api/endpoints/army";
import type { Troop, ArmyResponse } from "../../types/army";

interface ArmyPanelProps {
  open: boolean;
  onClose: () => void;
}

export default function ArmyPanel({ open, onClose }: ArmyPanelProps) {
  const { state, api, dispatch } = useGame();
  const [catalog, setCatalog] = useState<Troop[]>([]);
  const [army, setArmy] = useState<ArmyResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [training, setTraining] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [trainCounts, setTrainCounts] = useState<Record<string, number>>({});

  useEffect(() => {
    if (!open) return;
    let cancelled = false;
    setError(null);
    setLoading(true);

    Promise.all([
      armyApi.getCatalog(api),
      armyApi.getMyTroops(api),
    ]).then(([catalogRes, armyRes]) => {
      if (cancelled) return;
      if (catalogRes.success) setCatalog(catalogRes.data.troops);
      if (armyRes.success) setArmy(armyRes.data);
      if (!catalogRes.success || !armyRes.success) {
        setError(catalogRes.error?.message ?? armyRes.error?.message ?? "Failed to load army data");
      }
      setLoading(false);
    });

    return () => { cancelled = true; };
  }, [open, api]);

  const handleTrain = async (troopId: string) => {
    const count = trainCounts[troopId] ?? 1;
    setError(null);
    setTraining(troopId);
    try {
      const res = await armyApi.trainTroop(api, troopId, count);
      if (res.success) {
        const armyRes = await armyApi.getMyTroops(api);
        if (armyRes.success) {
          setArmy(armyRes.data);
          dispatch({ type: "SET_ARMY", payload: armyRes.data });
        }

        const econRes = await api.get<{ penitence: number; grace: number; max_penitence: number }>("/user/economy");
        if (econRes.success) {
          dispatch({ type: "SET_ECONOMY", payload: econRes.data });
        }
      } else {
        setError(res.error?.message ?? "Training failed");
      }
    } finally {
      setTraining(null);
    }
  };

  const handleDetrain = async (troopId: string) => {
    const count = trainCounts[troopId] ?? 1;
    setError(null);
    setTraining(troopId);
    try {
      const res = await armyApi.detrainTroop(api, troopId, count);
      if (res.success) {
        const armyRes = await armyApi.getMyTroops(api);
        if (armyRes.success) {
          setArmy(armyRes.data);
          dispatch({ type: "SET_ARMY", payload: armyRes.data });
        }

        const econRes = await api.get<{ penitence: number; grace: number; max_penitence: number }>("/user/economy");
        if (econRes.success) {
          dispatch({ type: "SET_ECONOMY", payload: econRes.data });
        }
      } else {
        setError(res.error?.message ?? "Detraining failed");
      }
    } finally {
      setTraining(null);
    }
  };

  const decrementCount = (troopId: string) => {
    setTrainCounts(prev => {
      const v = Math.max(1, (prev[troopId] ?? 1) - 1);
      return { ...prev, [troopId]: v };
    });
  };

  const incrementCount = (troopId: string) => {
    setTrainCounts(prev => {
      const v = (prev[troopId] ?? 1) + 1;
      return { ...prev, [troopId]: v };
    });
  };

  const penitence = state.economy?.penitence ?? 0;
  const used = army?.used_capacity ?? 0;
  const max = army?.max_capacity ?? 0;

  return (
    <div className={`absolute inset-0 z-40 transition-all duration-300 ${open ? "pointer-events-auto" : "pointer-events-none"}`}>
      {/* Backdrop */}
      <div
        className={`absolute inset-0 transition-opacity duration-300 ${open ? "opacity-100 pointer-events-auto" : "opacity-0 pointer-events-none"}`}
        onClick={onClose}
      />

      {/* Panel */}
      <div className={`absolute top-0 right-0 h-full w-96 bg-purgatory-card border-l border-purgatory-border shadow-2xl overflow-y-auto transition-transform duration-300 ease-out ${open ? "translate-x-0" : "translate-x-full"}`}>
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="font-serif text-xl font-bold tracking-wider text-gray-200">
              Legion
            </h2>
            <button onClick={onClose} disabled={training !== null} className="text-gray-500 hover:text-gray-300 p-1 transition-colors disabled:opacity-30 disabled:cursor-not-allowed">
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          {/* Capacity bar */}
          <div className="mb-6">
            <div className="flex justify-between text-[10px] uppercase tracking-widest text-gray-500 mb-1">
              <span>Capacity</span>
              <span>{used} / {max}</span>
            </div>
            <div className="h-2 bg-purgatory-input border border-purgatory-border rounded-full overflow-hidden">
              <div
                className="h-full bg-amber-500/70 rounded-full transition-all"
                style={{ width: `${max > 0 ? (used / max) * 100 : 0}%` }}
              />
            </div>
          </div>

          {/* Resources */}
          <div className="flex gap-4 mb-6">
            <div className="flex-1 bg-purgatory-input border border-purgatory-border rounded p-2">
              <div className="text-[9px] uppercase tracking-widest text-purple-400 font-bold">Penitence</div>
              <div className="text-gray-200 font-medium text-sm">{penitence.toLocaleString()}</div>
            </div>
          </div>

          {error && (
            <div className="mb-4 bg-red-900/20 border border-red-900/40 rounded p-3 text-sm text-red-300">
              {error}
            </div>
          )}

          {loading ? (
            <div className="flex justify-center py-12">
              <div className="w-6 h-6 border-2 border-amber-500/20 border-t-amber-500 rounded-full animate-spin" />
            </div>
          ) : (
            <div className="space-y-3">
              {catalog.map((troop) => {
                const owned = army?.troops[troop.id] ?? 0;
                const count = trainCounts[troop.id] ?? 1;
                const totalCost = troop.training_cost * count;
                const canAfford = penitence >= totalCost;
                const hasSpace = used + troop.space <= max;

                return (
                  <div
                    key={troop.id}
                    className="bg-purgatory-input border border-purgatory-border rounded p-4"
                  >
                    <div className="flex items-start justify-between mb-2">
                      <div>
                        <div className="text-gray-200 font-bold text-sm tracking-wide">
                          {troop.name}
                        </div>
                        <div className="text-[10px] uppercase tracking-widest text-gray-500 mt-0.5">
                          HP {troop.hp} · DPS {troop.dps} · Space {troop.space}
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="text-amber-500 font-bold text-sm">
                          {troop.training_cost}
                        </div>
                        <div className="text-[10px] uppercase tracking-wider text-amber-600/70">
                          penitence
                        </div>
                      </div>
                    </div>

                    {troop.preferred_target && (
                      <div className="text-[10px] text-gray-500 mb-2">
                        Targets: {troop.preferred_target}
                      </div>
                    )}

                    <div className="flex items-center justify-between">
                      <div className="text-[10px] text-gray-500">
                        Trained: {owned}
                      </div>

                      <div className="flex items-center gap-2">
                        <div className="flex items-center gap-1">
                          <button
                            onClick={() => decrementCount(troop.id)}
                            disabled={training !== null}
                            className="w-6 h-6 flex items-center justify-center text-xs border border-purgatory-border rounded text-gray-500 hover:text-gray-200 hover:border-gray-600 transition-all disabled:opacity-30 disabled:cursor-not-allowed"
                          >
                            −
                          </button>
                          <span className="text-xs text-gray-300 w-6 text-center font-bold">{count}</span>
                          <button
                            onClick={() => incrementCount(troop.id)}
                            disabled={training !== null}
                            className="w-6 h-6 flex items-center justify-center text-xs border border-purgatory-border rounded text-gray-500 hover:text-gray-200 hover:border-gray-600 transition-all disabled:opacity-30 disabled:cursor-not-allowed"
                          >
                            +
                          </button>
                        </div>
                        <button
                          onClick={() => handleTrain(troop.id)}
                          disabled={!canAfford || !hasSpace || training !== null}
                          className="text-xs uppercase tracking-widest font-bold px-4 py-1.5 rounded border transition-all min-w-[90px]
                            enabled:border-amber-500/40 enabled:text-amber-400 enabled:hover:bg-amber-500/10 enabled:hover:border-amber-400
                            disabled:border-gray-700 disabled:text-gray-600 disabled:cursor-not-allowed"
                        >
                          {training === troop.id ? (
                            <span className="flex items-center justify-center gap-1">
                              <span className="w-3 h-3 border border-amber-500/30 border-t-amber-400 rounded-full animate-spin inline-block" />
                              Training
                            </span>
                          ) : (
                            <>Train · {totalCost}</>
                          )}
                        </button>
                        {owned >= count && (
                          <button
                            onClick={() => handleDetrain(troop.id)}
                            disabled={training !== null}
                            className="text-[10px] uppercase tracking-widest font-bold px-3 py-1.5 rounded border transition-all
                              border-red-900/40 text-red-400 hover:bg-red-900/20 hover:border-red-600/60
                              disabled:border-gray-700 disabled:text-gray-600 disabled:cursor-not-allowed"
                          >
                            {training === troop.id ? (
                              <span className="flex items-center justify-center gap-1">
                                <span className="w-3 h-3 border border-red-400/30 border-t-red-400 rounded-full animate-spin inline-block" />
                                ...
                              </span>
                            ) : (
                              <>Detrain · {count}</>
                            )}
                          </button>
                        )}
                      </div>
                    </div>
                  </div>
                );
              })}

              {catalog.length === 0 && (
                <div className="text-center text-gray-500 text-sm py-6">
                  No troops available. Build a Barracks first.
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
