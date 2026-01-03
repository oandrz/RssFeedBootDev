package main

import (
	"bootDevGoRss/internal/config"
	"fmt"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	configData, err := config.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(configData)

	configData.SetUser("Oink")

	fmt.Println(configData)
}
