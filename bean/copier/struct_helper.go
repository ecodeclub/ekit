package copier

import (
	"errors"
	"reflect"
)

type structHelper struct {
	ptrOffset uintptr
	typ       reflect.Type
}

type structOffsets struct {
	helper []structHelper
	//Slice、Array、Map等需要深度拷贝的
	deepCopyOffsets map[reflect.Kind][]uintptr
}

var structOffsetsMapGlobal = make(map[string]*structOffsets)

var errorNotStruct = errors.New("the input must be struct")

func findOffsets(typ reflect.Type, offsets *structOffsets, lastOffset uintptr, structOffsetsMap map[string]*structOffsets) {
	kind := typ.Kind()
	if kind != reflect.Struct {
		//只存在于内部调用的情况，所以直接panic
		panic("findOffsets的reflect.Value类型一定为Struct")
	}
	name := typ.Name()
	// 若之前已经存过，直接加上offset后返回
	if storeStructOffsets, ok := structOffsetsMap[name]; ok {
		for _, subOffsets := range storeStructOffsets.helper {
			//对于上一层的struct，需要添加相应的偏移量
			offsets.helper = append(offsets.helper, structHelper{
				subOffsets.ptrOffset + lastOffset,
				subOffsets.typ,
			})
		}
		return
	}
	newStructOffsets := &structOffsets{
		helper: make([]structHelper, 0, 0),
	}
	for i := 0; i < typ.NumField(); i++ {
		subType := typ.Field(i).Type
		switch subType.Kind() {
		case reflect.Pointer, reflect.UnsafePointer:
			sSubType := subType.Elem()
			//指针解下去，直到找到一个struct
			newStructOffsets.helper = append(newStructOffsets.helper, structHelper{
				typ.Field(i).Offset,
				sSubType,
			})
			for sSubType.Kind() == reflect.Pointer || sSubType.Kind() == reflect.UnsafePointer {
				sSubType = sSubType.Elem()
			}
			if sSubType.Kind() == reflect.Struct {
				//指针类型结构体的具体内存分布和原结构体没有任何关系，所以要新建一个
				newSubStructOffsets := &structOffsets{
					helper: make([]structHelper, 0, 0),
				}
				findOffsets(sSubType, newSubStructOffsets, 0, structOffsetsMap)
			}
		case reflect.Struct:
			findOffsets(subType, newStructOffsets, typ.Field(i).Offset, structOffsetsMap)
		default:
			if isDeepCopyKind(subType.Kind()) {
				delayCreateOffsets(newStructOffsets, subType.Kind())
				delayCreateOffsets(offsets, subType.Kind())
				newStructOffsets.deepCopyOffsets[subType.Kind()] = append(newStructOffsets.deepCopyOffsets[subType.Kind()], typ.Field(i).Offset)
				offsets.deepCopyOffsets[subType.Kind()] = append(offsets.deepCopyOffsets[subType.Kind()], typ.Field(i).Offset+lastOffset)
			}
			continue
		}
	}
	for _, subOffsets := range newStructOffsets.helper {
		//对于上一层的struct，需要添加相应的偏移量
		offsets.helper = append(offsets.helper, structHelper{
			subOffsets.ptrOffset + lastOffset,
			subOffsets.typ,
		})
	}
	structOffsetsMap[name] = newStructOffsets
}

// 滞后的创建，并返回
func delayCreateOffsets(offset *structOffsets, kind reflect.Kind) {
	if offset.deepCopyOffsets == nil {
		offset.deepCopyOffsets = make(map[reflect.Kind][]uintptr)
	}
	if offset.deepCopyOffsets[kind] == nil {
		offset.deepCopyOffsets[kind] = make([]uintptr, 0, 1)
	}
}

// array| chan | map | slice
const deepCopyKind = (1 << 17) | (1 << 18) | (1 << 21) | (1 << 23)

func isDeepCopyKind(kind reflect.Kind) bool {
	return 1<<kind&deepCopyKind > 0
}

// 会自动维护一个全局的map,用于查询
func FindOffsetsDefault(inStruct any) (*structOffsets, error) {
	return findTypeOffsets(reflect.TypeOf(inStruct), structOffsetsMapGlobal)
}

func findValueOffsetsDefault(refVal reflect.Value) (*structOffsets, error) {
	return findTypeOffsets(refVal.Type(), structOffsetsMapGlobal)
}

func FindOffsets(inStruct any, structOffsetsMap map[string]*structOffsets) (*structOffsets, error) {
	return findTypeOffsets(reflect.TypeOf(inStruct), structOffsetsMap)
}

func findTypeOffsets(typ reflect.Type, structOffsetsMap map[string]*structOffsets) (*structOffsets, error) {
	if typ.Kind() != reflect.Struct {
		return nil, errorNotStruct
	}
	name := typ.Name()
	if storeStructOffsets, ok := structOffsetsMap[name]; ok {
		return storeStructOffsets, nil
	}
	newStructOffsets := &structOffsets{
		helper: make([]structHelper, 0, 0),
	}
	findOffsets(typ, newStructOffsets, 0, structOffsetsMap)
	return newStructOffsets, nil
}
