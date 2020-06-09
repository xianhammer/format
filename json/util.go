package json

import (
	"errors"
	"reflect"
	"unsafe"
)

var (
	ErrMustBePointer = errors.New("A pointer was expected")
	ErrInvalidType   = errors.New("Cannot convert to given type")
)

/* Unmarshal, example call
func main() {
	src := []byte(`{"a":"hello","b":{"c":-2,"d":[9,3.141]}}`)
	var out = struct {
		a string
		b struct {
			c int
			d []float32
		}
	}{}

	obj, _, err := json.Parse(src, nil)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(obj, &out)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("out = %v\n", out)
}
*/

// Unmarshal try to map values from content of an interface to content of a known struct
func Unmarshal(a, b interface{}) (err error) {
	vb := reflect.ValueOf(b)
	if vb.Kind() != reflect.Ptr {
		return ErrMustBePointer
	}

	switch va := a.(type) {
	case float64:
		switch vb.Elem().Kind() {
		default:
			err = ErrInvalidType
		case reflect.Int:
			(*b.(*int)) = int(va)
		case reflect.Int8:
			(*b.(*int8)) = int8(va)
		case reflect.Int16:
			(*b.(*int16)) = int16(va)
		case reflect.Int32:
			(*b.(*int32)) = int32(va)
		case reflect.Int64:
			(*b.(*int64)) = int64(va)
		case reflect.Uint:
			(*b.(*uint)) = uint(va)
		case reflect.Uint8:
			(*b.(*uint8)) = uint8(va)
		case reflect.Uint16:
			(*b.(*uint16)) = uint16(va)
		case reflect.Uint32:
			(*b.(*uint32)) = uint32(va)
		case reflect.Uint64:
			(*b.(*uint64)) = uint64(va)
		// case reflect.Uintptr:
		// 	(*b.(*int)) = int(va)
		case reflect.Float32:
			(*b.(*float32)) = float32(va)
		case reflect.Float64:
			(*b.(*float64)) = va
		}

	case string:
		if reflect.TypeOf(b) != reflect.PtrTo(reflect.TypeOf(a)) {
			return ErrInvalidType
		}
		(*b.(*string)) = va

	case bool:
		if reflect.TypeOf(b) != reflect.PtrTo(reflect.TypeOf(a)) {
			return ErrInvalidType
		}
		(*b.(*bool)) = va

	case []interface{}:
		slice := reflect.TypeOf(b).Elem()
		elemType := slice.Elem()

		switch elemType.Kind() {
		default:
			err = ErrInvalidType

		case reflect.Map:
		case reflect.Slice:

		case reflect.Int:
			a := make([]int, len(va))
			(*b.(*[]int)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		case reflect.Int8:
			a := make([]int8, len(va))
			(*b.(*[]int8)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		case reflect.Int16:
			a := make([]int16, len(va))
			(*b.(*[]int16)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		case reflect.Int32:
			a := make([]int32, len(va))
			(*b.(*[]int32)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		case reflect.Int64:
			a := make([]int64, len(va))
			(*b.(*[]int64)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		case reflect.Uint:
			a := make([]uint, len(va))
			(*b.(*[]uint)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		case reflect.Uint8:
			a := make([]uint8, len(va))
			(*b.(*[]uint8)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		case reflect.Uint16:
			a := make([]uint16, len(va))
			(*b.(*[]uint16)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		case reflect.Uint32:
			a := make([]uint32, len(va))
			(*b.(*[]uint32)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		case reflect.Uint64:
			a := make([]uint64, len(va))
			(*b.(*[]uint64)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		case reflect.Float32:
			a := make([]float32, len(va))
			(*b.(*[]float32)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		case reflect.Float64:
			a := make([]float64, len(va))
			(*b.(*[]float64)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		case reflect.Bool:
			a := make([]bool, len(va))
			(*b.(*[]bool)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		case reflect.String:
			a := make([]string, len(va))
			(*b.(*[]string)) = a
			for i, e := range va {
				if err = Unmarshal(e, &a[i]); err != nil {
					break
				}
			}
		}

	case map[string]interface{}:
		val := vb.Elem()
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			fv := val.Field(i)
			key := typ.Field(i).Name
			if value, found := va[key]; found {
				// See https://stackoverflow.com/questions/42664837/how-to-access-unexported-struct-fields-in-golang/43918797#43918797
				pv := reflect.NewAt(fv.Type(), unsafe.Pointer(fv.UnsafeAddr())) //.Elem()
				if err = Unmarshal(value, pv.Interface()); err != nil {
					break
				}
			}
		}

	default:
		// TODO Error?
	}

	return
}

func Equal(a, b interface{}) bool {
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)

	if ta != tb {
		return false
	}

	switch va := a.(type) {
	case float64:
		return va == b.(float64)

	case string:
		return va == b.(string)

	case bool:
		return va == b.(bool)

	case []interface{}:
		vb := b.([]interface{})
		if len(va) != len(vb) {
			return false
		}
		for i := 0; i < len(va); i++ {
			if !Equal(va[i], vb[i]) {
				return false
			}
		}

	case map[string]interface{}:
		vb := b.(map[string]interface{})
		if len(va) != len(vb) {
			return false
		}

		for key, value := range va {
			if !Equal(value, vb[key]) {
				return false
			}
		}

	default:
		return a == nil && b == nil
	}
	return true
}
