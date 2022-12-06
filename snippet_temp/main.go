package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
)

var (
	branch_commit = "develop"
	url_git       = "https://ee.ru"
	project       = "t/b"
	path_clone    = "D:/../tempRepo_"
)

func main() {
	//sys_clone_repo()
	gitStatus()
}

func gitStatus() {

	f, err := os.Open("snippet_temp/testD.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	wr := bytes.Buffer{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		wr.WriteString(sc.Text())
	}

	//tst_p(wr.String())
	r := bytes.NewReader(wr.Bytes())
	tst_p(r)

	// str_command := "git status"
	// system_exec(str_command)

}

func tst_p(st io.Reader) {

	// data, err := base64.StdEncoding.DecodeString(st)
	// if err != nil {
	// 	fmt.Println("-1-1-1-1-1-1-1-11-")
	// 	fmt.Println(err.Error())
	// 	return
	// }

	dec := base64.NewDecoder(base64.StdEncoding, st)
	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, dec)

	f, _ := os.Create("snippet_temp/data.epf")
	defer f.Close()

	// w := bufio.NewWriter(data)
	_, err = f.Write(buf.Bytes())
	//fmt.Fprintln(w, st)
	if err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
}

func sys_clone_repo() {
	//если нет репозитория - клонируем, переходим наветку
	//иначе перейти на ветку и pull из origin

	fullCommand := fmt.Sprintf("git clone %s/%s.git %s", url_git, project, path_clone)

	system_exec(fullCommand)
}

func system_exec(str_cmd string) error {
	//run command system
	c := exec.Command("cmd", "/C", str_cmd)

	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
		return err
	} else {
		fmt.Println("Done!")
		return nil
	}

}

func checkedOscript() error {
	return nil
}

func oscriptVersion() string {
	//oscript -v
	return ""

}
func checkedGit() error {
	return nil
}
func checkedPrecommit() error {
	//oscript -encoding=utf-8 "\_deploy\pre-commit\src\main.os" precommit ./ -source-dir "DataProcessorsExt"
	return nil
}
