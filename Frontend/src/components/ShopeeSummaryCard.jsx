import React from 'react';
import { ExternalLink, ShieldCheck, ShieldAlert, CheckCircle2, AlertTriangle, ShoppingBag, BarChart3, HelpCircle } from 'lucide-react';
import { formatConfidence } from '../utils/formatters';

const ShopeeSummaryCard = ({ productUrl, summary }) => {
  if (!summary) return null;

  const {
    total_review = 0,
    total_asli = 0,
    total_palsu = 0,
    percentage_asli = 0,
    percentage_palsu = 0,
    valid_judge = 0,
    need_review_judge = 0,
  } = summary;

  const isHighlyManipulated = percentage_palsu >= 50;

  return (
    <div className="w-full max-w-7xl mx-auto px-4 mb-8 space-y-6 animate-fade-in">
      {/* Upper E-commerce/Shopee Product Identification Card */}
      <div className="glass-panel rounded-2xl p-6 glow-indigo relative overflow-hidden">
        <div className="absolute top-0 bottom-0 left-0 w-[4px] bg-gradient-to-b from-blue-500 to-indigo-500"></div>
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
          <div className="space-y-1">
            <span className="text-[10px] font-extrabold text-indigo-400 uppercase tracking-widest block">
              SUMBER DATA ANALISIS
            </span>
            <div className="flex items-center gap-2">
              <ShoppingBag className="w-5 h-5 text-slate-300" />
              <h2 className="text-lg font-extrabold text-slate-100 tracking-tight truncate max-w-xl">
                Produk Shopee
              </h2>
            </div>
            {productUrl && (
              <a
                href={productUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="text-xs text-indigo-400 hover:text-indigo-300 transition-colors font-semibold flex items-center gap-1 hover:underline"
              >
                <span>{productUrl}</span>
                <ExternalLink className="w-3.5 h-3.5" />
              </a>
            )}
          </div>

          <div className="bg-slate-950/60 border border-slate-900 px-4 py-2.5 rounded-xl flex items-center gap-3">
            <BarChart3 className="w-5 h-5 text-indigo-400" />
            <div>
              <span className="text-[10px] font-bold text-slate-400 block uppercase">
                Ulasan Dianalisis
              </span>
              <span className="text-base font-extrabold text-slate-100">{total_review} Ulasan</span>
            </div>
          </div>
        </div>
      </div>

      {/* Main Grid Metrics: Progress indicators and statistic cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        
        {/* Card 1: Original Reviews */}
        <div className="glass-panel rounded-2xl p-5 glow-emerald flex flex-col justify-between relative overflow-hidden">
          <div className="absolute -top-12 -right-12 w-28 h-28 bg-emerald-500/5 rounded-full blur-2xl"></div>
          <div className="flex justify-between items-start mb-4">
            <div className="space-y-1">
              <span className="text-[10px] font-bold text-slate-400 uppercase">REVIEW ASLI</span>
              <div className="text-3xl font-extrabold text-emerald-400 font-mono">
                {percentage_asli}%
              </div>
            </div>
            <div className="p-2 rounded-xl bg-emerald-950/40 border border-emerald-500/25">
              <ShieldCheck className="w-6 h-6 text-emerald-400" />
            </div>
          </div>
          
          <div className="space-y-2 mt-2">
            <div className="w-full bg-slate-900 rounded-full h-2 overflow-hidden border border-slate-800">
              <div
                className="h-full bg-gradient-to-r from-teal-500 to-emerald-500 rounded-full transition-all duration-1000"
                style={{ width: `${percentage_asli}%` }}
              ></div>
            </div>
            <div className="flex justify-between items-center text-[11px] text-slate-400 font-bold">
              <span>Total Terindikasi Asli</span>
              <span className="text-emerald-400">{total_asli} ulasan</span>
            </div>
          </div>
        </div>

        {/* Card 2: Fake Reviews Indication */}
        <div className="glass-panel rounded-2xl p-5 glow-rose flex flex-col justify-between relative overflow-hidden">
          <div className="absolute -top-12 -right-12 w-28 h-28 bg-rose-500/5 rounded-full blur-2xl"></div>
          <div className="flex justify-between items-start mb-4">
            <div className="space-y-1">
              <span className="text-[10px] font-bold text-slate-400 uppercase">INDthreshold FAKE REVIEW</span>
              <div className="text-3xl font-extrabold text-rose-400 font-mono">
                {percentage_palsu}%
              </div>
            </div>
            <div className="p-2 rounded-xl bg-rose-950/40 border border-rose-500/25">
              <ShieldAlert className="w-6 h-6 text-rose-400" />
            </div>
          </div>

          <div className="space-y-2 mt-2">
            <div className="w-full bg-slate-900 rounded-full h-2 overflow-hidden border border-slate-800">
              <div
                className="h-full bg-gradient-to-r from-orange-500 to-rose-500 rounded-full transition-all duration-1000"
                style={{ width: `${percentage_palsu}%` }}
              ></div>
            </div>
            <div className="flex justify-between items-center text-[11px] text-slate-400 font-bold">
              <span>Total Terindikasi Palsu</span>
              <span className="text-rose-400">{total_palsu} ulasan</span>
            </div>
          </div>
        </div>

        {/* Card 3: AI Judge Valitadion Stats */}
        <div className="glass-panel rounded-2xl p-5 glow-indigo flex flex-col justify-between relative overflow-hidden">
          <div className="absolute -top-12 -right-12 w-28 h-28 bg-indigo-500/5 rounded-full blur-2xl"></div>
          <div className="flex justify-between items-start mb-4">
            <div className="space-y-1">
              <span className="text-[10px] font-bold text-slate-400 uppercase">REKAP VALIDASI AI JUDGE</span>
              <div className="text-sm font-bold text-slate-200 mt-1">
                Akurasi & Keselarasan
              </div>
            </div>
            <div className="p-2 rounded-xl bg-indigo-950/40 border border-indigo-500/25">
              <CheckCircle2 className="w-6 h-6 text-indigo-400" />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3 mt-2">
            <div className="bg-slate-950/40 border border-slate-900 p-2.5 rounded-xl flex items-center gap-2">
              <div className="w-1.5 h-1.5 rounded-full bg-emerald-500"></div>
              <div>
                <span className="text-[9px] text-slate-400 font-bold uppercase block">MENDUKUNG</span>
                <span className="text-xs font-bold text-emerald-400">{valid_judge} Review</span>
              </div>
            </div>
            <div className="bg-slate-950/40 border border-slate-900 p-2.5 rounded-xl flex items-center gap-2">
              <div className="w-1.5 h-1.5 rounded-full bg-amber-500"></div>
              <div>
                <span className="text-[9px] text-slate-400 font-bold uppercase block">PERLU TINJAU</span>
                <span className="text-xs font-bold text-amber-400">{need_review_judge} Review</span>
              </div>
            </div>
          </div>
        </div>

      </div>

      {/* Executive Conclusion Banner */}
      <div
        className={`glass-panel rounded-2xl p-6 relative overflow-hidden border ${
          isHighlyManipulated
            ? 'glow-rose border-rose-900/30 bg-rose-950/10'
            : 'glow-emerald border-emerald-900/30 bg-emerald-950/10'
        }`}
      >
        <div className="flex items-start gap-4">
          <div
            className={`p-3 rounded-2xl ${
              isHighlyManipulated
                ? 'bg-rose-950/60 border border-rose-500/20 text-rose-400'
                : 'bg-emerald-950/60 border border-emerald-500/20 text-emerald-400'
            }`}
          >
            {isHighlyManipulated ? (
              <ShieldAlert className="w-6 h-6" />
            ) : (
              <ShieldCheck className="w-6 h-6" />
            )}
          </div>
          <div className="space-y-1">
            <span className="text-[10px] font-extrabold text-slate-400 uppercase tracking-widest block">
              KESIMPULAN ANALISIS SHOPEE
            </span>
            <p className="text-sm md:text-base font-extrabold text-slate-200 leading-relaxed">
              {isHighlyManipulated
                ? 'Produk ini memiliki indikasi fake review yang cukup tinggi berdasarkan review yang dianalisis.'
                : 'Mayoritas review yang dianalisis terindikasi asli, namun tetap perhatikan review yang memerlukan tinjauan.'}
            </p>
            <p className="text-xs text-slate-400 leading-normal pt-1.5">
              💡 Keputusan pembelian disarankan didasarkan pada perbandingan ulasan dengan rating rendah (1-3 bintang) dan memperhatikan penilaian AI Judge yang berstatus "Perlu Ditinjau".
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ShopeeSummaryCard;
