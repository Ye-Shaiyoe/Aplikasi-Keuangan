import { useState, useEffect } from 'react';
import {
  Brain, TrendingUp, Sparkles, RefreshCw, AlertCircle, CheckCircle2,
  ChevronRight, Activity, Cpu, BarChart2,
} from 'lucide-react';
import {
  AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend,
} from 'recharts';
import { mlForecast, mlTrain, mlHealth, mlPredictCategory } from '../api/client';

// ── helpers ──────────────────────────────────────────────────────────────────
const fmt = (v) => new Intl.NumberFormat('id-ID').format(v || 0);
const fmtCompact = (v) => {
  if (Math.abs(v) >= 1_000_000) return `${(v / 1_000_000).toFixed(1)}jt`;
  if (Math.abs(v) >= 1_000) return `${(v / 1_000).toFixed(0)}rb`;
  return fmt(v);
};

const confidenceColor = (level) =>
  level === 'Tinggi'
    ? 'text-green-600 dark:text-green-400 bg-green-50 dark:bg-green-900/30'
    : level === 'Sedang'
    ? 'text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/30'
    : 'text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/30';

// ── component ─────────────────────────────────────────────────────────────────
export default function MLInsights() {
  const [health, setHealth] = useState(null);
  const [forecast, setForecast] = useState(null);
  const [training, setTraining] = useState(false);
  const [trainMsg, setTrainMsg] = useState(null);
  const [loadingForecast, setLoadingForecast] = useState(false);
  const [loadingHealth, setLoadingHealth] = useState(true);
  const [error, setError] = useState(null);

  // demo predict-category
  const [demoText, setDemoText] = useState('');
  const [demoResult, setDemoResult] = useState(null);
  const [predicting, setPredicting] = useState(false);

  // ── fetch health + forecast on mount ──────────────────────────────────────
  useEffect(() => {
    fetchAll();
  }, []);

  async function fetchAll() {
    setLoadingHealth(true);
    setError(null);
    try {
      const h = await mlHealth();
      setHealth(h);
    } catch (e) {
      const status = e?.response?.status;
      if (!status || e?.code === 'ERR_NETWORK') {
        setHealth({ _offline: true });
      } else {
        setHealth(null);
      }
    } finally {
      setLoadingHealth(false);
    }

    setLoadingForecast(true);
    try {
      const f = await mlForecast();
      setForecast(f);
    } catch (e) {
      const detail = e?.response?.data?.detail || e?.response?.data?.error;
      const status = e?.response?.status;
      if (!status || status === 0 || e?.code === 'ERR_NETWORK') {
        setError('Layanan ML tidak aktif. Jalankan: uvicorn main:app --reload --port 8000');
      } else {
        setError(detail || e?.message || 'Gagal memuat forecast');
      }
    } finally {
      setLoadingForecast(false);
    }
  }

  // ── train model ───────────────────────────────────────────────────────────
  async function handleTrain() {
    setTraining(true);
    setTrainMsg(null);
    try {
      const res = await mlTrain();
      setTrainMsg({ ok: true, text: res.message || 'Model berhasil di-train!' });
      // refresh health after training
      const h = await mlHealth();
      setHealth(h);
    } catch (e) {
      const msg = e?.response?.data?.detail || e?.message || 'Training gagal';
      setTrainMsg({ ok: false, text: msg });
    } finally {
      setTraining(false);
    }
  }

  // ── demo predict ──────────────────────────────────────────────────────────
  async function handlePredict() {
    if (!demoText.trim()) return;
    setPredicting(true);
    setDemoResult(null);
    try {
      const res = await mlPredictCategory(demoText.trim());
      setDemoResult(res);
    } catch (e) {
      const msg = e?.response?.data?.detail || e?.message || 'Prediksi gagal';
      setDemoResult({ error: msg });
    } finally {
      setPredicting(false);
    }
  }

  // ── chart data from projections ───────────────────────────────────────────
  const chartData = forecast?.projections?.map((p, i) => ({
    name: `Bulan +${p.month_offset}`,
    Pengeluaran: p.forecast_expense,
    Pemasukan: p.forecast_income,
    Saldo: p.forecast_balance,
  })) || [];

  // ── render ────────────────────────────────────────────────────────────────
  return (
    <div className="space-y-5 sm:space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl sm:text-2xl font-bold text-gray-800 dark:text-gray-100 flex items-center gap-2">
            <Brain size={24} className="text-indigo-500" />
            ML Insights
          </h1>
          <p className="text-xs text-gray-400 dark:text-gray-500 mt-0.5">
            Prediksi kategori &amp; forecast keuangan berbasis Machine Learning
          </p>
        </div>
        <button
          onClick={fetchAll}
          className="flex items-center gap-1.5 text-xs text-gray-500 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors px-3 py-1.5 rounded-xl border border-gray-200 dark:border-gray-700 hover:border-blue-300"
        >
          <RefreshCw size={14} />
          Refresh
        </button>
      </div>

      {/* Service status */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {/* Health card */}
        <div className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 p-4 sm:p-5">
          <div className="flex items-center gap-2 mb-3">
            <Activity size={18} className="text-indigo-500" />
            <span className="text-sm font-semibold text-gray-700 dark:text-gray-200">Status Layanan ML</span>
          </div>
          {loadingHealth ? (
            <div className="flex items-center gap-2 text-sm text-gray-400">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-indigo-500" />
              Memeriksa...
            </div>
          ) : health ? (
            <div className="space-y-2">
              {health._offline ? (
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-red-500" />
                  <span className="text-sm text-red-500 dark:text-red-400 font-medium">Offline</span>
                  <span className="text-xs text-gray-400">· ML service tidak berjalan</span>
                </div>
              ) : (
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
                  <span className="text-sm text-green-600 dark:text-green-400 font-medium">Online</span>
                  <span className="text-xs text-gray-400">· {health.service}</span>
                </div>
              )}
              {!health._offline && health.models && (
                <div className="mt-2 space-y-1">
                  {Object.entries(health.models).map(([name, status]) => (
                    <div key={name} className="flex items-center justify-between text-xs">
                      <span className="text-gray-500 dark:text-gray-400">{name}</span>
                      <span className={`px-2 py-0.5 rounded-full font-medium ${
                        status === 'ready'
                          ? 'bg-green-50 text-green-600 dark:bg-green-900/30 dark:text-green-400'
                          : 'bg-amber-50 text-amber-600 dark:bg-amber-900/30 dark:text-amber-400'
                      }`}>
                        {status === 'ready' ? 'Siap' : 'Belum dilatih'}
                      </span>
                    </div>
                  ))}
                </div>
              )}
              {health._offline && (
                <p className="text-xs text-gray-400 dark:text-gray-500 mt-1">
                  Jalankan: <code className="bg-gray-100 dark:bg-gray-700 px-1.5 py-0.5 rounded text-indigo-600 dark:text-indigo-400">uvicorn main:app --reload --port 8000</code>
                </p>
              )}
            </div>
          ) : (
            <div className="flex items-center gap-2 text-sm text-red-500 dark:text-red-400">
              <div className="w-2 h-2 rounded-full bg-red-500" />
              Layanan ML tidak tersedia
            </div>
          )}
        </div>

        {/* Train card */}
        <div className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 p-4 sm:p-5">
          <div className="flex items-center gap-2 mb-3">
            <Cpu size={18} className="text-indigo-500" />
            <span className="text-sm font-semibold text-gray-700 dark:text-gray-200">Training Model</span>
          </div>
          <p className="text-xs text-gray-400 dark:text-gray-500 mb-3 leading-relaxed">
            Latih ulang model kategorisasi menggunakan data transaksimu. Minimal 5 transaksi dengan deskripsi diperlukan.
          </p>
          <button
            onClick={handleTrain}
            disabled={training}
            className="flex items-center gap-2 bg-indigo-600 text-white text-sm px-4 py-2 rounded-xl hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
          >
            {training ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white" />
                Melatih...
              </>
            ) : (
              <>
                <Sparkles size={16} />
                Train Sekarang
              </>
            )}
          </button>
          {trainMsg && (
            <div className={`mt-2.5 flex items-start gap-2 text-xs px-3 py-2 rounded-xl ${
              trainMsg.ok
                ? 'bg-green-50 dark:bg-green-900/30 text-green-700 dark:text-green-400'
                : 'bg-red-50 dark:bg-red-900/30 text-red-600 dark:text-red-400'
            }`}>
              {trainMsg.ok
                ? <CheckCircle2 size={14} className="mt-0.5 shrink-0" />
                : <AlertCircle size={14} className="mt-0.5 shrink-0" />}
              {trainMsg.text}
            </div>
          )}
        </div>
      </div>

      {/* Demo: predict category */}
      <div className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 p-4 sm:p-5">
        <div className="flex items-center gap-2 mb-3">
          <Sparkles size={18} className="text-indigo-500" />
          <span className="text-sm font-semibold text-gray-700 dark:text-gray-200">Coba Prediksi Kategori</span>
        </div>
        <p className="text-xs text-gray-400 dark:text-gray-500 mb-3">
          Ketik deskripsi transaksi, lalu lihat prediksi kategori dari model ML.
          (Model harus sudah di-train terlebih dahulu.)
        </p>
        <div className="flex gap-2">
          <input
            type="text"
            value={demoText}
            onChange={(e) => setDemoText(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handlePredict()}
            placeholder="Contoh: makan siang di warteg, beli bensin, bayar listrik..."
            className="flex-1 border border-gray-200 dark:border-gray-600 rounded-xl px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 bg-white dark:bg-gray-700 dark:text-gray-200"
          />
          <button
            onClick={handlePredict}
            disabled={predicting || !demoText.trim()}
            className="flex items-center gap-1.5 bg-indigo-600 text-white text-sm px-4 py-2.5 rounded-xl hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
          >
            {predicting ? (
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white" />
            ) : (
              <ChevronRight size={16} />
            )}
            <span className="hidden sm:inline">Prediksi</span>
          </button>
        </div>

        {demoResult && (
          <div className="mt-3">
            {demoResult.error ? (
              <div className="flex items-center gap-2 text-sm text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/30 px-3 py-2.5 rounded-xl">
                <AlertCircle size={16} className="shrink-0" />
                {demoResult.error}
              </div>
            ) : (
              <div className="space-y-2">
                <div className="flex items-center justify-between bg-indigo-50 dark:bg-indigo-900/30 rounded-xl px-4 py-3">
                  <div>
                    <p className="text-xs text-gray-400 dark:text-gray-500 mb-0.5">Prediksi Terbaik</p>
                    <p className="text-sm font-semibold text-indigo-700 dark:text-indigo-300">
                      {demoResult.predicted_category_name}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="text-xs text-gray-400 dark:text-gray-500 mb-0.5">Confidence</p>
                    <p className="text-sm font-bold text-indigo-600 dark:text-indigo-400">
                      {(demoResult.confidence * 100).toFixed(1)}%
                    </p>
                  </div>
                </div>
                {demoResult.top_predictions?.length > 1 && (
                  <div>
                    <p className="text-xs text-gray-400 dark:text-gray-500 mb-1.5">Top Prediksi</p>
                    <div className="grid grid-cols-3 gap-2">
                      {demoResult.top_predictions.map((p, i) => (
                        <div
                          key={i}
                          className={`rounded-xl px-3 py-2 text-xs text-center border ${
                            i === 0
                              ? 'border-indigo-200 bg-indigo-50 dark:bg-indigo-900/30 dark:border-indigo-800'
                              : 'border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50'
                          }`}
                        >
                          <p className={`font-medium truncate ${i === 0 ? 'text-indigo-700 dark:text-indigo-300' : 'text-gray-600 dark:text-gray-300'}`}>
                            {p.category_name}
                          </p>
                          <p className="text-gray-400 dark:text-gray-500 mt-0.5">
                            {(p.confidence * 100).toFixed(1)}%
                          </p>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        )}
      </div>

      {/* Forecast section */}
      <div className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 p-4 sm:p-5">
        <div className="flex items-center gap-2 mb-4">
          <TrendingUp size={18} className="text-indigo-500" />
          <span className="text-sm font-semibold text-gray-700 dark:text-gray-200">Forecast 3 Bulan ke Depan</span>
        </div>

        {loadingForecast ? (
          <div className="flex items-center justify-center h-48">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500" />
          </div>
        ) : error ? (
          <div className="flex items-center gap-2 text-sm text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/30 px-4 py-3 rounded-xl">
            <AlertCircle size={16} className="shrink-0" />
            {error}
          </div>
        ) : forecast ? (
          <div className="space-y-5">
            {/* Warn if not enough data */}
            {forecast.data_points < 3 && (
              <div className="flex items-start gap-2 text-xs bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 text-amber-700 dark:text-amber-400 px-4 py-3 rounded-xl">
                <AlertCircle size={14} className="mt-0.5 shrink-0" />
                <span>
                  Data historis masih kurang ({forecast.data_points} bulan). Butuh minimal 3 bulan untuk forecast akurat.
                  Tambah lebih banyak transaksi untuk hasil yang lebih baik.
                </span>
              </div>
            )}
            {/* Confidence badge + meta */}
            <div className="flex flex-wrap items-center gap-2 text-xs">
              <span className={`px-2.5 py-1 rounded-full font-medium ${confidenceColor(forecast.confidence?.level)}`}>
                Akurasi: {forecast.confidence?.level}
              </span>
              <span className="text-gray-400 dark:text-gray-500">
                {forecast.data_points} bulan data · {forecast.model_type}
              </span>
              <span className="text-gray-400 dark:text-gray-500">
                R² pengeluaran: {(forecast.confidence?.expense_r2 * 100).toFixed(1)}% · R² pemasukan: {(forecast.confidence?.income_r2 * 100).toFixed(1)}%
              </span>
            </div>

            {/* Next-month spotlight cards */}
            <div className="grid grid-cols-3 gap-3">
              <div className="rounded-xl bg-red-50 dark:bg-red-900/20 border border-red-100 dark:border-red-900/40 px-3 py-3 text-center">
                <p className="text-[11px] text-gray-400 dark:text-gray-500 mb-1">Prediksi Keluar</p>
                <p className="text-sm sm:text-base font-bold text-red-600 dark:text-red-400">
                  Rp {fmtCompact(forecast.forecast_expense)}
                </p>
              </div>
              <div className="rounded-xl bg-green-50 dark:bg-green-900/20 border border-green-100 dark:border-green-900/40 px-3 py-3 text-center">
                <p className="text-[11px] text-gray-400 dark:text-gray-500 mb-1">Prediksi Masuk</p>
                <p className="text-sm sm:text-base font-bold text-green-600 dark:text-green-400">
                  Rp {fmtCompact(forecast.forecast_income)}
                </p>
              </div>
              <div className="rounded-xl bg-blue-50 dark:bg-blue-900/20 border border-blue-100 dark:border-blue-900/40 px-3 py-3 text-center">
                <p className="text-[11px] text-gray-400 dark:text-gray-500 mb-1">Prediksi Saldo</p>
                <p className={`text-sm sm:text-base font-bold ${
                  forecast.forecast_balance >= 0 ? 'text-blue-600 dark:text-blue-400' : 'text-red-600 dark:text-red-400'
                }`}>
                  Rp {fmtCompact(forecast.forecast_balance)}
                </p>
              </div>
            </div>

            {/* Area chart */}
            {chartData.length > 0 && (
              <div>
                <p className="text-xs text-gray-400 dark:text-gray-500 mb-2">Proyeksi 3 bulan</p>
                <ResponsiveContainer width="100%" height={220}>
                  <AreaChart data={chartData} margin={{ top: 5, right: 10, left: 0, bottom: 0 }}>
                    <defs>
                      <linearGradient id="colorExp" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#ef4444" stopOpacity={0.15} />
                        <stop offset="95%" stopColor="#ef4444" stopOpacity={0} />
                      </linearGradient>
                      <linearGradient id="colorInc" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#22c55e" stopOpacity={0.15} />
                        <stop offset="95%" stopColor="#22c55e" stopOpacity={0} />
                      </linearGradient>
                      <linearGradient id="colorBal" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#6366f1" stopOpacity={0.15} />
                        <stop offset="95%" stopColor="#6366f1" stopOpacity={0} />
                      </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" stroke="var(--tw-border-opacity, #f3f4f6)" opacity={0.5} />
                    <XAxis dataKey="name" tick={{ fontSize: 11 }} tickLine={false} axisLine={false} />
                    <YAxis
                      tick={{ fontSize: 10 }}
                      tickLine={false}
                      axisLine={false}
                      tickFormatter={(v) => fmtCompact(v)}
                      width={52}
                    />
                    <Tooltip
                      formatter={(value, name) => [`Rp ${fmt(value)}`, name]}
                      contentStyle={{
                        backgroundColor: 'var(--tooltip-bg, #fff)',
                        border: '1px solid #e5e7eb',
                        borderRadius: '12px',
                        fontSize: 12,
                      }}
                    />
                    <Legend wrapperStyle={{ fontSize: 12 }} />
                    <Area type="monotone" dataKey="Pengeluaran" stroke="#ef4444" fill="url(#colorExp)" strokeWidth={2} dot={{ r: 4 }} />
                    <Area type="monotone" dataKey="Pemasukan" stroke="#22c55e" fill="url(#colorInc)" strokeWidth={2} dot={{ r: 4 }} />
                    <Area type="monotone" dataKey="Saldo" stroke="#6366f1" fill="url(#colorBal)" strokeWidth={2} dot={{ r: 4 }} />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
            )}

            {/* Projections table */}
            <div>
              <p className="text-xs text-gray-400 dark:text-gray-500 mb-2">Detail proyeksi</p>
              <div className="overflow-x-auto rounded-xl border border-gray-100 dark:border-gray-700">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="bg-gray-50 dark:bg-gray-700/50">
                      <th className="text-left px-4 py-2.5 text-xs font-medium text-gray-400 dark:text-gray-500">Bulan Offset</th>
                      <th className="text-right px-4 py-2.5 text-xs font-medium text-gray-400 dark:text-gray-500">Pengeluaran</th>
                      <th className="text-right px-4 py-2.5 text-xs font-medium text-gray-400 dark:text-gray-500">Pemasukan</th>
                      <th className="text-right px-4 py-2.5 text-xs font-medium text-gray-400 dark:text-gray-500">Saldo</th>
                    </tr>
                  </thead>
                  <tbody>
                    {forecast.projections?.map((p) => (
                      <tr key={p.month_offset} className="border-t border-gray-50 dark:border-gray-700/50">
                        <td className="px-4 py-2.5 text-gray-600 dark:text-gray-300">+{p.month_offset} bulan</td>
                        <td className="px-4 py-2.5 text-right text-red-600 dark:text-red-400">Rp {fmt(p.forecast_expense)}</td>
                        <td className="px-4 py-2.5 text-right text-green-600 dark:text-green-400">Rp {fmt(p.forecast_income)}</td>
                        <td className={`px-4 py-2.5 text-right font-semibold ${
                          p.forecast_balance >= 0 ? 'text-blue-600 dark:text-blue-400' : 'text-red-600 dark:text-red-400'
                        }`}>
                          Rp {fmt(p.forecast_balance)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        ) : null}
      </div>

      {/* Info footer */}
      <p className="text-[11px] text-gray-400 dark:text-gray-500 text-center pb-2">
        Hasil prediksi bersifat estimasi. Akurasi meningkat seiring bertambahnya data historis.
      </p>
    </div>
  );
}
