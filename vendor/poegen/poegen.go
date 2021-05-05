package poegen

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	stoori   string = "moro %{ketk}s :D mitä %{ketk}s :D. tänää levutetaan %{skill}s %{suomiclass}s. tän buildin keksin ku olin %{mesta}s %{hommat}s. Isketään buildiin viel %{unique}s jonka löysin ku olin %{alueet}s farmaa tuomaan vähä semmost %{mausteet}s sekaa :D. kuten keisari izarokin sanoi: \"%{izaro}s\" nonii toivottavasti tykkäätte. mä rakastan teit %{ketk}s."
	datapath string = "datafiles"
	dmap     map[string][]string
)

func Tgen(format string, params map[string]string) string {
	for key, val := range params {
		format = strings.Replace(format, "%{"+key+"}s", val, -1)
	}
	return format
}

func populateKeys(files []fs.FileInfo) (map[string][]string, error) {
	ret := make(map[string][]string)
	for _, finfo := range files {
		n := strings.Split(finfo.Name(), ".")[0]
		p := datapath + "/" + finfo.Name()
		fmt.Println("read ", p)
		f, err := os.Open(p)
		defer f.Close()
		if err != nil {
			return nil, err
		}
		d, err := ioutil.ReadAll(f)

		ret[n] = strings.Split(string(d), "\n")
	}
	return ret, nil
}

func pickString(data map[string][]string) (map[string]string, error) {
	ret := make(map[string]string)
	for key, dlist := range data {
		if len(dlist) <= 0 {
			return nil, fmt.Errorf("no values for %s", key)
		}
		pick := dlist[rand.Int31n(int32(len(dlist)-1))]
		ret[key] = pick
		//fmt.Println("pick", pick)
	}
	return ret, nil
}

func init() {
	var err error
	datas, err := ioutil.ReadDir("datafiles")
	if err != nil {
		panic("datafiles")
	}
	dmap, err = populateKeys(datas)
	if err != nil {
		panic("pop map")
	}
}

func Generate() string {
	rand.Seed(time.Now().UnixNano())
	picks, err := pickString(dmap)
	fmt.Println(time.Now(), picks)
	if err != nil {
		panic("pick")
	}
	return Tgen(stoori, picks)
}
