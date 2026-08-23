package api

import (
	"encoding/json"
	"testing"
)

func TestChangeQuery(t *testing.T) {
	a := ChangeQuery{Query: "test", Requery: true}
	want := `{"method":"Flow.Launcher.ChangeQuery","parameters":["test",true]}`
	got := mustMarshal(t, a)
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestRestartApp(t *testing.T) {
	a := RestartApp{}
	want := `{"method":"Flow.Launcher.RestartApp"}`
	got := mustMarshal(t, a)
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestShellRun(t *testing.T) {
	a := ShellRun{Cmd: "echo hello", Filename: "powershell.exe"}
	want := `{"method":"Flow.Launcher.ShellRun","parameters":["echo hello","powershell.exe"]}`
	got := mustMarshal(t, a)
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestCopyToClipboard(t *testing.T) {
	a := CopyToClipboard{Text: "hello", DirectCopy: false, ShowDefaultNotification: true}
	want := `{"method":"Flow.Launcher.CopyToClipboard","parameters":["hello",false,true]}`
	got := mustMarshal(t, a)
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestShowMsg(t *testing.T) {
	a := ShowMsg{Title: "Hi", SubTitle: "Hello", IconPath: "icon.png"}
	want := `{"method":"Flow.Launcher.ShowMsg","parameters":["Hi","Hello","icon.png"]}`
	got := mustMarshal(t, a)
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestShowMsgWithOwner(t *testing.T) {
	a := ShowMsgWithOwner{Title: "Hi", SubTitle: "Hello", IconPath: "icon.png", UseMainWindowAsOwner: true}
	want := `{"method":"Flow.Launcher.ShowMsg","parameters":["Hi","Hello","icon.png",true]}`
	got := mustMarshal(t, a)
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestShowMsgError(t *testing.T) {
	a := ShowMsgError{Title: "Error", SubTitle: "Something went wrong"}
	want := `{"method":"Flow.Launcher.ShowMsgError","parameters":["Error","Something went wrong"]}`
	got := mustMarshal(t, a)
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestOpenUrl(t *testing.T) {
	a := OpenUrl{URL: "https://example.com", InPrivate: false}
	want := `{"method":"Flow.Launcher.OpenUrl","parameters":["https://example.com",false]}`
	got := mustMarshal(t, a)
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestOpenDirectory(t *testing.T) {
	a := OpenDirectory{DirectoryPath: "C:\\Users", FileNameOrFilePath: ""}
	want := `{"method":"Flow.Launcher.OpenDirectory","parameters":["C:\\Users",""]}`
	got := mustMarshal(t, a)
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestFuzzySearch(t *testing.T) {
	a := FuzzySearch{Query: "test", StringToCompare: "testing"}
	want := `{"method":"Flow.Launcher.FuzzySearch","parameters":["test","testing"]}`
	got := mustMarshal(t, a)
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestGetTranslation(t *testing.T) {
	a := GetTranslation{Key: "hello"}
	want := `{"method":"Flow.Launcher.GetTranslation","parameters":["hello"]}`
	got := mustMarshal(t, a)
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestLogDebug(t *testing.T) {
	a := LogDebug{ClassName: "Test", Message: "debug", MethodName: "TestLogDebug"}
	want := `{"method":"Flow.Launcher.LogDebug","parameters":["Test","debug","TestLogDebug"]}`
	got := mustMarshal(t, a)
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestEmptyStructs(t *testing.T) {
	tests := []struct {
		name string
		m    Method
		want string
	}{
		{"RestartApp", RestartApp{}, `{"method":"Flow.Launcher.RestartApp"}`},
		{"SaveAppAllSettings", SaveAppAllSettings{}, `{"method":"Flow.Launcher.SaveAppAllSettings"}`},
		{"CheckForNewUpdate", CheckForNewUpdate{}, `{"method":"Flow.Launcher.CheckForNewUpdate"}`},
		{"StartLoadingBar", StartLoadingBar{}, `{"method":"Flow.Launcher.StartLoadingBar"}`},
		{"StopLoadingBar", StopLoadingBar{}, `{"method":"Flow.Launcher.StopLoadingBar"}`},
		{"ReloadAllPluginData", ReloadAllPluginData{}, `{"method":"Flow.Launcher.ReloadAllPluginData"}`},
		{"GetAllPlugins", GetAllPlugins{}, `{"method":"Flow.Launcher.GetAllPlugins"}`},
		{"OpenSettingDialog", OpenSettingDialog{}, `{"method":"Flow.Launcher.OpenSettingDialog"}`},
		{"ToggleGameMode", ToggleGameMode{}, `{"method":"Flow.Launcher.ToggleGameMode"}`},
		{"IsGameModeOn", IsGameModeOn{}, `{"method":"Flow.Launcher.IsGameModeOn"}`},
		{"GetAvailableThemes", GetAvailableThemes{}, `{"method":"Flow.Launcher.GetAvailableThemes"}`},
		{"GetCurrentTheme", GetCurrentTheme{}, `{"method":"Flow.Launcher.GetCurrentTheme"}`},
		{"GetDataDirectory", GetDataDirectory{}, `{"method":"Flow.Launcher.GetDataDirectory"}`},
		{"GetLogDirectory", GetLogDirectory{}, `{"method":"Flow.Launcher.GetLogDirectory"}`},
		{"IsApplicationDarkTheme", IsApplicationDarkTheme{}, `{"method":"Flow.Launcher.IsApplicationDarkTheme"}`},
		{"HideMainWindow", HideMainWindow{}, `{"method":"Flow.Launcher.HideMainWindow"}`},
		{"ShowMainWindow", ShowMainWindow{}, `{"method":"Flow.Launcher.ShowMainWindow"}`},
		{"IsMainWindowVisible", IsMainWindowVisible{}, `{"method":"Flow.Launcher.IsMainWindowVisible"}`},
		{"BackToQueryResults", BackToQueryResults{}, `{"method":"Flow.Launcher.BackToQueryResults"}`},
		{"FocusQueryTextBox", FocusQueryTextBox{}, `{"method":"Flow.Launcher.FocusQueryTextBox"}`},
		{"SavePluginSettings", SavePluginSettings{}, `{"method":"Flow.Launcher.SavePluginSettings"}`},
		{"SavePluginCaches", SavePluginCaches{}, `{"method":"Flow.Launcher.SavePluginCaches"}`},
		{"GetPluginManifest", GetPluginManifest{}, `{"method":"Flow.Launcher.GetPluginManifest"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustMarshal(t, tt.m)
			if string(got) != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func mustMarshal(t *testing.T, m Method) []byte {
	t.Helper()
	mm := map[string]any{
		"method":     m.Method(),
		"parameters": m.Params(),
	}
	// Remove parameters key if nil
	if m.Params() == nil {
		delete(mm, "parameters")
	}
	b, err := json.Marshal(mm)
	if err != nil {
		t.Fatal(err)
	}
	return b
}