package torrent

import (
	"reflect"
	"strings"
)

type File struct {
	Length int64
	Path   []string
	Md5Sum *string
	Start  int64
	End    int64
	Mime   string
}
type Files []*File

func parseFiles(i interface{}) (files Files) {
	sum := int64(0)
	if reflect.TypeOf(i).Kind() == reflect.Slice {
		s := reflect.ValueOf(i)

		for i := 0; i < s.Len(); i++ {
			mapVal := reflect.ValueOf(s.Index(i).Interface())

			if mapVal.Kind() != reflect.Map {
				continue
			}

			f := parseFile(mapVal)
			f.Start = sum
			sum += f.Length
			f.End = sum

			files = append(files, f)
		}
	}

	return
}

func parseFile(mapVal reflect.Value) *File {
	file := File{}

	for _, k := range mapVal.MapKeys() {
		switch k.String() {
		case "path":
			i := mapVal.MapIndex(k).Interface()
			file.Path = parseFilePath(i)
		case "name" :
			i := mapVal.MapIndex(k).Interface()
			file.Path = []string{reflect.ValueOf(i).String()}
		case "length":
			file.Length = mapVal.MapIndex(k).Interface().(int64)
		case "md5sum":
			md := mapVal.MapIndex(k).Interface().(string)
			file.Md5Sum = &md
		}
	}

	return &file
}

func parseFilePath(i interface{}) (paths []string) {
	for index := 0; index < reflect.ValueOf(i).Len(); index++ {
		paths = append(paths, reflect.ValueOf(i).Index(index).Interface().(string))
	}

	return
}

func (f Files) hasSub() (bool, int) {
	for i, file := range f {
		if strings.Contains(file.Path[len(file.Path)-1], ".srt") {
			return true, i
		}
	}

	return false, -1
}


func (f Files) hasAudio() (bool, int) {
	for i, file := range f {
		if strings.Contains(strings.ToLower(file.Path[len(file.Path)-1]), ".mp3") {
			file.Mime = "audio/mpeg"
			return true, i
		}

		if strings.Contains(strings.ToLower(file.Path[len(file.Path)-1]), ".wav") {
			file.Mime = "audio/wav"
			return true, i
		}

		if strings.Contains(strings.ToLower(file.Path[len(file.Path)-1]), ".flac") {
			file.Mime = "audio/wav"
			return true, i
		}
	}

	return false, -1
}

func (f Files) hasVideo() (bool, int) {
	for i, file := range f {
		if strings.Contains(strings.ToLower(file.Path[len(file.Path)-1]), ".mkv") {
			file.Mime = "video/webm"
			return true, i
		}

		if strings.Contains(strings.ToLower(file.Path[len(file.Path)-1]), ".mp4") {
			file.Mime = "video/mp4"
			return true, i
		}

		if strings.Contains(strings.ToLower(file.Path[len(file.Path)-1]), ".avi") {
			file.Mime = "video/x-msvideo"
			return true, i
		}

		if strings.Contains(strings.ToLower(file.Path[len(file.Path)-1]), ".ogg") {
			file.Mime = "video/ogg"
			return true, i
		}

		if strings.Contains(strings.ToLower(file.Path[len(file.Path)-1]), ".vebm") {
			file.Mime = "video/webm"
			return true, i
		}
	}

	return false, -1
}
