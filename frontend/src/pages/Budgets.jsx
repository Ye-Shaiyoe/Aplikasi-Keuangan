import { useState, useEffect } from 'react';
import { Wallet, Calendar, Plus, ChevronLeft, ChevronRight, Edit3, Trash2, Sparkles, AlertCircle, CheckCircle2, Clock, Target, TrendingUp } from 'lucide-react';
import Modal from '../components/Modal';
import { getBudgetSummary, upsertBudget, deleteBudget, getCategories } from '../api/client';

function formatRupiah(n) {
  return new Intl.NumberFormat('id-ID').format(n || 0);
}

function formatCompact(n) {
  if (Math.abs(n) >= 1000000) return `${(n / 1000000).toFixed(n % 1000000 === 0 ? 0 : 1)}jt`;
  if (Math.abs(n) >= 1000) return `${(n / 1000).toFixed(0)}rb`;
  return new Intl.NumberFormat('id-ID').format(n || 0);
}

function getStatus(spent, budget) {
  if (budget === 0) return { color: '#94a3b8', bg: 'bg-slate-50', label: 'Belum diatur', icon: Clock };
  const pct = (spent / budget) * 100;
  if (pct > 100) return { color: '#ef4444', bg: 'bg-red-50', label: 'Over Budget', icon: AlertCircle };
  if (pct >= 80) return { color: '#f59e0b', bg: 'bg-amber-50', label: 'Hampir Habis', icon: AlertCircle };
  if (pct >= 50) return { color: '#f97316', bg: 'bg-orange-50', label: `${pct.toFixed(0)}% Terpakai`, icon: Clock };
  return { color: '#10b981', bg: 'bg-emerald-50', label: 'Aman', icon: CheckCircle2 };
}

// Circular Progress Ring
function ProgressRing({ percent, size = 120, strokeWidth = 10, color = '#6366f1' }) {
  const radius = (size - strokeWidth) / 2;
  const circumference = radius * 2 * Math.PI;
  const offset = circumference - (percent / 100) * circumference;
  const gradientId = `ring-grad-${Math.round(percent)}`;

  return (
    <svg width={size} height={size} className="transform -rotate-90">
      <defs>
        <linearGradient id={gradientId} x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor={color} />
          <stop offset="100%" stopColor={color} stopOpacity={0.6} />
        </linearGradient>
      </defs>
      <circle cx={size / 2} cy={size / 2} r={radius} stroke="#f1f5f9" strokeWidth={strokeWidth} fill="none" />
      <circle
        cx={size / 2} cy={size / 2} r={radius}
        stroke={`url(#${gradientId})`}
        strokeWidth={strokeWidth}
        fill="none"
        strokeLinecap="round"
        strokeDasharray={circumference}
        strokeDashoffset={offset}
        className="transition-all duration-1000 ease-out"
      />
    </svg>
  );
}

export default function Budgets() {
  const now = new Date();
  const [month, setMonth] = useState(now.getMonth() + 1);
  const [year, setYear] = useState(now.getFullYear());
  const [summary, setSummary] = useState(null);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [editCategory, setEditCategory] = useState(null);
  const [editBudget, setEditBudget] = useState(null);
  const [amount, setAmount] = useState('');
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState('');

  useEffect(() => { getCategories().then(setCategories); }, []);

  useEffect(() => {
    setLoading(true);
    getBudgetSummary(month, year).then((data) => {
      setSummary(data);
      setLoading(false);
    }).catch(() => setLoading(false));
  }, [month, year]);

  const openSetBudget = (category) => {
    const existing = summary?.budgets?.find(b => b.category_id === category.id);
    setEditCategory(category);
    setEditBudget(existing || null);
    setAmount(existing ? existing.amount.toString() : '');
    setFormError('');
    setModalOpen(true);
  };

  const handleSave = async () => {
    if (!amount || Number(amount) < 0) { setFormError('Masukkan jumlah yang valid'); return; }
    setSaving(true);
    setFormError('');
    try {
      await upsertBudget({ category_id: editCategory.id, month, year, amount: Number(amount) });
      setModalOpen(false);
      getBudgetSummary(month, year).then(setSummary);
    } catch (err) {
      setFormError(err.response?.data?.error || 'Gagal menyimpan');
    } finally { setSaving(false); }
  };

  const handleDelete = async (id) => {
    if (!confirm('Hapus anggaran ini?')) return;
    try {
      await deleteBudget(id);
      getBudgetSummary(month, year).then(setSummary);
    } catch (e) { alert('Gagal menghapus'); }
  };

  const prevMonth = () => { if (month === 1) { setMonth(12); setYear(year - 1); } else { setMonth(month - 1); } };
  const nextMonth = () => { if (month === 12) { setMonth(1); setYear(year + 1); } else { setMonth(month + 1); } };
  const monthName = new Date(year, month - 1).toLocaleString('id', { month: 'long' });

  const expenseCategories = categories.filter(c => c.type === 'expense');
  const budgetMap = {};
  summary?.budgets?.forEach(b => { budgetMap[b.category_id] = b; });

  const totalBudget = summary?.total_budget || 0;
  const totalSpent = summary?.total_spent || 0;
  const totalRemaining = totalBudget - totalSpent;
  const overallPct = totalBudget > 0 ? Math.min(Math.round((totalSpent / totalBudget) * 100), 100) : 0;
  const budgetedCount = summary?.budgets?.filter(b => b.amount > 0).length || 0;

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin w-8 h-8 border-4 border-indigo-500 border-t-transparent rounded-full" />
      </div>
    );
  }

  return (
    <div className="space-y-6 sm:space-y-8">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl sm:text-2xl font-bold text-gray-900 flex items-center gap-2.5">
            <div className="p-2 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl shadow-lg shadow-indigo-500/20">
              <Wallet className="text-white" size={20} />
            </div>
            <span>Anggaran</span>
          </h1>
          <p className="text-sm text-gray-400 mt-1 ml-12">Kelola pengeluaran bulananmu</p>
        </div>
      </div>

      {/* Month navigator */}
      <div className="inline-flex items-center bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
        <button onClick={prevMonth} className="p-3 hover:bg-gray-50 transition-colors border-r border-gray-100">
          <ChevronLeft size={18} className="text-gray-500" />
        </button>
        <div className="flex items-center gap-2 px-5">
          <Calendar size={16} className="text-indigo-500" />
          <span className="font-semibold text-gray-800">{monthName} {year}</span>
        </div>
        <button onClick={nextMonth} className="p-3 hover:bg-gray-50 transition-colors border-l border-gray-100">
          <ChevronRight size={18} className="text-gray-500" />
        </button>
      </div>

      {/* Overview Card - Premium Design */}
      <div className="relative overflow-hidden bg-gradient-to-br from-slate-900 via-indigo-950 to-slate-900 rounded-3xl p-6 sm:p-8 text-white">
        {/* Background decoration */}
        <div className="absolute top-0 right-0 w-64 h-64 bg-indigo-500/10 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2" />
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-purple-500/10 rounded-full blur-3xl translate-y-1/2 -translate-x-1/2" />

        <div className="relative flex flex-col sm:flex-row items-center gap-6 sm:gap-8">
          {/* Circular Progress Ring */}
          <div className="relative shrink-0">
            <ProgressRing
              percent={overallPct}
              size={140}
              strokeWidth={12}
              color={overallPct > 80 ? '#ef4444' : overallPct > 50 ? '#f59e0b' : '#6366f1'}
            />
            <div className="absolute inset-0 flex flex-col items-center justify-center">
              <span className="text-3xl font-bold">{overallPct}%</span>
              <span className="text-xs text-indigo-200 mt-0.5">terpakai</span>
            </div>
          </div>

          {/* Stats */}
          <div className="flex-1 w-full space-y-4">
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
              <div className="bg-white/5 backdrop-blur-sm rounded-2xl p-4 border border-white/10">
                <div className="flex items-center gap-2 mb-1.5">
                  <Target size={14} className="text-indigo-300" />
                  <span className="text-[11px] text-indigo-200 uppercase tracking-wider font-medium">Anggaran</span>
                </div>
                <p className="text-lg sm:text-xl font-bold">Rp {formatCompact(totalBudget)}</p>
              </div>
              <div className="bg-white/5 backdrop-blur-sm rounded-2xl p-4 border border-white/10">
                <div className="flex items-center gap-2 mb-1.5">
                  <TrendingUp size={14} className="text-rose-300" />
                  <span className="text-[11px] text-indigo-200 uppercase tracking-wider font-medium">Terpakai</span>
                </div>
                <p className="text-lg sm:text-xl font-bold text-rose-300">Rp {formatCompact(totalSpent)}</p>
              </div>
              <div className="bg-white/5 backdrop-blur-sm rounded-2xl p-4 border border-white/10 col-span-2 sm:col-span-1">
                <div className="flex items-center gap-2 mb-1.5">
                  <Sparkles size={14} className="text-emerald-300" />
                  <span className="text-[11px] text-indigo-200 uppercase tracking-wider font-medium">Sisa</span>
                </div>
                <p className={`text-lg sm:text-xl font-bold ${totalRemaining >= 0 ? 'text-emerald-300' : 'text-red-400'}`}>
                  {totalRemaining >= 0 ? `Rp ${formatCompact(totalRemaining)}` : `-Rp ${formatCompact(Math.abs(totalRemaining))}`}
                </p>
              </div>
            </div>

            {/* Progress bar */}
            <div className="space-y-2">
              <div className="w-full bg-white/10 rounded-full h-2 overflow-hidden">
                <div
                  className="h-full rounded-full transition-all duration-1000 ease-out"
                  style={{
                    width: `${overallPct}%`,
                    background: overallPct > 80
                      ? 'linear-gradient(90deg, #f87171, #ef4444)'
                      : overallPct > 50
                      ? 'linear-gradient(90deg, #fbbf24, #f59e0b)'
                      : 'linear-gradient(90deg, #818cf8, #6366f1)',
                  }}
                />
              </div>
              <div className="flex justify-between text-[11px] text-indigo-300">
                <span>{budgetedCount} dari {expenseCategories.length} kategori diatur</span>
                <span>{formatRupiah(totalSpent)} / {formatRupiah(totalBudget)}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Per-category budget cards */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-sm font-semibold text-gray-700 uppercase tracking-wider">Anggaran per Kategori</h2>
          <span className="text-xs text-gray-400">{budgetedCount} aktif</span>
        </div>

        {expenseCategories.length === 0 ? (
          <div className="bg-white rounded-2xl border border-gray-100 p-12 text-center shadow-sm">
            <div className="w-16 h-16 mx-auto mb-4 bg-indigo-50 rounded-2xl flex items-center justify-center">
              <Wallet size={28} className="text-indigo-400" />
            </div>
            <p className="text-gray-600 font-medium">Belum ada kategori pengeluaran</p>
            <p className="text-sm text-gray-400 mt-1">Buat kategori terlebih dahulu di halaman Kategori</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
            {expenseCategories.map((cat) => {
              const budget = budgetMap[cat.id];
              const spent = budget?.spent || 0;
              const budgetAmount = budget?.amount || 0;
              const remaining = budgetAmount - spent;
              const pct = budgetAmount > 0 ? Math.min(Math.round((spent / budgetAmount) * 100), 100) : 0;
              const status = getStatus(spent, budgetAmount);
              const StatusIcon = status.icon;

              return (
                <div
                  key={cat.id}
                  className={`group relative bg-white rounded-2xl border shadow-sm overflow-hidden transition-all hover:shadow-md ${
                    budget ? 'border-gray-100' : 'border-dashed border-gray-200'
                  }`}
                >
                  {/* Color accent stripe */}
                  {budget && (
                    <div className="absolute top-0 left-0 right-0 h-1" style={{ backgroundColor: cat.color, opacity: 0.8 }} />
                  )}

                  <div className="p-4 sm:p-5">
                    {/* Header row */}
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-3 min-w-0 flex-1">
                        <div
                          className="w-10 h-10 rounded-xl flex items-center justify-center text-white text-sm font-bold shrink-0 shadow-sm transition-transform group-hover:scale-105"
                          style={{ backgroundColor: cat.color }}
                        >
                          {cat.name.charAt(0)}
                        </div>
                        <div className="min-w-0">
                          <p className="text-sm font-semibold text-gray-800 truncate">{cat.name}</p>
                          {budget ? (
                            <div className="flex items-center gap-1.5 mt-0.5">
                              <StatusIcon size={12} style={{ color: status.color }} />
                              <span className="text-xs font-medium" style={{ color: status.color }}>{status.label}</span>
                            </div>
                          ) : (
                            <span className="text-xs text-gray-400 mt-0.5 block">Belum diatur</span>
                          )}
                        </div>
                      </div>

                      <div className="flex items-center gap-1 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button
                          onClick={() => openSetBudget(cat)}
                          className="p-2 bg-indigo-50 text-indigo-600 hover:bg-indigo-100 rounded-lg transition-colors"
                          title="Atur anggaran"
                        >
                          {budget ? <Edit3 size={14} /> : <Plus size={14} />}
                        </button>
                        {budget && (
                          <button
                            onClick={() => handleDelete(budget.id)}
                            className="p-2 bg-red-50 text-red-400 hover:bg-red-100 rounded-lg transition-colors"
                            title="Hapus"
                          >
                            <Trash2 size={14} />
                          </button>
                        )}
                      </div>
                      {/* Always visible + button on mobile */}
                      <button
                        onClick={() => openSetBudget(cat)}
                        className="sm:hidden p-2 bg-indigo-50 text-indigo-600 rounded-lg shrink-0"
                      >
                        {budget ? <Edit3 size={14} /> : <Plus size={14} />}
                      </button>
                    </div>

                    {budget ? (
                      <>
                        {/* Amount */}
                        <div className="mb-3">
                          <div className="flex items-baseline gap-1.5">
                            <span className="text-xl sm:text-2xl font-bold" style={{ color: status.color }}>
                              {formatCompact(spent)}
                            </span>
                            <span className="text-xs text-gray-400 font-medium">/ {formatCompact(budgetAmount)}</span>
                          </div>
                          <div className="flex items-center gap-2 mt-1">
                            <span className="text-xs text-gray-400">Rp {formatRupiah(spent)}</span>
                            <span className="text-xs text-gray-300">•</span>
                            <span className={`text-xs font-medium ${remaining >= 0 ? 'text-emerald-600' : 'text-red-500'}`}>
                              Sisa {remaining >= 0 ? `Rp ${formatCompact(remaining)}` : `-Rp ${formatCompact(Math.abs(remaining))}`}
                            </span>
                          </div>
                        </div>

                        {/* Progress bar */}
                        <div className="w-full bg-gray-100 rounded-full h-2 overflow-hidden">
                          <div
                            className="h-full rounded-full transition-all duration-700 ease-out"
                            style={{
                              width: `${pct}%`,
                              background: `linear-gradient(90deg, ${status.color}, ${status.color}dd)`,
                            }}
                          />
                        </div>
                      </>
                    ) : (
                      /* Empty state */
                      <button
                        onClick={() => openSetBudget(cat)}
                        className="w-full flex items-center justify-center gap-2 py-3 rounded-xl border-2 border-dashed border-gray-200 text-gray-400 hover:border-indigo-300 hover:text-indigo-500 hover:bg-indigo-50/50 transition-all text-sm"
                      >
                        <Plus size={16} />
                        <span>Atur Anggaran</span>
                      </button>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Set Budget Modal */}
      <Modal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        title={`Anggaran: ${editCategory?.name || ''}`}
      >
        <div className="space-y-5">
          {editCategory && (
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 rounded-xl flex items-center justify-center text-white text-sm font-bold shadow-md" style={{ backgroundColor: editCategory.color }}>
                {editCategory.name.charAt(0)}
              </div>
              <div>
                <p className="font-semibold text-gray-800">{editCategory.name}</p>
                <p className="text-xs text-gray-400">{monthName} {year}</p>
              </div>
            </div>
          )}

          {editBudget && editBudget.spent > 0 && (
            <div className="flex justify-between items-center text-sm bg-indigo-50 rounded-xl px-4 py-3 border border-indigo-100">
              <span className="text-indigo-600 font-medium">Terpakai bulan ini</span>
              <span className="font-bold text-indigo-700">Rp {formatRupiah(editBudget.spent)}</span>
            </div>
          )}

          {formError && (
            <div className="bg-red-50 text-red-600 text-sm rounded-xl px-4 py-3 border border-red-100 flex items-center gap-2">
              <AlertCircle size={16} />
              {formError}
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Jumlah Anggaran</label>
            <div className="relative">
              <span className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400 font-medium">Rp</span>
              <input
                type="number"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="5.000.000"
                min="0"
                className="w-full border border-gray-200 rounded-xl pl-12 pr-4 py-3.5 text-xl font-bold focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none bg-gray-50 focus:bg-white transition-colors"
                autoFocus
                inputMode="numeric"
              />
            </div>
            {/* Quick amounts */}
            <div className="grid grid-cols-4 gap-2 mt-3">
              {[500000, 1000000, 2000000, 5000000].map((amt) => (
                <button
                  key={amt}
                  type="button"
                  onClick={() => setAmount(amt.toString())}
                  className="py-2.5 bg-gray-100 hover:bg-indigo-50 hover:text-indigo-600 rounded-xl text-xs font-semibold text-gray-600 transition-colors"
                >
                  {amt >= 1000000 ? `${amt / 1000000}jt` : `${amt / 1000}rb`}
                </button>
              ))}
            </div>
          </div>

          <button
            onClick={handleSave}
            disabled={saving}
            className="w-full bg-gradient-to-r from-indigo-600 to-purple-600 text-white py-3.5 rounded-xl font-semibold hover:from-indigo-700 hover:to-purple-700 transition-all disabled:opacity-50 text-sm shadow-lg shadow-indigo-500/25"
          >
            {saving ? 'Menyimpan...' : editBudget ? 'Perbarui Anggaran' : 'Atur Anggaran'}
          </button>
        </div>
      </Modal>
    </div>
  );
}
