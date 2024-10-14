# Hackathon Guide: Build a Shutterized Rock Paper Scissors Game

## Introduction

Welcome to this hackathon guide! This guide will walk you through building a Rock Paper Scissors game utilizing NanoShutter's encryption/decryption API. The game leverages Shutter-like threshold encryption to ensure fairness, providing an engaging way to experience blockchain cryptography.

![ezgif-6-1044641a13](https://github.com/user-attachments/assets/9b3e82d5-dd4f-4190-aff1-8a102b8bf999)

### What You'll Learn:
- How to use the NanoShutter API for time-based encryption and decryption.
- Building a web application with JavaScript and HTML to implement a simple two-player game.
- Understanding the use of encryption in a decentralized, transparent manner.

## Prerequisites

- Basic knowledge of HTML, CSS, and JavaScript.
- Familiarity with REST APIs and making HTTP requests.
- Node.js and npm installed (for local development if needed).
- Basic understanding of cryptography concepts is helpful but not necessary.

## Getting Started

### 1. Use Our Provided API

NanoShutter provides hosted API endpoints that you can use directly. These include:

- **Encryption Endpoint** (with time-based access):
  - `https://nanoshutter.staging.shutter.network/encrypt/with_time`
- **Decryption Endpoint** (with time-based access):
  - `https://nanoshutter.staging.shutter.network/decrypt/with_time`

You can use these endpoints for testing and experimentation without needing to host your own server.

### 2. Set Up Your Project

- Create a new folder and open it in your code editor.
- Inside this folder, create a new HTML file called `index.html`. Copy and paste the provided HTML code for the Rock Paper Scissors game into this file.

**File Structure:**
```
rock-paper-scissors/
  ├── index.html
```

### 3. HTML Overview

The `index.html` file contains the entire UI for the game, including:
- Two player sections where players can select their moves.
- A button for each player to submit their move.
- Sections to display encrypted moves, decrypted moves, and the game result.
- A "Start New Game" button to reset the game.

### 4. Main Features of the Game
- **Player Interaction**: Each player can select their move (Rock, Paper, Scissors) and submit it.
- **Encryption/Decryption**: Each move is encrypted with a timestamp set to 20 seconds into the future. Decryption occurs after 20 seconds have elapsed.
- **Game Logic**: After both moves are decrypted, the game evaluates the winner and displays the result.

## Step-by-Step Guide to Running the Game

### Step 1: Integrate the NanoShutter API

In the provided JavaScript section of the HTML code, we make use of `axios` to send HTTP requests to the NanoShutter API for both encryption and decryption.

- **Encryption Example**:
  ```javascript
  const response = await axios.post(`${apiBaseUrl}/encrypt/with_time`, {
    cypher_text: move,
    timestamp: encryptionTimestamp
  });
  ```
  This encrypts the selected move using the given `encryptionTimestamp` (20 seconds from the current time).

- **Decryption Example**:
  ```javascript
  const response = await axios.post(`${apiBaseUrl}/decrypt/with_time`, {
    encrypted_msg: encryptedData.player1,
    timestamp: encryptionTimestamp
  });
  ```
  This decrypts the move after 20 seconds have elapsed.

### Step 2: Running the Game

1. **Open the `index.html` file in a browser**.
2. **Player 1 and Player 2**: Each player chooses their move and clicks the "Submit Move" button.
3. **Wait for the moves to be decrypted**: The game waits until the decryption timestamp to reveal the players' moves.
4. **View the Result**: Once both moves are decrypted, the game result (which player won) will be displayed.

### Step 3: Starting a New Game
- Click the "Start New Game" button to reset the game.

## Detailed Guide for Operating Your Own API

If you want to set up and operate your own version of the NanoShutter API instead of using the provided endpoints, follow these steps:

### 1. Clone the Repository
First, clone the NanoShutter repository from GitHub:
```sh
git clone https://github.com/shutter-network/nanoshutter.git
```

### 2. Install Dependencies
Navigate to the cloned folder
```sh
cd nanoshutter/rock-paper-scissors-demo
```

### 3. Run the API Server, or use the prdeployed one.

### 4. Update the HTML File to Use Your Local API
If you are running your own version of the API, change the `apiBaseUrl` in the HTML code to point to your local server:
```javascript
const apiBaseUrl = 'http://localhost:5000';
```
This will make sure all encryption and decryption requests are directed to your locally hosted API.

## Other Hackathon Ideas

NanoShutter is an ideal tool for experimenting with time-locked encryption in various hackathon scenarios. Below are some ideas to inspire you at your next hackathon:

### 1. **Shutterized Multiplayer Games**
   - Build games like Rock Paper Scissors, where players commit their moves and decrypt them simultaneously using NanoShutter, ensuring fairness.

### 2. **Sealed-Bid Auctions**
   - Create a decentralized auction where bids are encrypted and revealed only after a specific time, using the NanoShutter API for transparent and fair bidding.

### 3. **Front-Running Protection for DApps**
   - Integrate NanoShutter into a decentralized exchange (DEX) to enable front-running protection by encrypting trades before they are submitted for execution.

### 4. **Salary Negotiation Tool**
   - Develop a privacy-focused negotiation tool that allows individuals to submit encrypted salary offers that are revealed after both parties have committed, ensuring transparency and fairness.

### 5. **Blind Voting Application**
   - Build a simple voting application where votes are encrypted until the voting period ends, enabling privacy-preserving decision-making in decentralized environments.

### 6. **Token-Gated Content Release**
   - Use NanoShutter to create a platform for releasing premium content (e.g., music, art) only after a certain time or when conditions are met, enabling timed exclusive access.

### 7. **Blind Fundraising Campaigns**
   - Create a platform for charitable donations where donation amounts are encrypted and revealed only after the campaign ends, preventing strategic donations and promoting fairness.

### 8. **Encryption-Backed Trivia Game**
   - Design a trivia game where questions and answers are encrypted, and participants reveal their answers simultaneously after the question is closed, maintaining fairness and fun.

These hackathon ideas are just starting points. Feel free to modify, combine, or expand on them to create unique and impactful projects using NanoShutter!


## Tips for Hackathon Success

- **Understand Timing**: Ensure your timestamps for encryption and decryption match up to avoid decryption failures.
- **Debugging**: Use the browser console to inspect errors during encryption/decryption.
- **Collaboration**: Use version control (e.g., GitHub) to collaborate with teammates effectively.
- **Extend the Game**: You can also consider extending the game to add new features, such as player authentication or storing game history on-chain.

## Conclusion

By following this guide, you should be able to create a fully functioning "Shutterized Rock Paper Scissors" game that uses cryptographic primitives to ensure fair play. This is a great introduction to threshold cryptography and how it can be used in decentralized applications.

We can't wait to see what you build! If you have any questions, feel free to reach out to the NanoShutter community for support.

Happy Hacking!
