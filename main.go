package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/imdario/mergo"
	"k8s.io/apimachinery/pkg/util/sets"
)

type stringSetMerger struct {
}

func (t stringSetMerger) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	ts := typ.String()
	fmt.Println(ts)
	if typ == reflect.TypeOf([]string{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				if dst.Len() <= 1 && src.Len() == 0 {
					return nil
				}
				if dst.Len() == 0 && src.Len() == 1 {
					dst.Set(src)
					return nil
				}

				out := sets.NewString()
				for i := 0; i < dst.Len(); i++ {
					out.Insert(dst.Index(i).String())
				}
				for i := 0; i < src.Len(); i++ {
					out.Insert(src.Index(i).String())
				}
				dst.Set(reflect.ValueOf(out.List()))
			}
			return nil
		}
	}
	return nil
}

type timeTransformer struct {
}

func (t timeTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(time.Time{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				isZero := dst.MethodByName("IsZero")
				result := isZero.Call([]reflect.Value{})
				if result[0].Bool() {
					dst.Set(src)
				}
			}
			return nil
		}
	}
	return nil
}

type S2 string

const sw S2 = "a"

type Team struct {
	Names []string
}

func main() {
	src := Team{Names: []string{"a", "b", "a"}}
	dest := Team{Names: []string{"a", "b", "a"}}
	mergo.Merge(&dest, src, mergo.WithTransformers(timeTransformer{}), mergo.WithTransformers(stringSetMerger{}))
	fmt.Println(dest)
	// Will print
	// { 2018-01-12 01:15:00 +0000 UTC m=+0.000000001 }
}
