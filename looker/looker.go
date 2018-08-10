package looker

import (
	"path"
	"reflect"
)

type Package struct {
	ImportPath ImportElement
	Interfaces Interfaces
}

// ImportElement represents the imported package.
// Attribute `Alias` represents the optional alias for the package.
//		import foo "github.com/fooooo/baaaar-pppkkkkggg"
type ImportElement struct {
	Path  string
	Alias string
}

func (ie ImportElement) Name() string {
	if ie.Alias != "" {
		return ie.Alias
	}

	return path.Base(ie.Path)
}

func LookAtInterfaces(pkgPath string, is []reflect.Type) *Package {
	pkg := Package{
		ImportPath: ImportElement{Path: pkgPath},
		Interfaces: make(Interfaces, 0, len(is)),
	}
	for _, it := range is {
		intf := LookAtInterface(it)
		pkg.Interfaces = append(pkg.Interfaces, intf)
	}

	return &pkg
}

type Interface struct {
	Name    string
	Methods Methods
}

type Interfaces []*Interface

func LookAtInterface(typ reflect.Type) *Interface {
	intf := &Interface{
		Name:    typ.Name(),
		Methods: make(Methods, 0, typ.NumMethod()),
	}

	for i := 0; i < typ.NumMethod(); i++ {
		mt := typ.Method(i)
		m := Method{
			Name: mt.Name,
		}
		in, out := LookAtFuncParameters(typ.Method(i).Type)
		m.In = in
		m.Out = out

		intf.Methods = append(intf.Methods, &m)
	}
	return intf
}

type Method struct {
	Name string
	In   Parameters
	Out  Parameters
}

type Methods []*Method

type Parameter interface {
	Kind() reflect.Kind
	String() string
}

type Parameters []Parameter

func LookAtFuncParameters(mt reflect.Type) (Parameters, Parameters) {
	var in = make(Parameters, 0)
	for i := 0; i < mt.NumIn(); i++ {
		in = append(in, LookAtParameter(mt.In(i)))
	}

	var out = make(Parameters, 0)
	for i := 0; i < mt.NumOut(); i++ {
		out = append(out, LookAtParameter(mt.Out(i)))
	}

	return in, out
}

type StructElement struct {
	ImportPath ImportElement
	BaseType   reflect.Kind
	UserType   string
	Pointer    bool
	Fields     Fields
}

func (prm *StructElement) Kind() reflect.Kind {
	return reflect.Struct
}

func (prm *StructElement) PtrPrefix() string {
	if prm.Pointer {
		return "*"
	}
	return ""
}

func (prm *StructElement) String() string {
	return prm.PtrPrefix() + prm.ImportPath.Name() + "." + prm.UserType
}

func LookAtParameter(at reflect.Type) *StructElement {
	var pointer bool
	if at.Kind() == reflect.Ptr {
		at = at.Elem()
		pointer = true
	}
	prm := StructElement{
		ImportPath: ImportElement{Path: at.PkgPath()},
		BaseType:   at.Kind(),
		UserType:   at.Name(),
		Pointer:    pointer,
	}

	if prm.BaseType == reflect.Struct {
		prm.Fields = LookAtFields(at)
	}

	return &prm
}

type Field struct {
	Name       string
	ImportPath ImportElement
	BaseType   reflect.Kind
	UserType   string
	Anonymous  bool
}

type Fields []Field

func LookAtFields(st reflect.Type) Fields {
	fields := make(Fields, 0, st.NumField())
	for i := 0; i < st.NumField(); i++ {
		ft := st.Field(i)
		fields = append(fields, LookAtField(ft))
	}
	return fields
}

func LookAtField(ft reflect.StructField) Field {
	return Field{
		Name:       ft.Name,
		ImportPath: ImportElement{Path: ft.Type.PkgPath()},
		BaseType:   ft.Type.Kind(),
		UserType:   ft.Type.Name(),
		Anonymous:  ft.Anonymous,
	}
}
