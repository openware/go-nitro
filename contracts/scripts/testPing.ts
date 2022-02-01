import { Interface } from 'ethers/lib/utils';
import { ethers } from 'hardhat';

async function main() {
  const [wallet] = await ethers.getSigners();
  
  console.log("Pinging contracts with the account:", wallet.address);
  console.log("Account balance:", (await wallet.getBalance()).toString());

  let abi = new Interface([
    "function getNumber() public pure returns(uint256)"
  ]);
  let NitroAdjudicatorAddress = '0x724BA5e4507aD7501191eb77A900Ee8b4D7D44b9';
  let NitroAdjudicator = new ethers.Contract(NitroAdjudicatorAddress, abi, wallet);

  console.log("The number is:", (await NitroAdjudicator.getNumber()).toNumber());
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
