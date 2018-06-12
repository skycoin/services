# Bitmap Generation

`bitmaps\generate.py` is a Python script that searches for every PNG file in the current directory (usually `(...)/gen/bitmaps/`) and converts them into a `bitmaps.c` and `bitmaps.h` file.

## To add a bitmap image:

1. Add a PNG file to the bitmaps folder.
2. Ensure that Python 2 is installed and in the PATH: `python`
3. Set the working directory: `cd bitmaps`
4. Run the generate script: `python generate.py`


