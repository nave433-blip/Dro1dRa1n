import subprocess
from .utils import check_bootloader_status, get_device_model, detect_brand

def fastboot_remove_frp():
    """Attempt to remove FRP via Fastboot (requires unlocked bootloader)."""
    print("[*] Checking bootloader status...")
    if not check_bootloader_status():
        print("[-] Bootloader is locked. Fastboot FRP bypass may not work.")
        return False

    print("[*] Attempting Fastboot FRP bypass...")
    model = get_device_model()
    brand = detect_brand(model)

    # Universal Fastboot commands
    fastboot_commands = [
        "fastboot oem unlock",
        "fastboot flashing unlock",
        "fastboot erase frp",
        "fastboot erase config",
        "fastboot reboot"
    ]

    for cmd in fastboot_commands:
        print(f"[*] Executing: {cmd}")
        try:
            subprocess.run(cmd.split(), check=True)
        except subprocess.CalledProcessError as e:
            print(f"[-] Failed: {cmd} ({e})")

    # Brand-specific Fastboot commands
    if brand == "samsung":
        print("[*] Applying Samsung-specific Fastboot bypass...")
        samsung_commands = [
            "fastboot flash frp frp.img",  # Requires a dummy frp.img
            "fastboot oem frp-unlock"
        ]
        for cmd in samsung_commands:
            try:
                subprocess.run(cmd.split(), check=True)
            except:
                pass

    elif brand == "xiaomi":
        print("[*] Applying Xiaomi-specific Fastboot bypass...")
        xiaomi_commands = [
            "fastboot oem frp-reset",
            "fastboot oem unlock-go"
        ]
        for cmd in xiaomi_commands:
            try:
                subprocess.run(cmd.split(), check=True)
            except:
                pass

    print("[+] Fastboot FRP bypass attempt complete. Rebooting device...")
    subprocess.run(["fastboot", "reboot"], check=True)
    return True