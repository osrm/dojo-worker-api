package api

import (
	"github.com/gin-gonic/gin"
)

func LoginRoutes(router *gin.Engine) {
	apiV1 := router.Group("/api/v1")
	{
		worker := apiV1.Group("/worker")
		{
			worker.POST("/login/auth", WorkerLoginMiddleware(), WorkerLoginController)
			worker.POST("/partner", WorkerAuthMiddleware(), WorkerPartnerCreateController)
			worker.POST("/partner/disable", WorkerAuthMiddleware(), DisableMinerByWorkerController)
		}
        apiV1.PUT("/partner/edit", WorkerAuthMiddleware(), UpdateWorkerPartnerController)
		tasks := apiV1.Group("/tasks")
		{
			tasks.POST("/create-task", MinerAuthMiddleware(), CreateTaskController)
			tasks.PUT("/submit-result/:task-id", WorkerAuthMiddleware(), SubmitTaskResultController)
			tasks.GET("/:task-id", GetTaskByIdController)
			tasks.GET("/", GetTasksByPageController)
			tasks.GET("/get-results/:task-id", GetTaskResultsController)
		}

		miner := apiV1.Group("/miner")
		{
			miner.POST("/login/auth", MinerLoginMiddleware(), MinerLoginController)
			miner.GET("/info/:hotkey",MinerAuthMiddleware(), MinerInfoController)
			worker.POST("/partner/disable", MinerAuthMiddleware(), DisableWorkerByMinerController)
		}
	}
}
