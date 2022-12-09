import {ParamType} from 'ethers/lib/utils';
import {Wallet} from 'ethers';
import {ethers} from 'hardhat';

import {YellowAdjudicator} from '../../typechain-types/YellowAdjudicator';
import {bindSignatures, getFixedPart, getVariablePart, signStates, State} from '../../src';
import {FixedPart, SignedVariablePart, VariablePart} from '../../src/contract/state';

export function randomWallet(): Wallet {
  return ethers.Wallet.createRandom();
}

function randomAddress(): string {
  return randomWallet().address;
}

function randomNum(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

interface Liability {
  asset: string;
  amount: number;
}

interface SwapSpecs {
  brokerA: string;
  brokerB: string;
  swaps: Array<[Liability, Liability]>;
  swapSpecsFinalizationTimestamp: number;
}

const MIN_DEPOSIT_AMOUNT = 100;
const MAX_DEPOSIT_AMOUNT = 1000;

const MAX_SWAP_AMOUNT = 100;

export function constructSwapSpecs(
  brokerA: string,
  brokerB: string,
  swapsNum: number,
  finalizationTimestamp: number
): SwapSpecs {
  const specs: SwapSpecs = {
    brokerA,
    brokerB,
    swaps: [],
    swapSpecsFinalizationTimestamp: finalizationTimestamp,
  };

  return specs;
}

function randomDepositAmount(): number {
  return randomNum(MIN_DEPOSIT_AMOUNT, MAX_DEPOSIT_AMOUNT);
}

function randomSwapAmount(): number {
  return randomNum(1, MAX_SWAP_AMOUNT);
}

function encodeSwapSpecs(swapSpecs: SwapSpecs): string {
  return ethers.utils.defaultAbiCoder.encode(
    [
      {
        type: 'tuple',
        components: [
          {name: 'brokerA', type: 'address'},
          {name: 'brokerB', type: 'address'},
          {
            type: 'tuple[2][]',
            name: 'swaps',
            components: [
              {
                name: 'asset',
                type: 'address',
              },
              {
                name: 'amount',
                type: 'uint256',
              },
            ],
          } as ParamType,
          {name: 'swapSpecsFinalizationTimestamp', type: 'uint64'},
        ],
      } as ParamType,
    ],
    [swapSpecs]
  );
}

export async function depositForSwaps(
  yellowAdjudicator: YellowAdjudicator,
  brokerA: Wallet,
  brokerB: Wallet,
  swapNum: number
): Promise<string[]> {
  // create assets
  const assets = Array.from({length: swapNum}, () => randomAddress());

  // set assets to brokers
  await Promise.all(
    assets.map(async asset => {
      const amount1 = randomDepositAmount();
      await yellowAdjudicator.setDeposit(brokerA.address, asset, amount1);

      const amount2 = randomDepositAmount();
      await yellowAdjudicator.setDeposit(brokerB.address, asset, amount2);
    })
  );

  return assets;
}

export async function getSwapParams(
  brokerA: Wallet,
  brokerB: Wallet,
  assets: string[],
  swapNum: number
): Promise<[FixedPart, VariablePart, SignedVariablePart]> {
  // construct swaps
  const swaps = [] as Array<[Liability, Liability]>;

  for (let i = 0; i < swapNum; i += 2) {
    swaps.push([
      {asset: assets[i], amount: randomSwapAmount()},
      {asset: assets[i + 1], amount: randomSwapAmount()},
    ]);
  }

  const swapSpecsFinalizationTimestamp = Math.round(Date.now() / 1000) - 600;

  // construct SwapSpecs
  const swapSpecs: SwapSpecs = {
    brokerA: brokerA.address,
    brokerB: brokerB.address,
    swaps,
    swapSpecsFinalizationTimestamp,
  };

  // construct State
  const state: State = {
    chainId: '0x01',
    participants: [brokerA.address, brokerB.address],
    channelNonce: '0x01',
    appDefinition: randomAddress(),
    challengeDuration: 0xffff,
    turnNum: 1,
    isFinal: false,
    outcome: [],
    appData: encodeSwapSpecs(swapSpecs),
  };

  // construct swap parameters
  const fixedPart = getFixedPart(state);
  const variablePart = getVariablePart(state);

  const sigs = await signStates([state], [brokerA, brokerB], [0, 0]);

  const signedVariableParts = bindSignatures([variablePart], sigs, [0, 0]);

  // NOTE: as there are no checks on preSwap variable part yet, we can provide any variable part
  return [fixedPart, variablePart, signedVariableParts[0]];
}

export async function deployYellowAdjudicator(): Promise<YellowAdjudicator> {
  const YAFactory = await ethers.getContractFactory('YellowAdjudicator');
  const YA = (await YAFactory.deploy()) as YellowAdjudicator;
  await YA.deployed();
  return YA;
}
