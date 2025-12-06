package test

// nolint: staticcheck
func GetMalformedJson() []byte {
	return []byte(`{ "Records`)
}
