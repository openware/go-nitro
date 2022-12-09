import {writeFileSync} from 'fs';

import {gasUsed} from '../fixtures';

import {deployYellowAdjudicator, depositForSwaps, getSwapParams, randomWallet} from './fixtures';
import {BENCHMARK_STEPS, emptyYellowAdjudicatorGasResults} from './yellowAdjudicatorGas';

async function main() {
  const gasResults = emptyYellowAdjudicatorGasResults;

  await Promise.all(
    BENCHMARK_STEPS.map(async stepNum => {
      const yellowAdjudicator = await deployYellowAdjudicator();
      const brokerA = randomWallet();
      const brokerB = randomWallet();

      const assets = await depositForSwaps(yellowAdjudicator, brokerA, brokerB, stepNum * 2);

      gasResults.swap[`swaps_${stepNum}`] = await gasUsed(
        await yellowAdjudicator.swap(
          ...(await getSwapParams(brokerA, brokerB, assets, stepNum * 2)),
          {gasLimit: 30_000_000}
        )
      );
    })
  );

  writeFileSync(__dirname + '/gasResults.json', JSON.stringify(gasResults, null, 2));
  console.log('Benchmark results updated successfully!');
}

main().catch(error => {
  console.error(error);
  process.exitCode = 1;
});
