package tracker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (t *Tracker) announceHttp(endpoint string) (error, *Response) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err, nil
	}

	q := req.URL.Query()
	q.Add("info_hash", url.QueryEscape(string(t.State.PeerId[:20])))
	q.Add("peer_id", url.QueryEscape(string(t.State.InfoHash[:20])))
	q.Add("ip", t.State.Ip)
	q.Add("port", fmt.Sprintf("%d", t.State.Port))
	q.Add("downloaded", fmt.Sprintf("%d", t.State.Downloaded))
	q.Add("uploaded", fmt.Sprintf("%d", t.State.Uploaded))
	q.Add("left", fmt.Sprintf("%d", t.State.Left))

	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.RawQuery)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err, nil
	}

	fmt.Println(string(body))

	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("%s : %d", "invalid statuscode from tracker", res.StatusCode)), nil
	}

	return nil, &Response{}
}
