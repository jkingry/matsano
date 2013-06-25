package histogram

import (
	"bitbucket.org/jkingry/matsano/cmd"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

type Histogram map[byte]float64

var Letters Histogram = Histogram(map[byte]float64{
	'E': 12.02,
	'T': 9.10,
	'A': 8.12,
	'O': 7.68,
	'I': 7.31,
	'N': 6.95,
	'S': 6.28,
	'R': 6.02,
	'H': 5.92,
	'D': 4.32,
	'L': 3.98,
	'U': 2.88,
	'C': 2.71,
	'M': 2.61,
	'F': 2.30,
	'Y': 2.11,
	'W': 2.09,
	'G': 2.03,
	'P': 1.82,
	'B': 1.49,
	'V': 1.11,
	'K': 0.69,
	'X': 0.17,
	'Q': 0.11,
	'J': 0.10,
	'Z': 0.07,
})

func Make() Histogram {
	return Histogram(make(map[byte]float64))
}

func FromFile(path string) (h Histogram, e error) {
	var file *os.File
	if file, e = os.Open(path); e != nil {
		return
	}
	defer file.Close()

	read := 0
	buffer := make([]byte, 4096)
	h = Make()
	for {
		if read, e = file.Read(buffer); e != nil || read == 0 {
			if e == io.EOF {
				e = nil
			}
			return
		}

		for i := 0; i < read; i++ {
			h[buffer[i]] += 1
		}
	}
}

func New(in []byte) Histogram {
	frequency := Make()

	for _, v := range in {
		frequency[v] += 1
	}

	return frequency
}

func (x Histogram) Add(y Histogram) {
	for k, v := range y {
		x[k] += v
	}
}

func (h Histogram) String() string {
	data, _ := json.Marshal(h)
	return string(data)
}

func (h Histogram) Scan(state fmt.ScanState, verb rune) (e error) {
	braceCount := 0
	matchBrace := func(r rune) bool {
		if r == '{' {
			braceCount += 1
		} else if r == '{' {
			braceCount -= 1
		}

		return braceCount > 0
	}

	var data []byte
	if data, e = state.Token(true, matchBrace); e != nil {
		return
	}

	json.Unmarshal(data, &h)

	return
}

func (h Histogram) MarshalJSON() ([]byte, error) {
	export := make(map[string]float64, len(h))
	for k, v := range h {
		export[strconv.FormatUint(uint64(k), 10)] = v
	}

	return json.Marshal(export)
}

func (h Histogram) UnmarshalJSON(in []byte) (e error) {
	export := make(map[string]float64)
	e = json.Unmarshal(in, &export)
	if e != nil {
		return
	}

	for k, v := range export {
		var n uint64
		n, e = strconv.ParseUint(k, 10, 8)
		if e != nil {
			return
		}
		h[byte(n)] = v
	}

	return
}

func (h Histogram) UpperCase() {
	for k, v := range h {
		if k >= 'a' && k <= 'z' {
			h[k-32] += v
		}
		if k < 'A' || k > 'Z' {
			delete(h, k)
		}
	}
}

func (x Histogram) ChiSquaredDistance(y Histogram) float64 {
	sum := 0.0
	for i, _ := range x {
		sum += math.Pow(x[i]-y[i], 2) / (x[i] + y[i])
	}

	return sum / 2
}

func (h Histogram) Filter(target Histogram) {
	for k, _ := range h {
		if _, ok := target[k]; !ok {
			delete(h, k)
		}
	}
}

func (h Histogram) Normalize() {
	sum := 0.0
	for _, v := range h {
		sum += v
	}
	for k, v := range h {
		h[k] = v / sum
	}
}

var Commands *cmd.Command = cmd.NewCommand("histogram", "histogram commands")

func init() {
	create := Commands.Add("create", "[input]...")
	create.Command = func(args []string) {
		h := Make()
		for i := 0; i < len(args); i++ {
			hf := New([]byte(cmd.GetInput(args, i)))
			h.Add(hf)
		}

		fmt.Print(h)
	}

	norm := Commands.Add("normalize", "normalize histogram (from standard input)")
	norm.Command = func(args []string) {
		h := Make()
		if n, e := fmt.Scan(&h); n != 1 {
			fmt.Fprint(os.Stderr, e)
			return
		}

		h.Normalize()

		fmt.Print(h)
	}

	gofmt := Commands.Add("goFormat", "format histogram (from standard input) as GO structure")
	gofmt.Command = func(args []string) {
		h := Make()
		if n, e := fmt.Scan(&h); n != 1 {
			fmt.Fprint(os.Stderr, e)
			return
		}

		m := map[byte]float64(h)

		fmt.Printf("histogram.Histogram(%#v)", m)
	}
}
