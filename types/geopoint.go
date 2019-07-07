package types

// // GeoPoint :
// type GeoPoint struct {
// 	Latitude  float64
// 	Longitude float64
// }

// // DataType :
// func (gp GeoPoint) DataType(driver string, sf *reflext.StructField) columns.Column {
// 	return columns.Column{
// 		Name:     sf.Path,
// 		DataType: "POINT",
// 		Type:     "POINT",
// 		Nullable: sf.IsNullable,
// 	}
// }

// // Value :
// func (gp GeoPoint) Value() (driver.Value, error) {
// 	return gp.String(), nil
// }

// // Scan :
// func (gp *GeoPoint) Scan(it interface{}) error {
// 	switch vi := it.(type) {
// 	case []byte:
// 		return gp.unmarshal(util.UnsafeString(vi))

// 	case string:
// 		return gp.unmarshal(vi)

// 	case nil:
// 		return nil

// 	default:
// 		return xerrors.New("invalid date format")
// 	}
// }

// // String :
// func (gp GeoPoint) String() string {
// 	blr := util.AcquireString()
// 	defer util.ReleaseString(blr)
// 	gp.marshal(blr)
// 	return blr.String()
// }

// func (gp *GeoPoint) marshal(w writer) {
// 	w.WriteString("POINT")
// 	w.WriteByte('(')
// 	w.WriteString(strconv.FormatFloat(gp.Latitude, 'f', -1, 64))
// 	w.WriteByte(' ')
// 	w.WriteString(strconv.FormatFloat(gp.Longitude, 'f', -1, 64))
// 	w.WriteByte(')')
// }

// func (gp *GeoPoint) unmarshal(str string) (err error) {
// 	paths := strings.SplitN(str, ",", 2)
// 	if len(paths) != 2 {
// 		return xerrors.New("invalid value for GeoPoint")
// 	}

// 	gp.Latitude, err = strconv.ParseFloat(paths[0], 64)
// 	if err != nil {
// 		return
// 	}
// 	gp.Longitude, err = strconv.ParseFloat(paths[1], 64)
// 	if err != nil {
// 		return
// 	}
// 	return
// }
