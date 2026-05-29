import React from 'react';
import PredictionCard from './PredictionCard';
import JudgeCard from './JudgeCard';
import ConclusionCard from './ConclusionCard';
import HyDECard from './HyDECard';
import RetrievalResults from './RetrievalResults';

const SingleReviewResult = ({ result }) => {
  if (!result) return null;

  return (
    <div className="space-y-8 animate-fade-in">
      {/* The 3 Core Metric Cards (Prediction, Judge, Synthesis) */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <PredictionCard analysis={result.analysis} />
        <JudgeCard judge={result.judge} />
        <ConclusionCard analysis={result.analysis} judge={result.judge} />
      </div>

      {/* HyDE Hipotetik Card */}
      <HyDECard hydeDocument={result.analysis?.hyde_document} />

      {/* Retrieval Documents Pembanding list */}
      <RetrievalResults results={result.retrieval_results} />
    </div>
  );
};

export default SingleReviewResult;
