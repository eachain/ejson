package ejson

import (
	"fmt"
	"strconv"
	"strings"
)

type entryFunc func(*JSON) *JSON

func invalidFormat(i int, s string) error {
	return fmt.Errorf("ejson: smart key %vth part invalid format: '%v'", i, s)
}

func parseSmartKey(keys string) ([]entryFunc, error) {
	var entries []entryFunc
	parts := strings.Split(keys, ".")
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		l := strings.IndexByte(part, '[')
		r := strings.IndexByte(part, ']')
		if l < 0 {
			if r >= 0 {
				return nil, invalidFormat(i, part)
			}

			entries = append(entries, func(j *JSON) *JSON {
				if j.IsStr() && j.StrIsJSON() {
					j = j.StrJSON()
				}
				if j.IsArray() {
					if idx, err := strconv.Atoi(part); err == nil {
						return j.ArrayIndex(idx)
					}
				}
				return j.ObjectIndex(part)
			})
			continue
		}

		if l == 0 {
			if i > 0 { // "[0]"
				return nil, invalidFormat(i, part)
			}
			if r <= 0 {
				return nil, invalidFormat(i, part)
			}
		}
		if l > 0 {
			if (r < 0) || (r > 0 && r <= l) {
				return nil, invalidFormat(i, part)
			}

			entries = append(entries, func(j *JSON) *JSON {
				if j.IsStr() && j.StrIsJSON() {
					j = j.StrJSON()
				}
				if j.IsArray() {
					if idx, err := strconv.Atoi(part[:l]); err == nil {
						return j.ArrayIndex(idx)
					}
				}
				return j.ObjectIndex(part[:l])
			})
		}

		tmp := part[l:]
		for tmp != "" {
			if tmp[0] != '[' {
				return nil, invalidFormat(i, part)
			}
			r = strings.IndexByte(tmp, ']')
			if r < 0 {
				return nil, invalidFormat(i, part)
			}
			idx, err := strconv.Atoi(tmp[1:r])
			if err != nil {
				return nil, fmt.Errorf("ejson: smart key %vth invalid index: '%v'",
					i, part)
			}
			entries = append(entries, func(j *JSON) *JSON {
				if j.IsStr() && j.StrIsJSON() {
					j = j.StrJSON()
				}
				return j.ArrayIndex(idx)
			})

			tmp = tmp[r+1:]
		}
	}
	return entries, nil
}
