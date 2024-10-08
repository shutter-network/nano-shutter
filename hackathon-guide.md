# Hackathon Guide: Building a Decentralized Shutterized Rock Paper Scissors DApp

## Table of Contents

1. [Introduction](#introduction)
2. [Project Overview](#project-overview)
3. [Technology Stack](#technology-stack)
4. [Prerequisites](#prerequisites)
5. [Step-by-Step Guide](#step-by-step-guide)
   - [1. Setting Up the Development Environment](#1-setting-up-the-development-environment)
   - [2. Designing the Smart Contract](#2-designing-the-smart-contract)
   - [3. Implementing the Smart Contract](#3-implementing-the-smart-contract)
   - [4. Setting Up the Frontend](#4-setting-up-the-frontend)
   - [5. Integrating Web3 and Smart Contract](#5-integrating-web3-and-smart-contract)
   - [6. Implementing the Commit-Reveal Scheme](#6-implementing-the-commit-reveal-scheme)
   - [7. Testing the DApp](#7-testing-the-dapp)
6. [Enhancements and Additional Features](#enhancements-and-additional-features)
7. [Conclusion](#conclusion)
8. [Resources](#resources)

---

## Introduction

Welcome to this hackathon guide! In this tutorial, we'll walk you through building a decentralized application (DApp) version of the Shutterized Rock Paper Scissors game. This guide is designed to help you create an example DApp that demonstrates the use of cryptographic commit-reveal schemes on the blockchain, ensuring fairness and transparency in the game.

By the end of this guide, you'll have a working DApp where two players can securely play Rock Paper Scissors on a decentralized network.

---

## Project Overview

**Shutterized Rock Paper Scissors DApp** leverages blockchain technology to create a fair and transparent gaming experience. The game uses a commit-reveal scheme, where players commit to their moves in a hashed form and reveal them later. This ensures that neither player can cheat by changing their move after seeing the opponent's move.

**Key Features:**

- Decentralized gameplay using smart contracts.
- Secure commit-reveal mechanism to prevent cheating.
- Transparent and immutable game records on the blockchain.
- Frontend interface for players to interact with the game.

---

## Technology Stack

- **Blockchain Platform:** Ethereum (can be adapted to other EVM-compatible chains)
- **Smart Contract Language:** Solidity
- **Frontend Framework:** React.js (optional, you can use plain HTML/CSS/JS)
- **Web3 Integration:** Ethers.js or Web3.js
- **Development Environment:** Truffle or Hardhat for smart contract development
- **Cryptographic Functions:** Keccak256 hashing (built into Solidity)
- **Tools:** MetaMask for browser wallet integration

---

## Prerequisites

- **Basic Understanding of:**
  - Blockchain and smart contracts
  - Solidity programming
  - Web development (HTML, CSS, JavaScript)
- **Installed Software:**
  - Node.js and npm
  - Truffle or Hardhat
  - MetaMask browser extension
- **Accounts:**
  - Access to a test Ethereum network (e.g., Ganache, Ropsten, Rinkeby)

---

## Step-by-Step Guide

### 1. Setting Up the Development Environment

#### a. Install Node.js and npm

Download and install Node.js from the [official website](https://nodejs.org/).

#### b. Install Truffle or Hardhat

Choose one of the following frameworks for smart contract development:

- **Truffle:**

  ```bash
  npm install -g truffle
  ```

- **Hardhat:**

  ```bash
  npm install --save-dev hardhat
  ```

#### c. Install Ganache (Optional)

For local blockchain development, install Ganache:

```bash
npm install -g ganache-cli
```

#### d. Initialize the Project

Create a new project directory:

```bash
mkdir rock-paper-scissors-dapp
cd rock-paper-scissors-dapp
```

Initialize a new npm project:

```bash
npm init -y
```

### 2. Designing the Smart Contract

Our smart contract will handle:

- Player registration and move commitment.
- Move reveal and validation.
- Determining the winner and handling payouts.

#### Key Considerations:

- **Commit-Reveal Scheme:** Players commit to a hashed version of their move, combining their move with a secret nonce.
- **Security:** Prevent players from revealing their move before both have committed.
- **Fairness:** Ensure the game can be resolved even if a player fails to reveal their move.

### 3. Implementing the Smart Contract

Create a `contracts` directory and add a Solidity file:

```bash
mkdir contracts
cd contracts
touch RockPaperScissors.sol
```

#### **RockPaperScissors.sol**

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RockPaperScissors {
    enum Move { None, Rock, Paper, Scissors }
    enum GameState { WaitingForPlayers, MovesCommitted, MovesRevealed, GameOver }

    struct Player {
        address addr;
        bytes32 moveHash;
        Move move;
    }

    uint public betAmount;
    GameState public gameState;
    Player public player1;
    Player public player2;
    address public winner;

    constructor(uint _betAmount) {
        betAmount = _betAmount;
        gameState = GameState.WaitingForPlayers;
    }

    modifier onlyInState(GameState _state) {
        require(gameState == _state, "Invalid game state.");
        _;
    }

    modifier onlyPlayer() {
        require(msg.sender == player1.addr || msg.sender == player2.addr, "Not a player.");
        _;
    }

    function joinGame(bytes32 _moveHash) external payable onlyInState(GameState.WaitingForPlayers) {
        require(msg.value == betAmount, "Incorrect bet amount.");

        if (player1.addr == address(0)) {
            player1 = Player(msg.sender, _moveHash, Move.None);
        } else if (player2.addr == address(0)) {
            require(msg.sender != player1.addr, "Player already joined.");
            player2 = Player(msg.sender, _moveHash, Move.None);
            gameState = GameState.MovesCommitted;
        } else {
            revert("Game is full.");
        }
    }

    function revealMove(Move _move, string calldata _nonce) external onlyInState(GameState.MovesCommitted) onlyPlayer {
        bytes32 computedHash = keccak256(abi.encodePacked(_move, _nonce));

        if (msg.sender == player1.addr) {
            require(player1.moveHash == computedHash, "Invalid move or nonce.");
            player1.move = _move;
        } else if (msg.sender == player2.addr) {
            require(player2.moveHash == computedHash, "Invalid move or nonce.");
            player2.move = _move;
        }

        if (player1.move != Move.None && player2.move != Move.None) {
            determineWinner();
        }
    }

    function determineWinner() internal {
        // Logic to determine the winner
        if (player1.move == player2.move) {
            // It's a tie, refund bets
            payable(player1.addr).transfer(betAmount);
            payable(player2.addr).transfer(betAmount);
        } else if (
            (player1.move == Move.Rock && player2.move == Move.Scissors) ||
            (player1.move == Move.Paper && player2.move == Move.Rock) ||
            (player1.move == Move.Scissors && player2.move == Move.Paper)
        ) {
            // Player 1 wins
            winner = player1.addr;
            payable(winner).transfer(address(this).balance);
        } else {
            // Player 2 wins
            winner = player2.addr;
            payable(winner).transfer(address(this).balance);
        }

        gameState = GameState.GameOver;
    }

    function timeout() external onlyInState(GameState.MovesCommitted) onlyPlayer {
        // If one player fails to reveal within a certain time, allow the other to claim the pot
        // Implement timeout logic based on block.timestamp or block.number
    }

    // Additional helper functions as needed
}
```

**Explanation:**

- **Game States:** We define different states of the game to control the flow.
- **Player Structure:** Holds the player's address, committed move hash, and revealed move.
- **Join Game:** Players join by submitting the hash of their move and nonce.
- **Reveal Move:** Players reveal their move and nonce; the contract verifies the hash.
- **Determine Winner:** Once both moves are revealed, the contract determines the winner and transfers the pot.
- **Timeout Function:** Allows the game to proceed if a player fails to reveal (to be implemented).

### 4. Setting Up the Frontend

We'll create a simple frontend to interact with the smart contract.

#### a. Initialize React App (Optional)

```bash
npx create-react-app frontend
cd frontend
```

#### b. Install Dependencies

```bash
npm install ethers
npm install @openzeppelin/contracts
```

#### c. Set Up MetaMask

- Install the MetaMask extension in your browser.
- Create or import an account.
- Connect MetaMask to your test network (e.g., Ganache).

### 5. Integrating Web3 and Smart Contract

#### a. Compile and Deploy the Smart Contract

- For Truffle:

  ```bash
  truffle compile
  truffle migrate --network development
  ```

- For Hardhat:

  Create a deployment script and run:

  ```bash
  npx hardhat run scripts/deploy.js --network localhost
  ```

#### b. Obtain Contract ABI and Address

- After deployment, retrieve the contract's ABI and address.
- Copy the ABI file to your frontend project.

#### c. Connecting to the Contract in Frontend

```javascript
import { ethers } from "ethers";
import RockPaperScissors from "./RockPaperScissors.json"; // ABI file

const provider = new ethers.providers.Web3Provider(window.ethereum);
const signer = provider.getSigner();
const contractAddress = "YOUR_CONTRACT_ADDRESS";
const contract = new ethers.Contract(contractAddress, RockPaperScissors.abi, signer);
```

### 6. Implementing the Commit-Reveal Scheme

#### a. Player Move Commitment

- **Hashing the Move:**

  ```javascript
  const move = 1; // 1 for Rock, 2 for Paper, 3 for Scissors
  const nonce = "random_string";
  const moveHash = ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(["uint8", "string"], [move, nonce]));
  ```

- **Submitting the Hash to the Contract:**

  ```javascript
  await contract.joinGame(moveHash, { value: betAmount });
  ```

#### b. Player Move Reveal

- **Calling the Reveal Function:**

  ```javascript
  await contract.revealMove(move, nonce);
  ```

#### c. Frontend UI

- **Move Selection:**
  - Provide buttons or dropdown for players to select their move.
- **Nonce Generation:**
  - Automatically generate a random nonce or allow the player to input one.
- **Game Status Display:**
  - Show the current state of the game and any messages.
- **Winner Announcement:**
  - Display the winner once the game concludes.

### 7. Testing the DApp

#### a. Running a Local Blockchain

Start Ganache CLI or GUI to simulate a local blockchain:

```bash
ganache-cli
```

#### b. Interacting with the DApp

- Open your frontend application in the browser.
- Connect MetaMask to the local blockchain.
- Use two browser profiles or two different browsers to simulate two players.
- Follow the game flow:
  - Player 1 commits their move.
  - Player 2 commits their move.
  - Both players reveal their moves.
  - The winner is determined, and the pot is transferred.

#### c. Handling Edge Cases

- Test scenarios where a player does not reveal their move.
- Implement and test the timeout functionality to handle such cases.

---

## Enhancements and Additional Features

- **Timeout Mechanism:**
  - Implement a timeout period after which a player can forfeit if they fail to reveal.
- **Multiple Games:**
  - Modify the contract to handle multiple games simultaneously.
- **Betting System:**
  - Allow variable bet amounts and implement escrow.
- **User Interface Improvements:**
  - Enhance the frontend with better styling and user experience.
- **Security Audits:**
  - Review the contract for security vulnerabilities and optimize gas usage.
- **Integration with IPFS:**
  - Store game data on IPFS for decentralized storage.
- **Deployment to Testnet/Mainnet:**
  - Deploy the contract to a public testnet like Ropsten or Rinkeby.

---

## Conclusion

Congratulations! You've built a decentralized Rock Paper Scissors game using Ethereum smart contracts and a commit-reveal scheme. This DApp showcases how blockchain technology can be used to create fair and transparent applications.

This project serves as a foundation for further exploration into decentralized gaming and cryptographic protocols. Feel free to expand upon this guide by adding new features and enhancements.

---

## Resources

- **Ethereum Documentation:** [https://ethereum.org/en/developers/docs/](https://ethereum.org/en/developers/docs/)
- **Solidity Documentation:** [https://docs.soliditylang.org/](https://docs.soliditylang.org/)
- **Ethers.js Documentation:** [https://docs.ethers.io/](https://docs.ethers.io/)
- **Truffle Suite:** [https://www.trufflesuite.com/](https://www.trufflesuite.com/)
- **Hardhat:** [https://hardhat.org/](https://hardhat.org/)
- **MetaMask:** [https://metamask.io/](https://metamask.io/)
- **OpenZeppelin Contracts:** [https://docs.openzeppelin.com/contracts/](https://docs.openzeppelin.com/contracts/)
- **Ganache:** [https://www.trufflesuite.com/ganache](https://www.trufflesuite.com/ganache)

---

**Note:** Always exercise caution when deploying smart contracts to the mainnet. Thoroughly test your contracts and consider security audits for any production-level deployment.