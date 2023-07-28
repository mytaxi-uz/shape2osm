package tags

import (
	"strings"

	"github.com/mytaxi-uz/shape2osm/util/osm"
	"github.com/mytaxi-uz/shape2osm/util/shp"
)

func PoiAttrToOSMTag(num int, fields []shp.Field, reader *shp.Reader) (tags osm.Tags) {
	var key, value, str_type_uz, str_type_ru, str_type_en string

	for k, f := range fields {
		field := strings.ToUpper(f.String())
		switch field {
		case "TYP_STR_UZ":
			str_type_uz = reader.ReadAttribute(num, k)
		case "TYP_STR_RU":
			str_type_ru = reader.ReadAttribute(num, k)
		case "TYP_STR_EN":
			str_type_en = reader.ReadAttribute(num, k)
		}
	}

	for k, f := range fields {
		attr := reader.ReadAttribute(num, k)

		if attr == "" {
			continue
		}

		key = ""
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
		case "STREET_UZ":
			key = "addr:street:uz"
			value = attr
			if str_type_uz != "" {
				value += " " + str_type_uz
			}
		case "STREET_RU":
			key = "addr:street:ru"
			value = attr
			if str_type_ru != "" {
				value += " " + str_type_ru
			}
		case "STREET_EN":
			key = "addr:street:en"
			value = attr
			if str_type_en != "" {
				value += " " + str_type_en
			}
		case "ADDRESS_UZ":
			key = "addr:housenumber"
			value = attr
		// case "CLASS":
		// key = "amenity"
		// value = strings.ToLower(attr)
		case "ADDRESS":
			key = "addr:full"
			value = attr
		case "PHONE":
			key = "phone"
			value = attr
		case "WORKTIME":
			key = "opening_hours"
			value = attr
		case "HTTP":
			key = "website"
			value = attr
		case "EMAIL":
			key = "email"
			value = attr
		case "ZIPCODE":
			key = "addr:postcode"
			value = attr

		case "TYP_COD":
			switch attr {
			case "371":
				key = "aeroway"
				value = "aerodrome"
			case "727", "245", "248":
				key = "shop"
				value = "convenience"
			case "728":
				key = "shop"
				value = "clothes"
			case "734":
				key = "shop"
				value = "furniture"
			case "729":
				key = "shop"
				value = "hardware"
			case "171":
				key = "amenity"
				value = "bank"
			case "630":
				key = "amenity"
				value = "bar"
			case "613":
				key = "amenity"
				value = "hairdresser"
			case "243":
				key = "amenity"
				value = "cafe"
			case "609":
				key = "shop"
				value = "car"
			case "366":
				key = "shop"
				value = "car_repair"
			case "301":
				key = "amenity"
				value = "place_of_worship"
				tags = append(tags, osm.Tag{
					Key:   "religion",
					Value: "muslim",
				})
			case "303", "304", "622":
				key = "amenity"
				value = "place_of_worship"
				tags = append(tags, osm.Tag{
					Key:   "religion",
					Value: "christian",
				})
			case "327":
				key = "amenity"
				value = "cinema"
			case "213":
				key = "amenity"
				value = "college"
			case "615":
				key = "amenity"
				value = "dentist"
			case "362", "610":
				key = "amenity"
				value = "fuel"
			case "616", "624", "621":
				key = "office"
				value = "government"
			case "340", "350", "349":
				key = "amenity"
				value = "hospital"
			case "204":
				key = "amenity"
				value = "kindergarten"
			case "150":
				key = "tourism"
				value = "hotel"
			case "601", "250":
				key = "shop"
				value = "mall"
			case "333":
				key = "historic"
				value = "monument"
			case "329":
				key = "tourism"
				value = "museum"
			case "619", "709":
				key = "office"
				value = "company"
			case "103":
				key = "leisure"
				value = "park"
			case "240":
				key = "amenity"
				value = "pharmacy"
			case "247":
				key = "amenity"
				value = "restaurant"
			case "217", "218":
				key = "amenity"
				value = "school"
			case "733":
				key = "shop"
				value = "alcohol"
			case "731":
				key = "shop"
				value = "baby_goods"
			case "201":
				key = "amenity"
				value = "university"
			case "420", "440":
				key = "railway"
				value = "subway_entrance"
			case "325":
				key = "tourism"
				value = "zoo"
			case "271":
				key = "amenity"
				value = "swimming_pool"
			case "736":
				key = "shop"
				value = "computer"
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
