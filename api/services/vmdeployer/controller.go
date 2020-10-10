package vmdeployer

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/nomad/api"
	g "github.com/sdslabs/katana/configs"
)

func deployVMs(c *gin.Context) {

	conf := api.DefaultConfig()
	cli, _ := api.NewClient(conf)
	jobs := cli.Jobs()
	chImage := make(chan bool)
	chDisk := make(chan bool)
	go checkAndPullImage(chImage)
	go checkAndPullBootDisk(chDisk)

	count, err := strconv.Atoi(c.Param("replicas"))

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}

	job := getJob(count)

	if checkImage := <-chImage; checkImage == false {
		c.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   "Failed to deploy VMs",
		})
	}

	if checkBootDisk := <-chDisk; checkBootDisk == false {
		c.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   "Failed to deploy VMs",
		})
		return
	}

	jr, _, err := jobs.Register(&job, &api.WriteOptions{})
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   "Failed to deploy VMs",
		})
		return
	}

	eval := cli.Evaluations()
	evalu, _, err := eval.Info(jr.EvalID, &api.QueryOptions{})
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   "Failed to fetch evaluation information",
		})
		return
	}

	allocs, _, err := jobs.Allocations(evalu.JobID, true, &api.QueryOptions{})
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   "Failed to fetch Allocations",
		})
		return
	}

	allAllocations := cli.Allocations()
	ch := make(chan status)
	for _, elem := range allocs {
		go checkAllocation(ch, g.VMDeployerConfig.Taskname, elem.ID, *allAllocations)
	}

	var ips = make([]string, 256)
	var countvm int = 0
	for range allocs {
		resp := <-ch
		if resp.Error == nil {
			ips[countvm] = resp.Data.IP
			countvm++
		}
	}

	if count == 0 {
		c.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   "Failed to deploy VMs",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    ips,
	})

	return
}

func getJob(count int) api.Job {
	job := api.NewServiceJob(g.VMDeployerConfig.Jobname, g.VMDeployerConfig.Jobname, g.VMDeployerConfig.Region, g.VMDeployerConfig.Priority)
	job = job.AddDatacenter(g.VMDeployerConfig.Datacentre)
	taskGrp := getTaskGroup(count)
	job = job.AddTaskGroup(&taskGrp)
	return *job
}

func getTaskGroup(count int) api.TaskGroup {
	task := getTask()
	tg := &api.TaskGroup{
		Name:  &g.VMDeployerConfig.GroupName,
		Count: &count,
		RestartPolicy: &api.RestartPolicy{
			Attempts: &g.VMDeployerConfig.RestartAttempts,
			Mode:     &g.VMDeployerConfig.RestartMode,
		},

		Tasks: []*api.Task{&task},
	}

	return *tg
}

func getTask() api.Task {
	task := &api.Task{
		Name:   g.VMDeployerConfig.Taskname,
		Driver: g.VMDeployerConfig.Driver,
		Config: map[string]interface{}{
			"KernelImage": getKernelImagePath(),
			"Vcpus":       1,
			"Mem":         128,
			"BootDisk":    getBootDiskPath(),
			"Network":     g.VMDeployerConfig.NetworkName,
		},
		LogConfig: &api.LogConfig{
			MaxFiles:      &g.VMDeployerConfig.FileNum,
			MaxFileSizeMB: &g.VMDeployerConfig.FileSize,
		},
	}

	return *task
}
