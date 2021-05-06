package poegen

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

type items struct {
	Result []cat `json:"result"`
}
type cat struct {
	ID      string  `json:"id"`
	Label   string  `json:"label"`
	Entries []entry `json:"entries"`
}

type entry struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Text  string `json:"text"`
	Flags struct {
		Unique bool `json:"unique"`
	} `json:"flags"`
}

var (
	stoori   string = "moro %{ketk}s :D mitä %{ketk}s :D. tänää levutetaan %{suomiskill}s %{suomiclass}s. tän buildin keksin ku olin %{mesta}s %{hommat}s. buildin perusrunko on %{iteminosa}s jonka %{hankinta}s %d %{rahat}s. isketään buildiin viel %{unique}s jonka löysin ku olin %{alueet}s farmaa tuomaan vähä semmost %{mausteet}s sekaa :D. kuten keisari izarokin sanoi: \"%{izaro}s\" nonii toivottavasti tykkäätte. mä rakastan teit %{ketk}s."
	datapath string = "datafiles"
	dmap     map[string][]string
	uniqs    []string
)

func Tgen(format string, params map[string]string) string {
	for key, val := range params {
		format = strings.Replace(format, "%{"+key+"}s", val, -1)
	}
	return fmt.Sprintf(format, rand.Int31n(666))
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

// refresh eng itemnames
func parseItemJson(fname string) {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	var idata items
	d, err := ioutil.ReadAll(f)
	if err != nil {
		panic("json file")
	}
	f.Close()
	f, err = os.Create("engdata/items.txt")
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(d, &idata); err != nil {
		panic("json")
	}
	var data []cat
	data = append(data, idata.Result[0])
	data = append(data, idata.Result[1])
	data = append(data, idata.Result[8])
	var usons []entry
	for _, c := range data {
		for _, item := range c.Entries {
			if item.Flags.Unique {
				usons = append(usons, item)
			}
		}
	}
	for _, u := range usons {
		uniqs = append(uniqs, u.Name)
		fmt.Fprintln(f, u.Name)
	}
	f.Close()
	fmt.Println(uniqs)
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
	parseItemJson("items.json")
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
