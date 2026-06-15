import { useState, useEffect, useRef } from 'react';
import { Target, Plus, Edit2, Trash2, PiggyBank, TrendingUp, Calendar, ArrowUpCircle, ArrowDownCircle, ImagePlus, X } from 'lucide-react';
import Modal from '../components/Modal';
import { getSavingsGoals, createSavingsGoal, updateSavingsGoal, deleteSavingsGoal, depositToSavingsGoal, withdrawFromSavingsGoal } from '../api/client';

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#06b6d4', '#f97316'];

function formatRupiah(n) { return new Intl.NumberFormat('id-ID').format(n); }
function formatCompact(n) {
  if (n >= 1000000) return `${(n / 1000000).toFixed(n % 1000000 === 0 ? 0 : 1)}jt`;
  if (n >= 1000) return `${(n / 1000).toFixed(0)}rb`;
  return new Intl.NumberFormat('id-ID').format(n);
}
function getProgressPercent(current, target) { if (!target || target === 0) return 0; return Math.min(Math.round((current / target) * 100), 100); }
function getDaysRemaining(deadline) { if (!deadline) return null; return Math.ceil((new Date(deadline) - new Date()) / (1000 * 60 * 60 * 24)); }

// Compress image to base64 JPEG (max 400px wide, quality 0.7)
function compressImage(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const img = new Image();
      img.onload = () => {
        const MAX_W = 400;
        const scale = Math.min(MAX_W / img.width, 1);
        const w = Math.round(img.width * scale);
        const h = Math.round(img.height * scale);
        const canvas = document.createElement('canvas');
        canvas.width = w; canvas.height = h;
        const ctx = canvas.getContext('2d');
        ctx.drawImage(img, 0, 0, w, h);
        resolve(canvas.toDataURL('image/jpeg', 0.7));
      };
      img.onerror = reject;
      img.src = e.target.result;
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

export default function SavingsGoals() {
  const [goals, setGoals] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [fundModalOpen, setFundModalOpen] = useState(false);
  const [editing, setEditing] = useState(null);
  const [fundGoal, setFundGoal] = useState(null);
  const [fundMode, setFundMode] = useState('deposit');
  const [formError, setFormError] = useState('');
  const [saving, setSaving] = useState(false);

  const [name, setName] = useState('');
  const [targetAmount, setTargetAmount] = useState('');
  const [deadline, setDeadline] = useState('');
  const [color, setColor] = useState('#3b82f6');
  const [imageURL, setImageURL] = useState('');
  const [fundAmount, setFundAmount] = useState('');
  const [fundError, setFundError] = useState('');
  const fileInputRef = useRef(null);

  useEffect(() => { loadGoals(); }, []);

  async function loadGoals() {
    try { const data = await getSavingsGoals(); setGoals(data); }
    catch (e) { console.error(e); }
    finally { setLoading(false); }
  }

  function openCreate() {
    setEditing(null);
    setName(''); setTargetAmount(''); setDeadline(''); setColor('#3b82f6'); setImageURL(''); setFormError('');
    setModalOpen(true);
  }

  function openEdit(goal) {
    setEditing(goal);
    setName(goal.name); setTargetAmount(goal.target_amount.toString());
    setDeadline(goal.deadline || ''); setColor(goal.color); setImageURL(goal.image_url || ''); setFormError('');
    setModalOpen(true);
  }

  function openFund(goal, mode) {
    setFundGoal(goal); setFundMode(mode); setFundAmount(''); setFundError('');
    setFundModalOpen(true);
  }

  async function handleImageSelect(e) {
    const file = e.target.files?.[0];
    if (!file) return;
    try {
      const compressed = await compressImage(file);
      setImageURL(compressed);
    } catch (err) {
      console.error('Image compress failed', err);
    }
    e.target.value = ''; // reset so same file can be re-selected
  }

  async function handleSave(e) {
    e.preventDefault();
    if (!name.trim()) return setFormError('Nama target wajib diisi');
    if (!targetAmount || Number(targetAmount) < 1) return setFormError('Jumlah target harus lebih dari 0');
    setSaving(true); setFormError('');
    try {
      const payload = { name: name.trim(), target_amount: Number(targetAmount), deadline: deadline || null, color, image_url: imageURL };
      if (editing) await updateSavingsGoal(editing.id, payload);
      else await createSavingsGoal(payload);
      setModalOpen(false); loadGoals();
    } catch (err) { setFormError(err.response?.data?.error || 'Gagal menyimpan'); }
    finally { setSaving(false); }
  }

  async function handleFund(e) {
    e.preventDefault();
    if (!fundAmount || Number(fundAmount) < 1) return setFundError('Jumlah harus lebih dari 0');
    setSaving(true); setFundError('');
    try {
      const amt = Number(fundAmount);
      if (fundMode === 'deposit') await depositToSavingsGoal(fundGoal.id, amt);
      else await withdrawFromSavingsGoal(fundGoal.id, amt);
      setFundModalOpen(false); loadGoals();
    } catch (err) { setFundError(err.response?.data?.error || 'Gagal memproses'); }
    finally { setSaving(false); }
  }

  async function handleDelete(id) {
    if (!confirm('Hapus target tabungan ini?')) return;
    try { await deleteSavingsGoal(id); loadGoals(); }
    catch (e) { alert('Gagal menghapus'); }
  }

  const totalTarget = goals.reduce((s, g) => s + g.target_amount, 0);
  const totalSaved = goals.reduce((s, g) => s + g.current_amount, 0);
  const completedCount = goals.filter((g) => g.current_amount >= g.target_amount).length;

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
        <div>
          <h1 className="text-xl sm:text-2xl font-bold text-gray-800 flex items-center gap-2">
            <Target className="text-blue-500" size={24} />
            <span className="sm:hidden">Tabungan</span>
            <span className="hidden sm:inline">Target Tabungan</span>
          </h1>
          <p className="text-gray-400 text-xs sm:text-sm mt-0.5 sm:mt-1 hidden sm:block">Pantau progres tabunganmu menuju impian</p>
        </div>
        <button onClick={openCreate} className="flex items-center gap-1.5 sm:gap-2 bg-blue-600 text-white px-3.5 sm:px-4 py-2 sm:py-2.5 rounded-xl hover:bg-blue-700 transition-colors font-medium text-sm shadow-sm">
          <Plus size={18} />
          <span className="sm:hidden">Baru</span>
          <span className="hidden sm:inline">Target Baru</span>
        </button>
      </div>

      {/* Summary Cards */}
      {goals.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 sm:gap-4">
          <div className="bg-gradient-to-br from-blue-500 to-blue-600 rounded-2xl p-4 sm:p-5 text-white">
            <div className="flex items-center gap-2 opacity-80 text-xs sm:text-sm mb-1">
              <PiggyBank size={16} /> Total Terkumpul
            </div>
            <p className="text-lg sm:text-xl font-bold">Rp <span className="sm:hidden">{formatCompact(totalSaved)}</span><span className="hidden sm:inline">{formatRupiah(totalSaved)}</span></p>
            <div className="mt-2 w-full bg-white/20 rounded-full h-1.5 sm:h-2">
              <div className="bg-white rounded-full h-1.5 sm:h-2 transition-all duration-500" style={{ width: `${getProgressPercent(totalSaved, totalTarget)}%` }} />
            </div>
            <p className="text-[10px] sm:text-xs mt-1 opacity-70">dari Rp <span className="sm:hidden">{formatCompact(totalTarget)}</span><span className="hidden sm:inline">{formatRupiah(totalTarget)}</span></p>
          </div>
          <div className="hidden sm:block bg-white rounded-2xl p-5 border border-gray-100 shadow-sm">
            <div className="flex items-center gap-2 text-gray-500 text-sm mb-1"><Target size={16} />Target Aktif</div>
            <p className="text-2xl font-bold text-gray-800">{goals.length}</p>
          </div>
          <div className="hidden sm:block bg-white rounded-2xl p-5 border border-gray-100 shadow-sm">
            <div className="flex items-center gap-2 text-green-500 text-sm mb-1"><TrendingUp size={16} />Tercapai</div>
            <p className="text-2xl font-bold text-gray-800">{completedCount}</p>
          </div>
          <div className="sm:hidden flex gap-3">
            <div className="flex-1 bg-white rounded-xl p-3 border border-gray-100 text-center">
              <p className="text-lg font-bold text-gray-800">{goals.length}</p>
              <p className="text-[10px] text-gray-400">Target Aktif</p>
            </div>
            <div className="flex-1 bg-white rounded-xl p-3 border border-gray-100 text-center">
              <p className="text-lg font-bold text-green-600">{completedCount}</p>
              <p className="text-[10px] text-gray-400">Tercapai</p>
            </div>
          </div>
        </div>
      )}

      {/* Goal Cards */}
      {goals.length === 0 ? (
        <div className="text-center py-16 bg-white rounded-2xl border border-gray-100">
          <PiggyBank size={48} className="mx-auto text-gray-300 mb-4" />
          <p className="text-gray-500 font-medium text-sm">Belum ada target tabungan</p>
          <p className="text-gray-400 text-xs mt-1">Buat target pertamamu dan mulai menabung!</p>
        </div>
      ) : (
        <div className="space-y-3 sm:grid sm:grid-cols-2 sm:gap-4 sm:space-y-0">
          {goals.map((goal) => {
            const pct = getProgressPercent(goal.current_amount, goal.target_amount);
            const isComplete = pct >= 100;
            const daysLeft = getDaysRemaining(goal.deadline);
            const remaining = Math.max(goal.target_amount - goal.current_amount, 0);

            return (
              <div key={goal.id} className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden group">
                {/* Image header (if exists) */}
                {goal.image_url ? (
                  <div className="relative h-36 sm:h-44 overflow-hidden">
                    <img src={goal.image_url} alt={goal.name} className="w-full h-full object-cover" />
                    <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-black/10 to-transparent" />
                    <div className="absolute bottom-3 left-4 right-4">
                      <h3 className="font-bold text-white text-base sm:text-lg truncate drop-shadow">{goal.name}</h3>
                      {goal.deadline && (
                        <div className="flex items-center gap-1 text-[10px] sm:text-xs text-white/80 mt-0.5">
                          <Calendar size={11} />
                          {daysLeft !== null && daysLeft > 0 ? `${daysLeft} hari lagi` : daysLeft === 0 ? 'Hari ini!' : `Terlewat ${Math.abs(daysLeft)} hari`}
                        </div>
                      )}
                    </div>
                    {/* Actions overlay */}
                    <div className="absolute top-2 right-2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button onClick={() => openEdit(goal)} className="p-1.5 bg-white/90 backdrop-blur-sm rounded-lg text-gray-600 hover:text-blue-600 transition-colors"><Edit2 size={14} /></button>
                      <button onClick={() => handleDelete(goal.id)} className="p-1.5 bg-white/90 backdrop-blur-sm rounded-lg text-gray-600 hover:text-red-500 transition-colors"><Trash2 size={14} /></button>
                    </div>
                  </div>
                ) : (
                  <div className="h-1 sm:h-1.5" style={{ backgroundColor: goal.color }} />
                )}

                <div className="p-4 sm:p-5">
                  {/* Header (no image) */}
                  {!goal.image_url && (
                    <div className="flex items-start justify-between mb-3 sm:mb-4">
                      <div className="flex-1 min-w-0">
                        <h3 className="font-semibold text-gray-800 text-base sm:text-lg truncate">{goal.name}</h3>
                        {goal.deadline && (
                          <div className="flex items-center gap-1 text-[10px] sm:text-xs text-gray-400 mt-0.5">
                            <Calendar size={12} />
                            {daysLeft !== null && daysLeft > 0 ? `${daysLeft} hari lagi` : daysLeft === 0 ? 'Hari ini!' : `Terlewat ${Math.abs(daysLeft)} hari`}
                          </div>
                        )}
                      </div>
                      <div className="flex items-center gap-0.5 sm:gap-1 ml-2">
                        <button onClick={() => openFund(goal, 'deposit')} className="p-1.5 sm:p-2 text-green-500 hover:bg-green-50 rounded-lg transition-colors"><ArrowUpCircle size={18} /></button>
                        <button onClick={() => openFund(goal, 'withdraw')} className="p-1.5 sm:p-2 text-amber-500 hover:bg-amber-50 rounded-lg transition-colors"><ArrowDownCircle size={18} /></button>
                        <button onClick={() => openEdit(goal)} className="p-1 sm:p-1.5 text-gray-400 hover:bg-gray-50 rounded-lg transition-colors hidden sm:block"><Edit2 size={15} /></button>
                        <button onClick={() => handleDelete(goal.id)} className="p-1 sm:p-1.5 text-red-400 hover:bg-red-50 rounded-lg transition-colors hidden sm:block"><Trash2 size={15} /></button>
                      </div>
                    </div>
                  )}

                  {/* Fund actions (always visible) */}
                  {goal.image_url && (
                    <div className="flex items-center gap-2 mb-3">
                      <button onClick={() => openFund(goal, 'deposit')} className="flex-1 flex items-center justify-center gap-1.5 py-2 bg-green-50 text-green-600 hover:bg-green-100 rounded-xl transition-colors text-xs sm:text-sm font-medium">
                        <ArrowUpCircle size={16} /> Tambah
                      </button>
                      <button onClick={() => openFund(goal, 'withdraw')} className="flex-1 flex items-center justify-center gap-1.5 py-2 bg-amber-50 text-amber-600 hover:bg-amber-100 rounded-xl transition-colors text-xs sm:text-sm font-medium">
                        <ArrowDownCircle size={16} /> Tarik
                      </button>
                      {/* Mobile edit/delete */}
                      <button onClick={() => openEdit(goal)} className="sm:hidden p-2 text-gray-400 hover:bg-gray-50 rounded-lg"><Edit2 size={16} /></button>
                      <button onClick={() => handleDelete(goal.id)} className="sm:hidden p-2 text-red-400 hover:bg-red-50 rounded-lg"><Trash2 size={16} /></button>
                    </div>
                  )}

                  <div className="mb-2.5 sm:mb-3">
                    <div className="flex items-end justify-between mb-1.5 sm:mb-2">
                      <div>
                        <span className="text-base sm:text-xl font-bold" style={{ color: goal.color }}>
                          Rp <span className="sm:hidden">{formatCompact(goal.current_amount)}</span><span className="hidden sm:inline">{formatRupiah(goal.current_amount)}</span>
                        </span>
                        <span className="text-gray-400 text-xs sm:text-sm ml-1">
                          / Rp <span className="sm:hidden">{formatCompact(goal.target_amount)}</span><span className="hidden sm:inline">{formatRupiah(goal.target_amount)}</span>
                        </span>
                      </div>
                      <span className={`text-xs sm:text-sm font-semibold px-1.5 sm:px-2 py-0.5 rounded-full ${isComplete ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'}`}>
                        {pct}%
                      </span>
                    </div>
                    <div className="w-full bg-gray-100 rounded-full h-2 sm:h-3 overflow-hidden">
                      <div className="h-2 sm:h-3 rounded-full transition-all duration-700 ease-out" style={{ width: `${pct}%`, backgroundColor: goal.color }} />
                    </div>
                  </div>

                  {!isComplete ? (
                    <p className="text-[10px] sm:text-xs text-gray-400">
                      Sisa <span className="font-medium text-gray-600">Rp <span className="sm:hidden">{formatCompact(remaining)}</span><span className="hidden sm:inline">{formatRupiah(remaining)}</span></span> lagi
                    </p>
                  ) : (
                    <p className="text-[10px] sm:text-xs text-green-600 font-medium">Target tercapai!</p>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Create/Edit Modal */}
      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title={editing ? 'Edit Target' : 'Target Tabungan Baru'}>
        <form onSubmit={handleSave} className="space-y-4">
          {formError && <div className="bg-red-50 text-red-600 text-sm rounded-xl px-4 py-3 border border-red-100">{formError}</div>}

          {/* Image upload */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Foto Target (opsional)</label>
            {imageURL ? (
              <div className="relative rounded-xl overflow-hidden border border-gray-200">
                <img src={imageURL} alt="Preview" className="w-full h-40 object-cover" />
                <button type="button" onClick={() => setImageURL('')} className="absolute top-2 right-2 p-1.5 bg-black/50 text-white rounded-lg hover:bg-black/70 transition-colors">
                  <X size={14} />
                </button>
              </div>
            ) : (
              <button
                type="button"
                onClick={() => fileInputRef.current?.click()}
                className="w-full flex flex-col items-center gap-2 py-8 rounded-xl border-2 border-dashed border-gray-200 text-gray-400 hover:border-blue-300 hover:text-blue-500 hover:bg-blue-50/30 transition-all"
              >
                <ImagePlus size={28} />
                <span className="text-sm">Tambah foto (misal: foto barang yang mau dibeli)</span>
                <span className="text-xs text-gray-300">JPG, PNG - maks ~500KB setelah kompresi</span>
              </button>
            )}
            <input ref={fileInputRef} type="file" accept="image/*" className="hidden" onChange={handleImageSelect} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Nama Target</label>
            <input type="text" value={name} onChange={(e) => setName(e.target.value)} placeholder="Contoh: iPhone 16, Liburan ke Jepang..." className="w-full border border-gray-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm" autoFocus />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Jumlah Target (Rp)</label>
            <input type="number" value={targetAmount} onChange={(e) => setTargetAmount(e.target.value)} placeholder="10000000" min="1" className="w-full border border-gray-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm font-semibold" inputMode="numeric" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Deadline (opsional)</label>
            <input type="date" value={deadline} onChange={(e) => setDeadline(e.target.value)} className="w-full border border-gray-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Warna</label>
            <div className="flex gap-2.5 flex-wrap">
              {COLORS.map((c) => (
                <button key={c} type="button" onClick={() => setColor(c)} className={`w-9 h-9 sm:w-8 sm:h-8 rounded-full transition-all ${color === c ? 'ring-2 ring-offset-2 ring-gray-400 scale-110' : 'hover:scale-105'}`} style={{ backgroundColor: c }} />
              ))}
            </div>
          </div>
          <button type="submit" disabled={saving} className="w-full bg-blue-600 text-white py-3 rounded-xl font-semibold hover:bg-blue-700 transition-colors disabled:opacity-50 text-sm">
            {saving ? 'Menyimpan...' : editing ? 'Simpan Perubahan' : 'Buat Target'}
          </button>
        </form>
      </Modal>

      {/* Deposit/Withdraw Modal */}
      <Modal open={fundModalOpen} onClose={() => setFundModalOpen(false)} title={fundMode === 'deposit' ? `Tambah ke "${fundGoal?.name}"` : `Tarik dari "${fundGoal?.name}"`}>
        <form onSubmit={handleFund} className="space-y-4">
          {fundError && <div className="bg-red-50 text-red-600 text-sm rounded-xl px-4 py-3 border border-red-100">{fundError}</div>}
          {fundGoal && (
            <div className="bg-gray-50 rounded-xl p-3.5 sm:p-4">
              <div className="flex justify-between text-xs sm:text-sm text-gray-500 mb-1">
                <span>Terkumpul</span>
                <span className="font-medium text-gray-700">Rp {formatRupiah(fundGoal.current_amount)}</span>
              </div>
              <div className="flex justify-between text-xs sm:text-sm text-gray-500">
                <span>Target</span>
                <span className="font-medium text-gray-700">Rp {formatRupiah(fundGoal.target_amount)}</span>
              </div>
            </div>
          )}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">{fundMode === 'deposit' ? 'Jumlah Ditambahkan' : 'Jumlah Ditarik'} (Rp)</label>
            <input type="number" value={fundAmount} onChange={(e) => setFundAmount(e.target.value)} placeholder="500000" min="1" className="w-full border border-gray-200 rounded-xl px-4 py-3 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm font-semibold" autoFocus inputMode="numeric" />
            <div className="grid grid-cols-4 gap-2 mt-2.5">
              {[100000, 500000, 1000000, 5000000].map((amt) => (
                <button key={amt} type="button" onClick={() => setFundAmount(amt.toString())} className="py-2 bg-gray-100 hover:bg-gray-200 rounded-xl text-xs font-medium text-gray-600 transition-colors">
                  {amt >= 1000000 ? `${amt / 1000000}jt` : `${amt / 1000}rb`}
                </button>
              ))}
            </div>
          </div>
          {fundGoal && fundAmount && (
            <div className="bg-blue-50 rounded-xl p-3.5 text-sm">
              <p className="text-blue-700">
                {fundMode === 'deposit' ? 'Setelah ditambah: ' : 'Setelah ditarik: '}
                <span className="font-bold">Rp {formatRupiah(fundMode === 'deposit' ? fundGoal.current_amount + Number(fundAmount) : Math.max(fundGoal.current_amount - Number(fundAmount), 0))}</span>
                <span className="text-blue-400"> / Rp {formatRupiah(fundGoal.target_amount)}</span>
              </p>
            </div>
          )}
          <button type="submit" disabled={saving} className={`w-full text-white py-3 rounded-xl font-semibold transition-colors disabled:opacity-50 text-sm ${fundMode === 'deposit' ? 'bg-green-600 hover:bg-green-700' : 'bg-amber-500 hover:bg-amber-600'}`}>
            {saving ? 'Memproses...' : fundMode === 'deposit' ? 'Tambah Tabungan' : 'Tarik Tabungan'}
          </button>
        </form>
      </Modal>
    </div>
  );
}
