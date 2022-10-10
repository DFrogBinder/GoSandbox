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

func Extract_Parameters(Report ProtocolMap) []string {
	var FinalReportParamList []string

	for i := range Report.ProtocolData {
		var Test string = strings.Join(Report.ProtocolData[i].SequenceName, "")
		split := strings.Split(Test, `\`)
		ProtocolName := split[len(split)-1]
		for j := range Report.ProtocolData[i].SequenceParam {
			TempParameter := ProtocolName + " - " + Report.ProtocolData[i].SequenceParam[j] + " - " + Report.ProtocolData[i].SequenceVal[j]
			FinalReportParamList = append(FinalReportParamList, TempParameter)
		}
	}
	return FinalReportParamList
}

func main() {
	var GoldenStand, SuppliedRep string
	var FirstReportMapping ProtocolMap
	var SecondReportMapping ProtocolMap
	var SecondReportParamsList, FirstReportParamsList []string
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

	// Creates Parameter List for the first report
	FirstReportParamsList = Extract_Parameters(FirstReportMapping)

	// Creates Parameter list for the Second Report
	SecondReportParamsList = Extract_Parameters(SecondReportMapping)

	fmt.Println("Number of parameters in the first report: " + fmt.Sprint(len(FirstReportParamsList)))
	fmt.Println("Number of parameters in the first report: " + fmt.Sprint(len(SecondReportParamsList)))
}
