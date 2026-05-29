import React, { useState } from 'react';
import { ChevronDown, ChevronUp, Database, Star, Flame, Tag } from 'lucide-react';
import { formatSimilarity, formatStars } from '../utils/formatters';

const RetrievalResults = ({ results }) => {
  const [isOpen, setIsOpen] = useState(true);

  if (!results || results.length === 0) return null;

  const toggleOpen = () => setIsOpen(!isOpen);

  return (
    <div className="w-full max-w-7xl mx-auto px-4 mb-8">
      <div className="glass-panel rounded-2xl overflow-hidden glow-indigo">
        {/* Interactive Header for Collapse/Expand */}
        <button
          onClick={toggleOpen}
          className="w-full px-6 py-4 flex items-center justify-between bg-slate-900/40 hover:bg-slate-900/70 border-b border-slate-900/60 transition-colors focus:outline-none cursor-pointer"
        >
          <div className="flex items-center gap-3">
            <div className="p-2 rounded-lg bg-indigo-950/50 border border-indigo-500/20">
              <Database className="w-5 h-5 text-indigo-400" />
            </div>
            <div className="text-left">
              <h3 className="text-base font-bold text-slate-200">
                Dokumen Pembanding (RAG Retrieval Results)
              </h3>
              <p className="text-xs text-slate-400 font-medium">
                Ditemukan {results.length} ulasan serupa di Vector Database
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-xs font-semibold text-indigo-400 bg-indigo-950/60 border border-indigo-500/20 px-2 py-0.5 rounded-md hidden sm:inline-block">
              {isOpen ? 'Sembunyikan' : 'Tampilkan'}
            </span>
            {isOpen ? (
              <ChevronUp className="w-5 h-5 text-slate-400" />
            ) : (
              <ChevronDown className="w-5 h-5 text-slate-400" />
            )}
          </div>
        </button>

        {/* Collapsible Content */}
        {isOpen && (
          <div className="p-6 bg-slate-950/20">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {results.map((item, index) => {
                const isPalsu = item.label?.toLowerCase() === 'palsu';
                
                return (
                  <div
                    key={item.id || index}
                    className="glass-panel-hover glass-panel rounded-xl p-5 border border-slate-900 relative overflow-hidden flex flex-col justify-between"
                  >
                    {/* Rank Indicator and Similarity Tag */}
                    <div className="flex justify-between items-center mb-3">
                      <div className="flex items-center gap-2">
                        <span className="flex items-center justify-center w-6 h-6 rounded-md bg-indigo-950 border border-indigo-500/20 text-xs font-bold text-indigo-400">
                          #{index + 1}
                        </span>
                        {item.product_name && (
                          <span className="text-xs font-bold text-slate-400 truncate max-w-[150px]">
                            {item.product_name}
                          </span>
                        )}
                      </div>
                      
                      {/* Similarity Badge */}
                      <span className="inline-flex items-center gap-1 text-xs font-bold text-slate-300 bg-slate-900/90 px-2.5 py-1 rounded-lg border border-slate-800">
                        <Flame className="w-3 h-3 text-orange-400" />
                        <span>Sim: {formatSimilarity(item.similarity)}</span>
                      </span>
                    </div>

                    {/* Review Clean Content */}
                    <p className="text-sm text-slate-300 leading-relaxed font-medium mb-4 italic flex-grow">
                      "{item.clean_review || '-'}"
                    </p>

                    {/* Footer Details (Rating and Label) */}
                    <div className="pt-3 border-t border-slate-900/80 flex justify-between items-center mt-auto">
                      {/* Stars Rating */}
                      <div className="flex items-center gap-1.5">
                        <span className="text-[10px] font-bold text-slate-500 uppercase tracking-wider">Rating:</span>
                        <span className="text-xs" title={`Rating ${item.rating}/5`}>
                          {formatStars(item.rating)}
                        </span>
                      </div>

                      {/* Label Badge */}
                      <div className="flex items-center gap-1.5">
                        <span className="text-[10px] font-bold text-slate-500 uppercase tracking-wider">Status:</span>
                        <span
                          className={`inline-flex items-center gap-1 text-xs font-bold px-2 py-0.5 rounded-md ${
                            isPalsu
                              ? 'bg-rose-500/10 border border-rose-500/20 text-rose-400'
                              : 'bg-emerald-500/10 border border-emerald-500/20 text-emerald-400'
                          }`}
                        >
                          <Tag className="w-3 h-3" />
                          {item.label || '-'}
                        </span>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default RetrievalResults;
