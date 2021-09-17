package test

//nolint: stylecheck
func GetMalformedJson() []byte {
	return []byte(`{ "Records`)
}
