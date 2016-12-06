package main

import (
	"bytes"
	"encoding/xml"
	"log"
	"os"
	"time"
)

const s = `<styles>
	<item type="image" id="background" w="800" h="416" imageurl="@DNA/common/bg.png"/>
</styles>
`

func main() {
	f, err := os.Create(time.Now().Format("log/2006_0102_1504_05.log"))
	if err != nil {
		return
	}
	log.SetOutput(f)

	v := &Tag{}
	xml.NewDecoder(bytes.NewBuffer([]byte(s))).Decode(&v)

	log.Println(v.Name)
	log.Println(v.Attr)
	log.Println(v.Children)

}
