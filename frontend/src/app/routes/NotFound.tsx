import { Link } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

export default function NotFound() {
  const { user, accessToken, isLoading } = useAuth();
  const isAuthenticated = user && accessToken;

  const redirectPath = isAuthenticated ? '/game' : '/login';
  const buttonText = isAuthenticated ? 'RETURN TO THE ASCENT' : 'ENTER THE GATE';

  return (
    <div className="min-h-screen bg-purgatory-bg text-gray-200 flex flex-col justify-center items-center p-4 relative overflow-hidden">
      <div className="absolute top-0 left-1/4 w-96 h-96 bg-amber-500/10 rounded-full blur-[120px] opacity-20"></div>
      <div className="absolute bottom-10 right-1/4 w-96 h-96 bg-orange-700/5 rounded-full blur-[120px] opacity-20"></div>

      <div className="text-center space-y-6 relative z-10">
        <h1 className="font-serif text-8xl font-extrabold tracking-[0.2em] text-amber-500/80">
          404
        </h1>

        <p className="font-serif text-xl tracking-wider text-gray-400">
          Lost Soul
        </p>

        <div className="w-16 h-px bg-amber-500/20 mx-auto"></div>

        <p className="text-sm tracking-wide text-gray-500 max-w-sm mx-auto leading-relaxed">
          You have strayed from the mountain path. No terrace lies at this height<br />Only silence and shadow.
        </p>

        <Link
          to={redirectPath}
          className="inline-block font-serif bg-purgatory-border hover:bg-amber-950/40 border border-amber-500/30 hover:border-amber-500/60 text-gray-200 px-8 py-3 rounded text-sm tracking-widest font-semibold transition-all duration-300 shadow-lg mt-2"
        >
          {isLoading ? (
            <span className="inline-flex items-center gap-3">
              <span className="w-4 h-4 border-2 border-amber-500/20 border-t-amber-500 rounded-full animate-spin"></span>
              SEARCHING...
            </span>
          ) : (
            buttonText
          )}
        </Link>
      </div>

      <footer className="absolute bottom-8 text-[10px] font-serif font-semibold tracking-[0.3em] text-gray-600 text-center uppercase">
        Pure and disposed to mount unto the stars
      </footer>
    </div>
  );
}
