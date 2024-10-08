#! usr/bin/env bash

# requires: curl
# may require sudo

# ~
curl -L -o ghpm https://github.com/Neal-C/ghpm/releases/download/v0.1.0-rc/ghpm-linux-amd64

# ~
chmod +x ghpm

# ~
# might require to be prefixed with sudo
sudo mv ghpm usr/local/bin/ghpm