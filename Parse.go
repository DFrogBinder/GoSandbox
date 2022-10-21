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

func Compare_Inner_Structure(FR []string, SR []string, SimilarityList []int, WrongParameters []string) ([]int, []string) {
	if FR[1] == SR[1] && FR[2] == SR[2] {
		if FR[3] == SR[3] {
			SimilarityList = append(SimilarityList, 1)
		} else {
			msg := fmt.Sprintln("In: ", SR[0], SR[1], "=> Found: ", SR[2], SR[3], " => Should be:", FR[3])
			WrongParameters = append(WrongParameters, msg)
		}
	}
	return SimilarityList, WrongParameters
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
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
				TempParameter = ProtocolName + "$" + Report.ProtocolData[i].Card[j].Name + "$" + Report.ProtocolData[i].Card[j].SequenceParam[k] + "$" + Report.ProtocolData[i].Card[j].SequenceVal[k]
			}
			FinalReportParamList = append(FinalReportParamList, TempParameter)
		}
	}
	return FinalReportParamList
}

func CompareReports(FR, SR []string) ([]int, []string) {
	var SimilarityList []int
	var WrongParameters []string

	var T2StarAlias = []string{"T2STAR DIXON", "LMS_T2Star", "LMS T2STAR DIXON", "LMS_T2STAR", "LMS T2STAR"}
	var ExcludeAlias = []string{"Pancreas", "pancreas", "kidney", "Kidney", "localizer", "VOLUME"}
	var IdealALias = []string{"IDEAL", "LMS_IDEAL", "LMS IDEAL", "LMS Ideal", "LMS IDEAL"}
	var MolliAlias = []string{"MOLLI", "LMS_MOLLI", "LMS MOLLI", "LMS Molli", "LMS_Molli"}
	var MostAlias = []string{"MOST", "LMS MOST", "LMS Most", "LMS_MOST", "LMS_Most"}

	for i := range FR {
		for j := range SR {
			tFR := strings.Split(FR[i], "$")
			tSR := strings.Split(SR[j], "$")

			// T2Star Area
			if contains(T2StarAlias, tFR[0]) &&
				contains(T2StarAlias, tSR[0]) &&
				!contains(ExcludeAlias, tFR[0]) &&
				!contains(ExcludeAlias, tSR[0]) {
				SimilarityList, WrongParameters = Compare_Inner_Structure(tFR, tSR, SimilarityList, WrongParameters)
				// Ideal Area
			} else if contains(IdealALias, tFR[0]) &&
				contains(IdealALias, tSR[0]) &&
				!contains(ExcludeAlias, tFR[0]) &&
				!contains(ExcludeAlias, tSR[0]) {
				SimilarityList, WrongParameters = Compare_Inner_Structure(tFR, tSR, SimilarityList, WrongParameters)
				// Molli Area
			} else if contains(MolliAlias, tFR[0]) &&
				contains(MolliAlias, tSR[0]) &&
				!contains(ExcludeAlias, tFR[0]) &&
				!contains(ExcludeAlias, tSR[0]) {
				SimilarityList, WrongParameters = Compare_Inner_Structure(tFR, tSR, SimilarityList, WrongParameters)
			} else if contains(MostAlias, tFR[0]) &&
				contains(MostAlias, tSR[0]) &&
				!contains(ExcludeAlias, tFR[0]) &&
				!contains(ExcludeAlias, tSR[0]) {
				SimilarityList, WrongParameters = Compare_Inner_Structure(tFR, tSR, SimilarityList, WrongParameters)
			}
		}
	}

	return SimilarityList, WrongParameters
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

	Similarity, Message := CompareReports(FirstReportParamsList, SecondReportParamsList)
	fmt.Println("Correct Paramters: ", len(Similarity))
	fmt.Println(Message)
}
