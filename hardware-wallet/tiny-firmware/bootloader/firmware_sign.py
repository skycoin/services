#!/usr/bin/python
from __future__ import print_function
import argparse
import hashlib
import struct
import binascii
import skycoin_crypto
import random

try:
    raw_input
except:
    raw_input = input

SLOTS = 3

pubkeys = {
    1: '025839078e6c11c09ad4b00092f5feefd566f92af60f2c71facfb01d2b84c043d4',
    2: '029170192c2fdefb4d377af9e196e06d1176f86f733523d3950f90ff84c0cd02d3',
    3: '03338ffc0ff42df07d27b0b4131cd96ffdfa4685b5566aafc7aa71ed10fd1cbd6f',
    4: '039f12c93645e35e5274dc38f191be0b6d1321ec35d2d2a3ddf7d13ed12f6da85b',
    5: '03b17c7b7c564385be66f9c1b9da6a0b5aea56f0cb70548e6528a2f4f7b27245d8',
}

INDEXES_START = len('TRZR') + struct.calcsize('<I')
SIG_START = INDEXES_START + SLOTS + 1 + 52

def parse_args():
    parser = argparse.ArgumentParser(description='Commandline tool for signing Trezor firmware.')
    parser.add_argument('-f', '--file', dest='path', help="Firmware file to modify")
    parser.add_argument('-s', '--sign', dest='sign', action='store_true', help="Add signature to firmware slot")

    return parser.parse_args()

def prepare(data):
    # Takes raw OR signed firmware and clean out metadata structure
    # This produces 'clean' data for signing

    meta = b'TRZR'  # magic
    if data[:4] == b'TRZR':
        meta += data[4:4 + struct.calcsize('<I')]
    else:
        meta += struct.pack('<I', len(data))  # length of the code
    meta += b'\x00' * SLOTS  # signature index #1-#3
    meta += b'\x01'       # flags
    meta += b'\x00' * 52  # reserved
    meta += b'\x00' * 64 * SLOTS  # signature #1-#3

    if data[:4] == b'TRZR':
        # Replace existing header
        out = meta + data[len(meta):]
    else:
        # create data from meta + code
        out = meta + data

    return out

def check_signatures(data):
    # Analyses given firmware and prints out
    # status of included signatures

    try:
        indexes = [ ord(x) for x in data[INDEXES_START:INDEXES_START + SLOTS] ]
    except:
        indexes = [ x for x in data[INDEXES_START:INDEXES_START + SLOTS] ]

    to_sign = prepare(data)[256:] # without meta
    fingerprint = hashlib.sha256(to_sign).hexdigest()
    print("Firmware fingerprint:", fingerprint)

    used = []
    for x in range(SLOTS):
        signature = data[SIG_START + 64 * x:SIG_START + 64 * x + 64]

        if indexes[x] == 0:
            print("Slot #%d" % (x + 1), 'is empty')
        else:
            pk = pubkeys[indexes[x]]

            skycoin = skycoin_crypto.SkycoinCrypto()
            pubkey = skycoin.RecoverPubkeyFromSignature(binascii.unhexlify(fingerprint), signature)
            pubkey = binascii.hexlify(pubkey)
            
            if (pubkey == pk):
                if indexes[x] in used:
                    print("Slot #%d signature: DUPLICATE" % (x + 1), binascii.hexlify(signature))
                else:
                    used.append(indexes[x])
                    print("Slot #%d signature: VALID" % (x + 1), binascii.hexlify(signature))
            else:
                print("Slot #%d signature: INVALID" % (x + 1), binascii.hexlify(signature))


def modify(data, slot, index, signature):
    # Replace signature in data

    # Put index to data
    data = data[:INDEXES_START + slot - 1 ] + chr(index) + data[INDEXES_START + slot:]

    # Put signature to data
    data = data[:SIG_START + 64 * (slot - 1) ] + signature + data[SIG_START + 64 * slot:]

    return data

def sign(data):
    # Ask for index and private key and signs the firmware

    slot = int(raw_input('Enter signature slot (1-%d): ' % SLOTS))
    if slot < 1 or slot > SLOTS:
        raise Exception("Invalid slot")

    print("Paste SECEXP (in hex) and press Enter:")
    print("(blank private key removes the signature on given index)")
    secexp = raw_input()
    if secexp.strip() == '':
        # Blank key,let's remove existing signature from slot
        return modify(data, slot, 0, '\x00' * 64)
    skycoin = skycoin_crypto.SkycoinCrypto()
    seckey = binascii.unhexlify(secexp)
    pubkey = skycoin.GeneratePubkeyFromSeckey(seckey)
    pubkey = binascii.hexlify(pubkey.value)

    to_sign = prepare(data)[256:] # without meta
    fingerprint = hashlib.sha256(to_sign).hexdigest()
    print("Firmware fingerprint:", fingerprint)

    # Locate proper index of current signing key
    # pubkey = b'04' + binascii.hexlify(key.get_verifying_key().to_string())
    index = None
    for i, pk in pubkeys.items():
        if pk == pubkey:
            index = i
            break

    if index == None:
        raise Exception("Unable to find private key index. Unknown private key?")

    for i in range(10): #attempt 10 times until finding a signature with 64 bytes (the 65's byte will be assumed as 0)
        signature = skycoin.EcdsaSkycoinSign(binascii.unhexlify(fingerprint), seckey, random.randint(1, 0xFFFFFFFF))
        if len(signature.value) == 64:
            break

    print("Skycoin signature:", binascii.hexlify(signature.value))
    if len(signature.value) != 64:
        raise Exception("Signature lenght {} is not correct".format(len(signature.value)))

    return modify(data, slot, index, str(signature.value))

def main(args):

    if not args.path:
        raise Exception("-f/--file is required")

    data = open(args.path, 'rb').read()
    assert len(data) % 4 == 0

    if data[:4] != b'TRZR':
        print("Metadata has been added...")
        data = prepare(data)

    if data[:4] != b'TRZR':
        raise Exception("Firmware header expected")

    print("Firmware size %d bytes" % len(data))

    check_signatures(data)

    if args.sign:
        data = sign(data)
        check_signatures(data)

    fp = open(args.path, 'wb')
    fp.write(data)
    fp.close()

if __name__ == '__main__':
    args = parse_args()
    main(args)
