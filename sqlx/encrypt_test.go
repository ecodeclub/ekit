package sqlx

import (
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
		key            string
	}{
		{
			name:           "小于255",
			intVal:         234,
			wantEncryptErr: nil,
			wantDecryptErr: nil,
			key:            "ABCDABCDABCDABCD",
		},
		{
			name:           "16位key",
			intVal:         234,
			wantEncryptErr: nil,
			wantDecryptErr: nil,
			key:            "aBCDABCdABCDABCD",
		},
		{
			name:   "24位key",
			intVal: 234,
			key:    "ABCDABCDABCDABCDaBCDABCd",
		},
		{
			name:   "32位key",
			intVal: 234,
			key:    "ABCDABCDABCDABCDaBCDABCdABCDABCD",
		},
		{
			name:           "大数",
			intVal:         12324213,
			wantEncryptErr: nil,
			wantDecryptErr: nil,
			key:            "aBCDABCdABCDABCD",
		},
		{
			name:           "特大数",
			intVal:         12324213222223123,
			wantEncryptErr: nil,
			wantDecryptErr: nil,
			key:            "aBCDABCdABCDABCD",
		},
		{
			name:           "负数",
			intVal:         -12324213222223123,
			wantEncryptErr: nil,
			wantDecryptErr: nil,
			key:            "aBCDABCdABCDABCD",
		},
	}

	for _, ts := range testIntCases {
		t.Run(ts.name, func(t *testing.T) {
			cryptIntE := NewEncryptColumn[int](ts.key)
			cryptIntD := NewEncryptColumn[int](ts.key)
			cryptUint8E := NewEncryptColumn[uint8](ts.key)
			cryptUint8D := NewEncryptColumn[uint8](ts.key)
			cryptUint16E := NewEncryptColumn[uint16](ts.key)
			cryptUint16D := NewEncryptColumn[uint16](ts.key)
			cryptUint32E := NewEncryptColumn[uint32](ts.key)
			cryptUint32D := NewEncryptColumn[uint32](ts.key)
			cryptUint64E := NewEncryptColumn[uint64](ts.key)
			cryptUint64D := NewEncryptColumn[uint64](ts.key)
			cryptUIntE := NewEncryptColumn[uint](ts.key)
			cryptUIntD := NewEncryptColumn[uint](ts.key)

			cryptIntE.Val = int(ts.intVal)
			intVal, err := cryptIntE.Value()
			assert.Equal(t, notSpecifyInt, err)
			if err == nil {
				err = cryptIntD.Scan(intVal)
				assert.Equal(t, ts.wantDecryptErr, err)
				assert.Equal(t, int(ts.intVal), cryptIntD.Val)
			}

			cryptUint8E.Val = uint8(ts.intVal)
			uint8Val, err := cryptUint8E.Value()
			assert.Equal(t, ts.wantEncryptErr, err)
			if ts.wantEncryptErr == nil {
				err = cryptUint8D.Scan(uint8Val)
				assert.Equal(t, ts.wantDecryptErr, err)
				assert.Equal(t, uint8(ts.intVal), cryptUint8D.Val)
			}

			cryptUint16E.Val = uint16(ts.intVal)
			uint16Val, err := cryptUint16E.Value()
			assert.Equal(t, ts.wantEncryptErr, err)
			if ts.wantEncryptErr == nil {
				err = cryptUint16D.Scan(uint16Val)
				assert.Equal(t, ts.wantDecryptErr, err)
				assert.Equal(t, uint16(ts.intVal), cryptUint16D.Val)
			}

			cryptUint32E.Val = uint32(ts.intVal)
			uint32Val, err := cryptUint32E.Value()
			assert.Equal(t, ts.wantEncryptErr, err)
			if ts.wantEncryptErr == nil {
				err = cryptUint32D.Scan(uint32Val)
				assert.Equal(t, ts.wantDecryptErr, err)
				assert.Equal(t, uint32(ts.intVal), cryptUint32D.Val)
			}

			cryptUint64E.Val = uint64(ts.intVal)
			uint64Val, err := cryptUint64E.Value()
			assert.Equal(t, ts.wantEncryptErr, err)
			if ts.wantEncryptErr == nil {
				err = cryptUint64D.Scan(uint64Val)
				assert.Equal(t, ts.wantDecryptErr, err)
				assert.Equal(t, uint64(ts.intVal), cryptUint64D.Val)
			}

			cryptUIntE.Val = uint(ts.intVal)
			uintVal, err := cryptUIntE.Value()
			assert.Equal(t, notSpecifyInt, err)
			if err == nil {
				err = cryptUIntD.Scan(uintVal)
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
	}{
		{
			name:           "简单",
			val:            "大明教你学go",
			wantEncryptErr: nil,
			wantDecryptErr: nil,
		},
		{
			name:           "长一点",
			val:            "大明教你学go,的撒放你家看你反饥饿卡三顿饭呢asdfjk3jrfnjkdnsafjknbjhbasdf阿斯顿发",
			wantEncryptErr: nil,
			wantDecryptErr: nil,
		},
		{
			name:           "空",
			val:            "",
			wantEncryptErr: nil,
			wantDecryptErr: nil,
		},
	}
	cryptStringE := NewEncryptColumn[string]("ABCDABCDABCDABCD")
	cryptStringD := NewEncryptColumn[string]("ABCDABCDABCDABCD")

	for _, ts := range testStringCases {
		t.Run(ts.name, func(t *testing.T) {
			cryptStringE.Val = ts.val
			value, err := cryptStringE.Value()
			assert.Equal(t, ts.wantEncryptErr, err)
			if err == nil {
				fmt.Println(string(reflect.ValueOf(value).Bytes()))
				err = cryptStringD.Scan(value)
				assert.Equal(t, ts.wantDecryptErr, err)
				assert.Equal(t, ts.val, cryptStringD.Val)
			}
		})
	}
}

func TestEncryptColumn_Struct(t *testing.T) {
	key := "ABCDABCDABCDABCD"

	cryptSimple := NewEncryptColumn[Simple](key)
	cryptSimple.Val = Simple{
		A: 1,
		B: 1.2,
		D: false,
	}
	val, err := cryptSimple.Value()
	fmt.Println(string(reflect.ValueOf(val).Bytes()))
	assert.Equal(t, nil, err)

	cryptSimpleD := NewEncryptColumn[Simple](key)

	err = cryptSimpleD.Scan(val)
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

	cryptComposite := NewEncryptColumn[Composite](key)
	cryptComposite.Val = composite
	val, err = cryptComposite.Value()
	fmt.Println(string(reflect.ValueOf(val).Bytes()))
	assert.Equal(t, nil, err)

	cryptCompositeD := NewEncryptColumn[Composite](key)
	err = cryptCompositeD.Scan(val)

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

func BenchmarkEncryptColumn_ValueStructNoCopy(b *testing.B) {
	key := "ABCDABCDABCDABCD"
	cryptSimple := NewEncryptColumn[Simple](key)

	for i := 0; i < b.N; i++ {
		_, _ = cryptSimple.Value()
	}
}

func BenchmarkEncryptColumn_ValueStructCopy(b *testing.B) {
	key := "ABCDABCDABCDABCD"
	for i := 0; i < b.N; i++ {
		cryptSimple := NewEncryptColumn[Simple](key)
		_, _ = cryptSimple.Value()
	}
}

func BenchmarkEncryptColumn_ScanStructNoCopy(b *testing.B) {
	key := "ABCDABCDABCDABCD"
	cryptSimple := NewEncryptColumn[Simple](key)
	cBytes, _ := cryptSimple.Value()

	decryptSimple := NewEncryptColumn[Simple](key)
	for i := 0; i < b.N; i++ {
		decryptSimple.Scan(cBytes)
	}
}

func BenchmarkEncryptColumn_ScanStructCopy(b *testing.B) {
	key := "ABCDABCDABCDABCD"
	cryptSimple := NewEncryptColumn[Simple](key)
	cBytes, _ := cryptSimple.Value()

	for i := 0; i < b.N; i++ {
		decryptSimple := NewEncryptColumn[Simple](key)
		decryptSimple.Scan(cBytes)
	}
}
