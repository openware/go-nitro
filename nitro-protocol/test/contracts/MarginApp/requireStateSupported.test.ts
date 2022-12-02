import {expectRevert} from '@statechannels/devtools';
import {Contract, ethers, BigNumber} from 'ethers';

import MarginAppArtifact from '../../../artifacts/contracts/MarginApp.sol/MarginApp.json';
import {
  convertAddressToBytes32,
  encodeVoucherAmountAndSignature,
  getChannelId,
  signVoucher,
  Voucher,
} from '../../../src';
import {
  getFixedPart,
  getVariablePart,
  RecoveredVariablePart,
  State,
} from '../../../src/contract/state';
import {
  computeOutcome,
  generateParticipants,
  getTestProvider,
  setupContract,
} from '../../test-helpers';
const {HashZero} = ethers.constants;

let MarginApp: Contract;
const provider = getTestProvider();
const chainId = process.env.CHAIN_NETWORK_ID;

const nParticipants = 3;
const {wallets, participants} = generateParticipants(nParticipants);

const challengeDuration = 0x100;
const MAGIC_ETH_ADDRESS = '0x0000000000000000000000000000000000000000';

const baseState: State = {
  turnNum: 0,
  isFinal: false,
  chainId,
  channelNonce: '0x8',
  participants,
  challengeDuration,
  outcome: [],
  appData: HashZero,
  appDefinition: process.env.MARGIN_VOUCHERS_APP_ADDRESS,
};
const fixedPart = getFixedPart(baseState);
const channelId = getChannelId(fixedPart);

const initiator = convertAddressToBytes32(participants[0]); // NOTE: these destinations do not necessarily need to be related to participant addresses
const receiver = convertAddressToBytes32(participants[2]);

beforeAll(async () => {
  MarginApp = setupContract(provider, MarginAppArtifact, process.env.MARGIN_VOUCHERS_APP_ADDRESS);
});

describe('requireStateSupported (lone candidate route)', () => {
  interface TestCase {
    turnNum: number;
    isFinal: boolean;
    reason?: string;
  }

  const testcases: TestCase[] = [
    {turnNum: 0, isFinal: false, reason: undefined},
    {turnNum: 1, isFinal: false, reason: undefined},
    {turnNum: 2, isFinal: false, reason: 'bad candidate turnNum'},
    {turnNum: 3, isFinal: false, reason: '!final; turnNum=3 && |proof|=0'},
    {turnNum: 3, isFinal: true, reason: undefined},
    {turnNum: 4, isFinal: false, reason: 'bad candidate turnNum'},
  ];

  testcases.map(async tc => {
    it(`${tc.reason ? 'reverts        ' : 'does not revert'} for unaninmous consensus on ${
      tc.isFinal ? 'final' : 'nonfinal'
    } state with turnNum ${tc.turnNum}`, async () => {
      const state: State = {
        ...baseState,
        turnNum: tc.turnNum,
        isFinal: tc.isFinal,
      };

      const variablePart = getVariablePart(state);

      const candidate: RecoveredVariablePart = {
        variablePart,
        signedBy: BigNumber.from(0b111).toHexString(),
      };

      if (tc.reason) {
        await expectRevert(
          () => MarginApp.requireStateSupported(fixedPart, [], candidate),
          tc.reason
        );
      } else {
        await MarginApp.requireStateSupported(fixedPart, [], candidate);
      }
    });
  });
});

describe('requireStateSupported (candidate plus single proof state route)', () => {
  interface TestCase {
    proofTurnNum: number;
    candidateTurnNum: number;
    unanimityOnProof: boolean;
    receiverSignedCandidate: boolean;
    voucherForThisChannel: boolean;
    voucherSignedByInitiator: boolean;
    initiatorAdjustedCorrectly: boolean;
    receiverAdjustedCorrectly: boolean;
    nativeAsset: boolean;
    multipleAssets: boolean;
    initiatorUnderflow: boolean;
    reason?: string;
  }

  const vVR: TestCase = {
    // valid voucher redemption
    proofTurnNum: 1,
    candidateTurnNum: 2,
    unanimityOnProof: true,
    receiverSignedCandidate: true,
    voucherForThisChannel: true,
    voucherSignedByInitiator: true,
    initiatorAdjustedCorrectly: true,
    receiverAdjustedCorrectly: true,
    nativeAsset: true,
    multipleAssets: false,
    initiatorUnderflow: false,
    reason: undefined,
  };
  const testcases: TestCase[] = [
    vVR,
    {...vVR, proofTurnNum: 0, reason: 'bad proof[0].turnNum; |proof|=1'},
    {...vVR, unanimityOnProof: false, reason: 'postfund !unanimous; |proof|=1'},
    {...vVR, receiverSignedCandidate: false, reason: 'redemption not signed by Receiver'},
    {...vVR, voucherSignedByInitiator: false, reason: 'invalid signature for voucher'},
    {...vVR, voucherForThisChannel: false, reason: 'invalid signature for voucher'},
    {...vVR, initiatorAdjustedCorrectly: false, reason: 'Initiator not adjusted correctly'},
    {...vVR, receiverAdjustedCorrectly: false, reason: 'Receiver not adjusted correctly'},
    {...vVR, nativeAsset: false, reason: 'only native asset allowed'},
    {...vVR, multipleAssets: true, reason: 'only native asset allowed'},
    {...vVR, initiatorUnderflow: true, reason: ' '}, // we expect transaction to revert without a reason string
  ];

  testcases.map(async tc => {
    it(`${
      tc.reason ? 'reverts        ' : 'does not revert'
    } for a redemption transition with ${JSON.stringify(tc)}`, async () => {
      const proofState: State = {
        ...baseState,
        turnNum: tc.proofTurnNum,
        isFinal: false,
        outcome: computeOutcome({
          [MAGIC_ETH_ADDRESS]: {[initiator]: 10, [receiver]: 10},
        }),
      };

      // construct voucher with the (in)appropriate channelId
      const amount = tc.initiatorUnderflow
        ? BigNumber.from(999_999_999_999).toHexString() // much larger than Initiator's original balance
        : BigNumber.from(7).toHexString();

      const voucher: Voucher = {
        channelId: tc.voucherForThisChannel
          ? channelId
          : convertAddressToBytes32(MAGIC_ETH_ADDRESS),
        amount,
      };

      // make an (in)valid signature
      const signature = await signVoucher(voucher, wallets[0]);
      if (!tc.voucherSignedByInitiator) signature.s = signature.r; // (conditionally) corrupt the signature

      // embed voucher into candidate state
      const encodedVoucherAmountAndSignature = encodeVoucherAmountAndSignature(amount, signature);
      const candidateState: State = {
        ...proofState,
        outcome: computeOutcome({
          [MAGIC_ETH_ADDRESS]: {
            [initiator]: tc.initiatorAdjustedCorrectly ? 3 : 2,
            [receiver]: tc.receiverAdjustedCorrectly ? 7 : 99,
          },
        }),
        turnNum: tc.candidateTurnNum,
        appData: encodedVoucherAmountAndSignature,
      };

      if (!tc.nativeAsset)
        candidateState.outcome[0].asset = process.env.MARGIN_VOUCHERS_APP_ADDRESS;

      if (tc.multipleAssets) candidateState.outcome.push(candidateState.outcome[0]);

      // Sign the proof state (should be everyone)
      const proof: RecoveredVariablePart[] = [
        {
          variablePart: getVariablePart(proofState),
          signedBy: BigNumber.from(tc.unanimityOnProof ? 0b111 : 0b101).toHexString(),
        },
      ];

      // Sign the candidate state (should be just Receiver)
      const candidate: RecoveredVariablePart = {
        variablePart: getVariablePart(candidateState),
        signedBy: BigNumber.from(tc.receiverSignedCandidate ? 0b100 : 0b000).toHexString(), // 0b100 signed by Receiver only
      };

      if (tc.reason) {
        await expectRevert(
          () => MarginApp.requireStateSupported(fixedPart, proof, candidate),
          tc.reason
        );
      } else {
        await MarginApp.requireStateSupported(fixedPart, proof, candidate);
      }
    });
  });
});

describe('requireStateSupported (longer proof state route)', () => {
  it(`reverts for |support|>1`, async () => {
    const variablePart = getVariablePart(baseState);

    const candidate: RecoveredVariablePart = {
      variablePart,
      signedBy: BigNumber.from(0b111).toHexString(),
    };

    await expectRevert(
      () => MarginApp.requireStateSupported(fixedPart, [candidate, candidate], candidate),
      'bad proof length'
    );
  });
});
