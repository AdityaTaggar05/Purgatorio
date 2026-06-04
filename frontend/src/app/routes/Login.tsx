import { useState, type ChangeEvent, type SubmitEvent } from 'react';
import AuthLayout from '../../ui/layout/Auth';

export default function LoginPage() {
  const [showPassword, setShowPassword] = useState<boolean>(false);
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });

  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log('Logging in...', formData);
  };

  return (
    <AuthLayout title="PURGATORIO" subtitle="To see again the stars">
      <form onSubmit={handleSubmit} className="space-y-5">
        <div className="space-y-1.5">
          <label htmlFor="email" className="text-sm uppercase tracking-wider text-gray-400 font-medium">
            Email Address
          </label>
          <input
            type="email"
            id="email"
            name="email"
            required
            value={formData.email}
            onChange={handleInputChange}
            placeholder="name@domain.com"
            className="w-full bg-purgatory-input border border-purgatory-input-border rounded px-4 py-2.5 text-base text-gray-200 placeholder-gray-600 focus:outline-none focus:border-amber-500/50 transition-colors"
          />
        </div>

        <div className="space-y-1.5">
          <label htmlFor="password" className="text-sm uppercase tracking-wider text-gray-400 font-medium">
            Password
          </label>
          <div className="relative">
            <input
              type={showPassword ? 'text' : 'password'}
              id="password"
              name="password"
              required
              value={formData.password}
              onChange={handleInputChange}
              placeholder="••••••••"
              className="w-full bg-purgatory-input border border-purgatory-input-border rounded px-4 py-2.5 pr-12 text-base text-gray-200 placeholder-gray-600 focus:outline-none focus:border-amber-500/50 transition-colors"
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-300 focus:outline-none text-xs uppercase font-mono tracking-tighter hover:cursor-pointer"
            >
              {showPassword ? 'Hide' : 'Show'}
            </button>
          </div>
        </div>

        <button
          type="submit"
          className="w-full font-serif bg-purgatory-border hover:bg-amber-950/40 border border-amber-500/30 hover:border-amber-500/60 text-gray-200 py-3 rounded text-sm tracking-widest font-semibold transition-all duration-300 mt-4 shadow-lg hover:cursor-pointer"
        >
          ENTER THE GATE
        </button>
      </form>

      <div className="mt-8 pt-6 border-t border-purgatory-border text-center">
        <p className="text-sm tracking-wide text-gray-400">
          New to the journey?{' '}
          <a href="/register" className="text-amber-500/60 hover:text-amber-500/80 underline underline-offset-4 ml-1 transition-colors">
            Create an account
          </a>
        </p>
      </div>
    </AuthLayout>
  );
}
