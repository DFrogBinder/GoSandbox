package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type SequenceMap struct {
	SequenceName  []string `xml:"PrintProtocol>Protocol>SubStep>ProtHeaderInfo>HeaderProtPath"`
	SequenceParam []string `xml:"PrintProtocol>Protocol>SubStep>Card>ProtParameter>Label"`
	SequenceVal   []string `xml:"PrintProtocol>Protocol>SubStep>Card>ProtParameter>ValueAndUnit"`
}

func main() {
	var GoldenStand, SuppliedRep string
	var FirstReportMapping SequenceMap
	var SecondReportMapping SequenceMap
	var help bool

	flag.StringVar(&GoldenStand, "gold", "", "Specify the path to the golden standart report file")
	flag.StringVar(&SuppliedRep, "supp", "", "Specify the path to the supplied report file")
	flag.BoolVar(&help, "help", false, "Display usage information")

	flag.Parse()

	// Open our xmlFile
	FirstXmlFile, err := os.Open(GoldenStand)
	SecondXmlFile, err := os.Open(SuppliedRep)

	if err != nil {
		fmt.Println(err)
	}

	// read our opened xmlFile as a byte array.
	FbyteValue, _ := ioutil.ReadAll(FirstXmlFile)
	SbyteValue, _ := ioutil.ReadAll(SecondXmlFile)

	xml.Unmarshal(FbyteValue, &FirstReportMapping)
	xml.Unmarshal(SbyteValue, &SecondReportMapping)

	FirstXmlFile.Close()
	SecondXmlFile.Close()

	for i := range FirstReportMapping.SequenceName {
		fmt.Println(FirstReportMapping.SequenceName[i])
		fmt.Println(FirstReportMapping.SequenceParam[i])
		fmt.Println(FirstReportMapping.SequenceVal[i])

	}

}
