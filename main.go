package main

import (
	"errors"
	"fmt"
	"github.com/Oted/torrentz/lib/client"
	"github.com/Oted/torrentz/lib/torrent"
	tracker2 "github.com/Oted/torrentz/lib/tracker"
	"github.com/zeebo/bencode"
	"io/ioutil"
	"os"
)

type Input struct {
	Path string
}

func main() {
	if len(os.Args) < 2 {
		panic(errors.New("no path"))
	}

	input := NewInput(os.Args)

	var torrent *torrent.Torrent
	var err error

	if input.Path != "" {
		err, torrent = torrentFromPath(input.Path)
		if err != nil {
			panic(err)
		}
	}

	err, client := client.New()
	if err != nil {
		panic(err)
	}

	err, tracker := tracker2.Create(torrent, client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("+%v\n", client)
	fmt.Printf("+%v\n", torrent)
	fmt.Printf("+%v\n", tracker)
}

//implement parser
func NewInput(args []string) Input {
	return Input{
		Path: args[1],
	}
}

func torrentFromPath(path string) (error, *torrent.Torrent) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err, nil
	}

	var data map[string]interface{}

	err = bencode.DecodeBytes(file, &data)
	if err != nil {
		return err, nil
	}

	err, torrent := torrent.Create(data)
	if err != nil {
		return err, nil
	}

	return nil, torrent
}
