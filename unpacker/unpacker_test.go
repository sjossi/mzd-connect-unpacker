package unpacker

import (
	"compress/gzip"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	ini "github.com/ochinchina/go-ini"
	"gotest.tools/assert"
)

const MAIN_INSTRUCTIONS = "main_instructions.ini"
const EXECUTE_INSTRUCTIONS = "test/execute.ini.gz"
const TEST_DIR_VAR = "MZD_TESTS_DIR"
const GOTEST_DIR_VAR = "MZD_GOTEST_DIR"

var path_test_dir string
var path_test_want string
var path_gotest_dir string

func TestMain(m *testing.M) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	var env_testdir_set bool
	path_test_dir, env_testdir_set = os.LookupEnv(TEST_DIR_VAR)
	if !env_testdir_set {
		log.Println("Falling back to default folder \"../../mzd-connect-unpacker-assets/cmu150_ADR_70.00.367A_update/\"")
		// Find this exact file online and unpack the .up.
		// Please don't ask me for the password, it's not mine. Thank you
		path_test_dir = "../../mzd-connect-unpacker-assets/cmu150_ADR_70.00.367A_update/"
	}

	// This is where all the JSON representation of the "want" objects are
	// stored. This allows for sensitive data to not be checked into git
	var env_gotestpath_set bool
	path_gotest_dir, env_gotestpath_set = os.LookupEnv(GOTEST_DIR_VAR)
	if !env_gotestpath_set {
		path_gotest_dir = "../../mzd-connect-unpacker-assets/gotest/"
	}

	code := m.Run()

	os.Exit(code)
}

func TestExtractFiles(t *testing.T) {
	t.Log("TestExtractFiles")
	// TODO: make cobra
	t.Skip("Not really a test, but a CLI replacement. Should be moved.")

	mainInstructions := filepath.Join(path_test_dir, "main_instructions.ini")
	toBase := "./extracted_" + time.Now().Format("20060102150405") + ""

	t.Log("Parsing ini tree")
	got := ParseIniTree(mainInstructions)

	t.Log("Recursively parsing ini tree")
	for _, ini := range got {
		if strings.HasPrefix(ini.Filename, "files.ini") || strings.HasPrefix(ini.Filename, "execute.ini") {
			t.Logf("Extracting %s to %s", ini.Filename, toBase)
			ExtractFiles(ini, toBase)
		}
	}
}

func TestSimulateFullTree(t *testing.T) {
	t.Log("TestSimulateFullTree")
	// TODO: make cobra
	t.Skip("Not really a test, but a bad CLI replacement.")

	mainInstructions := filepath.Join(path_test_dir, "main_instructions.ini")
	got := ParseIniTree(mainInstructions)

	files := make([]string, 0)
	for _, v := range got {
		files = append(files, SimulateSteps(v)...)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	pretty_files, err := json.MarshalIndent(files, "", "\t")
	check(err)
	t.Logf("%s", pretty_files)

	outfile, err := os.Create("file_list.txt")
	check(err)
	defer outfile.Close()

	for _, file := range files {
		outfile.WriteString(file + "\n")
	}

	t.Logf("tree: %#v", got)
}

func TestSimulateExecute(t *testing.T) {
	t.Log("TestSimulateExecute")
	// t.Skip("Skipping TestSimulateExecute")

	file, err := os.Open(path_test_dir + "bootstrap/execute.ini.gz")
	check(err)

	reader, err := gzip.NewReader(file)
	check(err)

	ini := ini.Load(reader)
	in := ParseSubIni(ini)
	got := SimulateSteps(in)

	want := []string{}
	readJson("TestSimulateExecute.json", &want)

	assert.DeepEqual(t, got, want)
}

func TestParseIniTree(t *testing.T) {
	t.Log("TestParseIniTree")

	got := ParseIniTree(filepath.Join(path_test_dir, "main_instructions.ini"))

	want := make([]*Ini, 0)
	readJson("TestParseIniTree.json", &want)

	assert.DeepEqual(t, got, want)
}

func TestParseMainIni(t *testing.T) {
	t.Log("TestParseMainIni")

	have := ini.Load(filepath.Join(path_test_dir, "main_instructions.ini"))
	got := ParseMainIni(have)

	want := &Ini{}
	readJson("TestParseMainIni.json", want)

	assert.DeepEqual(t, got, want)
}

func TestParseSettingsIni(t *testing.T) {
	t.Log("TestParseSettingsIni()")

	have := ini.Load(filepath.Join(path_test_dir, "main_instructions.ini"))
	got := ParseSettings(have)

	want := Settings{}
	readJson("TestParseSettingsIni.json", &want)

	assert.DeepEqual(t, got, want)
}

func TestParseDataStorage(t *testing.T) {
	t.Log("TestParseDataStorage")

	have := ini.Load(filepath.Join(path_test_dir, "main_instructions.ini"))
	got := ParseDataStorage(have)

	want := DataStorage{}
	readJson("TestParseDataStorage.json", &want)

	assert.DeepEqual(t, got, want)
}

func TestParseExecuteIni(t *testing.T) {
	t.Log("TestParseExecuteIni")

	file, err := os.Open(filepath.Join(path_test_dir, "compactwnn/execute.ini.gz"))
	check(err)

	gz, err := gzip.NewReader(file)
	check(err)

	in := ini.Load(gz)
	got := ParseSubIni(in)

	want := &Ini{}
	readJson("TestParseExecuteIni.json", want)

	assert.DeepEqual(t, got, want)
}

func TestParseInstructions(t *testing.T) {
	t.Log("TestParseInstructions")

	in := ini.Load(filepath.Join(path_test_dir, "main_instructions.ini"))

	got := ParseInstructions(in, "Instructions", true)

	want := Instructions{}
	readJson("TestParseInstructions_1.json", &want)

	assert.DeepEqual(t, got, want)

	got = ParseInstructions(in, "Instructions_Ext", true)
	readJson("TestParseInstructions_2.json", &want)

	assert.DeepEqual(t, got, want)
}

/*
Dump a JSON representation of an object to a file.

This is used to export the "want" objects into a json file, so it stays out of
the source code.
*/
func dumpJson(filename string, v any) {
	j, err := json.Marshal(v)
	check(err)

	// log.Printf("dumping: %s", j)

	os.WriteFile(filepath.Join(path_gotest_dir, filename), j, 0755)
}

func readJson[T any](filename string, v *T) {
	contents, err := os.ReadFile(filepath.Join(path_gotest_dir, filename))
	check(err)

	err = json.Unmarshal(contents, v)
	check(err)
}
