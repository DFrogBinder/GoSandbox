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

func contains(s []string, str string) bool {
	for _, v := range s {
		if strings.Compare(strings.ReplaceAll(v, " ", ""), strings.ReplaceAll(str, " ", "")) == 0 {
			return true
		}
	}
	return false
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

func CompareReports(FR, SR []string) []string {
	var SimilarityList []string

	var T2StarAlias = []string{"T2Star", "T2_Star", "T2*", "T2STAR", "T2_STAR", "LMS_T2Star", "LMS T2Star", "LMS_T2STAR", "LMS T2STAR"}
	var ExcludeAlias = []string{"Pancreas", "pancreas", "kidney", "Kidney", "localizer"}
	var IdealALias = []string{"IDEAL", "Ideal", "LMS_IDEAL", "LMS IDEAL", "LMS Ideal", "LMS IDEAL"}
	var MolliAlias = []string{"MOLLI", "Molli", "LMS_MOLLI", "LMS MOLLI", "LMS Molli", "LMS_Molli"}

	for i := range FR {
		for j := range SR {
			tFR := strings.Split(FR[i], "-")
			tSR := strings.Split(SR[j], "-")
			// T2Star Area
			if contains(T2StarAlias, tFR[0]) &&
				contains(T2StarAlias, tSR[0]) &&
				!contains(ExcludeAlias, tFR[0]) &&
				!contains(ExcludeAlias, tSR[0]) {
				fmt.Println("Test Passed in T2Star Area")
				// Ideal Area
			} else if contains(IdealALias, tFR[0]) &&
				contains(IdealALias, tSR[0]) &&
				!contains(ExcludeAlias, tFR[0]) &&
				!contains(ExcludeAlias, tSR[0]) {
				fmt.Println("Test Passed in Ideal Area")
				// Molli Area
			} else if contains(MolliAlias, tFR[0]) &&
				contains(MolliAlias, tSR[0]) &&
				!contains(ExcludeAlias, tFR[0]) &&
				!contains(ExcludeAlias, tSR[0]) {
				fmt.Println("Test Passed in Molli Area")
			}
		}
	}

	return SimilarityList
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

	Similarity := CompareReports(FirstReportParamsList, SecondReportParamsList)
	fmt.Print(Similarity)
}
