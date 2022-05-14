package main

import (
	"fmt"
	"os"
)

var plugin_list []string

func ReloadPluginList(dir string) {
	plugin_list = []string{}
	fileslist, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Failed to read directory: %v\n", err)
		os.Exit(1)
	}
	for _, f := range fileslist {
		content, err := os.ReadFile(dir + "\\" + f.Name() + "\\manifest.json")
		if err != nil {
			fmt.Printf("Failed to read file: %v\n", err)
		} else {
			plugin_list = append(plugin_list, "{ \"id\": \""+f.Name()+"\", "+string(content[1:len(content)-1])+"}")
		}
	}
}
