import os
import subprocess
import re
from pyadb import ADB

def run_adb_command(command, device_serial=None):
    """Execute an ADB command and return output."""
    try:
        adb = ADB()
        result = adb.run_cmd(command, device_serial)
        return result.stdout.decode().strip()
    except Exception as e:
        print(f"[!] ADB Error: {e}")
        return None

def check_adb_devices():
    """Check if any devices are connected via ADB."""
    devices = run_adb_command("devices")
    if "no devices" in devices.lower():
        return False
    return True

def get_device_model():
    """Get the connected device model."""
    model = run_adb_command("shell getprop ro.product.model")
    return model if model else "Unknown"

def get_android_version():
    """Get the Android version of the connected device."""
    version = run_adb_command("shell getprop ro.build.version.release")
    return version if version else "Unknown"

def is_frp_locked():
    """Check if FRP is enabled on the device."""
    frp_status = run_adb_command("shell getprop ro.frp.pst")
    return "1" in frp_status if frp_status else False

def check_bootloader_status():
    """Check if bootloader is unlocked."""
    status = run_adb_command("shell getprop ro.boot.verifiedbootstate")
    return "orange" in status.lower() if status else False

def detect_brand(model):
    """Detect device brand from model name."""
    model = model.lower()
    if "samsung" in model or "sm-" in model:
        return "samsung"
    elif "xiaomi" in model or "redmi" in model or "mi " in model:
        return "xiaomi"
    elif "pixel" in model or "google" in model:
        return "google"
    elif "oneplus" in model:
        return "oneplus"
    elif "huawei" in model:
        return "huawei"
    else:
        return "generic"