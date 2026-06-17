import { useState, useEffect, useRef } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useTheme } from '../context/ThemeContext';
import {
  LayoutDashboard, ArrowLeftRight, Wallet, BarChart3, LogOut, UserCircle, Target, Sparkles,
  Moon, Sun, Menu, X,
} from 'lucide-react';

const navItems = [
  { to: '/', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/transactions', label: 'Transaksi', icon: ArrowLeftRight },
  { to: '/budgets', label: 'Anggaran', icon: Wallet },
  { to: '/savings', label: 'Tabungan', icon: Target },
  { to: '/insights', label: 'Insight', icon: Sparkles },
  { to: '/reports', label: 'Laporan', icon: BarChart3 },
];

// Items shown in bottom nav (first 5)
const bottomNavItems = navItems.slice(0, 5);
// Extra items shown in hamburger menu
const extraNavItems = navItems.slice(5);

export default function Layout({ children }) {
  const { user, logout } = useAuth();
  const { dark, toggle } = useTheme();
  const location = useLocation();
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef(null);

  const pageTitles = {
    '/': 'Dashboard',
    '/transactions': 'Transaksi',
    '/categories': 'Kategori',
    '/savings': 'Tabungan',
    '/insights': 'Insight',
    '/reports': 'Laporan',
  };

  const currentTitle = pageTitles[location.pathname] || 'Catatan Keuangan';

  useEffect(() => { setMenuOpen(false); }, [location.pathname]);

  useEffect(() => {
    if (!menuOpen) return;
    const handler = (e) => {
      if (menuRef.current && !menuRef.current.contains(e.target)) setMenuOpen(false);
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, [menuOpen]);

  const SidebarContent = ({ onItemClick }) => (
    <>
      <div className="p-5 pb-4">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white font-bold text-lg shadow-md shadow-blue-500/20">
            C
          </div>
          <div>
            <h1 className="text-base font-bold text-gray-800 dark:text-gray-100 tracking-tight">Catatan Keuangan</h1>
            <p className="text-[11px] text-gray-400 dark:text-gray-500">Personal Finance Tracker</p>
          </div>
        </div>
      </div>
      <nav className="flex-1 px-3 space-y-0.5">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.to === '/'}
            onClick={onItemClick}
            className={({ isActive }) =>
              `relative flex items-center gap-3 px-4 py-2.5 rounded-xl text-[13px] font-medium transition-all ${
                isActive
                  ? 'bg-blue-50 text-blue-700 font-semibold dark:bg-blue-900/40 dark:text-blue-400'
                  : 'text-gray-500 hover:bg-gray-50 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-700/50 dark:hover:text-gray-200'
              }`
            }
          >
            {({ isActive }) => (
              <>
                {isActive && (
                  <div className="absolute left-0 top-1/2 -translate-y-1/2 w-[3px] h-5 bg-blue-600 dark:bg-blue-400 rounded-r-full" />
                )}
                <item.icon size={19} strokeWidth={isActive ? 2.2 : 1.8} />
                <span>{item.label}</span>
                {item.to === '/insights' && !isActive && (
                  <span className="ml-auto text-[9px] font-bold bg-gradient-to-r from-amber-400 to-orange-500 text-white px-1.5 py-0.5 rounded-full">NEW</span>
                )}
              </>
            )}
          </NavLink>
        ))}
      </nav>
      <div className="px-3 pb-2">
        <button
          onClick={toggle}
          className="flex items-center gap-3 w-full px-4 py-2.5 rounded-xl text-[13px] font-medium text-gray-500 hover:bg-gray-50 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-700/50 dark:hover:text-gray-200 transition-all"
        >
          {dark ? <Sun size={19} /> : <Moon size={19} />}
          <span>{dark ? 'Mode Terang' : 'Mode Gelap'}</span>
        </button>
      </div>
      <div className="p-3 border-t border-gray-100 dark:border-gray-700">
        <div className="flex items-center gap-3 px-3 py-3 bg-gradient-to-r from-gray-50 to-slate-50 dark:from-gray-800 dark:to-gray-800 rounded-xl border border-gray-100/50 dark:border-gray-700">
          <div className="w-9 h-9 rounded-full bg-gradient-to-br from-blue-400 to-indigo-600 flex items-center justify-center text-white text-sm font-bold shadow-sm shadow-blue-500/20">
            {user?.name?.charAt(0)?.toUpperCase() || 'U'}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-[13px] font-medium text-gray-700 dark:text-gray-200 truncate">{user?.name}</p>
            <p className="text-[11px] text-gray-400 dark:text-gray-500 truncate">{user?.email}</p>
          </div>
          <button onClick={logout} className="p-2 hover:bg-red-50 dark:hover:bg-red-900/30 rounded-lg transition-colors group" title="Keluar">
            <LogOut size={16} className="text-gray-300 dark:text-gray-600 group-hover:text-red-500" />
          </button>
        </div>
      </div>
    </>
  );

  const isExtraActive = extraNavItems.some(i =>
    i.to === '/' ? location.pathname === '/' : location.pathname.startsWith(i.to)
  );

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex transition-colors">
      <aside className="hidden md:flex md:flex-col w-64 bg-white dark:bg-gray-800 border-r border-gray-100 dark:border-gray-700 fixed top-0 bottom-0 left-0 z-30 transition-colors">
        <SidebarContent onItemClick={() => {}} />
      </aside>

      <header className="md:hidden fixed top-0 left-0 right-0 z-50 bg-white/80 dark:bg-gray-800/80 backdrop-blur-xl border-b border-gray-100 dark:border-gray-700 transition-colors">
        <div className="flex items-center justify-between px-4 py-3">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center text-white font-bold text-sm shadow-sm">
              C
            </div>
            <h1 className="text-base font-bold text-gray-800 dark:text-gray-100">{currentTitle}</h1>
          </div>
          <div className="flex items-center gap-1">
            <button onClick={toggle} className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors" title={dark ? 'Mode Terang' : 'Mode Gelap'}>
              {dark ? <Sun size={18} className="text-yellow-400" /> : <Moon size={18} className="text-gray-400" />}
            </button>
            <button onClick={logout} className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors" title="Keluar">
              <LogOut size={18} className="text-gray-400 dark:text-gray-500" />
            </button>
          </div>
        </div>
      </header>

      <nav className="md:hidden fixed bottom-0 left-0 right-0 z-50 bg-white dark:bg-gray-800 border-t border-gray-100 dark:border-gray-700 safe-area-pb transition-colors">
        <div className="grid grid-cols-6 h-16 relative">
          {bottomNavItems.map((item) => {
            const isActive = item.to === '/' ? location.pathname === '/' : location.pathname.startsWith(item.to);
            return (
              <NavLink
                key={item.to}
                to={item.to}
                end={item.to === '/'}
                className={`flex flex-col items-center justify-center gap-0.5 transition-colors ${
                  isActive ? 'text-blue-600 dark:text-blue-400' : 'text-gray-400 dark:text-gray-500'
                }`}
              >
                <div className={`p-1 rounded-lg transition-all ${isActive ? 'scale-110' : ''}`}>
                  <item.icon size={22} strokeWidth={isActive ? 2.5 : 1.8} />
                </div>
                <span className={`text-[10px] ${isActive ? 'font-semibold' : 'font-medium'}`}>{item.label}</span>
                {isActive && <div className="absolute top-0 w-8 h-0.5 bg-blue-600 dark:bg-blue-400 rounded-full" />}
              </NavLink>
            );
          })}

          <div className="relative" ref={menuRef}>
            <button
              onClick={() => setMenuOpen(!menuOpen)}
              className={`w-full h-full flex flex-col items-center justify-center gap-0.5 transition-colors ${
                isExtraActive ? 'text-blue-600 dark:text-blue-400' : 'text-gray-400 dark:text-gray-500'
              }`}
            >
              <div className={`p-1 rounded-lg transition-all ${menuOpen ? 'scale-110' : ''}`}>
                {menuOpen ? <X size={22} strokeWidth={2} /> : <Menu size={22} strokeWidth={1.8} />}
              </div>
              <span className={`text-[10px] ${isExtraActive ? 'font-semibold' : 'font-medium'}`}>Lainnya</span>
              {isExtraActive && !menuOpen && <div className="absolute top-0 w-8 h-0.5 bg-blue-600 dark:bg-blue-400 rounded-full" />}
            </button>

            {menuOpen && (
              <div className="absolute bottom-full right-0 mb-2 w-48 bg-white dark:bg-gray-700 rounded-2xl shadow-xl border border-gray-100 dark:border-gray-600 overflow-hidden">
                <div className="p-1.5">
                  {extraNavItems.map((item) => {
                    const isActive = item.to === '/' ? location.pathname === '/' : location.pathname.startsWith(item.to);
                    return (
                      <NavLink
                        key={item.to}
                        to={item.to}
                        end={item.to === '/'}
                        onClick={() => setMenuOpen(false)}
                        className={`flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-colors ${
                          isActive
                            ? 'bg-blue-50 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400'
                            : 'text-gray-600 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-600'
                        }`}
                      >
                        <item.icon size={18} strokeWidth={isActive ? 2.2 : 1.8} />
                        <span>{item.label}</span>
                      </NavLink>
                    );
                  })}
                </div>
              </div>
            )}
          </div>
        </div>
      </nav>

      <main className="flex-1 md:ml-64">
        <div className="pt-14 pb-20 md:pt-0 md:pb-0">
          <div className="max-w-7xl mx-auto px-3 py-4 md:px-8 md:py-8">{children}</div>
        </div>
      </main>
    </div>
  );
}
