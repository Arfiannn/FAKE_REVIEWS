import { Edit3, Layers, Link as LinkIcon, Play, RefreshCw, Sparkles } from 'lucide-react';
import { useState } from 'react';

const ReviewInputCard = ({ onAnalyze, isLoading }) => {
  const [reviewText, setReviewText] = useState('');
  const [topK, setTopK] = useState(10);
  const [limit, setLimit] = useState(1);

  const exampleReview = 'produk bagus banget kualitas sesuai harga, pengiriman cepat dan barang sesuai deskripsi';
  const exampleLink = 'https://shopee.co.id/ROG-Xbox-Ally-X-(2025)-RC73XA-Z2EA35A3T-%E2%80%93-Black-AMD-Ryzen-AI-Z2-Extreme-Processor-AMD-Radeon-Graphics-24GB-1TB-7Inch-W11--i.698077090.41519305159?extraParams=%7B%7D';

  // Automatically detect input type
  const isShopee = reviewText.includes('shopee.co.id') || reviewText.includes('shp.ee');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!reviewText.trim() || isLoading) return;
    onAnalyze(reviewText, topK, limit);
  };

  return (
    <div className="w-full max-w-3xl mx-auto px-4 mb-8">
      <div className="glass-panel rounded-2xl p-6 md:p-8 glow-indigo shadow-2xl relative overflow-hidden">
        {/* Subtle top border accent */}
        <div className="absolute top-0 left-0 right-0 h-[2px] bg-gradient-to-r from-blue-500 via-indigo-500 to-purple-500"></div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="flex flex-col space-y-2">
            <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-2">
              <label
                htmlFor="review"
                className="text-sm font-semibold text-slate-200 flex items-center gap-2"
              >
                {isShopee ? (
                  <LinkIcon className="w-4 h-4 text-indigo-400 animate-pulse" />
                ) : (
                  <Sparkles className="w-4 h-4 text-indigo-400" />
                )}
                Masukkan Review atau Link Produk Shopee
              </label>
              
              <div className="flex flex-wrap gap-2">
                <button
                  type="button"
                  onClick={() => setReviewText(exampleReview)}
                  disabled={isLoading}
                  className="text-[11px] font-bold text-slate-400 hover:text-indigo-300 transition-colors bg-slate-900/60 hover:bg-slate-900/90 border border-slate-800/80 px-2.5 py-1 rounded-md cursor-pointer flex items-center gap-1"
                >
                  <Edit3 className="w-3 h-3 text-indigo-400" />
                  Gunakan Contoh Review
                </button>
                <button
                  type="button"
                  onClick={() => setReviewText(exampleLink)}
                  disabled={isLoading}
                  className="text-[11px] font-bold text-slate-400 hover:text-indigo-300 transition-colors bg-slate-900/60 hover:bg-slate-900/90 border border-slate-800/80 px-2.5 py-1 rounded-md cursor-pointer flex items-center gap-1"
                >
                  <LinkIcon className="w-3 h-3 text-indigo-400" />
                  Gunakan Contoh Link
                </button>
              </div>
            </div>
            
            <textarea
              id="review"
              rows={5}
              value={reviewText}
              onChange={(e) => setReviewText(e.target.value)}
              disabled={isLoading}
              placeholder={`Contoh review: ${exampleReview}\nContoh link: ${exampleLink}`}
              className="w-full px-4 py-3 rounded-xl bg-slate-950/70 border border-slate-800 text-slate-100 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 transition-all resize-none text-sm leading-relaxed font-medium"
            />
          </div>

          {/* Shopee Limit Selector */}
          {isShopee && (
            <div className="bg-slate-900/40 border border-indigo-950/40 rounded-xl p-5 space-y-4 animate-fade-in relative overflow-hidden">
              <div className="absolute top-0 bottom-0 left-0 w-[3px] bg-indigo-500"></div>
              
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                <div>
                  <h4 className="text-xs font-extrabold text-indigo-400 uppercase tracking-wider mb-0.5">
                    Jumlah Review Shopee
                  </h4>
                  <p className="text-[11px] text-slate-400 font-medium">
                    Pilih atau tentukan jumlah review yang akan ditarik dan dianalisis
                  </p>
                </div>
                
                {/* Preset Chips */}
                <div className="flex items-center gap-2">
                  {[1, 3, 5, 10].map((val) => (
                    <button
                      key={val}
                      type="button"
                      onClick={() => setLimit(val)}
                      disabled={isLoading}
                      className={`text-xs font-bold px-3 py-1 rounded-lg transition-all cursor-pointer ${
                        limit === val
                          ? 'bg-indigo-600 text-white shadow shadow-indigo-600/30'
                          : 'bg-slate-950 hover:bg-slate-900 text-slate-400 hover:text-slate-200 border border-slate-900'
                      }`}
                    >
                      {val}
                    </button>
                  ))}
                </div>
              </div>

              <div className="flex flex-col sm:flex-row sm:items-center gap-4 pt-1">
                <div className="flex items-center gap-3">
                  <span className="text-xs text-slate-400 font-semibold whitespace-nowrap">Input Manual:</span>
                  <input
                    type="number"
                    min={1}
                    max={50}
                    value={limit}
                    onChange={(e) => {
                      let val = parseInt(e.target.value, 10);
                      if (isNaN(val)) val = 1;
                      if (val > 50) val = 50;
                      if (val < 1) val = 1;
                      setLimit(val);
                    }}
                    disabled={isLoading}
                    className="w-20 px-3 py-1.5 rounded-lg bg-slate-950 border border-slate-800 text-slate-100 text-sm font-bold text-center focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500/30 transition-all"
                  />
                  <span className="text-xs text-slate-500 font-medium">
                    (Max 50)
                  </span>
                </div>
                
                <p className="text-[11px] text-amber-500/80 font-bold leading-normal">
                  ⚠️ Semakin banyak review, proses analisis akan lebih lama.
                </p>
              </div>
            </div>
          )}

          <div className="flex flex-col sm:flex-row gap-4 justify-between items-stretch sm:items-center">
            {/* Top-K Select Input */}
            <div className="flex items-center gap-3 bg-slate-900/60 border border-slate-800 rounded-xl px-4 py-2.5 sm:max-w-[200px]">
              <Layers className="w-4 h-4 text-slate-400 shrink-0" />
              <label htmlFor="top-k" className="text-xs font-semibold text-slate-400 whitespace-nowrap">
                Top-K RAG:
              </label>
              <select
                id="top-k"
                value={topK}
                onChange={(e) => setTopK(Number(e.target.value))}
                disabled={isLoading}
                className="bg-transparent text-sm font-bold text-slate-200 outline-none w-full cursor-pointer focus:text-indigo-400"
              >
                {[3, 5, 10, 15, 20].map((k) => (
                  <option key={k} value={k} className="bg-slate-950 text-slate-200">
                    {k}
                  </option>
                ))}
              </select>
            </div>

            {/* Submit Button */}
            <button
              type="submit"
              disabled={!reviewText.trim() || isLoading}
              className={`flex items-center justify-center gap-2.5 px-6 py-3 rounded-xl font-bold transition-all duration-300 select-none ${
                !reviewText.trim() || isLoading
                  ? 'bg-slate-800/40 text-slate-500 border border-slate-800/80 cursor-not-allowed'
                  : 'bg-indigo-600 hover:bg-indigo-500 text-white shadow-lg shadow-indigo-600/30 hover:shadow-indigo-500/40 transform active:scale-95 cursor-pointer'
              }`}
            >
              {isLoading ? (
                <>
                  <RefreshCw className="w-4 h-4 animate-spin text-white" />
                  <span>Menganalisis...</span>
                </>
              ) : (
                <>
                  <Play className="w-4 h-4 fill-current text-white" />
                  <span>Analisis</span>
                </>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ReviewInputCard;
