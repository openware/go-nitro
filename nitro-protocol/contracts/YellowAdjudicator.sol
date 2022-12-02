// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;
pragma experimental ABIEncoderV2;

import './NitroAdjudicator.sol';
import {NitroUtils} from './libraries/NitroUtils.sol';

contract YellowAdjudicator is NitroAdjudicator {
    struct BrokerLiability {
        address broker;
        address user;
        address asset;
        uint256 amount;
    }

    // TODO: rename
    struct SwapSpecs {
        BrokerLiability liabilityOfBrokerA;
        BrokerLiability liabilityOfBrokerB;
        uint64 swapSpecsFinalizationTimestamp; // guarantees swaps with equal amounts between same brokers are distinguishable
    }

    // broker => user => asset => balance
    mapping(address => mapping(address => mapping(address => uint256))) public deposits;

    // keep track of performed swaps to prevent using the same signatures twice
    // REVIEW: hashedPostSwapSpecs => swapWasPerformed
    mapping(bytes32 => bool) internal _swapPerformed;

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
        address brokerA = postSwapSpecs.liabilityOfBrokerA.broker;
        address brokerB = postSwapSpecs.liabilityOfBrokerB.broker;

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
        address userA = postSwapSpecs.liabilityOfBrokerA.user;
        address userB = postSwapSpecs.liabilityOfBrokerB.user;

        address assetA = postSwapSpecs.liabilityOfBrokerA.asset;
        address assetB = postSwapSpecs.liabilityOfBrokerB.asset;

        uint256 amountA = postSwapSpecs.liabilityOfBrokerA.amount;
        uint256 amountB = postSwapSpecs.liabilityOfBrokerB.amount;

        // alice - 55 WETH
        // alice - 5 WBTC

        // bob - 10 WETH
        // bob - 77 WBTC

        // swap: alice 2 WBTC
        //       bob 15 WETH

        // alice - 40 WETH
        // alice - 7 WBTC

        // bob - 25 WETH
        // bob - 75 WBTC

        deposits[brokerA][userA][assetA] -= amountB;
        deposits[brokerA][userA][assetB] += amountA;

        deposits[brokerB][userB][assetA] += amountB;
        deposits[brokerB][userB][assetB] -= amountA;
    }
}
