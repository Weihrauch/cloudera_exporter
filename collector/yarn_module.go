/*
 *
 * title           :collector/yarn_module.go
 * description     :Submodule Collector for the Cluster YARN metrics
 * author		       :Alejandro Villegas
 * co-author           :NTUMBA Phinées
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

    ACTIVE_APPLICATIONS                                ="SELECT LAST(active_applications) WHERE serviceType=\"YARN\""                                             
    ACTIVE_APPLICATIONS_CUMULATIVE                     ="SELECT LAST(active_applications_cumulative) WHERE serviceType=\"YARN\""                                  
    AGGREGATE_CONTAINERS_ALLOCATED_CUMULATIVE_RATE     ="SELECT LAST(aggregate_containers_allocated_cumulative_rate) WHERE serviceType=\"YARN\""                  
    AGGREGATE_CONTAINERS_ALLOCATED_RATE                ="SELECT LAST(aggregate_containers_allocated_rate) WHERE serviceType=\"YARN\""                             
    AGGREGATE_CONTAINERS_RELEASED_CUMULATIVE_RATE      ="SELECT LAST(aggregate_containers_released_cumulative_rate) WHERE serviceType=\"YARN\""                   
    AGGREGATE_CONTAINERS_RELEASED_RATE                 ="SELECT LAST(aggregate_containers_released_rate) WHERE serviceType=\"YARN\""                              
    ALERTS_RATE                                        ="SELECT LAST(alerts_rate) WHERE serviceType=\"YARN\""                                                     
    ALLOCATED_CONTAINERS                               ="SELECT LAST(allocated_containers) WHERE serviceType=\"YARN\""                                            
    ALLOCATED_CONTAINERS_CUMULATIVE                    ="SELECT LAST(allocated_containers_cumulative) WHERE serviceType=\"YARN\""                                 
    ALLOCATED_MEMORY_MB                                ="SELECT LAST(allocated_memory_mb) WHERE serviceType=\"YARN\""                                             
    ALLOCATED_MEMORY_MB_CUMULATIVE                     ="SELECT LAST(allocated_memory_mb_cumulative) WHERE serviceType=\"YARN\""                                  
    ALLOCATED_MEMORY_MB_WITH_PENDING_CONTAINERS        ="SELECT LAST(allocated_memory_mb_with_pending_containers) WHERE serviceType=\"YARN\""                     
    ALLOCATED_VCORES                                   ="SELECT LAST(allocated_vcores) WHERE serviceType=\"YARN\""                                                
    ALLOCATED_VCORES_CUMULATIVE                        ="SELECT LAST(allocated_vcores_cumulative) WHERE serviceType=\"YARN\""                                     
    ALLOCATED_VCORES_WITH_PENDING_CONTAINERS           ="SELECT LAST(allocated_vcores_with_pending_containers) WHERE serviceType=\"YARN\""                        
    APPLICATION_ADL_BYTES_READ_RATE                    ="SELECT LAST(yarn_application_adl_bytes_read_rate) WHERE serviceType=\"YARN\""                            
    APPLICATION_ADL_BYTES_WRITTEN_RATE                 ="SELECT LAST(yarn_application_adl_bytes_written_rate) WHERE serviceType=\"YARN\""                         
    APPLICATION_APPLICATION_DURATION_RATE              ="SELECT LAST(yarn_application_application_duration_rate) WHERE serviceType=\"YARN\""                      
    APPLICATION_CM_CPU_MILLISECONDS_RATE               ="SELECT LAST(yarn_application_cm_cpu_milliseconds_rate) WHERE serviceType=\"YARN\""                       
    APPLICATION_CPU_MILLISECONDS_RATE                  ="SELECT LAST(yarn_application_cpu_milliseconds_rate) WHERE serviceType=\"YARN\""                          
    APPLICATION_FILE_BYTES_READ_RATE                   ="SELECT LAST(yarn_application_file_bytes_read_rate) WHERE serviceType=\"YARN\""                           
    APPLICATION_FILE_BYTES_WRITTEN_RATE                ="SELECT LAST(yarn_application_file_bytes_written_rate) WHERE serviceType=\"YARN\""                        
    APPLICATION_HDFS_BYTES_READ_RATE                   ="SELECT LAST(yarn_application_hdfs_bytes_read_rate) WHERE serviceType=\"YARN\""                           
    APPLICATION_HDFS_BYTES_WRITTEN_RATE                ="SELECT LAST(yarn_application_hdfs_bytes_written_rate) WHERE serviceType=\"YARN\""                        
    APPLICATION_MAPS_RATE                              ="SELECT LAST(yarn_application_maps_rate) WHERE serviceType=\"YARN\""                                      
    APPLICATION_MB_MILLIS_MAPS_RATE                    ="SELECT LAST(yarn_application_mb_millis_maps_rate) WHERE serviceType=\"YARN\""                            
    APPLICATION_MB_MILLIS_REDUCES_RATE                 ="SELECT LAST(yarn_application_mb_millis_reduces_rate) WHERE serviceType=\"YARN\""                         
    APPLICATION_REDUCES_RATE                           ="SELECT LAST(yarn_application_reduces_rate) WHERE serviceType=\"YARN\""                                   
    APPLICATION_S3A_BYTES_READ_RATE                    ="SELECT LAST(yarn_application_s3a_bytes_read_rate) WHERE serviceType=\"YARN\""                            
    APPLICATION_S3A_BYTES_WRITTEN_RATE                 ="SELECT LAST(yarn_application_s3a_bytes_written_rate) WHERE serviceType=\"YARN\""                         
    APPLICATION_VCORES_MILLIS_MAPS_RATE                ="SELECT LAST(yarn_application_vcores_millis_maps_rate) WHERE serviceType=\"YARN\""                        
    APPLICATION_VCORES_MILLIS_REDUCES_RATE             ="SELECT LAST(yarn_application_vcores_millis_reduces_rate) WHERE serviceType=\"YARN\""                     
    APPS_COMPLETED_CUMULATIVE_RATE                     ="SELECT LAST(apps_completed_cumulative_rate) WHERE serviceType=\"YARN\""                                  
    APPS_COMPLETED_RATE                                ="SELECT LAST(apps_completed_rate) WHERE serviceType=\"YARN\""                                             
    APPS_FAILED_CUMULATIVE_RATE                        ="SELECT LAST(apps_failed_cumulative_rate) WHERE serviceType=\"YARN\""                                     
    APPS_FAILED_RATE                                   ="SELECT LAST(apps_failed_rate) WHERE serviceType=\"YARN\""                                                
    APPS_INGESTED_RATE                                 ="SELECT LAST(apps_ingested_rate) WHERE serviceType=\"YARN\""                                              
    APPS_KILLED_CUMULATIVE_RATE                        ="SELECT LAST(apps_killed_cumulative_rate) WHERE serviceType=\"YARN\""                                     
    APPS_KILLED_RATE                                   ="SELECT LAST(apps_killed_rate) WHERE serviceType=\"YARN\""                                                
    APPS_PENDING                                       ="SELECT LAST(apps_pending) WHERE serviceType=\"YARN\""                                                    
    APPS_PENDING_CUMULATIVE                            ="SELECT LAST(apps_pending_cumulative) WHERE serviceType=\"YARN\""                                         
    APPS_RUNNING                                       ="SELECT LAST(apps_running) WHERE serviceType=\"YARN\""                                                    
    APPS_RUNNING_BETWEEN_300TO1440_MINS                ="SELECT LAST(apps_running_between_300to1440_mins) WHERE serviceType=\"YARN\""                             
    APPS_RUNNING_BETWEEN_300TO1440_MINS_CUMULATIVE     ="SELECT LAST(apps_running_between_300to1440_mins_cumulative) WHERE serviceType=\"YARN\""                  
    APPS_RUNNING_BETWEEN_60TO300_MINS                  ="SELECT LAST(apps_running_between_60to300_mins) WHERE serviceType=\"YARN\""                               
    APPS_RUNNING_BETWEEN_60TO300_MINS_CUMULATIVE       ="SELECT LAST(apps_running_between_60to300_mins_cumulative) WHERE serviceType=\"YARN\""                    
    APPS_RUNNING_CUMULATIVE                            ="SELECT LAST(apps_running_cumulative) WHERE serviceType=\"YARN\""                                         
    APPS_RUNNING_OVER_1440_MINS                        ="SELECT LAST(apps_running_over_1440_mins) WHERE serviceType=\"YARN\""                                     
    APPS_RUNNING_OVER_1440_MINS_CUMULATIVE             ="SELECT LAST(apps_running_over_1440_mins_cumulative) WHERE serviceType=\"YARN\""                          
    APPS_RUNNING_WITHIN_60_MINS                        ="SELECT LAST(apps_running_within_60_mins) WHERE serviceType=\"YARN\""                                     
    APPS_RUNNING_WITHIN_60_MINS_CUMULATIVE             ="SELECT LAST(apps_running_within_60_mins_cumulative) WHERE serviceType=\"YARN\""                          
    APPS_SUBMITTED_CUMULATIVE_RATE                     ="SELECT LAST(apps_submitted_cumulative_rate) WHERE serviceType=\"YARN\""                                  
    APPS_SUBMITTED_RATE                                ="SELECT LAST(apps_submitted_rate) WHERE serviceType=\"YARN\""                                             
    AVAILABLE_MEMORY_MB                                ="SELECT LAST(available_memory_mb) WHERE serviceType=\"YARN\""                                             
    AVAILABLE_VCORES                                   ="SELECT LAST(available_vcores) WHERE serviceType=\"YARN\""                                                
    CONTAINER_WAIT_RATIO                               ="SELECT LAST(container_wait_ratio) WHERE serviceType=\"YARN\""                                            
    EVENTS_CRITICAL_RATE                               ="SELECT LAST(events_critical_rate) WHERE serviceType=\"YARN\""                                            
    EVENTS_IMPORTANT_RATE                              ="SELECT LAST(events_important_rate) WHERE serviceType=\"YARN\""                                           
    EVENTS_INFORMATIONAL_RATE                          ="SELECT LAST(events_informational_rate) WHERE serviceType=\"YARN\""                                       
    FAIR_SHARE_MB                                      ="SELECT LAST(fair_share_mb) WHERE serviceType=\"YARN\""                                                   
    FAIR_SHARE_MB_CUMULATIVE                           ="SELECT LAST(fair_share_mb_cumulative) WHERE serviceType=\"YARN\""                                        
    FAIR_SHARE_MB_WITH_PENDING_CONTAINERS              ="SELECT LAST(fair_share_mb_with_pending_containers) WHERE serviceType=\"YARN\""                           
    FAIR_SHARE_VCORES                                  ="SELECT LAST(fair_share_vcores) WHERE serviceType=\"YARN\""                                               
    FAIR_SHARE_VCORES_CUMULATIVE                       ="SELECT LAST(fair_share_vcores_cumulative) WHERE serviceType=\"YARN\""                                    
    FAIR_SHARE_VCORES_WITH_PENDING_CONTAINERS          ="SELECT LAST(fair_share_vcores_with_pending_containers) WHERE serviceType=\"YARN\""                       
    HEALTH_BAD_RATE                                    ="SELECT LAST(health_bad_rate) WHERE serviceType=\"YARN\""                                                 
    HEALTH_CONCERNING_RATE                             ="SELECT LAST(health_concerning_rate) WHERE serviceType=\"YARN\""                                          
    HEALTH_DISABLED_RATE                               ="SELECT LAST(health_disabled_rate) WHERE serviceType=\"YARN\""                                            
    HEALTH_GOOD_RATE                                   ="SELECT LAST(health_good_rate) WHERE serviceType=\"YARN\""                                                
    HEALTH_UNKNOWN_RATE                                ="SELECT LAST(health_unknown_rate) WHERE serviceType=\"YARN\""                                             
    MAX_SHARE_MB                                       ="SELECT LAST(max_share_mb) WHERE serviceType=\"YARN\""                                                    
    MAX_SHARE_MB_CUMULATIVE                            ="SELECT LAST(max_share_mb_cumulative) WHERE serviceType=\"YARN\""                                         
    MAX_SHARE_VCORES                                   ="SELECT LAST(max_share_vcores) WHERE serviceType=\"YARN\""                                                
    MAX_SHARE_VCORES_CUMULATIVE                        ="SELECT LAST(max_share_vcores_cumulative) WHERE serviceType=\"YARN\""                                     
    MIN_SHARE_MB                                       ="SELECT LAST(min_share_mb) WHERE serviceType=\"YARN\""                                                    
    MIN_SHARE_MB_CUMULATIVE                            ="SELECT LAST(min_share_mb_cumulative) WHERE serviceType=\"YARN\""                                         
    MIN_SHARE_VCORES                                   ="SELECT LAST(min_share_vcores) WHERE serviceType=\"YARN\""                                                
    MIN_SHARE_VCORES_CUMULATIVE                        ="SELECT LAST(min_share_vcores_cumulative) WHERE serviceType=\"YARN\""                                     
    PENDING_CONTAINERS                                 ="SELECT LAST(pending_containers) WHERE serviceType=\"YARN\""                                              
    PENDING_CONTAINERS_CUMULATIVE                      ="SELECT LAST(pending_containers_cumulative) WHERE serviceType=\"YARN\""                                   
    PENDING_MEMORY_MB                                  ="SELECT LAST(pending_memory_mb) WHERE serviceType=\"YARN\""                                               
    PENDING_MEMORY_MB_CUMULATIVE                       ="SELECT LAST(pending_memory_mb_cumulative) WHERE serviceType=\"YARN\""                                    
    PENDING_VCORES                                     ="SELECT LAST(pending_vcores) WHERE serviceType=\"YARN\""                                                  
    PENDING_VCORES_CUMULATIVE                          ="SELECT LAST(pending_vcores_cumulative) WHERE serviceType=\"YARN\""                                       
    QUERIES_INGESTED_RATE                              ="SELECT LAST(queries_ingested_rate) WHERE serviceType=\"YARN\""                                           
    REPORTS_CONTAINERS_ALLOCATED_MEMORY                ="SELECT LAST(yarn_reports_containers_allocated_memory) WHERE serviceType=\"YARN\""                        
    REPORTS_CONTAINERS_ALLOCATED_VCORES                ="SELECT LAST(yarn_reports_containers_allocated_vcores) WHERE serviceType=\"YARN\""                        
    REPORTS_CONTAINERS_USED_CPU                        ="SELECT LAST(yarn_reports_containers_used_cpu) WHERE serviceType=\"YARN\""                                
    REPORTS_CONTAINERS_USED_MEMORY                     ="SELECT LAST(yarn_reports_containers_used_memory) WHERE serviceType=\"YARN\""                             
    REPORTS_CONTAINERS_USED_VCORES                     ="SELECT LAST(yarn_reports_containers_used_vcores) WHERE serviceType=\"YARN\""                             
    REPORTS_USAGE_AGGREGATION_DURATION                 ="SELECT LAST(yarn_reports_usage_aggregation_duration) WHERE serviceType=\"YARN\""                         
    REPORTS_USAGE_APPS_WITH_METADATA                   ="SELECT LAST(yarn_reports_usage_apps_with_metadata) WHERE serviceType=\"YARN\""                           
    REPORTS_USAGE_APPS_WITHOUT_METADATA                ="SELECT LAST(yarn_reports_usage_apps_without_metadata) WHERE serviceType=\"YARN\""                        
    RESERVED_CONTAINERS                                ="SELECT LAST(reserved_containers) WHERE serviceType=\"YARN\""                                             
    RESERVED_CONTAINERS_CUMULATIVE                     ="SELECT LAST(reserved_containers_cumulative) WHERE serviceType=\"YARN\""                                  
    RESERVED_MEMORY_MB                                 ="SELECT LAST(reserved_memory_mb) WHERE serviceType=\"YARN\""                                              
    RESERVED_MEMORY_MB_CUMULATIVE                      ="SELECT LAST(reserved_memory_mb_cumulative) WHERE serviceType=\"YARN\""                                   
    RESERVED_VCORES                                    ="SELECT LAST(reserved_vcores) WHERE serviceType=\"YARN\""                                                 
    RESERVED_VCORES_CUMULATIVE                         ="SELECT LAST(reserved_vcores_cumulative) WHERE serviceType=\"YARN\""                                      
    STEADY_FAIR_SHARE_MB                               ="SELECT LAST(steady_fair_share_mb) WHERE serviceType=\"YARN\""                                            
    STEADY_FAIR_SHARE_MB_CUMULATIVE                    ="SELECT LAST(steady_fair_share_mb_cumulative) WHERE serviceType=\"YARN\""                                 
    STEADY_FAIR_SHARE_MB_WITH_PENDING_CONTAINERS       ="SELECT LAST(steady_fair_share_mb_with_pending_containers) WHERE serviceType=\"YARN\""                    
    STEADY_FAIR_SHARE_VCORES                           ="SELECT LAST(steady_fair_share_vcores) WHERE serviceType=\"YARN\""                                        
    STEADY_FAIR_SHARE_VCORES_CUMULATIVE                ="SELECT LAST(steady_fair_share_vcores_cumulative) WHERE serviceType=\"YARN\""                             
    STEADY_FAIR_SHARE_VCORES_WITH_PENDING_CONTAINERS   ="SELECT LAST(steady_fair_share_vcores_with_pending_containers) WHERE serviceType=\"YARN\""                
  )
  
  
  
  
  /* ======================================================================
   * Global variables
   * ====================================================================== */
  var (
    active_applications                                  =create_yarn_metric_struct("active_applications", "Number of YARN applications in this pool with unsatisfied resource requests.	Applications	CDH 5, CDH 6")
    active_applications_cumulative                       =create_yarn_metric_struct("active_applications_cumulative", "Number of YARN applications in this pool with unsatisfied resource requests. Includes this pool and any children.	Applications	CDH 5, CDH 6")
    aggregate_containers_allocated_cumulative_rate       =create_yarn_metric_struct("aggregate_containers_allocated_cumulative_rate", "Aggregate Containers Allocated (Cumulative)	containers per second	CDH 5, CDH 6")
    aggregate_containers_allocated_rate                  =create_yarn_metric_struct("aggregate_containers_allocated_rate", "Aggregate Containers Allocated	containers per second	CDH 5, CDH 6")
    aggregate_containers_released_cumulative_rate        =create_yarn_metric_struct("aggregate_containers_released_cumulative_rate", "Aggregate Containers Released (Cumulative)	containers per second	CDH 5, CDH 6")
    aggregate_containers_released_rate                   =create_yarn_metric_struct("aggregate_containers_released_rate", "Aggregate Containers Released	containers per second	CDH 5, CDH 6")
    alerts_rate                                          =create_yarn_metric_struct("alerts_rate", "The number of alerts.	events per second	CDH 5, CDH 6")
    allocated_containers                                 =create_yarn_metric_struct("allocated_containers", "Number of containers allocated in this pool.	containers	CDH 5, CDH 6")
    allocated_containers_cumulative                      =create_yarn_metric_struct("allocated_containers_cumulative", "Number of containers allocated in this pool. Includes this pool and its children.	containers	CDH 5, CDH 6")
    allocated_memory_mb                                  =create_yarn_metric_struct("allocated_memory_mb", "Allocated Memory	MB	CDH 5, CDH 6")
    allocated_memory_mb_cumulative                       =create_yarn_metric_struct("allocated_memory_mb_cumulative", "Allocated Memory (Cumulative)	MB	CDH 5, CDH 6")
    allocated_memory_mb_with_pending_containers          =create_yarn_metric_struct("allocated_memory_mb_with_pending_containers", "Allocated memory when pending containers is more than zero.	MB	CDH 5, CDH 6")
    allocated_vcores                                     =create_yarn_metric_struct("allocated_vcores", "Number of vcores allocated in this pool	VCores	CDH 5, CDH 6")
    allocated_vcores_cumulative                          =create_yarn_metric_struct("allocated_vcores_cumulative", "Number of vcores allocated in this pool. Includes this pool and its children.	VCores	CDH 5, CDH 6")
    allocated_vcores_with_pending_containers             =create_yarn_metric_struct("allocated_vcores_with_pending_containers", "Allocated vcores when pending containers is more than zero.	VCores	CDH 5, CDH 6")
    application_adl_bytes_read_rate                      =create_yarn_metric_struct("application_adl_bytes_read_rate", "ADL bytes read. Called 'adl_bytes_read' in searches.	bytes per second	CDH 5, CDH 6")
    application_adl_bytes_written_rate                   =create_yarn_metric_struct("application_adl_bytes_written_rate", "ADL bytes written. Called 'adl_bytes_written' in searches.	bytes per second	CDH 5, CDH 6")
    application_application_duration_rate                =create_yarn_metric_struct("application_application_duration_rate", "How long YARN took to execute this application. Called 'application_duration' in searches.	ms per second	CDH 5, CDH 6")
    application_cm_cpu_milliseconds_rate                 =create_yarn_metric_struct("application_cm_cpu_milliseconds_rate", "yarn.analysis.cm_cpu_milliseconds.description	ms per second	CDH 5, CDH 6")
    application_cpu_milliseconds_rate                    =create_yarn_metric_struct("application_cpu_milliseconds_rate", "CPU time. Called 'cpu_milliseconds' in searches.	ms per second	CDH 5, CDH 6")
    application_file_bytes_read_rate                     =create_yarn_metric_struct("application_file_bytes_read_rate", "File bytes read. Called 'file_bytes_read' in searches.	bytes per second	CDH 5, CDH 6")
    application_file_bytes_written_rate                  =create_yarn_metric_struct("application_file_bytes_written_rate", "File bytes written. Called 'file_bytes_written' in searches.	bytes per second	CDH 5, CDH 6")
    application_hdfs_bytes_read_rate                     =create_yarn_metric_struct("application_hdfs_bytes_read_rate", "HDFS bytes read. Called 'hdfs_bytes_read' in searches.	bytes per second	CDH 5, CDH 6")
    application_hdfs_bytes_written_rate                  =create_yarn_metric_struct("application_hdfs_bytes_written_rate", "HDFS bytes written. Called 'hdfs_bytes_written' in searches.	bytes per second	CDH 5, CDH 6")
    application_maps_rate                                =create_yarn_metric_struct("application_maps_rate", "The number of Map tasks in this MapReduce job. Called 'maps_total' in searches.	items per second	CDH 5, CDH 6")
    application_mb_millis_maps_rate                      =create_yarn_metric_struct("application_mb_millis_maps_rate", "Map memory allocation. Called 'mb_millis_maps' in searches.	items per second	CDH 5, CDH 6")
    application_mb_millis_reduces_rate                   =create_yarn_metric_struct("application_mb_millis_reduces_rate", "Reduce memory allocation. Called 'mb_millis_reduces' in searches.	items per second	CDH 5, CDH 6")
    application_reduces_rate                             =create_yarn_metric_struct("application_reduces_rate", "The number of reduce tasks in this MapReduce job. Called 'reduces_total' in searches.	items per second	CDH 5, CDH 6")
    application_s3a_bytes_read_rate                      =create_yarn_metric_struct("application_s3a_bytes_read_rate", "S3A bytes read. Called 's3a_bytes_read' in searches.	bytes per second	CDH 5, CDH 6")
    application_s3a_bytes_written_rate                   =create_yarn_metric_struct("application_s3a_bytes_written_rate", "S3A bytes written. Called 's3a_bytes_written' in searches.	bytes per second	CDH 5, CDH 6")
    application_vcores_millis_maps_rate                  =create_yarn_metric_struct("application_vcores_millis_maps_rate", "Map CPU allocation. Called 'vcores_millis_maps' in searches.	items per second	CDH 5, CDH 6")
    application_vcores_millis_reduces_rate               =create_yarn_metric_struct("application_vcores_millis_reduces_rate", "Reduce CPU allocation. Called 'vcores_millis_reduces' in searches.	items per second	CDH 5, CDH 6")
    apps_completed_cumulative_rate                       =create_yarn_metric_struct("apps_completed_cumulative_rate", "Number of completed YARN applications in this pool. Includes this pool and any children.	Applications per second	CDH 5, CDH 6")
    apps_completed_rate                                  =create_yarn_metric_struct("apps_completed_rate", "Number of completed YARN applications in this pool.	Applications per second	CDH 5, CDH 6")
    apps_failed_cumulative_rate                          =create_yarn_metric_struct("apps_failed_cumulative_rate", "Number of failed YARN applications in this pool. Includes this pool and any children.	Applications per second	CDH 5, CDH 6")
    apps_failed_rate                                     =create_yarn_metric_struct("apps_failed_rate", "Number of failed YARN applications in this pool.	Applications per second	CDH 5, CDH 6")
    apps_ingested_rate                                   =create_yarn_metric_struct("apps_ingested_rate", "YARN applications ingested by the Service Monitor	Applications per second	CDH 5, CDH 6")
    apps_killed_cumulative_rate                          =create_yarn_metric_struct("apps_killed_cumulative_rate", "Number of killed YARN applications in this pool. Includes this pool and any children.	Applications per second	CDH 5, CDH 6")
    apps_killed_rate                                     =create_yarn_metric_struct("apps_killed_rate", "Number of killed YARN applications in this pool.	Applications per second	CDH 5, CDH 6")
    apps_pending                                         =create_yarn_metric_struct("apps_pending", "Number of pending (queued) YARN applications in this pool.	Applications	CDH 5, CDH 6")
    apps_pending_cumulative                              =create_yarn_metric_struct("apps_pending_cumulative", "Number of pending (queued) YARN applications in this pool. Includes this pool and any children.	Applications	CDH 5, CDH 6")
    apps_running                                         =create_yarn_metric_struct("apps_running", "Number of running YARN applications in this pool.	Applications	CDH 5, CDH 6")
    apps_running_between_300to1440_mins                  =create_yarn_metric_struct("apps_running_between_300to1440_mins", "Number of running YARN applications in this pool for which elapsed time is between 300 and 1440 minutes.	Applications	CDH 5, CDH 6")
    apps_running_between_300to1440_mins_cumulative       =create_yarn_metric_struct("apps_running_between_300to1440_mins_cumulative", "Number of running YARN applications in this pool for which elapsed time is between 300 and 1440 minutes. Includes this pool and any children.	Applications	CDH 5, CDH 6")
    apps_running_between_60to300_mins                    =create_yarn_metric_struct("apps_running_between_60to300_mins", "Number of running YARN applications in this pool for which elapsed time is between 60 and 300 minutes.	Applications	CDH 5, CDH 6")
    apps_running_between_60to300_mins_cumulative         =create_yarn_metric_struct("apps_running_between_60to300_mins_cumulative", "Number of running YARN applications in this pool for which elapsed time is between 60 and 300 minutes. Includes this pool and any children.	Applications	CDH 5, CDH 6")
    apps_running_cumulative                              =create_yarn_metric_struct("apps_running_cumulative", "Number of running YARN applications in this pool. Includes this pool and any children.	Applications	CDH 5, CDH 6")
    apps_running_over_1440_mins                          =create_yarn_metric_struct("apps_running_over_1440_mins", "Number of running YARN applications in this pool for which elapsed time is more than 1440 minutes.	Applications	CDH 5, CDH 6")
    apps_running_over_1440_mins_cumulative               =create_yarn_metric_struct("apps_running_over_1440_mins_cumulative", "Number of running YARN applications in this pool for which elapsed time is more than 1440 minutes. Includes this pool and any children.	Applications	CDH 5, CDH 6")
    apps_running_within_60_mins                          =create_yarn_metric_struct("apps_running_within_60_mins", "Number of running YARN applications in this pool for which elapsed time is less than 60 minutes.	Applications	CDH 5, CDH 6")
    apps_running_within_60_mins_cumulative               =create_yarn_metric_struct("apps_running_within_60_mins_cumulative", "Number of running YARN applications in this pool for which elapsed time is less than 60 minutes. Includes this pool and any children.	Applications	CDH 5, CDH 6")
    apps_submitted_cumulative_rate                       =create_yarn_metric_struct("apps_submitted_cumulative_rate", "Number of submitted YARN applications in this pool. Includes this pool and any children.	Applications per second	CDH 5, CDH 6")
    apps_submitted_rate                                  =create_yarn_metric_struct("apps_submitted_rate", "Number of submitted YARN applications in this queue.	Applications per second	CDH 5, CDH 6")
    available_memory_mb                                  =create_yarn_metric_struct("available_memory_mb", "Memory not allocated to YARN containers.	MB	CDH 5, CDH 6")
    available_vcores                                     =create_yarn_metric_struct("available_vcores", "Available vcores that can be used for containers.	VCores	CDH 5, CDH 6")
    container_wait_ratio                                 =create_yarn_metric_struct("container_wait_ratio", "Percent of pending containers when pending containers is more than zero.	percent	CDH 5, CDH 6")
    events_critical_rate                                 =create_yarn_metric_struct("events_critical_rate", "The number of critical events.	events per second	CDH 5, CDH 6")
    events_important_rate                                =create_yarn_metric_struct("events_important_rate", "The number of important events.	events per second	CDH 5, CDH 6")
    events_informational_rate                            =create_yarn_metric_struct("events_informational_rate", "The number of informational events.	events per second	CDH 5, CDH 6")
    fair_share_mb                                        =create_yarn_metric_struct("fair_share_mb", "Fair share of memory in this pool.	MB	CDH 5, CDH 6")
    fair_share_mb_cumulative                             =create_yarn_metric_struct("fair_share_mb_cumulative", "Fair share of memory in this pool. Includes this pool and its children.	MB	CDH 5, CDH 6")
    fair_share_mb_with_pending_containers                =create_yarn_metric_struct("fair_share_mb_with_pending_containers", "Fair share of memory when pending containers is more than zero.	MB	CDH 5, CDH 6")
    fair_share_vcores                                    =create_yarn_metric_struct("fair_share_vcores", "Fair share of vcores in this pool.	VCores	CDH 5, CDH 6")
    fair_share_vcores_cumulative                         =create_yarn_metric_struct("fair_share_vcores_cumulative", "Fair share of vcores in this pool. Includes this pool and its children.	VCores	CDH 5, CDH 6")
    fair_share_vcores_with_pending_containers            =create_yarn_metric_struct("fair_share_vcores_with_pending_containers", "Fair share of vcores when pending containers is more than zero.	VCores	CDH 5, CDH 6")
    health_bad_rate                                      =create_yarn_metric_struct("health_bad_rate", "Percentage of Time with Bad Health	seconds per second	CDH 5, CDH 6")
    health_concerning_rate                               =create_yarn_metric_struct("health_concerning_rate", "Percentage of Time with Concerning Health	seconds per second	CDH 5, CDH 6")
    health_disabled_rate                                 =create_yarn_metric_struct("health_disabled_rate", "Percentage of Time with Disabled Health	seconds per second	CDH 5, CDH 6")
    health_good_rate                                     =create_yarn_metric_struct("health_good_rate", "Percentage of Time with Good Health	seconds per second	CDH 5, CDH 6")
    health_unknown_rate                                  =create_yarn_metric_struct("health_unknown_rate", "Percentage of Time with Unknown Health	seconds per second	CDH 5, CDH 6")
    max_share_mb                                         =create_yarn_metric_struct("max_share_mb", "Maximum share of memory configured for this pool.	MB	CDH 5, CDH 6")
    max_share_mb_cumulative                              =create_yarn_metric_struct("max_share_mb_cumulative", "Maximum share of memory configured for this pool. Includes this pool and its children.	MB	CDH 5, CDH 6")
    max_share_vcores                                     =create_yarn_metric_struct("max_share_vcores", "Maximum share of vcores configured for this pool.	VCores	CDH 5, CDH 6")
    max_share_vcores_cumulative                          =create_yarn_metric_struct("max_share_vcores_cumulative", "Maximum share of vcores configured for this pool. Includes this pool and its children.	VCores	CDH 5, CDH 6")
    min_share_mb                                         =create_yarn_metric_struct("min_share_mb", "Minimum share of memory configured for this pool.	MB	CDH 5, CDH 6")
    min_share_mb_cumulative                              =create_yarn_metric_struct("min_share_mb_cumulative", "Minimum share of memory configured for this pool. Includes this pool and its children.	MB	CDH 5, CDH 6")
    min_share_vcores                                     =create_yarn_metric_struct("min_share_vcores", "Minimum share of vcores configured for this pool.	VCores	CDH 5, CDH 6")
    min_share_vcores_cumulative                          =create_yarn_metric_struct("min_share_vcores_cumulative", "Minimum share of vcores configured for this pool. Includes this pool and its children.	VCores	CDH 5, CDH 6")
    pending_containers                                   =create_yarn_metric_struct("pending_containers", "Number of pending (queued) containers in this pool.	containers	CDH 5, CDH 6")
    pending_containers_cumulative                        =create_yarn_metric_struct("pending_containers_cumulative", "Number of pending (queued) containers in this pool. Includes this pool and any children.	containers	CDH 5, CDH 6")
    pending_memory_mb                                    =create_yarn_metric_struct("pending_memory_mb", "Sum of memory currently requested but not allocated for containers in this pool.	MB	CDH 5, CDH 6")
    pending_memory_mb_cumulative                         =create_yarn_metric_struct("pending_memory_mb_cumulative", "Sum of memory currently requested but not allocated for containers in this pool. Includes this pool and any children.	MB	CDH 5, CDH 6")
    pending_vcores                                       =create_yarn_metric_struct("pending_vcores", "Number of vcores requested but not yet allocated for this pool.	VCores	CDH 5, CDH 6")
    pending_vcores_cumulative                            =create_yarn_metric_struct("pending_vcores_cumulative", "Number of vcores requested but not yet allocated for this pool. Includes this pool and its children.	VCores	CDH 5, CDH 6")
    queries_ingested_rate                                =create_yarn_metric_struct("queries_ingested_rate", "Impala queries ingested by the Service Monitor	queries per second	CDH 5, CDH 6")
    reports_containers_allocated_memory                  =create_yarn_metric_struct("reports_containers_allocated_memory", "Memory allocated to YARN containers	MB seconds	CDH 5, CDH 6")
    reports_containers_allocated_vcores                  =create_yarn_metric_struct("reports_containers_allocated_vcores", "VCores allocated to YARN containers	VCore seconds	CDH 5, CDH 6")
    reports_containers_used_cpu                          =create_yarn_metric_struct("reports_containers_used_cpu", "CPU used by YARN containers	Percent seconds	CDH 5, CDH 6")
    reports_containers_used_memory                       =create_yarn_metric_struct("reports_containers_used_memory", "Memory used by YARN containers	MB seconds	CDH 5, CDH 6")
    reports_containers_used_vcores                       =create_yarn_metric_struct("reports_containers_used_vcores", "VCores used by YARN containers	VCore seconds	CDH 5, CDH 6")
    reports_usage_aggregation_duration                   =create_yarn_metric_struct("reports_usage_aggregation_duration", "The duration for generating YARN usage reporting metrics by aggregating YARN container usage metrics	ms	CDH 5, CDH 6")
    reports_usage_apps_with_metadata                     =create_yarn_metric_struct("reports_usage_apps_with_metadata", "YARN applications for which container usage was computed and also had metadata	Applications	CDH 5, CDH 6")
    reports_usage_apps_without_metadata                  =create_yarn_metric_struct("reports_usage_apps_without_metadata", "YARN applications for which container usage was computed but did not have metadata	Applications	CDH 5, CDH 6")
    reserved_containers                                  =create_yarn_metric_struct("reserved_containers", "Reserved containers for this pool.	containers	CDH 5, CDH 6")
    reserved_containers_cumulative                       =create_yarn_metric_struct("reserved_containers_cumulative", "Reserved containers for this pool. Includes this pool and any children.	containers	CDH 5, CDH 6")
    reserved_memory_mb                                   =create_yarn_metric_struct("reserved_memory_mb", "Reserved memory in this pool.	MB	CDH 5, CDH 6")
    reserved_memory_mb_cumulative                        =create_yarn_metric_struct("reserved_memory_mb_cumulative", "Reserved memory in this pool. Includes this pool and any children.	MB	CDH 5, CDH 6")
    reserved_vcores                                      =create_yarn_metric_struct("reserved_vcores", "Number of vcores set aside for this pool, but not used as a part of an allocation.	VCores	CDH 5, CDH 6")
    reserved_vcores_cumulative                           =create_yarn_metric_struct("reserved_vcores_cumulative", "Number of vcores set aside for this pool, but not used as a part of an allocation. Includes this pool and its children.	VCores	CDH 5, CDH 6")
    steady_fair_share_mb                                 =create_yarn_metric_struct("steady_fair_share_mb", "Steady fair share of memory in this pool.	MB	CDH 5, CDH 6")
    steady_fair_share_mb_cumulative                      =create_yarn_metric_struct("steady_fair_share_mb_cumulative", "Steady fair share of memory in this pool. Includes this pool and its children.	MB	CDH 5, CDH 6")
    steady_fair_share_mb_with_pending_containers         =create_yarn_metric_struct("steady_fair_share_mb_with_pending_containers", "Steady fair share of memory when pending containers is more than zero.	MB	CDH 5, CDH 6")
    steady_fair_share_vcores                             =create_yarn_metric_struct("steady_fair_share_vcores", "Steady fair share of vcores in this pool.	VCores	CDH 5, CDH 6")
    steady_fair_share_vcores_cumulative                  =create_yarn_metric_struct("steady_fair_share_vcores_cumulative", "Steady fair share of vcores in this pool. Includes this pool and its children.	VCores	CDH 5, CDH 6")
    steady_fair_share_vcores_with_pending_containers     =create_yarn_metric_struct("steady_fair_share_vcores_with_pending_containers", "Steady fair share of vcores when pending containers is more than zero.	VCores	CDH 5, CDH 6")
  )
  var yarn_query_variable_relationship = []relation {
    {ACTIVE_APPLICATIONS,                               *active_applications},                                                   
    {ACTIVE_APPLICATIONS_CUMULATIVE,                    *active_applications_cumulative},                                        
    {AGGREGATE_CONTAINERS_ALLOCATED_CUMULATIVE_RATE,    *aggregate_containers_allocated_cumulative_rate},                        
    {AGGREGATE_CONTAINERS_ALLOCATED_RATE,               *aggregate_containers_allocated_rate},                                   
    {AGGREGATE_CONTAINERS_RELEASED_CUMULATIVE_RATE,     *aggregate_containers_released_cumulative_rate},                         
    {AGGREGATE_CONTAINERS_RELEASED_RATE,                *aggregate_containers_released_rate},                                    
    {ALERTS_RATE,                                       *alerts_rate},                                                           
    {ALLOCATED_CONTAINERS,                              *allocated_containers},                                                  
    {ALLOCATED_CONTAINERS_CUMULATIVE,                   *allocated_containers_cumulative},                                       
    {ALLOCATED_MEMORY_MB,                               *allocated_memory_mb},                                                   
    {ALLOCATED_MEMORY_MB_CUMULATIVE,                    *allocated_memory_mb_cumulative},                                        
    {ALLOCATED_MEMORY_MB_WITH_PENDING_CONTAINERS,       *allocated_memory_mb_with_pending_containers},                           
    {ALLOCATED_VCORES,                                  *allocated_vcores},                                                      
    {ALLOCATED_VCORES_CUMULATIVE,                       *allocated_vcores_cumulative},                                           
    {ALLOCATED_VCORES_WITH_PENDING_CONTAINERS,          *allocated_vcores_with_pending_containers},                              
    {APPLICATION_ADL_BYTES_READ_RATE,                   *application_adl_bytes_read_rate},                                  
    {APPLICATION_ADL_BYTES_WRITTEN_RATE,                *application_adl_bytes_written_rate},                               
    {APPLICATION_APPLICATION_DURATION_RATE,             *application_application_duration_rate},                            
    {APPLICATION_CM_CPU_MILLISECONDS_RATE,              *application_cm_cpu_milliseconds_rate},                             
    {APPLICATION_CPU_MILLISECONDS_RATE,                 *application_cpu_milliseconds_rate},                                
    {APPLICATION_FILE_BYTES_READ_RATE,                  *application_file_bytes_read_rate},                                 
    {APPLICATION_FILE_BYTES_WRITTEN_RATE,               *application_file_bytes_written_rate},                              
    {APPLICATION_HDFS_BYTES_READ_RATE,                  *application_hdfs_bytes_read_rate},                                 
    {APPLICATION_HDFS_BYTES_WRITTEN_RATE,               *application_hdfs_bytes_written_rate},                              
    {APPLICATION_MAPS_RATE,                             *application_maps_rate},                                            
    {APPLICATION_MB_MILLIS_MAPS_RATE,                   *application_mb_millis_maps_rate},                                  
    {APPLICATION_MB_MILLIS_REDUCES_RATE,                *application_mb_millis_reduces_rate},                               
    {APPLICATION_REDUCES_RATE,                          *application_reduces_rate},                                         
    {APPLICATION_S3A_BYTES_READ_RATE,                   *application_s3a_bytes_read_rate},                                  
    {APPLICATION_S3A_BYTES_WRITTEN_RATE,                *application_s3a_bytes_written_rate},                               
    {APPLICATION_VCORES_MILLIS_MAPS_RATE,               *application_vcores_millis_maps_rate},                              
    {APPLICATION_VCORES_MILLIS_REDUCES_RATE,            *application_vcores_millis_reduces_rate},                           
    {APPS_COMPLETED_CUMULATIVE_RATE,                    *apps_completed_cumulative_rate},                                        
    {APPS_COMPLETED_RATE,                               *apps_completed_rate},                                                   
    {APPS_FAILED_CUMULATIVE_RATE,                       *apps_failed_cumulative_rate},                                           
    {APPS_FAILED_RATE,                                  *apps_failed_rate},                                                      
    {APPS_INGESTED_RATE,                                *apps_ingested_rate},                                                    
    {APPS_KILLED_CUMULATIVE_RATE,                       *apps_killed_cumulative_rate},                                           
    {APPS_KILLED_RATE,                                  *apps_killed_rate},                                                      
    {APPS_PENDING,                                      *apps_pending},                                                          
    {APPS_PENDING_CUMULATIVE,                           *apps_pending_cumulative},                                               
    {APPS_RUNNING,                                      *apps_running},                                                          
    {APPS_RUNNING_BETWEEN_300TO1440_MINS,               *apps_running_between_300to1440_mins},                                   
    {APPS_RUNNING_BETWEEN_300TO1440_MINS_CUMULATIVE,    *apps_running_between_300to1440_mins_cumulative},                        
    {APPS_RUNNING_BETWEEN_60TO300_MINS,                 *apps_running_between_60to300_mins},                                     
    {APPS_RUNNING_BETWEEN_60TO300_MINS_CUMULATIVE,      *apps_running_between_60to300_mins_cumulative},                          
    {APPS_RUNNING_CUMULATIVE,                           *apps_running_cumulative},                                               
    {APPS_RUNNING_OVER_1440_MINS,                       *apps_running_over_1440_mins},                                           
    {APPS_RUNNING_OVER_1440_MINS_CUMULATIVE,            *apps_running_over_1440_mins_cumulative},                                
    {APPS_RUNNING_WITHIN_60_MINS,                       *apps_running_within_60_mins},                                           
    {APPS_RUNNING_WITHIN_60_MINS_CUMULATIVE,            *apps_running_within_60_mins_cumulative},                                
    {APPS_SUBMITTED_CUMULATIVE_RATE,                    *apps_submitted_cumulative_rate},                                        
    {APPS_SUBMITTED_RATE,                               *apps_submitted_rate},                                                   
    {AVAILABLE_MEMORY_MB,                               *available_memory_mb},                                                   
    {AVAILABLE_VCORES,                                  *available_vcores},                                                      
    {CONTAINER_WAIT_RATIO,                              *container_wait_ratio},                                                  
    {EVENTS_CRITICAL_RATE,                              *events_critical_rate},                                                  
    {EVENTS_IMPORTANT_RATE,                             *events_important_rate},                                                 
    {EVENTS_INFORMATIONAL_RATE,                         *events_informational_rate},                                             
    {FAIR_SHARE_MB,                                     *fair_share_mb},                                                         
    {FAIR_SHARE_MB_CUMULATIVE,                          *fair_share_mb_cumulative},                                              
    {FAIR_SHARE_MB_WITH_PENDING_CONTAINERS,             *fair_share_mb_with_pending_containers},                                 
    {FAIR_SHARE_VCORES,                                 *fair_share_vcores},                                                     
    {FAIR_SHARE_VCORES_CUMULATIVE,                      *fair_share_vcores_cumulative},                                          
    {FAIR_SHARE_VCORES_WITH_PENDING_CONTAINERS,         *fair_share_vcores_with_pending_containers},                             
    {HEALTH_BAD_RATE,                                   *health_bad_rate},                                                       
    {HEALTH_CONCERNING_RATE,                            *health_concerning_rate},                                                
    {HEALTH_DISABLED_RATE,                              *health_disabled_rate},                                                  
    {HEALTH_GOOD_RATE,                                  *health_good_rate},                                                      
    {HEALTH_UNKNOWN_RATE,                               *health_unknown_rate},                                                   
    {MAX_SHARE_MB,                                      *max_share_mb},                                                          
    {MAX_SHARE_MB_CUMULATIVE,                           *max_share_mb_cumulative},                                               
    {MAX_SHARE_VCORES,                                  *max_share_vcores},                                                      
    {MAX_SHARE_VCORES_CUMULATIVE,                       *max_share_vcores_cumulative},                                           
    {MIN_SHARE_MB,                                      *min_share_mb},                                                          
    {MIN_SHARE_MB_CUMULATIVE,                           *min_share_mb_cumulative},                                               
    {MIN_SHARE_VCORES,                                  *min_share_vcores},                                                      
    {MIN_SHARE_VCORES_CUMULATIVE,                       *min_share_vcores_cumulative},                                           
    {PENDING_CONTAINERS,                                *pending_containers},                                                    
    {PENDING_CONTAINERS_CUMULATIVE,                     *pending_containers_cumulative},                                         
    {PENDING_MEMORY_MB,                                 *pending_memory_mb},                                                     
    {PENDING_MEMORY_MB_CUMULATIVE,                      *pending_memory_mb_cumulative},                                          
    {PENDING_VCORES,                                    *pending_vcores},                                                        
    {PENDING_VCORES_CUMULATIVE,                         *pending_vcores_cumulative},                                             
    {QUERIES_INGESTED_RATE,                             *queries_ingested_rate},                                                 
    {REPORTS_CONTAINERS_ALLOCATED_MEMORY,               *reports_containers_allocated_memory},                              
    {REPORTS_CONTAINERS_ALLOCATED_VCORES,               *reports_containers_allocated_vcores},                              
    {REPORTS_CONTAINERS_USED_CPU,                       *reports_containers_used_cpu},                                      
    {REPORTS_CONTAINERS_USED_MEMORY,                    *reports_containers_used_memory},                                   
    {REPORTS_CONTAINERS_USED_VCORES,                    *reports_containers_used_vcores},                                   
    {REPORTS_USAGE_AGGREGATION_DURATION,                *reports_usage_aggregation_duration},                               
    {REPORTS_USAGE_APPS_WITH_METADATA,                  *reports_usage_apps_with_metadata},                                 
    {REPORTS_USAGE_APPS_WITHOUT_METADATA,               *reports_usage_apps_without_metadata},                              
    {RESERVED_CONTAINERS,                               *reserved_containers},                                                   
    {RESERVED_CONTAINERS_CUMULATIVE,                    *reserved_containers_cumulative},                                        
    {RESERVED_MEMORY_MB,                                *reserved_memory_mb},                                                    
    {RESERVED_MEMORY_MB_CUMULATIVE,                     *reserved_memory_mb_cumulative},                                         
    {RESERVED_VCORES,                                   *reserved_vcores},                                                       
    {RESERVED_VCORES_CUMULATIVE,                        *reserved_vcores_cumulative},                                            
    {STEADY_FAIR_SHARE_MB,                              *steady_fair_share_mb},                                                  
    {STEADY_FAIR_SHARE_MB_CUMULATIVE,                   *steady_fair_share_mb_cumulative},                                       
    {STEADY_FAIR_SHARE_MB_WITH_PENDING_CONTAINERS,      *steady_fair_share_mb_with_pending_containers},                          
    {STEADY_FAIR_SHARE_VCORES,                          *steady_fair_share_vcores},                                              
    {STEADY_FAIR_SHARE_VCORES_CUMULATIVE,               *steady_fair_share_vcores_cumulative},                                   
    {STEADY_FAIR_SHARE_VCORES_WITH_PENDING_CONTAINERS,  *steady_fair_share_vcores_with_pending_containers},                      
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
 