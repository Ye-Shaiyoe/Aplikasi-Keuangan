import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { authRegister } from '../api/client';
import { Wallet, Mail, Lock, User, Eye, EyeOff } from 'lucide-react';

export default function Register() {
  const { register } = useAuth();
  const [form, setForm] = useState({ name: '', email: '', password: '', confirmPassword: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (form.password !== form.confirmPassword) {
      setError('Password tidak cocok');
      return;
    }
    if (form.password.length < 6) {
      setError('Password minimal 6 karakter');
      return;
    }

    setLoading(true);
    try {
      const { confirmPassword, ...payload } = form;
      const res = await authRegister(payload);
      register(res);
    } catch (err) {
      setError(err?.response?.data?.error || 'Registrasi gagal');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex">
      {/* Left Panel - Branding */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-emerald-500 via-teal-600 to-cyan-700 items-center justify-center p-12">
        <div className="max-w-md text-white">
          <div className="flex items-center gap-3 mb-8">
            <div className="p-3 bg-white/20 rounded-2xl backdrop-blur-sm">
              <Wallet size={36} />
            </div>
            <h1 className="text-3xl font-bold">Catatan Keuangan</h1>
          </div>
          <p className="text-emerald-100 text-lg leading-relaxed">
            Mulai perjalanan finansialmu hari ini. Buat akun gratis dan kelola keuanganmu dengan lebih cerdas.
          </p>
          <div className="mt-12 space-y-4">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center text-sm font-bold">1</div>
              <span className="text-emerald-100">Daftar akun gratis</span>
            </div>
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center text-sm font-bold">2</div>
              <span className="text-emerald-100">Catat transaksi harianmu</span>
            </div>
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center text-sm font-bold">3</div>
              <span className="text-emerald-100">Lihat laporan & analisis</span>
            </div>
          </div>
        </div>
      </div>

      {/* Right Panel - Register Form */}
      <div className="flex-1 flex items-center justify-center p-4 sm:p-6 bg-gray-50 dark:bg-gray-900 transition-colors">
        <div className="w-full max-w-md">
          {/* Mobile logo */}
          <div className="lg:hidden flex items-center gap-2 mb-6 sm:mb-8 justify-center">
            <div className="p-2 bg-emerald-600 rounded-xl text-white">
              <Wallet size={24} className="sm:hidden" />
              <Wallet size={28} className="hidden sm:block" />
            </div>
            <h1 className="text-xl sm:text-2xl font-bold text-gray-800 dark:text-gray-100">Catatan Keuangan</h1>
          </div>

          <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-5 sm:p-8 transition-colors">
            <h2 className="text-xl sm:text-2xl font-bold text-gray-800 dark:text-gray-100 mb-1">Buat akun baru</h2>
            <p className="text-sm sm:text-base text-gray-500 dark:text-gray-400 mb-6 sm:mb-8">Daftar gratis untuk memulai</p>

            {error && (
              <div className="bg-red-50 border border-red-200 text-red-600 text-sm rounded-xl px-4 py-3 mb-4">
                {error}
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Nama Lengkap</label>
                <div className="relative">
                  <User size={18} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
                  <input
                    type="text"
                    value={form.name}
                    onChange={(e) => setForm({ ...form, name: e.target.value })}
                    className="w-full border border-gray-200 dark:border-gray-600 rounded-xl pl-11 pr-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white dark:bg-gray-700 dark:text-gray-100"
                    placeholder="Nama lengkap"
                    required
                    minLength={2}
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Email</label>
                <div className="relative">
                  <Mail size={18} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
                  <input
                    type="email"
                    value={form.email}
                    onChange={(e) => setForm({ ...form, email: e.target.value })}
                    className="w-full border border-gray-200 dark:border-gray-600 rounded-xl pl-11 pr-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white dark:bg-gray-700 dark:text-gray-100"
                    placeholder="nama@email.com"
                    required
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Password</label>
                <div className="relative">
                  <Lock size={18} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
                  <input
                    type={showPassword ? 'text' : 'password'}
                    value={form.password}
                    onChange={(e) => setForm({ ...form, password: e.target.value })}
                    className="w-full border border-gray-200 dark:border-gray-600 rounded-xl pl-11 pr-11 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white dark:bg-gray-700 dark:text-gray-100"
                    placeholder="Minimal 6 karakter"
                    required
                    minLength={6}
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3.5 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                  >
                    {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                  </button>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Konfirmasi Password</label>
                <div className="relative">
                  <Lock size={18} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
                  <input
                    type={showPassword ? 'text' : 'password'}
                    value={form.confirmPassword}
                    onChange={(e) => setForm({ ...form, confirmPassword: e.target.value })}
                    className="w-full border border-gray-200 dark:border-gray-600 rounded-xl pl-11 pr-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white dark:bg-gray-700 dark:text-gray-100"
                    placeholder="Ulangi password"
                    required
                  />
                </div>
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full bg-blue-600 text-white py-3 rounded-xl hover:bg-blue-700 disabled:opacity-50 transition-colors text-sm font-semibold mt-2"
              >
                {loading ? 'Mendaftar...' : 'Daftar Sekarang'}
              </button>
            </form>
          </div>

          <p className="text-center mt-6 text-sm text-gray-500 dark:text-gray-400">
            Sudah punya akun?{' '}
            <Link to="/login" className="text-blue-600 font-medium hover:text-blue-700">
              Masuk di sini
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
