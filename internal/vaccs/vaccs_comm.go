package vaccs

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	vACK = iota
	vSTARTCFILES
	vCFILE
	vENDCFILES
	vSTARTPARAMS
	vPARAM
	vENDPARAMS
	vVACCSFILE
	vSTARTVACCSARGS
	vVACCSARGS
	vENDVACCSARGS
	vEXECUTABLE
	vSUCCESS
	vFAILURE_COMPILE
	vFAILURE_ANALYZE
	vFAILURE_PROTOCOL
	vFAILURE_TIMEOUT
	vFAILURE_SYSTEM
)

type VaccsComm struct {
	socket net.Conn
}

func NewVaccsComm() *VaccsComm {
	return &VaccsComm{}
}

func (vc *VaccsComm) Connect() error {
	var err error
	vc.socket, err = net.Dial("tcp", "localhost:3580")
	if err != nil {
		return err
	}

	vc.socket.SetReadDeadline(time.Time{})

	return nil
}

func (vc *VaccsComm) Close() {
	vc.socket.Close()
}

func (vc *VaccsComm) PerformAnalysis(filenames []string, cSources []string, vargs []string, executableName string, params []string) ([]byte, error) {
	err := vc.send(vc.makeSTARTCFILES(uint32(len(filenames))))

	if err != nil {
		log.Println("Error sending start cfiles")
		return nil, err
	}

	err = vc.receive_ack()
	if err != nil {
		log.Println("Error receiving ack")
		return nil, err
	}

	for i, filename := range filenames {
		err = vc.send(vc.makeCFILE(filename, cSources[i]))
		if err != nil {
			log.Println("Error sending cfile")
			return nil, err
		}
		err = vc.receive_ack()
		if err != nil {
			log.Println("Error receiving ack")
			return nil, err
		}
	}

	err = vc.send(vc.makeENDCFILES())
	if err != nil {
		log.Println("Error sending end cfiles")
		return nil, err
	}

	err = vc.receive_ack()
	if err != nil {
		log.Println("Error receiving ack")
		return nil, err
	}

	err = vc.send(vc.makeSTARTVACCSARG(uint32(len(vargs))))
	if err != nil {
		log.Println("Error sending start vaccsargs")
		return nil, err
	}

	err = vc.receive_ack()
	if err != nil {
		log.Println("Error receiving ack")
		return nil, err
	}

	for _, varg := range vargs {
		data := vc.makeVACCSARG(varg)
		err = vc.send(data)
		if err != nil {
			log.Println("Error sending vaccsarg")
			return nil, err
		}
		err = vc.receive_ack()
		if err != nil {
			log.Println("Error receiving ack")
			return nil, err
		}
	}

	err = vc.send(vc.makeENDVACCSARG())

	if err != nil {
		log.Println("Error sending end vaccsargs")
		return nil, err
	}
	err = vc.receive_ack()
	if err != nil {
		log.Println("Error receiving ack")
		return nil, err
	}

	err = vc.send(vc.makeEXECUTABLE(executableName))
	if err != nil {
		log.Println("Error sending executable name")
		return nil, err
	}

	err = vc.receive_ack()
	if err != nil {
		log.Println("Error receiving ack")
		return nil, err
	}

	err = vc.send(vc.makeSTARTPARAMS(uint32(len(params))))
	if err != nil {
		log.Println("Error sending start params")
		return nil, err
	}

	err = vc.receive_ack()
	if err != nil {
		log.Println("Error receiving ack")
		return nil, err
	}

	for _, param := range params {
		data := vc.makePARAM(param)
		err = vc.send(data)
		if err != nil {
			log.Println("Error sending param")
			return nil, err
		}
		err = vc.receive_ack()
		if err != nil {
			log.Println("Error receiving ack")
			return nil, err
		}
	}

	err = vc.send(vc.makeENDPARAMS())
	if err != nil {
		log.Println("Error sending end params")
		return nil, err
	}
	err = vc.receive_ack()
	if err != nil {
		log.Println("Error receiving ack")
		return nil, err
	}

	data, err := vc.receive_json()
	if err != nil {
		log.Println("Error receiving analysis data")
		return nil, err
	}
	err = vc.send_ack()
	if err != nil {
		log.Println("Error sending ack")
		return nil, err
	}

	return data, nil
}

func (vc *VaccsComm) send(data []byte) error {
	_, err := vc.socket.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func makeAnalysisError(tag byte, log_data string) error {
	if tag == vFAILURE_ANALYZE {
		return fmt.Errorf("The analyzer failed on this project: \n" + log_data)
	} else if tag == vFAILURE_COMPILE {
		return fmt.Errorf("The compiler failed on this project: \n" + log_data)
	} else if tag == vFAILURE_PROTOCOL {
		return fmt.Errorf("Communication with the analyzer failed on this project: \n" + log_data)
	} else if tag == vFAILURE_TIMEOUT {
		return fmt.Errorf("The analyzer timed out on this project: \n" + log_data)
	} else {
		return fmt.Errorf("A system error occured while analyzing this project: \n" + log_data)
	}
}

func (vc *VaccsComm) read_bytes(n int) ([]byte, error) {
	data := make([]byte, n)

	num_bytes, err := vc.socket.Read(data)

	if err != nil {
		return nil, err
	}
	if data[0] > vSUCCESS {
		return nil, makeAnalysisError(data[0], string(data[1:]))
	} else if num_bytes != n {
		return nil, fmt.Errorf("expected %d bytes, got %d", n, num_bytes)
	}

	return data, nil
}

func (vc *VaccsComm) send_ack() error {
	_, err := vc.socket.Write(vc.makeACK())
	if err != nil {
		return err
	}
	return nil
}

func (vc *VaccsComm) receive_ack() error {
	header, err := vc.read_bytes(1)
	if err != nil {
		return err
	}
	if header[0] != vACK {
		return fmt.Errorf("protocol error: expected vACK, got %d", header[0])
	}
	return nil
}

func (vc *VaccsComm) receive_json() ([]byte, error) {
	header, err := vc.read_bytes(5)
	if err != nil {
		log.Println("Error reading header")
		return nil, err
	}
	if header[0] != vVACCSFILE {
		return nil, fmt.Errorf("protocol error: expected vVACCSFILE, got %d", header[0])
	}

	sizeUint32 := uint32(header[1])<<24 | uint32(header[2])<<16 | uint32(header[3])<<8 | uint32(header[4])
	data, err := vc.read_bytes(int(sizeUint32))
	if err != nil {
		log.Println("Error reading analysis data")
		return nil, err
	}
	err = vc.send(vc.makeACK())
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (vc *VaccsComm) makeACK() []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(vACK))
	return buf.Bytes()
}

func writeInt(buf *bytes.Buffer, val int) []byte {
	buf.WriteByte(byte(val >> 24))
	buf.WriteByte(byte(val >> 16))
	buf.WriteByte(byte(val >> 8))
	buf.WriteByte(byte(val))
	return buf.Bytes()
}

func writeUint(buf *bytes.Buffer, val uint32) {
	buf.WriteByte(byte(val >> 24))
	buf.WriteByte(byte(val >> 16))
	buf.WriteByte(byte(val >> 8))
	buf.WriteByte(byte(val))
}

func (vc *VaccsComm) makeSTARTCFILES(numfiles uint32) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(vSTARTCFILES))
	writeUint(buf, numfiles)
	return buf.Bytes()
}

func (vc *VaccsComm) makeCFILE(filename string, source string) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(vCFILE))
	writeInt(buf, len(filename))
	buf.WriteString(filename)

	writeInt(buf, len(source))
	buf.WriteString(source)

	return buf.Bytes()
}

func (vc *VaccsComm) makeENDCFILES() []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(vENDCFILES))
	return buf.Bytes()
}

func (vc *VaccsComm) makeSTARTPARAMS(numparams uint32) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(vSTARTPARAMS))
	buf.WriteByte(byte(numparams >> 24))
	buf.WriteByte(byte(numparams >> 16))
	buf.WriteByte(byte(numparams >> 8))
	buf.WriteByte(byte(numparams))
	return buf.Bytes()
}

func (vc *VaccsComm) makePARAM(param string) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(vPARAM))
	buf.WriteByte(byte(len(param) >> 24))
	buf.WriteByte(byte(len(param) >> 16))
	buf.WriteByte(byte(len(param) >> 8))
	buf.WriteByte(byte(len(param)))
	buf.WriteString(param)
	return buf.Bytes()
}

func (vc *VaccsComm) makeENDPARAMS() []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(vENDPARAMS))
	return buf.Bytes()
}

func (vc *VaccsComm) makeSTARTVACCSARG(numargs uint32) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(vSTARTVACCSARGS))
	writeUint(buf, numargs)
	return buf.Bytes()
}

func (vc *VaccsComm) makeVACCSARG(arg string) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(vVACCSARGS))
	writeInt(buf, len(arg))
	buf.WriteString(arg)
	return buf.Bytes()
}

func (vc *VaccsComm) makeENDVACCSARG() []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(vENDVACCSARGS))
	return buf.Bytes()
}

func (vc *VaccsComm) makeEXECUTABLE(executable string) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(vEXECUTABLE))
	writeInt(buf, len(executable))
	buf.WriteString(executable)
	return buf.Bytes()
}
