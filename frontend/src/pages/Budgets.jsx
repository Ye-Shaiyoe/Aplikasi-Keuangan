import { useState, useEffect } from 'react';
import { Wallet, TrendingUp, TrendingDown, Calendar, Plus, ChevronLeft, ChevronRight, Edit3, PiggyBank } from 'lucide-react';
import Modal from '../components/Modal';
import { getBudgetSummary, upsertBudget, deleteBudget, getCategories } from '../api/client';

function formatRupiah(n) {
  return new Intl.NumberFormat('id-ID').format(n || 0);
}

function formatCompact(n) {
  if (n >= 1000000) return `${(n / 1000000).toFixed(n % 1000000 === 0 ? 0 : 1)}jt`;
  if (n >= 1000) return `${(n / 1000).toFixed(0)}rb`;
  return new Intl.NumberFormat('id-ID').format(n || 0);
}

function getStatus(spent, budget) {
  if (budget === 0) return { color: '#d1d5db', bg: 'bg-gray-100', label: 'Belum diatur', textColor: 'text-gray-400' };
  const pct = (spent / budget) * 100;
  if (pct > 100) return { color: '#ef4444', bg: 'bg-red-100', label: 'Over Budget!', textColor: 'text-red-600' };
  if (pct >= 80) return { color: '#f59e0b', bg: 'bg-amber-100', label: 'Hampir Habis', textColor: 'text-amber-600' };
  if (pct >= 50) return { color: '#f97316', bg: 'bg-orange-100', label: 'Terpakai ' + pct.toFixed(0) + '%', textColor: 'text-orange-600' };
  return { color: '#10b981', bg: 'bg-green-100', label: 'Aman', textColor: 'text-green-600' };
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

  useEffect(() => {
    getCategories().then(setCategories);
  }, []);

  useEffect(() => {
    setLoading(true);
    getBudgetSummary(month, year).then((data) => {
      setSummary(data);
      setLoading(false);
    }).catch(() => setLoading(false));
  }, [month, year]);

  const openSetBudget = (category) => {
    // Check if there's already a budget for this category
    const existing = summary?.budgets?.find(b => b.category_id === category.id);
    setEditCategory(category);
    setEditBudget(existing || null);
    setAmount(existing ? existing.amount.toString() : '');
    setFormError('');
    setModalOpen(true);
  };

  const handleSave = async () => {
    if (!amount || Number(amount) < 0) {
      setFormError('Masukkan jumlah yang valid');
      return;
    }
    setSaving(true);
    setFormError('');
    try {
      await upsertBudget({
        category_id: editCategory.id,
        month,
        year,
        amount: Number(amount),
      });
      setModalOpen(false);
      getBudgetSummary(month, year).then(setSummary);
    } catch (err) {
      setFormError(err.response?.data?.error || 'Gagal menyimpan');
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id) => {
    if (!confirm('Hapus anggaran ini?')) return;
    try {
      await deleteBudget(id);
      getBudgetSummary(month, year).then(setSummary);
    } catch (e) {
      alert('Gagal menghapus');
    }
  };

  const prevMonth = () => {
    if (month === 1) { setMonth(12); setYear(year - 1); }
    else { setMonth(month - 1); }
  };

  const nextMonth = () => {
    if (month === 12) { setMonth(1); setYear(year + 1); }
    else { setMonth(month + 1); }
  };

  const monthName = new Date(year, month - 1).toLocaleString('id', { month: 'long' });

  const expenseCategories = categories.filter(c => c.type === 'expense');
  const budgetMap = {};
  summary?.budgets?.forEach(b => { budgetMap[b.category_id] = b; });

  const totalBudget = summary?.total_budget || 0;
  const totalSpent = summary?.total_spent || 0;
  const totalRemaining = totalBudget - totalSpent;

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full" />
      </div>
    );
  }

  return (
    <div className="space-y-4 sm:space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl sm:text-2xl font-bold text-gray-800 flex items-center gap-2">
          <Wallet className="text-blue-500" size={24} />
          <span>Anggaran</span>
        </h1>
      </div>

      {/* Month navigator */}
      <div className="flex items-center justify-between bg-white rounded-2xl border border-gray-100 p-3 sm:p-4 shadow-sm">
        <button onClick={prevMonth} className="p-2 hover:bg-gray-50 rounded-xl transition-colors">
          <ChevronLeft size={20} className="text-gray-400" />
        </button>
        <div className="flex items-center gap-2">
          <Calendar size={18} className="text-blue-500" />
          <span className="font-semibold text-gray-800 text-base sm:text-lg">{monthName} {year}</span>
        </div>
        <button onClick={nextMonth} className="p-2 hover:bg-gray-50 rounded-xl transition-colors">
          <ChevronRight size={20} className="text-gray-400" />
        </button>
      </div>

      {/* Overview Card */}
      <div className="bg-gradient-to-br from-blue-500 to-indigo-600 rounded-2xl p-5 sm:p-6 text-white shadow-sm">
        <div className="flex items-start justify-between mb-4">
          <div>
            <p className="text-blue-100 text-xs sm:text-sm">Total Anggaran</p>
            <p className="text-xl sm:text-2xl font-bold mt-0.5">Rp {formatRupiah(totalBudget)}</p>
          </div>
          <div className="text-right">
            <p className="text-blue-100 text-xs sm:text-sm">Terpakai</p>
            <p className="text-xl sm:text-2xl font-bold mt-0.5">Rp {formatRupiah(totalSpent)}</p>
          </div>
        </div>
        {/* Total progress bar */}
        <div className="w-full bg-white/20 rounded-full h-3 sm:h-4 overflow-hidden">
          <div
            className="h-full rounded-full transition-all duration-500"
            style={{
              width: `${totalBudget > 0 ? Math.min((totalSpent / totalBudget) * 100, 100) : 0}%`,
              backgroundColor: totalBudget > 0 && totalSpent > totalBudget ? '#ef4444' : '#ffffff',
            }}
          />
        </div>
        <div className="flex justify-between mt-2 text-xs sm:text-sm text-blue-100">
          <span>Sisa: Rp {formatCompact(Math.max(totalRemaining, 0))}</span>
          <span>{totalBudget > 0 ? Math.min(Math.round((totalSpent / totalBudget) * 100), 100) : 0}%</span>
        </div>
      </div>

      {/* Per-category budget cards */}
      <div className="space-y-3">
        <h2 className="text-sm font-semibold text-gray-600">Anggaran per Kategori</h2>

        {expenseCategories.map((cat) => {
          const budget = budgetMap[cat.id];
          const spent = budget?.spent || 0;
          const budgetAmount = budget?.amount || 0;
          const remaining = budgetAmount - spent;
          const pct = budgetAmount > 0 ? Math.min(Math.round((spent / budgetAmount) * 100), 100) : 0;
          const status = getStatus(spent, budgetAmount);

          return (
            <div key={cat.id} className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
              <div className="p-4 sm:p-5">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3 min-w-0 flex-1">
                    <div
                      className="w-10 h-10 sm:w-12 sm:h-12 rounded-xl flex items-center justify-center text-white text-sm font-bold shrink-0 shadow-sm"
                      style={{ backgroundColor: cat.color }}
                    >
                      {cat.name.charAt(0)}
                    </div>
                    <div className="min-w-0">
                      <p className="text-sm sm:text-base font-semibold text-gray-800 truncate">{cat.name}</p>
                      <div className="flex items-center gap-2 mt-0.5">
                        <span className={`text-xs font-medium ${status.textColor}`}>
                          {budget ? status.label : 'Belum diatur'}
                        </span>
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-1.5 sm:gap-2 shrink-0">
                    <button
                      onClick={() => openSetBudget(cat)}
                      className="p-2 sm:p-2.5 bg-blue-50 text-blue-600 hover:bg-blue-100 rounded-xl transition-colors"
                      title="Atur anggaran"
                    >
                      {budget ? <Edit3 size={16} /> : <Plus size={16} />}
                    </button>
                    {budget && (
                      <button
                        onClick={() => handleDelete(budget.id)}
                        className="p-2 sm:p-2.5 bg-red-50 text-red-400 hover:bg-red-100 rounded-xl transition-colors"
                        title="Hapus"
                      >
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M3 6h18" /><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" /><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
                        </svg>
                      </button>
                    )}
                  </div>
                </div>

                {budget && (
                  <>
                    {/* Amount info */}
                    <div className="flex items-end justify-between mb-2">
                      <div className="space-y-0.5">
                        <p className="text-xs text-gray-400">Terpakai</p>
                        <p className="text-base sm:text-lg font-bold" style={{ color: status.color }}>
                          Rp {formatRupiah(spent)}
                          <span className="text-gray-400 text-xs sm:text-sm font-normal"> / Rp {formatRupiah(budgetAmount)}</span>
                        </p>
                      </div>
                      <div className="text-right">
                        <p className="text-xs text-gray-400">Sisa</p>
                        <p className={`text-sm font-bold ${remaining >= 0 ? 'text-gray-700' : 'text-red-600'}`}>
                          {remaining >= 0 ? `Rp ${formatCompact(remaining)}` : '−'}
                        </p>
                      </div>
                    </div>

                    {/* Progress bar */}
                    <div className="w-full bg-gray-100 rounded-full h-2.5 sm:h-3 overflow-hidden">
                      <div
                        className="h-full rounded-full transition-all duration-500"
                        style={{ width: `${pct}%`, backgroundColor: status.color }}
                      />
                    </div>

                    {/* Mini stats */}
                    <div className="grid grid-cols-3 gap-2 mt-3">
                      <div className={`rounded-lg px-2.5 py-1.5 text-center ${status.bg}`}>
                        <p className={`text-xs font-semibold ${status.textColor}`}>{pct}%</p>
                        <p className="text-[10px] text-gray-400">Terpakai</p>
                      </div>
                      <div className="bg-gray-50 rounded-lg px-2.5 py-1.5 text-center">
                        <p className="text-xs font-semibold text-gray-700">{spent > 0 && budgetAmount > 0 ? Math.round((budgetAmount - spent) / (spent / (spent / budgetAmount)) || 0) : '-'}</p>
                        <p className="text-[10px] text-gray-400">Hari</p>
                      </div>
                      <div className="bg-gray-50 rounded-lg px-2.5 py-1.5 text-center">
                        <p className={`text-xs font-semibold ${remaining >= 0 ? 'text-blue-600' : 'text-red-600'}`}>
                          {remaining >= 0 ? `Rp ${formatCompact(remaining)}` : '−'}
                        </p>
                        <p className="text-[10px] text-gray-400">Sisa</p>
                      </div>
                    </div>
                  </>
                )}

                {!budget && (
                  <div className="bg-gray-50 rounded-xl p-4 text-center">
                    <p className="text-sm text-gray-400 mb-2">Belum ada anggaran untuk kategori ini</p>
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>

      {/* Set Budget Modal */}
      <Modal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        title={`Anggaran: ${editCategory?.name || ''}`}
      >
        <div className="space-y-4">
          {editCategory && (
            <div className="flex items-center gap-3 bg-gray-50 rounded-xl p-4">
              <div className="w-10 h-10 rounded-xl flex items-center justify-center text-white text-sm font-bold" style={{ backgroundColor: editCategory.color }}>
                {editCategory.name.charAt(0)}
              </div>
              <div>
                <p className="font-medium text-gray-800">{editCategory.name}</p>
                <p className="text-xs text-gray-400">{monthName} {year}</p>
              </div>
            </div>
          )}

          {editBudget && (
            <div className="flex justify-between text-sm bg-blue-50 rounded-xl px-4 py-3">
              <span className="text-blue-600">Terpakai bulan ini</span>
              <span className="font-bold text-blue-700">Rp {formatRupiah(editBudget.spent)}</span>
            </div>
          )}

          {formError && (
            <div className="bg-red-50 text-red-600 text-sm rounded-xl px-4 py-3 border border-red-100">
              {formError}
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Jumlah Anggaran (Rp)</label>
            <input
              type="number"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              placeholder="5000000"
              min="0"
              className="w-full border border-gray-200 rounded-xl px-4 py-3 text-lg font-semibold focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
              autoFocus
              inputMode="numeric"
            />
            {/* Quick amounts */}
            <div className="grid grid-cols-4 gap-2 mt-2.5">
              {[500000, 1000000, 2000000, 5000000].map((amt) => (
                <button
                  key={amt}
                  type="button"
                  onClick={() => setAmount(amt.toString())}
                  className="py-2 bg-gray-100 hover:bg-gray-200 rounded-xl text-xs font-medium text-gray-600 transition-colors"
                >
                  {amt >= 1000000 ? `${amt / 1000000}jt` : `${amt / 1000}rb`}
                </button>
              ))}
            </div>
          </div>

          <button
            onClick={handleSave}
            disabled={saving}
            className="w-full bg-blue-600 text-white py-3 rounded-xl font-semibold hover:bg-blue-700 transition-colors disabled:opacity-50 text-sm"
          >
            {saving ? 'Menyimpan...' : editBudget ? 'Perbarui Anggaran' : 'Atur Anggaran'}
          </button>
        </div>
      </Modal>
    </div>
  );
}
