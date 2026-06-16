import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

export default function SplashScreen() {
  const { user, accessToken, isLoading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!isLoading) {
      if (user && accessToken) {
        navigate('/game', { replace: true });
      } else {
        navigate('/login', { replace: true });
      }
    }
  }, [isLoading, user, accessToken, navigate]);

  return (
    <div className="min-h-screen bg-purgatory-bg text-gray-200 flex flex-col justify-center items-center p-4 relative overflow-hidden">
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-amber-500/5 rounded-full blur-[140px] animate-pulse"></div>
      
      <div className="text-center space-y-6 relative z-10">
        <h1 className="font-serif text-5xl font-extrabold tracking-[0.3em] text-[#f8fafc] animate-pulse">
          PURGATORIO
        </h1>
        
        <div className="flex justify-center items-center pt-4">
          <div className="w-6 h-6 border-2 border-amber-500/20 border-t-amber-500 rounded-full animate-spin"></div>
        </div>

        <p className="font-serif text-xs uppercase tracking-[0.25em] text-gray-500 pt-2 animate-pulse">
          Preparing the Ascent...
        </p>
      </div>
    </div>
  );
}
