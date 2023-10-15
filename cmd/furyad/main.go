package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	furya "github.com/furya-labs/furya/v20/app"
	"github.com/furya-labs/furya/v20/app/params"
	"github.com/furya-labs/furya/v20/cmd/furyad/cmd"
)

func main() {
	params.SetAddressPrefixes()
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, furya.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
