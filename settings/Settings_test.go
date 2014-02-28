package settings

import "testing"

func TestSettingsSerializer(t *testing.T) {
	var settings1, settings2 BoolOptions
	settings1.NoParticlePhysics = true
	settings1.NoPlayerPhysics = true
	settings1.WaterHeight = 23.125

	settings1.Save()
	settings2.Load()

	if settings1 != settings2 {
		t.Errorf("settings1: %v\nsettings2: %v\n", settings1, settings2)
	}
}
