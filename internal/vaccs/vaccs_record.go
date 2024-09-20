package vaccs

import "encoding/json"

const VACCS_ARCH uint16 = 0         /* Information regarding target architecture */
const VACCS_ASM uint16 = 1          /* A line of assembly code */
const VACCS_BINARY uint16 = 2       /* Path to the executable */
const VACCS_CCODE uint16 = 3        /* Path to a C file */
const VACCS_CMD_LINE uint16 = 4     /* The command line for an execution */
const VACCS_DATA_REC uint16 = 5     /* The data records section of the file */
const VACCS_EVENT_REC uint16 = 6    /* The event records section of the file */
const VACCS_FUNCTION_INV uint16 = 7 /* A function invocation event */
const VACCS_OUTPUT uint16 = 8       /* A line of C ouptut */
const VACCS_REGISTER uint16 = 9     /* A register update event */
const VACCS_RETURN uint16 = 10      /* A function return event */
const VACCS_RETURN_ADDR uint16 = 11 /* The return address of a function invocation */
const VACCS_SD_CLEAR uint16 = 12    /* A sensitive data clear event */
const VACCS_SD_CZERO uint16 = 13    /* A sensitive data set core size to 0 event */
const VACCS_SD_LOCK uint16 = 14     /* A sensitive data mlock event */
const VACCS_SD_SET uint16 = 15      /* A sensitive data mset event */
const VACCS_SD_UNLOCK uint16 = 16   /* A sensitive data munlock event */
const VACCS_SECTION uint16 = 17     /* A data section in the binary */
const VACCS_VAR_ACCESS uint16 = 18  /* A variable access event */
const VACCS_MALLOC uint16 = 19      /* a malloc record */
const VACCS_FREE uint16 = 20        /* a free record */
const VACCS_FILE_OPEN uint16 = 21   /* a free record */
const VACCS_FILE_CLOSE uint16 = 22  /* a free record */
const VACCS_FILE_READ uint16 = 23   /* a free record */
const VACCS_FILE_WRITE uint16 = 24  /* a free record */
const VACCS_LINE_CHANGE uint16 = 25 /* a change in C line being executed */

const UINT8_BYTES uint64 = 1
const UINT16_BYTES uint64 = 2
const UINT32_BYTES uint64 = 4
const UINT64_BYTES uint64 = 8

const I386 uint8 = 0
const X86_64 uint8 = 1

const TypeSeparator string = "~!~"
const FieldSeparator string = "|"
const NULL_ADDRESS_STRING string = "--------"

const NO_EVENT int64 = -1

type Vaccs_record interface {
	get_json() ([]byte, error)
}

type Vaccs_architecture struct {
	Event_num   int64
	Arch_type   uint8
	Heap_start  uint64
	Heap_end    uint64
	Stack_start uint64
	Stack_end   uint64
}

type json_arch struct {
	VACCS_ARCH *Vaccs_architecture
}

type Vaccs_assembly struct {
	Event_num     int64
	Asm_line_num  uint16
	C_line_num    uint16
	Asm_file_name string
	C_file_name   string
	Inst          string
}

type json_assembly struct {
	VACCS_ASM *Vaccs_assembly
}

type Vaccs_binary struct {
	Event_num int64
	Path      string
}

type json_binary struct {
	VACCS_BINARY *Vaccs_binary
}

type Vaccs_c_source struct {
	Event_num   int64
	Line_num    uint16
	Path        string
	C_start_pos uint16
	Source      string
}

type json_c_source struct {
	VACCS_CCODE *Vaccs_c_source
}

type Vaccs_command_line struct {
	Event_num int64
	Cmd_line  string
}

type json_command_line struct {
	VACCS_CMD_LINE *Vaccs_command_line
}

type Vaccs_file_close struct {
	Event_num       int64
	Line_num        uint16
	Path            string
	File_descriptor uint64
}

type json_file_close struct {
	VACCS_FILE_CLOSE *Vaccs_file_close
}

type Vaccs_file_open struct {
	Event_num       int64
	Line_num        uint16
	Path            string
	File_descriptor uint64
	Flags           uint32
	Mode            uint32
	Is_sym_link     bool
}

type json_file_open struct {
	VACCS_FILE_OPEN *Vaccs_file_open
}

type Vaccs_file_read struct {
	Event_num       int64
	Line_num        uint16
	Path            string
	File_descriptor uint64
	Valid           bool
}

type json_file_read struct {
	VACCS_FILE_READ *Vaccs_file_read
}

type Vaccs_file_write struct {
	Event_num       int64
	Line_num        uint16
	Path            string
	File_descriptor uint64
	Valid           bool
}

type json_file_write struct {
	VACCS_FILE_WRITE *Vaccs_file_write
}

type Vaccs_free struct {
	Event_num int64
	Line_num  uint16
	Path      string
	Address   uint64
}

type json_free struct {
	VACCS_FREE *Vaccs_free
}

type Vaccs_function_invocation struct {
	Event_num    int64
	Line_num     uint16
	Path         string
	Inv_line_num uint16
	C_inv_path   string
	Func_name    string
	Address      uint64
}

type json_function_invocation struct {
	VACCS_FUNCTION_INV *Vaccs_function_invocation
}

type Vaccs_line_change struct {
	Event_num int64
	Line_num  uint16
	Path      string
}

type json_line_change struct {
	VACCS_LINE *Vaccs_line_change
}

type Vaccs_malloc struct {
	Event_num int64
	Line_num  uint16
	Path      string
	Address   uint64
	Num_bytes uint64
}

type json_malloc struct {
	VACCS_MALLOC *Vaccs_malloc
}

type Vaccs_output struct {
	Event_num int64
	Output    string
}

type json_output struct {
	VACCS_OUTPUT *Vaccs_output
}

type Vaccs_register struct {
	Event_num int64
	Line_num  uint16
	Path      string
	Register  string
	Value     int64
}

type json_register struct {
	VACCS_REGISTER *Vaccs_register
}

type Vaccs_return_address struct {
	Event_num      int64
	Dynamic_link   uint64
	Return_address uint64
	Function_name  string
}

type json_return_address struct {
	VACCS_RETURN_ADDR *Vaccs_return_address
}

type Vaccs_return struct {
	Event_num int64
}

type json_return struct {
	VACCS_RETURN *Vaccs_return
}

type Vaccs_section struct {
	Event_num int64
	Info      string
}

type json_section struct {
	VACCS_SECTION *Vaccs_section
}

type Vaccs_variable_access struct {
	Event_num int64
	Line_num  uint16
	Path      string
	Scope     string /* the function name for the scope of the variable */
	Address   uint64 /* the variable address */
	Var_type  string /* the variable type */
	Name      string /* the name of the variable */
	Points_to uint64 /* the address pointed to if a pointer */
	Value     string /* the variable value */
}

type json_variable_access struct {
	VACCS_VAR_ACCESS *Vaccs_variable_access
}

func (ja json_arch) get_json() ([]byte, error) {
	return json.Marshal(ja)
}

func (ja json_assembly) get_json() ([]byte, error) {
	return json.Marshal(ja)
}

func (jb json_binary) get_json() ([]byte, error) {
	return json.Marshal(jb)
}

func (jc json_c_source) get_json() ([]byte, error) {
	return json.Marshal(jc)
}

func (jcl json_command_line) get_json() ([]byte, error) {
	return json.Marshal(jcl)
}

func (jfc json_file_close) get_json() ([]byte, error) {
	return json.Marshal(jfc)
}

func (jfo json_file_open) get_json() ([]byte, error) {
	return json.Marshal(jfo)
}

func (jfr json_file_read) get_json() ([]byte, error) {
	return json.Marshal(jfr)
}

func (jfw json_file_write) get_json() ([]byte, error) {
	return json.Marshal(jfw)
}

func (jf json_free) get_json() ([]byte, error) {
	return json.Marshal(jf)
}

func (jfi json_function_invocation) get_json() ([]byte, error) {
	return json.Marshal(jfi)
}

func (jl json_line_change) get_json() ([]byte, error) {
	return json.Marshal(jl)
}

func (jm json_malloc) get_json() ([]byte, error) {
	return json.Marshal(jm)
}

func (jo json_output) get_json() ([]byte, error) {
	return json.Marshal(jo)
}

func (jr json_register) get_json() ([]byte, error) {
	return json.Marshal(jr)
}

func (jra json_return_address) get_json() ([]byte, error) {
	return json.Marshal(jra)
}

func (jret json_return) get_json() ([]byte, error) {
	return json.Marshal(jret)
}

func (js json_section) get_json() ([]byte, error) {
	return json.Marshal(js)
}

func (jva json_variable_access) get_json() ([]byte, error) {
	return json.Marshal(jva)
}
