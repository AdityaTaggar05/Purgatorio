import { useState, type ChangeEvent, type SubmitEvent } from 'react';
import AuthLayout from '../../ui/layout/Auth';

export default function RegisterPage() {
  const [showPassword, setShowPassword] = useState<boolean>(false);
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
  });

  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log('Registering...', formData);
  };

  return (
    <AuthLayout title="THE ASCENT" subtitle="Leave behind the dark valley">
      <form onSubmit={handleSubmit} className="space-y-5">
        <div className="space-y-1.5">
          <label htmlFor="username" className="text-xs uppercase tracking-wider text-gray-400 font-medium">
            Username
          </label>
          <input
            type="text"
            id="username"
            name="username"
            required
            value={formData.username}
            onChange={handleInputChange}
            placeholder="e.g., Beatrice"
            className="w-full bg-[#22252a] border border-[#343a40] rounded px-4 py-2.5 text-sm text-gray-200 placeholder-gray-600 focus:outline-none focus:border-amber-500/50 transition-colors"
          />
        </div>

        <div className="space-y-1.5">
          <label htmlFor="email" className="text-xs uppercase tracking-wider text-gray-400 font-medium">
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
            className="w-full bg-[#22252a] border border-[#343a40] rounded px-4 py-2.5 text-sm text-gray-200 placeholder-gray-600 focus:outline-none focus:border-amber-500/50 transition-colors"
          />
        </div>

        <div className="space-y-1.5">
          <label htmlFor="password" className="text-xs uppercase tracking-wider text-gray-400 font-medium block">
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
              className="w-full bg-[#22252a] border border-[#343a40] rounded px-4 py-2.5 pr-10 text-sm text-gray-200 placeholder-gray-600 focus:outline-none focus:border-amber-500/50 transition-colors"
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-300 focus:outline-none text-xs uppercase font-mono tracking-tighter"
            >
              {showPassword ? 'Hide' : 'Show'}
            </button>
          </div>
        </div>

        <button
          type="submit"
          className="w-full font-cinzel bg-[#2d3136] hover:bg-amber-950/40 border border-amber-500/30 hover:border-amber-500/60 text-[#e2e8f0] py-3 rounded text-sm tracking-widest font-semibold transition-all duration-300 mt-2 shadow-lg"
        >
          BEGIN PENANCE
        </button>
      </form>

      <div className="mt-8 pt-6 border-t border-[#2d3136] text-center">
        <p className="text-xs tracking-wide text-gray-400">
          Already walking the terraces?{' '}
          <a href="/login" className="text-amber-500/60 hover:text-amber-500/80 underline underline-offset-4 ml-1 transition-colors">
            Sign In
          </a>
        </p>
      </div>
    </AuthLayout>
  );
}
