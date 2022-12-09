export const BENCHMARK_STEPS = [1, 10, 100, 200, 300, 400, 500, 650, 800, 1000] as const;

interface EmptyYellowAdjudicatorGasResults {
  swap: {
    [key in typeof BENCHMARK_STEPS[number]]?: number;
  };
}

export const emptyYellowAdjudicatorGasResults: EmptyYellowAdjudicatorGasResults = {
  swap: {},
};

// BENCHMARK_STEPS.forEach(step => (emptyYellowAdjudicatorGasResults.swap[step] = 0));
