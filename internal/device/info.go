// Package device collects structured device information via adb.
package device

import (
	"fmt"
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
	}, nil
}

// BrandKey normalizes the brand for method matching.
func (i *Info) BrandKey() string {
	return strings.ToLower(strings.TrimSpace(i.Brand))
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
	row("android", i.Android+" (SDK "+i.SDK+")")
	row("sec. patch", i.SecurityPatch)
	row("build", i.BuildID)
	row("bootloader", i.Bootloader)
	row("verified boot", i.VerifiedBoot)
	return b.String()
}
