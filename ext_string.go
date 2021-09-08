package msgpack

import (
	"fmt"
	"reflect"

	"github.com/myhyh/msgpack/v5/msgpcode"
)

type extStringInfo struct {
	Type    reflect.Type
	Decoder func(d *Decoder, v reflect.Value, extLen int) error
}

var extStringTypes = make(map[string]*extInfo)


func RegisterExtString(extID string, value MarshalerUnmarshaler) {
	RegisterExtStringEncoder(extID, value, func(e *Encoder, v reflect.Value) ([]byte, error) {
		marshaler := v.Interface().(Marshaler)
		return marshaler.MarshalMsgpack()
	})
	RegisterExtStringDecoder(extID, value, func(d *Decoder, v reflect.Value, extLen int) error {
		b, err := d.readN(extLen)
		if err != nil {
			return err
		}
		return v.Interface().(Unmarshaler).UnmarshalMsgpack(b)
	})
}

func UnregisterExtString(extID string) {
	unregisterExtStringEncoder(extID)
	unregisterExtStringDecoder(extID)
}

func RegisterExtStringEncoder(
	extID string,
	value interface{},
	encoder func(enc *Encoder, v reflect.Value) ([]byte, error),
) {
	unregisterExtStringEncoder(extID)

	typ := reflect.TypeOf(value)
	extEncoder := makeExtStringEncoder(extID, typ, encoder)
	typeEncMap.Store(extID, typ)
	typeEncMap.Store(typ, extEncoder)
	if typ.Kind() == reflect.Ptr {
		typeEncMap.Store(typ.Elem(), makeExtStringEncoderAddr(extEncoder))
	}
}

func unregisterExtStringEncoder(extID string) {
	t, ok := typeEncMap.Load(extID)
	if !ok {
		return
	}
	typeEncMap.Delete(extID)
	typ := t.(reflect.Type)
	typeEncMap.Delete(typ)
	if typ.Kind() == reflect.Ptr {
		typeEncMap.Delete(typ.Elem())
	}
}

func makeExtStringEncoder(
	extID string,
	typ reflect.Type,
	encoder func(enc *Encoder, v reflect.Value) ([]byte, error),
) encoderFunc {
	nilable := typ.Kind() == reflect.Ptr

	return func(e *Encoder, v reflect.Value) error {
		if nilable && v.IsNil() {
			return e.EncodeNil()
		}

		b, err := encoder(e, v)
		if err != nil {
			return err
		}

		if err := e.EncodeExtStringHeader(extID, len(b)); err != nil {
			return err
		}

		return e.write(b)
	}
}

func makeExtStringEncoderAddr(extEncoder encoderFunc) encoderFunc {
	return func(e *Encoder, v reflect.Value) error {
		if !v.CanAddr() {
			return fmt.Errorf("msgpack: Decode(nonaddressable %T)", v.Interface())
		}
		return extEncoder(e, v.Addr())
	}
}

func RegisterExtStringDecoder(
	extID string,
	value interface{},
	decoder func(dec *Decoder, v reflect.Value, extLen int) error,
) {
	unregisterExtStringDecoder(extID)

	typ := reflect.TypeOf(value)
	extDecoder := makeExtStringDecoder(extID, typ, decoder)
	extStringTypes[extID] = &extInfo{
		Type:    typ,
		Decoder: decoder,
	}

	typeDecMap.Store(extID, typ)
	typeDecMap.Store(typ, extDecoder)
	if typ.Kind() == reflect.Ptr {
		typeDecMap.Store(typ.Elem(), makeExtStringDecoderAddr(extDecoder))
	}
}

func unregisterExtStringDecoder(extID string) {
	t, ok := typeDecMap.Load(extID)
	if !ok {
		return
	}
	typeDecMap.Delete(extID)
	delete(extStringTypes, extID)
	typ := t.(reflect.Type)
	typeDecMap.Delete(typ)
	if typ.Kind() == reflect.Ptr {
		typeDecMap.Delete(typ.Elem())
	}
}

func makeExtStringDecoder(
	wantedExtID string,
	typ reflect.Type,
	decoder func(d *Decoder, v reflect.Value, extLen int) error,
) decoderFunc {
	return nilAwareDecoder(typ, func(d *Decoder, v reflect.Value) error {
		extID, extLen, err := d.DecodeExtStringHeader()
		if err != nil {
			return err
		}
		if extID != wantedExtID {
			return fmt.Errorf("msgpack: got ext type=%d, wanted %d", extID, wantedExtID)
		}
		return decoder(d, v, extLen)
	})
}

func makeExtStringDecoderAddr(extDecoder decoderFunc) decoderFunc {
	return func(d *Decoder, v reflect.Value) error {
		if !v.CanAddr() {
			return fmt.Errorf("msgpack: Decode(nonaddressable %T)", v.Interface())
		}
		return extDecoder(d, v.Addr())
	}
}

func (e *Encoder) EncodeExtStringHeader(extID string, extLen int) error {
	if err := e.encodeExtStringLen(extLen); err != nil {
		return err
	}
	if err := e.EncodeString(extID); err != nil {
		return err
	}
	return nil
}

func (e *Encoder) encodeExtStringLen(l int) error {
	return e.write4(msgpcode.ExtStr, uint32(l))
}

func (d *Decoder) DecodeExtStringHeader() (extID string, extLen int, err error) {
	c, err := d.readCode()
	if err != nil {
		return
	}
	return d.extHeaderString(c)
}

func (d *Decoder) extHeaderString(c byte) (string, int, error) {
	extLen, err := d.parseExtStringLen(c)
	if err != nil {
		return "", 0, err
	}

	extID, err := d.DecodeString()
	if err != nil {
		return "", 0, err
	}

	return extID, extLen, nil
}

func (d *Decoder) parseExtStringLen(c byte) (int, error) {
	switch c {
	case msgpcode.ExtStr:
		n,err := d.uint32()
		return int(n),err
	default:
		return 0, fmt.Errorf("msgpack: invalid code=%x decoding ext len", c)
	}
}

func (d *Decoder) decodeInterfaceExtString(c byte) (interface{}, error) {
	extID, extLen, err := d.extHeaderString(c)
	if err != nil {
		return nil, err
	}

	info, ok := extStringTypes[extID]
	if !ok {
		return nil, fmt.Errorf("msgpack: unknown ext id=%d", extID)
	}

	v := reflect.New(info.Type).Elem()
	if nilable(v.Kind()) && v.IsNil() {
		v.Set(reflect.New(info.Type.Elem()))
	}

	if err := info.Decoder(d, v, extLen); err != nil {
		return nil, err
	}

	return v.Interface(), nil
}

func (d *Decoder) skipExtString(c byte) error {
	n, err := d.parseExtStringLen(c)
	if err != nil {
		return err
	}
	err = d.skipN(n)
	if err != nil {
		return err
	}
	_, err = d.DecodeString()
	return err
}

func (d *Decoder) skipextHeaderString(c byte) error {
	// Read ext type.
	_, err := d.readCode()
	if err != nil {
		return err
	}
	// Read ext body len.
	for i := 0; i < extHeaderStringLen(c); i++ {
		_, err := d.readCode()
		if err != nil {
			return err
		}
	}
	return nil
}

func extHeaderStringLen(c byte) int {
	switch c {
	case msgpcode.ExtStr:
		return 4
	}
	return 0
}
