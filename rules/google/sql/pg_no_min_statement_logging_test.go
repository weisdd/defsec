package sql

import (
	"testing"

	"github.com/aquasecurity/defsec/parsers/types"
	"github.com/aquasecurity/defsec/providers/google/sql"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/state"
	"github.com/stretchr/testify/assert"
)

func TestCheckPgNoMinStatementLogging(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.SQL
		expected bool
	}{
		{
			name: "Instance logging enabled for all statements",
			input: sql.SQL{
				Metadata: types.NewTestMetadata(),
				Instances: []sql.DatabaseInstance{
					{
						Metadata:        types.NewTestMetadata(),
						DatabaseVersion: types.String("POSTGRES_12", types.NewTestMetadata()),
						Settings: sql.Settings{
							Metadata: types.NewTestMetadata(),
							Flags: sql.Flags{
								Metadata:                types.NewTestMetadata(),
								LogMinDurationStatement: types.Int(0, types.NewTestMetadata()),
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "Instance logging disabled for all statements",
			input: sql.SQL{
				Metadata: types.NewTestMetadata(),
				Instances: []sql.DatabaseInstance{
					{
						Metadata:        types.NewTestMetadata(),
						DatabaseVersion: types.String("POSTGRES_12", types.NewTestMetadata()),
						Settings: sql.Settings{
							Metadata: types.NewTestMetadata(),
							Flags: sql.Flags{
								Metadata:                types.NewTestMetadata(),
								LogMinDurationStatement: types.Int(-1, types.NewTestMetadata()),
							},
						},
					},
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.Google.SQL = test.input
			results := CheckPgNoMinStatementLogging.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == rules.StatusFailed && result.Rule().LongID() == CheckPgNoMinStatementLogging.Rule().LongID() {
					found = true
				}
			}
			if test.expected {
				assert.True(t, found, "Rule should have been found")
			} else {
				assert.False(t, found, "Rule should not have been found")
			}
		})
	}
}
