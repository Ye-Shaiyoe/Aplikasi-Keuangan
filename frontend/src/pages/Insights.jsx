import { useState, useEffect } from 'react';
import {
  ComposedChart, Bar, Line, XAxis, YAxis, CartesianGrid, Tooltip,
  ResponsiveContainer, Area, AreaChart, PieChart, Pie, Cell
} from 'recharts';
import { TrendingUp, TrendingDown, Wallet, Calendar, ChevronLeft, ChevronRight, Sparkles, ArrowUpRight, ArrowDownRight } from 'lucide-react';
import { getYearlyTrend, getCategoryTrend, getSummary } from '../api/client';

const MONTH_NAMES = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des'];
const MONTH_FULL = ['Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni', 'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'];

function formatRupiah(n) {
  return new Intl.NumberFormat('id-ID').format(n || 0);
}

function formatCompact(n) {
  if (Math.abs(n) >= 1000000) return `${(n / 1000000).toFixed(n % 1000000 === 0 ? 0 : 1)}jt`;
  if (Math.abs(n) >= 1000) return `${(n / 1000).toFixed(0)}rb`;
  return new Intl.NumberFormat('id-ID').format(n);
}

// Mini sparkline component
function MiniSparkline({ data, color }) {
  return (
    <ResponsiveContainer width={80} height={32}>
      <AreaChart data={data} margin={{ top: 2, right: 2, bottom: 2, left: 2 }}>
        <defs>
          <linearGradient id={`grad-${color.replace('#', '')}`} x1="0" y1="0" x2="0" y2="1">
            <stop offset="5%" stopColor={color} stopOpacity={0.3} />
            <stop offset="95%" stopColor={color} stopOpacity={0} />
          </linearGradient>
        </defs>
        <Area
          type="monotone"
          dataKey="expense"
          stroke={color}
          strokeWidth={2}
          fill={`url(#grad-${color.replace('#', '')})`}
          dot={false}
        />
      </AreaChart>
    </ResponsiveContainer>
  );
}

// Custom tooltip for the mixed chart
function CustomTooltip({ active, payload, label }) {
  if (!active || !payload?.length) return null;
  return (
    <div className="bg-white/95 backdrop-blur-xl rounded-xl shadow-xl border border-gray-100 p-3 text-sm">
      <p className="font-semibold text-gray-700 mb-1.5">{label}</p>
      {payload.map((entry, i) => (
        <div key={i} className="flex items-center gap-2 py-0.5">
          <div className="w-2.5 h-2.5 rounded-full" style={{ backgroundColor: entry.color }} />
          <span className="text-gray-500">{entry.name}:</span>
          <span className="font-semibold text-gray-800">Rp {formatCompact(entry.value)}</span>
        </div>
      ))}
    </div>
  );
}

export default function Insights() {
  const now = new Date();
  const [year, setYear] = useState(now.getFullYear());
  const [month, setMonth] = useState(now.getMonth() + 1);
  const [yearlyTrend, setYearlyTrend] = useState(null);
  const [catTrend, setCatTrend] = useState(null);
  const [summary, setSummary] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    Promise.all([
      getYearlyTrend(year),
      getCategoryTrend(month, year),
      getSummary(month, year),
    ]).then(([yt, ct, s]) => {
      setYearlyTrend(yt);
      setCatTrend(ct);
      setSummary(s);
      setLoading(false);
    }).catch(() => setLoading(false));
  }, [year, month]);

  // Prepare mixed chart data
  const chartData = yearlyTrend?.months?.map((m) => ({
    name: MONTH_NAMES[m.month - 1],
    Pemasukan: m.income,
    Pengeluaran: m.expense,
    Saldo: m.balance,
  })) || [];

  // Summary stats
  const totalIncome = yearlyTrend?.total_income || 0;
  const totalExpense = yearlyTrend?.total_expense || 0;
  const yearBalance = yearlyTrend?.balance || 0;
  const monthIncome = summary?.total_income || 0;
  const monthExpense = summary?.total_expense || 0;

  // Category data for table with income/expense split
  const incomeCats = summary?.by_category?.filter(c => c.total > 0 && !c.category_name.includes('Pengeluaran') && !c.category_name.includes('Lainnya (Peng'))
    .filter(c => !['Makanan & Minuman', 'Transportasi', 'Belanja', 'Hiburan', 'Kesehatan', 'Tagihan', 'Pendidikan', 'Tempat Tinggal'].includes(c.category_name))
    .sort((a, b) => b.total - a.total) || [];
  const expenseCats = summary?.by_category?.filter(c => c.total > 0 && (
    c.category_name.includes('Pengeluaran') ||
    ['Makanan & Minuman', 'Transportasi', 'Belanja', 'Hiburan', 'Kesehatan', 'Tagihan', 'Pendidikan', 'Tempat Tinggal'].includes(c.category_name)
  )).sort((a, b) => b.total - a.total) || [];

  // Build weekly trend map for sparklines
  const weeklyMap = {};
  catTrend?.categories?.forEach(c => {
    weeklyMap[c.category_id] = c.weekly;
  });

  const totalCatExpense = expenseCats.reduce((s, c) => s + c.total, 0);
  const totalCatIncome = incomeCats.reduce((s, c) => s + c.total, 0);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-xl sm:text-2xl font-bold text-gray-800 flex items-center gap-2">
            <Sparkles className="text-amber-500" size={24} />
            Insight
          </h1>
          <p className="text-sm text-gray-400 mt-0.5">Analisis keuangan & tren pengeluaran</p>
        </div>
        <div className="flex items-center gap-2">
          <select
            value={month}
            onChange={(e) => setMonth(parseInt(e.target.value))}
            className="border border-gray-200 rounded-xl px-3 py-2 text-sm bg-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
          >
            {MONTH_FULL.map((m, i) => <option key={i} value={i + 1}>{m}</option>)}
          </select>
          <div className="flex items-center bg-white border border-gray-200 rounded-xl overflow-hidden">
            <button onClick={() => setYear(year - 1)} className="px-2.5 py-2 hover:bg-gray-50 text-gray-400">
              <ChevronLeft size={16} />
            </button>
            <span className="px-3 text-sm font-semibold text-gray-700 min-w-[60px] text-center">{year}</span>
            <button onClick={() => setYear(year + 1)} className="px-2.5 py-2 hover:bg-gray-50 text-gray-400">
              <ChevronRight size={16} />
            </button>
          </div>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 sm:gap-4">
        <div className="bg-white rounded-2xl border border-gray-100 p-4 sm:p-5 shadow-sm">
          <div className="flex items-center gap-2 mb-2">
            <div className="p-2 bg-green-50 rounded-lg"><ArrowUpRight size={16} className="text-green-500" /></div>
            <span className="text-xs text-gray-400">Pemasukan {year}</span>
          </div>
          <p className="text-base sm:text-xl font-bold text-green-600">Rp {formatCompact(totalIncome)}</p>
        </div>
        <div className="bg-white rounded-2xl border border-gray-100 p-4 sm:p-5 shadow-sm">
          <div className="flex items-center gap-2 mb-2">
            <div className="p-2 bg-red-50 rounded-lg"><ArrowDownRight size={16} className="text-red-500" /></div>
            <span className="text-xs text-gray-400">Pengeluaran {year}</span>
          </div>
          <p className="text-base sm:text-xl font-bold text-red-600">Rp {formatCompact(totalExpense)}</p>
        </div>
        <div className="bg-white rounded-2xl border border-gray-100 p-4 sm:p-5 shadow-sm">
          <div className="flex items-center gap-2 mb-2">
            <div className="p-2 bg-blue-50 rounded-lg"><Wallet size={16} className="text-blue-500" /></div>
            <span className="text-xs text-gray-400">Saldo {year}</span>
          </div>
          <p className={`text-base sm:text-xl font-bold ${yearBalance >= 0 ? 'text-blue-600' : 'text-red-600'}`}>
            Rp {formatCompact(yearBalance)}
          </p>
        </div>
        <div className="bg-gradient-to-br from-indigo-500 to-purple-600 rounded-2xl p-4 sm:p-5 shadow-sm text-white">
          <div className="flex items-center gap-2 mb-2">
            <div className="p-2 bg-white/20 rounded-lg"><Calendar size={16} className="text-white" /></div>
            <span className="text-xs text-white/70">Bulan ini</span>
          </div>
          <p className="text-base sm:text-xl font-bold">Rp {formatCompact(monthIncome - monthExpense)}</p>
          <p className="text-[10px] sm:text-xs text-white/60 mt-0.5">{MONTH_FULL[month - 1]}</p>
        </div>
      </div>

      {/* Mixed Chart - Yearly Trend */}
      <div className="bg-white rounded-2xl border border-gray-100 p-4 sm:p-6 shadow-sm">
        <div className="flex items-center justify-between mb-5">
          <h2 className="text-sm sm:text-lg font-semibold text-gray-800">Tren Keuangan Tahunan</h2>
          <div className="flex items-center gap-3 text-xs">
            <span className="flex items-center gap-1.5"><span className="w-3 h-3 rounded bg-green-400" /> Pemasukan</span>
            <span className="flex items-center gap-1.5"><span className="w-3 h-3 rounded bg-red-400" /> Pengeluaran</span>
            <span className="flex items-center gap-1.5 hidden sm:flex"><span className="w-3 h-1.5 rounded-full bg-blue-500" /> Saldo</span>
          </div>
        </div>
        {chartData.some(d => d.Pemasukan > 0 || d.Pengeluaran > 0) ? (
          <ResponsiveContainer width="100%" height={300}>
            <ComposedChart data={chartData} margin={{ top: 5, right: 5, bottom: 5, left: 5 }}>
              <defs>
                <linearGradient id="incomeGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stopColor="#10b981" stopOpacity={0.9} />
                  <stop offset="100%" stopColor="#10b981" stopOpacity={0.6} />
                </linearGradient>
                <linearGradient id="expenseGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stopColor="#ef4444" stopOpacity={0.9} />
                  <stop offset="100%" stopColor="#ef4444" stopOpacity={0.6} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" vertical={false} />
              <XAxis dataKey="name" tick={{ fontSize: 12 }} tickLine={false} axisLine={false} />
              <YAxis tick={{ fontSize: 11 }} tickFormatter={formatCompact} tickLine={false} axisLine={false} />
              <Tooltip content={<CustomTooltip />} />
              <Bar dataKey="Pemasukan" fill="url(#incomeGrad)" radius={[4, 4, 0, 0]} barSize={20} />
              <Bar dataKey="Pengeluaran" fill="url(#expenseGrad)" radius={[4, 4, 0, 0]} barSize={20} />
              <Line
                type="monotone"
                dataKey="Saldo"
                stroke="#6366f1"
                strokeWidth={2.5}
                dot={{ r: 3, fill: '#6366f1', strokeWidth: 0 }}
                activeDot={{ r: 5, strokeWidth: 2, stroke: '#fff' }}
              />
            </ComposedChart>
          </ResponsiveContainer>
        ) : (
          <div className="flex items-center justify-center h-[300px] text-gray-400 text-sm">
            Belum ada data untuk tahun {year}
          </div>
        )}
      </div>

      {/* Category Table */}
      <div className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
        <div className="p-4 sm:p-6 border-b border-gray-100">
          <h2 className="text-sm sm:text-lg font-semibold text-gray-800">
            Detail Pengeluaran & Pemasukan - {MONTH_FULL[month - 1]} {year}
          </h2>
        </div>

        {/* Income Section */}
        {incomeCats.length > 0 && (
          <div className="border-b border-gray-100">
            <div className="px-4 sm:px-6 py-3 bg-green-50/50 flex items-center gap-2">
              <ArrowUpRight size={16} className="text-green-500" />
              <span className="text-sm font-semibold text-green-700">Pemasukan</span>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-50 text-gray-400 text-xs">
                    <th className="text-left px-4 sm:px-6 py-2.5 font-medium">Kategori</th>
                    <th className="text-center px-3 py-2.5 font-medium hidden sm:table-cell">Transaksi</th>
                    <th className="text-right px-3 py-2.5 font-medium">Jumlah</th>
                    <th className="text-right px-3 py-2.5 font-medium hidden sm:table-cell">Persentase</th>
                    <th className="text-right px-3 sm:px-6 py-2.5 font-medium hidden md:table-cell">Tren Mingguan</th>
                  </tr>
                </thead>
                <tbody>
                  {incomeCats.map((c) => {
                    const pct = totalCatIncome > 0 ? Math.round((c.total / totalCatIncome) * 100) : 0;
                    const weekly = weeklyMap[c.category_id] || [{ expense: 0 }, { expense: 0 }, { expense: 0 }, { expense: 0 }];
                    const sparkData = weekly.map(w => ({ expense: w.expense }));
                    return (
                      <tr key={c.category_id} className="border-b border-gray-50 hover:bg-gray-50/50 transition-colors">
                        <td className="px-4 sm:px-6 py-3">
                          <div className="flex items-center gap-2.5">
                            <div className="w-3 h-3 rounded-full shrink-0" style={{ backgroundColor: c.category_color }} />
                            <span className="font-medium text-gray-700">{c.category_name}</span>
                          </div>
                        </td>
                        <td className="text-center px-3 py-3 text-gray-400 hidden sm:table-cell">{c.count}</td>
                        <td className="text-right px-3 py-3 font-semibold text-green-600">Rp {formatRupiah(c.total)}</td>
                        <td className="text-right px-3 py-3 hidden sm:table-cell">
                          <div className="flex items-center justify-end gap-2">
                            <div className="w-16 h-2 bg-gray-100 rounded-full overflow-hidden">
                              <div className="h-full rounded-full bg-green-400" style={{ width: `${pct}%` }} />
                            </div>
                            <span className="text-xs text-gray-400 w-8">{pct}%</span>
                          </div>
                        </td>
                        <td className="text-right px-3 sm:px-6 py-3 hidden md:table-cell">
                          <MiniSparkline data={sparkData} color={c.category_color} />
                        </td>
                      </tr>
                    );
                  })}
                  <tr className="bg-green-50/30 font-semibold">
                    <td className="px-4 sm:px-6 py-3 text-green-700">Total Pemasukan</td>
                    <td className="text-center px-3 py-3 text-green-600 hidden sm:table-cell">{incomeCats.reduce((s, c) => s + c.count, 0)}</td>
                    <td className="text-right px-3 py-3 text-green-700">Rp {formatRupiah(totalCatIncome)}</td>
                    <td className="text-right px-3 py-3 text-green-600 hidden sm:table-cell">100%</td>
                    <td className="px-3 sm:px-6 py-3 hidden md:table-cell" />
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        )}

        {/* Expense Section */}
        {expenseCats.length > 0 && (
          <div>
            <div className="px-4 sm:px-6 py-3 bg-red-50/50 flex items-center gap-2">
              <ArrowDownRight size={16} className="text-red-500" />
              <span className="text-sm font-semibold text-red-700">Pengeluaran</span>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-50 text-gray-400 text-xs">
                    <th className="text-left px-4 sm:px-6 py-2.5 font-medium">Kategori</th>
                    <th className="text-center px-3 py-2.5 font-medium hidden sm:table-cell">Transaksi</th>
                    <th className="text-right px-3 py-2.5 font-medium">Jumlah</th>
                    <th className="text-right px-3 py-2.5 font-medium hidden sm:table-cell">Persentase</th>
                    <th className="text-right px-3 sm:px-6 py-2.5 font-medium hidden md:table-cell">Tren Mingguan</th>
                  </tr>
                </thead>
                <tbody>
                  {expenseCats.map((c) => {
                    const pct = totalCatExpense > 0 ? Math.round((c.total / totalCatExpense) * 100) : 0;
                    const weekly = weeklyMap[c.category_id] || [{ expense: 0 }, { expense: 0 }, { expense: 0 }, { expense: 0 }];
                    const sparkData = weekly.map(w => ({ expense: w.expense }));
                    return (
                      <tr key={c.category_id} className="border-b border-gray-50 hover:bg-gray-50/50 transition-colors">
                        <td className="px-4 sm:px-6 py-3">
                          <div className="flex items-center gap-2.5">
                            <div className="w-3 h-3 rounded-full shrink-0" style={{ backgroundColor: c.category_color }} />
                            <span className="font-medium text-gray-700">{c.category_name}</span>
                          </div>
                        </td>
                        <td className="text-center px-3 py-3 text-gray-400 hidden sm:table-cell">{c.count}</td>
                        <td className="text-right px-3 py-3 font-semibold text-red-600">Rp {formatRupiah(c.total)}</td>
                        <td className="text-right px-3 py-3 hidden sm:table-cell">
                          <div className="flex items-center justify-end gap-2">
                            <div className="w-16 h-2 bg-gray-100 rounded-full overflow-hidden">
                              <div className="h-full rounded-full bg-red-400" style={{ width: `${pct}%` }} />
                            </div>
                            <span className="text-xs text-gray-400 w-8">{pct}%</span>
                          </div>
                        </td>
                        <td className="text-right px-3 sm:px-6 py-3 hidden md:table-cell">
                          <MiniSparkline data={sparkData} color={c.category_color} />
                        </td>
                      </tr>
                    );
                  })}
                  <tr className="bg-red-50/30 font-semibold">
                    <td className="px-4 sm:px-6 py-3 text-red-700">Total Pengeluaran</td>
                    <td className="text-center px-3 py-3 text-red-600 hidden sm:table-cell">{expenseCats.reduce((s, c) => s + c.count, 0)}</td>
                    <td className="text-right px-3 py-3 text-red-700">Rp {formatRupiah(totalCatExpense)}</td>
                    <td className="text-right px-3 py-3 text-red-600 hidden sm:table-cell">100%</td>
                    <td className="px-3 sm:px-6 py-3 hidden md:table-cell" />
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        )}

        {incomeCats.length === 0 && expenseCats.length === 0 && (
          <div className="flex items-center justify-center py-16 text-gray-400 text-sm">
            Belum ada data transaksi untuk bulan ini
          </div>
        )}
      </div>

      {/* Monthly Expense Distribution - Small pie */}
      {expenseCats.length > 0 && (
        <div className="bg-white rounded-2xl border border-gray-100 p-4 sm:p-6 shadow-sm">
          <h2 className="text-sm sm:text-lg font-semibold text-gray-800 mb-4">Distribusi Pengeluaran</h2>
          <div className="flex flex-col sm:flex-row items-center gap-6">
            <ResponsiveContainer width={180} height={180}>
              <PieChart>
                <Pie
                  data={expenseCats.map(c => ({ name: c.category_name, value: c.total, color: c.category_color }))}
                  cx="50%"
                  cy="50%"
                  innerRadius={50}
                  outerRadius={80}
                  paddingAngle={2}
                  dataKey="value"
                >
                  {expenseCats.map((c, i) => <Cell key={i} fill={c.category_color} />)}
                </Pie>
                <Tooltip formatter={(v) => `Rp ${formatRupiah(v)}`} />
              </PieChart>
            </ResponsiveContainer>
            <div className="flex-1 grid grid-cols-2 gap-x-4 gap-y-2 w-full">
              {expenseCats.slice(0, 8).map(c => (
                <div key={c.category_id} className="flex items-center gap-2 text-sm">
                  <div className="w-2.5 h-2.5 rounded-full shrink-0" style={{ backgroundColor: c.category_color }} />
                  <span className="text-gray-600 truncate flex-1">{c.category_name}</span>
                  <span className="text-gray-400 text-xs">{totalCatExpense > 0 ? Math.round((c.total / totalCatExpense) * 100) : 0}%</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
