package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type ProtocolMap struct {
	ProtocolData []SequenceMap `xml:"PrintProtocol>Protocol"`
}

type SequenceMap struct {
	SequenceName  []string `xml:"SubStep>ProtHeaderInfo>HeaderProtPath"`
	SequenceParam []string `xml:"SubStep>Card>ProtParameter>Label"`
	SequenceVal   []string `xml:"SubStep>Card>ProtParameter>ValueAndUnit"`
}

func main() {
	var GoldenStand, SuppliedRep string
	var FirstReportMapping SequenceMap
	var SecondReportMapping SequenceMap
	var WholeProtocol ProtocolMap
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
	xml.Unmarshal(FbyteValue, &WholeProtocol)

	FirstXmlFile.Close()
	SecondXmlFile.Close()

	for i := range WholeProtocol.ProtocolData {
		var Test string = strings.Join(WholeProtocol.ProtocolData[i].SequenceName, "")
		split := strings.Split(Test, `\`)
		ProtocolName := split[len(split)-1]
		for j := range WholeProtocol.ProtocolData[i].SequenceParam {
			fmt.Println(ProtocolName + " - " + WholeProtocol.ProtocolData[i].SequenceParam[j] + " - " + WholeProtocol.ProtocolData[i].SequenceVal[j])
		}
	}
}
