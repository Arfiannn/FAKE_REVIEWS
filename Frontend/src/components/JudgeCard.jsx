import React from 'react';
import { Award, AlertTriangle, CheckCircle2, FileQuestion } from 'lucide-react';

const JudgeCard = ({ judge }) => {
  if (!judge) return null;

  const { judge_score, judge_verdict, judge_comment } = judge;

  const isValid = judge_verdict?.toLowerCase() === 'valid';

  // Theme styling: Green for Valid (supporting), Amber/Orange for Invalid (needs review)
  const theme = isValid
    ? {
        accentClass: 'text-emerald-600 dark:text-emerald-400 border-emerald-200 dark:border-emerald-500/20 bg-emerald-50 dark:bg-emerald-950/20',
        badgeBg: 'bg-emerald-50 dark:bg-emerald-500/10 border border-emerald-250 dark:border-emerald-500/30 text-emerald-600 dark:text-emerald-400',
        progressBar: 'bg-gradient-to-r from-teal-500 to-emerald-500',
        glowClass: 'glow-emerald border-emerald-200 dark:border-emerald-900/30',
        icon: <CheckCircle2 className="w-8 h-8 text-emerald-600 dark:text-emerald-400 animate-pulse" />,
        shadowText: 'shadow-emerald-500/25',
        statusLabel: 'Validasi Mendukung',
      }
    : {
        accentClass: 'text-amber-600 dark:text-amber-400 border-amber-250 dark:border-amber-500/20 bg-amber-50 dark:bg-amber-950/20',
        badgeBg: 'bg-amber-50 dark:bg-amber-500/10 border border-amber-250 dark:border-amber-500/30 text-amber-600 dark:text-amber-400',
        progressBar: 'bg-gradient-to-r from-yellow-500 to-amber-500',
        glowClass: 'glow-amber border-amber-200 dark:border-amber-900/30',
        icon: <AlertTriangle className="w-8 h-8 text-amber-600 dark:text-amber-400 animate-bounce" style={{ animationDuration: '3s' }} />,
        shadowText: 'shadow-amber-500/25',
        statusLabel: 'Perlu Ditinjau',
      };

  return (
    <div className={`glass-panel rounded-2xl p-6 transition-all duration-300 ${theme.glowClass} flex flex-col h-full`}>
      {/* Header Info */}
      <div className="flex justify-between items-start mb-6">
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-xl bg-slate-50 border border-slate-200 dark:bg-slate-900/80 dark:border-slate-800">
            {theme.icon}
          </div>
          <div>
            <span className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider block">VALIDASI AI JUDGE</span>
            <span className={`text-2xl font-extrabold tracking-tight ${theme.accentClass.split(' ')[0]}`}>
              {theme.statusLabel}
            </span>
          </div>
        </div>
        
        {/* Verdict Badge */}
        <span className={`text-xs font-semibold px-2.5 py-1 rounded-full ${theme.badgeBg}`}>
          Verdict: {judge_verdict || '-'}
        </span>
      </div>

      {/* Progress Bar Container */}
      <div className="mb-6 bg-slate-50 dark:bg-slate-950/50 p-4 rounded-xl border border-slate-200 dark:border-slate-900">
        <div className="flex justify-between items-center text-xs font-semibold text-slate-500 dark:text-slate-400 mb-2">
          <span>Skor Kelayakan (Judge Score)</span>
          <span className={`font-mono text-sm ${theme.accentClass.split(' ')[0]}`}>
            {judge_score !== undefined ? `${judge_score}/100` : '-'}
          </span>
        </div>
        <div className="w-full bg-slate-100 dark:bg-slate-900 rounded-full h-3 overflow-hidden border border-slate-200 dark:border-slate-800 p-[1px]">
          <div
            className={`h-full rounded-full transition-all duration-1000 ${theme.progressBar} shadow-md ${theme.shadowText}`}
            style={{ width: `${Math.min(100, Math.max(0, judge_score || 0))}%` }}
          />
        </div>
      </div>

      {/* Comment Section */}
      <div className="flex-grow flex flex-col">
        <h4 className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-2.5 flex items-center gap-1.5">
          <Award className="w-3.5 h-3.5 text-indigo-600 dark:text-indigo-400" />
          Komentar AI Judge
        </h4>
        <div className="flex-grow bg-white dark:bg-slate-950/40 border border-slate-200 dark:border-slate-900/60 rounded-xl p-4 text-slate-700 dark:text-slate-300 text-sm leading-relaxed max-h-[180px] overflow-y-auto font-medium">
          {judge_comment || 'Tidak ada rincian evaluasi dari AI Judge.'}
        </div>
      </div>
    </div>
  );
};

export default JudgeCard;
