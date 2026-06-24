import { useState, useEffect, useRef } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useTheme } from '../context/ThemeContext';
import {
  LayoutDashboard, ArrowLeftRight, Wallet, BarChart3, LogOut, UserCircle, Target, Sparkles,
  Moon, Sun, Menu, X, BookOpen, Repeat, Table2, ChevronDown,
} from 'lucide-react';

const navItems = [
  { to: '/', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/transactions', label: 'Transaksi', icon: ArrowLeftRight },
  { to: '/budgets', label: 'Anggaran', icon: Wallet },
  { to: '/savings', label: 'Tabungan', icon: Target },
  { to: '/recurring', label: 'Transaksi Berulang', icon: Repeat },
  {
    to: '/insights', label: 'Insight', icon: Sparkles,
    children: [
      { to: '/insights/charts', label: 'Chart Data', icon: BarChart3 },
      { to: '/insights/analytics', label: 'Advanced Analytics', icon: Sparkles },
      { to: '/insights/tables', label: 'Table Data', icon: Table2 },
    ],
  },
  { to: '/reports', label: 'Laporan', icon: BarChart3 },
  { to: '/docs', label: 'Dokumentasi', icon: BookOpen },
];

// Mobile-only bottom navigation structure.
// Desktop uses navItems above; mobile uses a separate layout with dropdowns:
// Dashboard | Transaksi (dropdown) | Anggaran | Insight (dropdown) | Lainnya (dropdown)
const mobileNav = [
  { type: 'link', to: '/', label: 'Dashboard', icon: LayoutDashboard, end: true },
  {
    type: 'dropdown', key: 'transaksi', label: 'Transaksi', icon: ArrowLeftRight,
    items: [
      { to: '/transactions', label: 'Transaksi', icon: ArrowLeftRight },
      { to: '/recurring', label: 'Transaksi Berulang', icon: Repeat },
    ],
  },
  { type: 'link', to: '/budgets', label: 'Anggaran', icon: Wallet },
  {
    type: 'dropdown', key: 'insight', label: 'Insight', icon: Sparkles,
    items: [
      { to: '/insights/charts', label: 'Chart Data', icon: BarChart3 },
      { to: '/insights/analytics', label: 'Advanced Analytics', icon: Sparkles },
      { to: '/insights/tables', label: 'Table Data', icon: Table2 },
    ],
  },
  {
    type: 'dropdown', key: 'lainnya', label: 'Lainnya', icon: Menu,
    items: [
      { to: '/savings', label: 'Tabungan', icon: Target },
      { to: '/reports', label: 'Laporan', icon: BarChart3 },
      { to: '/docs', label: 'Dokumentasi', icon: BookOpen },
    ],
  },
];

export default function Layout({ children }) {
  const { user, logout } = useAuth();
  const { dark, toggle } = useTheme();
  const location = useLocation();
  const [openMenu, setOpenMenu] = useState(null); // 'transaksi' | 'insight' | 'lainnya' | null
  const [openSections, setOpenSections] = useState({});
  const bottomNavRef = useRef(null);

  const isSectionOpen = (to) => {
    const active = location.pathname.startsWith(to);
    return openSections[to] ?? active;
  };
  const toggleSection = (to) => setOpenSections((s) => ({ ...s, [to]: !isSectionOpen(to) }));

  const pageTitles = {
    '/': 'Dashboard',
    '/transactions': 'Transaksi',
    '/categories': 'Kategori',
    '/savings': 'Tabungan',
    '/recurring': 'Transaksi Berulang',
    '/insights/charts': 'Chart Data',
    '/insights/analytics': 'Advanced Analytics',
    '/insights/tables': 'Table Data',
    '/reports': 'Laporan',
    '/docs': 'Dokumentasi',
  };

  const currentTitle = pageTitles[location.pathname] || 'Catatan Keuangan';

  useEffect(() => { setOpenMenu(null); }, [location.pathname]);

  useEffect(() => {
    if (!openMenu) return;
    const handler = (e) => {
      if (bottomNavRef.current && !bottomNavRef.current.contains(e.target)) setOpenMenu(null);
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, [openMenu]);

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
        {navItems.map((item) => {
          if (item.children) {
            const isParentActive = location.pathname.startsWith(item.to);
            const open = isSectionOpen(item.to);
            return (
              <div key={item.to}>
                <button
                  onClick={() => toggleSection(item.to)}
                  className={`relative w-full flex items-center gap-3 px-4 py-2.5 rounded-xl text-[13px] font-medium transition-all ${
                    isParentActive
                      ? 'bg-blue-50 text-blue-700 font-semibold dark:bg-blue-900/40 dark:text-blue-400'
                      : 'text-gray-500 hover:bg-gray-50 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-700/50 dark:hover:text-gray-200'
                  }`}
                >
                  {isParentActive && (
                    <div className="absolute left-0 top-1/2 -translate-y-1/2 w-[3px] h-5 bg-blue-600 dark:bg-blue-400 rounded-r-full" />
                  )}
                  <item.icon size={19} strokeWidth={isParentActive ? 2.2 : 1.8} />
                  <span>{item.label}</span>
                  <div className="ml-auto flex items-center gap-1.5">
                    {!isParentActive && (
                      <span className="text-[9px] font-bold bg-gradient-to-r from-amber-400 to-orange-500 text-white px-1.5 py-0.5 rounded-full">NEW</span>
                    )}
                    <ChevronDown size={16} className={`transition-transform ${open ? 'rotate-180' : ''}`} />
                  </div>
                </button>
                {open && (
                  <div className="mt-0.5 ml-3 pl-3 border-l border-gray-100 dark:border-gray-700 space-y-0.5">
                    {item.children.map((child) => (
                      <NavLink
                        key={child.to}
                        to={child.to}
                        onClick={onItemClick}
                        className={({ isActive }) =>
                          `flex items-center gap-3 px-4 py-2 rounded-xl text-[13px] font-medium transition-all ${
                            isActive
                              ? 'bg-blue-50 text-blue-700 font-semibold dark:bg-blue-900/40 dark:text-blue-400'
                              : 'text-gray-500 hover:bg-gray-50 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-700/50 dark:hover:text-gray-200'
                          }`
                        }
                      >
                        {({ isActive }) => (
                          <>
                            <child.icon size={17} strokeWidth={isActive ? 2.2 : 1.8} />
                            <span>{child.label}</span>
                          </>
                        )}
                      </NavLink>
                    ))}
                  </div>
                )}
              </div>
            );
          }
          return (
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
                  {item.to === '/recurring' && !isActive && (
                    <span className="ml-auto text-[9px] font-bold bg-gradient-to-r from-amber-400 to-orange-500 text-white px-1.5 py-0.5 rounded-full">NEW</span>
                  )}
                </>
              )}
            </NavLink>
          );
        })}
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

  const isPathActive = (to) => (to === '/' ? location.pathname === '/' : location.pathname.startsWith(to));

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex transition-colors overflow-x-hidden">
      <aside className="hidden md:flex md:flex-col w-64 bg-white dark:bg-gray-800 border-r border-gray-100 dark:border-gray-700 fixed top-0 bottom-0 left-0 z-30 transition-colors print:hidden">
        <SidebarContent onItemClick={() => {}} />
      </aside>

      <header className="md:hidden fixed top-0 left-0 right-0 z-50 bg-white/80 dark:bg-gray-800/80 backdrop-blur-xl border-b border-gray-100 dark:border-gray-700 transition-colors print:hidden">
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

      <nav ref={bottomNavRef} className="md:hidden fixed bottom-0 left-0 right-0 z-50 bg-white dark:bg-gray-800 border-t border-gray-100 dark:border-gray-700 safe-area-pb transition-colors print:hidden">
        <div className="grid grid-cols-5 h-16 relative">
          {mobileNav.map((entry, idx) => {
            if (entry.type === 'link') {
              const isActive = isPathActive(entry.to);
              return (
                <NavLink
                  key={entry.to}
                  to={entry.to}
                  end={entry.end}
                  className={`relative flex flex-col items-center justify-center gap-0.5 transition-colors ${
                    isActive ? 'text-blue-600 dark:text-blue-400' : 'text-gray-400 dark:text-gray-500'
                  }`}
                >
                  <div className={`p-1 rounded-lg transition-all ${isActive ? 'scale-110' : ''}`}>
                    <entry.icon size={22} strokeWidth={isActive ? 2.5 : 1.8} />
                  </div>
                  <span className={`text-[10px] ${isActive ? 'font-semibold' : 'font-medium'}`}>{entry.label}</span>
                  {isActive && <div className="absolute top-0 w-8 h-0.5 bg-blue-600 dark:bg-blue-400 rounded-full" />}
                </NavLink>
              );
            }
            // dropdown tab
            const isActive = entry.items.some((i) => isPathActive(i.to));
            const isOpen = openMenu === entry.key;
            const alignRight = idx >= 2;
            return (
              <div key={entry.key} className="relative">
                <button
                  onClick={() => setOpenMenu(isOpen ? null : entry.key)}
                  className={`relative w-full h-full flex flex-col items-center justify-center gap-0.5 transition-colors ${
                    isActive || isOpen ? 'text-blue-600 dark:text-blue-400' : 'text-gray-400 dark:text-gray-500'
                  }`}
                >
                  <div className={`p-1 rounded-lg transition-all ${isOpen ? 'scale-110' : ''}`}>
                    {isOpen ? <X size={22} strokeWidth={2} /> : <entry.icon size={22} strokeWidth={isActive ? 2.5 : 1.8} />}
                  </div>
                  <span className={`text-[10px] ${isActive || isOpen ? 'font-semibold' : 'font-medium'}`}>{entry.label}</span>
                  {(isActive || isOpen) && <div className="absolute top-0 w-8 h-0.5 bg-blue-600 dark:bg-blue-400 rounded-full" />}
                </button>

                {isOpen && (
                  <div className={`absolute bottom-full mb-2 w-48 bg-white dark:bg-gray-700 rounded-2xl shadow-xl border border-gray-100 dark:border-gray-600 overflow-hidden ${alignRight ? 'right-0' : 'left-0'}`}>
                    <div className="p-1.5">
                      <div className="px-3.5 pt-1 pb-1.5 text-[10px] font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500">
                        {entry.label}
                      </div>
                      {entry.items.map((child) => {
                        const childActive = isPathActive(child.to);
                        return (
                          <NavLink
                            key={child.to}
                            to={child.to}
                            end={child.to === '/'}
                            onClick={() => setOpenMenu(null)}
                            className={`flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-colors ${
                              childActive
                                ? 'bg-blue-50 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400'
                                : 'text-gray-600 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-600'
                            }`}
                          >
                            <child.icon size={18} strokeWidth={childActive ? 2.2 : 1.8} />
                            <span>{child.label}</span>
                          </NavLink>
                        );
                      })}
                    </div>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </nav>

      <main className="flex-1 md:ml-64 print:ml-0 overflow-x-hidden">
        <div className="pt-14 pb-20 md:pt-0 md:pb-0 print:pt-0 print:pb-0">
          <div className="max-w-7xl mx-auto px-3 py-4 md:px-8 md:py-8 print:p-0 overflow-x-hidden">{children}</div>
        </div>
      </main>
    </div>
  );
}
