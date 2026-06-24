import { useState, useEffect } from 'react';
import {
  AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
} from 'recharts';
import {
  Sparkles, TrendingUp, Wallet, Percent, AlertCircle, CheckCircle2,
  Activity, ArrowUpRight, ArrowDownRight
} from 'lucide-react';
import Card from '../components/Card';
import { getAdvancedAnalytics } from '../api/client';
import { formatRupiah, formatCompact } from '../utils/format';

export default function AdvancedAnalytics() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    getAdvancedAnalytics()
      .then((res) => setData(res))
      .catch((err) => console.error(err))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
      </div>
    );
  }

  if (!data) {
    return (
      <div className="text-center py-12 text-gray-500">
        Gagal memuat analisis data keuangan. Silakan coba lagi nanti.
      </div>
    );
  }

  // Circular gauge color matching health score
  const getScoreColor = (score) => {
    if (score >= 80) return 'stroke-emerald-500 text-emerald-600 dark:text-emerald-400';
    if (score >= 60) return 'stroke-teal-500 text-teal-600 dark:text-teal-400';
    if (score >= 40) return 'stroke-amber-500 text-amber-600 dark:text-amber-400';
    return 'stroke-rose-500 text-rose-600 dark:text-rose-400';
  };

  const getScoreBg = (score) => {
    if (score >= 80) return 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950/30 dark:text-emerald-400';
    if (score >= 60) return 'bg-teal-50 text-teal-700 dark:bg-teal-950/30 dark:text-teal-400';
    if (score >= 40) return 'bg-amber-50 text-amber-700 dark:bg-amber-950/30 dark:text-amber-400';
    return 'bg-rose-50 text-rose-700 dark:bg-rose-950/30 dark:text-rose-400';
  };

  // SVG Gauge calculations
  const radius = 50;
  const circumference = 2 * Math.PI * radius;
  const strokeDashoffset = circumference - (data.health_score / 100) * circumference;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
        <div>
          <h1 className="text-xl sm:text-2xl font-bold text-gray-800 dark:text-gray-100 flex items-center gap-2">
            <Sparkles className="text-indigo-500 animate-pulse" size={24} />
            Analisis Finansial Cerdas
          </h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Analisis mendalam, metrik kesehatan, dan prediksi keuangan masa depan Anda.
          </p>
        </div>
      </div>

      {/* Main Grid: Health Score and Prediction */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Health Score Gauge */}
        <Card className="flex flex-col items-center justify-center p-6 text-center">
          <h3 className="text-sm font-semibold text-gray-500 dark:text-gray-400 mb-4">Skor Kesehatan Finansial</h3>
          <div className="relative flex items-center justify-center w-36 h-36">
            {/* Background Circle */}
            <svg className="w-full h-full transform -rotate-90">
              <circle
                cx="72"
                cy="72"
                r={radius}
                className="stroke-gray-100 dark:stroke-gray-800"
                strokeWidth="10"
                fill="transparent"
              />
              <circle
                cx="72"
                cy="72"
                r={radius}
                className={`transition-all duration-1000 ease-out ${getScoreColor(data.health_score)}`}
                strokeWidth="10"
                fill="transparent"
                strokeDasharray={circumference}
                strokeDashoffset={strokeDashoffset}
                strokeLinecap="round"
              />
            </svg>
            <div className="absolute flex flex-col items-center justify-center">
              <span className="text-3xl font-extrabold text-gray-800 dark:text-gray-100">{data.health_score}</span>
              <span className="text-[10px] text-gray-400 uppercase tracking-widest font-semibold">dari 100</span>
            </div>
          </div>
          <span className={`mt-4 px-3 py-1.5 rounded-full text-xs font-semibold tracking-wide ${getScoreBg(data.health_score)}`}>
            {data.health_rating}
          </span>
        </Card>

        {/* Prediction / Forecast Card */}
        <Card className="md:col-span-2 p-6 flex flex-col justify-between">
          <div>
            <h3 className="text-sm font-semibold text-gray-500 dark:text-gray-400 mb-3 flex items-center gap-2">
              <Activity size={18} className="text-indigo-500" />
              Proyeksi Keuangan Bulan Depan
            </h3>
            <p className="text-xs text-gray-400 mb-4">
              Estimasi berdasarkan tren pengeluaran dan pemasukan Anda selama beberapa bulan terakhir.
            </p>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div className="bg-emerald-50/50 dark:bg-emerald-950/10 p-4 rounded-xl border border-emerald-100/50 dark:border-emerald-900/30 flex items-center gap-3">
                <div className="p-2 bg-emerald-500 text-white rounded-lg">
                  <ArrowUpRight size={20} />
                </div>
                <div>
                  <p className="text-[10px] text-emerald-600 dark:text-emerald-400 uppercase font-semibold">Proyeksi Pemasukan</p>
                  <p className="text-lg font-bold text-emerald-700 dark:text-emerald-300">Rp {formatRupiah(data.forecast_income)}</p>
                </div>
              </div>

              <div className="bg-rose-50/50 dark:bg-rose-950/10 p-4 rounded-xl border border-rose-100/50 dark:border-rose-900/30 flex items-center gap-3">
                <div className="p-2 bg-rose-500 text-white rounded-lg">
                  <ArrowDownRight size={20} />
                </div>
                <div>
                  <p className="text-[10px] text-rose-600 dark:text-rose-400 uppercase font-semibold">Proyeksi Pengeluaran</p>
                  <p className="text-lg font-bold text-rose-700 dark:text-rose-300">Rp {formatRupiah(data.forecast_expense)}</p>
                </div>
              </div>
            </div>
          </div>

          <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-800 text-xs flex items-center gap-2 text-gray-500 dark:text-gray-400">
            <AlertCircle size={16} className="text-indigo-500 shrink-0" />
            <span>
              {data.forecast_income >= data.forecast_expense
                ? "Bagus! Proyeksi pemasukan Anda lebih besar dari perkiraan pengeluaran. Anda berpotensi menabung."
                : "Hati-hati: Perkiraan pengeluaran Anda melebihi perkiraan pemasukan. Segera tinjau kembali anggaran."}
            </span>
          </div>
        </Card>
      </div>

      {/* Metrics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Savings Rate Card */}
        <Card className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-semibold text-gray-500 dark:text-gray-400 flex items-center gap-2">
              <Percent size={18} className="text-blue-500" />
              Saving Rate Bulan Ini
            </h3>
            <span className={`text-lg font-bold ${data.savings_rate >= 10 ? 'text-emerald-600' : 'text-amber-500'}`}>
              {data.savings_rate.toFixed(1)}%
            </span>
          </div>
          <div className="w-full bg-gray-100 dark:bg-gray-800 rounded-full h-2.5 mb-3">
            <div
              className={`h-2.5 rounded-full transition-all duration-1000 ${data.savings_rate >= 20 ? 'bg-emerald-500' : data.savings_rate >= 10 ? 'bg-blue-500' : 'bg-amber-500'}`}
              style={{ width: `${Math.max(0, Math.min(100, data.savings_rate))}%` }}
            />
          </div>
          <p className="text-xs text-gray-400">
            Saving rate mengukur persentase pemasukan bersih yang berhasil Anda tabung. Angka ideal adalah di atas 10-20%.
          </p>
        </Card>

        {/* Budget Adherence Card */}
        <Card className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-semibold text-gray-500 dark:text-gray-400 flex items-center gap-2">
              <Wallet size={18} className="text-indigo-500" />
              Kepatuhan Anggaran (Budget Adherence)
            </h3>
            <span className={`text-lg font-bold ${data.budget_adherence >= 80 ? 'text-emerald-600' : 'text-rose-500'}`}>
              {data.budget_adherence.toFixed(0)}%
            </span>
          </div>
          <div className="w-full bg-gray-100 dark:bg-gray-800 rounded-full h-2.5 mb-3">
            <div
              className={`h-2.5 rounded-full transition-all duration-1000 ${data.budget_adherence >= 80 ? 'bg-emerald-500' : 'bg-rose-500'}`}
              style={{ width: `${data.budget_adherence}%` }}
            />
          </div>
          <p className="text-xs text-gray-400">
            Mengukur persentase kategori anggaran bulanan yang berhasil dijaga agar tidak melebihi batas yang ditentukan.
          </p>
        </Card>
      </div>

      {/* Historical Graph & Insights */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Historical Graph */}
        <Card className="lg:col-span-2 p-6">
          <h3 className="text-sm font-semibold text-gray-500 dark:text-gray-400 mb-4 flex items-center gap-2">
            <TrendingUp size={18} className="text-indigo-500" />
            Tren Keuangan Historis
          </h3>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={data.monthly_metrics} margin={{ top: 10, right: 10, left: -10, bottom: 0 }}>
                <defs>
                  <linearGradient id="analyticsInc" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#10b981" stopOpacity={0.2} />
                    <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
                  </linearGradient>
                  <linearGradient id="analyticsExp" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#ef4444" stopOpacity={0.2} />
                    <stop offset="95%" stopColor="#ef4444" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" vertical={false} />
                <XAxis dataKey="month_name" tick={{ fontSize: 10 }} tickLine={false} />
                <YAxis tick={{ fontSize: 10 }} tickFormatter={formatCompact} tickLine={false} />
                <Tooltip formatter={(value) => `Rp ${formatRupiah(value)}`} />
                <Area type="monotone" name="Pemasukan" dataKey="income" stroke="#10b981" fillOpacity={1} fill="url(#analyticsInc)" strokeWidth={2} />
                <Area type="monotone" name="Pengeluaran" dataKey="expense" stroke="#ef4444" fillOpacity={1} fill="url(#analyticsExp)" strokeWidth={2} />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </Card>

        {/* Actionable Insights Panel */}
        <Card className="p-6">
          <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-4 flex items-center gap-2 border-b border-gray-100 dark:border-gray-800 pb-3">
            <CheckCircle2 size={18} className="text-indigo-500" />
            Rekomendasi Keuangan
          </h3>
          <div className="space-y-4">
            {data.insights.map((insight, idx) => (
              <div key={idx} className="flex gap-2.5 items-start">
                <div className="w-1.5 h-1.5 rounded-full bg-indigo-500 shrink-0 mt-2" />
                <p className="text-xs sm:text-sm text-gray-600 dark:text-gray-300 leading-relaxed">
                  {insight}
                </p>
              </div>
            ))}
            {data.insights.length === 0 && (
              <p className="text-center py-6 text-xs text-gray-400">
                Data belum mencukupi untuk membuat analisis. Tambahkan transaksi Anda bulan ini!
              </p>
            )}
          </div>
        </Card>
      </div>
    </div>
  );
}
