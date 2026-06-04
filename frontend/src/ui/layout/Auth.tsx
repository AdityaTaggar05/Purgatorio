import type { ReactNode } from 'react';

interface AuthLayoutProps {
  children: ReactNode;
  title: string;
  subtitle: string;
}

export default function AuthLayout({ children, title, subtitle }: AuthLayoutProps) {
  return (
    <div className="min-h-screen flex flex-col justify-center items-center p-4 relative overflow-hidden">

      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full max-w-7xl h-full pointer-events-none opacity-20">
        <div className="absolute top-0 left-1/4 w-96 h-96 bg-amber-500/10 rounded-full blur-[120px]"></div>
        <div className="absolute bottom-10 right-1/4 w-96 h-96 bg-orange-700/5 rounded-full blur-[120px]"></div>
      </div>

      <div className="w-full max-w-md bg-purgatory-card border border-purgatory-border rounded-lg shadow-2xl p-8 relative z-10 backdrop-blur-sm">

        <div className="text-center mb-8">
          <h1 className="font-serif text-3xl font-bold tracking-wider text-[#e2e8f0] mb-2">
            {title}
          </h1>
          <p className="text-xs uppercase tracking-[0.2em] text-amber-500/70 font-serif">
            {subtitle}
          </p>
        </div>

        {children}

      </div>

      <footer className="mt-8 text-[10px] font-serif font-semibold tracking-[0.3em] text-gray-600 text-center uppercase relative z-10">
        Pure and disposed to mount unto the stars
      </footer>
    </div>
  );
}
