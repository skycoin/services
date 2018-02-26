package manifest

import (
	"io/ioutil"
	"fmt"
	_"log"
	"strconv"
	"bufio"
	"os"
	"strings"
	"encoding/csv"
	"time"
)


//SnapshotList show all snapshots in ./snapshots
func SnapshotList(flag int) []string {
	var str []string
	snapCount := 0
	files, err := ioutil.ReadDir("../snapshots")
	if err != nil {
		fmt.Println("can't find snapshots folder")
	} else {

		for _, f := range files {
			if f.Name()[len(f.Name())-3:] == "csv" {
				if flag == 1 {fmt.Println("Snapshot " + strconv.Itoa(snapCount) + ": " + f.Name()) }
				str = append(str, "../snapshots/"+f.Name())
				snapCount++
			}
		}
	}

	return str
}


//PromptCycle wait for user command
func PromptCycle() {
	for  {
		newCommand, args := InputFromCli()
		if newCommand == "" {
			continue
		}
		if newCommand == "merge" {
				snaps := SnapshotList(2)
				isCorrect, argsInt := CheckMergeArg(args, len(snaps))

				if isCorrect {
					MergeSnapshot(snaps, argsInt)
				} else {
					fmt.Println("bad args.")
				}
		}
	}
}

//InputFromCli format user command
func InputFromCli() (command string, args []string) {
	fmt.Println("Wait for command(example: merge 0 3)")
	command = ""
	args = []string{}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	splitInput := strings.Fields(input)
	if len(splitInput) == 0 {
		return
	}

	command = strings.Trim(splitInput[0], " ")
	if len(splitInput) > 1 {
		args = splitInput[1:]
	}
	return
}

//CheckMergeArg that args will be correct
func CheckMergeArg(args []string, snapCount int) (bool, []int) {
	var argsInt []int
	for _,arg := range args {
		k, e := strconv.Atoi(arg)
		if e != nil {
			return false, argsInt
		}
		if k > snapCount-1 || k < 0 {
			return false, argsInt
		}
		argsInt = append(argsInt, k)
	}

	return true, argsInt
}

//MergeSnapshot create new snapshot from two or more others
func MergeSnapshot(snaps []string,args []int) {
	mainSnap := ReadCVS(snaps[args[0]])


	args = args[1:]
	for _, arg := range args {
		nextSnap := ReadCVS(snaps[arg])
		mainSnap = Merge(mainSnap, nextSnap)
	}

	//write to csv merged info.
	csvFileName := time.Now().UTC().Format("2006-01-02 15:04:05") + ".csv"
	csvFile, err := os.Create("../snapshots/" + csvFileName)
	if err != nil {
		return
	}
	w := csv.NewWriter(csvFile)

	WriteIntoCSV(w, mainSnap)
	csvFile.Close()
}

//Merge two snapshots
func Merge(main []FileInfo, next []FileInfo) ([]FileInfo){
	flag := 0
	for _, newfile := range next {

		for i,oldfile := range main {
			if oldfile.Path == newfile.Path {
					main[i] = CompareFiles(oldfile, newfile)
					flag = 1
			}
		}
		if flag == 0 {
			main = append(main, newfile)
		} else {
			flag = 0
		}
	}


	return main
}

//CompareFiles check what file was modified early
func CompareFiles(old FileInfo, new FileInfo) FileInfo {
	t1, err := time.Parse("2006-01-02 15:04:05", old.Modified)
	if err != nil{
		fmt.Println("can't parse time")
	}
	t2, err := time.Parse("2006-01-02 15:04:05", new.Modified)
	if err != nil{
		fmt.Println("can't parse time")
	}
	if t1.Sub(t2) < 0*time.Second {
		return new
	}
	return old
}

//ReadCVS converts CSV into struct
func ReadCVS (path string) []FileInfo {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("bad path.")
	}
	r := csv.NewReader(file)
	records, _ := r.ReadAll()

	filesInfo := make([]FileInfo, len(records))
	for i, record := range records {
		filesInfo[i].Path = record[0]
		filesInfo[i].Name = record[1]
		filesInfo[i].Size = record[2]
		filesInfo[i].Modified = record[3]
		filesInfo[i].Hash = record[4]
	}

	return filesInfo
}

