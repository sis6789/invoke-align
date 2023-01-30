package handler

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/mitchellh/go-homedir"
)

var defaultValues = map[string]string{
	"adaptorP7":               "GATCGGAAGAGCACACGTCTGAACTCCAGTCAC",
	"bowtie2":                 "true",
	"clientFastqProcessRead":  "mongodb//localhost:27017",
	"clientFastqProcessWrite": "mongodb//localhost:27017",
	"clientReportRead":        "mongodb//localhost:27017",
	"clientReportWrite":       "mongodb//localhost:27017",
	"consensusMinimum":        "11",
	"cutPatternStrandMinus":   "CC[AT]",
	"cutPatternStrandPlus":    "[AT]GG",
	"email":                   "isseo@keyomics.co.kr",
	"env":                     "",
	"fastqLinePattern":        "^([ATGC]{8,12})(CTGACGTGAC([TA])GACT)([ACTG]{23,})$",
	"flatRatio":               "0.55",
	"maximumGapLength":        "151",
	"minimumFragmentLength":   "5",
	"minimumLength":           "21",
	"minimumStrandCount":      "2",
	"optimumGenomeLength":     "512",
	"run":                     "0",
	"serverAddress":           "localhost",
}

func Submit(c echo.Context) error {
	params := make(map[string]string)
	var err error
	err = c.Bind(&params)
	if err != nil {
		log.Printf("%v", err)
	}

	absPath := func(field string) {
		value := params[field]
		value, err = homedir.Expand(value)
		if err != nil {
			log.Printf("%v", err)
		}
		value, err = filepath.Abs(value)
		if err != nil {
			log.Printf("%v", err)
		}
		params[field] = value
	}
	if params["env"] != "" {
		absPath("env")
	}
	absPath("data")
	if params["folderBase"] != "" {
		absPath("folderBase")
	}

	var finalMap = make(map[string]string)
	for k, v := range params {
		defaultV, defaultExist := defaultValues[k]
		if defaultExist {
			if defaultV != v {
				finalMap[k] = v
			}
		} else {
			finalMap[k] = v
		}
	}
	var flagList []string
	for k := range finalMap {
		flagList = append(flagList, k)
	}
	sort.Strings(flagList)
	wString := ""
	for _, k := range flagList {
		v := finalMap[k]
		if v == "true" {
			wString += fmt.Sprintf(" --%s", k)
		} else {
			wString += fmt.Sprintf(" --%s \"%s\"", k, v)
		}
	}
	wString = wString[1:]

	if _, exist := params["env"]; exist {
		return c.JSONPretty(http.StatusOK, finalMap, "    ")
	} else {
		return c.String(http.StatusOK, wString)
	}
}
