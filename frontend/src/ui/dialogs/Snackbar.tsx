import { useEffect } from 'react';

interface SnackbarProps {
  message: string | null;
  onClose: () => void;
}

export default function Snackbar({ message, onClose }: SnackbarProps) {
  useEffect(() => {
    if (!message) return;

    const timer = setTimeout(() => {
      onClose();
    }, 4000);

    return () => clearTimeout(timer);
  }, [message, onClose]);

  if (!message) return null;

  return (
    <div className="fixed bottom-6 right-6 z-50 flex items-center justify-between w-full max-w-sm bg-[#1e1515] border border-red-900/60 rounded-lg p-4 shadow-[0_8px_30px_rgb(0,0,0,0.5)] animate-fade-in-up backdrop-blur-md">
      <div className="flex items-start space-x-3">
        <span className="text-red-500 font-serif font-bold mt-0.5 select-none">♦</span>
        <div>
          <p className="font-serif text-xs uppercase tracking-widest text-red-400 font-bold mb-0.5">
            The Gates Remain Shut
          </p>
          <p className="text-sm font-medium text-gray-300">
            {message}
          </p>
        </div>
      </div>

      <button
        onClick={onClose}
        className="text-gray-500 hover:text-gray-300 transition-colors ml-4 p-1 focus:outline-none"
        aria-label="Dismiss notification"
      >
        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.4} d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  );
}
