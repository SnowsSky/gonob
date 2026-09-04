package config

import (
	"gonob/translations"
	"os"
	"os/exec"

	scolor "github.com/SnowsSky/scolor/pkg"
)

func DetectPrivilegeEscalationTool() string {
	if _, err := exec.LookPath("sudo"); err == nil {
		return "sudo"
	}
	if _, err := exec.LookPath("sudo-rs"); err == nil {
		return "sudo"
	}
	if _, err := exec.LookPath("doas"); err == nil {
		return "doas"
	}
	scolor.BoldRed.DisplayText("==> " + translations.Translate("error_string") + " ")
	scolor.BoldWhite.DisplayText(" " + translations.Translate("unknown_privilege_escalation_tool") + "\n")
	os.Exit(1)
	return "unknown"
}

var PrivilegeEscalationTool = DetectPrivilegeEscalationTool()
