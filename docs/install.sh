#! usr/bin/env bash

# requires: curl unzip
# may require sudo

# ~
curl -L -o ghpm.gz  https://github.com/Neal-C/ghpm/releases/download/v0.1.0-rc/ghpm-linux-amd64.zip

# ~
gunzip ghpm.gz

# ~
chmod +x ghpm

# ~
# might require to be prefixed with sudo
sudo mv ghpm usr/local/bin/ghpm