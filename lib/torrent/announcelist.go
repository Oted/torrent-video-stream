package torrent

import "reflect"

func parseAnnounceList(i interface{}) (list []string) {
	if reflect.TypeOf(i).Kind() == reflect.Slice {
		s := reflect.ValueOf(i)

		for i := 0; i < s.Len(); i++ {
			val := reflect.ValueOf(s.Index(i).Interface())
			if val.Kind() != reflect.Slice {
				continue
			}

			elem := val.Index(0)
			if elem.Kind() == reflect.Interface {
				list = append(list, elem.Interface().(string))
			}
		}
	}

	return
}
