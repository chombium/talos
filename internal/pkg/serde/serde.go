package serde

import (
	"fmt"
)

// Serde describes a serializer/deserializer.
type Serde interface {
	Fields() []*Field
}

// FieldSerializerFunc is the func signature for serialization.
type FieldSerializerFunc = func([]byte, interface{}) error

// FieldDeserializerFunc is the func signature for deserialization.
type FieldDeserializerFunc = func(uint32, uint32, []byte, interface{}) ([]byte, error)

// Field represents a field in a datastructure.
type Field struct {
	Offset           uint32
	Length           uint32
	Contents         *[]byte
	SerializerFunc   FieldSerializerFunc
	DeserializerFunc FieldDeserializerFunc
}

// Ser serializes a field.
func Ser(t Serde, data []byte, offset uint32, opts interface{}) error {
	for _, field := range t.Fields() {
		if field.SerializerFunc == nil {
			return fmt.Errorf("the field is missing the serializer function")
		}

		if err := field.SerializerFunc(data[field.start(offset):field.end(offset)], opts); err != nil {
			return err
		}
	}

	return nil
}

// De deserializes a field.
func De(t Serde, data []byte, offset uint32, opts interface{}) error {
	for _, field := range t.Fields() {
		if field.DeserializerFunc == nil {
			return fmt.Errorf("the field is missing the deserializer function")
		}

		contents, err := field.DeserializerFunc(field.Offset, field.Length, data, opts)
		if err != nil {
			return err
		}

		if n := copy(data[field.start(offset):field.end(offset)], contents); uint32(n) != field.Length {
			return fmt.Errorf("expected to write %d elements, wrote %d", field.Length, n)
		}
	}

	return nil
}

func (fld *Field) start(offset uint32) uint32 {
	return fld.Offset + offset
}

func (fld *Field) end(offset uint32) uint32 {
	return fld.Offset + fld.Length + offset
}