package input

import (
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"github.com/zeebo/bencode"
	"io/ioutil"
)

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

//magnet:?xt=
// urn:btih:296ee7845e915a70e24df2b173c3f8c06fad98f5
// &dn=South.Park.S22E07.1080p.HDTV.x264-CRAVERS
// &tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969
// &tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80
// &tr=udp%3A%2F%2Fopen.demonii.com%3A1337
// &tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969

func (i *input) torrentFromMagnet(url string) (error, *torrent.Torrent) {
	return nil, nil
}
