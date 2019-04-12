package torrent

import (
	"crypto/sha1"
	"errors"
	"github.com/zeebo/bencode"
)

const MaxChunk = 16384


type Torrent struct {
	Announce     string
	AnnounceList []string
	Info         Info
	Comment      string
	CreatedBy    string
	CreatedAt    int64
	InfoHash     [20]byte
	Meta         struct {
		TargetIndex    int
		SubIndex       int
		HasVideo       bool
		HasAudio       bool
		IsSingle       bool
		HasSub         bool
		TargetFileName string
		Downloaded     int64 //bytes
		Uploaded       int64 //bytes
		SelectedPieces []*Piece
	}
}

func Create(data map[string]interface{}) (error, *Torrent) {
	torrent := Torrent{}

	for k, v := range data {
		switch k {
		case "info":
			bytes, err := bencode.EncodeBytes(v)
			if err != nil {
				return err, nil
			}

			torrent.InfoHash = sha1.Sum(bytes)
			torrent.Info = parseInfo(v)
		case "announce":
			torrent.Announce = v.(string)
		case "announce-list":
			torrent.AnnounceList = parseAnnounceList(v)
		case "comment":
			torrent.Comment = v.(string)
		case "creation date":
			torrent.CreatedAt = v.(int64)
		case "created by":
			torrent.CreatedBy = v.(string)
		}
	}

	err := postProcess(&torrent)
	if err != nil {
		return err, nil
	}

	return nil, &torrent
}

func (t *Torrent) SelectedFile() (*File) {
	return t.Info.Files[t.Meta.TargetIndex]
}

func postProcess(t *Torrent) error {
	t.Meta.Downloaded = 0
	t.Meta.Uploaded = 0

	hasVid, vidIndex := t.Info.Files.hasVideo()
	hasAudio, audIndex := t.Info.Files.hasAudio()

	if !hasAudio && !hasVid {
		return errors.New("no file found")
	}

	if hasVid {
		t.Meta.TargetIndex = vidIndex
		t.Meta.TargetFileName = t.Info.Files[vidIndex].Path[len(t.Info.Files[vidIndex].Path)-1]
		t.Meta.HasVideo = true
	}

	if hasAudio {
		t.Meta.TargetIndex = audIndex
		t.Meta.TargetFileName = t.Info.Files[vidIndex].Path[len(t.Info.Files[vidIndex].Path)-1]
		t.Meta.HasAudio = true
	}

	hasSub, subIndex := t.Info.Files.hasSub()
	if hasSub {
		t.Meta.SubIndex = subIndex
		t.Meta.HasSub = true
	}

	file := t.SelectedFile()
	filePos := int64(0)

	chunkSize := t.Info.PieceLength / MaxChunk

	for _, p := range t.Info.Pieces {
		p.ByteOffset = filePos

		filePos += t.Info.PieceLength

		if filePos >= file.Start && filePos <= file.End {
			t.Meta.SelectedPieces = append(t.Meta.SelectedPieces, p)
		}
	}

	return nil
}
