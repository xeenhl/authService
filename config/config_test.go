package config

import (
	"testing"
)

func TestConfigurationLoading(t *testing.T) {

	config = nil

	c, err := LoadConfiguration("./config_test.json")

	if err != nil {
		t.Errorf("File should be loaded without an error, but was: [%v]", err)
	}

	if c == nil {
		t.Errorf("Configuration should be loaded but was: nil")
	} else {
		if c.Port != 9999 {
			t.Errorf("Expected port [%v], but was: [%v]", 9999, c.Port)
		}
		if c.Auth.PublicKey != "PublicKeyPath" {
			t.Errorf("Expected Auth.PublicKey [%v], but was: [%v]", "PublicKeyPath", c.Auth.PublicKey)
		}
		if c.Auth.PrivateKey != "PrivateKeyPath" {
			t.Errorf("Expected Auth.PrivateKey [%v], but was: [%v]", "PrivateKeyPath", c.Auth.PrivateKey)
		}
	}

}

func TestGetConfigReturnsNilIfNoConfigLoaded(t *testing.T) {

	config = nil

	c, err := GetConfig()

	if err == nil {
		t.Errorf("[No Configuration loaded] error should be return if no configuration was loaded, but was: nil")
	}

	if c != nil {
		t.Errorf("No config mast be loaded before first LoadConfiguration call")
	}

}

func TestGetConfigReturnsLoadedConfigurationt(t *testing.T) {

	config = nil

	LoadConfiguration("./config_test.json")

	c, err := GetConfig()

	if err != nil {
		t.Errorf("Configuration mas returns without error, but was: [%v]", err)
	}

	if err != nil {
		t.Errorf("File should be loaded without an error, but was: [%v]", err)
	}

	if c == nil {
		t.Errorf("Configuration should be loaded but was: nil")
	} else {
		if c.Port != 9999 {
			t.Errorf("Expected port [%v], but was: [%v]", 9999, c.Port)
		}
		if c.Auth.PublicKey != "PublicKeyPath" {
			t.Errorf("Expected Auth.PublicKey [%v], but was: [%v]", "PublicKeyPath", c.Auth.PublicKey)
		}
		if c.Auth.PrivateKey != "PrivateKeyPath" {
			t.Errorf("Expected Auth.PrivateKey [%v], but was: [%v]", "PrivateKeyPath", c.Auth.PrivateKey)
		}
	}

}
