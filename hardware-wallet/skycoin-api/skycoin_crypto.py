from ctypes import cdll, c_char_p, c_uint32, addressof, create_string_buffer
import binascii
lib = cdll.LoadLibrary('./libskycoin-crypto.so')

class SkycoinCrypto(object):
    def __init__(self):
        pass

    def EcdsaSkycoinSign(self, digest, seckey, seed=1):
        signature = create_string_buffer(65)
        lib.ecdsa_skycoin_sign(c_uint32(seed), addressof(seckey), addressof(digest), signature)
        return signature.value
