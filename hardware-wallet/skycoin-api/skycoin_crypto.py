from ctypes import cdll, c_char_p, c_uint32, c_size_t, byref, addressof, create_string_buffer
import binascii

class SkycoinCrypto(object):
    def __init__(self):
        self.lib = cdll.LoadLibrary('./libskycoin-crypto.so')

    def EcdsaSkycoinSign(self, digest, seckey, seed=1):
        signature = create_string_buffer(65)
        self.lib.ecdsa_skycoin_sign(c_uint32(seed), seckey, digest, signature)
        return signature
    
    def ComputeSha256Sum(self, seed):
        digest = create_string_buffer(32)
        self.lib.compute_sha256sum(seed, digest, self.lib.strlen(seed))
        return digest

    def GeneratePubkeyFromSeckey(self, seckey):
        pubkey = create_string_buffer(33)
        self.lib.generate_pubkey_from_seckey(seckey, pubkey)
        return pubkey

    def Base58AddressFromPubkey(self, pubkey):
        address = create_string_buffer(36)
        address_size = c_size_t(36)
        self.lib.generate_base58_address_from_pubkey(pubkey, address, byref(address_size))
        return address

    def RecoverPubkeyFromSignature(self, message, signature):
        pubkey = create_string_buffer(33)
        self.lib.recover_pubkey_from_signed_message(message, signature, pubkey)
        return pubkey