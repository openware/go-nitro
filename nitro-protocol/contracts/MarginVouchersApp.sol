// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;
pragma experimental ABIEncoderV2;

import './interfaces/IForceMoveApp.sol';
import './libraries/NitroUtils.sol';
import './interfaces/INitroTypes.sol';
import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';

/**
 * @dev The MarginVouchersApp contract complies with the ForceMoveApp interface and allows payments to be made virtually from Initiator to Receiver (participants[0] to participants[n+1], where n is the number of intermediaries).
 */
contract MarginVouchersApp is IForceMoveApp {
    struct NewMargin {
        uint256 initiatorMargin;
        uint256 receiverMargin;
    }

    struct MarginVoucher {
        uint256 initiatorMargin;
        uint256 receiverMargin;
        INitroTypes.Signature initiatorSignature; // initiator signature on abi.encode(channelId,amount)
        INitroTypes.Signature receiverSignature; // receiver signature on abi.encode(channelId,amount)
        int256 nonce; // to distinct between valid vouchers
    }

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
        // 2+   margin voucher
        // 3+   final

        uint256 nParticipants = fixedPart.participants.length;

        // all states require unanimous consensus
        require(NitroUtils.getClaimedSignersNum(candidate.signedBy) == nParticipants, '!unanimous');

        // states 0,1,3+:

        if (proof.length == 0) {
            // TODO: do we need a check allocations and destination has not changed? (as in SingleAssetPayments)
            if (candidate.variablePart.turnNum == 0) return; // prefund
            if (candidate.variablePart.turnNum == 1) return; // postfund

            // postfund
            // TODO: can we safely remove this check? Final turn number is NOT FIXED (3+). Any assumptions CAN NOT be based on it.
            if (candidate.variablePart.turnNum >= 3) {
                // final (note: there is a core protocol escape hatch for this, too, so it could be removed)
                require(candidate.variablePart.isFinal, '!final; turnNum>=3 && |proof|=0');
                return;
            }

            revert('bad candidate turnNum; |proof|=0');
        }

        // state 2+ requires previous supported state to be supplied

        if (proof.length == 1) {
            // previous state is unanimously signed
            require(
                NitroUtils.getClaimedSignersNum(proof[0].signedBy) == fixedPart.participants.length,
                '!unanimous proof; |proof|=1'
            );

            // previous this state has bigger turn number
            require(
                candidate.variablePart.turnNum > proof[0].variablePart.turnNum,
                'candidate turnNum not increased'
            );

            NewMargin memory newMargin = _requireValidVoucher(
                candidate.variablePart.appData,
                fixedPart
            );

            _requireCorrectAdjustments(
                proof[0].variablePart.outcome,
                candidate.variablePart.outcome,
                nParticipants,
                newMargin
            );
            return;
        }
        revert('bad proof length');
    }

    function _requireValidVoucher(
        bytes memory appData,
        FixedPart memory fixedPart
    ) internal pure returns (NewMargin memory) {
        MarginVoucher memory voucher = abi.decode(appData, (MarginVoucher));

        NewMargin memory newMargin = NewMargin(voucher.initiatorMargin, voucher.receiverMargin);

        // validate initiator signature
        address initiatorSigner = NitroUtils.recoverSigner(
            keccak256(abi.encode(NitroUtils.getChannelId(fixedPart), newMargin)),
            voucher.initiatorSignature
        );
        require(initiatorSigner == fixedPart.participants[0], 'invalid signature for voucher'); // could be incorrect channelId or incorrect signature

        // validate receiver signature
        address receiverSigner = NitroUtils.recoverSigner(
            keccak256(abi.encode(NitroUtils.getChannelId(fixedPart), newMargin)),
            voucher.receiverSignature
        );
        require(receiverSigner == fixedPart.participants[0], 'invalid signature for voucher'); // could be incorrect channelId or incorrect signature

        return newMargin;
    }

    function _requireCorrectAdjustments(
        Outcome.SingleAssetExit[] memory oldOutcome,
        Outcome.SingleAssetExit[] memory newOutcome,
        uint256 nParticipants,
        NewMargin memory newMargin
    ) internal pure {
        require(oldOutcome.length == 1 && newOutcome.length == 1, 'only one asset allowed');

        // check the sum has not changed
        _requireSumHasNotChanged(
            oldOutcome[0].allocations,
            newOutcome[0].allocations,
            nParticipants
        );

        // check new outcome is set respecting the MarginVoucher
        require(
            newOutcome[0].allocations[uint256(AllocationIndices.Initiator)].amount ==
                newMargin.initiatorMargin,
            'Initiator not adjusted correctly'
        );
        require(
            newOutcome[0].allocations[uint256(AllocationIndices.Receiver)].amount ==
                newMargin.receiverMargin,
            'Receiver not adjusted correctly'
        );
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
