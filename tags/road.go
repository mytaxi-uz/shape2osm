package tags

import (
	"strings"

	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

func RoadAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	typCod := false
	for _, f := range fields {
		field := strings.ToUpper(f.String())
		if field == "TYP_COD" {
			typCod = true
			break
		}
	}

	if !typCod {
		tags = append(tags, osm.Tag{
			Key:   "highway",
			Value: "road",
		})
		return
	}

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key := ""
		value := ""
		field := strings.ToUpper(f.String())

		switch field {
		/*
			case "ID":
				key = "id"
				value = attr
		*/
		case "NAME_UZ":
			key = "name:uz"
			value = attr
		case "NAME", "NAME_RU":
			key = "name:ru"
			value = attr
		case "NAME_EN":
			key = "name:en"
			value = attr
		case "SPEED":
			key = "additional:maxspeed"
			value = attr
		case "DIRECTION":
			if attr == "FT" {
				key = "oneway"
				value = "yes"
			}
		case "TUNNEL":
			if attr == "2" {
				key = "tunnel"
				value = "yes"
				tags = append(tags, osm.Tag{
					Key:   "layer",
					Value: "-1",
				})
			} else if attr == "1" {
				key = "bridge"
				value = "yes"
				tags = append(tags, osm.Tag{
					Key:   "layer",
					Value: "1",
				})
			}
		case "TYP_COD":
			key = "highway"
			switch attr {
			case "34":
				value = "trunk"
			case "31":
				value = "trunk_link"
			case "49", "54":
				value = "primary"
			case "53":
				value = "secondary"
			case "453", "51", "55":
				value = "tertiary"
			case "45", "50":
				value = "residential"
			case "40", "46":
				value = "service"
			// case "33":
			// value = "pedestrian"
			// case "468", "36":
			// value = "track"
			case "146", "32":
				value = "footway"
			case "59": // crossing
				value = "footway"
				tags = append(tags, osm.Tag{
					Key:   "footway",
					Value: "crossing",
				})
			case "431": // footway tonnel
				value = "footway"
				tags = append(tags, osm.Tag{
					Key:   "layer",
					Value: "-1",
				})
				tags = append(tags, osm.Tag{
					Key:   "tunnel",
					Value: "yes",
				})
			case "145", "405": // steps
				value = "steps"
			default:
				value = "road"
			}
			/*
				case "CLASS":
					key = "highway"
					switch attr {
					case "1":
						value = "motorway"
					case "2":
						value = "trunk"
					case "3":
						value = "primary"
					case "4":
						value = "secondary"
					case "5":
						value = "tertiary"
					case "6":
						value = "unclassified"
					case "7":
						value = "residential"
					case "8":
						value = "living_street"
					}
			*/
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
