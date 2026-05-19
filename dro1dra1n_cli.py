import os
import sys
import time
import subprocess
import shutil
from colorama import Fore, Style, init

# Initialize colorama
init(autoreset=True)

BASE_DIR = os.path.dirname(os.path.abspath(__file__))

def clear_screen():
    os.system('cls' if os.name == 'nt' else 'clear')

def check_dependencies():
    deps = {
        "adb": shutil.which("adb"),
        "fastboot": shutil.which("fastboot"),
        "python3": shutil.which("python3"),
        "git": shutil.which("git"),
    }
    
    print(Fore.CYAN + "\n    [#] SYSTEM DEPENDENCY CHECK")
    print(Fore.BLUE + Style.DIM + "    " + "-" * 40)
    for name, path in deps.items():
        status = Fore.GREEN + "[OK]" if path else Fore.RED + "[MISSING]"
        print(f"    {name.ljust(10)} : {status} {Fore.WHITE + '(' + str(path) + ')' if path else ''}")
    print(Fore.BLUE + Style.DIM + "    " + "-" * 40 + "\n")
    return all(deps.values())

def get_usb_devices():
    print(Fore.CYAN + "    [#] CONNECTED USB DEVICES (ANDROID DETECTOR)")
    print(Fore.BLUE + Style.DIM + "    " + "-" * 40)
    
    try:
        if sys.platform == "darwin":
            cmd = "system_profiler SPUSBDataType"
            output = subprocess.check_output(cmd, shell=True).decode()
            devices = []
            current_device = None
            for line in output.split('\n'):
                if line.strip() and not line.startswith(' '):
                    current_device = line.strip()
                if "Vendor ID:" in line or "Manufacturer:" in line or "Serial Number:" in line:
                    if current_device:
                        devices.append(f"{current_device} | {line.strip()}")
                        current_device = None
            
            for dev in devices[:10]:
                if any(x in dev for x in ["Google", "Samsung", "Motorola", "Xiaomi"]):
                    print(Fore.GREEN + f"    [!] {dev}")
                else:
                    print(Fore.WHITE + f"    [-] {dev}")
        else:
            cmd = "lsusb"
            output = subprocess.check_output(cmd, shell=True).decode()
            for line in output.split('\n'):
                if line.strip():
                    if any(x in line for x in ["Google", "Samsung", "Motorola", "Xiaomi", "ADB"]):
                        print(Fore.GREEN + f"    [!] {line.strip()}")
                    else:
                        print(Fore.WHITE + f"    [-] {line.strip()}")
    except Exception as e:
        print(Fore.RED + f"    [!] Error reading USB: {e}")
    print(Fore.BLUE + Style.DIM + "    " + "-" * 40 + "\n")

def print_logo():
    # Green Droid with Umbrella in the Rain
    print(Fore.GREEN + r"""
              .        .        .        .
          .        .        .        .        .
        .      .        _________        .      .
            .      .   /         \   .      .
          .      .    /           \    .      .
        .      .     /_____________\     .      .
            .      .        | |        .      .
                    _____   | |   _____
                   /     \  | |  /     \
                  | () () | | | | () () |
                   \__^__/  | |  \__^__/
                  /|     |  | |  |     |\
                 / |     |  | |  |     | \
                   |_____|  | |  |_____|
                   /     \  |_|  /     \
                  /       \     /       \
    """)
    
    # Dro1dRa1n in Big Stylized Letters
    print(Fore.CYAN + Style.BRIGHT + r"""
      _____             _     _ ____         _           
     |  __ \           (_)   | |  _ \       (_)          
     | |  | |_ __ ___   _  __| | |_) | __ _ _ _ __      
     | |  | | '__/ _ \ | |/ _` |  _ < / _` | | '_ \     
     | |__| | | | (_) || | (_| | |_) | (_| | | | | |    
     |_____/|_|  \___/ |_|\__,_|____/ \__,_|_|_| |_|    
                                                         
    """)
    print(Fore.BLUE + Style.DIM + "      " + "~" * 50)
    print(Fore.WHITE + "        [ THE ULTIMATE ANDROID BYPASS & HACKING SUITE ]")
    print(Fore.BLUE + Style.DIM + "      " + "~" * 50 + "\n")

def print_fancy_menu():
    menu_items = [
        ("MODE", "Device Mode Manager"),
        ("FRP", "FRP Bypass Tool"),
        ("ROOT", "Rooting Assistant"),
        ("MTK", "MediaTek Toolkit"),
        ("BTLO", "Bootloader Unlock"),
        ("BBRC", "Baseband & Recovery"),
        ("ROM", "Custom ROM Tools"),
        ("RE", "Reverse Engineering"),
        ("HW", "Hardware Reference"),
        ("UPDT", "Update Suite"),
        ("EXIT", "Exit Suite")
    ]
    
    print(Fore.BLUE + "    ╔" + "═" * 60 + "╗")
    print(Fore.BLUE + "    ║" + Fore.YELLOW + Style.BRIGHT + " SELECT AN OPERATION MODE ".center(60) + Fore.BLUE + "║")
    print(Fore.BLUE + "    ╠" + "═" * 15 + "╦" + "═" * 44 + "╣")
    
    for i, (cmd, desc) in enumerate(menu_items, 1):
        num = Fore.GREEN + f"{str(i).rjust(2)}"
        command = Fore.CYAN + f"{cmd.ljust(10)}"
        description = Fore.WHITE + f"{desc.ljust(40)}"
        print(Fore.BLUE + f"    ║  {num}  {command} ║  {description}  ║")
    
    print(Fore.BLUE + "    ╚" + "═" * 15 + "╩" + "═" * 44 + "╝")

def mode_manager():
    clear_screen()
    print_logo()
    print(Fore.CYAN + Style.BRIGHT + "\n    [#] DEVICE MODE MANAGER")
    print(Fore.BLUE + Style.DIM + "    " + "=" * 60)
    print(Fore.WHITE + "    Select a target mode to reboot the connected device:")
    
    modes = [
        ("Normal", "reboot"),
        ("Recovery", "reboot recovery"),
        ("Bootloader / Fastboot", "reboot bootloader"),
        ("Download (Samsung)", "reboot download"),
        ("EDL (Qualcomm Emergency)", "reboot edl"),
        ("Sideload (ADB)", "reboot sideload")
    ]
    
    for i, (name, cmd) in enumerate(modes, 1):
        print(f"    {Fore.GREEN}[{i}] {Fore.CYAN}{name.ljust(25)} {Fore.WHITE}| {Fore.YELLOW}adb {cmd}")
    
    print(f"    {Fore.GREEN}[B] {Fore.WHITE}Back to Main Menu")
    
    choice = input(Fore.YELLOW + "\n    Mode Selection > ").strip().lower()
    
    if choice == 'b':
        return
    
    try:
        idx = int(choice) - 1
        if 0 <= idx < len(modes):
            name, cmd = modes[idx]
            print(Fore.GREEN + f"\n    [*] Attempting to reboot into {name}...")
            full_cmd = f"adb {cmd}"
            subprocess.run(full_cmd.split(), check=True)
            print(Fore.GREEN + "    [+] Command sent successfully.")
        else:
            print(Fore.RED + "    [!] Invalid choice.")
    except Exception as e:
        print(Fore.RED + f"    [!] Error: {e}")
    
    time.sleep(2)

def baseband_recovery_tools():
    clear_screen()
    print_logo()
    print(Fore.CYAN + Style.BRIGHT + "\n    [#] BASEBAND & CUSTOM RECOVERY TOOLS")
    print(Fore.BLUE + Style.DIM + "    " + "=" * 60)
    
    print(Fore.GREEN + "    [1] Flash TWRP / Custom Recovery")
    print(Fore.GREEN + "    [2] Backup Baseband / EFS (Root Required)")
    print(Fore.GREEN + "    [3] Restore Baseband / EFS")
    print(Fore.GREEN + "    [4] Flash Modem / Radio Firmware")
    print(Fore.GREEN + "    [B] Back")
    
    choice = input(Fore.YELLOW + "\n    Selection > ").strip().lower()
    
    if choice == '1':
        img = input(Fore.WHITE + "    Enter path to recovery .img file: ").strip()
        if os.path.exists(img):
            print(Fore.YELLOW + "    [!] Ensure device is in Bootloader mode.")
            input("    Press Enter to flash...")
            os.system(f"fastboot flash recovery {img}")
        else:
            print(Fore.RED + "    [-] File not found.")
    elif choice == '2':
        print(Fore.YELLOW + "    [*] Pulling EFS partition...")
        os.system("adb shell su -c 'dd if=/dev/block/bootdevice/by-name/efs of=/sdcard/efs.img'")
        os.system("adb pull /sdcard/efs.img ./efs_backup.img")
    elif choice == '4':
        img = input(Fore.WHITE + "    Enter path to modem/radio .img or .bin: ").strip()
        if os.path.exists(img):
            print(Fore.YELLOW + "    [!] Ensure device is in Bootloader/Fastboot mode.")
            part = input(Fore.WHITE + "    Enter partition name (modem/radio): ").strip()
            os.system(f"fastboot flash {part} {img}")
    
    time.sleep(2)

def show_hardware_info():
    clear_screen()
    print_logo()
    print(Fore.CYAN + Style.BRIGHT + "\n    [#] ANDROID HARDWARE & SOC REFERENCE")
    print(Fore.BLUE + Style.DIM + "    " + "=" * 60)
    
    info = {
        "Qualcomm Snapdragon": "ARMv8/v9, Kryo Cores, Adreno GPU. Uses EDL mode (9008) for low-level flashing.",
        "MediaTek (MTK)": "ARM Cortex cores, Mali/PowerVR GPU. Uses BROM mode for exploits (mtkclient).",
        "Samsung Exynos": "Custom Mongoose/Cortex cores. Uses Download Mode (Odin) for flashing.",
        "Unisoc (Spreadtrum)": "Budget SoCs. Uses ResearchDownload tool. Vulnerable to specific bootloader exploits.",
        "Google Tensor": "Exynos-based custom SoC for Pixel 6+. Strong Titan M2 security.",
        "ARM Architectures": "v7 (32-bit legacy), v8 (64-bit standard), v9 (Modern high-security)."
    }
    
    for soc, desc in info.items():
        print(Fore.GREEN + f"\n    ▶ {soc}")
        print(Fore.WHITE + f"      {desc}")
    
    input("\n    Press Enter to return...")

def auto_update():
    print(Fore.CYAN + "\n    [*] CHECKING FOR DRO1DRA1N UPDATES...")
    if shutil.which("git"):
        try:
            output = subprocess.check_output("git pull", shell=True).decode()
            if "Already up to date" in output:
                print(Fore.GREEN + "    [+] Dro1dRa1n is already at the latest version.")
            else:
                print(Fore.YELLOW + "    [!] Updates downloaded! Please restart Dro1dRa1n.")
                print(output)
        except Exception as e:
            print(Fore.RED + f"    [-] Git update failed: {e}")
    else:
        print(Fore.RED + "    [-] Git not installed.")
    time.sleep(2)

def main():
    while True:
        clear_screen()
        print_logo()
        check_dependencies()
        get_usb_devices()
        print_fancy_menu()
        
        choice = input(Fore.YELLOW + Style.BRIGHT + "\n    Dro1dRa1n > " + Style.RESET_ALL).strip().lower()
        
        if choice in ['1', 'mode']:
            mode_manager()
        elif choice in ['2', 'frp']:
            os.system(f"python3 {os.path.join(BASE_DIR, 'frp_bypass_tool/scripts/auto_frp_bypass.py')}")
        elif choice in ['3', 'root']:
            os.system(f"python3 {os.path.join(BASE_DIR, 'android_root_helper/android_root_helper.py')}")
        elif choice in ['4', 'mtk']:
            os.system(f"python3 {os.path.join(BASE_DIR, 'mtkclient/mtk_gui.py')}")
        elif choice in ['5', 'btlo']:
            os.system(f"bash {os.path.join(BASE_DIR, 'Bruteforce-Bootloader-Unlocker/bootloader_unlocker')}")
        elif choice in ['6', 'bbrc']:
            baseband_recovery_tools()
        elif choice in ['7', 'rom']:
            os.system(f"bash {os.path.join(BASE_DIR, 'RomTools/RomTools.sh')}")
        elif choice in ['8', 're']:
            os.system("objection --help")
            input("\n    Press Enter to return...")
        elif choice in ['9', 'hw']:
            show_hardware_info()
        elif choice in ['10', 'updt']:
            auto_update()
            sys.path.append(os.path.join(BASE_DIR, 'frp_bypass_tool'))
            try:
                from core.exploit_fetcher import GitHubExploitFetcher
                fetcher = GitHubExploitFetcher()
                exploits = fetcher.fetch_exploits()
                print(Fore.GREEN + f"    [+] Dynamic exploits updated: {len(exploits)} found.")
            except ImportError:
                print(Fore.RED + "    [-] Exploit fetcher core not found.")
            time.sleep(2)
        elif choice in ['11', 'exit', 'q']:
            print(Fore.YELLOW + "\n    [!] Shutting down Dro1dRa1n. Happy hunting!")
            sys.exit(0)
        else:
            print(Fore.RED + "\n    [!] Invalid Selection.")
            time.sleep(1)

if __name__ == "__main__":
    main()