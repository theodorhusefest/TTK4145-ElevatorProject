package Hardware

/*
   Translating C-functions to GOlang
*/
import "C"

func io_init() {
	return int(C.io_init())
}

func io_set_bit(channel int) {
	C.io_set_bit(C.int(channel))
}

func io_clear_bit(channel int) {
	C.io_clear_bit(C.int(channel))
}

func io_write_analog(channel int, value int) {
	C.io_write_analog(C.int(channel), C.int(value))
}

func io_read_bit(channel int) {
	return int(C.io_read_channel())
}

func io_read_analog(channel int) {
	return int(C.io_read_analog())
}
