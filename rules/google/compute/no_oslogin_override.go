package compute

import (
	"github.com/aquasecurity/defsec/providers"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/severity"
	"github.com/aquasecurity/defsec/state"
)

var CheckNoOsloginOverride = rules.Register(
	rules.Rule{
		AVDID:       "AVD-GCP-0036",
		Provider:    providers.GoogleProvider,
		Service:     "compute",
		ShortCode:   "no-oslogin-override",
		Summary:     "Instances should not override the project setting for OS Login",
		Impact:      "Access via SSH key cannot be revoked automatically when an IAM user is removed.",
		Resolution:  "Enable OS Login at project level and remove instance-level overrides",
		Explanation: `OS Login automatically revokes the relevant SSH keys when an IAM user has their access revoked.`,
		Links:       []string{},
		Terraform: &rules.EngineMetadata{
			GoodExamples:        terraformNoOsloginOverrideGoodExamples,
			BadExamples:         terraformNoOsloginOverrideBadExamples,
			Links:               terraformNoOsloginOverrideLinks,
			RemediationMarkdown: terraformNoOsloginOverrideRemediationMarkdown,
		},
		Severity: severity.Medium,
	},
	func(s *state.State) (results rules.Results) {
		for _, instance := range s.Google.Compute.Instances {
			if instance.IsUnmanaged() {
				continue
			}
			if instance.OSLoginEnabled.IsFalse() {
				results.Add(
					"Instance has OS Login disabled.",
					instance.OSLoginEnabled,
				)
			} else {
				results.AddPassed(&instance)
			}
		}
		return
	},
)
