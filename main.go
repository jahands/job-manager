package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

const JobPrefix = "job:"

func jobKey(jobId string) string {
	return JobPrefix + jobId
}
func main() {
	// Setup redis
	// ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_URL", "localhost:6379"),
		DB:   0, // use default DB
	})

	// Setup gin
	r := gin.Default()

	// Add auth
	r.Use(func(c *gin.Context) {
		if c.Request.URL.Query().Get("api_key") != getEnv("API_KEY", "x") {
			c.AbortWithError(401, fmt.Errorf("unauthorized"))
		}
	})

	// Create router
	v1 := r.Group("/v1")
	{
		// Add new job (or replace existing)
		v1.PUT("/jobs/:jobId", func(c *gin.Context) {
			jobId := c.Param("jobId")
			job := Job{
				JobKey:  jobId,
				Created: time.Now(),
				InUse:   false,
			}

			err := rdb.Set(jobKey(jobId), job, 0).Err()
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "error creating job"})
				return
			}
			c.JSON(200, gin.H{"result": "success"})
		})

		// Update job as still in use
		v1.POST("/jobs/:jobId", func(c *gin.Context) {
			jobId := c.Param("jobId")
			job := Job{}
			jsonStr, err := rdb.Get(jobKey(jobId)).Bytes()
			if err != nil {
				if err == redis.Nil {
					c.AbortWithStatusJSON(404, gin.H{"error": "job not found"})
					return
				}
				c.AbortWithStatusJSON(500, gin.H{"error": "error getting job"})
				return
			}
			if err = job.UnmarshalBinary(jsonStr); err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "error unmarshaling job"})
				return
			}
			job.InUse = true
			job.LastInUse = time.Now()
			err = rdb.Set(jobKey(jobId), job, 0).Err()
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "error updating job"})
				return
			}
			c.JSON(200, gin.H{"result": "success"})
		})

		// Get all jobs
		v1.GET("/jobs/", func(c *gin.Context) {
			jobs, err := getAllJobs(rdb)
			if err != nil {
				if _, ok := err.(*NotFoundError); ok {
					c.AbortWithStatusJSON(404, gin.H{"error": "no jobs found"})
					return
				}
				c.AbortWithStatusJSON(500, gin.H{"error": err})
				return
			} else if len(jobs) == 0 {
				c.AbortWithStatusJSON(404, gin.H{"error": "no jobs found"})
			}
			c.JSON(200, gin.H{"result": jobs})
		})

		// Get unused job and mark as in use
		// optional ?min_age=<minutes> parameter to specify minimum age if in use
		// If an in-use job is older than this, it will be assumed as not in use.
		v1.GET("/job", func(c *gin.Context) {
			// Find an unused job, lock it, and return it
			jobs, err := getAllJobs(rdb)
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": err})
				return
			}
			now := time.Now()
			minAgeStr := c.Request.URL.Query().Get("min_age")
			var minAge int
			if minAgeStr != "" {
				minAge, err = strconv.Atoi(minAgeStr)
				if err != nil {
					c.AbortWithStatusJSON(400, gin.H{"error": "error parsing min_age"})
				}
			}
			for _, job := range jobs {
				// Send back jobs that are not in use or haven't been last_used_on for more than min_age
				if !job.InUse || (minAge > 0 && now.Sub(job.LastInUse) > time.Minute*time.Duration(minAge)) {
					job.InUse = true
					job.LastInUse = time.Now()
					err := rdb.Set(jobKey(job.JobKey), job, 0).Err()
					if err != nil {
						c.AbortWithStatusJSON(500, gin.H{"error": err})
						return
					}
					c.JSON(200, gin.H{"result": job})
					return
				}
			}
			c.AbortWithStatusJSON(404, gin.H{"error": "no jobs available"})
		})

		// Get a job by id
		v1.GET("/jobs/:jobId", func(c *gin.Context) {
			jobId := c.Param("jobId")
			job := Job{}
			jsonStr, err := rdb.Get(JobPrefix + jobId).Bytes()
			if err != nil {
				if err == redis.Nil {
					c.AbortWithStatusJSON(404, gin.H{"error": "job not found"})
					return
				}
				c.AbortWithStatusJSON(500, gin.H{"error": "error getting job"})
				return
			}
			if err = job.UnmarshalBinary(jsonStr); err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "error unmarshaling job"})
				return
			}
			c.JSON(200, gin.H{"result": job})
		})
	}

	r.Run(":8080")
}

type Job struct {
	JobKey    string    `json:"job_key"`
	InUse     bool      `json:"in_use"`
	Created   time.Time `json:"created_on"`
	LastInUse time.Time `json:"last_used_on"`
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

func getAllJobs(rdb *redis.Client) ([]Job, error) {
	keys, err := rdb.Keys(JobPrefix + "*").Result()
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, &NotFoundError{}
	}
	jobs := make([]Job, len(keys))
	for i, key := range keys {
		jobStr, err := rdb.Get(key).Bytes()
		if err != nil {
			return nil, err
		}
		if err = jobs[i].UnmarshalBinary(jobStr); err != nil {
			return nil, err
		}
	}
	return jobs, nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
