# Automate Quick Start Guide for Complyctl

This quick start guide aims to set up the minimum environment to explore complyctl (the ComplyTime CLI).
If you intend to run complyctl on a real target system, the content will need to be updated.
Below are two options to quickly start using complyctl.

## Option 1: Run the quick_start.sh in ubi9 container
Assumed that you have installed [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) and [Podman](https://podman.io/docs/installation). If your system is Fedora, CentOS or RHEL you could install them directly:

```bash
dnf install git podman -y
```

1. Pull the complytl Git repository to your local system:

```bash
git clone https://github.com/complytime/complyctl.git
cd complyctl/scripts/quick_start
```

2. Set the RHEL_APPS_REPO variable to the app DNF repository URL:

```bash
export RHEL_APPS_REPO=${RHEL_APPS_REPO} # Change the ${RHEL_APPS_REPO} to the app dnf repo
```

3. Build the container image:

```bash
podman build --build-arg RHEL_APPS_REPO=${RHEL_APPS_REPO} . -t quick-start:latest
podman run -it quick-start:latest /bin/bash
```

## Option 2: Run the quick_start.sh in a fresh RHEL
Assume that you have already installed a fresh RHEL
1. Download or copy the [quick_start.sh](quick_start.sh) script in your fresh RHEL system.

2. Run the script on your fresh RHEL system:

```bash
chmod +x quick_start.sh
export RHEL_APPS_REPO=$RHEL_APPS_REPO # Change the ${RHEL_APPS_REPO} to the app dnf repo
sh quick_start.sh
```

Now you can explore complyctl commands:
```bash
complyctl list
complyctl plan anssi_bp28_minimal
complyctl generate
complyctl scan
```
