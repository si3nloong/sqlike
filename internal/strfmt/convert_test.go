package strfmt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnakeCase(t *testing.T) {
	require.Equal(t, `abcef_ghij_kk`, ToSnakeCase("abcef ghij kk"))
	require.Equal(t, `abcd_j`, ToSnakeCase("ABCD j"))
	require.Equal(t, `id`, ToSnakeCase(`id`))
	require.Equal(t, `const_k`, ToSnakeCase(`CONST_K`))
	require.Equal(t, `marshal_json`, ToSnakeCase(`MarshalJSON`))
}

func TestPascalCase(t *testing.T) {
	require.Equal(t, `AbcefGhijKk`, ToPascalCase("abcef ghij kk"))
	require.Equal(t, `AbcdJ`, ToPascalCase("ABCD j"))
	require.Equal(t, `Id`, ToPascalCase(`id`))
	require.Equal(t, `ConstK`, ToPascalCase(`CONST_K`))
	require.Equal(t, `MarshalJson`, ToPascalCase(`MarshalJSON`))
}

func TestCamelCase(t *testing.T) {
	require.Equal(t, `abcefGhijKk`, ToCamelCase("abcef ghij kk"))
	require.Equal(t, `abcdJ`, ToCamelCase("ABCD j"))
	require.Equal(t, `id`, ToCamelCase(`id`))
	require.Equal(t, `constK`, ToCamelCase(`CONST_K`))
	require.Equal(t, `marshalJson`, ToCamelCase(`MarshalJSON`))
}
