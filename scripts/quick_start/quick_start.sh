#!/bin/bash
# How to run the script
# Assume that you have already installed a fresh RHEL
# 1. Download/copy this script in your RHEL system
# 2. Run the script
# chmod +x quick_start.sh
# export RHEL_APPS_REPO=$RHEL_APPS_REPO
# sh quick_start.sh

set +e
# Check if the scap-security-guide package is available in the enabled repositories
dnf provides scap-security-guide

# Check the exit status of the previous command
if [ $? -ne 0 ]; then
    echo "No working repository is available to install scap-security-guide."

    # Check if RHEL_APPS_REPO variable is set
    if [ -z "$RHEL_APPS_REPO" ]; then
        echo "Error: RHEL_APPS_REPO is not set. Please set the variable and try again."
        exit 1
    else
        echo "Setting up CaC Apps repository..."
        cat > /etc/yum.repos.d/cac.repo <<EOF
[cac_apps_repo]
name=CaC Apps Repo
baseurl=${RHEL_APPS_REPO}
enabled=1
gpgcheck=0
EOF
        echo "CaC Apps repository has been added."
    fi
fi

echo "Installing dependencies..."
dnf update -y
dnf install git wget make scap-security-guide -y
rm -rf /usr/bin/go
go_mod="https://raw.githubusercontent.com/complytime/complytime/main/go.mod"
go_version=$(curl -s $go_mod | grep '^go' | awk '{print $2}')
go_tar_file=go$go_version.linux-amd64.tar.gz
wget https://go.dev/dl/$go_tar_file
tar -C /usr/local -xvzf $go_tar_file
rm -rf $go_tar_file
export PATH=$PATH:/usr/local/go/bin
source ~/.bash_profile

# Install and build complytime
echo "Cloning the Complytime repository..."
git clone https://github.com/complytime/complytime.git
cd complytime && make build && cp ./bin/complytime /usr/local/bin
echo "Complytime installed successfully!"
# Run complytime list to create the workspace
set +e
# Running list command that will fail due to no requirements files
echo "Attempting to run the command complytime list."
bin/complytime list 2>/dev/null
echo "An error occurred, but script continues after the complytime list."
# Copy the artifacts to workspace
cp docs/samples/sample-component-definition.json ~/.config/complytime/bundles
cp docs/samples/sample-profile.json ~/.config/complytime/controls

# Copy the plugins' files
cp -rp bin/openscap-plugin ~/.config/complytime/plugins
cp -rp cmd/openscap-plugin/openscap-plugin.yml ~/.config/complytime/plugins
checksum=$(sha256sum ~/.config/complytime/plugins/openscap-plugin| cut -d ' ' -f 1 )
cat > ~/.config/complytime/plugins/c2p-openscap-manifest.json << EOF
{
  "metadata": {
    "id": "openscap",
    "description": "My openscap plugin",
    "version": "0.0.1",
    "types": ["pvp"]
  },
  "executablePath": "openscap-plugin",
  "sha256": "$checksum"
}
EOF
