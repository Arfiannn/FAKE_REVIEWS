import { Award, Cpu, Database, ShieldCheck, Sparkles, Sun, Moon } from 'lucide-react';
import { useTheme } from './ThemeContext';

const Header = () => {
  const { theme, toggleTheme } = useTheme();

  return (
    <header className="relative w-full max-w-7xl mx-auto px-4 pt-8 pb-6 flex flex-col items-center text-center">
      {/* Decorative ambient glowing orbs in the background */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-72 h-72 bg-indigo-600/10 rounded-full blur-3xl -z-10 animate-pulse-slow"></div>
      
      {/* Theme Toggle Button */}
      <div className="absolute top-4 right-4 md:top-8 md:right-8 z-50">
        <button
          onClick={toggleTheme}
          className="p-2.5 rounded-xl bg-white/70 hover:bg-white border border-slate-200 dark:bg-slate-900/40 dark:hover:bg-slate-950/80 dark:border-slate-800/80 text-slate-700 hover:text-indigo-600 dark:text-slate-300 dark:hover:text-indigo-400 transition-all duration-300 shadow-md flex items-center justify-center cursor-pointer select-none group relative"
          title={theme === 'dark' ? 'Ubah ke mode terang' : 'Ubah ke mode gelap'}
          aria-label="Toggle Theme"
        >
          {theme === 'dark' ? (
            <Sun className="w-5 h-5 transition-transform duration-500 rotate-0 scale-100 group-hover:rotate-45" />
          ) : (
            <Moon className="w-5 h-5 transition-transform duration-500 rotate-0 scale-100 group-hover:-rotate-12" />
          )}

          {/* Tooltip */}
          <span className="pointer-events-none absolute top-full mt-2 right-0 w-max opacity-0 group-hover:opacity-100 transition-opacity duration-300 bg-slate-900 text-slate-100 dark:bg-slate-100 dark:text-slate-900 text-[10px] font-bold px-2 py-1 rounded shadow border border-slate-800 dark:border-slate-200">
            {theme === 'dark' ? 'Ubah ke mode terang' : 'Ubah ke mode gelap'}
          </span>
        </button>
      </div>

      {/* Premium Tech Badge */}
      <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-indigo-950/20 dark:bg-indigo-950/60 border border-indigo-500/20 dark:border-indigo-500/30 text-indigo-600 dark:text-indigo-300 text-xs font-semibold uppercase tracking-wider mb-4 animate-fade-in shadow-inner">
        <Sparkles className="w-3.5 h-3.5 text-indigo-500 dark:text-indigo-400 animate-spin" style={{ animationDuration: '4s' }} />
        <span>Hybrid AI System</span>
      </div>

      {/* Main Title with neat gradient */}
      <h1 className="text-4xl md:text-5xl font-extrabold tracking-tight bg-gradient-to-r from-slate-900 via-indigo-950 to-indigo-700 dark:from-white dark:via-slate-100 dark:to-indigo-300 bg-clip-text text-transparent drop-shadow-sm mb-4">
        Fake Review Detector
      </h1>

      {/* Subtitle describing advanced RAG architecture */}
      <p className="text-slate-600 dark:text-slate-400 text-base md:text-lg max-w-2xl font-medium leading-relaxed mb-6">
        Hybrid AI berbasis{' '}
        <span className="text-indigo-600 dark:text-indigo-400 font-semibold border-b border-indigo-500/20 pb-0.5">RAG</span>,{' '}
        <span className="text-indigo-600 dark:text-indigo-400 font-semibold border-b border-indigo-500/20 pb-0.5">HyDE</span>, dan{' '}
        <span className="text-indigo-600 dark:text-indigo-400 font-semibold border-b border-indigo-500/20 pb-0.5">LLM</span>{' '}
      </p>

      {/* Component Tech Badges grid */}
      <div className="flex flex-wrap justify-center gap-3 max-w-xl text-xs font-medium text-slate-700 dark:text-slate-300">
        <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-white/60 border border-slate-200/80 dark:bg-slate-900/60 dark:border-slate-800 shadow-sm">
          <Database className="w-3.5 h-3.5 text-blue-500 dark:text-blue-400" />
          <span>Vector Database (RAG)</span>
        </div>
        <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-white/60 border border-slate-200/80 dark:bg-slate-900/60 dark:border-slate-800 shadow-sm">
          <Cpu className="w-3.5 h-3.5 text-purple-500 dark:text-purple-400" />
          <span>HyDE Generation</span>
        </div>
        <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-white/60 border border-slate-200/80 dark:bg-slate-900/60 dark:border-slate-800 shadow-sm">
          <ShieldCheck className="w-3.5 h-3.5 text-indigo-500 dark:text-indigo-400" />
          <span>DeepSeek LLM</span>
        </div>
        <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-white/60 border border-slate-200/80 dark:bg-slate-900/60 dark:border-slate-800 shadow-sm">
          <Award className="w-3.5 h-3.5 text-amber-600 dark:text-amber-400" />
          <span>LLM-as-a-Judge</span>
        </div>
      </div>
    </header>
  );
};

export default Header;
