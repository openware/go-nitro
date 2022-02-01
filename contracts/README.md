# State channels Smart Contracts

## Download dependencies

1. Run `npm install` in this directory

## Deploy NitroAdjudicator

1. Run `npm run contract:node`
2. Don't close the console. While it is running, you can communicate with the contract deployed.

> NOTE: deployed contract addresses available in `addresses.json` file in such format:
```json
{
  [chainId: string]: [
    {
      chainId: string,
      name: string,
      contracts: {
        <contractName>: {
          address: hex
        }
      }
    }
  ]
}
```
