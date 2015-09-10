package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
)

type containers []*docker.Container

func (c containers) Filter(f filter) containers {
	ret := containers{}
	for _, container := range c {
		if f(container) {
			ret = append(ret, container)
		}
	}
	return ret
}

func (c containers) FilterByFinishedAt(d time.Duration) containers {
	return c.Filter(func(object interface{}) bool {
		container := object.(*docker.Container)
		return time.Now().Sub(container.State.FinishedAt) > d
	})
}

func listContainers(client *docker.Client) (containers, error) {
	containers := containers{}
	apiContainers, err := client.ListContainers(docker.ListContainersOptions{
		All:     true,
		Filters: map[string][]string{"status": {"exited"}},
	})
	if err != nil {
		return nil, err
	}

	for _, apiContainer := range apiContainers {
		container, err := client.InspectContainer(apiContainer.ID)
		if err == nil {
			containers = append(containers, container)
		}
	}

	return containers, nil
}

func doContainer(c *cli.Context) {
	duration, err := time.ParseDuration(c.String("duration"))
	if err != nil {
		log.Fatal(err)
	}

	client, err := docker.NewClient(c.GlobalString("endpoint"))
	if err != nil {
		log.Fatal(err)
	}

	containers, err := listContainers(client)
	if err != nil {
		log.Fatal(err)
	}

	ret := containers.FilterByFinishedAt(duration)

	for i := range ret {
		var err error
		run(c.Bool("dryrun"),
			func() {
				fmt.Println("dryrun: removed:", ret[i].ID, ret[i].Image, ret[i].Name)
			},
			func() {
				force := c.Bool("force")
				client.RemoveContainer(docker.RemoveContainerOptions{ID: ret[i].ID, Force: force})
				fmt.Println("removed:", ret[i].ID, ret[i].Image, ret[i].Name)
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: failed to delete a contaienr %v %s \n", err, ret[i].ID, ret[i].Image, ret[i].Name)
		}
	}
}
