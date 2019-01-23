package torrent

import (
	"reflect"
)

type Info struct {
	Files       Files
	Name        string
	PieceLength int64
	Pieces      Pieces
	Private     bool
}

func parseInfo(i interface{}) (info Info) {
	if reflect.TypeOf(i).Kind() == reflect.Map {
		info.Files = Files{}

		for _, k := range reflect.ValueOf(i).MapKeys() {
			val := reflect.ValueOf(i).MapIndex(k)

			switch k.String() {
			case "files":
				info.Files = parseFiles(val.Interface())
			case "piece length":
				info.PieceLength = reflect.ValueOf(val.Interface()).Int()
			case "pieces":
				info.Pieces = parsePieces(val.Interface())
			case "private":
				//TODO
			case "length": //if its a single file
				f := parseFile(reflect.ValueOf(i))
				f.Length = reflect.ValueOf(val.Interface()).Int()
				f.Start = 0
				f.End = f.Length
				info.Files = append(info.Files, f)
			}
		}
	}

	return
}
