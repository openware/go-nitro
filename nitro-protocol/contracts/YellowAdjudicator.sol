// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;
pragma experimental ABIEncoderV2;

import './NitroAdjudicator.sol';

contract YellowAdjudicator is NitroAdjudicator {
    struct ParticipantLiability {
        address participant;
        address token;
        uint256 amount;
    }

    // participant => asset => amount
    mapping(address => mapping(address => uint256)) public liabilities;

    function swap(
        ParticipantLiability memory a,
        ParticipantLiability memory b,
        bytes32 marginSig
    ) external {
        // TODO:
        // check sigs on marginSigs
        // store a fact that swap was performed (TODO: how to represent it?)
        // swap liabilities
    }
}
