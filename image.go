package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
)

type images []docker.APIImages

func (i images) Filter(f filter) images {
	ret := images{}
	for _, image := range i {
		if f(image) {
			ret = append(ret, image)
		}
	}
	return ret
}

func (i images) FilterByIncludeName(include []string) images {
	return i.Filter(func(object interface{}) bool {
		if len(include) == 0 {
			return true
		}

		image := object.(docker.APIImages)
		for i := range image.RepoTags {
			//include
			for _, name := range include {
				if strings.HasPrefix(image.RepoTags[i], name) {
					return true
				}
			}
		}
		return false
	})
}

func (i images) FilterByExclusiveName(exclusive []string) images {
	return i.Filter(func(object interface{}) bool {
		image := object.(docker.APIImages)
		for i := range image.RepoTags {
			//exclusive
			for _, name := range exclusive {
				if strings.HasPrefix(image.RepoTags[i], name) {
					return false
				}
			}
		}
		return true
	})
}

func (i images) FilterByCreatedAt(d time.Duration) images {
	return i.Filter(func(object interface{}) bool {
		image := object.(docker.APIImages)
		return time.Since(time.Unix(image.Created, 0)) > d
	})
}

func listImages(client *docker.Client) (images, error) {
	images := images{}
	apiImages, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return nil, err
	}
	for i := range apiImages {
		images = append(images, apiImages[i])
	}
	return apiImages, nil
}

func doImage(c *cli.Context) {
	duration, err := time.ParseDuration(c.String("duration"))
	if err != nil {
		log.Fatal(err)
	}

	client, err := docker.NewClient(c.GlobalString("endpoint"))
	if err != nil {
		log.Fatal(err)
	}

	images, err := listImages(client)
	if err != nil {
		log.Fatal(err)
	}

	include := c.StringSlice("include")
	exclusive := c.StringSlice("exclusive")

	ret := images.
		FilterByIncludeName(include).
		FilterByExclusiveName(exclusive).
		FilterByCreatedAt(duration)

	for i := range ret {
		var err error
		run(c.Bool("dryrun"),
			func() {
				fmt.Println("dryrun: removed:", ret[i].ID, ret[i].RepoTags)
			},
			func() {
				force := c.Bool("force")
				err = client.RemoveImageExtended(ret[i].ID, docker.RemoveImageOptions{Force: force})
				fmt.Println("removed:", ret[i].ID, ret[i].RepoTags)
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: failed to delete a image %v %s \n", err, ret[i].ID, ret[i].RepoTags)
		}
	}
}
