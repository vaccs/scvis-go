package vaccs

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestReadVaccsFile(t *testing.T) {

	vfr := (&Vaccs_file_reader{}).AddPath("t1.vaccs")

	json_vaccs, err := vfr.GetVaccsAnalysisRecords()

	if err != nil {
		t.Error(err)
	}
	output, _ := json.Marshal(json_vaccs)

	fmt.Println("json_vaccs: ", string(output))
}
