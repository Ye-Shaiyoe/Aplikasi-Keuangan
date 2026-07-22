import { X } from 'lucide-react';
import { useEffect, useRef } from 'react';

export default function Modal({ open, onClose, title, children }) {
  const contentRef = useRef(null);

  useEffect(() => {
    document.body.style.overflow = open ? 'hidden' : '';
    return () => { document.body.style.overflow = ''; };
  }, [open]);

  // Close on Escape
  useEffect(() => {
    if (!open) return;
    const handler = (e) => { if (e.key === 'Escape') onClose(); };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, [open, onClose]);

  if (!open) return null;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 z-[60] bg-black/50 backdrop-blur-[3px] animate-[fadeIn_0.18s_ease-out]"
        onClick={onClose}
      />

      {/* Wrapper */}
      <div className="fixed inset-0 z-[70] flex items-end sm:items-center justify-center pointer-events-none">

        {/* ── Desktop centered modal ── */}
        <div className="hidden sm:flex w-full max-w-md pointer-events-auto animate-[slideUp_0.22s_ease-out]">
          <div
            ref={contentRef}
            className="w-full bg-white dark:bg-gray-800 rounded-3xl shadow-2xl shadow-black/20 border border-gray-100 dark:border-gray-700 max-h-[88vh] overflow-hidden flex flex-col"
          >
            {/* Header */}
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-100 dark:border-gray-700 shrink-0">
              <h2 className="text-base font-semibold text-gray-800 dark:text-gray-100">{title}</h2>
              <button
                onClick={onClose}
                className="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-xl transition-colors text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
              >
                <X size={18} />
              </button>
            </div>
            {/* Body */}
            <div className="p-6 overflow-y-auto">{children}</div>
          </div>
        </div>

        {/* ── Mobile bottom sheet ── */}
        <div className="sm:hidden w-full pointer-events-auto animate-[slideUp_0.25s_cubic-bezier(0.32,0.72,0,1)]">
          <div className="bg-white dark:bg-gray-800 rounded-t-3xl shadow-2xl max-h-[93vh] overflow-hidden flex flex-col border-t border-gray-100 dark:border-gray-700">
            {/* Drag handle */}
            <div className="flex justify-center pt-3 pb-1.5 shrink-0">
              <div className="w-9 h-1 bg-gray-200 dark:bg-gray-600 rounded-full" />
            </div>
            {/* Header */}
            <div className="flex items-center justify-between px-5 py-3 border-b border-gray-100 dark:border-gray-700/50 shrink-0">
              <h2 className="text-sm font-semibold text-gray-800 dark:text-gray-100">{title}</h2>
              <button
                onClick={onClose}
                className="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-xl transition-colors text-gray-400"
              >
                <X size={17} />
              </button>
            </div>
            {/* Body */}
            <div className="px-5 py-4 overflow-y-auto pb-10">{children}</div>
          </div>
        </div>
      </div>
    </>
  );
}
