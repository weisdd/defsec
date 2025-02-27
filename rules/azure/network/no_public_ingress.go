package network

import (
	"github.com/aquasecurity/defsec/cidr"
	"github.com/aquasecurity/defsec/providers"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/severity"
	"github.com/aquasecurity/defsec/state"
)

var CheckNoPublicIngress = rules.Register(
	rules.Rule{
		AVDID:      "AVD-AZU-0047",
		Provider:   providers.AzureProvider,
		Service:    "network",
		ShortCode:  "no-public-ingress",
		Summary:    "An inbound network security rule allows traffic from /0.",
		Impact:     "The port is exposed for ingress from the internet",
		Resolution: "Set a more restrictive cidr range",
		Explanation: `Network security rules should not use very broad subnets.

Where possible, segments should be broken into smaller subnets.`,
		Links: []string{
			"https://docs.microsoft.com/en-us/azure/security/fundamentals/network-best-practices",
		},
		Terraform: &rules.EngineMetadata{
			GoodExamples:        terraformNoPublicIngressGoodExamples,
			BadExamples:         terraformNoPublicIngressBadExamples,
			Links:               terraformNoPublicIngressLinks,
			RemediationMarkdown: terraformNoPublicIngressRemediationMarkdown,
		},
		Severity: severity.Critical,
	},
	func(s *state.State) (results rules.Results) {
		for _, group := range s.Azure.Network.SecurityGroups {
			var failed bool
			for _, rule := range group.Rules {
				if rule.Outbound.IsTrue() || rule.Allow.IsFalse() {
					continue
				}
				for _, ip := range rule.SourceAddresses {
					if cidr.IsPublic(ip.Value()) {
						failed = true
						results.Add(
							"Security group rule allows ingress from public internet.",
							ip,
						)
					}
				}
			}
			if !failed {
				results.AddPassed(&group)
			}
		}
		return
	},
)
