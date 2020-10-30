package main

import (
	"gitlab.com/contextualcode/platform_cc/cmd"
)

func main() {

	if err := cmd.Execute(); err != nil {
		panic(err)
	}

	/*projectPath := path.Join("api", "test_data", "sample2")
	p, e := api.LoadProjectFromPath(projectPath)
	if e != nil {
		panic(e)
	}
	log.Println(p.ID)*/

	//d, _ := newDockerClient()

	/*d, _ := api.NewDockerClient()
	err := d.StartProject(p)
	if err != nil {
		panic(err)
	}*/
	//d.ShellContainer(p.GetAppContainerName(p.Apps[0]))
	//d.PurgeProject(p)

	//d.CreateProjectNetwork(&p)
	//d.CreateAppContainer(&p, &p.GetAppDef()[0])

}
