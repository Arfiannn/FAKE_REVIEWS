import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 900000, // 15 minutes (RAG retrieval + DeepSeek + LLM-as-a-Judge can take time)
});

/**
 * Classify a product review using the backend RAG + HyDE + LLM-as-a-Judge pipeline.
 * @param {string} review - The text of the review.
 * @param {number} topK - Top-K retrieval documents.
 * @returns {Promise<object>} The classified review result data.
 */
export const classifyReview = async (review, topK = 10) => {
  try {
    const response = await api.post('/rag/classify', {
      review,
      top_k: Number(topK),
    });
    return response.data;
  } catch (error) {
    // Provide structured error message for robust UI error handling
    if (error.response) {
      // Server responded with a status code outside the 2xx range
      throw new Error(error.response.data?.message || `Error Server: ${error.response.status}`);
    } else if (error.request) {
      // Request was made but no response was received
      throw new Error('Tidak dapat terhubung ke backend. Pastikan server Golang Anda berjalan di http://localhost:8080.');
    } else {
      // Something else happened setting up the request
      throw new Error(`Gagal memproses request: ${error.message}`);
    }
  }
};

/**
 * Analyze an input that can be either manual review text or a Shopee product link.
 * @param {string} input - Review text or Shopee product link.
 * @param {number} topK - Top-K RAG retrieval documents.
 * @param {number} limit - Review limit if input is a Shopee product link.
 * @returns {Promise<object>} The analysis result data.
 */
export const analyzeInput = async (input, topK = 10, limit = 1) => {
  try {
    const isShopee = input.includes('shopee.co.id') || input.includes('shp.ee');
    const payload = {
      input,
      top_k: Number(topK),
    };
    if (isShopee) {
      payload.limit = Number(limit);
    }

    const response = await api.post('/analyze', payload);
    return response.data;
  } catch (error) {
    if (error.response) {
      throw new Error(error.response.data?.message || `Error Server: ${error.response.status}`);
    } else if (error.request) {
      throw new Error('Tidak dapat terhubung ke backend. Pastikan server Golang Anda berjalan di http://localhost:8080.');
    } else {
      throw new Error(`Gagal memproses request: ${error.message}`);
    }
  }
};

export default api;
