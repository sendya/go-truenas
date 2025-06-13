#!/bin/bash
set -e

sudo apt-get update
sudo apt-get install -y \
    git \
    libglib2.0-dev \
    libfdt-dev \
    libpixman-1-dev \
    libslirp-dev \
    zlib1g-dev \
    ninja-build \
    build-essential

QEMU_VERSION="v10.0.2"
BUILD_DIR="$HOME/qemu-build"
INSTALL_PREFIX="/usr/local"

echo "Building QEMU ${QEMU_VERSION}..."

# Clone QEMU source
git clone --depth 1 --branch ${QEMU_VERSION} https://gitlab.com/qemu-project/qemu.git ~/qemu-src
cd ~/qemu-src

# Configure for minimal build - only x86_64 system emulator
./configure \
    --target-list=x86_64-softmmu \
    --prefix=${INSTALL_PREFIX} \
    --enable-slirp

# Build with all available cores
make -j$(nproc)

# Install to temporary directory for caching
mkdir -p ${BUILD_DIR}/bin ${BUILD_DIR}/share
make install DESTDIR=${BUILD_DIR}

echo "QEMU build completed successfully"
