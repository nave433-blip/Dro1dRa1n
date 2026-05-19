#!/bin/bash
# Emergency FRP Bypass for Rooted Devices
echo "[*] Checking for root..."
if ! su -c "id" | grep -q "uid=0"; then
    echo "[-] Device is not rooted. Aborting."
    exit 1
fi

echo "[*] Removing FRP via root..."
su -c "rm -f /data/system/frp/*.xml"
su -c "rm -f /data/system/users/0/frp/*.xml"
su -c "settings put global device_provisioned 1"
su -c "settings put global user_setup_complete 1"
su -c "pm clear com.google.android.gsf.login"
su -c "reboot"

echo "[+] FRP removed. Rebooting..."