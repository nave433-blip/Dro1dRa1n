# 🔓 Android FRP Bypass Tool (Open-Source)
**Factory Reset Protection Removal for Research & Recovery Purposes**

⚠️ **DISCLAIMER:**
This tool is for **educational and device recovery purposes only**. Bypassing FRP may violate **Google's Terms of Service** and local laws. Use at your own risk.

## ✅ Features
- **ADB-based FRP removal** (for devices with USB debugging enabled)
- **Fastboot/OEM unlock methods** (for unlocked bootloaders)
- **Brand-specific exploits** (Samsung, Xiaomi, Google Pixel, etc.)
- **APK injection** (for devices with locked ADB)
- **Automated device detection** (model, Android version, FRP status)

## 🛠️ Requirements
- **Python 3.8+**
- **ADB & Fastboot** (Android SDK Platform Tools)
- **USB Debugging enabled** (for ADB methods)
- **Unlocked Bootloader** (for Fastboot methods)

## 🚀 Installation
```bash
git clone https://github.com/yourusername/frp_bypass_tool.git
cd frp_bypass_tool
pip install -r requirements.txt
```

## 🎯 Usage
```bash
python scripts/auto_frp_bypass.py
```
Follow the interactive menu to select your device and method.

## 📂 Directory Structure
- `core/` – Core bypass logic
- `exploits/` – Brand-specific exploit payloads
- `scripts/` – Standalone execution scripts
- `payloads/` – Prebuilt APKs and binaries

## 🔧 Supported Devices
| Brand       | Method               | Status       |
|-------------|----------------------|--------------|
| Samsung     | ADB, Fastboot, APK   | ✅ Tested    |
| Xiaomi      | Fastboot, APK        | ✅ Tested    |
| Google Pixel| ADB, Fastboot        | ✅ Tested    |
| OnePlus     | Fastboot             | ⚠️ Partial   |
| Huawei      | APK Injection         | ❌ Unstable  |

## 📜 Legal Note
This tool is **not** intended for illegal use. Only use on **your own devices** or with **explicit permission**.