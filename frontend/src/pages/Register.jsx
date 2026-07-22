import { useState } from 'react';
import { Link } from 'react-router-dom';
import { GoogleLogin } from '@react-oauth/google';
import { useAuth } from '../context/AuthContext';
import { authRegister, authGoogle } from '../api/client';
import { Wallet, Mail, Lock, User, Eye, EyeOff, CheckCircle2 } from 'lucide-react';

const GOOGLE_ENABLED = !!import.meta.env.VITE_GOOGLE_CLIENT_ID;

const STEPS = [
  'Daftar akun gratis',
  'Catat transaksi harianmu',
  'Lihat laporan & analisis AI',
];

export default function Register() {
  const { register } = useAuth();
  const [form, setForm] = useState({ name: '', email: '', password: '', confirmPassword: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    if (form.password !== form.confirmPassword) { setError('Password tidak cocok'); return; }
    if (form.password.length < 6) { setError('Password minimal 6 karakter'); return; }
    setLoading(true);
    try {
      const { confirmPassword, ...payload } = form;
      register(await authRegister(payload));
    } catch (err) {
      setError(err?.response?.data?.error || 'Registrasi gagal');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex bg-gray-50 dark:bg-[#0f1117]">

      {/* ── Left branding panel ── */}
      <div className="hidden lg:flex lg:w-[45%] relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-br from-emerald-500 via-teal-600 to-cyan-700" />
        <div className="absolute -top-24 -left-24 w-96 h-96 bg-emerald-400/20 rounded-full blur-3xl" />
        <div className="absolute -bottom-24 -right-24 w-96 h-96 bg-cyan-500/20 rounded-full blur-3xl" />

        <div className="relative flex flex-col justify-center px-14 max-w-lg mx-auto">
          <div className="flex items-center gap-3 mb-10">
            <div className="p-2.5 bg-white/15 rounded-2xl backdrop-blur-sm border border-white/20">
              <Wallet size={28} className="text-white" />
            </div>
            <span className="text-2xl font-bold text-white tracking-tight">Catatan Keuangan</span>
          </div>

          <h2 className="text-4xl font-bold text-white leading-tight mb-4">
            Mulai perjalanan<br />finansialmu <span className="text-emerald-200">hari ini</span>
          </h2>
          <p className="text-emerald-100 text-base leading-relaxed mb-10">
            Gratis selamanya. Tanpa kartu kredit. Mulai dalam 30 detik.
          </p>

          <div className="space-y-5">
            {STEPS.map((text, i) => (
              <div key={i} className="flex items-center gap-4">
                <div className="w-8 h-8 rounded-full bg-white/15 border border-white/20 flex items-center justify-center text-white font-bold text-sm shrink-0">
                  {i + 1}
                </div>
                <span className="text-emerald-100 text-sm">{text}</span>
              </div>
            ))}
          </div>

          <div className="mt-12 p-5 bg-white/10 backdrop-blur-sm rounded-2xl border border-white/15">
            <div className="flex items-center gap-3">
              <CheckCircle2 size={20} className="text-emerald-300 shrink-0" />
              <p className="text-emerald-100 text-sm leading-relaxed">
                Datamu tersimpan aman di server pribadi — tidak dijual, tidak dibagikan.
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* ── Right form panel ── */}
      <div className="flex-1 flex items-center justify-center p-5 sm:p-8">
        <div className="w-full max-w-sm">

          {/* Mobile logo */}
          <div className="lg:hidden flex items-center gap-2.5 mb-8 justify-center">
            <div className="p-2 bg-emerald-600 rounded-xl text-white shadow-lg shadow-emerald-600/30">
              <Wallet size={22} />
            </div>
            <span className="text-lg font-bold text-gray-800 dark:text-gray-100">Catatan Keuangan</span>
          </div>

          <div className="mb-6">
            <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Buat akun baru ✨</h2>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">Daftar gratis, mulai hari ini</p>
          </div>

          {error && (
            <div className="flex items-center gap-2.5 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800/60 text-red-600 dark:text-red-400 text-sm rounded-2xl px-4 py-3 mb-5">
              <div className="w-1.5 h-1.5 rounded-full bg-red-500 shrink-0" />
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-3.5">
            {/* Name */}
            <div>
              <label className="block text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1.5">Nama Lengkap</label>
              <div className="relative">
                <User size={17} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-2xl pl-10 pr-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent focus:bg-white dark:focus:bg-gray-700/80 transition-all dark:text-gray-100"
                  placeholder="Nama lengkap"
                  required minLength={2}
                />
              </div>
            </div>

            {/* Email */}
            <div>
              <label className="block text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1.5">Email</label>
              <div className="relative">
                <Mail size={17} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
                <input
                  type="email"
                  value={form.email}
                  onChange={(e) => setForm({ ...form, email: e.target.value })}
                  className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-2xl pl-10 pr-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent focus:bg-white dark:focus:bg-gray-700/80 transition-all dark:text-gray-100"
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
                  className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-2xl pl-10 pr-11 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent focus:bg-white dark:focus:bg-gray-700/80 transition-all dark:text-gray-100"
                  placeholder="Min. 6 karakter"
                  required minLength={6}
                />
                <button type="button" onClick={() => setShowPassword(!showPassword)} className="absolute right-3.5 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors">
                  {showPassword ? <EyeOff size={17} /> : <Eye size={17} />}
                </button>
              </div>
            </div>

            {/* Confirm Password */}
            <div>
              <label className="block text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1.5">Konfirmasi Password</label>
              <div className="relative">
                <Lock size={17} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
                <input
                  type={showPassword ? 'text' : 'password'}
                  value={form.confirmPassword}
                  onChange={(e) => setForm({ ...form, confirmPassword: e.target.value })}
                  className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-2xl pl-10 pr-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent focus:bg-white dark:focus:bg-gray-700/80 transition-all dark:text-gray-100"
                  placeholder="Ulangi password"
                  required
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-emerald-600 hover:bg-emerald-700 active:bg-emerald-800 text-white py-3 rounded-2xl font-semibold text-sm transition-all shadow-lg shadow-emerald-600/25 disabled:opacity-60 disabled:cursor-not-allowed mt-1"
            >
              {loading ? (
                <span className="flex items-center justify-center gap-2">
                  <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  Mendaftar...
                </span>
              ) : 'Daftar Sekarang'}
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
                    try { register(await authGoogle(cr.credential)); }
                    catch (err) { setError(err?.response?.data?.error || 'Daftar dengan Google gagal'); }
                    finally { setLoading(false); }
                  }}
                  onError={() => setError('Daftar dengan Google gagal')}
                  theme="outline" size="large" shape="pill" text="signup_with"
                />
              </div>
            </>
          )}

          <p className="text-center mt-6 text-sm text-gray-500 dark:text-gray-400">
            Sudah punya akun?{' '}
            <Link to="/login" className="text-emerald-600 dark:text-emerald-400 font-semibold hover:text-emerald-700 transition-colors">
              Masuk di sini
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
