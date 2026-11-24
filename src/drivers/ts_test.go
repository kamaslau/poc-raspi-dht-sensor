package drivers

import "testing"

// A client that was never connected must surface an error, not a nil
// dereference: the loop calls TSWrite every minute and would otherwise panic
// long after start-up, with nothing in the log pointing at the configuration.
func TestTSWriteWithoutClient(t *testing.T) {
	previous := TSClient
	TSClient = nil
	t.Cleanup(func() { TSClient = previous })

	err := TSWrite("dht11", map[string]string{"type": "dht11"}, map[string]any{"temperature": 21.5})
	if err == nil {
		t.Fatal("want an error when the client is nil, got nil")
	}
}
