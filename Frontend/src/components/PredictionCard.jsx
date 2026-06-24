import React from 'react';
import { ShieldCheck, ShieldAlert, Cpu, Sparkles } from 'lucide-react';
import { formatConfidence } from '../utils/formatters';

const PredictionCard = ({ analysis }) => {
  if (!analysis) return null;

  const { prediction_label, confidence, confidence_score, reasoning } = analysis;

  const isPalsu = prediction_label?.toLowerCase() === 'palsu';
  
  // Theme colors based on predictions
  const theme = isPalsu
    ? {
        accentClass: 'text-rose-600 dark:text-rose-400 border-rose-200 dark:border-rose-500/20 bg-rose-50 dark:bg-rose-950/20',
        badgeBg: 'bg-rose-50 dark:bg-rose-500/10 border border-rose-200 dark:border-rose-500/30 text-rose-600 dark:text-rose-400',
        progressBar: 'bg-gradient-to-r from-orange-500 to-rose-500',
        glowClass: 'glow-rose border-rose-200 dark:border-rose-900/30',
        icon: <ShieldAlert className="w-8 h-8 text-rose-600 dark:text-rose-400" />,
        shadowText: 'shadow-rose-500/25',
      }
    : {
        accentClass: 'text-emerald-600 dark:text-emerald-400 border-emerald-200 dark:border-emerald-500/20 bg-emerald-50 dark:bg-emerald-950/20',
        badgeBg: 'bg-emerald-50 dark:bg-emerald-500/10 border border-emerald-200 dark:border-emerald-500/30 text-emerald-600 dark:text-emerald-400',
        progressBar: 'bg-gradient-to-r from-teal-500 to-emerald-500',
        glowClass: 'glow-emerald border-emerald-200 dark:border-emerald-900/30',
        icon: <ShieldCheck className="w-8 h-8 text-emerald-600 dark:text-emerald-400" />,
        shadowText: 'shadow-emerald-500/25',
      };

  const confidenceLabel = confidence ? confidence.charAt(0).toUpperCase() + confidence.slice(1) : '-';

  return (
    <div className={`glass-panel rounded-2xl p-6 transition-all duration-300 ${theme.glowClass} flex flex-col h-full`}>
      {/* Header Info */}
      <div className="flex justify-between items-start mb-6">
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-xl bg-slate-50 border border-slate-200 dark:bg-slate-900/80 dark:border-slate-800">
            {theme.icon}
          </div>
          <div>
            <span className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider block">PREDIKSI SISTEM</span>
            <span className={`text-2xl font-extrabold tracking-tight ${theme.accentClass.split(' ')[0]}`}>
              Ulasan {prediction_label || '-'}
            </span>
          </div>
        </div>
        
        {/* Confidence Badge */}
        <span className={`text-xs font-semibold px-2.5 py-1 rounded-full ${theme.badgeBg}`}>
          Confidence: {confidenceLabel}
        </span>
      </div>

      {/* Progress Bar Container */}
      <div className="mb-6 bg-slate-50 dark:bg-slate-950/50 p-4 rounded-xl border border-slate-200 dark:border-slate-900">
        <div className="flex justify-between items-center text-xs font-semibold text-slate-500 dark:text-slate-400 mb-2">
          <span>Skor Keyakinan (Confidence Score)</span>
          <span className={`font-mono text-sm ${theme.accentClass.split(' ')[0]}`}>
            {formatConfidence(confidence_score)}
          </span>
        </div>
        <div className="w-full bg-slate-100 dark:bg-slate-900 rounded-full h-3 overflow-hidden border border-slate-200 dark:border-slate-800 p-[1px]">
          <div
            className={`h-full rounded-full transition-all duration-1000 ${theme.progressBar} shadow-md ${theme.shadowText}`}
            style={{ width: `${Math.min(100, Math.max(0, confidence_score || 0))}%` }}
          />
        </div>
      </div>

      {/* Reasoning Section */}
      <div className="flex-grow flex flex-col">
        <h4 className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-2.5 flex items-center gap-1.5">
          <Cpu className="w-3.5 h-3.5 text-indigo-600 dark:text-indigo-400" />
          Alasan Klasifikasi
        </h4>
        <div className="flex-grow bg-white dark:bg-slate-950/40 border border-slate-200 dark:border-slate-900/60 rounded-xl p-4 text-slate-700 dark:text-slate-300 text-sm leading-relaxed max-h-[180px] overflow-y-auto font-medium">
          {reasoning || 'Tidak ada rincian alasan yang tersedia dari model.'}
        </div>
      </div>
    </div>
  );
};

export default PredictionCard;
