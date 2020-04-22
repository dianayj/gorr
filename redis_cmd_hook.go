package gorr

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/go-redis/redis"
)

// command hook
func statusCmdValue(cmd *redis.StatusCmd) string {
	value, err := getStoredValue(cmd.Args())
	if err != nil {
		return ""
	}
	return string(value)
}

func statusCmdResult(cmd *redis.StatusCmd) (string, error) {
	value, err := getStoredValue(cmd.Args())
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func stringCmdValue(cmd *redis.StringCmd) string {
	value, err := getStoredValue(cmd.Args())
	if err != nil {
		return ""
	}
	return string(value)
}

func stringSliceCmdValue(cmd *redis.StringSliceCmd) []string {
	ret, _ := stringSliceCmdResult(cmd)
	return ret
}

func stringSliceCmdResult(cmd *redis.StringSliceCmd) ([]string, error) {
	var ret []string
	value, err1 := getStoredValue(cmd.Args())
	if err1 != nil {
		return ret, errors.New("get value from db failed")
	}

	var buff bytes.Buffer
	sz, err2 := buff.Write(value)
	if err2 != nil || sz != len(value) {
		return ret, errors.New("write to buffer failed")
	}

	var slen int32
	err3 := binary.Read(&buff, binary.LittleEndian, &slen)
	if err3 != nil {
		return ret, errors.New("read size of slice from buffer failed")
	}

	sz -= 4

	for i := 0; i < int(slen); i++ {
		var sz2 int32
		err := binary.Read(&buff, binary.LittleEndian, &sz2)
		if err != nil || sz2 < 0 || int(sz2) > sz-4 {
			return ret, errors.New("read string size from buffer failed")
		}

		sz -= 4
		data := make([]byte, sz2)
		err = binary.Read(&buff, binary.LittleEndian, data)
		if err != nil {
			return ret, errors.New("read string from buffer failed")
		}

		sz -= int(sz2)
		ret = append(ret, string(data))
	}

	return ret, nil
}

func stringCmdResult(cmd *redis.StringCmd) (string, error) {
	value, err := getStoredValue(cmd.Args())
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func intCmdValue(cmd *redis.IntCmd) int64 {
	var buff bytes.Buffer
	value, err1 := getStoredValue(cmd.Args())

	if err1 != nil {
		// fmt.Printf("read int cmd from db failed, error:%s\n", err1.Error())
		return 0
	}

	sz, err2 := buff.Write(value)
	if err2 != nil || sz != len(value) {
		// fmt.Printf("read int cmd size failed, error:%s, sz:%d, len(value):%d\n", err2.Error(), sz, len(value))
		return 0.0
	}

	var ret int64
	err := binary.Read(&buff, binary.LittleEndian, &ret)
	if err != nil {
		// fmt.Printf("read int cmd failed, error:%s\n", err.Error())
		return 0
	}

	return ret
}

func floatCmdValue(cmd *redis.FloatCmd) float64 {
	var buff bytes.Buffer
	value, err1 := getStoredValue(cmd.Args())

	if err1 != nil {
		return 0.0
	}

	sz, err2 := buff.Write(value)
	if err2 != nil || sz != len(value) {
		return 0.0
	}

	var ret float64
	err := binary.Read(&buff, binary.LittleEndian, &ret)
	if err != nil {
		return 0.0
	}

	return ret
}

func stringStringMapCmdValue(cmd *redis.StringStringMapCmd) map[string]string {
	ret, _ := stringStringMapCmdResult(cmd)
	return ret
}

func stringStringMapCmdResult(cmd *redis.StringStringMapCmd) (map[string]string, error) {
	ret := make(map[string]string)
	value, err1 := getStoredValue(cmd.Args())
	if err1 != nil {
		return ret, errors.New("get value from db failed")
	}

	var buff bytes.Buffer
	sz, err2 := buff.Write(value)
	if err2 != nil || sz != len(value) {
		return ret, errors.New("write to buffer failed")
	}

	var mlen int32
	err3 := binary.Read(&buff, binary.LittleEndian, &mlen)
	if err3 != nil {
		return ret, errors.New("read size of map from buffer failed")
	}

	sz -= 4

	for i := 0; i < int(mlen); i++ {
		var sz2 int32
		err := binary.Read(&buff, binary.LittleEndian, &sz2)
		if err != nil || sz2 < 0 || int(sz2) > sz-4 {
			return ret, errors.New("read map key string size from buffer failed")
		}

		sz -= 4
		k := make([]byte, sz2)
		err = binary.Read(&buff, binary.LittleEndian, k)
		if err != nil {
			return ret, errors.New("read map key string from buffer failed")
		}
		sz -= int(sz2)

		var sz3 int32
		err = binary.Read(&buff, binary.LittleEndian, &sz3)
		if err != nil || sz3 < 0 || int(sz3) > sz-4 {
			return ret, errors.New("read map value string size from buffer failed")
		}

		sz -= 4
		v := make([]byte, sz3)
		err = binary.Read(&buff, binary.LittleEndian, v)
		if err != nil {
			return ret, errors.New("read map value string from buffer failed")
		}

		sz -= int(sz3)
		ret[string(k)] = string(v)
	}
	return ret, nil
}

func sliceCmdValue(cmd *redis.SliceCmd) []interface{} {
	ret, _ := sliceCmdResult(cmd)
	return ret
}

func sliceCmdResult(cmd *redis.SliceCmd) ([]interface{}, error) {
	var ret []interface{}
	value, err1 := getStoredValue(cmd.Args())
	if err1 != nil {
		return ret, errors.New("get value from db failed")
	}

	var buff bytes.Buffer
	sz, err2 := buff.Write(value)
	if err2 != nil || sz != len(value) {
		return ret, errors.New("write to buffer failed")
	}

	// 获取slice长度
	var slen int32
	err3 := binary.Read(&buff, binary.LittleEndian, &slen)
	if err3 != nil {
		return ret, errors.New("read size of slice from buffer failed")
	}

	sz -= 4

	// 获取元素类型长度
	var tlen int32
	err4 := binary.Read(&buff, binary.LittleEndian, &tlen)
	if err4 != nil {
		return ret, errors.New("read size of type from buffer failed")
	}
	sz -= 4

	// 获取元素类型
	tname := make([]byte, tlen)
	err5 := binary.Read(&buff, binary.LittleEndian, tname)
	if err5 != nil {
		return ret, errors.New("read type string from buffer failed")
	}

	sz -= int(tlen)

	// 依次获取元素
	for i := 0; i < int(slen); i++ {
		var sz2 int32
		err := binary.Read(&buff, binary.LittleEndian, &sz2)
		if err != nil || sz2 < 0 || int(sz2) > sz-4 {
			return ret, errors.New("read slice item size from buffer failed")
		}

		sz -= 4
		data := make([]byte, sz2)
		err = binary.Read(&buff, binary.LittleEndian, data)
		if err != nil {
			return ret, errors.New("read slice item from buffer failed")
		}

		// 根据类型获取值
		ud, _ := getValueByTypeName(string(tname), data)
		sz -= int(sz2)

		ret = append(ret, ud)
	}
	return ret, nil
}

func getValueByTypeName(name string, data []byte) (interface{}, error) {
	var err error
	switch name {
	case "string":
		var s string
		err = unmarshalValue(data, &s)
		return s, err
	case "int":
		var i int
		err = unmarshalValue(data, &i)
		return i, err
	case "int8":
		var i int8
		err = unmarshalValue(data, &i)
		return i, err
	case "int16":
		var i int16
		err = unmarshalValue(data, &i)
		return i, err
	case "int32":
		var i int32
		err = unmarshalValue(data, &i)
		return i, err
	case "int64":
		var i int64
		err = unmarshalValue(data, &i)
		return i, err
	case "uint":
		var ui uint
		err = unmarshalValue(data, &ui)
		return ui, err
	case "uint8":
		var ui uint8
		err = unmarshalValue(data, &ui)
		return ui, err
	case "uint16":
		var ui uint16
		err = unmarshalValue(data, &ui)
		return ui, err
	case "uint32":
		var ui uint32
		err = unmarshalValue(data, &ui)
		return ui, err
	case "uint64":
		var ui uint64
		err = unmarshalValue(data, &ui)
		return ui, err
	case "float32":
		var f float32
		err = unmarshalValue(data, &f)
		return f, err
	case "float64":
		var f float64
		err = unmarshalValue(data, &f)
		return f, err
	case "complex64":
		var c complex64
		err = unmarshalValue(data, &c)
		return c, err
	case "complex128":
		var c complex128
		err = unmarshalValue(data, &c)
		return c, err
	case "uintptr":
		var up uintptr
		err = unmarshalValue(data, &up)
		return up, err
	case "bool":
		var b bool
		err = unmarshalValue(data, &b)
		return b, err
	default:
		return nil, errors.New("unsupported type of slice item")
	}
}
