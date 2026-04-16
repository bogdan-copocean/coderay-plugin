package parsing

// ParseContext holds file path, UTF-8 source, and resolved language config.
type ParseContext struct {
	FilePath string
	Content  string
	LangCfg  *LanguageConfig
}

// NewParseContext returns a context for supported extensions, or nil.
func NewParseContext(filePath, content string) *ParseContext {
	cfg := GetLanguageForFile(filePath)
	if cfg == nil {
		return nil
	}
	return &ParseContext{
		FilePath: filePath,
		Content:  content,
		LangCfg:  cfg,
	}
}
