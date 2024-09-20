package vaccs

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type Vaccs_file_reader struct {
	arch_type    uint8
	path         string
	num_records  int64
	data         []byte
	position     uint64
	size         int64
	first_record bool
}

func (v *Vaccs_file_reader) GetPath() string {
	return v.path
}

func (v *Vaccs_file_reader) GetNumRecords() int64 {
	return v.num_records
}

func (v *Vaccs_file_reader) AddPath(p string) *Vaccs_file_reader {
	v.path = p
	return v
}

func (v *Vaccs_file_reader) GetVaccsAnalysisRecords() ([]Vaccs_record, error) {
	file, err := os.Open(v.path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	v.size = stats.Size()
	v.data = make([]byte, v.size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(v.data)

	if err != nil {
		return nil, err
	}

	vaccs_records := make([]Vaccs_record, 0)

	v.first_record = true
	var rec Vaccs_record
	for rec, err = v.getNextRecord(); rec != nil && err == nil; rec, err = v.getNextRecord() {
		vaccs_records = append(vaccs_records, rec)
	}

	return vaccs_records, err
}

func (v *Vaccs_file_reader) getBytes(num_bytes uint64) ([]uint8, error) {
	if v.position+num_bytes <= uint64(v.size) {
		buffer := make([]uint8, num_bytes)
		for i := 0; i < len(buffer); i++ {
			(buffer)[i] = v.data[v.position]
			v.position++
		}
		return buffer, nil
	} else {
		msg := fmt.Sprintf("Malformed Input: bytes to read = %d, bytes left = %d\n", num_bytes, uint64(v.size)-v.position+1)
		return nil, &Vaccs_error{msg}
	}
}

func getUint[T uint8 | uint16 | uint32 | uint64](v *Vaccs_file_reader, num_bytes uint64) (T, error) {
	data, err := v.getBytes(num_bytes)

	if err == nil {
		var val T
		b := bytes.NewReader(data)
		err = binary.Read(b, binary.LittleEndian, &val)
		if err == nil {
			return val, nil
		}
	}

	return 0, err
}

func (v *Vaccs_file_reader) getLong() (uint64, error) {
	var val uint64
	var err error

	if v.arch_type == X86_64 {
		val, err = getUint[uint64](v, UINT64_BYTES)
	} else {
		var tval uint32
		tval, err = getUint[uint32](v, UINT32_BYTES)
		val = uint64(tval)
	}
	return val, err
}

func (v *Vaccs_file_reader) getString() (string, error) {
	length, err := v.getLong()

	if err == nil {
		var bytes []byte
		bytes, err = v.getBytes(length)

		return string(bytes), err
	}
	return "", err
}

func (v *Vaccs_file_reader) getNextRecord() (Vaccs_record, error) {
	if v.position >= uint64(v.size) {
		return nil, nil
	}

	r_type, err := getUint[uint16](v, UINT16_BYTES)

	if err == nil {
		v.num_records++
		if v.first_record && r_type != VACCS_ARCH {
			return nil, &Vaccs_error{v.path + ": Malformed input, first record must identify the architecture type"}
		} else {
			v.first_record = false
		}

		switch r_type {
		case VACCS_ARCH:
			var vr json_arch
			vr, err = v.generate_vaccs_architecture()
			if err == nil {
				return vr, nil
			}
		case VACCS_ASM:
			var vr json_assembly
			vr, err = v.generate_vaccs_assembly()
			if err == nil {
				return vr, nil
			}
		case VACCS_BINARY:
			var vr json_binary
			vr, err = v.generate_vaccs_binary()
			if err == nil {
				return vr, nil
			}
		case VACCS_CCODE:
			var vr json_c_source
			vr, err = v.generate_vaccs_ccode()
			if err == nil {
				return vr, nil
			}
		case VACCS_CMD_LINE:
			var vr json_command_line
			vr, err = v.generate_vaccs_command_line()
			if err == nil {
				return vr, nil
			}
		case VACCS_FILE_CLOSE:
			var vr json_file_close
			vr, err = v.generate_vaccs_file_close()
			if err == nil {
				return vr, nil
			}
		case VACCS_FILE_OPEN:
			var vr json_file_open
			vr, err = v.generate_vaccs_file_open()
			if err == nil {
				return vr, nil
			}
		case VACCS_FILE_READ:
			var vr json_file_read
			vr, err = v.generate_vaccs_file_read()
			if err == nil {
				return vr, nil
			}
		case VACCS_FILE_WRITE:
			var vr json_file_write
			vr, err = v.generate_vaccs_file_write()
			if err == nil {
				return vr, nil
			}
		case VACCS_FREE:
			var vr json_free
			vr, err = v.generate_vaccs_free()
			if err == nil {
				return vr, nil
			}
		case VACCS_FUNCTION_INV:
			var vr json_function_invocation
			vr, err = v.generate_vaccs_function_invocation()
			if err == nil {
				return vr, nil
			}
		case VACCS_LINE_CHANGE:
			var vr json_line_change
			vr, err = v.generate_vaccs_line_change()
			if err == nil {
				return vr, nil
			}
		case VACCS_MALLOC:
			var vr json_malloc
			vr, err = v.generate_vaccs_malloc()
			if err == nil {
				return vr, nil
			}
		case VACCS_OUTPUT:
			var vr json_output
			vr, err = v.generate_vaccs_output()
			if err == nil {
				return vr, nil
			}
		case VACCS_REGISTER:
			var vr json_register
			vr, err = v.generate_vaccs_register()
			if err == nil {
				return vr, nil
			}
		case VACCS_RETURN_ADDR:
			var vr json_return_address
			vr, err = v.generate_vaccs_return_address()
			if err == nil {
				return vr, nil
			}
		case VACCS_RETURN:
			var vr json_return
			vr, err = v.generate_vaccs_return()
			if err == nil {
				return vr, nil
			}
		case VACCS_SECTION:
			var vr json_section
			vr, err = v.generate_vaccs_section()
			if err == nil {
				return vr, nil
			}
		case VACCS_VAR_ACCESS:
			var vr json_variable_access
			vr, err = v.generate_vaccs_variable_access()
			if err == nil {
				return vr, nil
			}
		default:
			msg := fmt.Sprintf("Invalid record type: %d\n", r_type)
			return nil, &Vaccs_error{msg}
		}
		return nil, err
	}

	return nil, err
}

func (v *Vaccs_file_reader) get_common_data() (int64, uint16, string, error) {
	var tevn uint64
	var err error

	tevn, err = v.getLong()
	if err == nil {
		event_num := int64(tevn)

		var c_line_num uint16
		c_line_num, err = getUint[uint16](v, UINT16_BYTES)
		if err == nil {
			var path string
			path, err = v.getString()
			if err == nil {
				return event_num, c_line_num, path, nil
			}
		}
	}

	return 0, 0, "", err
}

func (v *Vaccs_file_reader) generate_vaccs_architecture() (json_arch, error) {
	var arch_type uint8
	var err error

	arch_type, err = getUint[uint8](v, UINT8_BYTES)
	v.arch_type = arch_type

	if err == nil {
		var heap_start uint64
		heap_start, err = v.getLong()
		if err == nil {
			var heap_end uint64
			heap_end, err = v.getLong()
			if err == nil {
				var stack_start uint64
				stack_start, err = v.getLong()
				if err == nil {
					var stack_end uint64
					stack_end, err = v.getLong()
					if err == nil {
						return json_arch{&Vaccs_architecture{NO_EVENT, arch_type,
							heap_start, heap_end, stack_start, stack_end}}, nil
					}
				}
			}
		}
	}

	return json_arch{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_assembly() (json_assembly, error) {
	var asm_line_num uint16
	var err error

	asm_line_num, err = getUint[uint16](v, UINT16_BYTES)

	if err == nil {
		var c_line_num uint16
		c_line_num, err = getUint[uint16](v, UINT16_BYTES)
		if err == nil {
			var asm_file_name string
			asm_file_name, err = v.getString()
			if err == nil {
				var c_file_name string
				c_file_name, err = v.getString()
				if err == nil {
					var inst string
					inst, err = v.getString()
					if err == nil {
						return json_assembly{&Vaccs_assembly{NO_EVENT, asm_line_num,
							c_line_num, asm_file_name, c_file_name, inst}}, nil
					}
				}
			}
		}
	}

	return json_assembly{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_binary() (json_binary, error) {
	var path string
	var err error

	path, err = v.getString()
	if err == nil {
		return json_binary{&Vaccs_binary{NO_EVENT, path}}, nil
	}

	return json_binary{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_ccode() (json_c_source, error) {
	var c_line_num uint16
	var err error

	c_line_num, err = getUint[uint16](v, UINT16_BYTES)

	if err == nil {
		var start_pos uint16

		start_pos, err = getUint[uint16](v, UINT16_BYTES)
		if err == nil {
			var path string

			path, err = v.getString()
			if err == nil {
				var source_line string

				source_line, err = v.getString()
				if err == nil {
					return json_c_source{&Vaccs_c_source{NO_EVENT, c_line_num, path, start_pos, source_line}}, nil
				}
			}
		}
	}

	return json_c_source{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_command_line() (json_command_line, error) {
	var cmd_line string
	var err error

	cmd_line, err = v.getString()
	if err == nil {
		return json_command_line{&Vaccs_command_line{NO_EVENT, cmd_line}}, nil
	}

	return json_command_line{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_file_open() (json_file_open, error) {

	event_num, c_line_num, path, err := v.get_common_data()

	if err == nil {
		var file_descriptor uint64
		file_descriptor, err = v.getLong()
		if err == nil {
			var flags uint32
			flags, err = getUint[uint32](v, UINT32_BYTES)
			if err == nil {
				var mode uint32
				mode, err = getUint[uint32](v, UINT32_BYTES)
				if err == nil {
					var tsl uint8
					tsl, err = getUint[uint8](v, UINT8_BYTES)
					var is_sym_link bool = false
					if tsl == 1 {
						is_sym_link = true
					}
					if err == nil {
						return json_file_open{&Vaccs_file_open{event_num, c_line_num, path,
							file_descriptor, flags, mode, is_sym_link}}, nil
					}
				}
			}
		}
	}

	return json_file_open{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_file_close() (json_file_close, error) {

	event_num, c_line_num, path, err := v.get_common_data()

	if err == nil {
		var file_descriptor uint64
		file_descriptor, err = v.getLong()
		if err == nil {
			return json_file_close{&Vaccs_file_close{event_num, c_line_num, path, file_descriptor}}, nil
		}
	}

	return json_file_close{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_file_read() (json_file_read, error) {

	event_num, c_line_num, path, err := v.get_common_data()

	if err == nil {
		var file_descriptor uint64
		file_descriptor, err = v.getLong()
		if err == nil {
			var tv uint8
			tv, err = getUint[uint8](v, UINT8_BYTES)
			is_valid := true
			if tv == 0 {
				is_valid = false
			}
			if err == nil {
				return json_file_read{&Vaccs_file_read{event_num, c_line_num, path, file_descriptor, is_valid}}, nil
			}
		}
	}
	return json_file_read{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_file_write() (json_file_write, error) {

	event_num, c_line_num, path, err := v.get_common_data()

	if err == nil {
		var file_descriptor uint64
		file_descriptor, err = v.getLong()
		if err == nil {
			var tv uint8
			tv, err = getUint[uint8](v, UINT8_BYTES)
			is_valid := true
			if tv == 0 {
				is_valid = false
			}
			if err == nil {
				return json_file_write{&Vaccs_file_write{event_num, c_line_num, path, file_descriptor, is_valid}}, nil
			}
		}
	}

	return json_file_write{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_free() (json_free, error) {

	event_num, c_line_num, path, err := v.get_common_data()

	if err == nil {
		var address uint64
		address, err = v.getLong()
		if err == nil {
			return json_free{&Vaccs_free{event_num, c_line_num, path, address}}, nil
		}
	}

	return json_free{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_function_invocation() (json_function_invocation, error) {

	event_num, inv_line_num, inv_path, err := v.get_common_data()

	if err == nil {
		var func_line_num uint16
		func_line_num, err = getUint[uint16](v, UINT16_BYTES)
		if err == nil {
			var func_path string
			func_path, err = v.getString()
			if err == nil {
				var func_name string
				func_name, err = v.getString()
				if err == nil {
					var address uint64
					address, err = getUint[uint64](v, UINT64_BYTES)
					if err == nil {
						return json_function_invocation{&Vaccs_function_invocation{event_num, inv_line_num, inv_path,
							func_line_num, func_path, func_name, address}}, nil
					}
				}
			}
		}
	}

	return json_function_invocation{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_line_change() (json_line_change, error) {

	event_num, c_line_num, path, err := v.get_common_data()

	if err == nil {
		return json_line_change{&Vaccs_line_change{event_num, c_line_num, path}}, nil
	}

	return json_line_change{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_malloc() (json_malloc, error) {

	event_num, c_line_num, path, err := v.get_common_data()

	if err == nil {
		var address uint64
		address, err = v.getLong()
		if err == nil {
			var num_bytes uint64
			num_bytes, err = v.getLong()
			if err == nil {
				return json_malloc{&Vaccs_malloc{event_num, c_line_num, path, address, num_bytes}}, nil
			}
		}
	}

	return json_malloc{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_output() (json_output, error) {
	var output string
	var err error

	output, err = v.getString()
	if err == nil {
		return json_output{&Vaccs_output{NO_EVENT, output}}, nil
	}

	return json_output{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_register() (json_register, error) {

	event_num, c_line_num, path, err := v.get_common_data()

	if err == nil {
		var register string
		register, err = v.getString()
		if err == nil {
			var tv uint64
			tv, err = v.getLong()
			var value int64 = int64(tv)
			if err == nil {
				return json_register{&Vaccs_register{event_num, c_line_num, path, register, value}}, nil
			}
		}
	}

	return json_register{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_return_address() (json_return_address, error) {
	var err error
	var return_address uint64

	return_address, err = getUint[uint64](v, UINT64_BYTES)
	if err == nil {
		var dynamic_link uint64

		dynamic_link, err = getUint[uint64](v, UINT64_BYTES)
		if err == nil {
			var func_name string

			func_name, err = v.getString()
			if err == nil {
				return json_return_address{&Vaccs_return_address{NO_EVENT, dynamic_link, return_address, func_name}}, nil
			}
		}
	}

	return json_return_address{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_return() (json_return, error) {
	var event_num uint64
	var err error

	event_num, err = v.getLong()
	if err == nil {
		return json_return{&Vaccs_return{int64(event_num)}}, nil
	}

	return json_return{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_section() (json_section, error) {
	var section string
	var err error

	section, err = v.getString()
	if err == nil {
		return json_section{&Vaccs_section{NO_EVENT, section}}, nil
	}

	return json_section{nil}, err
}

func (v *Vaccs_file_reader) generate_vaccs_variable_access() (json_variable_access, error) {

	event_num, c_line_num, path, err := v.get_common_data()

	if err == nil {
		var scope string
		scope, err = v.getString()
		if err == nil {
			var address uint64
			address, err = v.getLong()
			if err == nil {
				var var_type string
				var_type, err = v.getString()
				if err == nil {
					var var_name string
					var_name, err = v.getString()
					if err == nil {
						var points_to uint64
						points_to, err = v.getLong()
						if err == nil {
							var value string
							value, err = v.getString()
							if err == nil {
								return json_variable_access{&Vaccs_variable_access{event_num, c_line_num, path,
									scope, address, var_type, var_name, points_to, value}}, nil
							}
						}
					}
				}
			}
		}
	}

	return json_variable_access{nil}, err
}
