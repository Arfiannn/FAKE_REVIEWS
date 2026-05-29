import React, { useState } from 'react';
import {
  ChevronDown,
  ChevronUp,
  ShieldCheck,
  ShieldAlert,
  CheckCircle2,
  AlertTriangle,
  User,
  Calendar,
  Store,
  MessageSquare,
  BookOpen,
  Layers,
  Flame,
  Star,
  Quote,
  Cpu,
  Bookmark
} from 'lucide-react';
import { formatSimilarity, formatStars, formatDate } from '../utils/formatters';

// Single Review Card Sub-component for independent state handling
const ShopeeSingleReviewCard = ({ index, review }) => {
  const [isAccordionOpen, setIsAccordionOpen] = useState(false);

  const {
    product_name = '-',
    shop_name = '-',
    username = 'Pembeli Shopee',
    rating = 5,
    date = '',
    raw_review = '',
    clean_review = '',
    analysis,
    judge,
    retrieval_results = [],
  } = review;

  const predictionLabel = analysis?.prediction_label || '-';
  const confidenceScore = analysis?.confidence_score || 0;
  const reasoning = analysis?.reasoning || 'Alasan tidak tersedia.';

  const judgeScore = judge?.judge_score;
  const judgeVerdict = judge?.judge_verdict || '-';
  const judgeComment = judge?.judge_comment || 'Evaluasi AI Judge tidak tersedia.';

  const isPalsu = predictionLabel.toLowerCase() === 'palsu';
  const isJudgeValid = judgeVerdict.toLowerCase() === 'valid';

  return (
    <div className={`glass-panel rounded-2xl p-6 relative overflow-hidden transition-all duration-300 ${
      isPalsu ? 'glow-rose border-rose-950/40' : 'glow-emerald border-emerald-950/40'
    }`}>
      {/* Decorative vertical band indicator */}
      <div className={`absolute top-0 bottom-0 left-0 w-[4px] ${
        isPalsu ? 'bg-rose-500' : 'bg-emerald-500'
      }`}></div>

      {/* Header Info */}
      <div className="flex flex-col lg:flex-row justify-between items-start gap-4 pb-4 border-b border-slate-900/60 mb-5">
        
        {/* User, Shop, Product information */}
        <div className="space-y-2 w-full lg:max-w-[70%]">
          <div className="flex items-center gap-2.5">
            <span className="flex items-center justify-center w-7 h-7 rounded-lg bg-indigo-950 border border-indigo-500/20 text-xs font-black text-indigo-400">
              #{index}
            </span>
            <div className="flex items-center gap-1.5 text-sm font-extrabold text-slate-100">
              <User className="w-3.5 h-3.5 text-indigo-400" />
              <span>{username}</span>
            </div>
            <div className="flex items-center text-xs">
              {formatStars(rating)}
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 text-xs font-semibold text-slate-400">
            <div className="flex items-center gap-1.5 truncate">
              <Store className="w-3.5 h-3.5 text-indigo-400/80 shrink-0" />
              <span>Toko: <strong className="text-slate-300 font-bold">{shop_name}</strong></span>
            </div>
            <div className="flex items-center gap-1.5 truncate">
              <Bookmark className="w-3.5 h-3.5 text-indigo-400/80 shrink-0" />
              <span>Produk: <strong className="text-slate-300 font-bold">{product_name}</strong></span>
            </div>
          </div>
        </div>

        {/* Date / Metadata */}
        <div className="flex items-center gap-1.5 text-xs font-semibold text-slate-400 shrink-0 self-end lg:self-start">
          <Calendar className="w-3.5 h-3.5 text-slate-500" />
          <span>{formatDate(date)}</span>
        </div>

      </div>

      {/* Main Analytical Badges Row */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-5">
        
        {/* Metric 1: System Prediction Box */}
        <div className={`p-4 rounded-xl border flex flex-col justify-between ${
          isPalsu
            ? 'bg-rose-950/10 border-rose-500/10'
            : 'bg-emerald-950/10 border-emerald-500/10'
        }`}>
          <div className="flex justify-between items-center mb-3">
            <span className="text-[10px] font-extrabold text-slate-400 uppercase tracking-wider">
              Prediksi Sistem
            </span>
            <span className={`inline-flex items-center gap-1 text-[11px] font-black px-2.5 py-0.5 rounded-full ${
              isPalsu
                ? 'bg-rose-500/10 border border-rose-500/25 text-rose-400'
                : 'bg-emerald-500/10 border border-emerald-500/25 text-emerald-400'
            }`}>
              {isPalsu ? <ShieldAlert className="w-3 h-3" /> : <ShieldCheck className="w-3 h-3" />}
              {predictionLabel}
            </span>
          </div>

          <div className="space-y-1.5">
            <div className="flex justify-between items-center text-[10px] font-bold text-slate-400">
              <span>Confidence Score</span>
              <span className={isPalsu ? 'text-rose-400' : 'text-emerald-400'}>
                {formatSimilarity(confidenceScore)}
              </span>
            </div>
            <div className="w-full bg-slate-950 rounded-full h-1.5 overflow-hidden border border-slate-900">
              <div
                className={`h-full rounded-full transition-all duration-1000 ${
                  isPalsu ? 'bg-gradient-to-r from-orange-500 to-rose-500' : 'bg-gradient-to-r from-teal-500 to-emerald-500'
                }`}
                style={{ width: `${confidenceScore}%` }}
              ></div>
            </div>
          </div>
        </div>

        {/* Metric 2: AI Judge Validation Box */}
        <div className={`p-4 rounded-xl border flex flex-col justify-between ${
          isJudgeValid
            ? 'bg-emerald-950/10 border-emerald-500/10'
            : 'bg-amber-950/10 border-amber-500/10'
        }`}>
          <div className="flex justify-between items-center mb-3">
            <span className="text-[10px] font-extrabold text-slate-400 uppercase tracking-wider">
              Validasi AI Judge
            </span>
            <span className={`inline-flex items-center gap-1 text-[11px] font-black px-2.5 py-0.5 rounded-full ${
              isJudgeValid
                ? 'bg-emerald-500/10 border border-emerald-500/25 text-emerald-400'
                : 'bg-amber-500/10 border border-amber-500/25 text-amber-400'
            }`}>
              {isJudgeValid ? <CheckCircle2 className="w-3 h-3" /> : <AlertTriangle className="w-3 h-3" />}
              {isJudgeValid ? 'Validasi Mendukung' : 'Perlu Ditinjau'}
            </span>
          </div>

          <div className="space-y-1.5">
            <div className="flex justify-between items-center text-[10px] font-bold text-slate-400">
              <span>Judge Score</span>
              <span className={isJudgeValid ? 'text-emerald-400' : 'text-amber-400'}>
                {judgeScore !== undefined ? `${judgeScore}/100` : '-'}
              </span>
            </div>
            <div className="w-full bg-slate-950 rounded-full h-1.5 overflow-hidden border border-slate-900">
              <div
                className={`h-full rounded-full transition-all duration-1000 ${
                  isJudgeValid ? 'bg-gradient-to-r from-teal-500 to-emerald-500' : 'bg-gradient-to-r from-yellow-500 to-amber-500'
                }`}
                style={{ width: `${judgeScore || 0}%` }}
              ></div>
            </div>
          </div>
        </div>

      </div>

      {/* Review Text Area */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        
        {/* Raw Text Card */}
        <div className="bg-slate-950/50 rounded-xl p-4 border border-slate-900/80">
          <span className="text-[9px] font-extrabold text-indigo-400/90 uppercase tracking-widest block mb-2 flex items-center gap-1">
            <Quote className="w-2.5 h-2.5" /> Ulasan Asli (Raw)
          </span>
          <p className="text-xs text-slate-300 leading-relaxed font-semibold italic">
            "{raw_review || '-'}"
          </p>
        </div>

        {/* Clean Text Card */}
        <div className="bg-slate-950/50 rounded-xl p-4 border border-slate-900/80">
          <span className="text-[9px] font-extrabold text-indigo-400/90 uppercase tracking-widest block mb-2 flex items-center gap-1">
            <BookOpen className="w-2.5 h-2.5" /> Ulasan Terbersih (Preprocessed)
          </span>
          <p className="text-xs text-slate-300 leading-relaxed font-semibold">
            {clean_review || '-'}
          </p>
        </div>

      </div>

      {/* Explanatory Explanations Console Block */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-5">
        
        {/* reasoning */}
        <div className="bg-slate-950/30 rounded-xl border border-slate-900 p-4">
          <span className="text-[9px] font-extrabold text-slate-400 uppercase tracking-widest block mb-1.5 flex items-center gap-1">
            <Cpu className="w-3 h-3 text-indigo-400" /> Alasan Klasifikasi AI
          </span>
          <p className="text-xs text-slate-400 leading-relaxed font-semibold max-h-24 overflow-y-auto pr-1">
            {reasoning}
          </p>
        </div>

        {/* judge comment */}
        <div className="bg-slate-950/30 rounded-xl border border-slate-900 p-4">
          <span className="text-[9px] font-extrabold text-slate-400 uppercase tracking-widest block mb-1.5 flex items-center gap-1">
            <MessageSquare className="w-3 h-3 text-indigo-400" /> Catatan AI Judge
          </span>
          <p className="text-xs text-slate-400 leading-relaxed font-semibold max-h-24 overflow-y-auto pr-1">
            {judgeComment}
          </p>
        </div>

      </div>

      {/* Accordion trigger for Retrieval RAG Results */}
      {retrieval_results && retrieval_results.length > 0 && (
        <div className="pt-2 border-t border-slate-900/60 mt-4">
          <button
            type="button"
            onClick={() => setIsAccordionOpen(!isAccordionOpen)}
            className="w-full flex items-center justify-between py-2 text-xs font-extrabold text-indigo-400 hover:text-indigo-300 transition-colors cursor-pointer select-none focus:outline-none"
          >
            <span className="flex items-center gap-2">
              <Layers className="w-4 h-4" />
              Dokumen Pembanding (RAG Match - {retrieval_results.length} Ulasan)
            </span>
            <span className="flex items-center gap-1.5">
              <span className="text-[10px] font-bold px-2 py-0.5 rounded bg-indigo-950/60 border border-indigo-500/10">
                {isAccordionOpen ? 'Sembunyikan' : 'Tampilkan'}
              </span>
              {isAccordionOpen ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
            </span>
          </button>

          {/* Collapsible Retrieval Content */}
          {isAccordionOpen && (
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-3 pb-1 animate-fade-in">
              {retrieval_results.map((ragItem, ragIdx) => {
                const isRagPalsu = ragItem.label?.toLowerCase() === 'palsu';
                return (
                  <div
                    key={ragIdx}
                    className="bg-slate-950/80 border border-slate-900 rounded-xl p-3.5 flex flex-col justify-between"
                  >
                    <div className="flex justify-between items-center gap-2 mb-2 pb-1.5 border-b border-slate-900">
                      <span className="flex items-center justify-center w-5 h-5 rounded bg-indigo-950 border border-indigo-500/15 text-[10px] font-black text-indigo-400">
                        #{ragIdx + 1}
                      </span>
                      
                      <div className="flex items-center gap-2 shrink-0">
                        {/* Rating stars */}
                        <div className="flex items-center text-[10px] shrink-0">
                          {formatStars(ragItem.rating)}
                        </div>

                        {/* Similarity percentage */}
                        <span className="inline-flex items-center gap-0.5 text-[10px] font-extrabold text-slate-300 bg-slate-900 border border-slate-800 px-1.5 py-0.5 rounded">
                          <Flame className="w-2.5 h-2.5 text-orange-400" />
                          <span>{formatSimilarity(ragItem.similarity)}</span>
                        </span>
                      </div>
                    </div>

                    <p className="text-xs text-slate-300 italic font-semibold leading-relaxed mb-3">
                      "{ragItem.clean_review || '-'}"
                    </p>

                    <div className="flex items-center justify-between pt-1 text-[10px] font-bold text-slate-500">
                      <span>Label Pembanding:</span>
                      <span className={`px-1.5 py-0.5 rounded font-black uppercase text-[9px] ${
                        isRagPalsu ? 'text-rose-400 bg-rose-950/20' : 'text-emerald-400 bg-emerald-950/20'
                      }`}>
                        {ragItem.label || '-'}
                      </span>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      )}

    </div>
  );
};

// Main List Container Component
const ShopeeReviewResultCard = ({ results }) => {
  if (!results || results.length === 0) {
    return (
      <div className="w-full max-w-7xl mx-auto px-4 mb-8 text-center py-10 bg-slate-900/30 border border-slate-800 rounded-2xl">
        <p className="text-slate-400 text-sm font-semibold">Tidak ada ulasan Shopee yang ditemukan.</p>
      </div>
    );
  }

  return (
    <div className="w-full max-w-7xl mx-auto px-4 mb-8 space-y-6 animate-fade-in">
      <div className="flex items-center gap-2 mb-2">
        <MessageSquare className="w-5 h-5 text-indigo-400" />
        <h3 className="text-base font-extrabold text-slate-200 uppercase tracking-wider">
          Daftar Rincian Review Shopee ({results.length})
        </h3>
      </div>
      <div className="space-y-6">
        {results.map((review, idx) => (
          <ShopeeSingleReviewCard key={idx} index={idx + 1} review={review} />
        ))}
      </div>
    </div>
  );
};

export default ShopeeReviewResultCard;
