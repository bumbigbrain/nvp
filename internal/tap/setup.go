package tap

import (
	"fmt"

	"github.com/songgao/water"
)

// Setup creates and configures a new TAP interface
func Setup(name string) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TAP,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: name,
		},
	}

	// Create a new TAP interface
	ifce, err := water.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create TAP interface: %w", err)
	}

	return ifce, nil
}
