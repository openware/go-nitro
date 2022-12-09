// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;
pragma experimental ABIEncoderV2;

import './interfaces/IForceMoveApp.sol';
import './libraries/NitroUtils.sol';
import './interfaces/INitroTypes.sol';
import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';

// NOTE: Attack:
// Bob can submit a convenient candidate, when Alice in trouble (Way back machine attack)

// Possible solutions:
// 1: Alice does checkpoint periodically
// 2: Alice hire a WatchTower, which replicates Alice's states,
// and challenge in the case of challenge event and missing heartbeat

/**
 * @dev The MarginApp contract complies with the ForceMoveApp interface and allows payments to be made virtually from Initiator to Receiver (participants[0] to participants[n+1], where n is the number of intermediaries).
 */
// rename to VirtualMarginApp
contract MarginApp is IForceMoveApp {
    enum AllocationIndices {
        Initiator,
        Receiver
    }

    /**
     * @notice Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @dev Encodes application-specific rules for a particular ForceMove-compliant state channel.
     * @param fixedPart Fixed part of the state channel.
     * @param proof Array of recovered variable parts which constitutes a support proof for the candidate.
     * @param candidate Recovered variable part the proof was supplied for.
     */
    function requireStateSupported(
        FixedPart calldata fixedPart,
        RecoveredVariablePart[] calldata proof,
        RecoveredVariablePart calldata candidate
    ) external pure override {
        // This channel has only 4 states which can be supported:
        // 0    prefund
        // 1    postfund
        // 2+   margin change
        // 3+   final

        uint8 nParticipants = uint8(fixedPart.participants.length);

        // states 0,1,3+:

        if (proof.length == 0) {
            require(
                NitroUtils.getClaimedSignersNum(candidate.signedBy) == nParticipants,
                '!unanimous'
            );

            if (candidate.variablePart.turnNum == 0) return; // prefund
            if (candidate.variablePart.turnNum == 1) return; // postfund

            // postfund
            if (candidate.variablePart.turnNum >= 3) {
                // final (note: there is a core protocol escape hatch for this, too, so it could be removed)
                require(candidate.variablePart.isFinal, '!final; turnNum>=3 && |proof|=0');
                return;
            }

            revert('bad candidate turnNum; |proof|=0');
        }

        // state 2+ requires previous supported state to be supplied

        if (proof.length == 1) {
            require(candidate.variablePart.turnNum >= 2, 'turnNum < 2; |proof|=1');

            require(
                NitroUtils.isClaimedSignedBy(candidate.signedBy, 0),
                'redemption not signed by Leader'
            );

            require(
                NitroUtils.isClaimedSignedBy(candidate.signedBy, nParticipants - 1),
                'redemption not signed by Receiver'
            );

            // previous state if postfund
            require(proof[0].variablePart.turnNum == 1, 'proof[0].turnNum != 1; |proof|=1');

            // previous state is unanimously signed postfund
            require(
                NitroUtils.getClaimedSignersNum(proof[0].signedBy) == fixedPart.participants.length,
                '!unanimous proof; |proof|=1'
            );

            // TODO: check only 2 assets with only 1 destination each
            _requireSumHasNotChanged(
                proof[0].variablePart.outcome[0].allocations,
                candidate.variablePart.outcome[0].allocations,
                nParticipants
            );

            return;
        }
        revert('bad proof length');
    }

    function _requireCorrectAssets() internal pure {
        // require(oldOutcome.length == 2 && newOutcome.length == 2, 'invalid number of assets');
        // TODO: Add later getter and setter, for Fee and collateral currencies
        // oldOutcome[0].asset == ASSET_FEE_ADDRESS &&
        // newOutcome[0].asset == ASSET_COLLATERAL_ADDRESS,
    }

    function _requireSumHasNotChanged(
        Outcome.Allocation[] memory oldAllocations,
        Outcome.Allocation[] memory newAllocations,
        uint256 nParticipants
    ) internal pure {
        uint256 oldAllocationSum;
        uint256 newAllocationSum;
        for (uint256 i = 0; i < nParticipants; i++) {
            oldAllocationSum += oldAllocations[i].amount;
            newAllocationSum += newAllocations[i].amount;
        }
        require(oldAllocationSum == newAllocationSum, 'Total allocated cannot change');
    }
}
