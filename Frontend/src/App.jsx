import React, { useState, useEffect } from 'react';
import Header from './components/Header';
import ReviewInputCard from './components/ReviewInputCard';
import SingleReviewResult from './components/SingleReviewResult';
import ShopeeSummaryCard from './components/ShopeeSummaryCard';
import ShopeeReviewResultCard from './components/ShopeeReviewResultCard';
import { analyzeInput } from './api/fakeReviewApi';
import { AlertTriangle, Bot, Sparkles, RefreshCw, Cpu, Layers } from 'lucide-react';

function App() {
  const [result, setResult] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [loadingStep, setLoadingStep] = useState(0);
  const [loadingIsShopee, setLoadingIsShopee] = useState(false);

  const loadingSteps = [
    'Mengekstrak data teks ulasan...',
    'Menghasilkan Dokumen Hipotetis (HyDE) dengan LLM...',
    'Melakukan pencarian kemiripan di Vector Database RAG...',
    'Menghitung skor probabilitas dengan DeepSeek...',
    'Menjalankan Validasi Akhir dengan LLM-as-a-Judge...',
  ];

  // Cycling the high-tech loading phases for maximum visual feedback
  useEffect(() => {
    let interval;
    if (isLoading) {
      setLoadingStep(0);
      interval = setInterval(() => {
        setLoadingStep((prev) => (prev < loadingSteps.length - 1 ? prev + 1 : prev));
      }, 2500);
    }
    return () => clearInterval(interval);
  }, [isLoading]);

  const handleAnalyzeInput = async (input, topK, limit) => {
    setIsLoading(true);
    setError(null);
    setResult(null);

    const isShopee = input.includes('shopee.co.id') || input.includes('shp.ee');
    setLoadingIsShopee(isShopee);

    try {
      const response = await analyzeInput(input, topK, limit);
      
      // Strict response validation
      if (response && response.success && response.data) {
        setResult(response.data);
      } else {
        setError('Gagal memproses hasil analisis dari server.');
      }
    } catch (err) {
      const errMsg = err.message || '';
      // Map Shopee scrap errors as requested
      if (isShopee && (
        errMsg.toLowerCase().includes('scrape') || 
        errMsg.toLowerCase().includes('scraping') || 
        errMsg.toLowerCase().includes('shopee') || 
        errMsg.toLowerCase().includes('fetch') || 
        errMsg.toLowerCase().includes('timeout')
      )) {
        setError('Review dari link Shopee belum berhasil diambil. Pastikan link benar dan coba ulang.');
      } else {
        setError(errMsg || 'Terjadi kesalahan saat menghubungi backend.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen pb-16 flex flex-col justify-start">
      {/* Universal Page Header */}
      <Header />

      {/* Input Section */}
      <ReviewInputCard onAnalyze={handleAnalyzeInput} isLoading={isLoading} />

      {/* Main Content Area */}
      <main className="flex-grow w-full max-w-7xl mx-auto px-4">
        {/* Loading State */}
        {isLoading && (
          <div className="w-full max-w-3xl mx-auto glass-panel rounded-2xl p-12 text-center flex flex-col items-center justify-center glow-indigo shadow-xl mb-8 animate-pulse-slow">
            <div className="p-4 rounded-full bg-indigo-950/60 border border-indigo-500/20 mb-6 relative">
              <div className="absolute inset-0 rounded-full border border-indigo-500/30 animate-ping"></div>
              <RefreshCw className="w-8 h-8 text-indigo-400 animate-spin" />
            </div>
            
            <h3 className="text-lg font-bold text-slate-100 mb-2">
              {loadingIsShopee ? 'Memproses Link Shopee' : 'Menganalisis Review Manual'}
            </h3>
            
            <p className="text-sm text-indigo-400 font-bold mb-4 px-6 leading-relaxed">
              {loadingIsShopee
                ? 'Mengambil review dari Shopee dan menganalisis hasil fake review...'
                : 'Menganalisis review...'}
            </p>

            <span className="text-xs text-slate-400 font-medium italic mb-5 block">
              Proses aktif: {loadingSteps[loadingStep]}
            </span>
            
            {/* Elegant multi-step progress indicators */}
            <div className="flex gap-2.5 mt-2 justify-center">
              {loadingSteps.map((_, idx) => (
                <div
                  key={idx}
                  className={`h-1.5 rounded-full transition-all duration-500 ${
                    idx === loadingStep
                      ? 'w-8 bg-indigo-500'
                      : idx < loadingStep
                      ? 'w-3 bg-emerald-500/80'
                      : 'w-3 bg-slate-800'
                  }`}
                />
              ))}
            </div>

            {/* Timeout warning notice */}
            {loadingIsShopee && (
              <div className="mt-8 px-4 py-3 rounded-xl bg-amber-500/5 border border-amber-500/20 text-xs text-amber-400/90 font-semibold max-w-md mx-auto leading-normal">
                ⚠️ Proses analisis link Shopee dapat memakan waktu lebih lama, terutama jika jumlah review besar.
              </div>
            )}
          </div>
        )}

        {/* Error Handling Alert */}
        {error && (
          <div className="w-full max-w-3xl mx-auto mb-8 bg-rose-500/10 border border-rose-500/30 rounded-2xl p-5 flex items-start gap-4 shadow-lg shadow-rose-950/20">
            <div className="p-2 rounded-xl bg-rose-950/50 border border-rose-500/20 text-rose-400 shrink-0">
              <AlertTriangle className="w-6 h-6" />
            </div>
            <div className="space-y-1 text-left">
              <h4 className="text-sm font-extrabold text-rose-400 uppercase tracking-wider">Kesalahan Deteksi</h4>
              <p className="text-slate-300 text-sm font-semibold">{error}</p>
            </div>
          </div>
        )}

        {/* Empty State / Welcome Screen */}
        {!isLoading && !result && !error && (
          <div className="w-full max-w-3xl mx-auto glass-panel rounded-2xl p-8 md:p-12 text-center glow-indigo shadow-xl mb-8 relative overflow-hidden">
            {/* Backdrop highlights */}
            <div className="absolute -top-12 -left-12 w-48 h-48 bg-indigo-500/5 rounded-full blur-3xl"></div>
            <div className="absolute -bottom-12 -right-12 w-48 h-48 bg-purple-500/5 rounded-full blur-3xl"></div>

            <div className="flex justify-center mb-6">
              <div className="p-3.5 rounded-2xl bg-indigo-950/60 border border-indigo-500/25">
                <Bot className="w-10 h-10 text-indigo-400" />
              </div>
            </div>

            <h2 className="text-xl md:text-2xl font-extrabold text-slate-200 mb-3 tracking-tight">
              Siap Menganalisis Ulasan Produk Anda
            </h2>
            <p className="text-slate-400 text-sm md:text-base leading-relaxed max-w-lg mx-auto mb-8 font-medium">
              Gunakan teknologi kecerdasan buatan hibrida untuk mendeteksi ulasan manipulatif secara akurat dengan mencocokkan pola RAG, HyDE, dan validasi evaluator LLM independen.
            </p>

            {/* Quick explanation flow cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-left">
              <div className="bg-slate-950/40 border border-slate-900 rounded-xl p-4 flex gap-3">
                <Layers className="w-5 h-5 text-blue-400 shrink-0 mt-0.5" />
                <div>
                  <h4 className="text-xs font-bold text-slate-200 uppercase mb-1">1. RAG Retrieval</h4>
                  <p className="text-[11px] text-slate-400 leading-normal">Mencari ribuan ulasan pembanding di Vector Database.</p>
                </div>
              </div>
              <div className="bg-slate-950/40 border border-slate-900 rounded-xl p-4 flex gap-3">
                <Cpu className="w-5 h-5 text-purple-400 shrink-0 mt-0.5" />
                <div>
                  <h4 className="text-xs font-bold text-slate-200 uppercase mb-1">2. DeepSeek AI</h4>
                  <p className="text-[11px] text-slate-400 leading-normal">Klasifikasi semantik tingkat lanjut untuk mendeteksi anomali kata.</p>
                </div>
              </div>
              <div className="bg-slate-950/40 border border-slate-900 rounded-xl p-4 flex gap-3">
                <Sparkles className="w-5 h-5 text-amber-400 shrink-0 mt-0.5" />
                <div>
                  <h4 className="text-xs font-bold text-slate-200 uppercase mb-1">3. AI Judge</h4>
                  <p className="text-[11px] text-slate-400 leading-normal">Penilaian objektivitas akhir untuk memastikan validasi tanpa bias.</p>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Results Panel with routing logic */}
        {!isLoading && result && (
          <div className="animate-fade-in">
            {result.type === 'single_review' ? (
              <SingleReviewResult result={result.result} />
            ) : (
              <div className="space-y-8">
                <ShopeeSummaryCard productUrl={result.product_url} summary={result.summary} />
                <ShopeeReviewResultCard results={result.results} />
              </div>
            )}
          </div>
        )}
      </main>
    </div>
  );
}

export default App;
