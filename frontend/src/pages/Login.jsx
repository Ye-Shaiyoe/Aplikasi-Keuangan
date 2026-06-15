import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { authLogin } from '../api/client';
import { Wallet, Mail, Lock, Eye, EyeOff } from 'lucide-react';

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
      setError(err?.response?.data?.error || 'Login gagal');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex">
      {/* Left Panel - Branding */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-blue-600 via-blue-700 to-indigo-800 items-center justify-center p-12">
        <div className="max-w-md text-white">
          <div className="flex items-center gap-3 mb-8">
            <div className="p-3 bg-white/20 rounded-2xl backdrop-blur-sm">
              <Wallet size={36} />
            </div>
            <h1 className="text-3xl font-bold">Catatan Keuangan</h1>
          </div>
          <p className="text-blue-100 text-lg leading-relaxed">
            Kelola keuangan pribadimu dengan mudah. Catat pemasukan, pengeluaran, dan lihat laporan keuanganmu dalam satu aplikasi.
          </p>
          <div className="mt-12 flex gap-6">
            <div className="text-center">
              <div className="text-3xl font-bold">100%</div>
              <div className="text-blue-200 text-sm">Gratis</div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold">24/7</div>
              <div className="text-blue-200 text-sm">Akses</div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold">100%</div>
              <div className="text-blue-200 text-sm">Privasi</div>
            </div>
          </div>
        </div>
      </div>

      {/* Right Panel - Login Form */}
      <div className="flex-1 flex items-center justify-center p-4 sm:p-6 bg-gray-50">
        <div className="w-full max-w-md">
          {/* Mobile logo */}
          <div className="lg:hidden flex items-center gap-2 mb-6 sm:mb-8 justify-center">
            <div className="p-2 bg-blue-600 rounded-xl text-white">
              <Wallet size={24} className="sm:hidden" />
              <Wallet size={28} className="hidden sm:block" />
            </div>
            <h1 className="text-xl sm:text-2xl font-bold text-gray-800">Catatan Keuangan</h1>
          </div>

          <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-5 sm:p-8">
            <h2 className="text-xl sm:text-2xl font-bold text-gray-800 mb-1">Selamat datang kembali</h2>
            <p className="text-sm sm:text-base text-gray-500 mb-6 sm:mb-8">Masuk ke akunmu untuk melanjutkan</p>

            {error && (
              <div className="bg-red-50 border border-red-200 text-red-600 text-sm rounded-xl px-4 py-3 mb-4">
                {error}
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-5">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">Email</label>
                <div className="relative">
                  <Mail size={18} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
                  <input
                    type="email"
                    value={form.email}
                    onChange={(e) => setForm({ ...form, email: e.target.value })}
                    className="w-full border border-gray-200 rounded-xl pl-11 pr-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="nama@email.com"
                    required
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">Password</label>
                <div className="relative">
                  <Lock size={18} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
                  <input
                    type={showPassword ? 'text' : 'password'}
                    value={form.password}
                    onChange={(e) => setForm({ ...form, password: e.target.value })}
                    className="w-full border border-gray-200 rounded-xl pl-11 pr-11 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Masukkan password"
                    required
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

              <button
                type="submit"
                disabled={loading}
                className="w-full bg-blue-600 text-white py-3 rounded-xl hover:bg-blue-700 disabled:opacity-50 transition-colors text-sm font-semibold"
              >
                {loading ? 'Memproses...' : 'Masuk'}
              </button>
            </form>
          </div>

          <p className="text-center mt-6 text-sm text-gray-500">
            Belum punya akun?{' '}
            <Link to="/register" className="text-blue-600 font-medium hover:text-blue-700">
              Daftar sekarang
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
