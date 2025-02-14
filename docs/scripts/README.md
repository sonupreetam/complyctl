# Automate Quick Start Guide for ComplyTime CLI

This quick start guide aims to set up the minimum environment to explore the ComplyTime CLI.
If you intend to run the ComplyTime CLI on a real target system, the content will need to be updated.
Below are two options to quickly start using the ComplyTime CLI.

## Option 1: Run the quick_start.sh in ubi9 container
1. Pull the ComplyTime Git repository to your local system.
Note: It’s assumed that you have [Podman installed](https://podman.io/docs/installation) on your system.
```bash
git clone https://github.com/complytime/complytime.git
cd complytime/docs/scripts/
```
2. Set the APPS_REPO variable to the app Yum repository URL:
```bash
export APPS_REPO=${APPS_REPO} # Change the ${APPS_REPO} to the app yum repo
```
3. Build the container image:
```bash
podman build --build-arg APPS_REPO=${APPS_REPO} . -t quick-start:latest
podman run -it quick-start:latest /bin/bash
```
## Option 2: Run the quick_start.sh in a fresh RHEL
Assume that you have already installed a fresh RHEL
1. Download or copy the quick_start.sh script in your refresh RHEL system.
2. Run the script on your refresh RHEL system:
```bash
chmod +x quick_start.sh
export APPS_REPO=$APPS_REPO # Change the ${APPS_REPO} to the app yum repo
sh quick_start.sh
```


Now you could explore Complytime CLIs:
```bash
complytime list
complytime plan example
complytime generate
complytime scan
```
