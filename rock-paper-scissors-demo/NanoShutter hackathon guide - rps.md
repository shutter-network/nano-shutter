# Hackathon Guide: Building a Shutterized Rock Paper Scissors Game

Welcome to this hackathon guide! We'll walk you through building a **Shutterized Rock Paper Scissors** game using the NanoShutter API. NanoShutter provides a simplified version of Shutter's encryption, ideal for building dApps where trust and secrecy are required, such as Rock Paper Scissors.

## What You'll Learn
- How to use NanoShutter's time-based encryption/decryption API.
- Building an interactive web application with two-player input.
- Ensuring fairness in player moves by delaying decryption until both players have submitted their actions.

## Prerequisites
- Basic knowledge of JavaScript, HTML, and CSS.
- Node.js installed (for setting up local servers and proxies).
- Basic understanding of HTTP and REST API concepts.

## Project Setup

### Step 1: Clone the Starter Code
To get started, clone the provided starter code or create a basic HTML file as shown below.
This guide assumes you'll be using a local server to run your app.

Create a new folder called `shutterized-rps` and navigate to it:
```sh
mkdir shutterized-rps
cd shutterized-rps
```

Create an `index.html` file and add the following:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Shutterized Rock Paper Scissors</title>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap.min.css">
</head>
<body>
  <!-- Game Interface -->
  <div class="container">
    <h1>Shutterized Rock Paper Scissors Game</h1>
    <div id="playerInputs"></div>
    <button id="newGame" class="btn btn-primary">New Game</button>
  </div>
  <script src="https://cdn.jsdelivr.net/npm/axios/dist/axios.min.js"></script>
  <script src="script.js"></script>
</body>
</html>
```

### Step 2: Write the JavaScript for Game Logic
Create a `script.js` file. This file will handle encryption and decryption requests using the NanoShutter API. It will also manage the game logic to ensure that players can't cheat by decrypting moves prematurely.

**Create the script.js file:**

```javascript
const apiBaseUrl = 'https://nanoshutter.staging.shutter.network';
let gameInProgress = true;
let encryptedMoves = {};

// Player 1 & 2 Move Submission
const submitMove = async (playerId) => {
  const move = document.getElementById(`move${playerId}`).value;
  const timestamp = Math.floor(Date.now() / 1000);

  try {
    const response = await axios.post(`${apiBaseUrl}/encrypt/with_time`, {
      cypher_text: move,
      timestamp: timestamp
    });
    encryptedMoves[playerId] = response.data.message;
    document.getElementById('status').textContent = `${playerId} submitted.`;
  } catch (error) {
    console.error('Error encrypting:', error);
  }
};

// Start New Game
const resetGame = () => {
  encryptedMoves = {};
  gameInProgress = true;
  document.getElementById('status').textContent = 'Waiting for players...';
};

document.getElementById('newGame').addEventListener('click', resetGame);
```

### Step 3: Set Up Your Local Server
To avoid CORS errors, serve the game through a local server.

1. **Using Python:**

   ```sh
   python -m http.server 8000
   ```

2. **Access the App:**

   Go to `http://localhost:8000` in your web browser.

### Step 4: Encrypt and Decrypt Moves
We use the `/encrypt/with_time` and `/decrypt/with_time` API endpoints to ensure both players' moves remain secret until both moves are submitted and the decryption key is available.

Add the decryption function in `script.js`:

```javascript
const decryptMove = async (playerId) => {
  try {
    const response = await axios.post(`${apiBaseUrl}/decrypt/with_time`, {
      encrypted_msg: encryptedMoves[playerId],
      timestamp: Math.floor(Date.now() / 1000)
    });
    document.getElementById('decryptedMoves').textContent += `${playerId}: ${response.data.message}\n`;
  } catch (error) {
    console.error('Error decrypting:', error);
  }
};
```

### Step 5: Determine the Winner
Once moves are decrypted, evaluate the game results.

```javascript
const determineWinner = (move1, move2) => {
  if (move1 === move2) return "It's a tie!";
  if ((move1 === 'Rock' && move2 === 'Scissors') ||
      (move1 === 'Paper' && move2 === 'Rock') ||
      (move1 === 'Scissors' && move2 === 'Paper')) {
    return 'Player 1 wins!';
  }
  return 'Player 2 wins!';
};
```

### Step 6: Complete and Play!
Once you've integrated encryption, decryption, and winner logic, complete the HTML structure to handle both players. Run it on your local server, and start a new game using the `New Game` button whenever you want.

### Optional Enhancements
- **UI Improvement**: Add better visuals using libraries like Bootstrap or Material UI.
- **Host Online**: Deploy to a service like Netlify, and ensure that the backend allows CORS requests from the web app domain.
- **Multiplayer**: Extend the game to allow multiple matches simultaneously by using session management.

## Summary
In this guide, you've built a secure and fair **Rock Paper Scissors** game using the NanoShutter API. You've learned how to integrate time-based encryption and decryption to ensure players cannot cheat by viewing moves prematurely.

**Remember**, NanoShutter is perfect for hackathons and projects that require secrecy in a distributed manner. Enjoy building, and happy hacking!

