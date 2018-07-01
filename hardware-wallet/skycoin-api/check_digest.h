/* Code inspired by 
 * https://github.com/kokke/tiny-regex-c/blob/master/re.c
 * */


static bool matchhexdigit(char c)
{
  return ((c >= '0') && (c <= '9')) || ((c >= 'a') && (c <= 'f')) || ((c >= 'A') && (c <= 'F'));
}

int is_digest(char* digest) {
    if (strlen(digest) != 64) {
        return false;
    }
    bool bDigest = true;
    for (int i = 0; i < 64 && bDigest; ++i) {
        bDigest &= matchhexdigit(digest[i]);
    }
    return bDigest;
}