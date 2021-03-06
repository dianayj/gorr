package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestScan(t *testing.T) {
	*TestCaseConfigPattern = "config.json"
	items, err := ScanTestData("./testdata")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(items))

	assert.Equal(t, "input/global1.db", items[0].DB[0])
	assert.Equal(t, "input/global2.db", items[1].DB[0])

	assert.Equal(t, 1, len(items[0].Input))
	assert.Equal(t, 2, len(items[1].Input))
	assert.Equal(t, 2, len(items[0].TestCases))
	assert.Equal(t, 2, len(items[1].TestCases))

	assert.Equal(t, "input/pos1.data", items[0].Input[0].Src)
	assert.Equal(t, "/var/data/conf/ad_pos.done", items[0].Input[0].Dst)

	assert.Equal(t, "input/pos2.data", items[1].Input[0].Src)
	assert.Equal(t, "/var/data/conf/ad_pos.done", items[1].Input[0].Dst)
	assert.Equal(t, "input", items[1].Input[1].Src)
	assert.Equal(t, "output", items[1].Input[1].Dst)

	assert.Equal(t, "access_client1", items[0].TestCases[0].Runner)
	assert.Equal(t, "input/access_req1.json", items[0].TestCases[0].Req)
	assert.Equal(t, "output/access_rsp1.dat", items[0].TestCases[0].Rsp)
	assert.Equal(t, "mixer_client1", items[0].TestCases[1].Runner)
	assert.Equal(t, "input/access_req2.json", items[0].TestCases[1].Req)
	assert.Equal(t, "output/access_rsp2.dat", items[0].TestCases[1].Rsp)

	assert.Equal(t, "access_client2", items[1].TestCases[0].Runner)
	assert.Equal(t, "input/access_req1.json", items[1].TestCases[0].Req)
	assert.Equal(t, "output/access_rsp1.dat", items[1].TestCases[0].Rsp)
	assert.Equal(t, "mixer_client2", items[1].TestCases[1].Runner)
	assert.Equal(t, "input/access_req2.json", items[1].TestCases[1].Req)
	assert.Equal(t, "output/access_rsp2.dat", items[1].TestCases[1].Rsp)
}

func TestRunCase(t *testing.T) {
	*TestCaseConfigPattern = "config.json"
	item := &TestItem{
		Path:  "./dummy.config",
		DB:    []string{"testdata/dummy.rsp"},
		Flags: []string{"-h"},
		Input: []MoveData{MoveData{Src: "testdata/dummy.rsp", Dst: "/tmp/dummy.rsp"}},
		TestCases: []*TestCase{
			&TestCase{Req: "dummy.req", Rsp: "testdata/dummy.rsp", Desc: "dummy test", Runner: "/bin/echo"}, // succ
			&TestCase{Req: "dummy.req", Rsp: "dummy2.rsp", Desc: "dummy test", Runner: "/bin/echo"},         // fail
			&TestCase{Req: "dummy.req", Rsp: "dummy.rsp", Desc: "dummy test", Runner: "/bin/echo2"},         // fail
		},
	}

	nt, ret := RunTestCase("./rdiff", "/bin/ls", "/bin/ls", "1.1.1.1:233", "./testdata", "/tmp", "./testdata/f.test.flag", item)

	assert.Equal(t, 2, ret.Fail)

	fmt.Printf("total err:%d\nmsg:%+v\ndiff:%+v", ret.Fail, ret.Msg, ret.Diff)

	os.Remove(item.Path)
	os.Remove("./testdata/f.test.flag")

	for _, t := range nt {
		os.Remove(t.Path)
	}
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags

	_, err := exec.Command("/bin/sh", "-c", "go build -o rdiff ./diff/diff.go ./diff/diffValue.go").Output()
	if err != nil {
		fmt.Printf("build diff tool failed\n")
		os.Exit(23)
	}

	ret := m.Run()
	func() { os.Remove("./rdiff") }()

	os.Exit(ret)
}
