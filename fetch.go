// +build !appengine

package gosubway

import (
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	_ "github.com/jprobinson/gtfs/nyct_subway_proto"
	"github.com/jprobinson/gtfs/transit_realtime"
	"golang.org/x/net/context"
)

// GetFeed takes an API key generated from http://datamine.mta.info/user/register
// and a boolean specifying which feed (1,2,3,4,5,6,S trains OR L train) and
// it will return a transit_realtime.FeedMessage with NYCT extensions.
func GetFeed(_ context.Context, key string, lTrain bool) (*FeedMessage, error) {
	url := "http://datamine.mta.info/mta_esi.php?key=" + key
	if lTrain {
		url = url + "&feed_id=2"
	} else {
		url = url + "&feed_id=1"
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	transit := &transit_realtime.FeedMessage{}
	err = proto.Unmarshal(body, transit)
	if err != nil {
		return nil, err
	}
	return &FeedMessage{*transit}, nil
}
