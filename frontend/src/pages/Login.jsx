import { useState } from 'react';
import { Link } from 'react-router-dom';
import { GoogleLogin } from '@react-oauth/google';
import { useAuth } from '../context/AuthContext';
import { authLogin, authGoogle } from '../api/client';
import { Wallet, Mail, Lock, Eye, EyeOff, TrendingUp, ShieldCheck, Sparkles } from 'lucide-react';

const GOOGLE_ENABLED = !!import.meta.env.VITE_GOOGLE_CLIENT_ID;

const FEATURES = [
  { icon: TrendingUp,   text: 'Pantau tren keuangan secara real-time' },
  { icon: ShieldCheck,  text: 'Data privat & aman, hanya milikmu' },
  { icon: Sparkles,     text: 'Prediksi keuangan berbasis AI' },
];

export default function Login() {
  const { login } = useAuth();
  const [form, setForm] = useState({ email: '', password: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await authLogin(form);
      login(res);
    } catch (err) {
      setError(err?.response?.data?.error || 'Email atau password salah');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex bg-gray-50 dark:bg-[#0f1117]">

      {/* ── Left branding panel ──────────────────────────── */}
      <div className="hidden lg:flex lg:w-[45%] relative overflow-hidden">
        {/* Gradient background */}
        <div className="absolute inset-0 bg-gradient-to-br from-blue-600 via-indigo-700 to-violet-800" />
        {/* Decorative blobs */}
        <div className="absolute -top-24 -left-24 w-96 h-96 bg-blue-400/20 rounded-full blur-3xl" />
        <div className="absolute -bottom-24 -right-24 w-96 h-96 bg-violet-500/20 rounded-full blur-3xl" />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-64 h-64 bg-indigo-400/10 rounded-full blur-2xl" />

        <div className="relative flex flex-col justify-center px-14 max-w-lg mx-auto">
          {/* Logo */}
          <div className="flex items-center gap-3 mb-10">
            <div className="p-2.5 bg-white/15 rounded-2xl backdrop-blur-sm border border-white/20">
              <Wallet size={28} className="text-white" />
            </div>
            <span className="text-2xl font-bold text-white tracking-tight">Catatan Keuangan</span>
          </div>

          <h2 className="text-4xl font-bold text-white leading-tight mb-4">
            Kelola keuangan<br />dengan lebih <span className="text-blue-200">cerdas</span>
          </h2>
          <p className="text-blue-100 text-base leading-relaxed mb-10">
            Catat setiap rupiah, pantau tren, dan raih tujuan finansialmu — semua dalam satu platform.
          </p>

          <div className="space-y-4">
            {FEATURES.map(({ icon: Icon, text }, i) => (
              <div key={i} className="flex items-center gap-3">
                <div className="p-2 bg-white/10 rounded-xl border border-white/15 shrink-0">
                  <Icon size={18} className="text-blue-200" />
                </div>
                <span className="text-blue-100 text-sm">{text}</span>
              </div>
            ))}
          </div>

          {/* Stats row */}
          <div className="flex gap-8 mt-12 pt-8 border-t border-white/15">
            {[['100%', 'Gratis'], ['24/7', 'Akses'], ['0', 'Iklan']].map(([val, label]) => (
              <div key={label}>
                <p className="text-2xl font-bold text-white">{val}</p>
                <p className="text-blue-200 text-xs mt-0.5">{label}</p>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* ── Right form panel ─────────────────────────────── */}
      <div className="flex-1 flex items-center justify-center p-5 sm:p-8">
        <div className="w-full max-w-sm">

          {/* Mobile logo */}
          <div className="lg:hidden flex items-center gap-2.5 mb-8 justify-center">
            <div className="p-2 bg-blue-600 rounded-xl text-white shadow-lg shadow-blue-600/30">
              <Wallet size={22} />
            </div>
            <span className="text-lg font-bold text-gray-800 dark:text-gray-100">Catatan Keuangan</span>
          </div>

          <div className="mb-6">
            <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Selamat datang kembali 👋</h2>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">Masuk untuk melanjutkan</p>
          </div>

          {error && (
            <div className="flex items-center gap-2.5 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800/60 text-red-600 dark:text-red-400 text-sm rounded-2xl px-4 py-3 mb-5">
              <div className="w-1.5 h-1.5 rounded-full bg-red-500 shrink-0" />
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Email */}
            <div>
              <label className="block text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1.5">Email</label>
              <div className="relative">
                <Mail size={17} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
                <input
                  type="email"
                  value={form.email}
                  onChange={(e) => setForm({ ...form, email: e.target.value })}
                  className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-2xl pl-10 pr-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent focus:bg-white dark:focus:bg-gray-700/80 transition-all dark:text-gray-100"
                  placeholder="nama@email.com"
                  required
                />
              </div>
            </div>

            {/* Password */}
            <div>
              <label className="block text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1.5">Password</label>
              <div className="relative">
                <Lock size={17} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
                <input
                  type={showPassword ? 'text' : 'password'}
                  value={form.password}
                  onChange={(e) => setForm({ ...form, password: e.target.value })}
                  className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-2xl pl-10 pr-11 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent focus:bg-white dark:focus:bg-gray-700/80 transition-all dark:text-gray-100"
                  placeholder="••••••••"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3.5 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
                >
                  {showPassword ? <EyeOff size={17} /> : <Eye size={17} />}
                </button>
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white py-3 rounded-2xl font-semibold text-sm transition-all shadow-lg shadow-blue-600/25 disabled:opacity-60 disabled:cursor-not-allowed mt-1"
            >
              {loading ? (
                <span className="flex items-center justify-center gap-2">
                  <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  Memproses...
                </span>
              ) : 'Masuk'}
            </button>
          </form>

          {GOOGLE_ENABLED && (
            <>
              <div className="flex items-center gap-3 my-5">
                <div className="flex-1 h-px bg-gray-200 dark:bg-gray-700" />
                <span className="text-xs text-gray-400 dark:text-gray-500 font-medium">atau</span>
                <div className="flex-1 h-px bg-gray-200 dark:bg-gray-700" />
              </div>
              <div className="flex justify-center">
                <GoogleLogin
                  onSuccess={async (cr) => {
                    setError(''); setLoading(true);
                    try { login(await authGoogle(cr.credential)); }
                    catch (err) { setError(err?.response?.data?.error || 'Login dengan Google gagal'); }
                    finally { setLoading(false); }
                  }}
                  onError={() => setError('Login dengan Google gagal')}
                  theme="outline" size="large" shape="pill" text="signin_with"
                />
              </div>
            </>
          )}

          <p className="text-center mt-6 text-sm text-gray-500 dark:text-gray-400">
            Belum punya akun?{' '}
            <Link to="/register" className="text-blue-600 dark:text-blue-400 font-semibold hover:text-blue-700 transition-colors">
              Daftar gratis
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
