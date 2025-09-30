package setup_test 

import(
	"testing"

	"github.com/hihikaAAa/meeting-events/internal/lib/logger/setup"
)

func TestSetupLogger_NotNil(t *testing.T){
	for _, env := range []string{"local","dev","prod","unknown"} {
		if l := setup.SetupLogger(env); l == nil {
			t.Fatalf("logger is nil for env=%s", env)
		}
	}
}