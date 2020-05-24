package main

/*
#cgo pkg-config: luajit
#include "client.h"

LUALIB_API int luaopen_client(lua_State* L) {
	luaL_Reg regs[] = {
		{ "poll", luatc_poll },
		{ "httpraw", luatc_httpraw },
		{ "connect", luatc_connect },
		{ "read", luatc_read },
		{ "write", luatc_write },
		{ NULL, NULL },
	};
	luaL_newlib(L, regs);
	return 1;
}
*/
import "C"

// main is the pseudo main function that will be simply
// ignored in the buildmode c-shared.
func main() {
}
