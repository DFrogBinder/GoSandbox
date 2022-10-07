package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type root struct {
	XMLName xml.Name      `xml:"root"`
	Files   []ProtocolMap `xml:",any"`
}

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

	// for i := range WholeProtocol.ProtocolData {
	// 	var Test string = strings.Join(WholeProtocol.ProtocolData[i].SequenceName, ",")
	// 	fmt.Println(Test)
	// 	for j := range WholeProtocol.ProtocolData[i].SequenceParam {
	// 		fmt.Println(WholeProtocol.ProtocolData[i].SequenceName[0] + " - " + WholeProtocol.ProtocolData[i].SequenceParam[j] + " - " + WholeProtocol.ProtocolData[i].SequenceVal[j])
	// 	}
	// }
	Delimiter := "\""
	Test := strings.Join(WholeProtocol.ProtocolData[0].SequenceName, "")
	split := strings.Split(Test, Delimiter)
	fmt.Println(len(split))
	//fmt.Println(len(FirstReportMapping.SequenceParam))
}
