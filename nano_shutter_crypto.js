// nano_shutter_crypto.js

class NanoShutterCrypto {
    constructor(apiUrl) {
      this.apiUrl = apiUrl;
      this.eonPublicKey = null;
      this.epochDuration = null;
      this.currentEpoch = null;
    }
  
    async initialize() {
      // Wait for sodium to be ready
      await sodium.ready;
    }
  
    async fetchEonKey() {
      const response = await axios.get(`${this.apiUrl}/eon-key`);
      this.eonPublicKey = sodium.from_hex(response.data.eon_public_key);
  
      // Ensure eonPublicKey is a Uint8Array of length 32
      if (this.eonPublicKey.length !== sodium.crypto_scalarmult_BYTES) {
        throw new Error('Invalid Eon public key length');
      }
  
      this.epochDuration = response.data.epoch_duration;
      this.currentEpoch = response.data.current_epoch;
    }
  
    getCurrentEpoch() {
      if (!this.epochDuration) {
        throw new Error('Epoch duration is not defined. Call fetchEonKey() first.');
      }
      // Update the current epoch based on the current time
      const currentTime = Math.floor(Date.now() / 1000);
      const epochNumber = Math.floor(currentTime / this.epochDuration);
      return epochNumber;
    }
  
    async encryptMessage(message) {
      if (!this.eonPublicKey || !this.epochDuration) {
        await this.fetchEonKey();
      }
  
      const epoch = this.getCurrentEpoch();
  
      // Fetch the epoch public key
      const epochPubKeyResponse = await axios.get(`${this.apiUrl}/epoch-public-key`, {
        params: {
          epoch: epoch,
        },
      });
      const epochPublicKey = sodium.from_hex(epochPubKeyResponse.data.epoch_public_key);
  
      // Generate ephemeral key pair
      // Use crypto_box_keypair for compatibility with libsodium-wrappers 0.5.2
      const ephemeralKeyPair = sodium.crypto_box_keypair();
      const ephemeralPrivateKey = ephemeralKeyPair.privateKey;
      const ephemeralPublicKey = ephemeralKeyPair.publicKey;
  
      console.log('Ephemeral private key (client side, encryption):', sodium.to_hex(ephemeralPrivateKey));
      console.log('Ephemeral public key (client side, encryption):', sodium.to_hex(ephemeralPublicKey));
  
      // Compute shared secret using the epoch public key
      const sharedSecret = sodium.crypto_scalarmult(ephemeralPrivateKey, epochPublicKey);
      console.log('Shared secret derived (client side, encryption):', sodium.to_hex(sharedSecret));
  
      const epochBytes = new Uint8Array(8);
      const dataView = new DataView(epochBytes.buffer);
      dataView.setBigUint64(0, BigInt(epoch), false);
      console.log('Epoch bytes (client side, encryption):', sodium.to_hex(epochBytes));
  
      // Combine sharedSecret with epochBytes and hash
      const combinedInput = concatUint8Arrays([sharedSecret, epochBytes]);
      console.log('Combined input for hash (client side, encryption):', sodium.to_hex(combinedInput));
  
      const combinedSecretHash = sodium.crypto_generichash(32, combinedInput);
      console.log('Combined secret hash (client side, encryption):', sodium.to_hex(combinedSecretHash));
  
      const nonce = sodium.randombytes_buf(sodium.crypto_secretbox_NONCEBYTES);
      console.log('Nonce (client side, encryption):', sodium.to_hex(nonce)); // Log the nonce for debugging
  
      const messageUint8 = sodium.from_string(message);
      const encryptedMessage = sodium.crypto_secretbox_easy(messageUint8, nonce, combinedSecretHash);
  
      return {
        ephemeralPublicKey: sodium.to_hex(ephemeralPublicKey),
        nonce: sodium.to_hex(nonce),
        ciphertext: sodium.to_hex(encryptedMessage),
        epoch: epoch, // Include the epoch number
      };
    }
  
    async decryptMessage(encryptedData) {
      if (!this.eonPublicKey || !this.epochDuration) {
        await this.fetchEonKey();
      }
  
      const epoch = Number(encryptedData.epoch);
      console.log('Attempting to decrypt message for epoch:', epoch);
  
      try {
        // Fetch the decryption key for the message's epoch
        const response = await axios.get(`${this.apiUrl}/decryption-key`, {
          params: {
            epoch: epoch,
            ephemeral_public_key: encryptedData.ephemeralPublicKey, // Send ephemeral public key for server-side derivation
          },
        });
        const epochPrivateKey = sodium.from_hex(response.data.epoch_private_key);
        console.log('Epoch private key (client side, decryption):', sodium.to_hex(epochPrivateKey));
        console.log('Ephemeral public key (client side, decryption):', encryptedData.ephemeralPublicKey);
  
        const sharedSecret = sodium.crypto_scalarmult(
          epochPrivateKey,
          sodium.from_hex(encryptedData.ephemeralPublicKey)
        );
  
        console.log('Shared secret derived (client side, decryption):', sodium.to_hex(sharedSecret));
  
        const epochBytes = new Uint8Array(8);
        const dataView = new DataView(epochBytes.buffer);
        dataView.setBigUint64(0, BigInt(epoch), false);
        console.log('Epoch bytes (client side, decryption):', sodium.to_hex(epochBytes));
  
        // Combine sharedSecret with epochBytes and hash
        const combinedInput = concatUint8Arrays([sharedSecret, epochBytes]);
        console.log('Combined input for hash (client side, decryption):', sodium.to_hex(combinedInput));
  
        const combinedSecretHash = sodium.crypto_generichash(32, combinedInput);
        console.log('Combined secret hash (client side, decryption):', sodium.to_hex(combinedSecretHash));
  
        console.log('Nonce used for decryption:', encryptedData.nonce); // Log nonce for verification
  
        const decryptedMessage = sodium.crypto_secretbox_open_easy(
          sodium.from_hex(encryptedData.ciphertext),
          sodium.from_hex(encryptedData.nonce),
          combinedSecretHash
        );
  
        if (!decryptedMessage) {
          throw new Error('Failed to decrypt message.');
        }
  
        return sodium.to_string(decryptedMessage);
      } catch (error) {
        if (error.response && error.response.status === 404) {
          console.log('Decryption key not available yet.');
          throw new Error('Decryption key not available yet.');
        } else {
          console.error('Error during decryption:', error);
          throw error;
        }
      }
    }
  }
  
  // Helper function to concatenate Uint8Arrays
  function concatUint8Arrays(arrays) {
    // Calculate total length of all arrays
    let totalLength = arrays.reduce((sum, arr) => sum + arr.length, 0);
  
    // Create a new Uint8Array with the total length
    let result = new Uint8Array(totalLength);
  
    // Concatenate all arrays into the result
    let offset = 0;
    for (let arr of arrays) {
      result.set(arr, offset);
      offset += arr.length;
    }
  
    return result;
  }
  
  // Expose the class to the global scope
  window.NanoShutterCrypto = NanoShutterCrypto;
  