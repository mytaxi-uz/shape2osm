package tags

import (
	"strings"

	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

func RailwayAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key := ""
		value := ""
		field := strings.ToUpper(f.String())

		switch field {
		case "NAME_UZ":
			key = "name:uz"
			value = attr
		case "NAME", "NAME_RU":
			key = "name:ru"
			value = attr
		case "NAME_EN":
			key = "name:en"
			value = attr
		case "TYP_COD":
			key = "railway"
			switch attr {
			case "420":
				value = "subway"
			case "442":
				value = "rail"
			default:
				key = ""
			}
		}

		if key != "" {
			tags = append(tags, osm.Tag{
				Key:   key,
				Value: value,
			})
		}
	}

	return
}
