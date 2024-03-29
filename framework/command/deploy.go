package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/cobra"
	"github.com/26huitailang/yogo/framework/contract"
	"github.com/26huitailang/yogo/framework/provider/ssh"
	"github.com/26huitailang/yogo/framework/util"
	"github.com/google/shlex"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
)

// initDeployCommand 为自动化部署的命令
func initDeployCommand() *cobra.Command {
	deployCommand.AddCommand(deployFrontendCommand)
	deployCommand.AddCommand(deployBackendCommand)
	deployCommand.AddCommand(deployAllCommand)
	deployCommand.AddCommand(deployRollbackCommand)
	deployCommand.AddCommand(deployCleanCommand)
	return deployCommand
}

// deployCommand 一级命令，显示帮助信息
var deployCommand = &cobra.Command{
	Use:   "deploy",
	Short: "部署相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

// deployFrontendCommand 部署前端
var deployFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "部署前端",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()

		// 创建部署文件夹
		deployFolder, err := createDeployFolder(container)
		if err != nil {
			return err
		}

		// 编译前端到部署文件夹
		if err := deployBuildFrontend(c, deployFolder); err != nil {
			return err
		}

		// 上传部署文件夹并执行对应的shell
		return deployUploadAction(deployFolder, container, "frontend")
	},
}

// deployBackendCommand 部署后端
var deployBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "部署后端",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()

		// 创建部署文件夹
		deployFolder, err := createDeployFolder(container)
		if err != nil {
			return err
		}

		// 编译后端到部署文件夹
		if err := deployBuildBackend(c, deployFolder); err != nil {
			return err
		}

		// 上传部署文件夹并执行对应的shell
		return deployUploadAction(deployFolder, container, "backend")
	},
}

var deployAllCommand = &cobra.Command{
	Use:   "all",
	Short: "全部部署",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()

		deployFolder, err := createDeployFolder(container)
		if err != nil {
			return err
		}

		// 编译前端
		if err := deployBuildFrontend(c, deployFolder); err != nil {
			return err
		}

		// 编译后端
		if err := deployBuildBackend(c, deployFolder); err != nil {
			return err
		}

		// 上传前端+后端，并执行对应的shell
		return deployUploadAction(deployFolder, container, "all")
	},
}

// deployRollbackCommand 部署回滚
var deployRollbackCommand = &cobra.Command{
	Use:   "rollback",
	Short: "部署回滚",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()

		if len(args) != 2 {
			return errors.New("参数错误,请按照参数进行回滚 ./yogo deploy rollback [version] [frontend/backend/all]")
		}

		version := args[0]
		end := args[1]

		// 获取版本信息
		appService := container.MustMake(contract.AppKey).(contract.App)
		deployFolder := filepath.Join(appService.DeployFolder(), version)

		// 上传部署文件夹并执行对应的shell
		return deployUploadAction(deployFolder, container, end)
	},
}

// deployFrontendCommand 部署前端
var deployCleanCommand = &cobra.Command{
	Use:   "clean",
	Short: "清理deploy文件夹下的历史文件",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		deployFolder := appService.DeployFolder()
		if !util.Exists(deployFolder) {
			return errors.New(fmt.Sprintf("folder not exists: %s", deployFolder))
		}
		files, err := os.ReadDir(deployFolder)
		if err != nil {
			return err
		}
		for _, f := range files {
			if !f.IsDir() {
				continue
			}
			if err != nil {
				return err
			}
			path2Del := path.Join(deployFolder, f.Name())
			if err = os.RemoveAll(path2Del); err != nil {
				return err
			}
			fmt.Println("deleted deploy version:", f.Name())
		}
		return nil
	},
}

func deployBuildBackend(c *cobra.Command, deployFolder string) error {
	container := c.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	appService := container.MustMake(contract.AppKey).(contract.App)
	envService := container.MustMake(contract.EnvKey).(contract.Env)
	logger := container.MustMake(contract.LogKey).(contract.Log)

	env := envService.AppEnv()

	// 组装命令
	binFile := "yogo"
	if configService.GetString("deploy.backend.bin") != "" {
		binFile = configService.GetString("deploy.backend.bin")
	}
	useDocker := configService.GetBool("deploy.backend.use_docker")

	path := "/usr/local/go/bin/go"
	deployBinFile := filepath.Join(deployFolder, binFile)
	if !useDocker {
		lpath, err := exec.LookPath("go")
		path = lpath
		if err != nil {
			log.Fatalln("yogo go: 请在Path路径中先安装go")
		}
	} else {
		version := filepath.Base(deployFolder)
		deployFolderInDocker := filepath.Join("/code", "deploy", version)
		deployBinFile = filepath.Join(deployFolderInDocker, binFile)
	}

	cmd := exec.Command(path, "build", "-o", deployBinFile, "./")
	cmd.Env = os.Environ()
	// 设置GOOS和GOARCH
	if configService.GetString("deploy.backend.goos") != "" {
		cmd.Env = append(cmd.Env, "GOOS="+configService.GetString("deploy.backend.goos"))
	}
	if configService.GetString("deploy.backend.goarch") != "" {
		cmd.Env = append(cmd.Env, "GOARCH="+configService.GetString("deploy.backend.goarch"))
	}
	if configService.GetString("deploy.backend.cgo") != "" {
		cmd.Env = append(cmd.Env, "CGO_ENABLED="+configService.GetString("deploy.backend.cgo"))
	}
	if configService.GetString("deploy.backend.cc") != "" {
		cmd.Env = append(cmd.Env, "CC="+configService.GetString("deploy.backend.cc"))
	}

	if useDocker {
		image := configService.GetString("deploy.backend.image")
		dcmd := fmt.Sprintf("run --rm -v %s:/code %s bash -c", filepath.Dir(appService.BaseFolder()), image)
		dcmds, _ := shlex.Split(dcmd)
		dcmds = append(dcmds, fmt.Sprintf("ls /code && cd /code && %s", cmd.String()))
		dockerPath, _ := exec.LookPath("docker")
		cmd = exec.Command(dockerPath, dcmds...)
	}

	// 执行命令
	ctx := context.Background()
	logger.Info(ctx, "开始执行命令:", map[string]interface{}{"command": cmd.String()})
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(ctx, "go build err", map[string]interface{}{
			"err": err,
			"out": string(out),
		})
		return err
	}
	logger.Info(ctx, "编译成功", map[string]interface{}{"out": string(out)})

	// 复制.env文件
	if util.Exists(filepath.Join(appService.BaseFolder(), ".env")) {
		if err := util.CopyFile(filepath.Join(appService.BaseFolder(), ".env"), filepath.Join(deployFolder, ".env")); err != nil {
			return err
		}
	}

	// 复制config文件
	deployConfigFolder := filepath.Join(deployFolder, "config", env)
	if !util.Exists(deployConfigFolder) {
		if err := os.MkdirAll(deployConfigFolder, os.ModePerm); err != nil {
			return err
		}
	}
	if err := util.CopyFolder(filepath.Join(appService.ConfigFolder(), env), deployConfigFolder); err != nil {
		return err
	}

	logger.Info(ctx, "build local ok", nil)
	return nil
}

// 上传部署文件夹，并且执行对应的前置和后置的shell
func deployUploadAction(deployFolder string, container framework.Container, end string) error {
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	sshService := container.MustMake(contract.SSHKey).(contract.SSHService)
	logger := container.MustMake(contract.LogKey).(contract.Log)

	// 遍历所有deploy的服务器
	deployNodes := configService.GetStringSlice("deploy.connections")
	if len(deployNodes) == 0 {
		return errors.New("deploy connections len is zero")
	}
	remoteFolder := configService.GetString("deploy.remote_folder")
	if remoteFolder == "" {
		return errors.New("remote folder is empty")
	}

	preActions := make([]string, 0, 1)
	postActions := make([]string, 0, 1)

	if end == "frontend" || end == "all" {
		preActions = append(preActions, configService.GetStringSlice("deploy.frontend.pre_action")...)
		postActions = append(postActions, configService.GetStringSlice("deploy.frontend.post_action")...)
	}
	if end == "backend" || end == "all" {
		preActions = append(preActions, configService.GetStringSlice("deploy.backend.pre_action")...)
		postActions = append(postActions, configService.GetStringSlice("deploy.backend.post_action")...)
	}

	// 对每个远端服务
	for _, node := range deployNodes {
		sshClient, err := sshService.GetClient(ssh.WithConfigPath(node))
		if err != nil {
			return err
		}
		client, err := sftp.NewClient(sshClient)
		if err != nil {
			return err
		}

		// 执行所有的前置命令
		for _, action := range preActions {
			// 创建session
			session, err := sshClient.NewSession()
			if err != nil {
				return err
			}
			logger.Info(context.Background(), "execute pre action start", map[string]interface{}{
				"cmd":        action,
				"connection": node,
			})
			// 执行命令，并且等待返回
			bts, err := session.CombinedOutput(action)
			if err != nil {
				session.Close()
				return err
			}
			session.Close()
			// 执行前置命令成功
			logger.Info(context.Background(), "execute pre action", map[string]interface{}{
				"cmd":        action,
				"connection": node,
				"out":        strings.ReplaceAll(string(bts), "\n", ""),
			})
		}

		if err := uploadFolderToSFTP(container, deployFolder, remoteFolder, client); err != nil {
			logger.Info(context.Background(), "upload folder failed", map[string]interface{}{
				"remoteFolder": remoteFolder,
				"deployFolder": deployFolder,
				"err":          err,
			})
			return err
		}
		logger.Info(context.Background(), "upload folder success", nil)

		for _, action := range postActions {
			session, err := sshClient.NewSession()
			if err != nil {
				return err
			}
			logger.Info(context.Background(), "execute post action start", map[string]interface{}{
				"cmd":        action,
				"connection": node,
			})
			bts, err := session.CombinedOutput(action)
			if err != nil {
				session.Close()
				return err
			}
			logger.Info(context.Background(), "execute post action finish", map[string]interface{}{
				"cmd":        action,
				"connection": node,
				"out":        strings.ReplaceAll(string(bts), "\n", ""),
			})
			session.Close()
		}
	}
	return nil
}

// 上传部署文件夹
func uploadFolderToSFTP(container framework.Container, localFolder, remoteFolder string, client *sftp.Client) error {
	logger := container.MustMake(contract.LogKey).(contract.Log)
	// 遍历本地文件
	return filepath.Walk(localFolder, func(path string, info os.FileInfo, err error) error {
		// 获取除了folder前缀的后续文件名称
		relPath := strings.Replace(path, localFolder, "", 1)
		if relPath == "" {
			return nil
		}
		// 如果是遍历到了一个目录
		if info.IsDir() {
			logger.Info(context.Background(), "mkdir: "+filepath.Join(remoteFolder, relPath), nil)
			// 创建这个目录
			return client.MkdirAll(filepath.Join(remoteFolder, relPath))
		}

		// 打开本地的文件
		rf, err := os.Open(filepath.Join(localFolder, relPath))
		if err != nil {
			return errors.New("read file " + filepath.Join(localFolder, relPath) + " error:" + err.Error())
		}
		defer rf.Close()
		// 检查文件大小
		rfStat, err := rf.Stat()
		if err != nil {
			return err
		}
		// 打开/创建远端文件
		f, err := client.Create(filepath.Join(remoteFolder, relPath))
		if err != nil {
			return errors.New("create file " + filepath.Join(remoteFolder, relPath) + " error:" + err.Error())
		}
		defer f.Close()

		// 大于2M的文件显示进度
		if rfStat.Size() > 2*1024*1024 {
			logger.Info(context.Background(), "upload local file: "+filepath.Join(localFolder, relPath)+
				" to remote file: "+filepath.Join(remoteFolder, relPath)+" start", nil)
			// 开启一个goroutine来不断计算进度
			go func(localFile, remoteFile string) {
				// 每10s计算一次
				ticker := time.NewTicker(2 * time.Second)
				for range ticker.C {
					// 获取远端文件信息
					remoteFileInfo, err := client.Stat(remoteFile)
					if err != nil {
						logger.Error(context.Background(), "stat error", map[string]interface{}{
							"err":         err,
							"remote_file": remoteFile,
						})
						continue
					}
					// 如果远端文件大小等于本地文件大小，说明已经结束了
					size := remoteFileInfo.Size()
					if size >= rfStat.Size() {
						break
					}
					// 计算进度并且打印进度
					percent := int(size * 100 / rfStat.Size())
					logger.Info(context.Background(), "upload local file: "+filepath.Join(localFolder, relPath)+
						" to remote file: "+filepath.Join(remoteFolder, relPath)+fmt.Sprintf(" %v%% %v/%v", percent, size, rfStat.Size()), nil)
				}
			}(filepath.Join(localFolder, relPath), filepath.Join(remoteFolder, relPath))
		}

		// 将本地文件并发读取到远端文件
		if _, err := f.ReadFromWithConcurrency(rf, 10); err != nil {
			return errors.New("Write file " + filepath.Join(remoteFolder, relPath) + " error:" + err.Error())
		}
		// 记录成功信息
		logger.Info(context.Background(), "upload local file: "+filepath.Join(localFolder, relPath)+
			" to remote file: "+filepath.Join(remoteFolder, relPath)+" finish", nil)
		return nil
	})
}

func deployBuildFrontend(c *cobra.Command, deployFolder string) error {
	container := c.GetContainer()
	appService := container.MustMake(contract.AppKey).(contract.App)

	// 编译前端
	if err := buildFrontendCommand.RunE(c, []string{}); err != nil {
		return err
	}

	// 复制前端文件到deploy文件夹
	frontendFolder := filepath.Join(deployFolder, "dist")
	if err := os.Mkdir(frontendFolder, os.ModePerm); err != nil {
		return err
	}

	buildFolder := filepath.Join(appService.BaseFolder(), "dist")
	if err := util.CopyFolder(buildFolder, frontendFolder); err != nil {
		return err
	}
	return nil
}

// 创建部署的folder
func createDeployFolder(c framework.Container) (string, error) {
	appService := c.MustMake(contract.AppKey).(contract.App)
	deployFolder := appService.DeployFolder()

	deployVersion := time.Now().Format("20060102150405")
	versionFolder := filepath.Join(deployFolder, deployVersion)
	if !util.Exists(versionFolder) {
		return versionFolder, os.Mkdir(versionFolder, os.ModePerm)
	}
	return versionFolder, nil
}
