// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;
pragma experimental ABIEncoderV2;

import './NitroAdjudicator.sol';
import {NitroUtils} from './libraries/NitroUtils.sol';

contract YellowAdjudicator is NitroAdjudicator {
    struct Liability {
        address asset;
        uint256 amount;
    }

    struct SwapSpecs {
        address brokerA;
        address brokerB;
        Liability[2][] swaps;
        uint64 swapSpecsFinalizationTimestamp; // guarantees swaps with equal amounts between same brokers are distinguishable
    }

    // broker => asset => balance
    mapping(address => mapping(address => uint256)) public deposits;

    // keep track of performed swaps to prevent using the same signatures twice
    // REVIEW: hashedPostSwapSpecs => swapWasPerformed
    mapping(bytes32 => bool) internal _swapPerformed;

    function test(SwapSpecs memory swapSpecs) external pure {}

    function swap(
        FixedPart memory fixedPart,
        VariablePart memory preSwapVP,
        SignedVariablePart memory postSwapSVP
    ) external {
        SwapSpecs memory postSwapSpecs = abi.decode(postSwapSVP.variablePart.appData, (SwapSpecs));

        // check this swap has not been performed
        require(
            _swapPerformed[keccak256(abi.encode(postSwapSpecs))] == false,
            'swap already performed'
        );

        // check finalizationTimestamp is < now and != 0
        require(postSwapSpecs.swapSpecsFinalizationTimestamp != 0, 'swap specs not finalized yet');
        require(
            postSwapSpecs.swapSpecsFinalizationTimestamp < block.timestamp,
            'future swap specs finalized'
        );

        // REVIEW: what check on outcome (margin) should be performed (check outcome sums equal)
        // REVIEW: should we check if guarantee allocations (margin) exist?

        // check sigs on postSwapSpecs
        bytes32 postSwapSpecsHash = NitroUtils.hashState(fixedPart, postSwapSVP.variablePart);
        address brokerA = postSwapSpecs.brokerA;
        address brokerB = postSwapSpecs.brokerB;

        require(
            NitroUtils.isSignedBy(postSwapSpecsHash, postSwapSVP.sigs[0], brokerA),
            'not signed by brokerA'
        );
        require(
            NitroUtils.isSignedBy(postSwapSpecsHash, postSwapSVP.sigs[1], brokerB),
            'not signed by brokerB'
        );

        // mark swap as performed
        _swapPerformed[keccak256(abi.encode(postSwapSpecs))] = true;

        // perform swap
        // REVIEW: how to improve readability?

        for (uint256 i = 0; i < postSwapSpecs.swaps.length; i++) {
            address assetA = postSwapSpecs.swaps[i][0].asset;
            address assetB = postSwapSpecs.swaps[i][1].asset;

            uint256 amountA = postSwapSpecs.swaps[i][0].amount;
            uint256 amountB = postSwapSpecs.swaps[i][1].amount;

            // broker1 - 55 WETH
            // broker1 - 5 WBTC

            // broker2 - 10 WETH
            // broker2 - 77 WBTC

            // swap: broker1 2 WBTC
            //       broker2 15 WETH

            // broker1 - 40 WETH
            // broker1 - 7 WBTC

            // broker2 - 25 WETH
            // broker2 - 75 WBTC

            deposits[brokerA][assetA] -= amountB;
            deposits[brokerA][assetB] += amountA;

            deposits[brokerB][assetA] += amountB;
            deposits[brokerB][assetB] -= amountA;
        }
    }

    function setDeposit(address broker, address asset, uint256 amount) external {
        deposits[broker][asset] = amount;
    }
}
