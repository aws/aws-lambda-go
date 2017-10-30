package test

func GetMalformedJson() []byte {
	return []byte(`{ "Records`)
}
