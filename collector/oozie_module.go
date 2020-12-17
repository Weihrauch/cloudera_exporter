/*
 *
 * title           :collector/oozie_module.go
 * description     :Submodule Collector for the Cluster OOZIE metrics
 * author		       :Ntumba Phinées
 * date            :2020/12/09
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
    pool "keedio/cloudera_exporter/pool"
  )
  
  
  
  
  /* ======================================================================
   * Data Structs
   * ====================================================================== */
  // None
  
  
  
  /* ======================================================================
   * Constants with the Host module TSquery sentences
   * ====================================================================== */
  const OOZIE_SCRAPER_NAME = "oozie"
  var (
    // Agent Queries
OOZIE_ALERTS_RATE                                                                                              ="SELECT LAST(alerts_rate) WHERE serviceType=\"OOZIE\""
OOZIE_CALLABLEQUEUE_ITEMS_EXECUTED_RATE                                                                        ="SELECT LAST(callablequeue_items_executed_rate) WHERE serviceType=\"OOZIE\""
OOZIE_CALLABLEQUEUE_ITEMS_FAILED_RATE                                                                          ="SELECT LAST(callablequeue_items_failed_rate) WHERE serviceType=\"OOZIE\""
OOZIE_CALLABLEQUEUE_ITEMS_QUEUED_RATE                                                                          ="SELECT LAST(callablequeue_items_queued_rate) WHERE serviceType=\"OOZIE\""
OOZIE_CGROUP_CPU_SYSTEM_RATE                                                                                   ="SELECT LAST(cgroup_cpu_system_rate) WHERE serviceType=\"OOZIE\""
OOZIE_CGROUP_CPU_USER_RATE                                                                                     ="SELECT LAST(cgroup_cpu_user_rate) WHERE serviceType=\"OOZIE\""
OOZIE_CGROUP_MEM_PAGE_CACHE                                                                                    ="SELECT LAST(cgroup_mem_page_cache) WHERE serviceType=\"OOZIE\""
OOZIE_CGROUP_MEM_RSS                                                                                           ="SELECT LAST(cgroup_mem_rss) WHERE serviceType=\"OOZIE\""
OOZIE_CGROUP_MEM_SWAP                                                                                          ="SELECT LAST(cgroup_mem_swap) WHERE serviceType=\"OOZIE\""
OOZIE_CGROUP_READ_BYTES_RATE                                                                                   ="SELECT LAST(cgroup_read_bytes_rate) WHERE serviceType=\"OOZIE\""
OOZIE_CGROUP_READ_IOS_RATE                                                                                     ="SELECT LAST(cgroup_read_ios_rate) WHERE serviceType=\"OOZIE\""
OOZIE_CGROUP_WRITE_BYTES_RATE                                                                                  ="SELECT LAST(cgroup_write_bytes_rate) WHERE serviceType=\"OOZIE\""
OOZIE_CGROUP_WRITE_IOS_RATE                                                                                    ="SELECT LAST(cgroup_write_ios_rate) WHERE serviceType=\"OOZIE\""
OOZIE_CPU_SYSTEM_RATE                                                                                          ="SELECT LAST(cpu_system_rate) WHERE serviceType=\"OOZIE\""
OOZIE_CPU_USER_RATE                                                                                            ="SELECT LAST(cpu_user_rate) WHERE serviceType=\"OOZIE\""
OOZIE_EVENTS_CRITICAL_RATE                                                                                     ="SELECT LAST(events_critical_rate) WHERE serviceType=\"OOZIE\""
OOZIE_EVENTS_IMPORTANT_RATE                                                                                    ="SELECT LAST(events_important_rate) WHERE serviceType=\"OOZIE\""
OOZIE_EVENTS_INFORMATIONAL_RATE                                                                                ="SELECT LAST(events_informational_rate) WHERE serviceType=\"OOZIE\""
OOZIE_FD_OPEN                                                                                                  ="SELECT LAST(fd_open) WHERE serviceType=\"OOZIE\""
OOZIE_JVM_TOTAL_MEMORY                                                                                         ="SELECT LAST(jvm_total_memory) WHERE serviceType=\"OOZIE\""
OOZIE_MEM_RSS                                                                                                  ="SELECT LAST(mem_rss) WHERE serviceType=\"OOZIE\""
OOZIE_MEM_SWAP                                                                                                 ="SELECT LAST(mem_swap) WHERE serviceType=\"OOZIE\""
OOZIE_MEM_VIRTUAL                                                                                              ="SELECT LAST(mem_virtual) WHERE serviceType=\"OOZIE\""
OOZIE_OOM_EXITS_RATE                                                                                           ="SELECT LAST(oom_exits_rate) WHERE serviceType=\"OOZIE\""
OOZIE_BUNDLE_ACTION_QUERY_EXECUTOR_GET_BUNDLE_ACTIONS_FOR_BUNDLE_DURATION_TIMER_AVG                            ="SELECT LAST(oozie_bundle_action_query_executor_get_bundle_actions_for_bundle_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_BUNDLE_ACTION_QUERY_EXECUTOR_GET_BUNDLE_ACTIONS_FOR_BUNDLE_DURATION_TIMER_RATE                           ="SELECT LAST(oozie_bundle_action_query_executor_get_bundle_actions_for_bundle_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_BUNDLE_JOB_QUERY_EXECUTOR_GET_BUNDLE_JOB_DURATION_TIMER_AVG                                              ="SELECT LAST(oozie_bundle_job_query_executor_get_bundle_job_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_BUNDLE_JOB_QUERY_EXECUTOR_GET_BUNDLE_JOB_DURATION_TIMER_RATE                                             ="SELECT LAST(oozie_bundle_job_query_executor_get_bundle_job_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_CALLBACK_REST_CALL_DURATION_TIMER_AVG                                                                    ="SELECT LAST(oozie_callback_rest_call_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_CALLBACK_REST_CALL_DURATION_TIMER_RATE                                                                   ="SELECT LAST(oozie_callback_rest_call_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_COORD_ACTION_QUERY_EXECUTOR_GET_COORD_ACTION_DURATION_TIMER_AVG                                          ="SELECT LAST(oozie_coord_action_query_executor_get_coord_action_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_COORD_ACTION_QUERY_EXECUTOR_GET_COORD_ACTION_DURATION_TIMER_RATE                                         ="SELECT LAST(oozie_coord_action_query_executor_get_coord_action_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_COORD_ACTION_QUERY_EXECUTOR_GET_COORD_ACTIVE_ACTIONS_BY_JOBID_DURATION_TIMER_RATE                        ="SELECT LAST(oozie_coord_action_query_executor_get_coord_active_actions_by_jobid_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_COORD_ACTION_QUERY_EXECUTOR_GET_COORD_ACTIVE_ACTIONS_COUNT_BY_JOBID_DURATION_TIMER_AVG                   ="SELECT LAST(oozie_coord_action_query_executor_get_coord_active_actions_count_by_jobid_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOB_ACTION_READY_DURATION_TIMER_AVG                                   ="SELECT LAST(oozie_coord_job_query_executor_get_coord_job_action_ready_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOB_ACTION_READY_DURATION_TIMER_RATE                                  ="SELECT LAST(oozie_coord_job_query_executor_get_coord_job_action_ready_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOB_DURATION_TIMER_AVG                                                ="SELECT LAST(oozie_coord_job_query_executor_get_coord_job_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOB_DURATION_TIMER_RATE                                               ="SELECT LAST(oozie_coord_job_query_executor_get_coord_job_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOB_USER_APPNAME_DURATION_TIMER_AVG                                   ="SELECT LAST(oozie_coord_job_query_executor_get_coord_job_user_appname_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOB_USER_APPNAME_DURATION_TIMER_RATE                                  ="SELECT LAST(oozie_coord_job_query_executor_get_coord_job_user_appname_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOBS_CHANGED_DURATION_TIMER_AVG                                       ="SELECT LAST(oozie_coord_job_query_executor_get_coord_jobs_changed_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOBS_CHANGED_DURATION_TIMER_RATE                                      ="SELECT LAST(oozie_coord_job_query_executor_get_coord_jobs_changed_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_JVM_PAUSE_TIME_RATE                                                                                      ="SELECT LAST(oozie_jvm_pause_time_rate) WHERE serviceType=\"OOZIE\""
OOZIE_JVM_PAUSES_INFO_THRESHOLD_RATE                                                                           ="SELECT LAST(oozie_jvm_pauses_info_threshold_rate) WHERE serviceType=\"OOZIE\""
OOZIE_JVM_PAUSES_WARN_THRESHOLD_RATE                                                                           ="SELECT LAST(oozie_jvm_pauses_warn_threshold_rate) WHERE serviceType=\"OOZIE\""
OOZIE_MEMORY_HEAP_COMMITTED                                                                                    ="SELECT LAST(oozie_memory_heap_committed) WHERE serviceType=\"OOZIE\""
OOZIE_MEMORY_HEAP_INIT                                                                                         ="SELECT LAST(oozie_memory_heap_init) WHERE serviceType=\"OOZIE\""
OOZIE_MEMORY_HEAP_USED                                                                                         ="SELECT LAST(oozie_memory_heap_used) WHERE serviceType=\"OOZIE\""
OOZIE_MEMORY_NON_HEAP_COMMITTED                                                                                ="SELECT LAST(oozie_memory_non_heap_committed) WHERE serviceType=\"OOZIE\""
OOZIE_MEMORY_NON_HEAP_INIT                                                                                     ="SELECT LAST(oozie_memory_non_heap_init) WHERE serviceType=\"OOZIE\""
OOZIE_MEMORY_NON_HEAP_USED                                                                                     ="SELECT LAST(oozie_memory_non_heap_used) WHERE serviceType=\"OOZIE\""
OOZIE_MEMORY_TOTAL_COMMITTED                                                                                   ="SELECT LAST(oozie_memory_total_committed) WHERE serviceType=\"OOZIE\""
OOZIE_MEMORY_TOTAL_INIT                                                                                        ="SELECT LAST(oozie_memory_total_init) WHERE serviceType=\"OOZIE\""
OOZIE_MEMORY_TOTAL_USED                                                                                        ="SELECT LAST(oozie_memory_total_used) WHERE serviceType=\"OOZIE\""
OOZIE_V1JOB_REST_CALL_DURATION_TIMER_AVG                                                                       ="SELECT LAST(oozie_v1job_rest_call_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_V1JOB_REST_CALL_DURATION_TIMER_RATE                                                                      ="SELECT LAST(oozie_v1job_rest_call_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_V1JOBS_REST_CALL_DURATION_TIMER_AVG                                                                      ="SELECT LAST(oozie_v1jobs_rest_call_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_V1JOBS_REST_CALL_DURATION_TIMER_RATE                                                                     ="SELECT LAST(oozie_v1jobs_rest_call_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_V2JOB_REST_CALL_DURATION_TIMER_AVG                                                                       ="SELECT LAST(oozie_v2job_rest_call_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_V2JOB_REST_CALL_DURATION_TIMER_RATE                                                                      ="SELECT LAST(oozie_v2job_rest_call_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_VERSION_REST_CALL_DURATION_TIMER_AVG                                                                     ="SELECT LAST(oozie_version_rest_call_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_VERSION_REST_CALL_DURATION_TIMER_RATE                                                                    ="SELECT LAST(oozie_version_rest_call_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_WEBSERVICES_REQUESTS_RATE                                                                                ="SELECT LAST(oozie_webservices_requests_rate) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_ACTION_CHECK_DURATION_TIMER_AVG                                       ="SELECT LAST(oozie_workflow_action_query_executor_get_action_check_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_ACTION_CHECK_DURATION_TIMER_RATE                                      ="SELECT LAST(oozie_workflow_action_query_executor_get_action_check_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_ACTION_COMPLETED_DURATION_TIMER_AVG                                   ="SELECT LAST(oozie_workflow_action_query_executor_get_action_completed_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_ACTION_COMPLETED_DURATION_TIMER_RATE                                  ="SELECT LAST(oozie_workflow_action_query_executor_get_action_completed_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_ACTION_DURATION_TIMER_AVG                                             ="SELECT LAST(oozie_workflow_action_query_executor_get_action_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_ACTION_DURATION_TIMER_RATE                                            ="SELECT LAST(oozie_workflow_action_query_executor_get_action_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_PENDING_ACTIONS_DURATION_TIMER_AVG                                    ="SELECT LAST(oozie_workflow_action_query_executor_get_pending_actions_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_PENDING_ACTIONS_DURATION_TIMER_RATE                                   ="SELECT LAST(oozie_workflow_action_query_executor_get_pending_actions_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_RUNNING_ACTIONS_DURATION_TIMER_AVG                                    ="SELECT LAST(oozie_workflow_action_query_executor_get_running_actions_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_RUNNING_ACTIONS_DURATION_TIMER_RATE                                   ="SELECT LAST(oozie_workflow_action_query_executor_get_running_actions_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_JOB_QUERY_EXECUTOR_GET_WORKFLOW_DURATION_TIMER_AVG                                              ="SELECT LAST(oozie_workflow_job_query_executor_get_workflow_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_JOB_QUERY_EXECUTOR_GET_WORKFLOW_DURATION_TIMER_RATE                                             ="SELECT LAST(oozie_workflow_job_query_executor_get_workflow_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_JOB_QUERY_EXECUTOR_GET_WORKFLOW_START_END_TIME_DURATION_TIMER_AVG                               ="SELECT LAST(oozie_workflow_job_query_executor_get_workflow_start_end_time_duration_timer_avg) WHERE serviceType=\"OOZIE\""
OOZIE_WORKFLOW_JOB_QUERY_EXECUTOR_GET_WORKFLOW_START_END_TIME_DURATION_TIMER_RATE                              ="SELECT LAST(oozie_workflow_job_query_executor_get_workflow_start_end_time_duration_timer_rate) WHERE serviceType=\"OOZIE\""
OOZIE_PAUSE_TIME_RATE                                                                                          ="SELECT LAST(pause_time_rate) WHERE serviceType=\"OOZIE\""
OOZIE_PAUSES_RATE                                                                                              ="SELECT LAST(pauses_rate) WHERE serviceType=\"OOZIE\""
OOZIE_READ_BYTES_RATE                                                                                          ="SELECT LAST(read_bytes_rate) WHERE serviceType=\"OOZIE\""
OOZIE_UNEXPECTED_EXITS_RATE                                                                                    ="SELECT LAST(unexpected_exits_rate) WHERE serviceType=\"OOZIE\""
OOZIE_UPTIME                                                                                                   ="SELECT LAST(uptime) WHERE serviceType=\"OOZIE\""
OOZIE_WEB_METRICS_COLLECTION_DURATION                                                                          ="SELECT LAST(web_metrics_collection_duration) WHERE serviceType=\"OOZIE\""
OOZIE_WRITE_BYTES_RATE                                                                                         ="SELECT LAST(write_bytes_rate) WHERE serviceType=\"OOZIE\""
)
  
  
  
  
  /* ======================================================================
   * Global variables
   * ====================================================================== */
  var (
    oozie_alerts_rate                                                                                                       =create_oozie_metric_struct("oozie_alerts_rate", "The number of alerts.	events per second")
    oozie_callablequeue_items_executed_rate                                                                                 =create_oozie_metric_struct("oozie_callablequeue_items_executed_rate", "Callable Queue Executed Items	items per second")
    oozie_callablequeue_items_failed_rate                                                                                   =create_oozie_metric_struct("oozie_callablequeue_items_failed_rate", "Callable Queue Failed Items	items per second")
    oozie_callablequeue_items_queued_rate                                                                                   =create_oozie_metric_struct("oozie_callablequeue_items_queued_rate", "Callable Queue Queued Items	items per second")
    oozie_cgroup_cpu_system_rate                                                                                            =create_oozie_metric_struct("oozie_cgroup_cpu_system_rate", "CPU usage of the role's cgroup	seconds per second")
    oozie_cgroup_cpu_user_rate                                                                                              =create_oozie_metric_struct("oozie_cgroup_cpu_user_rate", "User Space CPU usage of the role's cgroup	seconds per second")
    oozie_cgroup_mem_page_cache                                                                                             =create_oozie_metric_struct("oozie_cgroup_mem_page_cache", "Page cache usage of the role's cgroup	bytes")
    oozie_cgroup_mem_rss                                                                                                    =create_oozie_metric_struct("oozie_cgroup_mem_rss", "Resident memory of the role's cgroup	bytes")
    oozie_cgroup_mem_swap                                                                                                   =create_oozie_metric_struct("oozie_cgroup_mem_swap", "Swap usage of the role's cgroup	bytes")
    oozie_cgroup_read_bytes_rate                                                                                            =create_oozie_metric_struct("oozie_cgroup_read_bytes_rate", "Bytes read from all disks by the role's cgroup	bytes per second")
    oozie_cgroup_read_ios_rate                                                                                              =create_oozie_metric_struct("oozie_cgroup_read_ios_rate", "Number of read I/O operations from all disks by the role's cgroup	ios per second")
    oozie_cgroup_write_bytes_rate                                                                                           =create_oozie_metric_struct("oozie_cgroup_write_bytes_rate", "Bytes written to all disks by the role's cgroup	bytes per second")
    oozie_cgroup_write_ios_rate                                                                                             =create_oozie_metric_struct("oozie_cgroup_write_ios_rate", "Number of write I/O operations to all disks by the role's cgroup	ios per second")
    oozie_cpu_system_rate                                                                                                   =create_oozie_metric_struct("oozie_cpu_system_rate", "Total System CPU	seconds per second")
    oozie_cpu_user_rate                                                                                                     =create_oozie_metric_struct("oozie_cpu_user_rate", "Total CPU user time	seconds per second")
    oozie_events_critical_rate                                                                                              =create_oozie_metric_struct("oozie_events_critical_rate", "The number of critical events.	events per second")
    oozie_events_important_rate                                                                                             =create_oozie_metric_struct("oozie_events_important_rate", "The number of important events.	events per second")
    oozie_events_informational_rate                                                                                         =create_oozie_metric_struct("oozie_events_informational_rate", "The number of informational events.	events per second")
    oozie_fd_open                                                                                                           =create_oozie_metric_struct("oozie_fd_open", "Open file descriptors.	file descriptors")
    oozie_jvm_total_memory                                                                                                  =create_oozie_metric_struct("oozie_jvm_total_memory", "The total amount of memory in the Java virtual machine.	bytes")
    oozie_mem_rss                                                                                                           =create_oozie_metric_struct("oozie_mem_rss", "Resident memory used	bytes")
    oozie_mem_swap                                                                                                          =create_oozie_metric_struct("oozie_mem_swap", "Amount of swap memory used by this role's process.	bytes")
    oozie_mem_virtual                                                                                                       =create_oozie_metric_struct("oozie_mem_virtual", "Virtual memory used	bytes")
    oozie_bundle_action_query_executor_get_bundle_actions_for_bundle_duration_timer_avg                                     =create_oozie_metric_struct("oozie_bundle_action_query_executor_get_bundle_actions_for_bundle_duration_timer_avg", "message.metrics.oozie_bundle_action_query_executor_get_bundle_actions_for_bundle_duration_timer_avg.desc	ms")
    oozie_bundle_action_query_executor_get_bundle_actions_for_bundle_duration_timer_rate                                    =create_oozie_metric_struct("oozie_bundle_action_query_executor_get_bundle_actions_for_bundle_duration_timer_rate", "message.metrics.oozie_bundle_action_query_executor_get_bundle_actions_for_bundle_duration_timer_count.desc	message.units.executions per second")
    oozie_bundle_job_query_executor_get_bundle_job_duration_timer_avg                                                       =create_oozie_metric_struct("oozie_bundle_job_query_executor_get_bundle_job_duration_timer_avg", "message.metrics.oozie_bundle_job_query_executor_get_bundle_job_duration_timer_avg.desc	ms")
    oozie_bundle_job_query_executor_get_bundle_job_duration_timer_rate                                                      =create_oozie_metric_struct("oozie_bundle_job_query_executor_get_bundle_job_duration_timer_rate", "message.metrics.oozie_bundle_job_query_executor_get_bundle_job_duration_timer_count.desc	message.units.executions per second")
    oozie_callback_rest_call_duration_timer_avg                                                                             =create_oozie_metric_struct("oozie_callback_rest_call_duration_timer_avg", "message.metrics.oozie_callback_rest_call_duration_timer_avg.desc	ms")
    oozie_callback_rest_call_duration_timer_rate                                                                            =create_oozie_metric_struct("oozie_callback_rest_call_duration_timer_rate", "message.metrics.oozie_callback_rest_call_duration_timer_count.desc	Calls per second")
    oozie_coord_action_query_executor_get_coord_action_duration_timer_avg                                                   =create_oozie_metric_struct("oozie_coord_action_query_executor_get_coord_action_duration_timer_avg", "message.metrics.oozie_coord_action_query_executor_get_coord_action_duration_timer_avg.desc	ms")
    oozie_coord_action_query_executor_get_coord_action_duration_timer_rate                                                  =create_oozie_metric_struct("oozie_coord_action_query_executor_get_coord_action_duration_timer_rate", "message.metrics.oozie_coord_action_query_executor_get_coord_action_duration_timer_count.desc	message.units.executions per second")
    oozie_coord_action_query_executor_get_coord_active_actions_by_jobid_duration_timer_rate                                 =create_oozie_metric_struct("oozie_coord_action_query_executor_get_coord_active_actions_by_jobid_duration_timer_rate", "message.metrics.oozie_coord_action_query_executor_get_coord_active_actions_count_by_jobid_duration_timer_count.desc	message.units.executions per second")
    oozie_coord_action_query_executor_get_coord_active_actions_count_by_jobid_duration_timer_avg                            =create_oozie_metric_struct("oozie_coord_action_query_executor_get_coord_active_actions_count_by_jobid_duration_timer_avg", "message.metrics.oozie_coord_action_query_executor_get_coord_active_actions_count_by_jobid_duration_timer_avg.desc	ms")
    oozie_coord_job_query_executor_get_coord_job_action_ready_duration_timer_avg                                            =create_oozie_metric_struct("oozie_coord_job_query_executor_get_coord_job_action_ready_duration_timer_avg", "message.metrics.oozie_coord_job_query_executor_get_coord_job_action_ready_duration_timer_avg.desc	ms")
    oozie_coord_job_query_executor_get_coord_job_action_ready_duration_timer_rate                                           =create_oozie_metric_struct("oozie_coord_job_query_executor_get_coord_job_action_ready_duration_timer_rate", "message.metrics.oozie_coord_job_query_executor_get_coord_job_action_ready_duration_timer_count.desc	message.units.executions per second")
    oozie_coord_job_query_executor_get_coord_job_duration_timer_avg                                                         =create_oozie_metric_struct("oozie_coord_job_query_executor_get_coord_job_duration_timer_avg", "message.metrics.oozie_coord_job_query_executor_get_coord_job_duration_timer_avg.desc	ms")
    oozie_coord_job_query_executor_get_coord_job_duration_timer_rate                                                        =create_oozie_metric_struct("oozie_coord_job_query_executor_get_coord_job_duration_timer_rate", "message.metrics.oozie_coord_job_query_executor_get_coord_job_duration_timer_count.desc	message.units.executions per second")
    oozie_coord_job_query_executor_get_coord_job_user_appname_duration_timer_avg                                            =create_oozie_metric_struct("oozie_coord_job_query_executor_get_coord_job_user_appname_duration_timer_avg", "message.metrics.oozie_coord_job_query_executor_get_coord_job_user_appname_duration_timer_avg.desc	ms")
    oozie_coord_job_query_executor_get_coord_job_user_appname_duration_timer_rate                                           =create_oozie_metric_struct("oozie_coord_job_query_executor_get_coord_job_user_appname_duration_timer_rate", "message.metrics.oozie_coord_job_query_executor_get_coord_job_user_appname_duration_timer_count.desc	message.units.executions per second")
    oozie_coord_job_query_executor_get_coord_jobs_changed_duration_timer_avg                                                =create_oozie_metric_struct("oozie_coord_job_query_executor_get_coord_jobs_changed_duration_timer_avg", "message.metrics.oozie_coord_job_query_executor_get_coord_jobs_changed_duration_timer_avg.desc	ms")
    oozie_coord_job_query_executor_get_coord_jobs_changed_duration_timer_rate                                               =create_oozie_metric_struct("oozie_coord_job_query_executor_get_coord_jobs_changed_duration_timer_rate", "message.metrics.oozie_coord_job_query_executor_get_coord_jobs_changed_duration_timer_count.desc	message.units.executions per second")
    oozie_jvm_pause_time_rate                                                                                               =create_oozie_metric_struct("oozie_jvm_pause_time_rate", "message.metrics.oozie_jvm_pause_time.desc	ms per second")
    oozie_jvm_pauses_info_threshold_rate                                                                                    =create_oozie_metric_struct("oozie_jvm_pauses_info_threshold_rate", "message.metrics.oozie_jvm_pauses_info_threshold_count.desc	pauses per second")
    oozie_jvm_pauses_warn_threshold_rate                                                                                    =create_oozie_metric_struct("oozie_jvm_pauses_warn_threshold_rate", "message.metrics.oozie_jvm_pauses_warn_threshold_count.desc	pauses per second")
    oozie_memory_heap_committed                                                                                             =create_oozie_metric_struct("oozie_memory_heap_committed", "message.metrics.oozie_memory_heap_committed.desc	bytes")
    oozie_memory_heap_init                                                                                                  =create_oozie_metric_struct("oozie_memory_heap_init", "message.metrics.oozie_memory_heap_init.desc	bytes")
    oozie_memory_heap_used                                                                                                  =create_oozie_metric_struct("oozie_memory_heap_used", "message.metrics.oozie_memory_heap_used.desc	bytes")
    oozie_memory_non_heap_committed                                                                                         =create_oozie_metric_struct("oozie_memory_non_heap_committed", "message.metrics.oozie_memory_non_heap_committed.desc	bytes")
    oozie_memory_non_heap_init                                                                                              =create_oozie_metric_struct("oozie_memory_non_heap_init", "message.metrics.oozie_memory_non_heap_init.desc	bytes")
    oozie_memory_non_heap_used                                                                                              =create_oozie_metric_struct("oozie_memory_non_heap_used", "message.metrics.oozie_memory_non_heap_used.desc	bytes")
    oozie_memory_total_committed                                                                                            =create_oozie_metric_struct("oozie_memory_total_committed", "message.metrics.oozie_memory_total_committed.desc	bytes")
    oozie_memory_total_init                                                                                                 =create_oozie_metric_struct("oozie_memory_total_init", "message.metrics.oozie_memory_total_init.desc	bytes")
    oozie_memory_total_used                                                                                                 =create_oozie_metric_struct("oozie_memory_total_used", "message.metrics.oozie_memory_total_used.desc	bytes")
    oozie_v1job_rest_call_duration_timer_avg                                                                                =create_oozie_metric_struct("oozie_v1job_rest_call_duration_timer_avg", "message.metrics.oozie_v1job_rest_call_duration_timer_avg.desc	ms")
    oozie_v1job_rest_call_duration_timer_rate                                                                               =create_oozie_metric_struct("oozie_v1job_rest_call_duration_timer_rate", "message.metrics.oozie_v1job_rest_call_duration_timer_count.desc	Calls per second")
    oozie_v1jobs_rest_call_duration_timer_avg                                                                               =create_oozie_metric_struct("oozie_v1jobs_rest_call_duration_timer_avg", "message.metrics.oozie_v1jobs_rest_call_duration_timer_avg.desc	ms")
    oozie_v1jobs_rest_call_duration_timer_rate                                                                              =create_oozie_metric_struct("oozie_v1jobs_rest_call_duration_timer_rate", "message.metrics.oozie_v1jobs_rest_call_duration_timer_count.desc	Calls per second")
    oozie_v2job_rest_call_duration_timer_avg                                                                                =create_oozie_metric_struct("oozie_v2job_rest_call_duration_timer_avg", "message.metrics.oozie_v2job_rest_call_duration_timer_avg.desc	ms")
    oozie_v2job_rest_call_duration_timer_rate                                                                               =create_oozie_metric_struct("oozie_v2job_rest_call_duration_timer_rate", "message.metrics.oozie_v2job_rest_call_duration_timer_count.desc	Calls per second")
    oozie_version_rest_call_duration_timer_avg                                                                              =create_oozie_metric_struct("oozie_version_rest_call_duration_timer_avg", "message.metrics.oozie_version_rest_call_duration_timer_avg.desc	ms")
    oozie_version_rest_call_duration_timer_rate                                                                             =create_oozie_metric_struct("oozie_version_rest_call_duration_timer_rate", "message.metrics.oozie_version_rest_call_duration_timer_count.desc	Calls per second")
    oozie_webservices_requests_rate                                                                                         =create_oozie_metric_struct("oozie_webservices_requests_rate", "message.metrics.oozie_webservices_requests.desc	requests per second")
    oozie_workflow_action_query_executor_get_action_check_duration_timer_avg                                                =create_oozie_metric_struct("oozie_workflow_action_query_executor_get_action_check_duration_timer_avg", "message.metrics.oozie_workflow_action_query_executor_get_action_check_duration_timer_avg.desc	ms")
    oozie_workflow_action_query_executor_get_action_check_duration_timer_rate                                               =create_oozie_metric_struct("oozie_workflow_action_query_executor_get_action_check_duration_timer_rate", "message.metrics.oozie_workflow_action_query_executor_get_action_check_duration_timer_count.desc	message.units.executions per second")
    oozie_workflow_action_query_executor_get_action_completed_duration_timer_avg                                            =create_oozie_metric_struct("oozie_workflow_action_query_executor_get_action_completed_duration_timer_avg", "message.metrics.oozie_workflow_action_query_executor_get_action_completed_duration_timer_avg.desc	ms")
    oozie_workflow_action_query_executor_get_action_completed_duration_timer_rate                                           =create_oozie_metric_struct("oozie_workflow_action_query_executor_get_action_completed_duration_timer_rate", "message.metrics.oozie_workflow_action_query_executor_get_action_completed_duration_timer_count.desc	message.units.executions per second")
    oozie_workflow_action_query_executor_get_action_duration_timer_avg                                                      =create_oozie_metric_struct("oozie_workflow_action_query_executor_get_action_duration_timer_avg", "message.metrics.oozie_workflow_action_query_executor_get_action_duration_timer_avg.desc	ms")
    oozie_workflow_action_query_executor_get_action_duration_timer_rate                                                     =create_oozie_metric_struct("oozie_workflow_action_query_executor_get_action_duration_timer_rate", "message.metrics.oozie_workflow_action_query_executor_get_action_duration_timer_count.desc	message.units.executions per second")
    oozie_workflow_action_query_executor_get_pending_actions_duration_timer_avg                                             =create_oozie_metric_struct("oozie_workflow_action_query_executor_get_pending_actions_duration_timer_avg", "message.metrics.oozie_workflow_action_query_executor_get_pending_actions_duration_timer_avg.desc	ms")
    oozie_workflow_action_query_executor_get_pending_actions_duration_timer_rate                                            =create_oozie_metric_struct("oozie_workflow_action_query_executor_get_pending_actions_duration_timer_rate", "message.metrics.oozie_workflow_action_query_executor_get_pending_actions_duration_timer_count.desc	message.units.executions per second")
    oozie_workflow_action_query_executor_get_running_actions_duration_timer_avg                                             =create_oozie_metric_struct("oozie_workflow_action_query_executor_get_running_actions_duration_timer_avg", "message.metrics.oozie_workflow_action_query_executor_get_running_actions_duration_timer_avg.desc	ms")
    oozie_workflow_action_query_executor_get_running_actions_duration_timer_rate                                            =create_oozie_metric_struct("oozie_workflow_action_query_executor_get_running_actions_duration_timer_rate", "message.metrics.oozie_workflow_action_query_executor_get_running_actions_duration_timer_count.desc	message.units.executions per second")
    oozie_workflow_job_query_executor_get_workflow_duration_timer_avg                                                       =create_oozie_metric_struct("oozie_workflow_job_query_executor_get_workflow_duration_timer_avg", "message.metrics.oozie_workflow_job_query_executor_get_workflow_duration_timer_avg.desc	ms")
    oozie_workflow_job_query_executor_get_workflow_duration_timer_rate                                                      =create_oozie_metric_struct("oozie_workflow_job_query_executor_get_workflow_duration_timer_rate", "message.metrics.oozie_workflow_job_query_executor_get_workflow_duration_timer_count.desc	message.units.executions per second")
    oozie_workflow_job_query_executor_get_workflow_start_end_time_duration_timer_avg                                        =create_oozie_metric_struct("oozie_workflow_job_query_executor_get_workflow_start_end_time_duration_timer_avg", "message.metrics.oozie_workflow_job_query_executor_get_workflow_start_end_time_duration_timer_avg.desc	ms")
    oozie_workflow_job_query_executor_get_workflow_start_end_time_duration_timer_rate                                       =create_oozie_metric_struct("oozie_workflow_job_query_executor_get_workflow_start_end_time_duration_timer_rate", "message.metrics.oozie_workflow_job_query_executor_get_workflow_start_end_time_duration_timer_count.desc	message.units.executions per second")
    oozie_pause_time_rate                                                                                                   =create_oozie_metric_struct("oozie_pause_time_rate", "Total time spent paused. This is the total extra time the pause monitor thread spent sleeping on top of the requested 500 ms.	ms per second")
    oozie_pauses_rate                                                                                                       =create_oozie_metric_struct("oozie_pauses_rate", "Number of pauses detected. The pause monitor thread sleeps for 500 ms and calculates the extra time it spent paused on top of the sleep time. If the extra sleep time exceeds 1 second, it treats it as one pause.	pauses per second")
    oozie_read_bytes_rate                                                                                                   =create_oozie_metric_struct("oozie_read_bytes_rate", "The number of bytes read from the device	bytes per second")
    oozie_unexpected_exits_rate                                                                                             =create_oozie_metric_struct("oozie_unexpected_exits_rate", "The number of times the role's backing process exited unexpectedly.	exits per second")
    oozie_uptime                                                                                                            =create_oozie_metric_struct("oozie_uptime", "For a host, the amount of time since the host was booted. For a role, the uptime of the backing process.	seconds")
    oozie_web_metrics_collection_duration                                                                                   =create_oozie_metric_struct("oozie_web_metrics_collection_duration", "Web Server Responsiveness	ms")
    oozie_write_bytes_rate                                                                                                  =create_oozie_metric_struct("oozie_write_bytes_rate", "The number of bytes written to the device	bytes per second")
  )

  var oozie_query_variable_relationship = []relation {
    {OOZIE_ALERTS_RATE,                                                                                               *oozie_alerts_rate},
    {OOZIE_CALLABLEQUEUE_ITEMS_EXECUTED_RATE,                                                                         *oozie_callablequeue_items_executed_rate},
    {OOZIE_CALLABLEQUEUE_ITEMS_FAILED_RATE,                                                                           *oozie_callablequeue_items_failed_rate},
    {OOZIE_CALLABLEQUEUE_ITEMS_QUEUED_RATE,                                                                           *oozie_callablequeue_items_queued_rate},
    {OOZIE_CGROUP_CPU_SYSTEM_RATE,                                                                                    *oozie_cgroup_cpu_system_rate},
    {OOZIE_CGROUP_CPU_USER_RATE,                                                                                      *oozie_cgroup_cpu_user_rate},
    {OOZIE_CGROUP_MEM_PAGE_CACHE,                                                                                     *oozie_cgroup_mem_page_cache},
    {OOZIE_CGROUP_MEM_RSS,                                                                                            *oozie_cgroup_mem_rss},
    {OOZIE_CGROUP_MEM_SWAP,                                                                                           *oozie_cgroup_mem_swap},
    {OOZIE_CGROUP_READ_BYTES_RATE,                                                                                    *oozie_cgroup_read_bytes_rate},
    {OOZIE_CGROUP_READ_IOS_RATE,                                                                                      *oozie_cgroup_read_ios_rate},
    {OOZIE_CGROUP_WRITE_BYTES_RATE,                                                                                   *oozie_cgroup_write_bytes_rate},
    {OOZIE_CGROUP_WRITE_IOS_RATE,                                                                                     *oozie_cgroup_write_ios_rate},
    {OOZIE_CPU_SYSTEM_RATE,                                                                                           *oozie_cpu_system_rate},
    {OOZIE_CPU_USER_RATE,                                                                                             *oozie_cpu_user_rate},
    {OOZIE_EVENTS_CRITICAL_RATE,                                                                                      *oozie_events_critical_rate},
    {OOZIE_EVENTS_IMPORTANT_RATE,                                                                                     *oozie_events_important_rate},
    {OOZIE_EVENTS_INFORMATIONAL_RATE,                                                                                 *oozie_events_informational_rate},
    {OOZIE_FD_OPEN,                                                                                                   *oozie_fd_open},
    {OOZIE_JVM_TOTAL_MEMORY,                                                                                          *oozie_jvm_total_memory},
    {OOZIE_MEM_RSS,                                                                                                   *oozie_mem_rss},
    {OOZIE_MEM_SWAP,                                                                                                  *oozie_mem_swap},
    {OOZIE_MEM_VIRTUAL,                                                                                               *oozie_mem_virtual},
    {OOZIE_OOM_EXITS_RATE,                                                                                            *oozie_oom_exits_rate},
    {OOZIE_BUNDLE_ACTION_QUERY_EXECUTOR_GET_BUNDLE_ACTIONS_FOR_BUNDLE_DURATION_TIMER_AVG,                             *oozie_bundle_action_query_executor_get_bundle_actions_for_bundle_duration_timer_avg},
    {OOZIE_BUNDLE_ACTION_QUERY_EXECUTOR_GET_BUNDLE_ACTIONS_FOR_BUNDLE_DURATION_TIMER_RATE,                            *oozie_bundle_action_query_executor_get_bundle_actions_for_bundle_duration_timer_rate},
    {OOZIE_BUNDLE_JOB_QUERY_EXECUTOR_GET_BUNDLE_JOB_DURATION_TIMER_AVG,                                               *oozie_bundle_job_query_executor_get_bundle_job_duration_timer_avg},
    {OOZIE_BUNDLE_JOB_QUERY_EXECUTOR_GET_BUNDLE_JOB_DURATION_TIMER_RATE,                                              *oozie_bundle_job_query_executor_get_bundle_job_duration_timer_rate},
    {OOZIE_CALLBACK_REST_CALL_DURATION_TIMER_AVG,                                                                     *oozie_callback_rest_call_duration_timer_avg},
    {OOZIE_CALLBACK_REST_CALL_DURATION_TIMER_RATE,                                                                    *oozie_callback_rest_call_duration_timer_rate},
    {OOZIE_COORD_ACTION_QUERY_EXECUTOR_GET_COORD_ACTION_DURATION_TIMER_AVG,                                           *oozie_coord_action_query_executor_get_coord_action_duration_timer_avg},
    {OOZIE_COORD_ACTION_QUERY_EXECUTOR_GET_COORD_ACTION_DURATION_TIMER_RATE,                                          *oozie_coord_action_query_executor_get_coord_action_duration_timer_rate},
    {OOZIE_COORD_ACTION_QUERY_EXECUTOR_GET_COORD_ACTIVE_ACTIONS_BY_JOBID_DURATION_TIMER_RATE,                         *oozie_coord_action_query_executor_get_coord_active_actions_by_jobid_duration_timer_rate},
    {OOZIE_COORD_ACTION_QUERY_EXECUTOR_GET_COORD_ACTIVE_ACTIONS_COUNT_BY_JOBID_DURATION_TIMER_AVG,                    *oozie_coord_action_query_executor_get_coord_active_actions_count_by_jobid_duration_timer_avg},
    {OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOB_ACTION_READY_DURATION_TIMER_AVG,                                    *oozie_coord_job_query_executor_get_coord_job_action_ready_duration_timer_avg},
    {OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOB_ACTION_READY_DURATION_TIMER_RATE,                                   *oozie_coord_job_query_executor_get_coord_job_action_ready_duration_timer_rate},
    {OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOB_DURATION_TIMER_AVG,                                                 *oozie_coord_job_query_executor_get_coord_job_duration_timer_avg},
    {OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOB_DURATION_TIMER_RATE,                                                *oozie_coord_job_query_executor_get_coord_job_duration_timer_rate},
    {OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOB_USER_APPNAME_DURATION_TIMER_AVG,                                    *oozie_coord_job_query_executor_get_coord_job_user_appname_duration_timer_avg},
    {OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOB_USER_APPNAME_DURATION_TIMER_RATE,                                   *oozie_coord_job_query_executor_get_coord_job_user_appname_duration_timer_rate},
    {OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOBS_CHANGED_DURATION_TIMER_AVG,                                        *oozie_coord_job_query_executor_get_coord_jobs_changed_duration_timer_avg},
    {OOZIE_COORD_JOB_QUERY_EXECUTOR_GET_COORD_JOBS_CHANGED_DURATION_TIMER_RATE,                                       *oozie_coord_job_query_executor_get_coord_jobs_changed_duration_timer_rate},
    {OOZIE_JVM_PAUSE_TIME_RATE,                                                                                       *oozie_jvm_pause_time_rate},
    {OOZIE_JVM_PAUSES_INFO_THRESHOLD_RATE,                                                                            *oozie_jvm_pauses_info_threshold_rate},
    {OOZIE_JVM_PAUSES_WARN_THRESHOLD_RATE,                                                                            *oozie_jvm_pauses_warn_threshold_rate},
    {OOZIE_MEMORY_HEAP_COMMITTED,                                                                                     *oozie_memory_heap_committed},
    {OOZIE_MEMORY_HEAP_INIT,                                                                                          *oozie_memory_heap_init},
    {OOZIE_MEMORY_HEAP_USED,                                                                                          *oozie_memory_heap_used},
    {OOZIE_MEMORY_NON_HEAP_COMMITTED,                                                                                 *oozie_memory_non_heap_committed},
    {OOZIE_MEMORY_NON_HEAP_INIT,                                                                                      *oozie_memory_non_heap_init},
    {OOZIE_MEMORY_NON_HEAP_USED,                                                                                      *oozie_memory_non_heap_used},
    {OOZIE_MEMORY_TOTAL_COMMITTED,                                                                                    *oozie_memory_total_committed},
    {OOZIE_MEMORY_TOTAL_INIT,                                                                                         *oozie_memory_total_init},
    {OOZIE_MEMORY_TOTAL_USED,                                                                                         *oozie_memory_total_used},
    {OOZIE_V1JOB_REST_CALL_DURATION_TIMER_AVG,                                                                        *oozie_v1job_rest_call_duration_timer_avg},
    {OOZIE_V1JOB_REST_CALL_DURATION_TIMER_RATE,                                                                       *oozie_v1job_rest_call_duration_timer_rate},
    {OOZIE_V1JOBS_REST_CALL_DURATION_TIMER_AVG,                                                                       *oozie_v1jobs_rest_call_duration_timer_avg},
    {OOZIE_V1JOBS_REST_CALL_DURATION_TIMER_RATE,                                                                      *oozie_v1jobs_rest_call_duration_timer_rate},
    {OOZIE_V2JOB_REST_CALL_DURATION_TIMER_AVG,                                                                        *oozie_v2job_rest_call_duration_timer_avg},
    {OOZIE_V2JOB_REST_CALL_DURATION_TIMER_RATE,                                                                       *oozie_v2job_rest_call_duration_timer_rate},
    {OOZIE_VERSION_REST_CALL_DURATION_TIMER_AVG,                                                                      *oozie_version_rest_call_duration_timer_avg},
    {OOZIE_VERSION_REST_CALL_DURATION_TIMER_RATE,                                                                     *oozie_version_rest_call_duration_timer_rate},
    {OOZIE_WEBSERVICES_REQUESTS_RATE,                                                                                 *oozie_webservices_requests_rate},
    {OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_ACTION_CHECK_DURATION_TIMER_AVG,                                        *oozie_workflow_action_query_executor_get_action_check_duration_timer_avg},
    {OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_ACTION_CHECK_DURATION_TIMER_RATE,                                       *oozie_workflow_action_query_executor_get_action_check_duration_timer_rate},
    {OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_ACTION_COMPLETED_DURATION_TIMER_AVG,                                    *oozie_workflow_action_query_executor_get_action_completed_duration_timer_avg},
    {OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_ACTION_COMPLETED_DURATION_TIMER_RATE,                                   *oozie_workflow_action_query_executor_get_action_completed_duration_timer_rate},
    {OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_ACTION_DURATION_TIMER_AVG,                                              *oozie_workflow_action_query_executor_get_action_duration_timer_avg},
    {OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_ACTION_DURATION_TIMER_RATE,                                             *oozie_workflow_action_query_executor_get_action_duration_timer_rate},
    {OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_PENDING_ACTIONS_DURATION_TIMER_AVG,                                     *oozie_workflow_action_query_executor_get_pending_actions_duration_timer_avg},
    {OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_PENDING_ACTIONS_DURATION_TIMER_RATE,                                    *oozie_workflow_action_query_executor_get_pending_actions_duration_timer_rate},
    {OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_RUNNING_ACTIONS_DURATION_TIMER_AVG,                                     *oozie_workflow_action_query_executor_get_running_actions_duration_timer_avg},
    {OOZIE_WORKFLOW_ACTION_QUERY_EXECUTOR_GET_RUNNING_ACTIONS_DURATION_TIMER_RATE,                                    *oozie_workflow_action_query_executor_get_running_actions_duration_timer_rate},
    {OOZIE_WORKFLOW_JOB_QUERY_EXECUTOR_GET_WORKFLOW_DURATION_TIMER_AVG,                                               *oozie_workflow_job_query_executor_get_workflow_duration_timer_avg},
    {OOZIE_WORKFLOW_JOB_QUERY_EXECUTOR_GET_WORKFLOW_DURATION_TIMER_RATE,                                              *oozie_workflow_job_query_executor_get_workflow_duration_timer_rate},
    {OOZIE_WORKFLOW_JOB_QUERY_EXECUTOR_GET_WORKFLOW_START_END_TIME_DURATION_TIMER_AVG,                                *oozie_workflow_job_query_executor_get_workflow_start_end_time_duration_timer_avg},
    {OOZIE_WORKFLOW_JOB_QUERY_EXECUTOR_GET_WORKFLOW_START_END_TIME_DURATION_TIMER_RATE,                               *oozie_workflow_job_query_executor_get_workflow_start_end_time_duration_timer_rate},
    {OOZIE_PAUSE_TIME_RATE,                                                                                           *oozie_pause_time_rate},
    {OOZIE_PAUSES_RATE,                                                                                               *oozie_pauses_rate},
    {OOZIE_READ_BYTES_RATE,                                                                                           *oozie_read_bytes_rate},
    {OOZIE_UNEXPECTED_EXITS_RATE,                                                                                     *oozie_unexpected_exits_rate},
    {OOZIE_UPTIME,                                                                                                    *oozie_uptime},
    {OOZIE_WEB_METRICS_COLLECTION_DURATION,                                                                           *oozie_web_metrics_collection_duration},
    {OOZIE_WRITE_BYTES_RATE,                                                                                          *oozie_write_bytes_rate},
  }
 
 /* ======================================================================
  * Functions
  * ====================================================================== */
 // Create and returns a prometheus descriptor for a oozie metric. 
 // The "metric_name" parameter its mandatory
 // If the "description" parameter is empty, the function assings it with the
 // value of the name of the metric in uppercase and separated by spaces
 func create_oozie_metric_struct(metric_name string, description string) *prometheus.Desc {
   // Correct "description" parameter if is empty
   if len(description) == 0 {
     description = strings.Replace(strings.ToUpper(metric_name), "_", " ", 0)
   }
 
   // return prometheus descriptor
   return prometheus.NewDesc(
     prometheus.BuildFQName(namespace, OOZIE_SCRAPER_NAME, metric_name),
     description,
     []string{"cluster", "entityName","hostname"},
     nil,
   )
 }
 
 
 // Generic function to extract de metadata associated with the query value
 // Only for OOZIE metric type
 func create_oozie_metric (ctx context.Context, config Collector_connection_data, query string, metric_struct prometheus.Desc, ch chan<- prometheus.Metric, pclient *pool.PClient) bool {
   // Make the query
   json_parsed, err := make_and_parse_timeseries_query(ctx, config, query, pclient)
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
 // ScrapeOOZIE struct
 type ScrapeOOZIE struct{}
 
 // Name of the Scraper. Should be unique.
 func (ScrapeOOZIE) Name() string {
   return OOZIE_SCRAPER_NAME
 }
 
 // Help describes the role of the Scraper.
 func (ScrapeOOZIE) Help() string {
   return "OOZIE Metrics"
 }
 
 // Version.
 func (ScrapeOOZIE) Version() float64 {
   return 1.0
 }
 
 // Scrape generic function. Override for host module.
 func (ScrapeOOZIE) Scrape (ctx context.Context, config *Collector_connection_data, ch chan<- prometheus.Metric) error {
   log.Debug_msg("Executing OOZIE Metrics Scraper")
 
   // Queries counters
   success_queries := 0
   error_queries := 0
 
   pclient := pool.NewPClient()
   // Execute the generic funtion for creation of metrics with the pairs (QUERY, PROM:DESCRIPTOR)
   for i:=0 ; i < len(oozie_query_variable_relationship) ; i++ {
     if create_oozie_metric(ctx, *config, oozie_query_variable_relationship[i].Query, oozie_query_variable_relationship[i].Metric_struct, ch, pclient) {
       success_queries += 1
     } else {
       error_queries += 1
     }
   }
   log.Info_msg("In the OOZIE Module has been executed %d queries. %d success and %d with errors", success_queries + error_queries, success_queries, error_queries)
   return nil
 }
 
 var _ Scraper = ScrapeOOZIE{}
 