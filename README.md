# NanoShutter

NanoShutter is a simplified implementation of the [Shutter Network](https://shutter.network/) protocol designed specifically for hackathons and educational purposes. It provides a lightweight framework for developers to integrate time-locked encryption into their decentralized applications (DApps) or other projects without the complexity of the full Shutter protocol. It's fully centralized and unsafe to use in any kind of production environment.

NanoShutter allows you to encrypt messages or transactions that can only be decrypted after a certain time (epoch), enabling use cases like front-running protection, sealed-bid auctions, or fair multiplayer games.

![ezgif-3-0169c86877](https://github.com/user-attachments/assets/61c6573e-ede3-4fcc-82e8-39d61a92af40)


---

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Architecture Overview](#architecture-overview)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [1. Running the NanoShutter API Server](#1-running-the-nanoshutter-api-server)
  - [2. Integrating the Client-Side Library](#2-integrating-the-client-side-library)
  - [3. Example Application: Shutterized Rock Paper Scissors](#3-example-application-shutterized-rock-paper-scissors)
- [API Reference](#api-reference)
- [How It Works](#how-it-works)
- [Customization](#customization)
- [Limitations](#limitations)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgments](#acknowledgments)

---

## Introduction

**NanoShutter** is tailored to help developers quickly prototype and demonstrate applications that require time-locked encryption. By abstracting the complexities of cryptographic operations and key management, NanoShutter enables you to focus on building your application's logic.

Whether you're participating in a hackathon, teaching cryptography concepts, or exploring new ideas, NanoShutter provides an accessible entry point to time-based encryption mechanisms.

---

## Features

- **Simplified API:** Easy-to-use endpoints for fetching public keys and decryption keys.
- **Client-Side Library:** JavaScript library (`nano_shutter_crypto.js`) for encryption and decryption operations.
- **Time-Locked Encryption:** Encrypt messages that can only be decrypted after a specified epoch.
- **Lightweight Server:** Minimal Flask server (`nano_shutter_api.py`) to manage keys and epochs.
- **Educational Tool:** Ideal for learning and teaching cryptographic concepts.
- **Hackathon Ready:** Quickly integrate into projects to add security and fairness features.

---

## Architecture Overview

NanoShutter consists of two main components:

1. **NanoShutter API Server (`nano_shutter_api.py`):**
   - Manages cryptographic keys (eon keys and epoch keys).
   - Provides endpoints to fetch the eon public key, epoch public keys, and decryption keys.
   - Simulates the behavior of a keyper committee in the Shutter Network.

2. **Client-Side Library (`nano_shutter_crypto.js`):**
   - Handles encryption and decryption of messages.
   - Communicates with the NanoShutter API Server to obtain necessary keys.
   - Provides functions to integrate time-locked encryption into your application.

---

## Getting Started

### Prerequisites

- **Python 3.6+**
- **Node.js and npm** (for frontend development)
- **Flask** and **PyNaCl** Python packages
- Modern web browser (Google Chrome, Mozilla Firefox, etc.)

### Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/your-username/nanoshutter.git
   cd nanoshutter
   ```

2. **Install Python Dependencies:**

   ```bash
   pip install flask pynacl
   ```

3. **Install Frontend Dependencies (if using the example application):**

   ```bash
   cd frontend
   npm install
   ```

---

## Usage

### 1. Running the NanoShutter API Server

Start the NanoShutter API server to manage cryptographic keys:

```bash
python nano_shutter_api.py
```

By default, the server runs on `http://localhost:5000`.

### 2. Integrating the Client-Side Library

Include the `nano_shutter_crypto.js` script in your HTML file:

```html
<script src="nano_shutter_crypto.js"></script>
```

Initialize the NanoShutter client in your JavaScript code:

```javascript
const nanoShutter = new NanoShutterCrypto('http://localhost:5000');
await nanoShutter.initialize();
await nanoShutter.fetchEonKey();
```

Use the provided methods to encrypt and decrypt messages:

```javascript
// Encrypt a message
const encryptedData = await nanoShutter.encryptMessage("Your secret message");

// Decrypt a message
const decryptedMessage = await nanoShutter.decryptMessage(encryptedData);
```

### 3. Example Application: Shutterized Rock Paper Scissors

We've included an example application demonstrating how to use NanoShutter in a real-world scenario.

#### Running the Example

1. **Start the API Server:**

   ```bash
   python nano_shutter_api.py
   ```

2. **Open the `index.html` File:**

   - Open `index.html` in your web browser.
   - Alternatively, serve it using a local server:

     ```bash
     cd frontend
     npx http-server -p 8000
     ```

     Navigate to `http://localhost:8000` in your browser.

#### How the Example Works

- Two players select and submit their moves (Rock, Paper, or Scissors).
- Moves are encrypted using NanoShutter and displayed on the page.
- The API server releases decryption keys after the epoch ends.
- The client decrypts the moves and determines the winner.

---

## API Reference

### Endpoints

- **GET `/eon-key`**

  Retrieves the eon public key and epoch duration.

  **Response:**

  ```json
  {
    "eon_public_key": "hexadecimal_string",
    "epoch_duration": 10,
    "current_epoch": 123456
  }
  ```

- **GET `/epoch-public-key?epoch=<epoch_number>`**

  Retrieves the epoch public key for the specified epoch.

  **Response:**

  ```json
  {
    "epoch_public_key": "hexadecimal_string"
  }
  ```

- **GET `/decryption-key?epoch=<epoch_number>&ephemeral_public_key=<hex_string>`**

  Retrieves the decryption key for the specified epoch and ephemeral public key.

  **Response:**

  ```json
  {
    "epoch_private_key": "hexadecimal_string"
  }
  ```

---

## How It Works

1. **Eon Key Generation:**

   - The server generates a long-term eon key pair (public and private keys).

2. **Epoch Key Derivation:**

   - Every epoch (fixed time interval), the server derives a new epoch private key from the eon private key and the current epoch number.
   - The epoch public key is derived from the epoch private key.

3. **Message Encryption:**

   - The client fetches the current epoch public key.
   - The client generates an ephemeral key pair.
   - A shared secret is derived using the client's ephemeral private key and the server's epoch public key.
   - The shared secret and epoch number are hashed to produce a symmetric key.
   - The message is encrypted using this symmetric key.

4. **Message Decryption:**

   - After the epoch ends, the server releases the epoch private key.
   - The client re-derives the shared secret using the epoch private key and the ephemeral public key.
   - The symmetric key is recomputed, and the message is decrypted.

---

## Customization

- **Epoch Duration:**

  Adjust the `EPOCH_DURATION` variable in `nano_shutter_api.py` to change the duration of each epoch.

  ```python
  EPOCH_DURATION = 10  # seconds
  ```

- **Key Management:**

  For more advanced use cases, modify the key derivation methods to suit your needs.

- **Client-Side Modifications:**

  Extend the `NanoShutterCrypto` class with additional methods or integrate it into frameworks like React or Angular.

---

## Limitations

- **Security:**

  NanoShutter is intended for educational and demonstration purposes. It lacks the robustness and security features of the full Shutter protocol.

- **Single Point of Failure:**

  The API server acts as a centralized component, which is not suitable for production environments requiring decentralization.

- **No Consensus Mechanism:**

  NanoShutter does not implement a consensus mechanism or handle adversarial conditions.

---

## Contributing

We welcome contributions! Please follow these steps:

1. **Fork the Repository**

2. **Create a Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Commit Your Changes**

   ```bash
   git commit -m "Add your message"
   ```

4. **Push to the Branch**

   ```bash
   git push origin feature/your-feature-name
   ```

5. **Open a Pull Request**

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Acknowledgments

- [Shutter Network](https://shutter.network/) for the inspiration and underlying concepts.
- [libsodium](https://libsodium.gitbook.io/doc/) and [PyNaCl](https://pynacl.readthedocs.io/) for cryptographic functions.
- Participants and organizers of hackathons promoting innovation and learning.

---

**Disclaimer:** NanoShutter is not affiliated with or endorsed by the Shutter Network. It is a simplified tool created for learning and hackathon projects.

**Contact:** For questions or support, please open an issue on the repository.

---

Happy hacking!
