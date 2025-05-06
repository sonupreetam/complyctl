# Installation

## Binary

- The latest binary release can be downloaded from <https://github.com/complytime/complytime/releases/latest>.
- The release signature can be verified with:
  ```
  cosign verify-blob --certificate complytime_*_checksums.txt.pem --signature complytime_*_checksums.txt.sig complytime_*_checksums.txt --certificate-oidc-issuer=https://token.actions.githubusercontent.com --certificate-identity=https://github.com/complytime/complytime/.github/workflows/release.yml@refs/heads/main
  ```


## From Source

### Prerequisites

- **Go** version 1.20 or higher
- **Make** (optional, for using the `Makefile` if included)
- **pandoc** (optional, for generating man pages using the `make man`)

### Clone the repository

```bash
git clone https://github.com/complytime/complytime.git
cd complytime
```

### Build Instructions
To compile complytime and openscap-plugin:

```bash
make build
```

The binaries can be found in the `bin/` directory in the local repo. Add it to your PATH and you are all set!
