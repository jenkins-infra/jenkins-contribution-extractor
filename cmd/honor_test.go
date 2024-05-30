/*
Copyright © 2024 Jean-Marc Meessen jean-marc@meessen-web.org

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_performHonorContributorSelection_params(t *testing.T) {
	type args struct {
		dataDir           string
		outputFileName    string
		monthToSelectFrom string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"inexistent data directory",
			args{
				dataDir:           "inexistentDir",
				monthToSelectFrom: "2024-04",
			},
			true,
		},
		{
			"valid data directory and month",
			args{
				dataDir:           "../test-data",
				monthToSelectFrom: "2024-04",
			},
			false,
		},
		{
			"invalid month",
			args{
				monthToSelectFrom: "junkMonth",
				dataDir:           "../test-data",
			},
			true,
		},
		{
			"invalid header in input file",
			args{
				dataDir:           "../test-data",
				monthToSelectFrom: "2024-03",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := performHonorContributorSelection(tt.args.dataDir, tt.args.outputFileName, tt.args.monthToSelectFrom); (err != nil) != tt.wantErr {
				t.Errorf("performHonorContributorSelection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_honorCommand_paramCheck_noMonth(t *testing.T) {
	//Setup environment
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	var commandArguments []string
	commandArguments = append(commandArguments, "honor", "--data_dir=../test-data")
	rootCmd.SetArgs(commandArguments)

	// execute command
	error := rootCmd.Execute()

	// check results
	assert.ErrorContains(t, error, "\"month\" argument is missing.", "Call should have failed with expected error.")
}

func Test_honorCommand_paramCheck_invalidMonth(t *testing.T) {
	//Setup environment
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	var commandArguments []string
	commandArguments = append(commandArguments, "honor", "junkMonth", "--data_dir=../test-data")
	rootCmd.SetArgs(commandArguments)

	// execute command
	error := rootCmd.Execute()

	// check results
	assert.ErrorContains(t, error, "\"junkMonth\" is not a valid month.", "Call should have failed with expected error.")
}

func Test_honorCommand_integrationTest_verbose(t *testing.T) {

	// Setup test environment
	tempDir := t.TempDir()
	// duplicate the file but keep the original filename
	dataFilename, err := duplicateFile("../test-data/pr_per_submitter-2024-04.csv", tempDir, false)

	assert.NoError(t, err, "Unexpected data file duplication error")
	assert.NotEmpty(t, dataFilename, "Failure to copy data file")

	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	var commandArguments []string
	commandArguments = append(commandArguments, "honor", "2024-04", "--data_dir="+tempDir, "--verbose")
	rootCmd.SetArgs(commandArguments)

	// execute command
	error := rootCmd.Execute()

	// check results
	assert.NoError(t, error, "Call should not have failed")
	assert.NotEmpty(t, filepath.Join(tempDir, "honored_contributor.csv"), "Failure to generate target file")
	//TODO: check that it has the correct header
	//TODO: check that the data (second line) has usable data (is this worth it?)

}

func Test_stringifySlice(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"happy case",
			args{s: []string{"aaa", "bbb", "ccc"}},
			"aaa bbb ccc",
		},
		{
			"Single item case",
			args{s: []string{"aaa"}},
			"aaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringifySlice(tt.args.s); got != tt.want {
				t.Errorf("stringifySlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateHonoredContributorDataAsCSV(t *testing.T) {
	type args struct {
		contributorData HonoredContributorData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"typical case",
			args{
				contributorData: HonoredContributorData{
					handle:            "GH_handle",
					fullName:          "author_fullName",
					authorURL:         "author_url",
					authorAvatarUrl:   "author_avatar",
					authorCompany:     "a_company",
					month:             "a_month",
					totalPRs_found:    "PR_found",
					totalPRs_expected: "PR_expected",
					repositories:      "repositories",
				},
			},
			"\"a_month\", \"GH_handle\", \"author_fullName\", \"a_company\", \"author_url\", \"author_avatar\", \"PR_found\", \"repositories\"",
		},
		{
			"with empty fields",
			args{
				contributorData: HonoredContributorData{
					handle:            "GH_handle",
					fullName:          "",
					authorURL:         "author_url",
					authorAvatarUrl:   "author_avatar",
					authorCompany:     "",
					month:             "a_month",
					totalPRs_found:    "PR_found",
					totalPRs_expected: "PR_expected",
					repositories:      "repositories",
				},
			},
			"\"a_month\", \"GH_handle\", \"\", \"\", \"author_url\", \"author_avatar\", \"PR_found\", \"repositories\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateHonoredContributorDataAsCSV(tt.args.contributorData); got != tt.want {
				t.Errorf("generateHonoredContributorDataAsCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

//Test the whole input file for query mismatches
// func Test_getSubmitterPRsForBasil(t *testing.T) {
// 	check_getSubmittersPRfromGH(t, "basil", "69", "2024-04")
// }

// func Test_getSubmitterPRsForgounthar(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"gounthar","40","2024-04")
//  }
//  func Test_getSubmitterPRsForlemeurherve(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"lemeurherve","33","2024-04")
//  }
//  func Test_getSubmitterPRsForsmerle33(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"smerle33","31","2024-04")
//  }
//  func Test_getSubmitterPRsForMarkEWaite(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"MarkEWaite","31","2024-04")
//  }
//  func Test_getSubmitterPRsFordduportal(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"dduportal","28","2024-04")
//  }
//  func Test_getSubmitterPRsForjanfaracik(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"janfaracik","25","2024-04")
//  }
//  func Test_getSubmitterPRsForjonesbusy(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"jonesbusy","20","2024-04")
//  }
//  func Test_getSubmitterPRsFordanielBeck(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"daniel-beck","19","2024-04")
//  }
//  func Test_getSubmitterPRsFortimja(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"timja","16","2024-04")
//  }
//  func Test_getSubmitterPRsFormawinter69(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"mawinter69","15","2024-04")
//  }
//  func Test_getSubmitterPRsForjglick(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"jglick","15","2024-04")
//  }
//  func Test_getSubmitterPRsForNotMyFault(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"NotMyFault","14","2024-04")
//  }
//  func Test_getSubmitterPRsFormichaelDoubez(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"michael-doubez","12","2024-04")
//  }
//  func Test_getSubmitterPRsForkmartens27(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"kmartens27","12","2024-04")
//  }
//  func Test_getSubmitterPRsForuhafner(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"uhafner","11","2024-04")
//  }
//  func Test_getSubmitterPRsForkrisstern(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"krisstern","11","2024-04")
//  }
//  func Test_getSubmitterPRsForzbynek(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"zbynek","9","2024-04")
//  }
//  func Test_getSubmitterPRsFornikitaTkachenkoDatadog(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"nikita-tkachenko-datadog","9","2024-04")
//  }
//  func Test_getSubmitterPRsForalecharp(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"alecharp","9","2024-04")
//  }
//  func Test_getSubmitterPRsForsusmitagorai29(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"susmitagorai29","8","2024-04")
//  }
//  func Test_getSubmitterPRsForjtnord(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"jtnord","8","2024-04")
//  }
//  func Test_getSubmitterPRsForolamy(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"olamy","7","2024-04")
//  }
//  func Test_getSubmitterPRsForalextu(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"alextu","7","2024-04")
//  }
//  func Test_getSubmitterPRsForStefanSpieker(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"StefanSpieker","7","2024-04")
//  }
//  func Test_getSubmitterPRsFortamarleviCm(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"tamarleviCm","6","2024-04")
//  }
//  func Test_getSubmitterPRsForstrangelookingnerd(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"strangelookingnerd","6","2024-04")
//  }
//  func Test_getSubmitterPRsForjanasrikanth(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"janasrikanth","6","2024-04")
//  }
//  func Test_getSubmitterPRsFordwnusbaum(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"dwnusbaum","6","2024-04")
//  }
//  func Test_getSubmitterPRsFordamianszczepanik(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"damianszczepanik","6","2024-04")
//  }
//  func Test_getSubmitterPRsForAniketNS(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"AniketNS","6","2024-04")
//  }
//  func Test_getSubmitterPRsForwaltwilo(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"waltwilo","5","2024-04")
//  }
//  func Test_getSubmitterPRsFormPokornyETM(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"mPokornyETM","5","2024-04")
//  }
//  func Test_getSubmitterPRsForhashar(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"hashar","5","2024-04")
//  }
//  func Test_getSubmitterPRsForBobDu(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"BobDu","5","2024-04")
//  }
//  func Test_getSubmitterPRsForysmaoui(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"ysmaoui","4","2024-04")
//  }
//  func Test_getSubmitterPRsForvishalhcl5960(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"vishalhcl-5960","4","2024-04")
//  }
//  func Test_getSubmitterPRsForrahulkaukuntla(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"rahulkaukuntla","4","2024-04")
//  }
//  func Test_getSubmitterPRsFormikecirioli(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"mikecirioli","4","2024-04")
//  }
//  func Test_getSubmitterPRsFormaksudurRahmanMaruf(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"maksudur-rahman-maruf","4","2024-04")
//  }
//  func Test_getSubmitterPRsForarturmelanchyk(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"arturmelanchyk","4","2024-04")
//  }
//  func Test_getSubmitterPRsForthomasvincent(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"thomasvincent","3","2024-04")
//  }
//  func Test_getSubmitterPRsForsridamul(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"sridamul","3","2024-04")
//  }
//  func Test_getSubmitterPRsForrsandell(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"rsandell","3","2024-04")
//  }
//  func Test_getSubmitterPRsForrkosegi(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"rkosegi","3","2024-04")
//  }
//  func Test_getSubmitterPRsForpboLinaro(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"pbo-linaro","3","2024-04")
//  }
//  func Test_getSubmitterPRsForjulieheard(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"julieheard","3","2024-04")
//  }
//  func Test_getSubmitterPRsForfabiodcasilva(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"fabiodcasilva","3","2024-04")
//  }
//  func Test_getSubmitterPRsForbzzitsme(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"bzzitsme","3","2024-04")
//  }
//  func Test_getSubmitterPRsForandreibangau99(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"andreibangau99","3","2024-04")
//  }
//  func Test_getSubmitterPRsForWaschndolos(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Waschndolos","3","2024-04")
//  }
//  func Test_getSubmitterPRsForRomainGeissler1A(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Romain-Geissler-1A","3","2024-04")
//  }
//  func Test_getSubmitterPRsForslide(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"slide","2","2024-04")
//  }
//  func Test_getSubmitterPRsForrepolevedavaj(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"repolevedavaj","2","2024-04")
//  }
//  func Test_getSubmitterPRsForpfeuffer(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"pfeuffer","2","2024-04")
//  }
//  func Test_getSubmitterPRsForhakre(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"hakre","2","2024-04")
//  }
//  func Test_getSubmitterPRsForhaidao247(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"haidao247","2","2024-04")
//  }
//  func Test_getSubmitterPRsForgvazquezmorean(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"gvazquezmorean","2","2024-04")
//  }
//  func Test_getSubmitterPRsForgabrielCheck24(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"gabriel-check24","2","2024-04")
//  }
//  func Test_getSubmitterPRsForfroque(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"froque","2","2024-04")
//  }
//  func Test_getSubmitterPRsForclayburn(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"clayburn","2","2024-04")
//  }
//  func Test_getSubmitterPRsForcdgopal(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"cdgopal","2","2024-04")
//  }
//  func Test_getSubmitterPRsForcarRoll(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"car-roll","2","2024-04")
//  }
//  func Test_getSubmitterPRsForawangParasoft(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"awang-parasoft","2","2024-04")
//  }
//  func Test_getSubmitterPRsForamuniz(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"amuniz","2","2024-04")
//  }
//  func Test_getSubmitterPRsForal3xanndru(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"al3xanndru","2","2024-04")
//  }
//  func Test_getSubmitterPRsForSOOSMMalony(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"SOOS-MMalony","2","2024-04")
//  }
//  func Test_getSubmitterPRsForPierreBtz(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"PierreBtz","2","2024-04")
//  }
//  func Test_getSubmitterPRsForMunishh992(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Munishh992","2","2024-04")
//  }
//  func Test_getSubmitterPRsForMarkRx(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"MarkRx","2","2024-04")
//  }
//  func Test_getSubmitterPRsForMarioFuchsTT(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"MarioFuchsTT","2","2024-04")
//  }
//  func Test_getSubmitterPRsForLmhJava(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Lmh-java","2","2024-04")
//  }
//  func Test_getSubmitterPRsForKiryushinAndrey(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Kiryushin-Andrey","2","2024-04")
//  }
//  func Test_getSubmitterPRsForFedeLo13(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"FedeLo13","2","2024-04")
//  }
//  func Test_getSubmitterPRsForDominikRusso(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"DominikRusso","2","2024-04")
//  }
//  func Test_getSubmitterPRsForDohbedoh(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Dohbedoh","2","2024-04")
//  }
//  func Test_getSubmitterPRsForAnski1(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Anski1","2","2024-04")
//  }
//  func Test_getSubmitterPRsForyyuyanyu(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"yyuyanyu","1","2024-04")
//  }
//  func Test_getSubmitterPRsForwollow12(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"wollow12","1","2024-04")
//  }
//  func Test_getSubmitterPRsForwelandaz(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"welandaz","1","2024-04")
//  }
//  func Test_getSubmitterPRsForvigneshtestsigma(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"vigneshtestsigma","1","2024-04")
//  }
//  func Test_getSubmitterPRsForvemulaanvesh(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"vemulaanvesh","1","2024-04")
//  }
//  func Test_getSubmitterPRsForvahidsh1(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"vahidsh1","1","2024-04")
//  }
//  func Test_getSubmitterPRsFortomasbjerre(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"tomasbjerre","1","2024-04")
//  }
//  func Test_getSubmitterPRsFortimbrown5(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"timbrown5","1","2024-04")
//  }
//  func Test_getSubmitterPRsFortilalx(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"tilalx","1","2024-04")
//  }
//  func Test_getSubmitterPRsForthyldrm(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"thyldrm","1","2024-04")
//  }
//  func Test_getSubmitterPRsFortherealsujitk(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"therealsujitk","1","2024-04")
//  }
//  func Test_getSubmitterPRsForswatipersistent(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"swatipersistent","1","2024-04")
//  }
//  func Test_getSubmitterPRsForstuartrowe(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"stuartrowe","1","2024-04")
//  }
//  func Test_getSubmitterPRsForskillcoder(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"skillcoder","1","2024-04")
//  }
//  func Test_getSubmitterPRsForsephirothj(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"sephiroth-j","1","2024-04")
//  }
//  func Test_getSubmitterPRsForrmartineias(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"rmartine-ias","1","2024-04")
//  }
//  func Test_getSubmitterPRsForreinhapa(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"reinhapa","1","2024-04")
//  }
//  func Test_getSubmitterPRsForraulArabaolaza(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"raul-arabaolaza","1","2024-04")
//  }
//  func Test_getSubmitterPRsForrahmanTiobe(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"rahman-tiobe","1","2024-04")
//  }
//  func Test_getSubmitterPRsForradhatiwari01(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"radhatiwari01","1","2024-04")
//  }
//  func Test_getSubmitterPRsForpyieh(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"pyieh","1","2024-04")
//  }
//  func Test_getSubmitterPRsForpurushotham99(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"purushotham99","1","2024-04")
//  }
//  func Test_getSubmitterPRsForpreyankababu(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"preyankababu","1","2024-04")
//  }
//  func Test_getSubmitterPRsForppettina(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"ppettina","1","2024-04")
//  }
//  func Test_getSubmitterPRsForpaulsavoie(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"paulsavoie","1","2024-04")
//  }
//  func Test_getSubmitterPRsForpatrikcerbak(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"patrikcerbak","1","2024-04")
//  }
//  func Test_getSubmitterPRsForowenmartinToast(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"owenmartin-toast","1","2024-04")
//  }
//  func Test_getSubmitterPRsForoospinar(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"oospinar","1","2024-04")
//  }
//  func Test_getSubmitterPRsForoffa(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"offa","1","2024-04")
//  }
//  func Test_getSubmitterPRsFornsBliu(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"ns-bliu","1","2024-04")
//  }
//  func Test_getSubmitterPRsFornmcc1212(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"nmcc1212","1","2024-04")
//  }
//  func Test_getSubmitterPRsFornitin6542(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"nitin-6542","1","2024-04")
//  }
//  func Test_getSubmitterPRsFornikhildabhade(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"nikhil-dabhade","1","2024-04")
//  }
//  func Test_getSubmitterPRsFornghiadhd2702(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"nghiadhd-2702","1","2024-04")
//  }
//  func Test_getSubmitterPRsFornattofriends(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"nattofriends","1","2024-04")
//  }
//  func Test_getSubmitterPRsFormjeanson(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"mjeanson","1","2024-04")
//  }
//  func Test_getSubmitterPRsFormguillem(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"mguillem","1","2024-04")
//  }
//  func Test_getSubmitterPRsFormayukothule(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"mayukothule","1","2024-04")
//  }
//  func Test_getSubmitterPRsFormattheimer(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"mattheimer","1","2024-04")
//  }
//  func Test_getSubmitterPRsForlpb1(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"lpb1","1","2024-04")
//  }
//  func Test_getSubmitterPRsForljackiewicz(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"ljackiewicz","1","2024-04")
//  }
//  func Test_getSubmitterPRsForlaudrup(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"laudrup","1","2024-04")
//  }
//  func Test_getSubmitterPRsForlangyizhao(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"langyizhao","1","2024-04")
//  }


 func Test_getSubmitterPRsForkyleleonhard(t *testing.T) {
	check_getSubmittersPRfromGH(t,"kyle-leonhard","1","2024-04")
 }


//  func Test_getSubmitterPRsForkvanzuijlen(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"kvanzuijlen","1","2024-04")
//  }
//  func Test_getSubmitterPRsForkothulemayur(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"kothulemayur","1","2024-04")
//  }
//  func Test_getSubmitterPRsForkaushalgupta88(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"kaushalgupta88","1","2024-04")
//  }
//  func Test_getSubmitterPRsForjwojnarowicz(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"jwojnarowicz","1","2024-04")
//  }
//  func Test_getSubmitterPRsForjudovana(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"judovana","1","2024-04")
//  }
//  func Test_getSubmitterPRsForjuanmafabbri(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"juanmafabbri","1","2024-04")
//  }
//  func Test_getSubmitterPRsForjondaley(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"jondaley","1","2024-04")
//  }
//  func Test_getSubmitterPRsForjayvirtanen(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"jayvirtanen","1","2024-04")
//  }
//  func Test_getSubmitterPRsForjandroav(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"jandroav","1","2024-04")
//  }
//  func Test_getSubmitterPRsForjamiejackson(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"jamiejackson","1","2024-04")
//  }
//  func Test_getSubmitterPRsForjahid1209(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"jahid1209","1","2024-04")
//  }
//  func Test_getSubmitterPRsForjluong(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"j-luong","1","2024-04")
//  }
//  func Test_getSubmitterPRsForimonteroperez(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"imonteroperez","1","2024-04")
//  }
//  func Test_getSubmitterPRsForharshanabandara(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"harshanabandara","1","2024-04")
//  }
//  func Test_getSubmitterPRsForgrowfrow(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"growfrow","1","2024-04")
//  }
//  func Test_getSubmitterPRsForgabrieleara(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"gabrieleara","1","2024-04")
//  }
//  func Test_getSubmitterPRsForfranknarf8(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"franknarf8","1","2024-04")
//  }
//  func Test_getSubmitterPRsForfrancisf(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"francisf","1","2024-04")
//  }
//  func Test_getSubmitterPRsForeduardtita(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"eduard-tita","1","2024-04")
//  }
//  func Test_getSubmitterPRsForduyluonganh(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"duyluonganh","1","2024-04")
//  }
//  func Test_getSubmitterPRsFordorin7bogdan(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"dorin7bogdan","1","2024-04")
//  }
//  func Test_getSubmitterPRsFordelineasagar(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"delinea-sagar","1","2024-04")
//  }
//  func Test_getSubmitterPRsFordarpanLalwani(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"darpanLalwani","1","2024-04")
//  }
//  func Test_getSubmitterPRsForcperrin88(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"cperrin88","1","2024-04")
//  }
//  func Test_getSubmitterPRsForcodervijay143(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"codervijay143","1","2024-04")
//  }
//  func Test_getSubmitterPRsForckullabosch(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"ckullabosch","1","2024-04")
//  }
//  func Test_getSubmitterPRsForckpattar(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"ckpattar","1","2024-04")
//  }
//  func Test_getSubmitterPRsForcconnert(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"cconnert","1","2024-04")
//  }
//  func Test_getSubmitterPRsForc0d3m0nky(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"c0d3m0nky","1","2024-04")
//  }
//  func Test_getSubmitterPRsForavivbs96(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"avivbs96","1","2024-04")
//  }
//  func Test_getSubmitterPRsForasimell(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"asimell","1","2024-04")
//  }
//  func Test_getSubmitterPRsForanniechellah(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"anniechellah","1","2024-04")
//  }
//  func Test_getSubmitterPRsForankitpatilhubs(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"ankit-patil-hubs","1","2024-04")
//  }
//  func Test_getSubmitterPRsForaneveux(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"aneveux","1","2024-04")
//  }
//  func Test_getSubmitterPRsForampuscas(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"ampuscas","1","2024-04")
//  }
//  func Test_getSubmitterPRsForalarreine(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"alarreine","1","2024-04")
//  }
//  func Test_getSubmitterPRsForakhilkolu(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"akhilkolu","1","2024-04")
//  }
//  func Test_getSubmitterPRsForabhishekshahqmetry(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"abhishekshah-qmetry","1","2024-04")
//  }
//  func Test_getSubmitterPRsForVlatombe(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Vlatombe","1","2024-04")
//  }
//  func Test_getSubmitterPRsForTobiX(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"TobiX","1","2024-04")
//  }
//  func Test_getSubmitterPRsForTheJonesFoundation(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"TheJonesFoundation","1","2024-04")
//  }
//  func Test_getSubmitterPRsForTheJonsey(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"The-Jonsey","1","2024-04")
//  }
//  func Test_getSubmitterPRsForTWestling(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"TWestling","1","2024-04")
//  }
//  func Test_getSubmitterPRsForSofiaVBuitrago(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"SofiaVBuitrago","1","2024-04")
//  }
//  func Test_getSubmitterPRsForRoyLu(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Roy-Lu","1","2024-04")
//  }
//  func Test_getSubmitterPRsForRirishi(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Ririshi","1","2024-04")
//  }
//  func Test_getSubmitterPRsForMunishkumar92(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Munishkumar92","1","2024-04")
//  }
//  func Test_getSubmitterPRsForMartinvH(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Martin-vH","1","2024-04")
//  }
//  func Test_getSubmitterPRsForLuisGuga(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Luis-Guga","1","2024-04")
//  }
//  func Test_getSubmitterPRsForLouey11(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Louey11","1","2024-04")
//  }
//  func Test_getSubmitterPRsForKevinCB(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Kevin-CB","1","2024-04")
//  }
//  func Test_getSubmitterPRsForJaspreet1601(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Jaspreet1601","1","2024-04")
//  }
//  func Test_getSubmitterPRsForIbraheemHaseeb7(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"IbraheemHaseeb7","1","2024-04")
//  }
//  func Test_getSubmitterPRsForGoooler(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Goooler","1","2024-04")
//  }
//  func Test_getSubmitterPRsForGerkinDev(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"GerkinDev","1","2024-04")
//  }
//  func Test_getSubmitterPRsForGOptimistic(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"GOptimistic","1","2024-04")
//  }
//  func Test_getSubmitterPRsForFrogDevelopper(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"FrogDevelopper","1","2024-04")
//  }
//  func Test_getSubmitterPRsForEfrenRey(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"EfrenRey","1","2024-04")
//  }
//  func Test_getSubmitterPRsForCJkrishnan(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"CJkrishnan","1","2024-04")
//  }
//  func Test_getSubmitterPRsForAviGupta1(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"Avi-Gupta1","1","2024-04")
//  }
//  func Test_getSubmitterPRsForAshitaSingamsetty(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"AshitaSingamsetty","1","2024-04")
//  }
//  func Test_getSubmitterPRsForAnoojM(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"AnoojM","1","2024-04")
//  }
//  func Test_getSubmitterPRsForAamirahP(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"AamirahP","1","2024-04")
//  }
//  func Test_getSubmitterPRsFor007PRAKHAR(t *testing.T) {
// 	check_getSubmittersPRfromGH(t,"007-PRAKHAR","1","2024-04")
//  }

func check_getSubmittersPRfromGH(t *testing.T, submittersName string, submittersPRs string, monthToSelectFrom string) {
	err, contributorData := getSubmittersPRfromGH(submittersName, submittersPRs, monthToSelectFrom)
	assert.NoError(t,err,"call failed for %s", submittersName)

	fmt.Printf("Getting details for %s (%s == %s) was successful", submittersName, submittersPRs, contributorData.totalPRs_found)
}
