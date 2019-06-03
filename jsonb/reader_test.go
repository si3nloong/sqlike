package jsonb

import (
	"log"
	"testing"
)

func TestReader(t *testing.T) {
	data := []byte(`
	{         
		"Test" :"hello world!!" ,
		"Test2"   : "x1#$%^&*xx",
		"Test4": {
			"Test" :"hello world!!" ,
			"Test2" :"hello world!!" ,
			"Testxx" :"hello world!!" , 
			"nested"  : {
				"deep" : {
					"value" : 199303.00
				}
			}
		},
		"Test0": 100.111,
		"Test99": 6000,
		"Bool": true
	}   		
	
	`)

	var o struct {
		Test  string
		Test0 float64
		Test2 string
		Test4 struct {
			Test   string
			Test2  string
			Nested struct {
				Deep struct {
					Value float64 `sqlike:"value"`
				} `sqlike:"deep"`
			} `sqlike:"nested"`
		}
		Test99 int
		Bool   bool
	}

	if err := Unmarshal(data, &o); err != nil {
		panic(err)
	}
	log.Println("Result :", o)
	log.Println("Result :", o.Test4.Nested.Deep.Value)

}
