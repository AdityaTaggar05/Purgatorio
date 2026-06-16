import { useEffect, useState } from "react";

interface SnackbarProps {
  message: string;
  onDone: () => void;
  duration?: number;
}

export default function UpgradeSnackbar({ message, onDone, duration = 5000 }: SnackbarProps) {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    requestAnimationFrame(() => setVisible(true));
    const timer = setTimeout(() => {
      setVisible(false);
      setTimeout(onDone, 300);
    }, duration);
    return () => clearTimeout(timer);
  }, [duration, onDone]);

  return (
    <div className="fixed bottom-24 left-1/2 -translate-x-1/2 z-50 pointer-events-none">
      <div
        className={`pointer-events-auto bg-teal-900/90 backdrop-blur-md border border-teal-500/40 rounded-lg px-6 py-3 shadow-2xl transition-all duration-300
          ${visible ? "opacity-100 translate-y-0" : "opacity-0 translate-y-4"}`}
      >
        <div className="flex items-center gap-3">
          <svg className="w-5 h-5 text-teal-400 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 12.75l6 6 9-13.5" />
          </svg>
          <span className="text-sm text-teal-100 font-medium">{message}</span>
        </div>
      </div>
    </div>
  );
}
