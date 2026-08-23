package manifest

import (
	"encoding/json"
	"testing"
)

func TestNewWithDefaults(t *testing.T) {
	m := New(Config{
		ID:            "test-id",
		ActionKeyword: "t",
		Name:          "Test",
		Description:   "A test plugin",
		Author:        "me",
		ExeFileName:   "test.exe",
	})
	if m.Language != "executable" {
		t.Errorf("got Language=%s, want executable", m.Language)
	}
	if m.Version != "1.0.0" {
		t.Errorf("got Version=%s, want 1.0.0", m.Version)
	}
	if m.Schema == "" {
		t.Error("schema should not be empty")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		m       *Manifest
		wantErr bool
	}{
		{"valid", &Manifest{ID: "a", ActionKeyword: "a", Name: "a", Description: "a", Author: "a", Version: "a", Language: "a", ExecuteFileName: "a"}, false},
		{"missing ID", &Manifest{ActionKeyword: "a", Name: "a", Description: "a", Author: "a", Version: "a", Language: "a", ExecuteFileName: "a"}, true},
		{"missing Name", &Manifest{ID: "a", ActionKeyword: "a", Description: "a", Author: "a", Version: "a", Language: "a", ExecuteFileName: "a"}, true},
		{"missing ExecuteFileName", &Manifest{ID: "a", ActionKeyword: "a", Name: "a", Description: "a", Author: "a", Version: "a", Language: "a"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	m := New(Config{
		ID:            "test-id",
		ActionKeyword: "t",
		Name:          "Test",
		Description:   "A test plugin",
		Author:        "me",
		ExeFileName:   "test.exe",
	})
	data, err := m.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	if result["ID"] != "test-id" {
		t.Errorf("got ID=%v", result["ID"])
	}
}

func TestUnmarshal(t *testing.T) {
	data := `{
		"ID": "test-id",
		"ActionKeyword": "t",
		"Name": "Test",
		"Description": "A test plugin",
		"Author": "me",
		"Version": "1.0.0",
		"Language": "executable",
		"ExecuteFileName": "test.exe"
	}`
	m, err := Unmarshal([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	if m.ID != "test-id" {
		t.Errorf("got ID=%s, want test-id", m.ID)
	}
	if m.Language != "executable" {
		t.Errorf("got Language=%s", m.Language)
	}
}

func TestUnmarshal_InvalidJSON(t *testing.T) {
	_, err := Unmarshal([]byte(`{invalid`))
	if err == nil {
		t.Fatal("expected error")
	}
}