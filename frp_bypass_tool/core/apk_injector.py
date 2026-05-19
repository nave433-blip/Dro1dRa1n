import os
from .utils import run_adb_command, get_device_model

def inject_frp_bypass_apk():
    """Sideload a universal FRP bypass APK."""
    apk_path = os.path.join(os.path.dirname(__file__), "../payloads/frp_bypass.apk")
    if not os.path.exists(apk_path):
        print("[-] frp_bypass.apk not found in payloads/")
        return False

    print("[*] Installing FRP bypass APK...")
    install_output = run_adb_command(f"install {apk_path}")
    if "Success" in install_output:
        print("[+] APK installed successfully.")
        run_adb_command("shell am start -n com.frp.bypass/.MainActivity")
        print("[+] Launched FRP bypass app. Follow on-screen instructions.")
        return True
    else:
        print("[-] APK installation failed.")
        return False