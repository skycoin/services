
## Prerequisites

Before using the Go bindings, you must install the libpostal C library. Make sure you have the following prerequisites:

**On Ubuntu/Debian**
```
sudo apt-get install curl autoconf automake libtool pkg-config
```

**On CentOS/RHEL**
```
sudo yum install curl autoconf automake libtool pkgconfig
```

**On Mac OSX**
```
sudo brew install curl autoconf automake libtool pkg-config
```

**Installing libpostal**

```bash
git clone https://github.com/openvenues/libpostal
cd libpostal
./bootstrap.sh
./configure --datadir=[...some dir with a few GB of space...]
make
sudo make install

# On Linux it's probably a good idea to run
sudo ldconfig
```

**Installing cli tool**

```
cd services/address-parser
go install
```

**Usage**

*Parse from a text file*
As a file is used a text file separated by '\n'.
```
address file <path-to-file>
```

Example:
```bash
address-parser file ~/tmp/addrs.txt
```
Result:
```json
[
    {
        "original": "781 Franklin Ave Crown Heights Brooklyn NY 11216 USA",
        "parsed": [
            "781 franklin avenue crown heights brooklyn ny 11216 usa",
            "781 franklin avenue crown heights brooklyn new york 11216 usa"
        ]
    }
]
```


*Parse from command line*
```
address-parser address <query>
```

Example:
```bash
address-parser address "781 Franklin Ave Crown Heights Brooklyn NY 11216 USA"
```
Result:
```json
[
    {
        "original": "781 Franklin Ave Crown Heights Brooklyn NY 11216 USA",
        "parsed": [
            "781 franklin avenue crown heights brooklyn ny 11216 usa",
            "781 franklin avenue crown heights brooklyn new york 11216 usa"
        ]
    }
]
```




