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
func ParseItemJson(fname string) {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	d, err := ioutil.ReadAll(f)
	if err != nil {
		panic("json file")
	}
	f.Close()
	var idata items
	if err = json.Unmarshal(d, &idata); err != nil {
		panic("json")
	}
	writeUniqItemFile(idata, "engdata/unique.txt")
	writeMapsNameFile(idata, "engdata/maps.txt")
}

func writeUniqItemFile(d items, outfile string) {
	i := []int{0, 1, 8}
	writeDataFile(d, i, outfile)
}

func writeMapsNameFile(d items, outfile string) {
	i := []int{7}
	writeDataFile(d, i, outfile)
}

func writeDataFile(d items, e []int, fname string) {
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	var data []cat
	for _, idx := range e {
		data = append(data, d.Result[idx])
	}
	var es []entry
	for _, c := range data {
		for _, item := range c.Entries {
			if item.Flags.Unique {
				es = append(es, item)
			}
		}
	}
	for _, u := range es {
		fmt.Fprintln(f, strings.ToLower(u.Text))
	}
	f.Close()
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
