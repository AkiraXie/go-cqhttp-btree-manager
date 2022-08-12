package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AkiraXie/go-cqhttp-btree-manager/base"
)

var sets []*flag.FlagSet

func NewFlagSet(name string, errorHandling flag.ErrorHandling) (set *flag.FlagSet) {
	set = flag.NewFlagSet(name, errorHandling)
	sets = append(sets, set)
	return
}

func main() {

	h := flag.Bool("h", false, "usage of command")
	insertCmd := NewFlagSet("insert", flag.ExitOnError)
	insertImg := insertCmd.String("f", "", "cache file path")
	insertDb := insertCmd.String("o", "", "cache db path")

	selectCmd := NewFlagSet("select", flag.ExitOnError)
	selectImg := selectCmd.String("i", "", "cache md5")
	selectDb := selectCmd.String("o", "", "cache db path")

	showImgCmd := NewFlagSet("showimg", flag.ExitOnError)
	showImgImg := showImgCmd.String("f", "", "cache file path")

	saveCmd := NewFlagSet("export", flag.ExitOnError)
	saveImg := saveCmd.String("i", "", "cache md5")
	saveDb := saveCmd.String("o", "", "cache db path")
	saveDst := saveCmd.String("f", "", "cache file path")

	showAllCmd := NewFlagSet("showall", flag.ExitOnError)
	showAllDb := showAllCmd.String("o", "", "cache db path")

	dumpCmd := NewFlagSet("dump", flag.ExitOnError)
	dumpSrcDb := dumpCmd.String("s", "", "source cache db")
	dumpDstDb := dumpCmd.String("d", "", "destination cache db")
	flag.Parse()
	if *h {
		flag.Usage()
		fmt.Println("Subcommand Usage:")
		for _, set := range sets {
			fmt.Print(" ")
			set.Usage()
		}
		return
	}
	if len(os.Args) < 2 {
		fmt.Println("command invalid,please use -h")
		os.Exit(1)
	}

	var err error

	switch os.Args[1] {
	case "insert":
		insertCmd.Parse(os.Args[2:])
		err = base.InsertCacheToDb(*insertImg, *insertDb)
		if err == nil {
			fmt.Printf("Insert image %s to %s success\n", *insertImg, *insertDb)
		} else {
			fmt.Println(err)
		}
	case "select":
		selectCmd.Parse(os.Args[2:])
		res, err := base.ShowImageFromDb(*selectImg, *selectDb)
		if err == nil {
			fmt.Println(res)
		} else {
			fmt.Println(err)
		}
	case "showimage":
		showImgCmd.Parse(os.Args[2:])
		res, err := base.ShowImg(*showImgImg)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}
	case "export":
		saveCmd.Parse(os.Args[2:])
		res, err := base.SaveImageFromDb(*saveImg, *saveDb, *saveDst)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Save md5 %s from %s to %s success\n", *saveImg, *saveDb, res)
		}
	case "showall":
		showAllCmd.Parse(os.Args[2:])
		num, res := base.ShowAllFromDb(*showAllDb)
		if num == 0 {
			fmt.Printf("db %s is empty!\n", *showAllDb)
		} else {
			fmt.Println(res)
		}
	case "dump":
		dumpCmd.Parse(os.Args[2:])
		num := base.DumpAllToDb(*dumpSrcDb, *dumpDstDb)
		if num == 0 {
			fmt.Printf("db %s is empty!\n", *dumpSrcDb)
		} else {
			fmt.Printf("dump from %s to %s success,total %d images\n", *dumpSrcDb, *dumpDstDb, num)
		}
	default:
		fmt.Println("command invalid,please use -h")
		os.Exit(1)
	}
}
