from core.utils import check_adb_devices, get_device_model, get_android_version, is_frp_locked, check_bootloader_status, detect_brand

def main():
    if not check_adb_devices():
        print("[!] No ADB devices detected.")
        return

    model = get_device_model()
    brand = detect_brand(model)
    android_version = get_android_version()
    frp_status = is_frp_locked()
    bootloader_status = check_bootloader_status()

    print("\n=== Device Information ===")
    print(f"Model: {model}")
    print(f"Brand: {brand}")
    print(f"Android Version: {android_version}")
    print(f"FRP Status: {'ENABLED' if frp_status else 'DISABLED'}")
    print(f"Bootloader: {'UNLOCKED' if bootloader_status else 'LOCKED'}")

if __name__ == "__main__":
    main()