package msgpack_test

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"

	"code.byted.org/ad/msgpack_extstr"
	"code.byted.org/ad/msgpack_extstr/msgpcode"
	"github.com/stretchr/testify/require"
)

func init() {
	msgpack.RegisterExtString("ext_stringX", (*ExtStringTest)(nil))
}

type ExtStringTest struct {
	S string
}

var (
	_ msgpack.Marshaler   = (*ExtStringTest)(nil)
	_ msgpack.Unmarshaler = (*ExtStringTest)(nil)
)

func (ext ExtStringTest) MarshalMsgpack() ([]byte, error) {
	return msgpack.Marshal("hello " + ext.S)
}

func (ext *ExtStringTest) UnmarshalMsgpack(b []byte) error {
	return msgpack.Unmarshal(b, &ext.S)
}

func TestEncodeDecodeExtStringHeader(t *testing.T) {
	v := &ExtStringTest{"world"}

	payload, err := v.MarshalMsgpack()
	require.Nil(t, err)

	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	err = enc.EncodeExtStringHeader("ext_stringX", len(payload))
	require.Nil(t, err)

	_, err = buf.Write(payload)
	require.Nil(t, err)

	var dst interface{}
	err = msgpack.Unmarshal(buf.Bytes(), &dst)
	require.Nil(t, err)

	v = dst.(*ExtStringTest)
	wanted := "hello world"
	require.Equal(t, v.S, wanted)

	dec := msgpack.NewDecoder(&buf)
	extID, extLen, err := dec.DecodeExtStringHeader()
	require.Nil(t, err)
	require.Equal(t, "ext_stringX", extID)
	require.Equal(t, len(payload), extLen)

	data := make([]byte, extLen)
	err = dec.ReadFull(data)
	require.Nil(t, err)

	v = &ExtStringTest{}
	err = v.UnmarshalMsgpack(data)
	require.Nil(t, err)
	require.Equal(t, wanted, v.S)
}

func TestExtString(t *testing.T) {
	v := &ExtStringTest{"world"}
	b, err := msgpack.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	var dst interface{}
	err = msgpack.Unmarshal(b, &dst)
	if err != nil {
		t.Fatal(err)
	}

	v, ok := dst.(*ExtStringTest)
	if !ok {
		t.Fatalf("got %#v, wanted ExtStringTest", dst)
	}

	wanted := "hello world"
	if v.S != wanted {
		t.Fatalf("got %q, wanted %q", v.S, wanted)
	}

	ext := new(ExtStringTest)
	err = msgpack.Unmarshal(b, &ext)
	if err != nil {
		t.Fatal(err)
	}
	if ext.S != wanted {
		t.Fatalf("got %q, wanted %q", ext.S, wanted)
	}
}

func TestUnknownExtString(t *testing.T) {
	b := []byte{byte(msgpcode.FixExt1), 2, 0}

	var dst interface{}
	err := msgpack.Unmarshal(b, &dst)
	if err == nil {
		t.Fatalf("got nil, wanted error")
	}
	got := err.Error()
	wanted := "msgpack: unknown ext id=2"
	if got != wanted {
		t.Fatalf("got %q, wanted %q", got, wanted)
	}
}

func TestSliceOfTimeString(t *testing.T) {
	in := []interface{}{time.Now()}
	b, err := msgpack.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}

	var out []interface{}
	err = msgpack.Unmarshal(b, &out)
	if err != nil {
		t.Fatal(err)
	}

	outTime := out[0].(time.Time)
	inTime := in[0].(time.Time)
	if outTime.Unix() != inTime.Unix() {
		t.Fatalf("got %v, wanted %v", outTime, inTime)
	}
}

type customPayloadString struct {
	payload []byte
}

func (cp *customPayloadString) MarshalMsgpack() ([]byte, error) {
	return cp.payload, nil
}

func (cp *customPayloadString) UnmarshalMsgpack(b []byte) error {
	cp.payload = b
	return nil
}

func TestDecodeCustomPayloadString(t *testing.T) {
	b, err := hex.DecodeString("c70500c09eec3100")
	if err != nil {
		t.Fatal(err)
	}

	msgpack.RegisterExt(0, (*customPayload)(nil))

	var cp *customPayload
	err = msgpack.Unmarshal(b, &cp)
	if err != nil {
		t.Fatal(err)
	}

	payload := hex.EncodeToString(cp.payload)
	wanted := "c09eec3100"
	if payload != wanted {
		t.Fatalf("got %q, wanted %q", payload, wanted)
	}
}
