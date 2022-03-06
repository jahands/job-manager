package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	docs "github.com/jahands/job-manager/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Job Manager API
func main() {
	// Setup gin
	r := gin.Default()
	// Add auth
	r.Use(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
			if c.Request.URL.Query().Get("api_key") != getEnv("API_KEY", "x") {
				c.AbortWithError(401, fmt.Errorf("unauthorized"))
				return
			}
			c.Header("cache-control", "no-cache")
		}
	})
	// Create router
	v1 := r.Group("/v1")
	{
		// Add new job
		v1.PUT("/:namespace/jobs/:jobId", PutJobById)
		// Update job as still in use
		v1.POST("/:namespace/jobs/:jobId", PostJobStillInUse)
		// Get all jobs
		v1.GET("/:namespace/jobs/", GetAllJobs)
		// Get unused job and mark as in use
		// optional ?min_age=<minutes> parameter to specify minimum age if in use
		// If an in-use job is older than this, it will be assumed as not in use.
		v1.GET("/:namespace/job", GetUnusedJob)
		// Get a job by id
		v1.GET("/:namespace/jobs/:jobId", GetJobById)
		// Delete a job by id
		v1.DELETE("/:namespace/jobs/:jobId", DeleteJobById)
	}
	// Add docs
	docs.SwaggerInfo.BasePath = "/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}

type JobsResponse struct {
	Result []Job `json:"result"`
}
type JobResponse struct {
	Result Job `json:"result"`
}
type SuccessResponse struct {
	Result string `json:"result"`
}
type ResultResponse struct {
	Result interface{} `json:"result"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
type Job struct {
	JobKey    string      `json:"job_key"`
	InUse     bool        `json:"in_use"`
	Created   time.Time   `json:"created_on"`
	LastInUse time.Time   `json:"last_used_on"`
	InUseBy   string      `json:"in_use_by,omitempty"`
	Meta      interface{} `json:"meta,omitempty"`
}

func (i Job) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (i *Job) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	return nil
}

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "not found"
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
