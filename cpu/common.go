package cpu

import (
	"strconv"
	"strings"
)

func parseListString(s string) []uint16 {
	var res []uint16
	for _, c := range strings.Split(s, ",") {
		chunk := strings.TrimSpace(c)

		if strings.Contains(chunk, "-") {
			r := parseRangeString(chunk)
			if r != nil {
				res = append(res, r...)
			}
		} else {
			num, err := strconv.ParseInt(chunk, 10, 16)
			if err != nil {
				log.Warnf(err.Error())
				continue
			}
			res = append(res, uint16(num))
		}
	}
	return res
}

func parseRangeString(s string) []uint16 {
	e := strings.Split(s, "-")
	if len(e) != 2 {
		log.Warnf("invalid range format: %s", s)
		return nil
	}

	start, err := strconv.ParseInt(e[0], 10, 16)
	if err != nil {
		log.Warnf(err.Error())
		return nil
	}

	end, err := strconv.ParseInt(e[1], 10, 16)
	if err != nil {
		log.Warnf(err.Error())
		return nil
	}

	if start == end {
		return []uint16{uint16(start)}
	}

	if start > end {
		tmp := start
		start = end
		end = tmp
	}

	var res []uint16
	for i := start; i <= end; i++ {
		res = append(res, uint16(i))
	}
	return res
}
