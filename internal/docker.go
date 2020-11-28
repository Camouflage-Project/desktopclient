package internal

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
)

type ExecResult struct {
	StdOut string
	StdErr string
	ExitCode int
}

func SpinUpContainer() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}

	reader, err := cli.ImagePull(ctx, "docker.io/nvidia/cuda", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	containerCreateResp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "nvidia/cuda",
		Cmd:   nil,
		Tty:   true,
		ExposedPorts: nat.PortSet{
			"22/tcp": struct{}{},
		},
	},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"22/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "22",
					},
				},
			},
		}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, containerCreateResp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	command := []string{"bash", "-c", "ln -snf /usr/share/zoneinfo/Europe/Kiev /etc/localtime && echo Europe/Kiev > /etc/timezone && apt-get update && apt-get install -y openssh-server && mkdir /var/run/sshd && echo 'root:THEPASSWORDYOUCREATED' | chpasswd && sed -i 's/#*PermitRootLogin prohibit-password/PermitRootLogin yes/g' /etc/ssh/sshd_config && sed -i 's@session\\s*required\\s*pam_loginuid.so@session optional pam_loginuid.so@g' /etc/pam.d/sshd && service ssh start"}

	config :=  types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd: command,
	}

	execResp, err := cli.ContainerExecCreate(ctx, containerCreateResp.ID, config)
	if err != nil {
		panic(err)
	}

	_, err = cli.ContainerExecAttach(context.Background(), execResp.ID, types.ExecConfig{})
	if err != nil {
		fmt.Println(err)
	}

	//stopContainer(ctx, cli, containerCreateResp.ID)
	//removeContainer(ctx, cli, containerCreateResp.ID)
}

func stopContainer(ctx context.Context, cli *client.Client, id string) {
	if err := cli.ContainerStop(ctx, id, nil); err != nil {
		panic(err)
	}
}

func removeContainer(ctx context.Context, cli *client.Client, id string) {
	if err := cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}
}