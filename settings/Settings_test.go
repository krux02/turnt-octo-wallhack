package settings

import "testing"

func TestSettingsSerializer(t *testing.T) {
	var settings1, settings2 BoolOptions
	settings1.NoParticlePhysics = true
	settings1.NoPlayerPhysics = true

	settings1.Save()
	settings2.Load()

	if settings1.NoParticleRender != settings2.NoParticleRender ||
		settings1.NoParticlePhysics != settings2.NoParticlePhysics ||
		settings1.NoWorldRender != settings2.NoWorldRender ||
		settings1.NoTreeRender != settings2.NoTreeRender ||
		settings1.NoPlayerPhysics != settings2.NoPlayerPhysics ||
		settings1.Wireframe != settings2.Wireframe {
		t.Errorf("settings1: %v\nsettings2: %v\n", settings1, settings2)
	}
}
