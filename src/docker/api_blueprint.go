package docker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"bufio"
	"strconv"
	"os"
	"os/user"
	"errors"
	"github.com/docker/go-connections/nat"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	conttype "github.com/docker/docker/api/types/container"
	doctype "github.com/docker/docker/api/types"
	strslice "github.com/docker/docker/api/types/strslice"
	volumetype "github.com/docker/docker/api/types/volume"

	"github.com/azukaar/cosmos-server/src/utils"
)

type ContainerCreateRequestContainer struct {
	Name 			string            `json:"container_name"`
	Image       string            `json:"image"`
	Environment []string `json:"environment"`
	Labels      map[string]string `json:"labels"`
	Ports       []string          `json:"ports"`
	Volumes     []mount.Mount          `json:"volumes"`
	Networks    map[string]struct {
		Aliases []string `json:"aliases,omitempty"`
		IPV4Address string `json:"ipv4_address,omitempty"`
		IPV6Address string `json:"ipv6_address,omitempty"`
	} `json:"networks"`
	Routes 			   []utils.ProxyRouteConfig          `json:"routes"`

	RestartPolicy  string            `json:"restart,omitempty"`
	Devices        []string          `json:"devices"`
	Expose 		     []string          `json:"expose"`
	DependsOn      []string          `json:"depends_on"`
	Tty            bool              `json:"tty,omitempty"`
	StdinOpen      bool              `json:"stdin_open,omitempty"`

	Command string `json:"command,omitempty"`
	Entrypoint string `json:"entrypoint,omitempty"`
	WorkingDir string `json:"working_dir,omitempty"`
	User string `json:"user,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Domainname string `json:"domainname,omitempty"`
	MacAddress string `json:"mac_address,omitempty"`
	Privileged bool `json:"privileged,omitempty"`
	NetworkMode string `json:"network_mode,omitempty"`
	StopSignal string `json:"stop_signal,omitempty"`
	StopGracePeriod int `json:"stop_grace_period,omitempty"`
	HealthCheck struct {
		Test []string `json:"test"`
		Interval int `json:"interval"`
		Timeout int `json:"timeout"`
		Retries int `json:"retries"`
		StartPeriod int `json:"start_period"`
	} `json:"healthcheck,omitempty"`
	DNS []string `json:"dns,omitempty"`
	DNSSearch []string `json:"dns_search,omitempty"`
	ExtraHosts []string `json:"extra_hosts,omitempty"`
	Links []string `json:"links,omitempty"`
	SecurityOpt []string `json:"security_opt,omitempty"`
	StorageOpt map[string]string `json:"storage_opt,omitempty"`
	Sysctls map[string]string `json:"sysctls,omitempty"`
	Isolation string `json:"isolation,omitempty"`

	CapAdd []string `json:"cap_add,omitempty"`
	CapDrop []string `json:"cap_drop,omitempty"`
	SysctlsMap map[string]string `json:"sysctls,omitempty"`
}

type ContainerCreateRequestVolume struct {
	// name must be unique
	Name string `json:"name"`
	Driver string `json:"driver"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type ContainerCreateRequestNetwork struct {
	// name must be unique
	Name string `json:"name"`
	Driver string `json:"driver"`
	Attachable bool `json:"attachable"`
	Internal bool `json:"internal"`
	EnableIPv6 bool `json:"enable_ipv6"`
	IPAM struct {
		Driver string `json:"driver"`
		Config []struct {
			Subnet string `json:"subnet"`
		} `json:"config"`
	} `json:"ipam"`
}

type DockerServiceCreateRequest struct {
	Services map[string]ContainerCreateRequestContainer `json:"services"`
	Volumes []ContainerCreateRequestVolume `json:"volumes"`
	Networks map[string]ContainerCreateRequestNetwork `json:"networks"`
}

type DockerServiceCreateRollback struct {
	// action: disconnect, remove, etc...
	Action string `json:"action"`
	// type: container, volume, network
	Type string `json:"type"`
	// name: container name, volume name, network name
	Name string `json:"name"`
}

func Rollback(actions []DockerServiceCreateRollback , w http.ResponseWriter, flusher http.Flusher) {
	for i := len(actions) - 1; i >= 0; i-- {
		action := actions[i]
		switch action.Type {
		case "container":
			DockerClient.ContainerKill(DockerContext, action.Name, "SIGKILL")
			err := DockerClient.ContainerRemove(DockerContext, action.Name, doctype.ContainerRemoveOptions{})
			if err != nil {
				utils.Error("Rollback: Container", err)
			} else {
				utils.Log(fmt.Sprintf("Rolled back container %s", action.Name))
				fmt.Fprintf(w, "Rolled back container %s\n", action.Name)
				flusher.Flush()
			}
		case "volume":
			err := DockerClient.VolumeRemove(DockerContext, action.Name, true)
			if err != nil {
				utils.Error("Rollback: Volume", err)
			} else {
				utils.Log(fmt.Sprintf("Rolled back volume %s", action.Name))
				fmt.Fprintf(w, "Rolled back volume %s\n", action.Name)
				flusher.Flush()
			}
		case "network":
			err := DockerClient.NetworkRemove(DockerContext, action.Name)
			if err != nil {
				utils.Error("Rollback: Network", err)
			} else {
				utils.Log(fmt.Sprintf("Rolled back network %s", action.Name))
				fmt.Fprintf(w, "Rolled back network %s\n", action.Name)
				flusher.Flush()
			}
		}
	}
	
	// After all operations
	utils.Error("CreateService", fmt.Errorf("Operation failed. Changes have been rolled back."))
	fmt.Fprintf(w, "[OPERATION FAILED]. CHANGES HAVE BEEN ROLLEDBACK.\n")
	flusher.Flush()
}

func CreateServiceRoute(w http.ResponseWriter, req *http.Request) {
	if utils.AdminOnly(w, req) != nil {
		return
	}

	errD := Connect()
	if errD != nil {
		utils.Error("CreateService", errD)
		utils.HTTPError(w, "Internal server error: " + errD.Error(), http.StatusInternalServerError, "DS002")
		return
	}

	if req.Method == "POST" {
		// Enable streaming of response by setting appropriate headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Transfer-Encoding", "chunked")

		flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
        return
    }
		
		utils.ConfigLock.Lock()
		defer utils.ConfigLock.Unlock()
		
		utils.Log("Starting creation of new service...")
		fmt.Fprintf(w, "Starting creation of new service...\n")
		flusher.Flush()
		
		config := utils.ReadConfigFromFile()
		configRoutes := config.HTTPConfig.ProxyConfig.Routes
				
		decoder := json.NewDecoder(req.Body)
		var serviceRequest DockerServiceCreateRequest
		err := decoder.Decode(&serviceRequest)
		if err != nil {
			utils.Error("CreateService", err)
			utils.HTTPError(w, "Bad request: "+err.Error(), http.StatusBadRequest, "DS003")
			return
		}

		var rollbackActions []DockerServiceCreateRollback

		// Create networks
		for networkToCreateName, networkToCreate := range serviceRequest.Networks {
			utils.Log(fmt.Sprintf("Creating network %s...", networkToCreateName))
			fmt.Fprintf(w, "Creating network %s...\n", networkToCreateName)
			flusher.Flush()

			// check if network already exists
			_, err = DockerClient.NetworkInspect(DockerContext, networkToCreateName, doctype.NetworkInspectOptions{})
			if err == nil {
				utils.Error("CreateService: Network", err)
				fmt.Fprintf(w, "[ERROR] Network %s already exists\n", networkToCreateName)
				flusher.Flush()
				Rollback(rollbackActions, w, flusher)
				return
			}

			ipamConfig := make([]network.IPAMConfig, len(networkToCreate.IPAM.Config))
			if networkToCreate.IPAM.Config != nil {
				for i, config := range networkToCreate.IPAM.Config {
						ipamConfig[i] = network.IPAMConfig{
								Subnet: config.Subnet,
						}
				}
			}
			
			_, err = DockerClient.NetworkCreate(DockerContext, networkToCreateName, doctype.NetworkCreate{
					Driver:     networkToCreate.Driver,
					Attachable: networkToCreate.Attachable,
					Internal:   networkToCreate.Internal,
					EnableIPv6: networkToCreate.EnableIPv6,
					IPAM: &network.IPAM{
							Driver: networkToCreate.IPAM.Driver,
							Config: ipamConfig,
					},
			})

			if err != nil {
				utils.Error("CreateService: Rolling back changes because of -- Network", err)
				fmt.Fprintf(w, "[ERROR] Rolling back changes because of -- Network creation error: "+err.Error())
				flusher.Flush()
				Rollback(rollbackActions, w, flusher)
				return
			}

			rollbackActions = append(rollbackActions, DockerServiceCreateRollback{
				Action: "remove",
				Type:   "network",
				Name:   networkToCreateName,
			})
		
			// Write a response to the client
			utils.Log(fmt.Sprintf("Network %s created", networkToCreateName))
			fmt.Fprintf(w, "Network %s created\n", networkToCreateName)
			flusher.Flush()
		}

		// Create volumes
		for _, volume := range serviceRequest.Volumes {
			utils.Log(fmt.Sprintf("Creating volume %s...", volume.Name))
			fmt.Fprintf(w, "Creating volume %s...\n", volume.Name)
			flusher.Flush()
			
			_, err = DockerClient.VolumeCreate(DockerContext, volumetype.CreateOptions{
				Driver:     volume.Driver,
				Name:       volume.Name,
			})

			if err != nil {
				utils.Error("CreateService: Rolling back changes because of -- Volume", err)
				fmt.Fprintf(w, "[ERROR] Rolling back changes because of -- Volume creation error: "+err.Error())
				flusher.Flush()
				Rollback(rollbackActions, w, flusher)
				return
			}

			rollbackActions = append(rollbackActions, DockerServiceCreateRollback{
				Action: "remove",
				Type:   "volume",
				Name:   volume.Name,
			})

			// Write a response to the client
			utils.Log(fmt.Sprintf("Volume %s created", volume.Name))
			fmt.Fprintf(w, "Volume %s created\n", volume.Name)
			flusher.Flush()
		}

		// pull images
		for _, container := range serviceRequest.Services {
			// Write a response to the client
			utils.Log(fmt.Sprintf("Pulling image %s", container.Image))
			fmt.Fprintf(w, "Pulling image %s\n", container.Image)
			flusher.Flush()

			out, err := DockerClient.ImagePull(DockerContext, container.Image, doctype.ImagePullOptions{})
			if err != nil {
				utils.Error("CreateService: Rolling back changes because of -- Image pull", err)
				fmt.Fprintf(w, "[ERROR] Rolling back changes because of -- Image pull error: "+err.Error())
				flusher.Flush()
				Rollback(rollbackActions, w, flusher)
				return
			}
			defer out.Close()

			// wait for image pull to finish
			scanner := bufio.NewScanner(out)
			for scanner.Scan() {
				fmt.Fprintf(w, "%s\n", scanner.Text())
				flusher.Flush()
			}
			
			// Write a response to the client
			utils.Log(fmt.Sprintf("Image %s pulled", container.Image))
			fmt.Fprintf(w, "Image %s pulled\n", container.Image)
			flusher.Flush()
		}

		// Create containers
		for _, container := range serviceRequest.Services {
			utils.Log(fmt.Sprintf("Creating container %s...", container.Name))
			fmt.Fprintf(w, "Creating container %s...\n", container.Name)
			flusher.Flush()

			containerConfig := &conttype.Config{
				Image:        container.Image,
				Env:          container.Environment,
				Labels:       container.Labels,
				ExposedPorts: nat.PortSet{},
				WorkingDir:   container.WorkingDir,
				User:         container.User,
				Hostname:     container.Hostname,
				Domainname:   container.Domainname,
				MacAddress:   container.MacAddress,
				StopSignal:   container.StopSignal,
				StopTimeout:  &container.StopGracePeriod,
				Tty:          container.Tty,
				OpenStdin:    container.StdinOpen,
			}
			
			if container.Command != "" {
				containerConfig.Cmd = strings.Fields(container.Command)
			}

			if container.Entrypoint != "" {
				containerConfig.Entrypoint = strslice.StrSlice(strings.Fields(container.Entrypoint))
			}

			// For Expose / Ports
			for _, expose := range container.Expose {
				exposePort := nat.Port(expose)
				containerConfig.ExposedPorts[exposePort] = struct{}{}
			}

			PortBindings := nat.PortMap{}

			for _, port := range container.Ports {
				portContainer := strings.Split(port, ":")[0]
				portHost := strings.Split(port, ":")[1]			

				containerConfig.ExposedPorts[nat.Port(portContainer)] = struct{}{}
				PortBindings[nat.Port(portContainer)] = []nat.PortBinding{
					{
						HostIP:   "",
						HostPort: portHost,
					},
				}
			}

			// Create missing folders for bind mounts
			for _, newmount := range container.Volumes {
				if newmount.Type == mount.TypeBind {
					if _, err := os.Stat(newmount.Source); os.IsNotExist(err) {
						err := os.MkdirAll(newmount.Source, 0755)
						if err != nil {
							utils.Error("CreateService: Unable to create directory for bind mount", err)
							fmt.Fprintf(w, "[ERROR] Unable to create directory for bind mount: "+err.Error())
							flusher.Flush()
							Rollback(rollbackActions, w, flusher)
							return
						}
			
						if container.User != "" {
							// Change the ownership of the directory to the container.User
							userInfo, err := user.Lookup(container.User)
							if err != nil {
								utils.Error("CreateService: Unable to lookup user", err)
								fmt.Fprintf(w, "[ERROR] Unable to lookup user " + container.User + "." +err.Error())
								flusher.Flush()
							} else {
								uid, _ := strconv.Atoi(userInfo.Uid)
								gid, _ := strconv.Atoi(userInfo.Gid)
								err = os.Chown(newmount.Source, uid, gid)
								if err != nil {
									utils.Error("CreateService: Unable to change ownership of directory", err)
									fmt.Fprintf(w, "[ERROR] Unable to change ownership of directory: "+err.Error())
									flusher.Flush()
								}
							}	
						}
					}
				}
			}

			hostConfig := &conttype.HostConfig{
				PortBindings: PortBindings,
				Mounts:       container.Volumes,
				RestartPolicy: conttype.RestartPolicy{
					Name: container.RestartPolicy,
				},
				Privileged:   container.Privileged,
				NetworkMode:  conttype.NetworkMode(container.NetworkMode),
				DNS:         container.DNS,
				DNSSearch:   container.DNSSearch,
				ExtraHosts:  container.ExtraHosts,
				Links:       container.Links,
				SecurityOpt: container.SecurityOpt,
				StorageOpt:  container.StorageOpt,
				Sysctls:     container.Sysctls,
				Isolation:   conttype.Isolation(container.Isolation),
				CapAdd:      container.CapAdd,
				CapDrop:     container.CapDrop,
			}

			// For Healthcheck
			if len(container.HealthCheck.Test) > 0 {
				containerConfig.Healthcheck = &conttype.HealthConfig{
					Test: container.HealthCheck.Test,
					Interval: time.Duration(container.HealthCheck.Interval) * time.Second,
					Timeout: time.Duration(container.HealthCheck.Timeout) * time.Second,
					StartPeriod: time.Duration(container.HealthCheck.StartPeriod) * time.Second,
					Retries: container.HealthCheck.Retries,
				}
			}

			// For Devices
			devices := []conttype.DeviceMapping{}

			for _, device := range container.Devices {
				deviceSplit := strings.Split(device, ":")
				devices = append(devices, conttype.DeviceMapping{
					PathOnHost:        deviceSplit[0],
					PathInContainer:   deviceSplit[1],
					CgroupPermissions: "rwm", // This can be "r", "w", "m", or any combination
				})
			}

			hostConfig.Devices = devices


			networkingConfig := &network.NetworkingConfig{
				EndpointsConfig: make(map[string]*network.EndpointSettings),
			}
			
			for netName, netConfig := range container.Networks {
					networkingConfig.EndpointsConfig[netName] = &network.EndpointSettings{
							Aliases:     netConfig.Aliases,
							IPAddress:   netConfig.IPV4Address,
							GlobalIPv6Address: netConfig.IPV6Address,
					}
			}
	
			_, err = DockerClient.ContainerCreate(DockerContext, containerConfig, hostConfig, networkingConfig, nil, container.Name)
			
			if err != nil {
				utils.Error("CreateService: Rolling back changes because of -- Container", err)
				fmt.Fprintf(w, "[ERROR] Rolling back changes because of -- Container creation error: "+err.Error())
				flusher.Flush()
				Rollback(rollbackActions, w, flusher)
				return
			}
			
			rollbackActions = append(rollbackActions, DockerServiceCreateRollback{
				Action: "remove",
				Type:   "container",
				Name:   container.Name,
			})

			// add routes 
			for _, route := range container.Routes {
				// check if route already exists
				exists := false
				for _, configRoute := range configRoutes {
					if configRoute.Name == route.Name {
						exists = true
						break
					}
				}
				if !exists {
					configRoutes = append([]utils.ProxyRouteConfig{(utils.ProxyRouteConfig)(route)}, configRoutes...)
				} else {
					utils.Error("CreateService: Rolling back changes because of -- Route already exist", nil)
					fmt.Fprintf(w, "[ERROR] Rolling back changes because of -- Route already exist")
					flusher.Flush()
					Rollback(rollbackActions, w, flusher)
					return
				}
			}
			
			// Write a response to the client
			utils.Log(fmt.Sprintf("Container %s created", container.Name))
			fmt.Fprintf(w, "Container %s created\n", container.Name)
			flusher.Flush()
		}

		// re-order containers dpeneding on depends_on
		startOrder, err := ReOrderServices(serviceRequest.Services)
		if err != nil {
			utils.Error("CreateService: Rolling back changes because of -- Container", err)
			fmt.Fprintf(w, "[ERROR] Rolling back changes because of -- Container creation error: "+err.Error())
			flusher.Flush()
			Rollback(rollbackActions, w, flusher)
			return
		}

		// Start all the newly created containers
		for _, container := range startOrder {
			err = DockerClient.ContainerStart(DockerContext, container.Name, doctype.ContainerStartOptions{})
			if err != nil {
				utils.Error("CreateService: Start Container", err)
				fmt.Fprintf(w, "[ERROR] Rolling back changes because of -- Container start error: "+err.Error())
				flusher.Flush()
				Rollback(rollbackActions, w, flusher)
				return
			}

			// Write a response to the client
			utils.Log(fmt.Sprintf("Container %s started", container.Name))
			fmt.Fprintf(w, "Container %s started\n", container.Name)
			flusher.Flush()
		}
		
		// Save the route configs 
		config.HTTPConfig.ProxyConfig.Routes = configRoutes
		utils.SaveConfigTofile(config)
		utils.NeedsRestart = true

		// After all operations
		utils.Log("CreateService: Operation succeeded. SERVICE STARTED")
		fmt.Fprintf(w, "[OPERATION SUCCEEDED]. SERVICE STARTED\n")
		flusher.Flush()
	} else {
		utils.Error("CreateService: Method not allowed" + req.Method, nil)
		utils.HTTPError(w, "Method not allowed", http.StatusMethodNotAllowed, "HTTP001")
		return
	}
}

func ReOrderServices(serviceMap map[string]ContainerCreateRequestContainer) ([]ContainerCreateRequestContainer, error) {
	startOrder := []ContainerCreateRequestContainer{}

	for len(serviceMap) > 0 {
		// Keep track of whether we've added any services in this iteration
		changed := false

		for name, service := range serviceMap {
			// Check if all dependencies are already in startOrder
			allDependenciesStarted := true
			for _, dependency := range service.DependsOn {
				dependencyStarted := false
				for _, startedService := range startOrder {
					if startedService.Name == dependency {
						dependencyStarted = true
						break
					}
				}
				if !dependencyStarted {
					allDependenciesStarted = false
					break
				}
			}

			// If all dependencies are started, we can add this service to startOrder
			if allDependenciesStarted {
				startOrder = append(startOrder, service)
				delete(serviceMap, name)
				changed = true
			}
		}

		// If we haven't added any services in this iteration, then there must be a circular dependency
		if !changed {
			break
		}
	}

	// If there are any services left in serviceMap, they couldn't be started due to unsatisfied dependencies or circular dependencies
	if len(serviceMap) > 0 {
		errorMessage := "Could not start all services due to unsatisfied dependencies or circular dependencies:\n"
		for name, _ := range serviceMap {
			errorMessage += "Could not start service: " + name + "\n"
			errorMessage += "Unsatisfied dependencies:\n"
			for _, dependency := range serviceMap[name].DependsOn {
				_, ok := serviceMap[dependency]
				if ok {
					errorMessage += dependency + "\n"
				}
			}
		}
		return nil, errors.New(errorMessage)
	}

	return startOrder, nil
}