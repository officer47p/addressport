# Addressport

## Introduction

**Yo, what the heck is Addressport?**
Addressport is an investigation tool for analyzing transactions on different blockchains. Basically, it maps blockchain transactions into graphs, using addresses as nodes and transactions as links (edges).

**Yo cool idea, but what am I going to use it for?**
That's a great question, which I don't know the exact answer to! 😄

Basically, there was a guy trying to scam me that I wanted to track down. I wanted to find out at least which local exchange he was going to use to exchange his/her crypto so I could contact them and let them know. So, I built this tool.

I hope you never need to use this project, but as grandma always says, it's gonna rain 💦.

## How to Start the Project

1. Clone the repository:

   ```
   git clone https://github.com/officer47p/addressport.git
   ```

2. Navigate to the project directory:

   ```
   cd addressport
   ```

3. Create a `.env` file and add your Etherscan API key(use the `.env.example` file to find all possible environment variables)

4. Start the server (default port is 3000):

   ```
   make start
   ```

5. Open your browser and visit:
   http://localhost:3000

## Features and Components

Addressport currently has one feature: "Address Info Graph"

### Address Info Graph

Access this feature via the following URL:
`/api/v1/investigation/tools/transaction-association/<TARGET-ADDRESS>?depth=<DEPTH-OF-SEARCH>&format=<RESPONSE-FORMAT>`

#### Parameters

- `<TARGET-ADDRESS>`: Any EVM account address

  - Addressport currently searches the Ethereum network, but adding any EVM-based blockchain should be a piece of 🍰

- `<DEPTH-OF-SEARCH>`: The level of transaction search

  - 1: Direct transactions **from** or **to** the target address
  - 2+: Includes transactions from addresses found in previous levels

- `<RESPONSE-FORMAT>`:
  - `graph` (default): Returns results in graph format
  - `nodesandlinks`: Returns results in JSON format for custom plotting

#### Example Scenario

1. Scammer creates a fresh address
2. Uses it for some purposes
3. Sends the amount to another fresh address
4. Sends the amount to an exchange user address
5. Exchange sends the amount to their hot wallet (often labeled on the internet)

Using a search depth of 4 in this scenario would provide valuable insights into the transaction flow.
