package wordpress

import (
	"fmt"

	"github.com/ChrisWiegman/kana/internal/docker"
)

func NewWordPress(siteName string, controller *docker.Controller) {

	wordPressContainers := []docker.ContainerConfig{
		{
			Image:       "mariadb",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_mariadb", siteName),
			Env: []string{
				"MARIADB_ROOT_PASSWORD=password",
				"MARIADB_DATABASE=wordpress",
				"MARIADB_USER=wordpress",
				"MARIADB_PASSWORD=wordpress",
			},
		},
		{
			Image:       "wordpress",
			NetworkName: "kana",
			HostName:    fmt.Sprintf("kana_%s_wordpress", siteName),
			Env: []string{
				fmt.Sprintf("WORDPRESS_DB_HOST=kana_%s_mariadb", siteName),
				"WORDPRESS_DB_USER=wordpress",
				"WORDPRESS_DB_PASSWORD=wordpress",
				"WORDPRESS_DB_NAME=wordpress",
			},
			Labels: map[string]string{
				"traefik.enable": "true",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.entrypoints", siteName): "web",
				fmt.Sprintf("traefik.http.routers.wordpress-%s-http.rule", siteName):        fmt.Sprintf("Host(`%s.sites.cfw.li`)", siteName),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.entrypoints", siteName):      "websecure",
				fmt.Sprintf("traefik.http.routers.wordpress-%s.rule", siteName):             fmt.Sprintf("Host(`%s.sites.cfw.li`)", siteName),
				fmt.Sprintf("traefik.http.routers.wordpress-%s.tls", siteName):              "true",
			},
		},
	}

	for _, container := range wordPressContainers {
		err := controller.EnsureImage(container.Image)
		if err != nil {
			panic(err)
		}

		_, err = controller.ContainerRun(container)
		if err != nil {
			panic(err)
		}
	}
}
