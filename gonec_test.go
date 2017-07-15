package gonec

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func Test_interpreter_ParseAndRun(t *testing.T) {
	ti1 := Interpreter()
	in1 := strings.NewReader(`
	// This is scanned code.
	Пакет Основной

		//перем дд,вв;
		
		Функция а(б,в,г) экспОрт
			если б<>в тогда
				д=б
				д=в
			иначе
				д=0
			конецЕсли
			возврат д
		КонецФункции

		б = а(1,2,3)
		Сообщить(б)
	`)
	w := &bytes.Buffer{}
	err := ti1.ParseAndRun(in1, w)
	fmt.Println(w.String())
	if err != nil {
		fmt.Println(err.Error())
	}
}
