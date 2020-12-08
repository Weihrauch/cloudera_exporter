/*
 *
 * title           :collector/yarn_module.go
 * description     :Submodule Collector for the Cluster YARN metrics
 * author		       :Alejandro Villegas
 * date            :2019/02/04
 * version         :1.0
 *
 */
 package collector



/* ======================================================================
 * Dependencies and libraries
 * ====================================================================== */
 import (
    // Go Default libraries
    "context"
    "strings"
  
    // Own libraries
    jp "keedio/cloudera_exporter/json_parser"
    log "keedio/cloudera_exporter/logger"
  
    // Go Prometheus libraries
    "github.com/prometheus/client_golang/prometheus"
  )
  
  
  
  
  /* ======================================================================
   * Data Structs
   * ====================================================================== */
  // None
  
  
  
  /* ======================================================================
   * Constants with the Host module TSquery sentences
   * ====================================================================== */
  const YARN_SCRAPER_NAME = "yarn"
  var (
    // Agent Queries
    YARN_ACTIVE_APPLICATIONS                                        ="SELECT LAST(yarn_active_applications) WHERE serviceType=\"YARN\""
    YARN_ACTIVE_APPLICATIONS_CUMULATIVE                             ="SELECT LAST(yarn_active_applications_cumulative) WHERE serviceType=\"YARN\""
    YARN_AGGREGATE_CONTAINERS_ALLOCATED_CUMULATIVE_RATE             ="SELECT LAST(yarn_aggregate_containers_allocated_cumulative_rate) WHERE serviceType=\"YARN\""
    YARN_AGGREGATE_CONTAINERS_ALLOCATED_RATE                        ="SELECT LAST(yarn_aggregate_containers_allocated_rate) WHERE serviceType=\"YARN\""
    YARN_AGGREGATE_CONTAINERS_RELEASED_CUMULATIVE_RATE              ="SELECT LAST(yarn_aggregate_containers_released_cumulative_rate) WHERE serviceType=\"YARN\""
    YARN_AGGREGATE_CONTAINERS_RELEASED_RATE                         ="SELECT LAST(yarn_aggregate_containers_released_rate) WHERE serviceType=\"YARN\""
    YARN_ALLOCATED_CONTAINERS                                       ="SELECT LAST(yarn_allocated_containers) WHERE serviceType=\"YARN\""
    YARN_ALLOCATED_CONTAINERS_CUMULATIVE                            ="SELECT LAST(yarn_allocated_containers_cumulative) WHERE serviceType=\"YARN\""
    YARN_ALLOCATED_MEMORY_MB                                        ="SELECT LAST(yarn_allocated_memory_mb) WHERE serviceType=\"YARN\""
    YARN_ALLOCATED_MEMORY_MB_CUMULATIVE                             ="SELECT LAST(yarn_allocated_memory_mb_cumulative) WHERE serviceType=\"YARN\""
    YARN_ALLOCATED_MEMORY_MB_WITH_PENDING_CONTAINERS                ="SELECT LAST(yarn_allocated_memory_mb_with_pending_containers) WHERE serviceType=\"YARN\""
    YARN_ALLOCATED_VCORES                                           ="SELECT LAST(yarn_allocated_vcores) WHERE serviceType=\"YARN\""
    YARN_ALLOCATED_VCORES_CUMULATIVE                                ="SELECT LAST(yarn_allocated_vcores_cumulative) WHERE serviceType=\"YARN\""
    YARN_ALLOCATED_VCORES_WITH_PENDING_CONTAINERS                   ="SELECT LAST(yarn_allocated_vcores_with_pending_containers) WHERE serviceType=\"YARN\""
    YARN_APPS_COMPLETED_CUMULATIVE_RATE                             ="SELECT LAST(yarn_apps_completed_cumulative_rate) WHERE serviceType=\"YARN\""
    YARN_APPS_COMPLETED_RATE                                        ="SELECT LAST(yarn_apps_completed_rate) WHERE serviceType=\"YARN\""
    YARN_APPS_FAILED_CUMULATIVE_RATE                                ="SELECT LAST(yarn_apps_failed_cumulative_rate) WHERE serviceType=\"YARN\""
    YARN_APPS_FAILED_RATE                                           ="SELECT LAST(yarn_apps_failed_rate) WHERE serviceType=\"YARN\""
    YARN_APPS_INGESTED_RATE                                         ="SELECT LAST(yarn_apps_ingested_rate) WHERE serviceType=\"YARN\""
    YARN_APPS_KILLED_CUMULATIVE_RATE                                ="SELECT LAST(yarn_apps_killed_cumulative_rate) WHERE serviceType=\"YARN\""
    YARN_APPS_KILLED_RATE                                           ="SELECT LAST(yarn_apps_killed_rate) WHERE serviceType=\"YARN\""
    YARN_APPS_PENDING                                               ="SELECT LAST(yarn_apps_pending) WHERE serviceType=\"YARN\""
    YARN_APPS_PENDING_CUMULATIVE                                    ="SELECT LAST(yarn_apps_pending_cumulative) WHERE serviceType=\"YARN\""
    YARN_APPS_RUNNING                                               ="SELECT LAST(yarn_apps_running) WHERE serviceType=\"YARN\""
    YARN_APPS_RUNNING_BETWEEN_300TO1440_MINS                        ="SELECT LAST(yarn_apps_running_between_300to1440_mins) WHERE serviceType=\"YARN\""
    YARN_APPS_RUNNING_BETWEEN_300TO1440_MINS_CUMULATIVE             ="SELECT LAST(yarn_apps_running_between_300to1440_mins_cumulative) WHERE serviceType=\"YARN\""
    YARN_APPS_RUNNING_BETWEEN_60TO300_MINS                          ="SELECT LAST(yarn_apps_running_between_60to300_mins) WHERE serviceType=\"YARN\""
    YARN_APPS_RUNNING_BETWEEN_60TO300_MINS_CUMULATIVE               ="SELECT LAST(yarn_apps_running_between_60to300_mins_cumulative) WHERE serviceType=\"YARN\""
    YARN_APPS_RUNNING_CUMULATIVE                                    ="SELECT LAST(yarn_apps_running_cumulative) WHERE serviceType=\"YARN\""
    YARN_APPS_RUNNING_OVER_1440_MINS                                ="SELECT LAST(yarn_apps_running_over_1440_mins) WHERE serviceType=\"YARN\""
    YARN_APPS_RUNNING_OVER_1440_MINS_CUMULATIVE                     ="SELECT LAST(yarn_apps_running_over_1440_mins_cumulative) WHERE serviceType=\"YARN\""
    YARN_APPS_RUNNING_WITHIN_60_MINS                                ="SELECT LAST(yarn_apps_running_within_60_mins) WHERE serviceType=\"YARN\""
    YARN_APPS_RUNNING_WITHIN_60_MINS_CUMULATIVE                     ="SELECT LAST(yarn_apps_running_within_60_mins_cumulative) WHERE serviceType=\"YARN\""
    YARN_APPS_SUBMITTED_CUMULATIVE_RATE                             ="SELECT LAST(yarn_apps_submitted_cumulative_rate) WHERE serviceType=\"YARN\""
    YARN_APPS_SUBMITTED_RATE                                        ="SELECT LAST(yarn_apps_submitted_rate) WHERE serviceType=\"YARN\""
    YARN_AVAILABLE_MEMORY_MB                                        ="SELECT LAST(yarn_available_memory_mb) WHERE serviceType=\"YARN\""
    YARN_AVAILABLE_VCORES                                           ="SELECT LAST(yarn_available_vcores) WHERE serviceType=\"YARN\""
    YARN_CONTAINER_WAIT_RATIO                                       ="SELECT LAST(yarn_container_wait_ratio) WHERE serviceType=\"YARN\""
    YARN_FAIR_SHARE_MB                                              ="SELECT LAST(yarn_fair_share_mb) WHERE serviceType=\"YARN\""
    YARN_FAIR_SHARE_MB_CUMULATIVE                                   ="SELECT LAST(yarn_fair_share_mb_cumulative) WHERE serviceType=\"YARN\""
    YARN_FAIR_SHARE_MB_WITH_PENDING_CONTAINERS                      ="SELECT LAST(yarn_fair_share_mb_with_pending_containers) WHERE serviceType=\"YARN\""
    YARN_FAIR_SHARE_VCORES                                          ="SELECT LAST(yarn_fair_share_vcores) WHERE serviceType=\"YARN\""
    YARN_FAIR_SHARE_VCORES_CUMULATIVE                               ="SELECT LAST(yarn_fair_share_vcores_cumulative) WHERE serviceType=\"YARN\""
    YARN_FAIR_SHARE_VCORES_WITH_PENDING_CONTAINERS                  ="SELECT LAST(yarn_fair_share_vcores_with_pending_containers) WHERE serviceType=\"YARN\""
    YARN_MAX_SHARE_MB                                               ="SELECT LAST(yarn_max_share_mb) WHERE serviceType=\"YARN\""
    YARN_MAX_SHARE_MB_CUMULATIVE                                    ="SELECT LAST(yarn_max_share_mb_cumulative) WHERE serviceType=\"YARN\""
    YARN_MAX_SHARE_VCORES                                           ="SELECT LAST(yarn_max_share_vcores) WHERE serviceType=\"YARN\""
    YARN_MAX_SHARE_VCORES_CUMULATIVE                                ="SELECT LAST(yarn_max_share_vcores_cumulative) WHERE serviceType=\"YARN\""
    YARN_MIN_SHARE_MB                                               ="SELECT LAST(yarn_min_share_mb) WHERE serviceType=\"YARN\""
    YARN_MIN_SHARE_MB_CUMULATIVE                                    ="SELECT LAST(yarn_min_share_mb_cumulative) WHERE serviceType=\"YARN\""
    YARN_MIN_SHARE_VCORES                                           ="SELECT LAST(yarn_min_share_vcores) WHERE serviceType=\"YARN\""
    YARN_MIN_SHARE_VCORES_CUMULATIVE                                ="SELECT LAST(yarn_min_share_vcores_cumulative) WHERE serviceType=\"YARN\""
    YARN_PENDING_CONTAINERS                                         ="SELECT LAST(yarn_pending_containers) WHERE serviceType=\"YARN\""
    YARN_PENDING_CONTAINERS_CUMULATIVE                              ="SELECT LAST(yarn_pending_containers_cumulative) WHERE serviceType=\"YARN\""
    YARN_PENDING_MEMORY_MB                                          ="SELECT LAST(yarn_pending_memory_mb) WHERE serviceType=\"YARN\""
    YARN_PENDING_MEMORY_MB_CUMULATIVE                               ="SELECT LAST(yarn_pending_memory_mb_cumulative) WHERE serviceType=\"YARN\""
    YARN_PENDING_VCORES                                             ="SELECT LAST(yarn_pending_vcores) WHERE serviceType=\"YARN\""
    YARN_PENDING_VCORES_CUMULATIVE                                  ="SELECT LAST(yarn_pending_vcores_cumulative) WHERE serviceType=\"YARN\""
    YARN_QUERIES_INGESTED_RATE                                      ="SELECT LAST(yarn_queries_ingested_rate) WHERE serviceType=\"YARN\""
    YARN_RESERVED_CONTAINERS                                        ="SELECT LAST(yarn_reserved_containers) WHERE serviceType=\"YARN\""
    YARN_RESERVED_CONTAINERS_CUMULATIVE                             ="SELECT LAST(yarn_reserved_containers_cumulative) WHERE serviceType=\"YARN\""
    YARN_RESERVED_MEMORY_MB                                         ="SELECT LAST(yarn_reserved_memory_mb) WHERE serviceType=\"YARN\""
    YARN_RESERVED_MEMORY_MB_CUMULATIVE                              ="SELECT LAST(yarn_reserved_memory_mb_cumulative) WHERE serviceType=\"YARN\""
    YARN_RESERVED_VCORES                                            ="SELECT LAST(yarn_reserved_vcores) WHERE serviceType=\"YARN\""
    YARN_RESERVED_VCORES_CUMULATIVE                                 ="SELECT LAST(yarn_reserved_vcores_cumulative) WHERE serviceType=\"YARN\""
    YARN_STEADY_FAIR_SHARE_MB                                       ="SELECT LAST(yarn_steady_fair_share_mb) WHERE serviceType=\"YARN\""
    YARN_STEADY_FAIR_SHARE_MB_CUMULATIVE                            ="SELECT LAST(yarn_steady_fair_share_mb_cumulative) WHERE serviceType=\"YARN\""
    YARN_STEADY_FAIR_SHARE_MB_WITH_PENDING_CONTAINERS               ="SELECT LAST(yarn_steady_fair_share_mb_with_pending_containers) WHERE serviceType=\"YARN\""
    YARN_STEADY_FAIR_SHARE_VCORES                                   ="SELECT LAST(yarn_steady_fair_share_vcores) WHERE serviceType=\"YARN\""
    YARN_STEADY_FAIR_SHARE_VCORES_CUMULATIVE                        ="SELECT LAST(yarn_steady_fair_share_vcores_cumulative) WHERE serviceType=\"YARN\""
    YARN_STEADY_FAIR_SHARE_VCORES_WITH_PENDING_CONTAINERS           ="SELECT LAST(yarn_steady_fair_share_vcores_with_pending_containers) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_APPLICATION_DURATION_RATE                      ="SELECT LAST(yarn_application_application_duration_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_REDUCES_RATE                                   ="SELECT LAST(yarn_application_reduces_rate) WHERE serviceType=\"YARN\""

    ALERTS_RATE                                                     ="SELECT LAST(alerts_rate) WHERE serviceType=\"YARN\""
    APPS_INGESTED_RATE                                              ="SELECT LAST(apps_ingested_rate) WHERE serviceType=\"YARN\""
    EVENTS_CRITICAL_RATE                                            ="SELECT LAST(events_critical_rate) WHERE serviceType=\"YARN\""
    EVENTS_IMPORTANT_RATE                                           ="SELECT LAST(events_important_rate) WHERE serviceType=\"YARN\""
    EVENTS_INFORMATIONAL_RATE                                       ="SELECT LAST(events_informational_rate) WHERE serviceType=\"YARN\""
    HEALTH_BAD_RATE                                                 ="SELECT LAST(health_bad_rate) WHERE serviceType=\"YARN\""
    HEALTH_CONCERNING_RATE                                          ="SELECT LAST(health_concerning_rate) WHERE serviceType=\"YARN\""
    HEALTH_DISABLED_RATE                                            ="SELECT LAST(health_disabled_rate) WHERE serviceType=\"YARN\""
    HEALTH_GOOD_RATE                                                ="SELECT LAST(health_good_rate) WHERE serviceType=\"YARN\""
    HEALTH_UNKNOWN_RATE                                             ="SELECT LAST(health_unknown_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_ADL_BYTES_READ_RATE                            ="SELECT LAST(yarn_application_adl_bytes_read_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_ADL_BYTES_WRITTEN_RATE                         ="SELECT LAST(yarn_application_adl_bytes_written_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_CM_CPU_MILLISECONDS_RATE                       ="SELECT LAST(yarn_application_cm_cpu_milliseconds_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_CPU_MILLISECONDS_RATE                          ="SELECT LAST(yarn_application_cpu_milliseconds_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_FILE_BYTES_READ_RATE                           ="SELECT LAST(yarn_application_file_bytes_read_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_FILE_BYTES_WRITTEN_RATE                        ="SELECT LAST(yarn_application_file_bytes_written_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_HDFS_BYTES_READ_RATE                           ="SELECT LAST(yarn_application_hdfs_bytes_read_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_HDFS_BYTES_WRITTEN_RATE                        ="SELECT LAST(yarn_application_hdfs_bytes_written_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_MAPS_RATE                                      ="SELECT LAST(yarn_application_maps_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_MB_MILLIS_MAPS_RATE                            ="SELECT LAST(yarn_application_mb_millis_maps_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_MB_MILLIS_REDUCES_RATE                         ="SELECT LAST(yarn_application_mb_millis_reduces_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_S3A_BYTES_READ_RATE                            ="SELECT LAST(yarn_application_s3a_bytes_read_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_S3A_BYTES_WRITTEN_RATE                         ="SELECT LAST(yarn_application_s3a_bytes_written_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_VCORES_MILLIS_MAPS_RATE                        ="SELECT LAST(yarn_application_vcores_millis_maps_rate) WHERE serviceType=\"YARN\""
    YARN_APPLICATION_VCORES_MILLIS_REDUCES_RATE                     ="SELECT LAST(yarn_application_vcores_millis_reduces_rate) WHERE serviceType=\"YARN\""
    YARN_REPORTS_CONTAINERS_ALLOCATED_MEMORY                        ="SELECT LAST(yarn_reports_containers_allocated_memory) WHERE serviceType=\"YARN\""
    YARN_REPORTS_CONTAINERS_ALLOCATED_VCORES                        ="SELECT LAST(yarn_reports_containers_allocated_vcores) WHERE serviceType=\"YARN\""
    YARN_REPORTS_CONTAINERS_USED_CPU                                ="SELECT LAST(yarn_reports_containers_used_cpu) WHERE serviceType=\"YARN\""
    YARN_REPORTS_CONTAINERS_USED_MEMORY                             ="SELECT LAST(yarn_reports_containers_used_memory) WHERE serviceType=\"YARN\""
    YARN_REPORTS_CONTAINERS_USED_VCORES                             ="SELECT LAST(yarn_reports_containers_used_vcores) WHERE serviceType=\"YARN\""
    YARN_REPORTS_USAGE_AGGREGATION_DURATION                         ="SELECT LAST(yarn_reports_usage_aggregation_duration) WHERE serviceType=\"YARN\""
    YARN_REPORTS_USAGE_APPS_WITH_METADATA                           ="SELECT LAST(yarn_reports_usage_apps_with_metadata) WHERE serviceType=\"YARN\""
    YARN_REPORTS_USAGE_APPS_WITHOUT_METADATA                        ="SELECT LAST(yarn_reports_usage_apps_without_metadata) WHERE serviceType=\"YARN\""
  )
  
  
  
  
  /* ======================================================================
   * Global variables
   * ====================================================================== */
  var (
    yarn_active_applications                                      =create_yarn_metric_struct("yarn_active_applications",	"Number of YARN applications in this pool with unsatisfied resource requests.	Applications")
    yarn_active_applications_cumulative                           =create_yarn_metric_struct("yarn_active_applications_cumulative",	"Number of YARN applications in this pool with unsatisfied resource requests. Includes this pool and any children.	Applications")
    yarn_aggregate_containers_allocated_cumulative_rate           =create_yarn_metric_struct("yarn_aggregate_containers_allocated_cumulative_rate",	"Aggregate Containers Allocated (Cumulative)	containers per second")
    yarn_aggregate_containers_allocated_rate                      =create_yarn_metric_struct("yarn_aggregate_containers_allocated_rate",	"Aggregate Containers Allocated	containers per second")
    yarn_aggregate_containers_released_cumulative_rate            =create_yarn_metric_struct("yarn_aggregate_containers_released_cumulative_rate",	"Aggregate Containers Released (Cumulative)	containers per second")
    yarn_aggregate_containers_released_rate                       =create_yarn_metric_struct("yarn_aggregate_containers_released_rate",	"Aggregate Containers Released	containers per second")
    yarn_allocated_containers                                     =create_yarn_metric_struct("yarn_allocated_containers",	"Number of containers allocated in this pool.	containers")
    yarn_allocated_containers_cumulative                          =create_yarn_metric_struct("yarn_allocated_containers_cumulative",	"Number of containers allocated in this pool. Includes this pool and its children.	containers")
    yarn_allocated_memory_mb                                      =create_yarn_metric_struct("yarn_allocated_memory_mb",	"Allocated Memory	MB")
    yarn_allocated_memory_mb_cumulative                           =create_yarn_metric_struct("yarn_allocated_memory_mb_cumulative",	"Allocated Memory (Cumulative)	MB")
    yarn_allocated_memory_mb_with_pending_containers              =create_yarn_metric_struct("yarn_allocated_memory_mb_with_pending_containers",	"Allocated memory when pending containers is more than zero.	MB")
    yarn_allocated_vcores                                         =create_yarn_metric_struct("yarn_allocated_vcores",	"Number of vcores allocated in this pool	VCores")
    yarn_allocated_vcores_cumulative                              =create_yarn_metric_struct("yarn_allocated_vcores_cumulative",	"Number of vcores allocated in this pool. Includes this pool and its children.	VCores")
    yarn_allocated_vcores_with_pending_containers                 =create_yarn_metric_struct("yarn_allocated_vcores_with_pending_containers",	"Allocated vcores when pending containers is more than zero.	VCores")
    yarn_apps_completed_cumulative_rate                           =create_yarn_metric_struct("yarn_apps_completed_cumulative_rate",	"Number of completed YARN applications in this pool. Includes this pool and any children.	Applications per second")
    yarn_apps_completed_rate                                      =create_yarn_metric_struct("yarn_apps_completed_rate",	"Number of completed YARN applications in this pool.	Applications per second")
    yarn_apps_failed_cumulative_rate                              =create_yarn_metric_struct("yarn_apps_failed_cumulative_rate",	"Number of failed YARN applications in this pool. Includes this pool and any children.	Applications per second")
    yarn_apps_failed_rate                                         =create_yarn_metric_struct("yarn_apps_failed_rate",	"Number of failed YARN applications in this pool.	Applications per second")
    yarn_apps_ingested_rate                                       =create_yarn_metric_struct("yarn_apps_ingested_rate",	"YARN applications ingested by the Service Monitor	Applications per second")
    yarn_apps_killed_cumulative_rate                              =create_yarn_metric_struct("yarn_apps_killed_cumulative_rate",	"Number of killed YARN applications in this pool. Includes this pool and any children.	Applications per second")
    yarn_apps_killed_rate                                         =create_yarn_metric_struct("yarn_apps_killed_rate",	"Number of killed YARN applications in this pool.	Applications per second")
    yarn_apps_pending                                             =create_yarn_metric_struct("yarn_apps_pending",	"Number of pending (queued) YARN applications in this pool.	Applications")
    yarn_apps_pending_cumulative                                  =create_yarn_metric_struct("yarn_apps_pending_cumulative",	"Number of pending (queued) YARN applications in this pool. Includes this pool and any children.	Applications")
    yarn_apps_running                                             =create_yarn_metric_struct("yarn_apps_running",	"Number of running YARN applications in this pool.	Applications")
    yarn_apps_running_between_300to1440_mins                      =create_yarn_metric_struct("yarn_apps_running_between_300to1440_mins",	"Number of running YARN applications in this pool for which elapsed time is between 300 and 1440 minutes.	Applications")
    yarn_apps_running_between_300to1440_mins_cumulative           =create_yarn_metric_struct("yarn_apps_running_between_300to1440_mins_cumulative",	"Number of running YARN applications in this pool for which elapsed time is between 300 and 1440 minutes. Includes this pool and any children.	Applications")
    yarn_apps_running_between_60to300_mins                        =create_yarn_metric_struct("yarn_apps_running_between_60to300_mins",	"Number of running YARN applications in this pool for which elapsed time is between 60 and 300 minutes.	Applications")
    yarn_apps_running_between_60to300_mins_cumulative             =create_yarn_metric_struct("yarn_apps_running_between_60to300_mins_cumulative",	"Number of running YARN applications in this pool for which elapsed time is between 60 and 300 minutes. Includes this pool and any children.	Applications")
    yarn_apps_running_cumulative                                  =create_yarn_metric_struct("yarn_apps_running_cumulative",	"Number of running YARN applications in this pool. Includes this pool and any children.	Applications")
    yarn_apps_running_over_1440_mins                              =create_yarn_metric_struct("yarn_apps_running_over_1440_mins",	"Number of running YARN applications in this pool for which elapsed time is more than 1440 minutes.	Applications")
    yarn_apps_running_over_1440_mins_cumulative                   =create_yarn_metric_struct("yarn_apps_running_over_1440_mins_cumulative",	"Number of running YARN applications in this pool for which elapsed time is more than 1440 minutes. Includes this pool and any children.	Applications")
    yarn_apps_running_within_60_mins                              =create_yarn_metric_struct("yarn_apps_running_within_60_mins",	"Number of running YARN applications in this pool for which elapsed time is less than 60 minutes.	Applications")
    yarn_apps_running_within_60_mins_cumulative                   =create_yarn_metric_struct("yarn_apps_running_within_60_mins_cumulative",	"Number of running YARN applications in this pool for which elapsed time is less than 60 minutes. Includes this pool and any children.	Applications")
    yarn_apps_submitted_cumulative_rate                           =create_yarn_metric_struct("yarn_apps_submitted_cumulative_rate",	"Number of submitted YARN applications in this pool. Includes this pool and any children.	Applications per second")
    yarn_apps_submitted_rate                                      =create_yarn_metric_struct("yarn_apps_submitted_rate",	"Number of submitted YARN applications in this queue.	Applications per second")
    yarn_available_memory_mb                                      =create_yarn_metric_struct("yarn_available_memory_mb",	"Memory not allocated to YARN containers.	MB")
    yarn_available_vcores                                         =create_yarn_metric_struct("yarn_available_vcores",	"Available vcores that can be used for containers.	VCores")
    yarn_container_wait_ratio                                     =create_yarn_metric_struct("yarn_container_wait_ratio",	"Percent of pending containers when pending containers is more than zero.	percent")
    yarn_fair_share_mb                                            =create_yarn_metric_struct("yarn_fair_share_mb",	"Fair share of memory in this pool.	MB")
    yarn_fair_share_mb_cumulative                                 =create_yarn_metric_struct("yarn_fair_share_mb_cumulative",	"Fair share of memory in this pool. Includes this pool and its children.	MB")
    yarn_fair_share_mb_with_pending_containers                    =create_yarn_metric_struct("yarn_fair_share_mb_with_pending_containers",	"Fair share of memory when pending containers is more than zero.	MB")
    yarn_fair_share_vcores                                        =create_yarn_metric_struct("yarn_fair_share_vcores",	"Fair share of vcores in this pool.	VCores")
    yarn_fair_share_vcores_cumulative                             =create_yarn_metric_struct("yarn_fair_share_vcores_cumulative",	"Fair share of vcores in this pool. Includes this pool and its children.	VCores")
    yarn_fair_share_vcores_with_pending_containers                =create_yarn_metric_struct("yarn_fair_share_vcores_with_pending_containers",	"Fair share of vcores when pending containers is more than zero.	VCores")
    yarn_max_share_mb                                             =create_yarn_metric_struct("yarn_max_share_mb",	"Maximum share of memory configured for this pool.	MB")
    yarn_max_share_mb_cumulative                                  =create_yarn_metric_struct("yarn_max_share_mb_cumulative",	"Maximum share of memory configured for this pool. Includes this pool and its children.	MB")
    yarn_max_share_vcores                                         =create_yarn_metric_struct("yarn_max_share_vcores",	"Maximum share of vcores configured for this pool.	VCores")
    yarn_max_share_vcores_cumulative                              =create_yarn_metric_struct("yarn_max_share_vcores_cumulative",	"Maximum share of vcores configured for this pool. Includes this pool and its children.	VCores")
    yarn_min_share_mb                                             =create_yarn_metric_struct("yarn_min_share_mb",	"Minimum share of memory configured for this pool.	MB")
    yarn_min_share_mb_cumulative                                  =create_yarn_metric_struct("yarn_min_share_mb_cumulative",	"Minimum share of memory configured for this pool. Includes this pool and its children.	MB")
    yarn_min_share_vcores                                         =create_yarn_metric_struct("yarn_min_share_vcores",	"Minimum share of vcores configured for this pool.	VCores")
    yarn_min_share_vcores_cumulative                              =create_yarn_metric_struct("yarn_min_share_vcores_cumulative",	"Minimum share of vcores configured for this pool. Includes this pool and its children.	VCores")
    yarn_pending_containers                                       =create_yarn_metric_struct("yarn_pending_containers",	"Number of pending (queued) containers in this pool.	containers")
    yarn_pending_containers_cumulative                            =create_yarn_metric_struct("yarn_pending_containers_cumulative",	"Number of pending (queued) containers in this pool. Includes this pool and any children.	containers")
    yarn_pending_memory_mb                                        =create_yarn_metric_struct("yarn_pending_memory_mb",	"Sum of memory currently requested but not allocated for containers in this pool.	MB")
    yarn_pending_memory_mb_cumulative                             =create_yarn_metric_struct("yarn_pending_memory_mb_cumulative",	"Sum of memory currently requested but not allocated for containers in this pool. Includes this pool and any children.	MB")
    yarn_pending_vcores                                           =create_yarn_metric_struct("yarn_pending_vcores",	"Number of vcores requested but not yet allocated for this pool.	VCores")
    yarn_pending_vcores_cumulative                                =create_yarn_metric_struct("yarn_pending_vcores_cumulative",	"Number of vcores requested but not yet allocated for this pool. Includes this pool and its children.	VCores")
    yarn_queries_ingested_rate                                    =create_yarn_metric_struct("yarn_queries_ingested_rate",	"Impala queries ingested by the Service Monitor	queries per second")
    yarn_reserved_containers                                      =create_yarn_metric_struct("yarn_reserved_containers",	"Reserved containers for this pool.	containers")
    yarn_reserved_containers_cumulative                           =create_yarn_metric_struct("yarn_reserved_containers_cumulative",	"Reserved containers for this pool. Includes this pool and any children.	containers")
    yarn_reserved_memory_mb                                       =create_yarn_metric_struct("yarn_reserved_memory_mb",	"Reserved memory in this pool.	MB")
    yarn_reserved_memory_mb_cumulative                            =create_yarn_metric_struct("yarn_reserved_memory_mb_cumulative",	"Reserved memory in this pool. Includes this pool and any children.	MB")
    yarn_reserved_vcores                                          =create_yarn_metric_struct("yarn_reserved_vcores",	"Number of vcores set aside for this pool, but not used as a part of an allocation.	VCores")
    yarn_reserved_vcores_cumulative                               =create_yarn_metric_struct("yarn_reserved_vcores_cumulative",	"Number of vcores set aside for this pool, but not used as a part of an allocation. Includes this pool and its children.	VCores")
    yarn_steady_fair_share_mb                                     =create_yarn_metric_struct("yarn_steady_fair_share_mb",	"Steady fair share of memory in this pool.	MB")
    yarn_steady_fair_share_mb_cumulative                          =create_yarn_metric_struct("yarn_steady_fair_share_mb_cumulative",	"Steady fair share of memory in this pool. Includes this pool and its children.	MB")
    yarn_steady_fair_share_mb_with_pending_containers             =create_yarn_metric_struct("yarn_steady_fair_share_mb_with_pending_containers",	"Steady fair share of memory when pending containers is more than zero.	MB")
    yarn_steady_fair_share_vcores                                 =create_yarn_metric_struct("yarn_steady_fair_share_vcores",	"Steady fair share of vcores in this pool.	VCores")
    yarn_steady_fair_share_vcores_cumulative                      =create_yarn_metric_struct("yarn_steady_fair_share_vcores_cumulative",	"Steady fair share of vcores in this pool. Includes this pool and its children.	VCores")
    yarn_steady_fair_share_vcores_with_pending_containers         =create_yarn_metric_struct("yarn_steady_fair_share_vcores_with_pending_containers",	"Steady fair share of vcores when pending containers is more than zero.	VCores")
    yarn_application_reduces_rate                                 =create_yarn_metric_struct("yarn_application_reduces_rate",	"The number of reduce tasks in this MapReduce job. Called 'reduces_total' in searches.	items per second")
    alerts_rate                                                   =create_yarn_metric_struct("alerts_rate",	"The number of alerts.	events per second")
    apps_ingested_rate                                            =create_yarn_metric_struct("apps_ingested_rate",	"YARN applications ingested by the Service Monitor	Applications per second")
    events_critical_rate                                          =create_yarn_metric_struct("events_critical_rate",	"The number of critical events.	events per second")
    events_important_rate                                         =create_yarn_metric_struct("events_important_rate",	"The number of important events.	events per second")
    events_informational_rate                                     =create_yarn_metric_struct("events_informational_rate",	"The number of informational events.	events per second")
    health_bad_rate                                               =create_yarn_metric_struct("health_bad_rate",	"Percentage of Time with Bad Health	seconds per second")
    health_concerning_rate                                        =create_yarn_metric_struct("health_concerning_rate",	"Percentage of Time with Concerning Health	seconds per second")
    health_disabled_rate                                          =create_yarn_metric_struct("health_disabled_rate",	"Percentage of Time with Disabled Health	seconds per second")
    health_good_rate                                              =create_yarn_metric_struct("health_good_rate",	"Percentage of Time with Good Health	seconds per second")
    health_unknown_rate                                           =create_yarn_metric_struct("health_unknown_rate",	"Percentage of Time with Unknown Health	seconds per second")
    yarn_application_adl_bytes_read_rate                          =create_yarn_metric_struct("yarn_application_adl_bytes_read_rate",	"ADL bytes read. Called 'adl_bytes_read' in searches.	bytes per second")
    yarn_application_adl_bytes_written_rate                       =create_yarn_metric_struct("yarn_application_adl_bytes_written_rate",	"ADL bytes written. Called 'adl_bytes_written' in searches.	bytes per second")
    yarn_application_application_duration_rate                    =create_yarn_metric_struct("yarn_application_application_duration_rate",	"How long YARN took to execute this application. Called 'application_duration' in searches.	ms per second")
    yarn_application_cm_cpu_milliseconds_rate                     =create_yarn_metric_struct("yarn_application_cm_cpu_milliseconds_rate",	"yarn.analysis.cm_cpu_milliseconds.description	ms per second")
    yarn_application_cpu_milliseconds_rate                        =create_yarn_metric_struct("yarn_application_cpu_milliseconds_rate",	"CPU time. Called 'cpu_milliseconds' in searches.	ms per second")
    yarn_application_file_bytes_read_rate                         =create_yarn_metric_struct("yarn_application_file_bytes_read_rate",	"File bytes read. Called 'file_bytes_read' in searches.	bytes per second")
    yarn_application_file_bytes_written_rate                      =create_yarn_metric_struct("yarn_application_file_bytes_written_rate",	"File bytes written. Called 'file_bytes_written' in searches.	bytes per second")
    yarn_application_hdfs_bytes_read_rate                         =create_yarn_metric_struct("yarn_application_hdfs_bytes_read_rate",	"HDFS bytes read. Called 'hdfs_bytes_read' in searches.	bytes per second")
    yarn_application_hdfs_bytes_written_rate                      =create_yarn_metric_struct("yarn_application_hdfs_bytes_written_rate",	"HDFS bytes written. Called 'hdfs_bytes_written' in searches.	bytes per second")
    yarn_application_maps_rate                                    =create_yarn_metric_struct("yarn_application_maps_rate",	"The number of Map tasks in this MapReduce job. Called 'maps_total' in searches.	items per second")
    yarn_application_mb_millis_maps_rate                          =create_yarn_metric_struct("yarn_application_mb_millis_maps_rate",	"Map memory allocation. Called 'mb_millis_maps' in searches.	items per second")
    yarn_application_mb_millis_reduces_rate                       =create_yarn_metric_struct("yarn_application_mb_millis_reduces_rate",	"Reduce memory allocation. Called 'mb_millis_reduces' in searches.	items per second")
    yarn_application_s3a_bytes_read_rate                          =create_yarn_metric_struct("yarn_application_s3a_bytes_read_rate",	"S3A bytes read. Called 's3a_bytes_read' in searches.	bytes per second")
    yarn_application_s3a_bytes_written_rate                       =create_yarn_metric_struct("yarn_application_s3a_bytes_written_rate",	"S3A bytes written. Called 's3a_bytes_written' in searches.	bytes per second")
    yarn_application_vcores_millis_maps_rate                      =create_yarn_metric_struct("yarn_application_vcores_millis_maps_rate",	"Map CPU allocation. Called 'vcores_millis_maps' in searches.	items per second")
    yarn_application_vcores_millis_reduces_rate                   =create_yarn_metric_struct("yarn_application_vcores_millis_reduces_rate",	"Reduce CPU allocation. Called 'vcores_millis_reduces' in searches.	items per second")
    yarn_reports_containers_allocated_memory                      =create_yarn_metric_struct("yarn_reports_containers_allocated_memory",	"Memory allocated to YARN containers	MB seconds")
    yarn_reports_containers_allocated_vcores                      =create_yarn_metric_struct("yarn_reports_containers_allocated_vcores",	"VCores allocated to YARN containers	VCore seconds")
    yarn_reports_containers_used_cpu                              =create_yarn_metric_struct("yarn_reports_containers_used_cpu",	"CPU used by YARN containers	Percent seconds")
    yarn_reports_containers_used_memory                           =create_yarn_metric_struct("yarn_reports_containers_used_memory",	"Memory used by YARN containers	MB seconds")
    yarn_reports_containers_used_vcores                           =create_yarn_metric_struct("yarn_reports_containers_used_vcores",	"VCores used by YARN containers	VCore seconds")
    yarn_reports_usage_aggregation_duration                       =create_yarn_metric_struct("yarn_reports_usage_aggregation_duration",	"The duration for generating YARN usage reporting metrics by aggregating YARN container usage metrics	ms")
    yarn_reports_usage_apps_with_metadata                         =create_yarn_metric_struct("yarn_reports_usage_apps_with_metadata",	"YARN applications for which container usage was computed and also had metadata	Applications")
    yarn_reports_usage_apps_without_metadata                      =create_yarn_metric_struct("yarn_reports_usage_apps_without_metadata",	"YARN applications for which container usage was computed but did not have metadata	Applications")
  )
  var yarn_query_variable_relationship = []relation {
    {YARN_ACTIVE_APPLICATIONS,                                        *yarn_active_applications},
    {YARN_ACTIVE_APPLICATIONS_CUMULATIVE,                             *yarn_active_applications_cumulative},
    {YARN_AGGREGATE_CONTAINERS_ALLOCATED_CUMULATIVE_RATE,             *yarn_aggregate_containers_allocated_cumulative_rate},
    {YARN_AGGREGATE_CONTAINERS_ALLOCATED_RATE,                        *yarn_aggregate_containers_allocated_rate},
    {YARN_AGGREGATE_CONTAINERS_RELEASED_CUMULATIVE_RATE,              *yarn_aggregate_containers_released_cumulative_rate},
    {YARN_AGGREGATE_CONTAINERS_RELEASED_RATE,                         *yarn_aggregate_containers_released_rate},
    {YARN_ALLOCATED_CONTAINERS,                                       *yarn_allocated_containers},
    {YARN_ALLOCATED_CONTAINERS_CUMULATIVE,                            *yarn_allocated_containers_cumulative},
    {YARN_ALLOCATED_MEMORY_MB,                                        *yarn_allocated_memory_mb},
    {YARN_ALLOCATED_MEMORY_MB_CUMULATIVE,                             *yarn_allocated_memory_mb_cumulative},
    {YARN_ALLOCATED_MEMORY_MB_WITH_PENDING_CONTAINERS,                *yarn_allocated_memory_mb_with_pending_containers},
    {YARN_ALLOCATED_VCORES,                                           *yarn_allocated_vcores},
    {YARN_ALLOCATED_VCORES_CUMULATIVE,                                *yarn_allocated_vcores_cumulative},
    {YARN_ALLOCATED_VCORES_WITH_PENDING_CONTAINERS,                   *yarn_allocated_vcores_with_pending_containers},
    {YARN_APPS_COMPLETED_CUMULATIVE_RATE,                             *yarn_apps_completed_cumulative_rate},
    {YARN_APPS_COMPLETED_RATE,                                        *yarn_apps_completed_rate},
    {YARN_APPS_FAILED_CUMULATIVE_RATE,                                *yarn_apps_failed_cumulative_rate},
    {YARN_APPS_FAILED_RATE,                                           *yarn_apps_failed_rate},
    {YARN_APPS_INGESTED_RATE,                                         *yarn_apps_ingested_rate},
    {YARN_APPS_KILLED_CUMULATIVE_RATE,                                *yarn_apps_killed_cumulative_rate},
    {YARN_APPS_KILLED_RATE,                                           *yarn_apps_killed_rate},
    {YARN_APPS_PENDING,                                               *yarn_apps_pending},
    {YARN_APPS_PENDING_CUMULATIVE,                                    *yarn_apps_pending_cumulative},
    {YARN_APPS_RUNNING,                                               *yarn_apps_running},
    {YARN_APPS_RUNNING_BETWEEN_300TO1440_MINS,                        *yarn_apps_running_between_300to1440_mins},
    {YARN_APPS_RUNNING_BETWEEN_300TO1440_MINS_CUMULATIVE,             *yarn_apps_running_between_300to1440_mins_cumulative},
    {YARN_APPS_RUNNING_BETWEEN_60TO300_MINS,                          *yarn_apps_running_between_60to300_mins},
    {YARN_APPS_RUNNING_BETWEEN_60TO300_MINS_CUMULATIVE,               *yarn_apps_running_between_60to300_mins_cumulative},
    {YARN_APPS_RUNNING_CUMULATIVE,                                    *yarn_apps_running_cumulative},
    {YARN_APPS_RUNNING_OVER_1440_MINS,                                *yarn_apps_running_over_1440_mins},
    {YARN_APPS_RUNNING_OVER_1440_MINS_CUMULATIVE,                     *yarn_apps_running_over_1440_mins_cumulative},
    {YARN_APPS_RUNNING_WITHIN_60_MINS,                                *yarn_apps_running_within_60_mins},
    {YARN_APPS_RUNNING_WITHIN_60_MINS_CUMULATIVE,                     *yarn_apps_running_within_60_mins_cumulative},
    {YARN_APPS_SUBMITTED_CUMULATIVE_RATE,                             *yarn_apps_submitted_cumulative_rate},
    {YARN_APPS_SUBMITTED_RATE,                                        *yarn_apps_submitted_rate},
    {YARN_AVAILABLE_MEMORY_MB,                                        *yarn_available_memory_mb},
    {YARN_AVAILABLE_VCORES,                                           *yarn_available_vcores},
    {YARN_CONTAINER_WAIT_RATIO,                                       *yarn_container_wait_ratio},
    {YARN_FAIR_SHARE_MB,                                              *yarn_fair_share_mb},
    {YARN_FAIR_SHARE_MB_CUMULATIVE,                                   *yarn_fair_share_mb_cumulative},
    {YARN_FAIR_SHARE_MB_WITH_PENDING_CONTAINERS,                      *yarn_fair_share_mb_with_pending_containers},
    {YARN_FAIR_SHARE_VCORES,                                          *yarn_fair_share_vcores},
    {YARN_FAIR_SHARE_VCORES_CUMULATIVE,                               *yarn_fair_share_vcores_cumulative},
    {YARN_FAIR_SHARE_VCORES_WITH_PENDING_CONTAINERS,                  *yarn_fair_share_vcores_with_pending_containers},
    {YARN_MAX_SHARE_MB,                                               *yarn_max_share_mb},
    {YARN_MAX_SHARE_MB_CUMULATIVE,                                    *yarn_max_share_mb_cumulative},
    {YARN_MAX_SHARE_VCORES,                                           *yarn_max_share_vcores},
    {YARN_MAX_SHARE_VCORES_CUMULATIVE,                                *yarn_max_share_vcores_cumulative},
    {YARN_MIN_SHARE_MB,                                               *yarn_min_share_mb},
    {YARN_MIN_SHARE_MB_CUMULATIVE,                                    *yarn_min_share_mb_cumulative},
    {YARN_MIN_SHARE_VCORES,                                           *yarn_min_share_vcores},
    {YARN_MIN_SHARE_VCORES_CUMULATIVE,                                *yarn_min_share_vcores_cumulative},
    {YARN_PENDING_CONTAINERS,                                         *yarn_pending_containers},
    {YARN_PENDING_CONTAINERS_CUMULATIVE,                              *yarn_pending_containers_cumulative},
    {YARN_PENDING_MEMORY_MB,                                          *yarn_pending_memory_mb},
    {YARN_PENDING_MEMORY_MB_CUMULATIVE,                               *yarn_pending_memory_mb_cumulative},
    {YARN_PENDING_VCORES,                                             *yarn_pending_vcores},
    {YARN_PENDING_VCORES_CUMULATIVE,                                  *yarn_pending_vcores_cumulative},
    {YARN_QUERIES_INGESTED_RATE,                                      *yarn_queries_ingested_rate},
    {YARN_RESERVED_CONTAINERS,                                        *yarn_reserved_containers},
    {YARN_RESERVED_CONTAINERS_CUMULATIVE,                             *yarn_reserved_containers_cumulative},
    {YARN_RESERVED_MEMORY_MB,                                         *yarn_reserved_memory_mb},
    {YARN_RESERVED_MEMORY_MB_CUMULATIVE,                              *yarn_reserved_memory_mb_cumulative},
    {YARN_RESERVED_VCORES,                                            *yarn_reserved_vcores},
    {YARN_RESERVED_VCORES_CUMULATIVE,                                 *yarn_reserved_vcores_cumulative},
    {YARN_STEADY_FAIR_SHARE_MB,                                       *yarn_steady_fair_share_mb},
    {YARN_STEADY_FAIR_SHARE_MB_CUMULATIVE,                            *yarn_steady_fair_share_mb_cumulative},
    {YARN_STEADY_FAIR_SHARE_MB_WITH_PENDING_CONTAINERS,               *yarn_steady_fair_share_mb_with_pending_containers},
    {YARN_STEADY_FAIR_SHARE_VCORES,                                   *yarn_steady_fair_share_vcores},
    {YARN_STEADY_FAIR_SHARE_VCORES_CUMULATIVE,                        *yarn_steady_fair_share_vcores_cumulative},
    {YARN_STEADY_FAIR_SHARE_VCORES_WITH_PENDING_CONTAINERS,           *yarn_steady_fair_share_vcores_with_pending_containers},
    {YARN_APPLICATION_CM_CPU_MILLISECONDS_RATE,                       *yarn_application_cm_cpu_milliseconds_rate},
    {YARN_APPLICATION_MB_MILLIS_REDUCES_RATE,                         *yarn_application_mb_millis_reduces_rate},
    {YARN_APPLICATION_REDUCES_RATE,                                   *yarn_application_reduces_rate},    
    {ALERTS_RATE,                                                     *alerts_rate},
    {APPS_INGESTED_RATE,                                              *apps_ingested_rate},
    {EVENTS_CRITICAL_RATE,                                            *events_critical_rate},
    {EVENTS_IMPORTANT_RATE,                                           *events_important_rate},
    {EVENTS_INFORMATIONAL_RATE,                                       *events_informational_rate},
    {HEALTH_BAD_RATE,                                                 *health_bad_rate},
    {HEALTH_CONCERNING_RATE,                                          *health_concerning_rate},
    {HEALTH_DISABLED_RATE,                                            *health_disabled_rate},
    {HEALTH_GOOD_RATE,                                                *health_good_rate},
    {HEALTH_UNKNOWN_RATE,                                             *health_unknown_rate},
    {YARN_APPLICATION_ADL_BYTES_READ_RATE,                            *yarn_application_adl_bytes_read_rate},
    {YARN_APPLICATION_ADL_BYTES_WRITTEN_RATE,                         *yarn_application_adl_bytes_written_rate},
    {YARN_APPLICATION_APPLICATION_DURATION_RATE,                      *yarn_application_application_duration_rate},
    {YARN_APPLICATION_CPU_MILLISECONDS_RATE,                          *yarn_application_cpu_milliseconds_rate},
    {YARN_APPLICATION_FILE_BYTES_READ_RATE,                           *yarn_application_file_bytes_read_rate},
    {YARN_APPLICATION_FILE_BYTES_WRITTEN_RATE,                        *yarn_application_file_bytes_written_rate},
    {YARN_APPLICATION_HDFS_BYTES_READ_RATE,                           *yarn_application_hdfs_bytes_read_rate},
    {YARN_APPLICATION_HDFS_BYTES_WRITTEN_RATE,                        *yarn_application_hdfs_bytes_written_rate},
    {YARN_APPLICATION_MAPS_RATE,                                      *yarn_application_maps_rate},
    {YARN_APPLICATION_MB_MILLIS_MAPS_RATE,                            *yarn_application_mb_millis_maps_rate},
    {YARN_APPLICATION_S3A_BYTES_READ_RATE,                            *yarn_application_s3a_bytes_read_rate},
    {YARN_APPLICATION_S3A_BYTES_WRITTEN_RATE,                         *yarn_application_s3a_bytes_written_rate},
    {YARN_APPLICATION_VCORES_MILLIS_MAPS_RATE,                        *yarn_application_vcores_millis_maps_rate},
    {YARN_APPLICATION_VCORES_MILLIS_REDUCES_RATE,                     *yarn_application_vcores_millis_reduces_rate},
    {YARN_REPORTS_CONTAINERS_ALLOCATED_MEMORY,                        *yarn_reports_containers_allocated_memory},
    {YARN_REPORTS_CONTAINERS_ALLOCATED_VCORES,                        *yarn_reports_containers_allocated_vcores},
    {YARN_REPORTS_CONTAINERS_USED_CPU,                                *yarn_reports_containers_used_cpu},
    {YARN_REPORTS_CONTAINERS_USED_MEMORY,                             *yarn_reports_containers_used_memory},
    {YARN_REPORTS_CONTAINERS_USED_VCORES,                             *yarn_reports_containers_used_vcores},
    {YARN_REPORTS_USAGE_AGGREGATION_DURATION,                         *yarn_reports_usage_aggregation_duration},
    {YARN_REPORTS_USAGE_APPS_WITH_METADATA,                           *yarn_reports_usage_apps_with_metadata},
    {YARN_REPORTS_USAGE_APPS_WITHOUT_METADATA,                        *yarn_reports_usage_apps_without_metadata},
  }
 
 /* ======================================================================
  * Functions
  * ====================================================================== */
 // Create and returns a prometheus descriptor for a yarn metric. 
 // The "metric_name" parameter its mandatory
 // If the "description" parameter is empty, the function assings it with the
 // value of the name of the metric in uppercase and separated by spaces
 func create_yarn_metric_struct(metric_name string, description string) *prometheus.Desc {
   // Correct "description" parameter if is empty
   if len(description) == 0 {
     description = strings.Replace(strings.ToUpper(metric_name), "_", " ", 0)
   }
 
   // return prometheus descriptor
   return prometheus.NewDesc(
     prometheus.BuildFQName(namespace, YARN_SCRAPER_NAME, metric_name),
     description,
     []string{"cluster", "entityName","hostname"},
     nil,
   )
 }
 
 
 // Generic function to extract de metadata associated with the query value
 // Only for YARN metric type
 func create_yarn_metric (ctx context.Context, config Collector_connection_data, query string, metric_struct prometheus.Desc, ch chan<- prometheus.Metric) bool {
   // Make the query
   json_parsed, err := make_and_parse_timeseries_query(ctx, config, query)
   if err != nil {
     return false
   }
 
   // Get the num of hosts in the cluster or clusters
   num_ts_series, err := jp.Get_timeseries_num(json_parsed)
   if err != nil {
     return false
   }
 
   // Extract Metadata for each TimeSerie
   for ts_index := 0; ts_index < num_ts_series; ts_index ++ {
     // Get the Cluster Name
     cluster_name := jp.Get_timeseries_query_cluster(json_parsed, ts_index)
     // Get Host Name
     host_name := jp.Get_timeseries_query_host_name(json_parsed, ts_index)
     // Get the entity Name
     entity_name := jp.Get_timeseries_query_entity_name(json_parsed, ts_index)
     // Get Query LAST value
     value, err := jp.Get_timeseries_query_value(json_parsed, ts_index)
     if err != nil {
       continue
     }
     // Assing the data to the Prometheus descriptor
     ch <- prometheus.MustNewConstMetric(&metric_struct, prometheus.GaugeValue, value, cluster_name, entity_name, host_name)
   }
   return true
 }
 
 
 
 
 /* ======================================================================
  * Scrape "Class"
  * ====================================================================== */
 // ScrapeYARN struct
 type ScrapeYARN struct{}
 
 // Name of the Scraper. Should be unique.
 func (ScrapeYARN) Name() string {
   return YARN_SCRAPER_NAME
 }
 
 // Help describes the role of the Scraper.
 func (ScrapeYARN) Help() string {
   return "YARN Metrics"
 }
 
 // Version.
 func (ScrapeYARN) Version() float64 {
   return 1.0
 }
 
 // Scrape generic function. Override for host module.
 func (ScrapeYARN) Scrape (ctx context.Context, config *Collector_connection_data, ch chan<- prometheus.Metric) error {
   log.Debug_msg("Ejecutando YARN Metrics Scraper")
 
   // Queries counters
   success_queries := 0
   error_queries := 0
 
   // Execute the generic funtion for creation of metrics with the pairs (QUERY, PROM:DESCRIPTOR)
   for i:=0 ; i < len(yarn_query_variable_relationship) ; i++ {
     if create_yarn_metric(ctx, *config, yarn_query_variable_relationship[i].Query, yarn_query_variable_relationship[i].Metric_struct, ch) {
       success_queries += 1
     } else {
       error_queries += 1
     }
   }
   log.Debug_msg("In the YARN Module has been executed %d queries. %d success and %d with errors", success_queries + error_queries, success_queries, error_queries)
   return nil
 }
 
 var _ Scraper = ScrapeYARN{}
 