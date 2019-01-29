package torrent

import (
	"crypto/sha1"
	"errors"
	"github.com/zeebo/bencode"
)

type Torrent struct {
	Announce     string
	AnnounceList []string
	Info         Info
	Comment      string
	CreatedBy    string
	CreatedAt    int64
	InfoHash     [20]byte
	Meta         struct {
		VidIndex   int
		SubIndex   int
		IsSingle   bool
		HasSub     bool
		VidFile    string
		Downloaded int64 //bytes
		Uploaded   int64 //bytes
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

	err := meta(&torrent)
	if err != nil {
		return err, nil
	}

	return nil, &torrent
}

func (t *Torrent) SelectVideoPieces() (err error, pi []*Piece) {
	file := t.SelectVideoFile()
	filePos := int64(0)

	for _, p := range t.Info.Pieces {
		if filePos >= file.Start && filePos <= file.End {
			pi = append(pi, p)
		}

		filePos += t.Info.PieceLength
	}

	return
}

func (t *Torrent) SelectVideoFile() (*File) {
	return t.Info.Files[t.Meta.VidIndex]
}

func meta(t *Torrent) error {
	t.Meta.Downloaded = 0
	t.Meta.Uploaded = 0

	hasVid, vidIndex := t.Info.Files.hasVideo()
	if !hasVid {
		return errors.New("no video found")
	}

	t.Meta.VidIndex = vidIndex
	t.Meta.VidFile = t.Info.Files[vidIndex].Path[len(t.Info.Files[vidIndex].Path)-1]

	hasSub, subIndex := t.Info.Files.hasSub()
	if hasSub {
		t.Meta.SubIndex = subIndex
		t.Meta.HasSub = true
	}

	return nil
}
