package main
 
import (
	"golang.org/x/net/websocket"
	// "encoding/json"
	// "log"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"container/list"
	"sync"
	"strconv"
	"errors"
	// "time"
	// "fmt"
)

/*
#include "client.h"

static void luatc_newtable(lua_State* L) {
	lua_newtable(L);
}
*/
import "C"

var websocketInstance *websocket.Conn;
// type KeyInput struct {
// 	T int `json:"t"`
// 	K int `json:"k"`
// }

var lastErr error
var remoteEventListLock *sync.Mutex
var remoteEventList = list.New()
// var listSize int = 0
// var lastRemoteEvent KeyInput

func eventPush(evt string) {
	remoteEventListLock.Lock()
	defer remoteEventListLock.Unlock()
	remoteEventList.PushBack(evt)
	// listSize++
}

//export luatc_connect
func luatc_connect(L *C.lua_State) C.int {
	if luaTypeOf(L, 1) != luaTypeString {
		luaStackTopSet(L, 0)
		luaNilPush(L)
		luaStringPush(L, "no url provided")
	}
	origin := "http://127.0.0.1/"  //客户端地址
	// url := "ws://120.79.44.197:8966/"  //服务器地址
	url := luaStringGet(L, 1)
	ws, err := websocket.Dial(url, "", origin) //第二个参数是websocket子协议，可以为空
	websocketInstance = ws
	remoteEventListLock = new(sync.Mutex)

	go func() {
		var buf= make([]byte, 100)
		// var evt KeyInput
		for {
			n, err := ws.Read(buf)
			if err != nil {
				lastErr = err
				// log.Fatal(err)
			}
			messageString := string(buf[:n])
			eventPush(messageString)
			// jsonValue = gjson.Parse()
			// evt := KeyInput{}
			// fmt.Printf("receive: %v\n", string(buf[:n]))
			// err = json.Unmarshal(buf[:n], &evt)
			if err != nil {
				lastErr = err
				// log.Fatal(err)
			}
			
			// fmt.Printf("recv: %+v\n", lastRemoteEvent)
		}
	}()

	if err != nil {
		luaStackTopSet(L, 0)
		luaNilPush(L)
		luaStringPush(L, err.Error())
		return C.int(2)
	}

	luaStackTopSet(L, 0)
	luaIntegerPush(L, 1)
	luaNilPush(L)
	return C.int(2)
}

func luaNewTable(L *C.lua_State) {
	C.luatc_newtable(L)
}

func luaTableStringInt(L *C.lua_State, key string, value int) {
	luaStringPush(L, key)
	luaIntegerPush(L, value)
	C.lua_settable(L, -3)
}

func luaTableStringString(L *C.lua_State, key string, value string) {
	luaStringPush(L, key)
	luaStringPush(L, value)
	C.lua_settable(L, -3)
}

// func luaTablePushIntKey(L *C.lua_State, index int, t int, k int) {
// 	luaIntegerPush(L, index)

// 	C.luatc_newtable(L)
// 	luaTablePushStringInt(L, "t", t)
// 	luaTablePushStringInt(L, "k", k)

// 	C.lua_settable(L, -3)
// }

//export luatc_read
func luatc_read(L *C.lua_State) C.int {

	if lastErr != nil {
		luaStackTopSet(L, 0)
		luaNilPush(L)
		luaStringPush(L, lastErr.Error())
		return C.int(2)
	}

	remoteEventListLock.Lock()
	defer remoteEventListLock.Unlock()

	luaNewTable(L)

	i := 1
	var next *list.Element
    for e := remoteEventList.Front(); e != nil; e = next {
		messageString, ok := e.Value.(string)
		if ok {
			// key
			luaIntegerPush(L, i)
			i++

			// value
			luaNewTable(L)
			jsonValue := gjson.Parse(messageString)
			jsonValue.ForEach(func (key gjson.Result, value gjson.Result) bool {
				itemKey := key.String()
				intValue := int(value.Int())
				strValue := value.String()
				conValue := strconv.Itoa(intValue)
				// fmt.Printf("<%v><%v>\n", conValue, strValue)
				if conValue == strValue {
					luaTableStringInt(L, itemKey, intValue)
				} else {
					luaTableStringString(L, itemKey, strValue)
				}
				return true
			})

			// close table item
			C.lua_settable(L, -3)
		} else {
			lastErr = errors.New("message not string")
		}

        next = e.Next()
        remoteEventList.Remove(e)
	}

	if lastErr == nil {
		luaNilPush(L)
	} else {
		luaStringPush(L, lastErr.Error())
	}

	return C.int(2)
}

//export luatc_write
func luatc_write(L *C.lua_State) C.int {
	if websocketInstance == nil {
		luaStackTopSet(L, 0)
		luaStringPush(L, "not connected")
		return C.int(1)
	}

	luaNilPush(L)

	var invald = false
	// var evt KeyInput

	var jsonStr string = "{}"
	var jsonErr error = nil

	for int(C.lua_next(L, C.int(-2))) != 0 {
		tKey := int(C.lua_type(L, C.int(-2)))
		tValue := int(C.lua_type(L, C.int(-1)))

		if tKey == 4 {
			keyStr := luaStringGet(L, -2)
			if tValue == 3 {
				valueNum := int(C.lua_tonumber(L, -1))
				jsonStr, jsonErr = sjson.Set(jsonStr, keyStr, valueNum)
			} else if tValue == 4 {
				valueStr := luaStringGet(L, -1)
				jsonStr, jsonErr = sjson.Set(jsonStr, keyStr, valueStr)
			} else {
				invald = true
			}
		} else {
			invald = true
		}

		// switch tKey {
		// case 3:
		// 	keyNum := C.int(C.lua_tonumber(L, -2))
		// 	fmt.Printf("key=%v", keyNum)
		// 	break
		// case 4:
		// 	keyStr := luaStringGet(L, -2)
		// 	fmt.Printf("key=%v", keyStr)
		// 	break
		// default:
		// 	break
		// }

		// switch tValue {
		// case 3:
		// 	valueNum := C.int(C.lua_tonumber(L, -1))
		// 	fmt.Printf("value=%v", valueNum)
		// 	break
		// case 4:
		// 	valueStr := luaStringGet(L, -1)
		// 	fmt.Printf("value=%v", valueStr)
		// 	break
		// default:
		// 	break
		// }
		luaStackPop(L, 1)
	}

	
	if invald {
		luaStackTopSet(L, 0)
		luaStringPush(L, "invalid input")
		return C.int(1)
	}

	if jsonErr != nil {
		luaStackTopSet(L, 0)
		luaStringPush(L, jsonErr.Error())
		return C.int(1)
	}

	// data, err := json.Marshal(evt)
	// if err != nil {
	// 	luaStackTopSet(L, 0)
	// 	luaStringPush(L, err.Error())
	// 	return C.int(1)
	// }

	_, writeErr := websocketInstance.Write([]byte(jsonStr))

	if writeErr != nil {
		luaStackTopSet(L, 0)
		luaStringPush(L, writeErr.Error())
		return C.int(1)
	}

	luaStackTopSet(L, 0)
	luaNilPush(L)
	return C.int(1)
}
