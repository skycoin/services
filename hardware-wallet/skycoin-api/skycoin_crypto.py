from ctypes import cdll, c_char_p, c_uint32, c_size_t, addressof, create_string_buffer
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