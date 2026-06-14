import { useState, useEffect } from "react";
import { useGame } from "../../hooks/useGame";
import * as shopApi from "../../api/endpoints/shop";
import type { ShopItem } from "../../types/building";

interface ShopPanelProps {
  open: boolean;
  onClose: () => void;
}

export default function ShopPanel({ open, onClose }: ShopPanelProps) {
  const { state, api, dispatch } = useGame();
  const [items, setItems] = useState<ShopItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!open) return;
    let cancelled = false;
    setError(null);
    setLoading(true);

    shopApi.getShop(api).then((res) => {
      if (cancelled) return;
      if (res.success) {
        setItems(res.data.items);
      } else {
        setError(res.error?.message ?? "Failed to load shop");
      }
      setLoading(false);
    });

    return () => { cancelled = true; };
  }, [open, api]);

  const handleBuy = async (buildingId: string) => {
    setError(null);
    const res = await shopApi.buyBuilding(api, buildingId);
    if (res.success) {
      const economyRes = await api.get<{ penitence: number; grace: number; max_penitence: number }>("/user/economy");
      if (economyRes.success) {
        dispatch({ type: "SET_ECONOMY", payload: economyRes.data });
      }
      shopApi.getShop(api).then((r) => {
        if (r.success) setItems(r.data.items);
      });
    } else {
      setError(res.error?.message ?? "Purchase failed");
    }
  };

  const penitence = state.economy?.penitence ?? 0;
  const grace = state.economy?.grace ?? 0;

  return (
    <div className={`absolute inset-0 z-30 transition-all duration-300 ${open ? "pointer-events-auto" : "pointer-events-none"}`}>
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
              Altar of Exchange
            </h2>
            <button
              onClick={onClose}
              className="text-gray-500 hover:text-gray-300 transition-colors p-1"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          {/* Resource display */}
          <div className="flex gap-4 mb-6">
            <div className="flex-1 bg-purgatory-input border border-purgatory-border rounded p-2">
              <div className="text-[9px] uppercase tracking-widest text-purple-400 font-bold">Penitence</div>
              <div className="text-gray-200 font-medium text-sm">{penitence.toLocaleString()}</div>
            </div>
            <div className="flex-1 bg-purgatory-input border border-purgatory-border rounded p-2">
              <div className="text-[9px] uppercase tracking-widest text-teal-400 font-bold">Grace</div>
              <div className="text-gray-200 font-medium text-sm">{grace.toLocaleString()}</div>
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
              {items.map((item) => (
                <div
                  key={item.building.id}
                  className="bg-purgatory-input border border-purgatory-border rounded p-4"
                >
                  <div className="flex items-start justify-between mb-2">
                    <div>
                      <div className="text-gray-200 font-bold text-sm tracking-wide">
                        {item.building.name}
                      </div>
                      <div className="text-[10px] uppercase tracking-widest text-gray-500 mt-0.5">
                        {item.building.category} · Size {item.building.size}
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="text-amber-500 font-bold text-sm">
                        {item.building.price}
                      </div>
                      <div className="text-[10px] uppercase tracking-wider text-amber-600/70">
                        {item.building.currency}
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="text-[10px] text-gray-500">
                      Owned: {item.current_owned} / {item.max_allowed}
                    </div>
                    <button
                      onClick={() => handleBuy(item.building.id)}
                      disabled={!item.can_buy}
                      className="text-xs uppercase tracking-widest font-bold px-4 py-1.5 rounded border transition-all
                        enabled:border-amber-500/40 enabled:text-amber-400 enabled:hover:bg-amber-500/10 enabled:hover:border-amber-400
                        disabled:border-gray-700 disabled:text-gray-600 disabled:cursor-not-allowed"
                    >
                      Buy
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
