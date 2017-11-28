package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	protoFile string
	outputDir string
)

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.PersistentFlags().StringVar(&protoFile, "proto", "", "protobuf file path")
	genCmd.PersistentFlags().StringVar(&outputDir, "output", "", "output dir,Default is current dir")
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "gen framework",
	Run:   gen,
}

func gen(cmd *cobra.Command, args []string) {
	if protoFile == "" {
		panic("No protobuf file path specified")
	}
	if outputDir == "" {
		outputDir = "."
	}
	outputDir = strings.TrimRight(outputDir, "/")
	//判断outputDir文件内是否不为空
	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		panic(err)
	}
	_, err = os.Stat(absOutputDir)
	if err != nil && os.IsNotExist(err) {
		//mkdir
		err := os.MkdirAll(absOutputDir, 0777)
		if err != nil {
			panic(fmt.Sprintf("mkdir %s failed,err=%v", absOutputDir, err))
		}
	}
	s, _ := ioutil.ReadDir(absOutputDir)
	if len(s) != 0 {
		panic(fmt.Sprintf("%s is not empty.", absOutputDir))
	}

	pbCmd := exec.Command("protoc", "-I.", protoFile, "--nest_out=output="+outputDir+":.")
	out, err := pbCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error: %s", string(out))
		panic(err)
	}
	//mkdir
	err = os.MkdirAll(outputDir+"/proto", 0777)
	if err != nil {
		panic(err)
	}
	pbCmd = exec.Command("protoc", "-I.", protoFile, "--go_out=plugins=grpc:"+outputDir+"/proto")
	out, err = pbCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error: %s", string(out))
		panic(err)
	}

}
