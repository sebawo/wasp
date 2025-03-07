package wasp

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type ResponseSample struct {
	Data string
}

type StatsSample struct {
	CallTimeout      float64
	CurrentInstances float64
	CurrentRPS       float64
	RunFailed        bool
	RunStopped       bool
	Success          float64
	Failed           float64
}

type LokiSamplesAssertions struct {
	ResponsesSamples []ResponseSample
	StatsSamples     []StatsSample
}

func assertSamples(t *testing.T, samples []PromtailSendResult, a LokiSamplesAssertions) {
	var cd CallResult
	for i, s := range samples[0:2] {
		err := json.Unmarshal([]byte(s.Entry), &cd)
		require.NoError(t, err)
		require.NotEmpty(t, cd.Duration)
		require.Equal(t, cd.Data, a.ResponsesSamples[i].Data)
	}
	// marshal to map because atomic can't be marshalled
	var ls map[string]interface{}
	for i, s := range samples[2:4] {
		err := json.Unmarshal([]byte(s.Entry), &ls)
		require.NoError(t, err)
		require.Equal(t, ls["callTimeout"], a.StatsSamples[i].CallTimeout)
		require.Equal(t, ls["current_instances"], a.StatsSamples[i].CurrentInstances)
		require.Equal(t, ls["current_rps"], a.StatsSamples[i].CurrentRPS)
		require.Equal(t, ls["run_failed"], a.StatsSamples[i].RunFailed)
		require.Equal(t, ls["run_stopped"], a.StatsSamples[i].RunStopped)
		require.Equal(t, ls["success"], a.StatsSamples[i].Success)
		require.Equal(t, ls["failed"], a.StatsSamples[i].Failed)
	}
}

func TestLokiSamples(t *testing.T) {
	defaultLabels := map[string]string{
		"cluster":    "test_cluster",
		"namespace":  "test_namespace",
		"app":        "test_app",
		"test_group": "test_group",
		"test_id":    "test_id",
	}

	type test struct {
		name       string
		genCfg     *Config
		assertions LokiSamplesAssertions
	}

	tests := []test{
		{
			name: "successful RPS run should contain at least 2 response samples without errors and 2 stats samples",
			genCfg: &Config{
				T: t,
				// empty URL is a special case for mocked client
				LokiConfig: NewDefaultLokiConfig("", ""),
				Labels:     defaultLabels,
				LoadType:   RPSScheduleType,
				Schedule:   Plain(1, 55*time.Millisecond),
				Gun: NewMockGun(&MockGunConfig{
					CallSleep: 50 * time.Millisecond,
				}),
			},
			assertions: LokiSamplesAssertions{
				ResponsesSamples: []ResponseSample{
					{
						Data: "successCallData",
					},
					{
						Data: "successCallData",
					},
				},
				StatsSamples: []StatsSample{
					{
						CallTimeout:      0,
						CurrentInstances: 0,
						CurrentRPS:       1,
						RunFailed:        false,
						RunStopped:       false,
						Success:          2,
						Failed:           0,
					},
					{
						CallTimeout:      0,
						CurrentInstances: 0,
						CurrentRPS:       1,
						RunFailed:        false,
						RunStopped:       false,
						Success:          2,
						Failed:           0,
					},
				},
			}},
	}

	for _, tc := range tests {
		gen, err := NewGenerator(tc.genCfg)
		require.NoError(t, err)
		gen.Run(true)
		assertSamples(t, gen.loki.AllHandleResults(), tc.assertions)
	}

	t.Parallel()
}
