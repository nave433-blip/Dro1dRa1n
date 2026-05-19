import os
import sys
import time
from colorama import Fore, Style, init
from core.utils import check_adb_devices, get_device_model, is_frp_locked, detect_brand
from core.adb_bypass import adb_remove_frp
from core.fastboot_bypass import fastboot_remove_frp
from core.exploit_bypass import apply_exploit
from core.apk_injector import inject_frp_bypass_apk
from core.exploit_fetcher import GitHubExploitFetcher

# Initialize colorama
init(autoreset=True)
fetcher = GitHubExploitFetcher()

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

def print_menu():
    menu_items = [
        "Check FRP Status",
        "ADB FRP Bypass (USB Debug)",
        "Fastboot FRP Bypass",
        "Brand-Specific Exploit",
        "APK Injection",
        "Update Exploits from GitHub",
        "Exit Tool"
    ]
    
    print(Fore.BLUE + "    ╔" + "═" * 60 + "╗")
    print(Fore.BLUE + "    ║" + Fore.YELLOW + Style.BRIGHT + " FRP BYPASS OPERATIONS ".center(60) + Fore.BLUE + "║")
    print(Fore.BLUE + "    ╠" + "═" * 4 + "╦" + "═" * 55 + "╣")
    
    for i, desc in enumerate(menu_items, 1):
        num = Fore.GREEN + f"{i}"
        description = Fore.WHITE + f"{desc.ljust(53)}"
        print(Fore.BLUE + f"    ║ {num} ║ {description} ║")
    
    print(Fore.BLUE + "    ╚" + "═" * 4 + "╩" + "═" * 55 + "╝")

def main():
    # Start auto-updater in background
    fetcher.start_auto_updater()
    
    while True:
        os.system('cls' if os.name == 'nt' else 'clear')
        print_logo()
        
        if not check_adb_devices():
            print(Fore.RED + "    [!] No ADB devices detected. Connect device and enable USB debugging.")
        else:
            model = get_device_model()
            brand = detect_brand(model)
            print(Fore.GREEN + f"    [*] Connected: {model} ({brand})")

        print_menu()
        choice = input(Fore.YELLOW + Style.BRIGHT + "\n    FRP > " + Style.RESET_ALL).strip()

        if choice == "1":
            if is_frp_locked():
                print(Fore.RED + "    [-] FRP is ENABLED on this device.")
            else:
                print(Fore.GREEN + "    [+] FRP is DISABLED on this device.")

        elif choice == "2":
            print(Fore.BLUE + "\n    [*] Starting ADB FRP Bypass...")
            adb_remove_frp()

        elif choice == "3":
            print(Fore.BLUE + "\n    [*] Starting Fastboot FRP Bypass...")
            fastboot_remove_frp()

        elif choice == "4":
            model = get_device_model()
            brand = detect_brand(model)
            print(Fore.BLUE + f"\n    [*] Applying {brand.capitalize()} FRP Exploit...")
            apply_exploit(brand)

        elif choice == "5":
            print(Fore.BLUE + "\n    [*] Starting APK Injection...")
            inject_frp_bypass_apk()

        elif choice == "6":
            print(Fore.BLUE + "\n    [*] Updating exploits from GitHub...")
            exploits = fetcher.fetch_exploits()
            print(Fore.GREEN + f"    [+] Loaded {len(exploits)} dynamic exploits.")

        elif choice == "7" or choice.lower() == 'q':
            break

        else:
            print(Fore.RED + "    [-] Invalid choice. Try again.")

        time.sleep(2)

if __name__ == "__main__":
    main()