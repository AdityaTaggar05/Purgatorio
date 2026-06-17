interface EntityBarProps {
  hp: number;
  maxHp: number;
  width?: number;
}

export default function EntityBar({ hp, maxHp, width = 36 }: EntityBarProps) {
  const pct = maxHp > 0 ? Math.max(0, Math.min(100, (hp / maxHp) * 100)) : 0;
  const color = pct > 60 ? "bg-green-500" : pct > 30 ? "bg-yellow-500" : "bg-red-500";

  return (
    <div
      className="h-1.5 bg-black/70 border border-black/50 rounded-sm overflow-hidden"
      style={{ width }}
    >
      <div
        className={`h-full ${color} transition-all duration-300 ease-out`}
        style={{ width: `${pct}%` }}
      />
    </div>
  );
}
