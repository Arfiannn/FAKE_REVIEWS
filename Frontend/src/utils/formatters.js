/**
 * Formats a similarity score (typically 0-1 float) to a percentage string.
 * @param {number} value - The similarity score.
 * @returns {string} Formatted percentage.
 */
export const formatSimilarity = (value) => {
  if (value === undefined || value === null) return '0%';
  // If the value is already in 0-100 range (e.g. > 1)
  const percent = value <= 1 ? value * 100 : value;
  return `${Math.round(percent)}%`;
};

/**
 * Formats a confidence score (typically 0-100 or 0-1 float) to a percentage string.
 * @param {number} value - The confidence score.
 * @returns {string} Formatted percentage.
 */
export const formatConfidence = (value) => {
  if (value === undefined || value === null) return '0%';
  // If the score is in 0-1 range, convert to 0-100
  const percent = value <= 1 ? value * 100 : value;
  return `${Math.round(percent)}%`;
};

/**
 * Formats an ISO date string to a human-readable format.
 * @param {string} dateString - The ISO date string.
 * @returns {string} Formatted date.
 */
export const formatDate = (dateString) => {
  if (!dateString) return '-';
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return dateString;
    return date.toLocaleDateString('id-ID', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  } catch (error) {
    return dateString;
  }
};

/**
 * Generates an array of stars for rating display.
 * @param {number} rating - Product rating (1-5).
 * @returns {string} Star string.
 */
export const formatStars = (rating) => {
  const rounded = Math.round(rating || 0);
  return '⭐'.repeat(Math.max(0, Math.min(5, rounded)));
};
