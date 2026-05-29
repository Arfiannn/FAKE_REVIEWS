import { Award, Cpu, Database, ShieldCheck, Sparkles } from 'lucide-react';

const Header = () => {
  return (
    <header className="relative w-full max-w-7xl mx-auto px-4 pt-8 pb-6 flex flex-col items-center text-center">
      {/* Decorative ambient glowing orbs in the background */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-72 h-72 bg-indigo-600/10 rounded-full blur-3xl -z-10 animate-pulse-slow"></div>
      
      {/* Premium Tech Badge */}
      <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-indigo-950/60 border border-indigo-500/30 text-indigo-300 text-xs font-semibold uppercase tracking-wider mb-4 animate-fade-in shadow-inner">
        <Sparkles className="w-3.5 h-3.5 text-indigo-400 animate-spin" style={{ animationDuration: '4s' }} />
        <span>Hybrid AI System</span>
      </div>

      {/* Main Title with neat gradient */}
      <h1 className="text-4xl md:text-5xl font-extrabold tracking-tight bg-gradient-to-r from-white via-slate-100 to-indigo-300 bg-clip-text text-transparent drop-shadow-sm mb-4">
        Fake Review Detector
      </h1>

      {/* Subtitle describing advanced RAG architecture */}
      <p className="text-slate-400 text-base md:text-lg max-w-2xl font-medium leading-relaxed mb-6">
        Hybrid AI berbasis{' '}
        <span className="text-indigo-400 font-semibold border-b border-indigo-500/20 pb-0.5">RAG</span>,{' '}
        <span className="text-indigo-400 font-semibold border-b border-indigo-500/20 pb-0.5">HyDE</span>, dan{' '}
        <span className="text-indigo-400 font-semibold border-b border-indigo-500/20 pb-0.5">LLM</span>{' '}
      </p>

      {/* Component Tech Badges grid */}
      <div className="flex flex-wrap justify-center gap-3 max-w-xl text-xs font-medium text-slate-300">
        <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-900/60 border border-slate-800">
          <Database className="w-3.5 h-3.5 text-blue-400" />
          <span>Vector Database (RAG)</span>
        </div>
        <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-900/60 border border-slate-800">
          <Cpu className="w-3.5 h-3.5 text-purple-400" />
          <span>HyDE Generation</span>
        </div>
        <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-900/60 border border-slate-800">
          <ShieldCheck className="w-3.5 h-3.5 text-indigo-400" />
          <span>DeepSeek LLM</span>
        </div>
        <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-900/60 border border-slate-800">
          <Award className="w-3.5 h-3.5 text-amber-400" />
          <span>LLM-as-a-Judge</span>
        </div>
      </div>
    </header>
  );
};

export default Header;
