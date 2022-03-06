package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

const JobPrefix = "job:"

func jobKey(jobId string, namespace string) string {
	if namespace == "" {
		namespace = "default"
	}
	return JobPrefix + namespace + ":" + jobId
}

type PutJobByIdBody struct {
	Meta interface{} `json:"meta,omitempty"`
}

// PutJobById godoc
// @Summary Insert/replace job with specified id
// @Description Overwrites job if it exists.
// @Tags root
// @Param namespace path string true "Namespace of job(s)"
// @Param jobId path string true "Job ID"
// @Param api_key query string true "API Key"
// @Param in_use_by query string false "Who's using this job (optional) eg. hostname of machine using it."
// @Param body body PutJobByIdBody false "Job metadata (optional) - arbitrary json can be stored in {meta: {...}}"
// @Accept application/json, */*
// @Produce json
// @Success 200 {object} SuccessResponse
// @Router /{namespace}/jobs/{jobId} [put]
func PutJobById(c *gin.Context) {
	jobId := c.Param("jobId")
	namespace := c.Param("namespace")
	inUseBy := c.Query("in_use_by")
	jsonDataStr, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(500, ErrorResponse{"error reading request body"})
		return
	}
	fmt.Println(string(jsonDataStr))
	var metaData PutJobByIdBody
	err = json.Unmarshal(jsonDataStr, &metaData)
	if err != nil {
		c.AbortWithStatusJSON(500, ErrorResponse{"error unmarshaling json body"})
		return
	}
	job := Job{
		JobKey:  jobId,
		Created: time.Now(),
		InUse:   false,
		InUseBy: inUseBy,
		Meta:    metaData.Meta,
	}

	err = rdb.Set(jobKey(jobId, namespace), job, 0).Err()
	if err != nil {
		c.AbortWithStatusJSON(500, ErrorResponse{"error creating job"})
		return
	}
	c.JSON(200, SuccessResponse{"success"})
}

// DeleteJobById godoc
// @Summary Delete a job with specified id
// @Description Deletes a job with specified id
// @Tags root
// @Param namespace path string true "Namespace of job(s)"
// @Param jobId path string true "Job ID"
// @Param api_key query string true "API Key"
// @Accept */*
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /{namespace}/jobs/{jobId} [delete]
func DeleteJobById(c *gin.Context) {
	jobId := c.Param("jobId")
	namespace := c.Param("namespace")
	err := rdb.Get(jobKey(jobId, namespace)).Err()
	if err != nil {
		if err == redis.Nil {
			c.AbortWithStatusJSON(404, ErrorResponse{"job not found"})
			return
		}
		c.AbortWithStatusJSON(500, ErrorResponse{"error getting job"})
		return
	}
	if err = rdb.Del(jobKey(jobId, namespace)).Err(); err != nil {
		c.AbortWithStatusJSON(500, ErrorResponse{"error deleting job"})
		return
	}
	c.JSON(200, SuccessResponse{"success"})
}

// GetJobById godoc
// @Summary Get a job with specified id
// @Description Gets a job with specified id
// @Tags root
// @Param namespace path string true "Namespace of job(s)"
// @Param jobId path string true "Job ID"
// @Param api_key query string true "API Key"
// @Accept */*
// @Produce json
// @Success 200 {object} JobResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /{namespace}/jobs/{jobId} [get]
func GetJobById(c *gin.Context) {
	jobId := c.Param("jobId")
	namespace := c.Param("namespace")
	job := Job{}
	jsonStr, err := rdb.Get(jobKey(jobId, namespace)).Bytes()
	if err != nil {
		if err == redis.Nil {
			c.AbortWithStatusJSON(404, ErrorResponse{"job not found"})
			return
		}
		c.AbortWithStatusJSON(500, ErrorResponse{"error getting job"})
		return
	}
	if err = job.UnmarshalBinary(jsonStr); err != nil {
		c.AbortWithStatusJSON(500, ErrorResponse{"error unmarshaling job"})
		return
	}
	c.JSON(200, JobResponse{job})
}

// GetUnusedJob godoc
// @Summary Get an unused job and lock it as in-use
// @Description Finds a job that either is not in-use or has been inactive for more than the specified time.
// @Tags root
// @Param namespace path string true "Namespace of job(s)"
// @Param min_age path int false "Minimum age of job (last_used_on) in minutes before assuming it's no longer in use (optional, defaults to never)"
// @Param api_key query string true "API Key"
// @Accept */*
// @Produce json
// @Success 200 {object} JobResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /{namespace}/job [get]
func GetUnusedJob(c *gin.Context) {
	namespace := c.Param("namespace")
	// Find an unused job, lock it, and return it
	jobs, err := getAllJobs(rdb, namespace)
	if err != nil {
		c.AbortWithStatusJSON(500, ErrorResponse{err.Error()})
		return
	}
	now := time.Now()
	minAgeStr := c.Request.URL.Query().Get("min_age")
	var minAge int
	if minAgeStr != "" {
		minAge, err = strconv.Atoi(minAgeStr)
		if err != nil {
			c.AbortWithStatusJSON(400, ErrorResponse{"error parsing min_age"})
		}
	}
	for _, job := range jobs {
		// Send back jobs that are not in use or haven't been last_used_on for more than min_age
		if !job.InUse || (minAge > 0 && now.Sub(job.LastInUse) > time.Minute*time.Duration(minAge)) {
			job.InUse = true
			job.LastInUse = time.Now()
			err := rdb.Set(jobKey(job.JobKey, namespace), job, 0).Err()
			if err != nil {
				c.AbortWithStatusJSON(500, ErrorResponse{err.Error()})
				return
			}
			c.JSON(200, JobResponse{job})
			return
		}
	}
	c.AbortWithStatusJSON(404, ErrorResponse{"no jobs available"})
}

// GetAllJobs godoc
// @Summary Get all jobs
// @Description Gets all jobs
// @Tags root
// @Param namespace path string true "Namespace of job(s)"
// @Param api_key query string true "API Key"
// @Accept */*
// @Produce json
// @Success 200 {object} JobsResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /{namespace}/jobs/ [get]
func GetAllJobs(c *gin.Context) {
	namespace := c.Param("namespace")
	jobs, err := getAllJobs(rdb, namespace)
	if err != nil {
		if _, ok := err.(*NotFoundError); ok {
			c.AbortWithStatusJSON(404, ErrorResponse{"no jobs found"})
			return
		}
		c.AbortWithStatusJSON(500, ErrorResponse{err.Error()})
		return
	} else if len(jobs) == 0 {
		c.AbortWithStatusJSON(404, ErrorResponse{"no jobs found"})
	}
	c.JSON(200, JobsResponse{jobs})
}

// PostJobStillInUse godoc
// @Summary Update a job, marking it as still in use
// @Description Used together with GET /job's min_age parameter so that inactive jobs can be reused. Fails if the job doesn't exist
// @Tags root
// @Param namespace path string true "Namespace of job(s)"
// @Param jobId path string true "Job ID"
// @Param api_key query string true "API Key"
// @Accept */*
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /{namespace}/jobs/{jobId} [post]
func PostJobStillInUse(c *gin.Context) {
	jobId := c.Param("jobId")
	namespace := c.Param("namespace")
	job := Job{}
	jsonStr, err := rdb.Get(jobKey(jobId, namespace)).Bytes()
	if err != nil {
		if err == redis.Nil {
			c.AbortWithStatusJSON(404, ErrorResponse{"job not found"})
			return
		}
		c.AbortWithStatusJSON(500, ErrorResponse{"error getting job"})
		return
	}
	if err = job.UnmarshalBinary(jsonStr); err != nil {
		c.AbortWithStatusJSON(500, ErrorResponse{"error unmarshaling job"})
		return
	}
	job.InUse = true
	job.LastInUse = time.Now()
	err = rdb.Set(jobKey(jobId, namespace), job, 0).Err()
	if err != nil {
		c.AbortWithStatusJSON(500, ErrorResponse{"error updating job"})
		return
	}
	c.JSON(200, SuccessResponse{"success"})
}

func getAllJobs(rdb *redis.Client, namespace string) ([]Job, error) {
	keys, err := rdb.Keys(JobPrefix + namespace + ":" + "*").Result()
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
