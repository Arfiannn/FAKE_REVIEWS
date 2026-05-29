import React, { useState } from 'react';
import { Lightbulb, Info, HelpCircle } from 'lucide-react';

const HyDECard = ({ hydeDocument }) => {
  const [showExplanation, setShowExplanation] = useState(false);

  if (!hydeDocument) return null;

  return (
    <div className="w-full max-w-7xl mx-auto px-4 mb-8">
      <div className="glass-panel rounded-2xl p-6 glow-indigo relative overflow-hidden">
        {/* Decorative corner glow */}
        <div className="absolute top-0 right-0 w-24 h-24 bg-amber-500/5 rounded-full blur-xl"></div>

        {/* Card Title & Info trigger */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4">
          <div className="flex items-center gap-2.5">
            <div className="p-2 rounded-xl bg-amber-950/40 border border-amber-500/20">
              <Lightbulb className="w-5 h-5 text-amber-400" />
            </div>
            <div>
              <h3 className="text-base font-bold text-slate-200">Dokumen Hipotetik HyDE</h3>
              <p className="text-xs text-slate-400 font-medium">Ulasan ideal buatan LLM untuk memandu pencarian RAG</p>
            </div>
          </div>

          <button
            type="button"
            onClick={() => setShowExplanation(!showExplanation)}
            className="self-start sm:self-center flex items-center gap-1.5 text-xs text-indigo-400 hover:text-indigo-300 font-semibold transition-colors bg-indigo-950/30 hover:bg-indigo-950/60 px-2.5 py-1 rounded-md border border-indigo-500/10 cursor-pointer"
          >
            <HelpCircle className="w-3.5 h-3.5" />
            <span>{showExplanation ? 'Sembunyikan Info' : 'Apa itu HyDE?'}</span>
          </button>
        </div>

        {/* Dynamic Scientific Explanation of HyDE */}
        {showExplanation && (
          <div className="mb-4 bg-slate-900/60 border border-slate-800 rounded-xl p-4 text-xs text-slate-300 leading-relaxed font-medium transition-all duration-300 flex gap-2.5">
            <Info className="w-4 h-4 text-indigo-400 shrink-0 mt-0.5" />
            <div>
              <strong className="text-indigo-300 block mb-0.5">Hypothetical Document Embeddings (HyDE)</strong>
              HyDE bekerja dengan meminta LLM menulis "ulasan palsu/asli hipotetis" terlebih dahulu sebelum melakukan pencarian. Dokumen hipotetis ini kemudian dijadikan representasi vektor untuk mencari ulasan yang mirip di Vector Database. Ini meningkatkan akurasi pencarian secara signifikan dibanding pencarian kata kunci langsung.
            </div>
          </div>
        )}

        {/* Content Box */}
        <div className="bg-slate-950/60 rounded-xl border border-slate-900 p-4 md:p-5 relative">
          <p className="text-sm text-slate-300 leading-relaxed font-medium italic">
            "{hydeDocument}"
          </p>
        </div>
      </div>
    </div>
  );
};

export default HyDECard;
