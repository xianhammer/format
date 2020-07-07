package json

import "errors"

var ErrBadType = errors.New("Bad type")

/*type Object interface {
	Bool(key interface{}) (bool, error)
	Number(key interface{}) (float64, error)
	String(key interface{}) (string, error)
	Array(key interface{}) ([]Object, error)
	Object(key interface{}) (Object, error)
}
*/
func New(root interface{}) (o *Object) {
	var ok bool
	o, ok = root.(*Object)
	if ok || o != nil {
		return
	}

	a_, _ := root.([]interface{})
	o_, _ := root.(map[string]interface{})
	o = &Object{a_, o_}
	return
}

type Object struct {
	array  []interface{}
	object map[string]interface{}
}

func (o *Object) Bool(key interface{}) (v bool, err error) {
	var ok bool
	if k, ok_ := key.(int); ok_ {
		v, ok = o.array[k].(bool)
	} else if k, ok_ := key.(string); ok_ {
		v, ok = o.object[k].(bool)
	}

	if !ok {
		err = ErrBadType
	}
	return
}

func (o *Object) Number(key interface{}) (v float64, err error) {
	var ok bool
	if k, ok_ := key.(int); ok_ {
		v, ok = o.array[k].(float64)
	} else if k, ok_ := key.(string); ok_ {
		v, ok = o.object[k].(float64)
	}

	if !ok {
		err = ErrBadType
	}
	return
}

func (o *Object) String(key interface{}) (v string, err error) {
	var ok bool
	if k, ok_ := key.(int); ok_ {
		v, ok = o.array[k].(string)
	} else if k, ok_ := key.(string); ok_ {
		v, ok = o.object[k].(string)
	}

	if !ok {
		err = ErrBadType
	}

	return
}

func (o *Object) Array(key interface{}) (v []*Object, err error) {
	var ok bool
	var v_ []interface{}
	if k, ok_ := key.(int); ok_ {
		v_, ok = o.array[k].([]interface{})
	} else if k, ok_ := key.(string); ok_ {
		v_, ok = o.object[k].([]interface{})
	}

	if !ok {
		return nil, ErrBadType
	}

	v = make([]*Object, len(v_))
	for i, e := range v_ {
		v[i] = New(e)
	}

	return
}

func (o *Object) Object(key interface{}) (v *Object, err error) {
	var ok bool
	var v_ map[string]interface{}
	if k, ok_ := key.(int); ok_ {
		v_, ok = o.array[k].(map[string]interface{})
	} else if k, ok_ := key.(string); ok_ {
		v_, ok = o.object[k].(map[string]interface{})
	}

	if !ok {
		return nil, ErrBadType
	}

	return New(v_), nil
}

func (o *Object) BoolPanic(key interface{}) (v bool) {
	if k, ok_ := key.(int); ok_ {
		return o.array[k].(bool)
	}
	if k, ok_ := key.(string); ok_ {
		return o.object[k].(bool)
	}
	return
}

func (o *Object) NumberPanic(key interface{}) (v float64) {
	if k, ok_ := key.(int); ok_ {
		return o.array[k].(float64)
	}
	if k, ok_ := key.(string); ok_ {
		return o.object[k].(float64)
	}
	return
}

func (o *Object) StringPanic(key interface{}) (v string) {
	if k, ok_ := key.(int); ok_ {
		return o.array[k].(string)
	}
	if k, ok_ := key.(string); ok_ {
		return o.object[k].(string)
	}
	return
}

func (o *Object) ArrayPanic(key interface{}) (v []*Object) {
	var v_ []interface{}
	if k, ok_ := key.(int); ok_ {
		v_ = o.array[k].([]interface{})
	} else if k, ok_ := key.(string); ok_ {
		v_ = o.object[k].([]interface{})
	}

	v = make([]*Object, len(v_))
	for i, e := range v_ {
		v[i] = New(e)
	}

	return
}

func (o *Object) ObjectPanic(key interface{}) (v *Object) {
	var v_ map[string]interface{}
	if k, ok_ := key.(int); ok_ {
		v_ = o.array[k].(map[string]interface{})
	} else if k, ok_ := key.(string); ok_ {
		v_ = o.object[k].(map[string]interface{})
	}

	return New(v_)
}
