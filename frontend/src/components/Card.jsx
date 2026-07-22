export default function Card({ children, className = '', hover = false }) {
  return (
    <div
      className={`
        bg-white dark:bg-gray-800/90
        rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700/60
        p-5 sm:p-6 transition-all duration-200
        ${hover ? 'hover:shadow-md hover:-translate-y-0.5 cursor-pointer' : ''}
        ${className}
      `}
    >
      {children}
    </div>
  );
}

export function StatCard({ title, value, icon: Icon, color = 'blue', trend }) {
  const colorMap = {
    blue:   { icon: 'bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400',   val: 'text-blue-600 dark:text-blue-400',   glow: 'shadow-blue-500/20' },
    green:  { icon: 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400', val: 'text-emerald-600 dark:text-emerald-400', glow: 'shadow-emerald-500/20' },
    red:    { icon: 'bg-red-50 dark:bg-red-900/30 text-red-500 dark:text-red-400',        val: 'text-red-500 dark:text-red-400',      glow: 'shadow-red-500/20' },
    purple: { icon: 'bg-purple-50 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400', val: 'text-purple-600 dark:text-purple-400', glow: 'shadow-purple-500/20' },
  };

  const c = colorMap[color] || colorMap.blue;

  const fmt = (val) => {
    if (val >= 1_000_000_000) return `${(val / 1_000_000_000).toFixed(1)}M`;
    if (val >= 1_000_000)     return `${(val / 1_000_000).toFixed(val % 1_000_000 === 0 ? 0 : 1)}jt`;
    if (val >= 1_000)         return `${(val / 1_000).toFixed(0)}rb`;
    return new Intl.NumberFormat('id-ID').format(val);
  };

  return (
    <div className="group bg-white dark:bg-gray-800/90 rounded-2xl border border-gray-100 dark:border-gray-700/60 p-3 sm:p-5 transition-all duration-200 hover:shadow-md hover:-translate-y-0.5">
      {/* Mobile layout */}
      <div className="flex flex-col sm:hidden">
        <div className="flex items-center justify-between mb-2">
          <p className="text-[10px] font-medium text-gray-400 dark:text-gray-500 uppercase tracking-wide">{title}</p>
          {Icon && (
            <div className={`p-1.5 rounded-xl ${c.icon} transition-transform group-hover:scale-110`}>
              <Icon size={16} />
            </div>
          )}
        </div>
        <p className={`text-base font-bold ${c.val}`}>{fmt(value)}</p>
      </div>

      {/* Desktop layout */}
      <div className="hidden sm:flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-gray-400 dark:text-gray-500 mb-1">{title}</p>
          <p className="text-2xl font-bold text-gray-800 dark:text-gray-100">
            {new Intl.NumberFormat('id-ID').format(value)}
          </p>
          {trend !== undefined && (
            <p className={`text-xs mt-1 ${trend >= 0 ? 'text-emerald-500' : 'text-red-400'}`}>
              {trend >= 0 ? '▲' : '▼'} {Math.abs(trend)}% dari bulan lalu
            </p>
          )}
        </div>
        {Icon && (
          <div className={`p-3 rounded-2xl ${c.icon} shadow-lg ${c.glow} transition-transform group-hover:scale-110`}>
            <Icon size={26} />
          </div>
        )}
      </div>
    </div>
  );
}
