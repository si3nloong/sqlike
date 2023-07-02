package strfmt

func ToSnakeCase(s string) string {
	return toScreamingDelimited(s, '_', "", false)
}

func ToPascalCase(s string) string {
	return toCamelInitCase(s, true)
}

func ToCamelCase(s string) string {
	return toCamelInitCase(s, false)
}
