package sqldump

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/si3nloong/sqlike/util"
)

var null = []byte(`null`)

// Parser :
type Parser func([]byte) string

func byteToString(data []byte) string {
	if len(data) == 0 {
		return `""`
	}
	return strconv.Quote(string(data))
}

func numToString(data []byte) string {
	return string(data)
}

func tsToString(data []byte) string {
	t, _ := time.Parse(time.RFC3339, string(data))
	return t.UTC().Format(`"2006-01-02 15:04:05.999999999"`)
}

func dateToString(data []byte) string {
	t, _ := time.Parse(time.RFC3339, string(data))
	return t.UTC().Format(`"2006-01-02"`)
}

func jsonToString(data []byte) string {
	if bytes.Equal(data, null) {
		return `"null"`
	}
	return strconv.Quote(string(data))
}

func setToString(data []byte) string {
	str := string(data)
	vals := strings.Split(str, ",")
	blr := util.AcquireString()
	defer util.ReleaseString(blr)
	blr.WriteString(`'`)
	for i, v := range vals {
		if i > 0 {
			blr.WriteByte(',')
		}
		blr.WriteString(v)
	}
	blr.WriteString(`'`)
	return blr.String()
}

// func StringPoint(data []byte) string {
// 	// if len(data) == 42 {
// 	// 	dst := make([]byte, 21)
// 	// 	_, err := hex.Decode(dst, data)
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// 	data = dst
// 	// }

// 	scanner := wkb.Scanner(nil)
// 	scanner.Scan(data)
// 	// if len(data) == 21 {
// 	// 	// the length of a point type in WKB
// 	// 	return scan.Scan(data[:])
// 	// }

// 	// if len(data) == 25 {
// 	// 	// Most likely MySQL's SRID+WKB format.
// 	// 	// However, could be a line string or multipoint with only one point.
// 	// 	// But those would be invalid for parsing a point.
// 	// 	// return p.unmarshalWKB(data[4:])
// 	// 	if err := scanner.Scan(data[4:]); err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// }

// 	// log.Println("debug =>", scanner.Geometry.GeoJSONType())
// 	// log.Println(scanner.Geometry.Bound())
// 	return ""
// 	// return fmt.Sprintf(`ST_PointFromText("POINT(%v %v)")`, p.Lon(), p.Lat())
// }
