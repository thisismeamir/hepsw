# HEP SW initialization script
# ==== What will this script do? ====
# This script is a general environment setup before
# starting the build process, or environment setups.
# 
# It is the critical starting point of hepsw project.
# 
# ==== Usage ====
# just source this shell script and informations would pop-up.



echo "======================"
echo "= ╻ ╻┏━╸┏━┓   ┏━┓╻ ╻ ="
echo "= ┣━┫┣╸ ┣━┛   ┗━┓┃╻┃ ="
echo "= ╹ ╹┗━╸╹     ┗━┛┗┻┛ ="
echo "= version      0.0.0 ="
echo "======================"                                                           

echo "[INFO] Initializing HEP SW"
echo ""
echo "======== BASIC PATHS AND ENV VARIABLES ========"
export HEPSW=$(pwd)
export HEPSW_SOURCES=$HEPSW/sources
export HEPSW_BUILDS=$HEPSW/builds
export HEPSW_ISNTALL=$HEPSW/install
echo "[INFO] Setting home directory to: $HEPSW"
echo "[INFO] Setting sources directory to: $HEPSW_SOURCES"
echo "[INFO] Setting builds directory to $HEPSW_BUILDS"
echo "[INFO] Setting install directory to $HEPSW_INSTALL"

echo ""
echo "======== SYSTEM INFORMATION ========"
export OS_NAME=$(uname -s)
export NODE_NAME=$(uname -n)
export KERNEL_RELEASE=$(uname -r)
export KERNEL_VERSION=$(uname -r)
export ARCHITECTURE=$(uname -m)
echo "[INFO] OS detected: $OS_NAME"
echo "[INFO] Distro: $NODE_NAME"
echo "[INFO] Kernel: $KERNEL_RELEASE $KERNEL_VERSION"
echo "[INFO] Architecture: $ARCHITECTURE"

