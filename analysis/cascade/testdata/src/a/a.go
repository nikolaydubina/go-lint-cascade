package a

type Inner struct {
	Value int
}

func (s Inner) WithDefaults() Inner {
	if s.Value == 0 {
		s.Value = 42
	}
	return s
}

type Outer struct {
	Inner Inner
	Name  string
}

func (s Outer) WithDefaults() Outer { // want `Outer.WithDefaults\(\) does not call Inner.WithDefaults\(\)`
	if s.Name == "" {
		s.Name = "default"
	}
	// Missing: s.Inner = s.Inner.WithDefaults()
	return s
}

type OuterCorrect struct {
	Inner Inner
	Name  string
}

func (s OuterCorrect) WithDefaults() OuterCorrect {
	s.Inner = s.Inner.WithDefaults()
	if s.Name == "" {
		s.Name = "default"
	}
	return s
}

type NoDefaults struct {
	Value int
}

type OuterWithNoDefaultsField struct {
	NoDefaults NoDefaults
	Name       string
}

func (s OuterWithNoDefaultsField) WithDefaults() OuterWithNoDefaultsField {
	// NoDefaults doesn't have WithDefaults(), so this is fine
	if s.Name == "" {
		s.Name = "default"
	}
	return s
}

type DeepNested struct {
	A NestedA
}

type NestedA struct {
	B NestedB
}

func (s NestedA) WithDefaults() NestedA {
	s.B = s.B.WithDefaults()
	return s
}

type NestedB struct {
	Value int
}

func (s NestedB) WithDefaults() NestedB {
	if s.Value == 0 {
		s.Value = 100
	}
	return s
}

func (s DeepNested) WithDefaults() DeepNested { // want `DeepNested.WithDefaults\(\) does not call A.WithDefaults\(\)`
	return s
}

type PointerField struct {
	Inner *Inner
}

func (s PointerField) WithDefaults() PointerField { // want `PointerField.WithDefaults\(\) does not call Inner.WithDefaults\(\)`
	return s
}
