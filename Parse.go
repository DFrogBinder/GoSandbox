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

type Protocol struct {
	XMLName      xml.Name `xml:"Protocol"`
	ID           string   `xml:"Id,attr"`
	SequenceName []string `xml:"SubStep>ProtHeaderInfo>HeaderProtPath"`
	Card         []Card   `xml:"SubStep>Card"`
}

type Card struct {
	Name          string   `xml:"name,attr"`
	SequenceParam []string `xml:"ProtParameter>Label"`
	SequenceVal   []string `xml:"ProtParameter>ValueAndUnit"`
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
	var TempParameter string

	for i := range Report.ProtocolData {
		var Test string = strings.Join(Report.ProtocolData[i].SequenceName, "")
		split := strings.Split(Test, `\`)
		ProtocolName := split[len(split)-1]
		for j := range Report.ProtocolData[i].Card {
			for k := range Report.ProtocolData[i].Card[j].SequenceParam {
				TempParameter = ProtocolName + "_" + Report.ProtocolData[i].Card[j].Name + "_" + Report.ProtocolData[i].Card[j].SequenceParam[k] + "_" + Report.ProtocolData[i].Card[j].SequenceVal[k]
			}
			FinalReportParamList = append(FinalReportParamList, TempParameter)
		}
	}
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
			} else if contains(MostAlias, tFR[0]) &&
				contains(MostAlias, tSR[0]) &&
				!contains(ExcludeAlias, tFR[0]) &&
				!contains(ExcludeAlias, tSR[0]) {
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

	if FbyteValue != nil && SbyteValue != nil {
		xml.Unmarshal(FbyteValue, &FirstReportMapping)
		xml.Unmarshal(SbyteValue, &SecondReportMapping)
	} else {
		fmt.Println("Problem is reading the input files - value is <nil>. Terminating...")
		return
	}

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
