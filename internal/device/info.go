// Package device collects structured device information via adb.
package device

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nave433-blip/Dro1dRa1n/internal/adb"
)

// Info is a structured dump of device properties relevant to repair work.
type Info struct {
	Serial        string
	Brand         string
	Manufacturer  string
	Model         string
	DeviceName    string
	Product       string
	Android       string
	SDK           string
	SecurityPatch string
	BuildID       string
	Fingerprint   string
	Bootloader    string // ro.boot.bootloader
	VerifiedBoot  string // ro.boot.verifiedbootstate
	FRPPst        string // ro.frp.pst (when exposed)
	OEMUnlock     string // sys.oem_unlock_allowed (1 = unlocking allowed)
	KnoxWarranty  string // ro.boot.warranty_bit (0 = Knox intact, 1 = tripped)
	AntiRollback  string // ro.build.version.anti (Xiaomi anti-rollback index)
	FlashLocked   string // ro.boot.flash.locked (1 = bootloader locked)
}

// Collect gathers Info for the given serial using an adb Runner.
func Collect(r *adb.Runner, serial string) (*Info, error) {
	out, err := r.Run(serial, "shell", "getprop")
	if err != nil {
		return nil, fmt.Errorf("getprop failed: %v", err)
	}
	p := adb.ParseGetprop(out)
	return &Info{
		Serial:        serial,
		Brand:         p["ro.product.brand"],
		Manufacturer:  p["ro.product.manufacturer"],
		Model:         p["ro.product.model"],
		DeviceName:    p["ro.product.device"],
		Product:       p["ro.product.name"],
		Android:       p["ro.build.version.release"],
		SDK:           p["ro.build.version.sdk"],
		SecurityPatch: p["ro.build.version.security_patch"],
		BuildID:       p["ro.build.display.id"],
		Fingerprint:   p["ro.build.fingerprint"],
		Bootloader:    p["ro.boot.bootloader"],
		VerifiedBoot:  p["ro.boot.verifiedbootstate"],
		FRPPst:        p["ro.frp.pst"],
		OEMUnlock:     p["sys.oem_unlock_allowed"],
		KnoxWarranty:  p["ro.boot.warranty_bit"],
		AntiRollback:  p["ro.build.version.anti"],
		FlashLocked:   p["ro.boot.flash.locked"],
	}, nil
}

// BrandKey normalizes the brand for method matching.
func (i *Info) BrandKey() string {
	return strings.ToLower(strings.TrimSpace(i.Brand))
}

// AndroidMajor parses the major version from e.g. "13" or "12.0"; 0 = unknown.
func (i *Info) AndroidMajor() int {
	rel := strings.TrimSpace(i.Android)
	if rel == "" {
		return 0
	}
	v, err := strconv.Atoi(strings.SplitN(rel, ".", 2)[0])
	if err != nil {
		return 0
	}
	return v
}

// Summary renders a human-readable block.
func (i *Info) Summary() string {
	var b strings.Builder
	row := func(k, v string) {
		if v == "" {
			v = "(unknown)"
		}
		fmt.Fprintf(&b, "  %-15s %s\n", k+":", v)
	}
	row("serial", i.Serial)
	row("brand", i.Brand)
	row("model", i.Model)
	row("device", i.DeviceName)
	android := i.Android
	if android != "" && i.SDK != "" {
		android += " (SDK " + i.SDK + ")"
	} else if i.SDK != "" {
		android = "SDK " + i.SDK
	}
	row("android", android)
	row("sec. patch", i.SecurityPatch)
	row("build", i.BuildID)
	row("bootloader", i.Bootloader)
	row("verified boot", i.VerifiedBoot)
	if i.OEMUnlock != "" {
		v := i.OEMUnlock
		if v == "1" {
			v = "1 (OEM unlocking allowed)"
		} else if v == "0" {
			v = "0 (OEM unlocking NOT allowed)"
		}
		row("oem unlock", v)
	}
	if i.FlashLocked != "" {
		v := i.FlashLocked
		if v == "1" {
			v = "1 (bootloader LOCKED)"
		} else if v == "0" {
			v = "0 (bootloader UNLOCKED)"
		}
		row("flash lock", v)
	}
	if i.KnoxWarranty != "" {
		v := i.KnoxWarranty
		if v == "0" {
			v = "0 (Knox intact)"
		} else if v == "1" {
			v = "1 (Knox TRIPPED — warranty void)"
		}
		row("knox warranty", v)
	}
	if i.AntiRollback != "" {
		row("anti-rollback", i.AntiRollback)
	}
	return b.String()
}
