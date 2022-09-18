package sqlx

import (
	"crypto/aes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestEncryptColumn_Basic(t *testing.T) {

	testIntCases := []struct {
		name           string
		intVal         int64
		wantEncryptErr error
		wantDecryptErr error
		key            []byte
	}{
		{
			name:           "小于255",
			intVal:         234,
			wantEncryptErr: nil,
			wantDecryptErr: nil,
			key:            []byte("ABCDABCDABCDABCD"),
		},
		{
			name:           "16位key",
			intVal:         234,
			wantEncryptErr: nil,
			wantDecryptErr: nil,
			key:            []byte("aBCDABCdABCDABCD"),
		},
		{
			name:   "24位key",
			intVal: 234,
			key:    []byte("ABCDABCDABCDABCDaBCDABCd"),
		},
		{
			name:   "32位key",
			intVal: 234,
			key:    []byte("ABCDABCDABCDABCDaBCDABCdABCDABCD"),
		},
		{
			name:           "错误长度key",
			intVal:         234,
			wantEncryptErr: aes.KeySizeError(len("aBCDABCdABCDABCD1234f")),
			key:            []byte("aBCDABCdABCDABCD1234f"),
		},
		{
			name:           "大数",
			intVal:         12324213,
			wantEncryptErr: nil,
			wantDecryptErr: nil,
			key:            []byte("aBCDABCdABCDABCD"),
		},
		{
			name:           "特大数",
			intVal:         12324213222223123,
			wantEncryptErr: nil,
			wantDecryptErr: nil,
			key:            []byte("aBCDABCdABCDABCD"),
		},
		{
			name:           "负数",
			intVal:         -12324213222223123,
			wantEncryptErr: nil,
			wantDecryptErr: nil,
			key:            []byte("aBCDABCdABCDABCD"),
		},
	}

	cryptIntE := EncryptColumn[int]{}
	cryptIntD := EncryptColumn[int]{}
	cryptUint8E := EncryptColumn[uint8]{}
	cryptUint8D := EncryptColumn[uint8]{}
	cryptUint16E := EncryptColumn[uint16]{}
	cryptUint16D := EncryptColumn[uint16]{}
	cryptUint32E := EncryptColumn[uint32]{}
	cryptUint32D := EncryptColumn[uint32]{}
	cryptUint64E := EncryptColumn[uint64]{}
	cryptUint64D := EncryptColumn[uint64]{}
	cryptUIntE := EncryptColumn[uint]{}
	cryptUIntD := EncryptColumn[uint]{}

	for _, ts := range testIntCases {
		t.Run(ts.name, func(t *testing.T) {
			cryptIntE.Val = int(ts.intVal)
			intVal, err := cryptIntE.Value(ts.key)
			assert.Equal(t, ts.wantEncryptErr, err)
			if ts.wantEncryptErr == nil {
				err = cryptIntD.Scan(intVal, ts.key)
				assert.Equal(t, ts.wantDecryptErr, err)
				assert.Equal(t, int(ts.intVal), cryptIntD.Val)
			}

			cryptUint8E.Val = uint8(ts.intVal)
			uint8Val, err := cryptUint8E.Value(ts.key)
			assert.Equal(t, ts.wantEncryptErr, err)
			if ts.wantEncryptErr == nil {
				err = cryptUint8D.Scan(uint8Val, ts.key)
				assert.Equal(t, ts.wantDecryptErr, err)
				assert.Equal(t, uint8(ts.intVal), cryptUint8D.Val)
			}

			cryptUint16E.Val = uint16(ts.intVal)
			uint16Val, err := cryptUint16E.Value(ts.key)
			assert.Equal(t, ts.wantEncryptErr, err)
			if ts.wantEncryptErr == nil {
				err = cryptUint16D.Scan(uint16Val, ts.key)
				assert.Equal(t, ts.wantDecryptErr, err)
				assert.Equal(t, uint16(ts.intVal), cryptUint16D.Val)
			}

			cryptUint32E.Val = uint32(ts.intVal)
			uint32Val, err := cryptUint32E.Value(ts.key)
			assert.Equal(t, ts.wantEncryptErr, err)
			if ts.wantEncryptErr == nil {
				err = cryptUint32D.Scan(uint32Val, ts.key)
				assert.Equal(t, ts.wantDecryptErr, err)
				assert.Equal(t, uint32(ts.intVal), cryptUint32D.Val)
			}

			cryptUint64E.Val = uint64(ts.intVal)
			uint64Val, err := cryptUint64E.Value(ts.key)
			assert.Equal(t, ts.wantEncryptErr, err)
			if ts.wantEncryptErr == nil {
				err = cryptUint64D.Scan(uint64Val, ts.key)
				assert.Equal(t, ts.wantDecryptErr, err)
				assert.Equal(t, uint64(ts.intVal), cryptUint64D.Val)
			}

			cryptUIntE.Val = uint(ts.intVal)
			uintVal, err := cryptUIntE.Value(ts.key)
			assert.Equal(t, ts.wantEncryptErr, err)
			if ts.wantEncryptErr == nil {
				err = cryptUIntD.Scan(uintVal, ts.key)
				assert.Equal(t, ts.wantDecryptErr, err)
				assert.Equal(t, uint(ts.intVal), cryptUIntD.Val)
			}
		})
	}

	testStringCases := []struct {
		name           string
		val            string
		wantEncryptErr error
		wantDecryptErr error
		key            []byte
	}{
		{
			name:           "简单",
			val:            "大明教你学go",
			wantEncryptErr: nil,
			wantDecryptErr: nil,
			key:            []byte("ABCDABCDABCDABCD"),
		},
		{
			name:           "长一点",
			val:            "大明教你学go,的撒放你家看你反饥饿卡三顿饭呢asdfjk3jrfnjkdnsafjknbjhbasdf阿斯顿发",
			wantEncryptErr: nil,
			wantDecryptErr: nil,
			key:            []byte("ABCDABCDABCDABCD"),
		},
		{
			name:           "空",
			val:            "",
			wantEncryptErr: nil,
			wantDecryptErr: nil,
			key:            []byte("ABCDABCDABCDABCD"),
		},
	}
	cryptStringE := EncryptColumn[string]{}
	cryptStringD := EncryptColumn[string]{}

	for _, ts := range testStringCases {
		t.Run(ts.name, func(t *testing.T) {
			cryptStringE.Val = ts.val
			value, err := cryptStringE.Value(ts.key)
			assert.Equal(t, ts.wantEncryptErr, err)
			if err == nil {
				fmt.Println(string(reflect.ValueOf(value).Bytes()))
				err = cryptStringD.Scan(value, ts.key)
				assert.Equal(t, ts.wantDecryptErr, err)
				assert.Equal(t, ts.val, cryptStringD.Val)
			}
		})
	}
}

func TestEncryptColumn_Struct(t *testing.T) {
	key := []byte("ABCDABCDABCDABCD")

	cryptSimple := EncryptColumn[Simple]{
		Val: Simple{
			A: 1,
			B: 1.2,
			D: false,
		},
	}
	val, err := cryptSimple.Value(key)
	fmt.Println(string(reflect.ValueOf(val).Bytes()))
	assert.Equal(t, nil, err)

	cryptSimpleD := EncryptColumn[Simple]{}
	err = cryptSimpleD.Scan(val, key)
	assert.Equal(t, nil, err)
	assert.Equal(t, Simple{
		A: 1,
		B: 1.2,
		D: false}, cryptSimpleD.Val)

	composite := Composite{
		E: 2,
		F: 2.1,
		G: true,
		H: "abc",
		Simple: Simple{
			A: 1,
			B: 1.2,
			D: false,
		},
	}

	cryptComposite := EncryptColumn[Composite]{
		Val: composite,
	}
	val, err = cryptComposite.Value(key)
	fmt.Println(string(reflect.ValueOf(val).Bytes()))
	assert.Equal(t, nil, err)

	cryptCompositeD := EncryptColumn[Composite]{}
	err = cryptCompositeD.Scan(val, key)
	assert.Equal(t, nil, err)
	assert.Equal(t, composite, cryptCompositeD.Val)

}

type Simple struct {
	A int
	B float32
	D bool
}

type Composite struct {
	E int
	F float32
	G bool
	H string
	Simple
}

func TestEncryptColumn_Error(t *testing.T) {

}
