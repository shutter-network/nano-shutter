# nano_shutter_api.py

from flask import Flask, jsonify, request
from flask_cors import CORS
from datetime import datetime, timezone
import nacl.utils
from nacl.public import PrivateKey, PublicKey
import nacl.encoding
import hashlib
import logging
from nacl.bindings import crypto_scalarmult

app = Flask(__name__)
CORS(app)

# Enable detailed logging for debugging purposes
logging.basicConfig(level=logging.DEBUG)

EPOCH_DURATION = 10  # seconds

# Generate Eon key at startup
eon_private_key = PrivateKey.generate()
eon_public_key = eon_private_key.public_key

def get_current_epoch():
    current_time = datetime.now(timezone.utc)
    epoch_number = int(current_time.timestamp() // EPOCH_DURATION)
    return epoch_number

def derive_epoch_private_key(epoch):
    try:
        epoch_bytes = epoch.to_bytes(8, byteorder='big')
        eon_private_bytes = eon_private_key.encode()  # 32 bytes

        # Derive the epoch private key using HMAC-SHA512
        epoch_key_material = hashlib.sha512(eon_private_bytes + epoch_bytes).digest()
        epoch_private_key_bytes = epoch_key_material[:32]  # Use first 32 bytes

        # Log the derived key material for debugging
        logging.debug(f"Derived epoch private key (bytes): {epoch_private_key_bytes.hex()}")

        # Create a PrivateKey object for returning the private key
        epoch_private_key = PrivateKey(epoch_private_key_bytes)
        return epoch_private_key
    except Exception as e:
        logging.error(f"Error deriving epoch private key: {e}")
        raise

@app.route('/eon-key', methods=['GET'])
def get_eon_key():
    try:
        eon_public_hex = eon_public_key.encode(encoder=nacl.encoding.HexEncoder).decode('utf-8')
        current_epoch = get_current_epoch()
        return jsonify({
            'eon_public_key': eon_public_hex,
            'epoch_duration': EPOCH_DURATION,
            'current_epoch': current_epoch
        })
    except Exception as e:
        logging.error(f"Error in /eon-key endpoint: {e}")
        return jsonify({'error': 'Internal Server Error'}), 500

@app.route('/epoch-public-key', methods=['GET'])
def get_epoch_public_key():
    try:
        epoch = request.args.get('epoch', type=int)
        if epoch is None:
            return jsonify({'error': 'Epoch parameter is required.'}), 400
        # Derive epoch private key
        epoch_private_key = derive_epoch_private_key(epoch)
        epoch_public_key = epoch_private_key.public_key
        epoch_public_hex = epoch_public_key.encode(encoder=nacl.encoding.HexEncoder).decode('utf-8')
        return jsonify({
            'epoch_public_key': epoch_public_hex,
            'epoch': epoch
        })
    except Exception as e:
        logging.error(f"Error in /epoch-public-key endpoint: {e}")
        return jsonify({'error': 'Internal Server Error'}), 500

@app.route('/decryption-key', methods=['GET'])
def get_decryption_key():
    try:
        epoch = request.args.get('epoch', type=int)
        ephemeral_public_key_hex = request.args.get('ephemeral_public_key')
        if epoch is None or ephemeral_public_key_hex is None:
            logging.warning('Epoch or ephemeral public key parameter is missing or invalid.')
            return jsonify({'error': 'Epoch and ephemeral public key parameters are required.'}), 400

        current_epoch = get_current_epoch()
        if epoch >= current_epoch:
            time_remaining = (epoch + 1) * EPOCH_DURATION - datetime.now(timezone.utc).timestamp()
            logging.info(f'Decryption key for epoch {epoch} not available yet. Time remaining: {time_remaining} seconds.')
            return jsonify({
                'error': 'Decryption key not available yet.',
                'seconds_until_available': int(time_remaining)
            }), 404

        # Convert the ephemeral public key to bytes
        ephemeral_public_key_bytes = bytes.fromhex(ephemeral_public_key_hex)
        logging.debug(f'Received ephemeral public key (server side): {ephemeral_public_key_hex}')

        # Derive epoch private key
        epoch_private_key = derive_epoch_private_key(epoch)
        epoch_private_key_bytes = epoch_private_key.encode()  # Convert to bytes
        logging.debug(f'Epoch private key (server side): {epoch_private_key_bytes.hex()}')

        # Derive shared secret using scalar multiplication
        shared_secret = crypto_scalarmult(epoch_private_key_bytes, ephemeral_public_key_bytes)
        logging.debug(f'Shared secret derived (server side, scalar multiplication inputs: '
              f'epoch_private_key: {epoch_private_key_bytes.hex()}, ephemeral_public_key: {ephemeral_public_key_hex}): {shared_secret.hex()}')

        # Return the epoch private key for client-side decryption
        epoch_private_hex = epoch_private_key_bytes.hex()
        logging.info(f'Decryption key for epoch {epoch} released.')
        return jsonify({'epoch_private_key': epoch_private_hex})

    except Exception as e:
        logging.error(f"Error in /decryption-key endpoint: {e}")
        return jsonify({'error': 'Internal Server Error'}), 500

if __name__ == '__main__':
    # Run the Flask app
    app.run(host='0.0.0.0', port=5000, threaded=True)
