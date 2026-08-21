package flowapi

// ── Flow control ──────────────────────────────────────────────────────────────

// ChangeQuery changes Flow Launcher's query text.
type ChangeQuery struct {
	Query   string
	Requery bool
}

func (c ChangeQuery) Method() string { return "Flow.Launcher.ChangeQuery" }
func (c ChangeQuery) Params() []any  { return []any{c.Query, c.Requery} }

// RestartApp restarts Flow Launcher.
type RestartApp struct{}

func (RestartApp) Method() string { return "Flow.Launcher.RestartApp" }
func (RestartApp) Params() []any  { return nil }

// ReQuery reloads the current query.
type ReQuery struct {
	Reselect bool
}

func (r ReQuery) Method() string { return "Flow.Launcher.ReQuery" }
func (r ReQuery) Params() []any  { return []any{r.Reselect} }

// BackToQueryResults returns to query results from context menu or history.
type BackToQueryResults struct{}

func (BackToQueryResults) Method() string { return "Flow.Launcher.BackToQueryResults" }
func (BackToQueryResults) Params() []any  { return nil }

// FocusQueryTextBox focuses the query text box.
type FocusQueryTextBox struct{}

func (FocusQueryTextBox) Method() string { return "Flow.Launcher.FocusQueryTextBox" }
func (FocusQueryTextBox) Params() []any  { return nil }

// ── Window visibility ─────────────────────────────────────────────────────────

// ShowMainWindow shows the main window.
type ShowMainWindow struct{}

func (ShowMainWindow) Method() string { return "Flow.Launcher.ShowMainWindow" }
func (ShowMainWindow) Params() []any  { return nil }

// HideMainWindow hides the main window.
type HideMainWindow struct{}

func (HideMainWindow) Method() string { return "Flow.Launcher.HideMainWindow" }
func (HideMainWindow) Params() []any  { return nil }

// IsMainWindowVisible checks if the main window is visible.
type IsMainWindowVisible struct{}

func (IsMainWindowVisible) Method() string { return "Flow.Launcher.IsMainWindowVisible" }
func (IsMainWindowVisible) Params() []any  { return nil }

// ── Shell & clipboard ─────────────────────────────────────────────────────────

// ShellRun runs a shell command.
type ShellRun struct {
	Cmd      string
	Filename string // default "cmd.exe"
}

func (s ShellRun) Method() string { return "Flow.Launcher.ShellRun" }
func (s ShellRun) Params() []any  { return []any{s.Cmd, s.Filename} }

// CopyToClipboard copies text to the clipboard.
type CopyToClipboard struct {
	Text                     string
	DirectCopy               bool
	ShowDefaultNotification  bool
}

func (c CopyToClipboard) Method() string { return "Flow.Launcher.CopyToClipboard" }
func (c CopyToClipboard) Params() []any  { return []any{c.Text, c.DirectCopy, c.ShowDefaultNotification} }

// ── Settings & save ───────────────────────────────────────────────────────────

// SaveAppAllSettings saves all Flow Launcher and plugin settings.
type SaveAppAllSettings struct{}

func (SaveAppAllSettings) Method() string { return "Flow.Launcher.SaveAppAllSettings" }
func (SaveAppAllSettings) Params() []any  { return nil }

// SavePluginSettings saves all plugin settings.
type SavePluginSettings struct{}

func (SavePluginSettings) Method() string { return "Flow.Launcher.SavePluginSettings" }
func (SavePluginSettings) Params() []any  { return nil }

// SavePluginCaches saves all plugin caches.
type SavePluginCaches struct{}

func (SavePluginCaches) Method() string { return "Flow.Launcher.SavePluginCaches" }
func (SavePluginCaches) Params() []any  { return nil }

// OpenSettingDialog opens the setting dialog.
type OpenSettingDialog struct{}

func (OpenSettingDialog) Method() string { return "Flow.Launcher.OpenSettingDialog" }
func (OpenSettingDialog) Params() []any  { return nil }

// OpenPluginSettingsWindow opens the settings window for a specific plugin.
type OpenPluginSettingsWindow struct {
	PluginID string
}

func (o OpenPluginSettingsWindow) Method() string { return "Flow.Launcher.OpenPluginSettingsWindow" }
func (o OpenPluginSettingsWindow) Params() []any  { return []any{o.PluginID} }

// ── Plugin management ─────────────────────────────────────────────────────────

// GetAllPlugins gets all loaded plugins.
type GetAllPlugins struct{}

func (GetAllPlugins) Method() string { return "Flow.Launcher.GetAllPlugins" }
func (GetAllPlugins) Params() []any  { return nil }

// GetAllInitializedPlugins gets all initialized plugins.
type GetAllInitializedPlugins struct {
	IncludeFailed bool
}

func (g GetAllInitializedPlugins) Method() string { return "Flow.Launcher.GetAllInitializedPlugins" }
func (g GetAllInitializedPlugins) Params() []any  { return []any{g.IncludeFailed} }

// AddActionKeyword adds an action keyword to a plugin.
type AddActionKeyword struct {
	PluginID         string
	NewActionKeyword string
}

func (a AddActionKeyword) Method() string { return "Flow.Launcher.AddActionKeyword" }
func (a AddActionKeyword) Params() []any  { return []any{a.PluginID, a.NewActionKeyword} }

// RemoveActionKeyword removes an action keyword from a plugin.
type RemoveActionKeyword struct {
	PluginID         string
	OldActionKeyword string
}

func (r RemoveActionKeyword) Method() string { return "Flow.Launcher.RemoveActionKeyword" }
func (r RemoveActionKeyword) Params() []any  { return []any{r.PluginID, r.OldActionKeyword} }

// ActionKeywordAssigned checks if an action keyword is assigned (obsolete, always returns false).
type ActionKeywordAssigned struct {
	ActionKeyword string
}

func (a ActionKeywordAssigned) Method() string { return "Flow.Launcher.ActionKeywordAssigned" }
func (a ActionKeywordAssigned) Params() []any  { return []any{a.ActionKeyword} }

// ReloadAllPluginData reloads all plugin data.
type ReloadAllPluginData struct{}

func (ReloadAllPluginData) Method() string { return "Flow.Launcher.ReloadAllPluginData" }
func (ReloadAllPluginData) Params() []any  { return nil }

// ── Notifications ─────────────────────────────────────────────────────────────

// ShowMsg shows a message box (3-param overload).
type ShowMsg struct {
	Title    string
	SubTitle string
	IconPath string
}

func (s ShowMsg) Method() string { return "Flow.Launcher.ShowMsg" }
func (s ShowMsg) Params() []any  { return []any{s.Title, s.SubTitle, s.IconPath} }

// ShowMsgWithOwner shows a message box with owner window option (4-param overload).
type ShowMsgWithOwner struct {
	Title                string
	SubTitle             string
	IconPath             string
	UseMainWindowAsOwner bool
}

func (s ShowMsgWithOwner) Method() string { return "Flow.Launcher.ShowMsg" }
func (s ShowMsgWithOwner) Params() []any  { return []any{s.Title, s.SubTitle, s.IconPath, s.UseMainWindowAsOwner} }

// ShowMsgError shows an error message using Flow's standard error icon.
type ShowMsgError struct {
	Title    string
	SubTitle string
}

func (s ShowMsgError) Method() string { return "Flow.Launcher.ShowMsgError" }
func (s ShowMsgError) Params() []any  { return []any{s.Title, s.SubTitle} }

// ── Updates ───────────────────────────────────────────────────────────────────

// CheckForNewUpdate checks for a new Flow Launcher update.
type CheckForNewUpdate struct{}

func (CheckForNewUpdate) Method() string { return "Flow.Launcher.CheckForNewUpdate" }
func (CheckForNewUpdate) Params() []any  { return nil }

// UpdatePluginManifestAsync updates the plugin manifest.
type UpdatePluginManifestAsync struct {
	UsePrimaryURLOnly bool
}

func (u UpdatePluginManifestAsync) Method() string { return "Flow.Launcher.UpdatePluginManifestAsync" }
func (u UpdatePluginManifestAsync) Params() []any  { return []any{u.UsePrimaryURLOnly} }

// PluginModified checks if a plugin has been modified.
type PluginModified struct {
	ID string
}

func (p PluginModified) Method() string { return "Flow.Launcher.PluginModified" }
func (p PluginModified) Params() []any  { return []any{p.ID} }

// GetPluginManifest gets the current plugin manifest.
type GetPluginManifest struct{}

func (GetPluginManifest) Method() string { return "Flow.Launcher.GetPluginManifest" }
func (GetPluginManifest) Params() []any  { return nil }

// ── HTTP ──────────────────────────────────────────────────────────────────────

// HttpGetStringAsync downloads a URL as string.
type HttpGetStringAsync struct {
	URL string
}

func (h HttpGetStringAsync) Method() string { return "Flow.Launcher.HttpGetStringAsync" }
func (h HttpGetStringAsync) Params() []any  { return []any{h.URL} }

// HttpGetStreamAsync downloads a URL as stream.
type HttpGetStreamAsync struct {
	URL string
}

func (h HttpGetStreamAsync) Method() string { return "Flow.Launcher.HttpGetStreamAsync" }
func (h HttpGetStreamAsync) Params() []any  { return []any{h.URL} }

// ── Search ────────────────────────────────────────────────────────────────────

// FuzzySearch performs a fuzzy search.
type FuzzySearch struct {
	Query          string
	StringToCompare string
}

func (f FuzzySearch) Method() string { return "Flow.Launcher.FuzzySearch" }
func (f FuzzySearch) Params() []any  { return []any{f.Query, f.StringToCompare} }

// GetTranslation gets a translation for the current language.
type GetTranslation struct {
	Key string
}

func (g GetTranslation) Method() string { return "Flow.Launcher.GetTranslation" }
func (g GetTranslation) Params() []any  { return []any{g.Key} }

// ── Logging ───────────────────────────────────────────────────────────────────

// LogDebug logs a debug message.
type LogDebug struct {
	ClassName  string
	Message    string
	MethodName string
}

func (l LogDebug) Method() string { return "Flow.Launcher.LogDebug" }
func (l LogDebug) Params() []any  { return []any{l.ClassName, l.Message, l.MethodName} }

// LogInfo logs an info message.
type LogInfo struct {
	ClassName  string
	Message    string
	MethodName string
}

func (l LogInfo) Method() string { return "Flow.Launcher.LogInfo" }
func (l LogInfo) Params() []any  { return []any{l.ClassName, l.Message, l.MethodName} }

// LogWarn logs a warning message.
type LogWarn struct {
	ClassName  string
	Message    string
	MethodName string
}

func (l LogWarn) Method() string { return "Flow.Launcher.LogWarn" }
func (l LogWarn) Params() []any  { return []any{l.ClassName, l.Message, l.MethodName} }

// LogError logs an error message.
type LogError struct {
	ClassName  string
	Message    string
	MethodName string
}

func (l LogError) Method() string { return "Flow.Launcher.LogError" }
func (l LogError) Params() []any  { return []any{l.ClassName, l.Message, l.MethodName} }

// ── UI ────────────────────────────────────────────────────────────────────────

// OpenDirectory opens a directory in the file explorer.
type OpenDirectory struct {
	DirectoryPath     string
	FileNameOrFilePath string
}

func (o OpenDirectory) Method() string { return "Flow.Launcher.OpenDirectory" }
func (o OpenDirectory) Params() []any  { return []any{o.DirectoryPath, o.FileNameOrFilePath} }

// OpenUrl opens a URL in the browser.
type OpenUrl struct {
	URL      string
	InPrivate bool // nil means use default
}

func (o OpenUrl) Method() string { return "Flow.Launcher.OpenUrl" }
func (o OpenUrl) Params() []any  { return []any{o.URL, o.InPrivate} }

// OpenWebUrl opens a URL in the browser even if it's a local file.
type OpenWebUrl struct {
	URL      string
	InPrivate bool
}

func (o OpenWebUrl) Method() string { return "Flow.Launcher.OpenWebUrl" }
func (o OpenWebUrl) Params() []any  { return []any{o.URL, o.InPrivate} }

// OpenAppUri opens an application URI.
type OpenAppUri struct {
	AppURI string
}

func (o OpenAppUri) Method() string { return "Flow.Launcher.OpenAppUri" }
func (o OpenAppUri) Params() []any  { return []any{o.AppURI} }

// ── Loading bar ───────────────────────────────────────────────────────────────

// StartLoadingBar starts the loading animation.
type StartLoadingBar struct{}

func (StartLoadingBar) Method() string { return "Flow.Launcher.StartLoadingBar" }
func (StartLoadingBar) Params() []any  { return nil }

// StopLoadingBar stops the loading animation.
type StopLoadingBar struct{}

func (StopLoadingBar) Method() string { return "Flow.Launcher.StopLoadingBar" }
func (StopLoadingBar) Params() []any  { return nil }

// ── Game mode ─────────────────────────────────────────────────────────────────

// ToggleGameMode toggles game mode on/off.
type ToggleGameMode struct{}

func (ToggleGameMode) Method() string { return "Flow.Launcher.ToggleGameMode" }
func (ToggleGameMode) Params() []any  { return nil }

// SetGameMode sets game mode to a specific value.
type SetGameMode struct {
	Value bool
}

func (s SetGameMode) Method() string { return "Flow.Launcher.SetGameMode" }
func (s SetGameMode) Params() []any  { return []any{s.Value} }

// IsGameModeOn checks if game mode is on.
type IsGameModeOn struct{}

func (IsGameModeOn) Method() string { return "Flow.Launcher.IsGameModeOn" }
func (IsGameModeOn) Params() []any  { return nil }

// ── Theme ─────────────────────────────────────────────────────────────────────

// GetAvailableThemes gets all available themes.
type GetAvailableThemes struct{}

func (GetAvailableThemes) Method() string { return "Flow.Launcher.GetAvailableThemes" }
func (GetAvailableThemes) Params() []any  { return nil }

// GetCurrentTheme gets the current theme.
type GetCurrentTheme struct{}

func (GetCurrentTheme) Method() string { return "Flow.Launcher.GetCurrentTheme" }
func (GetCurrentTheme) Params() []any  { return nil }

// ── System ────────────────────────────────────────────────────────────────────

// GetDataDirectory gets the Flow Launcher data directory.
type GetDataDirectory struct{}

func (GetDataDirectory) Method() string { return "Flow.Launcher.GetDataDirectory" }
func (GetDataDirectory) Params() []any  { return nil }

// GetLogDirectory gets the Flow Launcher log directory.
type GetLogDirectory struct{}

func (GetLogDirectory) Method() string { return "Flow.Launcher.GetLogDirectory" }
func (GetLogDirectory) Params() []any  { return nil }

// IsApplicationDarkTheme checks if the application is using a dark theme.
type IsApplicationDarkTheme struct{}

func (IsApplicationDarkTheme) Method() string { return "Flow.Launcher.IsApplicationDarkTheme" }
func (IsApplicationDarkTheme) Params() []any  { return nil }