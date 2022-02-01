import { ethers } from "hardhat";

async function main() {
  const [deployer] = await ethers.getSigners();
  
  console.log("Deploying contracts with the account:", deployer.address);
  console.log("Account balance:", (await deployer.getBalance()).toString());
  
  const NAFactory = await ethers.getContractFactory("NitroAdjudicator");
  const NitroAdjudicator = await NAFactory.deploy();

  await NitroAdjudicator.deployed();

  console.log("NitroAdjudicator address:", NitroAdjudicator.address);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
