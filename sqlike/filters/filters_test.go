package filters_test

import (
	"log"
	"regexp"
	"testing"
)

func TestFilter(t *testing.T) {
	query := "user.id==133,category==;asdas,asdas"
	reg := regexp.MustCompile(`(\,|\;)`)
	log.Println(reg.FindAllString(query, -1))
	idxs := reg.FindAllStringIndex(query, -1)
	for i, pos := range idxs {
		a := 0
		if i > 0 {
			a = idxs[i-1][1]
		}
		log.Println(i, pos, query[pos[0]:pos[1]], query[a:pos[0]])
	}
	log.Println(reg.FindAllStringSubmatch(query, -1))
	paths := reg.Split(query, -1)
	log.Println(paths, len(paths))
	for i, p := range paths {
		log.Println(i+1, ":", p)
	}
	// query = `user.id%3D%3D133,;`
	// // query = `admin.id==133`
	// url.Parse(``)
	// v, err := url.ParseQuery(query)
	// log.Println(v, err)
	// url.ParseRequestURI(``)
}
