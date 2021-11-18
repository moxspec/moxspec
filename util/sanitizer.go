package util

import (
	"regexp"
	"strings"
)

// SanitizeString returns string with all unuseful words
func SanitizeString(str string) string {
	s := str
	s = strings.Replace(s, "(TM)", "", -1)
	s = strings.Replace(s, "(R)", "", -1)
	s = strings.Replace(s, "\x00", "", -1)
	s = strings.Replace(s, "To be filled by O.E.M.", "", -1)
	s = strings.Replace(s, "To Be Filled By O.E.M.", "", -1)
	return strings.TrimSpace(s)
}

// ShortenVendorName returns short version vendor name
func ShortenVendorName(str string) string {
	s := str
	s = strings.Replace(s, "Broadcom Limited", "Broadcom", -1)
	s = strings.Replace(s, "Broadcom and subsidiaries", "Broadcom", -1)
	s = strings.Replace(s, "Broadcom Inc. and subsidiaries", "Broadcom", -1)
	s = strings.Replace(s, "Advanced Micro Devices, Inc. [AMD]", "AMD", -1)
	s = strings.Replace(s, "LSI Logic / Symbios Logic", "LSI/Symbios", -1)
	s = strings.Replace(s, "American Megatrends Inc.", "AMI", -1)
	s = strings.Replace(s, "Hewlett-Packard Company", "HP", -1)
	s = strings.Replace(s, "Mellanox Technologies", "Mellanox", -1)
	s = strings.Replace(s, "Quanta Computer Inc", "Quanta", -1)
	s = strings.Replace(s, "Hynix Semiconductor", "Hynix", -1)
	s = strings.Replace(s, "PenguinComputing", "Penguin Computing", -1)
	s = strings.Replace(s, " Electronics Systems Ltd.", "", -1)
	s = strings.Replace(s, " Technologies Co., Ltd.", "", -1)
	s = strings.Replace(s, " Electronics Co Ltd", "", -1)
	s = strings.Replace(s, " Technology, Inc.", "", -1)
	s = strings.Replace(s, " Technologies LTD", "", -1)
	s = strings.Replace(s, " Corporation", "", -1)
	s = strings.Replace(s, " Corp.", "", -1)
	s = strings.Replace(s, ", Inc.", "", -1)
	s = strings.Replace(s, ", Inc", "", -1)
	s = strings.Replace(s, " Inc.", "", -1)
	return SanitizeString(s)
}

// ShortenProcName returns short version processor name
func ShortenProcName(str string) string {
	s := str
	s = strings.Replace(s, " CPU", "", -1)
	s = strings.Replace(s, " 0 @ ", " ", -1)
	s = strings.Replace(s, " @ ", " ", -1)
	s = regexp.MustCompile(` [0-9]+-Core Processor`).ReplaceAllString(s, "") // Zen
	s = regexp.MustCompile(`  +`).ReplaceAllString(s, " ")                   // Westmere
	return SanitizeString(s)
}
