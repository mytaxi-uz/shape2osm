package tags

import (
	"strings"

	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

func WaterAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	var key, value string

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key = ""
		field := strings.ToUpper(f.String())

		switch field {
		case "TYP_COD":
			key = "water"
			switch attr {
			case "62":
				value = "reservoir"
				tag := osm.Tag{
					Key:   "natural",
					Value: "water",
				}
				tags = append(tags, tag)
			case "65", "64":
				value = "canal"
			case "61", "68", "74":
				value = "river"
			case "75":
				value = "stream"
			case "271", "72", "73":
				key = "leisure"
				value = "swimming_pool"
			default:
				key = ""
			}
		case "NAME_UZ":
			key = "name:uz"
			value = attr
		case "NAME", "NAME_RU":
			key = "name:ru"
			value = attr
		case "NAME_EN":
			key = "name:en"
			value = attr
		}

		if key != "" {
			tag := osm.Tag{
				Key:   key,
				Value: value,
			}
			tags = append(tags, tag)
		}
	}

	return
}
