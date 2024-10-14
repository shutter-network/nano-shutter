# NanoShutter

NanoShutter is a simplified implementation of the [Shutter Network](https://shutter.network/) protocol designed specifically for hackathons and educational purposes. It provides a lightweight framework for developers to integrate time-locked encryption into their decentralized applications (DApps) or other projects without the complexity of the full Shutter protocol. It's fully centralized and unsafe to use in any kind of production environment.

NanoShutter allows you to encrypt messages or transactions that can only be decrypted after a certain time (epoch), enabling use cases like front-running protection, sealed-bid auctions, or fair multiplayer games.

![ezgif-6-1044641a13](https://github.com/user-attachments/assets/e701a287-4f8c-4872-aca1-f6eb5de63f8a)


## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [API Endpoints](#api-endpoints)
- [Getting Started](#getting-started)
  - [Using Our Provided API](#using-our-provided-api)
  - [Running Your Own API](#running-your-own-api)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Encrypting and Decrypting](#encrypting-and-decrypting)
- [Development](#development)
  - [Key Components](#key-components)
  - [Docker Deployment](#docker-deployment)
- [Contributing](#contributing)
- [License](#license)

---

## Introduction

**NanoShutter** offers a simplified version of Shutter Network for developers to quickly prototype and learn. It is a Go-based API server providing encryption and decryption capabilities through four RESTful endpoints. These endpoints use time-based or custom epochs to manage cryptographic operations securely.

## Features

- **Time-Based Encryption/Decryption**: Encrypt messages that can only be decrypted after a certain time.
- **Custom Epoch Encryption/Decryption**: Users can specify custom epochs for cryptographic operations.
- **REST API Interface**: Four straightforward API endpoints to interact with.
- **Lightweight and Hackathon-Friendly**: Simplified setup, perfect for experimentation and learning.

## API Endpoints

The NanoShutter API provides four main endpoints:

1. **Time-Based Encryption**: `/encrypt/with_time`
   - Encrypts a message for a future decryption, based on a timestamp in the future.

2. **Time-Based Decryption**: `/decrypt/with_time`
   - Decrypts a message after the specified time has elapsed.

3. **Custom Epoch Encryption**: `/encrypt/custom`
   - Encrypts a message using a custom epoch ID provided by the user.

4. **Custom Epoch Decryption**: `/decrypt/custom`
   - Decrypts a message based on a custom epoch ID.

### Testing URLs

We provide the following testing URLs that you can use to interact with the NanoShutter API without needing to deploy your own:

- **Time-Based Encryption/Decryption**
  - `https://nanoshutter.staging.shutter.network/encrypt/with_time`
  - `https://nanoshutter.staging.shutter.network/decrypt/with_time`
  
- **Custom Epoch Encryption/Decryption**
  - `https://nanoshutter.staging.shutter.network/encrypt/custom`
  - `https://nanoshutter.staging.shutter.network/decrypt/custom`


## Hackathon Ideas

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


## Getting Started

### Using Our Provided API

If you want to start quickly without setting up your own API server, you can use the provided endpoints:

1. Use any REST client (e.g., Postman, Insomnia, or cURL) to interact with the NanoShutter endpoints.
2. Send POST requests to the provided URLs with the required JSON body (see examples in [Usage](#usage)).
3. Decrypt the messages when the required time or epoch conditions are met.

This is the fastest way to get started and understand how NanoShutter works without dealing with infrastructure setup.

### Running Your Own API

If you want to operate your own instance of the NanoShutter API, follow these steps:

#### Prerequisites

- **Go** (version 1.16 or higher)
- **Docker** (optional, for containerized deployment)

#### Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/your-username/nanoshutter.git
   cd nanoshutter/api
   ```

2. **Install Dependencies:**

   Use `go mod` to install required Go modules:

   ```bash
   go mod tidy
   ```

#### Running the Service

##### Running Locally

1. **Environment Configuration:**

   Copy `.env.example` to `.env` and modify environment variables as needed:

   ```bash
   cp .env.example .env
   ```

2. **Run the Server:**

   ```bash
   go run main.go
   ```

   By default, the server runs on `http://localhost:8080`.

##### Running with Docker

1. **Build Docker Image:**

   ```bash
   docker build -t nanoshutter .
   ```

2. **Run the Docker Container:**

   ```bash
   docker-compose up -d
   ```

   This will start the NanoShutter service on the specified port defined in the `docker-compose.yml`.

## Configuration

Use the `.env` file to configure the application settings, such as:

- **PORT**: The port on which the service runs.

### Example `.env` File:

```env
PORT=8080
```

## Usage

### Encrypting and Decrypting

#### 1. **Time-Based Encryption**

Send a POST request to `/encrypt/with_time` with the message and timestamp to encrypt (e.g. 60 seconds in the future).

**Request Body Example:**

```json
{
  "cypher_text": "This is a secret message",
  "timestamp": 1696787387
}
```

**Response Example:**

```json
{
  "message": "encrypted_hex_string"
}
```

#### 2. **Time-Based Decryption**

Send a POST request to `/decrypt/with_time` with the encrypted message and the same timestamp used previously.

**Request Body Example:**

```json
{
  "encrypted_msg": "encrypted_hex_string",
  "timestamp": 1696787387
}
```

**Response Example:**

```json
{
  "message": "This is a secret message"
}
```

#### 3. **Custom Epoch Encryption**

Send a POST request to `/encrypt/custom` with the message and a custom epoch ID.

**Request Body Example:**

```json
{
  "cypher_text": "Custom epoch secret message",
  "epoch_id": "12345"
}
```

**Response Example:**

```json
{
  "message": "encrypted_hex_string"
}
```

#### 4. **Custom Epoch Decryption**

Send a POST request to `/decrypt/custom` with the encrypted message and epoch ID.

**Request Body Example:**

```json
{
  "encrypted_msg": "encrypted_hex_string",
  "epoch_id": "12345"
}
```

**Response Example:**

```json
{
  "message": "Custom epoch secret message"
}
```

## Development

### Key Components

- **`main.go`**: The main entry point of the service.
- **`service/service.go`**: Implements core business logic for the encryption/decryption API.
- **`router/router.go`**: Defines the HTTP routes and handlers.
- **`dkg/dkg.go`**: Handles Distributed Key Generation (DKG) functionality.
- **`internal/error/error.go`**: Centralized error handling.
- **`internal/middleware/error.go`**: Middleware for handling request errors.
- **Docker and `.env` Files**: To assist with containerized deployments.

### Docker Deployment

To deploy the service using Docker:

1. **Build and Run Using Docker Compose:**

   ```bash
   docker-compose up --build
   ```

   This will build the image and run the container.

## Contributing

Contributions are welcome! Please follow these steps:

1. **Fork the Project**
2. **Create a Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Commit Your Changes**

   ```bash
   git commit -m "Add your feature"
   ```

4. **Push to Branch**

   ```bash
   git push origin feature/your-feature-name
   ```

5. **Open a Pull Request**

---

## License

This project is licensed under the MIT License - see the LICENSE file for details.
