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
	ProtocolData []Protocol `xml:"PrintProtocol>Protocol"`
}

type Card struct {
	ID            string   `xml:"ID,attr"`
	SequenceParam []string `xml:"ProtParameter>Label"`
	SequenceVal   []string `xml:"ProtParameter>ValueAndUnit"`
}

type Protocol struct {
	ID string `xml:"Id,attr"`
	//Title        string `xml:"HeaderTitle"`
	SequenceName string `xml:"SubStep>ProtHeaderInfo>HeaderProtPath"`
	SequeceCard  []Card `xml:"SubStep>Card"`
}

func stcmp(s []string, str string) bool {
	for _, v := range s {
		if strings.Compare(strings.ReplaceAll(v, " ", ""), strings.ReplaceAll(str, " ", "")) == 0 {
			return true
		}
	}
	return false
}

func Extract_Parameters(Report ProtocolMap) []string {
	var FinalReportParamList []string

	// for i := range Report.ProtocolData {
	// 	var Test string = strings.Join(Report.ProtocolData[i].SequenceName, "")
	// 	split := strings.Split(Test, `\`)
	// 	ProtocolName := split[len(split)-1]
	// 	for j := range Report.ProtocolData[i].SequenceParam {
	// 		TempParameter := ProtocolName + " - " + Report.ProtocolData[i].SequenceParam[j] + " - " + Report.ProtocolData[i].SequenceVal[j]
	// 		FinalReportParamList = append(FinalReportParamList, TempParameter)
	// 	}
	// }
	return FinalReportParamList
}

func CompareReports(FR, SR []string) []string {
	var SimilarityList []string

	var T2StarAlias = []string{"LMS_T2Star", "LMS T2STAR DIXON", "LMS_T2STAR", "LMS T2STAR"}
	var ExcludeAlias = []string{"Pancreas", "pancreas", "kidney", "Kidney", "localizer_haste_bh"}
	var IdealALias = []string{"LMS_IDEAL", "LMS IDEAL", "LMS Ideal", "LMS IDEAL"}
	var MolliAlias = []string{"LMS_MOLLI", "LMS MOLLI", "LMS Molli", "LMS_Molli"}
	var MostAlias = []string{"LMS MOST", "LMS Most", "LMS_MOST", "LMS_Most"}

	for i := range FR {
		for j := range SR {
			tFR := strings.Split(FR[i], "-")
			tSR := strings.Split(SR[j], "-")
			// T2Star Area
			if stcmp(T2StarAlias, tFR[0]) &&
				stcmp(T2StarAlias, tSR[0]) &&
				!stcmp(ExcludeAlias, tFR[0]) &&
				!stcmp(ExcludeAlias, tSR[0]) {
				fmt.Println("Test Passed in T2Star Area")
				// Ideal Area
			} else if stcmp(IdealALias, tFR[0]) &&
				stcmp(IdealALias, tSR[0]) &&
				!stcmp(ExcludeAlias, tFR[0]) &&
				!stcmp(ExcludeAlias, tSR[0]) {
				fmt.Println("Test Passed in Ideal Area")
				// Molli Area
			} else if stcmp(MolliAlias, tFR[0]) &&
				stcmp(MolliAlias, tSR[0]) &&
				!stcmp(ExcludeAlias, tFR[0]) &&
				!stcmp(ExcludeAlias, tSR[0]) {
				fmt.Println("Test Passed in Molli Area")
			} else if stcmp(MostAlias, tFR[0]) &&
				stcmp(MostAlias, tSR[0]) &&
				!stcmp(ExcludeAlias, tFR[0]) &&
				!stcmp(ExcludeAlias, tSR[0]) {
				fmt.Println("Test Passed in MOST Area")
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

	fmt.Println(xml.Unmarshal(FbyteValue, &FirstReportMapping))
	xml.Unmarshal(SbyteValue, &SecondReportMapping)

	FirstXmlFile.Close()
	SecondXmlFile.Close()

	// Creates Parameter List for the first report
	FirstReportParamsList = Extract_Parameters(FirstReportMapping)

	// Creates Parameter list for the Second Report
	SecondReportParamsList = Extract_Parameters(SecondReportMapping)

	fmt.Println("Number of parameters in the first report: " + fmt.Sprint(len(FirstReportParamsList)))
	fmt.Println("Number of parameters in the first report: " + fmt.Sprint(len(SecondReportParamsList)))

	// Similarity := CompareReports(FirstReportParamsList, SecondReportParamsList)
	// fmt.Print(Similarity)
}
