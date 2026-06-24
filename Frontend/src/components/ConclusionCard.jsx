import React from 'react';
import { Info, HelpCircle, CheckCircle2, AlertTriangle, AlertCircle } from 'lucide-react';

const ConclusionCard = ({ analysis, judge }) => {
  if (!analysis || !judge) return null;

  const { prediction_label } = analysis;
  const { judge_verdict } = judge;

  const isValid = judge_verdict?.toLowerCase() === 'valid';
  const isPalsu = prediction_label?.toLowerCase() === 'palsu';

  // Determine conclusion text, icon, and colors
  let cardTitle = '';
  let cardDescription = '';
  let themeClass = '';
  let icon = null;

  if (isValid) {
    cardTitle = 'Analisis Selesai & Terverifikasi';
    themeClass = isPalsu
      ? 'glow-rose border-rose-200 dark:border-rose-900/30 bg-rose-50 dark:bg-rose-950/10'
      : 'glow-emerald border-emerald-200 dark:border-emerald-900/30 bg-emerald-50 dark:bg-emerald-950/10';
    icon = isPalsu ? (
      <AlertCircle className="w-10 h-10 text-rose-600 dark:text-rose-400 shrink-0" />
    ) : (
      <CheckCircle2 className="w-10 h-10 text-emerald-600 dark:text-emerald-400 shrink-0" />
    );

    cardDescription = (
      <span>
        Review ini terindikasi{' '}
        <span className={`font-extrabold px-1.5 py-0.5 rounded ${isPalsu ? 'text-rose-600 bg-rose-50 dark:text-rose-400 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-500/20' : 'text-emerald-600 bg-emerald-50 dark:text-emerald-400 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-500/20'}`}>
          {prediction_label || '-'}
        </span>{' '}
        dan hasil validasi AI Judge{' '}
        <span className="font-semibold text-slate-800 dark:text-slate-100">mendukung prediksi sistem</span> secara kuat.
      </span>
    );
  } else {
    // judge_verdict is NOT Valid ("Tidak Valid")
    cardTitle = 'Rekomendasi Peninjauan Ulang';
    themeClass = 'glow-amber border-amber-200 dark:border-amber-900/30 bg-amber-50 dark:bg-amber-950/10';
    icon = <AlertTriangle className="w-10 h-10 text-amber-600 dark:text-amber-400 shrink-0" />;

    cardDescription = (
      <span>
        Review ini terindikasi{' '}
        <span className={`font-extrabold px-1.5 py-0.5 rounded ${isPalsu ? 'text-rose-600 bg-rose-50 dark:text-rose-400 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-500/20' : 'text-emerald-600 bg-emerald-50 dark:text-emerald-400 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-500/20'}`}>
          {prediction_label || '-'}
        </span>
        , namun hasil <span className="text-amber-700 dark:text-amber-400 font-bold underline decoration-amber-500/40 underline-offset-4">perlu ditinjau kembali</span> karena validasi AI Judge menunjukkan dukungan konteks masih{' '}
        <span className="font-semibold text-slate-800 dark:text-slate-100">lemah atau tidak selaras</span>.
      </span>
    );
  }

  return (
    <div className={`glass-panel rounded-2xl p-6 transition-all duration-300 ${themeClass} flex flex-col h-full justify-between`}>
      <div className="space-y-4">
        {/* Card Header */}
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-xl bg-slate-50 border border-slate-200 dark:bg-slate-900/80 dark:border-slate-800">
            {icon}
          </div>
          <div>
            <span className="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider block">KESIMPULAN AKHIR</span>
            <h3 className="text-lg font-bold text-slate-800 dark:text-slate-200">{cardTitle}</h3>
          </div>
        </div>

        {/* Dynamic User-Friendly Text */}
        <div className="bg-white dark:bg-slate-950/60 rounded-xl p-4 md:p-5 border border-slate-200 dark:border-slate-900 text-slate-700 dark:text-slate-300 text-sm md:text-base leading-relaxed font-medium">
          {cardDescription}
        </div>
      </div>

      {/* Simplified User-Friendly Guide/Tip */}
      <div className="mt-6 pt-4 border-t border-slate-200 dark:border-slate-900/60 flex items-start gap-2.5 text-xs text-slate-500 dark:text-slate-400">
        <Info className="w-4 h-4 text-indigo-600 dark:text-indigo-400 shrink-0 mt-0.5" />
        <p className="leading-normal">
          {isValid
            ? 'Hasil analisis memiliki keandalan tinggi karena RAG database memiliki ulasan pembanding yang relevan dan disetujui oleh evaluator LLM.'
            : 'Perbedaan prediksi dapat terjadi jika ulasan memiliki struktur bahasa ganda atau jika ulasan pembanding di Vector Database kurang merepresentasikan konteks.'}
        </p>
      </div>
    </div>
  );
};

export default ConclusionCard;
