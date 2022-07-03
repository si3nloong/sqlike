package jsonb

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/si3nloong/sqlike/v2/x/reflext"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

type CustomString string

func TestDecodeByte(t *testing.T) {
	var (
		dec = DefaultDecoder{}
		r   *Reader
		x   []byte
		b   []byte
		err error
	)

	v := reflext.ValueOf(&x).Elem()

	r = NewReader([]byte(`""`))
	err = dec.DecodeByte(r, v)
	require.NoError(t, err)
	require.Equal(t, make([]byte, 0), x)

	r = NewReader([]byte(`null`))
	err = dec.DecodeByte(r, v)
	require.NoError(t, err)
	require.Equal(t, []byte(nil), x)

	b = []byte(`"VGhlIGlubGluZSB0YWJsZXMgYWJvdmUgYXJlIGlkZW50aWNhbCB0byB0aGUgZm9sbG93aW5nIHN0YW5kYXJkIHRhYmxlIGRlZmluaXRpb25zOg=="`)
	r = NewReader(b)
	err = dec.DecodeByte(r, v)
	require.NoError(t, err)
	require.Equal(t, []byte(`The inline tables above are identical to the following standard table definitions:`), x)
}

func TestDecodeLanguage(t *testing.T) {
	var (
		dec = DefaultDecoder{}
		r   *Reader
		// tag language.Tag
		x   language.Tag
		err error
	)

	v := reflext.ValueOf(&x).Elem()

	// r = NewReader([]byte(`""`))
	// err = dec.DecodeTime(r, v)
	// require.NoError(t, err)
	// require.Equal(t, time.Time{}, x)

	{
		r = NewReader([]byte(`"en"`))
		err = dec.DecodeLanguage(r, v)
		require.NoError(t, err)
		require.Equal(t, language.English, x)
	}
}

func TestCurrency(t *testing.T) {
	var (
		dec = DefaultDecoder{registry: buildDefaultRegistry()}
		r   *Reader
		x   currency.Unit
		err error
	)

	v := reflext.ValueOf(&x).Elem()

	t.Run("Decode with null", func(ti *testing.T) {
		r = NewReader([]byte(`null`))
		err = dec.DecodeCurrency(r, v)
		require.NoError(t, err)
		require.Equal(t, currency.Unit{}, x)
	})

	t.Run("Decode with value", func(ti *testing.T) {
		r = NewReader([]byte(`"USD"`))
		err = dec.DecodeCurrency(r, v)
		require.NoError(t, err)
		require.Equal(t, currency.USD, x)
	})

	t.Run("Decode with invalid value", func(ti *testing.T) {
		r = NewReader([]byte(`"USDT"`))
		err = dec.DecodeCurrency(r, v)
		require.Error(t, err)
	})
}

func TestDecodeTime(t *testing.T) {
	var (
		dec = DefaultDecoder{}
		r   *Reader
		dt  time.Time
		x   time.Time
		err error
	)

	v := reflext.ValueOf(&x).Elem()

	r = NewReader([]byte(`""`))
	err = dec.DecodeTime(r, v)
	require.NoError(t, err)
	require.Equal(t, time.Time{}, x)

	dt, _ = time.Parse("2006-01-02", "2018-01-02")
	r = NewReader([]byte(`"2018-01-02"`))
	err = dec.DecodeTime(r, v)
	require.NoError(t, err)
	require.Equal(t, dt, x)

	dt, _ = time.Parse("2006-01-02 15:04:05", "2018-01-02 13:58:26")
	r = NewReader([]byte(`"2018-01-02 13:58:26"`))
	err = dec.DecodeTime(r, v)
	require.NoError(t, err)
	require.Equal(t, dt, x)

	t.Run("Decode Time w invalid format", func(it *testing.T) {
		r = NewReader([]byte(`"2018-01-02 13:65:66"`))
		err = dec.DecodeTime(r, v)
		require.Error(it, err)

		r = NewReader([]byte(`2018-01-02 13:65:66"`))
		err = dec.DecodeTime(r, v)
		require.Error(it, err)

		r = NewReader([]byte(`"2018-01-02 13:65:66`))
		err = dec.DecodeTime(r, v)
		require.Error(it, err)

		r = NewReader([]byte(``))
		err = dec.DecodeTime(r, v)
		require.Error(it, err)
	})

}

func TestDecodeJSONRaw(t *testing.T) {
	var (
		dec = DefaultDecoder{registry: buildDefaultRegistry()}
		r   *Reader
		x   json.RawMessage
		err error
	)

	v := reflext.ValueOf(&x).Elem()

	t.Run("Decode with null", func(ti *testing.T) {
		r = NewReader([]byte(`null`))
		err = dec.DecodeJSONRaw(r, v)
		require.NoError(t, err)
		require.Equal(t, json.RawMessage(`null`), x)
	})

	t.Run("Decode with number", func(ti *testing.T) {
		r = NewReader([]byte(`98`))
		err = dec.DecodeJSONRaw(r, v)
		require.NoError(t, err)
		require.Equal(t, json.RawMessage(`98`), x)
	})
}

func TestDecodeJSONNumber(t *testing.T) {
	var (
		dec = DefaultDecoder{registry: buildDefaultRegistry()}
		r   *Reader
		x   json.Number
		err error
	)

	v := reflext.ValueOf(&x).Elem()

	t.Run("Decode with null", func(ti *testing.T) {
		r = NewReader([]byte(`null`))
		err = dec.DecodeJSONNumber(r, v)
		require.NoError(t, err)
		require.Equal(t, json.Number("0"), x)
	})

	t.Run("Decode with position integer", func(ti *testing.T) {
		r = NewReader([]byte(`88`))
		err = dec.DecodeJSONNumber(r, v)
		require.NoError(t, err)
		require.Equal(t, json.Number("88"), x)
		require.Equal(t, "88", x.String())
		i64, _ := x.Int64()
		require.Equal(t, int64(88), i64)
		f64, _ := x.Float64()
		require.Equal(t, float64(88), f64)
	})

	t.Run("Decode json.Number w invalid value", func(ti *testing.T) {
		r = NewReader([]byte(`hsdjdkd`))
		err = dec.DecodeJSONNumber(r, v)
		require.Error(t, err)

		// r = NewReader([]byte(`10.3a`))
		// err = dec.DecodeJSONNumber(r, v)
		// require.Error(t, err)
	})
}

func TestDecodeMap(t *testing.T) {
	var (
		dec = DefaultDecoder{registry: buildDefaultRegistry()}
		r   *Reader
		x   map[string]any
		err error
	)

	v := reflext.ValueOf(&x).Elem()

	t.Run("Decode with null", func(ti *testing.T) {
		r = NewReader([]byte(`null`))
		err = dec.DecodeMap(r, v)
		require.NoError(t, err)
		require.Equal(t, map[string]any(nil), x)
	})

	t.Run("Decode with empty object", func(ti *testing.T) {
		r = NewReader([]byte(`{}`))
		err = dec.DecodeMap(r, v)
		require.NoError(t, err)
		require.Equal(t, make(map[string]any), x)
	})

	t.Run("Decode to map<string,any>", func(ti *testing.T) {
		r = NewReader([]byte(`
		{
			"a":"123", 
			"b":   108213312, 
			"c": true, 
			"d": "alSLKaj28173-021@#$%^&*\"",
			"e": 0.3127123
		}`))
		err = dec.DecodeMap(r, v)
		require.NoError(t, err)
		require.Equal(t, map[string]any{
			"a": "123",
			"b": float64(108213312),
			"c": true,
			"d": `alSLKaj28173-021@#$%^&*"`,
			"e": float64(0.3127123),
		}, x)

	})

	t.Run("Decode to map<string,string>", func(ti *testing.T) {
		r = NewReader([]byte(`
		{
			"number":      "1234567890", 
			"b":"abcdefghijklmnopqrstuvwxyz",
			"emoji": "😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊",
			"japanese": "福岡市美術館で夜間開館スタート！7月～10月の金曜日と土曜日は20時まで延長開館"
		}`))
		m := make(map[string]string)
		v := reflect.ValueOf(&m)
		err = dec.DecodeMap(r, v.Elem())
		require.NoError(ti, err)
		require.Equal(ti, map[string]string{
			"number":   "1234567890",
			"b":        "abcdefghijklmnopqrstuvwxyz",
			"emoji":    "😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊",
			"japanese": "福岡市美術館で夜間開館スタート！7月～10月の金曜日と土曜日は20時まで延長開館",
		}, m)
	})

	t.Run("Decode to map<string,bool>", func(ti *testing.T) {
		r = NewReader([]byte(`
		{
			"true":     true, 
			"false": false
		}`))
		m := make(map[string]bool)
		v := reflect.ValueOf(&m)
		err = dec.DecodeMap(r, v.Elem())
		require.NoError(ti, err)
		require.Equal(ti, map[string]bool{
			"true":  true,
			"false": false,
		}, m)
	})

	t.Run("Decode to map<string,int>", func(ti *testing.T) {
		r = NewReader([]byte(`
		{
			"minus-one": -1,
			"negative": -31231237,
			"one":      1, 
			"two":2,
			"eleven": 11,
			"hundred": 100
		}`))
		m := make(map[string]int)
		v := reflect.ValueOf(&m)
		err = dec.DecodeMap(r, v.Elem())
		require.NoError(ti, err)
		require.Equal(ti, map[string]int{
			"minus-one": -1,
			"negative":  -31231237,
			"one":       1,
			"two":       2,
			"eleven":    11,
			"hundred":   100,
		}, m)
	})

	t.Run("Decode to map<string,uint8>", func(ti *testing.T) {
		r = NewReader([]byte(`
		{
			"one":      1, 
			"two":2,
			"eleven": 11,
			"hundred": 100
		}`))
		m := make(map[string]uint8)
		v := reflect.ValueOf(&m)
		err = dec.DecodeMap(r, v.Elem())
		require.NoError(ti, err)
		require.Equal(ti, map[string]uint8{
			"one":     1,
			"two":     2,
			"eleven":  11,
			"hundred": 100,
		}, m)
	})

	t.Run("Decode to map<string,float32>", func(ti *testing.T) {
		r = NewReader([]byte(`
		{
			"minus-one": -1,
			"negative":  -31231237,
			"one":      1, 
			"two":2,
			"eleven": 11,
			"hundred": 100,
			"number":    3123123799213,
		}`))
		m := make(map[string]float32)
		v := reflect.ValueOf(&m)
		err = dec.DecodeMap(r, v.Elem())
		require.NoError(ti, err)
		require.Equal(ti, map[string]float32{
			"minus-one": -1,
			"negative":  -31231237,
			"one":       1,
			"two":       2,
			"eleven":    11,
			"hundred":   100,
			"number":    3123123799213,
		}, m)
	})

	t.Run("Decode to map<string,float64>", func(ti *testing.T) {
		r = NewReader([]byte(`
		{
			"minus-one": -1,
			"negative":  -3123123799213,
			"one":      1, 
			"two":2,
			"eleven": 11,
			"hundred": 100,
			"number":    3123123799213,
		}`))
		m := make(map[string]float64)
		v := reflect.ValueOf(&m)
		err = dec.DecodeMap(r, v.Elem())
		require.NoError(ti, err)
		require.Equal(ti, map[string]float64{
			"minus-one": -1,
			"negative":  -3123123799213,
			"one":       1,
			"two":       2,
			"eleven":    11,
			"hundred":   100,
			"number":    3123123799213,
		}, m)
	})

	t.Run("Decode to map<string,any>", func(ti *testing.T) {
		r = NewReader([]byte(`
		{
			"negative": -183,
			"string": "textasjdhasljdlasjkdjlsa:'dasdas",
			"number":    3123123799213,
			"nested": {
				"k": {
					"bool": true,
					"no": 10,
					"string": "😀😁😂"
				}
			}
		}`))
		m := make(map[string]any)
		v := reflect.ValueOf(&m)
		err = dec.DecodeMap(r, v.Elem())
		require.NoError(ti, err)
		require.Equal(ti, map[string]any{
			"negative": float64(-183),
			"string":   "textasjdhasljdlasjkdjlsa:'dasdas",
			"number":   float64(3123123799213),
			"nested": map[string]any{
				"k": map[string]any{
					"bool":   true,
					"no":     float64(10),
					"string": "😀😁😂",
				},
			},
		}, m)
	})

	t.Run("Decode to map[CustomString]*string", func(ti *testing.T) {
		var m map[CustomString]*string

		r = NewReader([]byte(`
		{
			"0": "zero",
			"1": "one",
			"2": "two",
			"3": "three"
		}`))

		v := reflect.ValueOf(&m)
		err = dec.DecodeMap(r, v.Elem())
		require.NoError(ti, err)
		zero, one, two, three := "zero", "one", "two", "three"
		require.Equal(t, map[CustomString]*string{
			CustomString("0"): &zero,
			CustomString("1"): &one,
			CustomString("2"): &two,
			CustomString("3"): &three,
		}, m)
	})

	t.Run("Decode with unsupported data type", func(ti *testing.T) {
		var mx map[*CustomString]string
		r = NewReader([]byte(`
		{
			"0": "zero",
			"1": "one",
			"2": "two",
			"3": "three"
		}`))
		v = reflect.ValueOf(&mx)
		err = dec.DecodeMap(r, v.Elem())
		require.Error(ti, err)
	})

}

func TestDecodeArray(t *testing.T) {
	var (
		err error
	)

	t.Run("Decode to [2]string", func(ti *testing.T) {
		arr := [2]string{}
		err = Unmarshal([]byte(`["test", "abc", "ddd"]`), &arr)
		require.Error(ti, err)

		err = Unmarshal([]byte(`["京都着物レンタル夢館", "aBcdEfgHiJklmnO"]`), &arr)
		require.NoError(ti, err)
		require.ElementsMatch(ti, [...]string{
			"京都着物レンタル夢館",
			"aBcdEfgHiJklmnO",
		}, arr)
	})

	t.Run("Decode to []int", func(ti *testing.T) {
		arr := []int{}
		err = Unmarshal([]byte(`[1, "xdd", 3]`), &arr)
		require.Error(ti, err)

		err = Unmarshal([]byte(`[1, 88, 3, -1992, 9999]`), &arr)
		require.NoError(ti, err)
		require.ElementsMatch(t, []int{1, 88, 3, -1992, 9999}, arr)
	})

	t.Run("Decode to []struct{}", func(it *testing.T) {
		type Region struct {
			DialingCode int
			CountryCode string
		}

		type Contact struct {
			Name        string
			LangCode    language.Tag
			PhoneNumber string
			IsPrimary   bool
			Region      Region
		}

		type Address struct {
			Remark         string
			Address        string
			Contact        Contact
			FloorNumber    string
			BuildingNumber string
			EntranceNumber string
		}

		addrs := []Address{}
		b := []byte(`{hsh:""}`)
		err = Unmarshal(b, &addrs)
		require.Error(t, err)

		b = []byte(`[
			{
				"Remark": "One Utama, 1 UTAMA SHOPPING CENTRE,  LEBUH BANDAR UTAMA,  BANDAR UTAMA",
				"Address": "1 UTAMA SHOPPING CENTRE,  LEBUH BANDAR UTAMA,  BANDAR UTAMA, KOTA DAMANSARA PETALING JAYA, 47800, Petaling Jaya, Selangor, Malaysia",
				"Contact": {
					"Name": "One Utama",
					"PhoneNumber": "60176473298"
				},
				"FloorNumber": "",
				"BuildingNumber": "",
				"EntranceNumber": ""
			},
			{
				"Remark": "",
				"Address": "Gugusan Melur, Kota Damansara, 47810, Petaling Jaya, Selangor, Malaysia",
				"Contact": {
					"Name": "Mohamed Yussuf",
					"PhoneNumber": "60176473298",
					"IsPrimary": true,
					"Region": {
						"DialingCode": 60,
						"CountryCode": "MY"
					}
				},
				"FloorNumber": "",
				"BuildingNumber": "",
				"EntranceNumber": ""
			}
		]`)

		addrs = []Address{} // reset address
		err = Unmarshal(b, &addrs)
		require.NoError(t, err)

		require.ElementsMatch(t, []Address{
			{
				Address: "1 UTAMA SHOPPING CENTRE,  LEBUH BANDAR UTAMA,  BANDAR UTAMA, KOTA DAMANSARA PETALING JAYA, 47800, Petaling Jaya, Selangor, Malaysia",
				Remark:  "One Utama, 1 UTAMA SHOPPING CENTRE,  LEBUH BANDAR UTAMA,  BANDAR UTAMA",
				Contact: Contact{
					Name:        "One Utama",
					PhoneNumber: "60176473298",
				},
			},
			{
				Address: "Gugusan Melur, Kota Damansara, 47810, Petaling Jaya, Selangor, Malaysia",
				Contact: Contact{
					Name:        "Mohamed Yussuf",
					PhoneNumber: "60176473298",
					IsPrimary:   true,
					Region: Region{
						DialingCode: 60,
						CountryCode: "MY",
					},
				},
			},
		}, addrs)
	})
}
