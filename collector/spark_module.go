/*
 *
 * title           :collector/spark_module.go
 * description     :Submodule Collector for the Cluster SPARK metrics
 * author               :Raul Barroso
 * co-author           :Alejandro Villegas
 * date            :2019/02/11
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
  * Constants with the Host module TSquery sentences
  * ====================================================================== */
  const SPARK_SCRAPER_NAME = "spark"
  var (
    // Agent Queries
    SPARK_CATALOG_JVM_COMITTED_BYTES =           "SELECT LAST(spark_catalogserver_jvm_heap_committed_usage_bytes) WHERE serviceType = \"SPARK\""
    SPARK_CATALOG_JVM_CURRENT_BYTES =            "SELECT LAST(spark_catalogserver_jvm_heap_current_usage_bytes) WHERE serviceType = \"SPARK\""
    SPARK_CATALOG_JVM_INIT_BYTES =               "SELECT LAST(spark_catalogserver_jvm_heap_init_usage_bytes) WHERE serviceType = \"SPARK\""
    SPARK_CATALOG_JVM_MAX_BYTES =                "SELECT LAST(spark_catalogserver_jvm_heap_max_usage_bytes) WHERE serviceType = \"SPARK\""
    SPARK_CGROUP_MEM_PAGE_CACHE =                "SELECT LAST(cgroup_mem_page_cache) WHERE serviceType = \"SPARK\""
    SPARK_CGROUP_MEM_RSS =                       "SELECT LAST(cgroup_mem_rss) WHERE serviceType = \"SPARK\""
    SPARK_CGROUP_MEM_SWAP =                      "SELECT LAST(cgroup_mem_swap) WHERE serviceType = \"SPARK\""
    SPARK_CGROUP_READ_IOSRATE =                  "SELECT LAST(INTEGRAL(cgroup_read_ios_rate)) WHERE serviceType = \"SPARK\""
    SPARK_CGROUP_READ_RATE =                     "SELECT LAST(INTEGRAL(cgroup_read_bytes_rate)) WHERE serviceType = \"SPARK\""
    SPARK_CGROUP_SYSTEM_RATE =                   "SELECT LAST(INTEGRAL(cgroup_cpu_system_rate)) WHERE serviceType = \"SPARK\""
    SPARK_CGROUP_USER_RATE =                     "SELECT LAST(INTEGRAL(cgroup_cpu_user_rate)) WHERE serviceType = \"SPARK\""
    SPARK_CGROUP_WRITE_IOSRATE =                 "SELECT LAST(INTEGRAL(cgroup_write_ios_rate)) WHERE serviceType = \"SPARK\""
    SPARK_CGROUP_WRITE_RATE =                    "SELECT LAST(INTEGRAL(cgroup_write_bytes_rate)) WHERE serviceType = \"SPARK\""
    SPARK_MEM_RSS =                              "SELECT LAST(mem_rss) WHERE serviceType = \"SPARK\""
    SPARK_MEM_SWAP =                             "SELECT LAST(mem_swap) WHERE serviceType = \"SPARK\""
    SPARK_MEM_VIRT =                             "SELECT LAST(mem_virtual) WHERE serviceType = \"SPARK\""
    SPARK_OOMEXIT =                              "SELECT LAST(INTEGRAL(oom_exits_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_ADMISSION_WAIT_RATE =            "SELECT LAST(INTEGRAL(spark_query_admission_wait_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_BYTES_HDFS_READ_RATE =           "SELECT LAST(INTEGRAL(spark_query_hdfs_bytes_read_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_BYTES_HDFS_WRITTE_RATE =         "SELECT LAST(INTEGRAL(spark_query_hdfs_bytes_written_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_BYTES_STREAMED_RATE =            "SELECT LAST(INTEGRAL(spark_query_bytes_streamed_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_CM_CPU =                         "SELECT LAST(INTEGRAL(spark_query_cm_cpu_milliseconds_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_DURATION_RATE =                  "SELECT LAST(INTEGRAL(spark_query_query_duration_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_INGESTED_RATE =                  "SELECT LAST(INTEGRAL(queries_ingested_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_MEM_ACCRUAL_RATE =               "SELECT LAST(INTEGRAL(spark_query_memory_accrual_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_MEM_SPILLED_RATE =               "SELECT LAST(INTEGRAL(spark_query_memory_spilled_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_OOMRATE =                        "SELECT LAST(INTEGRAL(queries_oom_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_REJECTED_RATE =                  "SELECT LAST(INTEGRAL(queries_rejected_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_SPILLED_RATE =                   "SELECT LAST(INTEGRAL(queries_spilled_memory_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_SUCCESSFUL_RATE =                "SELECT LAST(INTEGRAL(queries_successful_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_THREAD_CPU_RATE =                "SELECT LAST(INTEGRAL(spark_query_thread_cpu_time_rate)) WHERE serviceType = \"SPARK\""
    SPARK_QUERY_TIME_OUT_RATE =                  "SELECT LAST(INTEGRAL(queries_timed_out_rate)) WHERE serviceType = \"SPARK\""
    SPARK_READ_RATE =                            "SELECT LAST(INTEGRAL(read_bytes_rate)) WHERE serviceType = \"SPARK\""
    SPARK_STATE_STORE_CACHE_TOTAL_CLIENTS =      "SELECT LAST(statestore_subscriber_statestore_client_cache_total_clients) WHERE serviceType = \"SPARK\""
    SPARK_STATE_STORE_CLIENTS_IN_USE =           "SELECT LAST(statestore_subscriber_statestore_client_cache_clients_in_use) WHERE serviceType = \"SPARK\""
    SPARK_STATE_STORE_HEART_BEAT_LAST =          "SELECT LAST(statestore_subscriber_heartbeat_interval_time_last) WHERE serviceType = \"SPARK\""
    SPARK_STATE_STORE_HEART_BEAT_MAX =           "SELECT LAST(statestore_subscriber_heartbeat_interval_time_max) WHERE serviceType = \"SPARK\""
    SPARK_STATE_STORE_HEART_BEAT_MEAN =          "SELECT LAST(statestore_subscriber_heartbeat_interval_time_mean) WHERE serviceType = \"SPARK\""
    SPARK_STATE_STORE_HEART_BEAT_MIN =           "SELECT LAST(statestore_subscriber_heartbeat_interval_time_min) WHERE serviceType = \"SPARK\""
    SPARK_STATE_STORE_HEART_BEAT_RATE =          "SELECT LAST(INTEGRAL(statestore_subscriber_heartbeat_interval_time_rate)) WHERE serviceType = \"SPARK\""
    SPARK_STATE_STORE_HEART_BEAT_STDDEV =        "SELECT LAST(statestore_subscriber_heartbeat_interval_time_stddev) WHERE serviceType = \"SPARK\""
    SPARK_STATE_STORE_LAST_RECOVERY_DURATION =   "SELECT LAST(statestore_subscriber_last_recovery_duration) WHERE serviceType = \"SPARK\""
    SPARK_TCMALLOC_FREE_BYTES =                  "SELECT LAST(tcmalloc_pageheap_free_bytes) WHERE serviceType = \"SPARK\""
    SPARK_TCMALLOC_PHYSICAL_RESERVED_BYTES =     "SELECT LAST(tcmalloc_physical_bytes_reserved) WHERE serviceType = \"SPARK\""
    SPARK_TCMALLOC_TOTAL_RESERVED_BYTES =        "SELECT LAST(tcmalloc_total_bytes_reserved) WHERE serviceType = \"SPARK\""
    SPARK_TCMALLOC_UNMAPPED_BYTES =              "SELECT LAST(tcmalloc_pageheap_unmapped_bytes) WHERE serviceType = \"SPARK\""
    SPARK_TCMALLOC_USED_BYTES =                  "SELECT LAST(tcmalloc_bytes_in_use) WHERE serviceType = \"SPARK\""
    SPARK_THRIFT_CONNECTIONS_RATE =              "SELECT LAST(INTEGRAL(thrift_server_catalog_service_connections_rate)) WHERE serviceType = \"SPARK\""
    SPARK_THRIFT_CONNECTIONS_USED =              "SELECT LAST(thrift_server_catalog_service_connections_in_use) WHERE serviceType = \"SPARK\""
    SPARK_WRITE_RATE =                           "SELECT LAST(INTEGRAL(write_bytes_rate)) WHERE serviceType = \"SPARK\""
SPARK_YARN_HISTORY_SERVER =                           "SELECT cpu_system_rate + cpu_user_rate where serviceName=\"SPARK\""
    
  )
 
 
 
 
 /* ======================================================================
  * Global variables
  * ====================================================================== */
 var (
   spark_catalog_jvm_comitted_bytes =          create_spark_metric_struct("spark_catalogserver_jvm_heap_committed_usage_bytes", "Jvm heap Committed Usage in Bytes.")
   spark_catalog_jvm_current_bytes =           create_spark_metric_struct("spark_catalogserver_jvm_heap_current_usage_bytes", "Jvm heap Current Usage in Bytes.")
   spark_catalog_jvm_init_bytes =              create_spark_metric_struct("spark_catalogserver_jvm_heap_init_usage_bytes", "JVM heap Init Usage in Bytes.")
   spark_catalog_jvm_max_bytes =               create_spark_metric_struct("spark_catalogserver_jvm_heap_max_usage_bytes", "JVM heap Max Usage in Bytes.")
   spark_cgroup_mem_page_cache =               create_spark_metric_struct("cgroup_mem_page_cache", "Page cache usage of the role's cgroup in Bytes.")
   spark_cgroup_mem_rss =                      create_spark_metric_struct("cgroup_mem_rss", "Resident memory of the role's cgroup in Bytes.")
   spark_cgroup_mem_swap =                     create_spark_metric_struct("cgroup_mem_swap", "Swap usage of the role's cgroup in Bytes.")
   spark_cgroup_read_iosrate =                 create_spark_metric_struct("cgroup_read_ios_rate", "Number of read I/O operations from all disks by the role's cgroup.")
   spark_cgroup_read_rate =                    create_spark_metric_struct("cgroup_read_bytes_rate", "Bytes read from all disks by the role's cgroup in Bytes per second.")
   spark_cgroup_system_rate =                  create_spark_metric_struct("cgroup_cpu_system_rate", "CPU usage of the role's cgroup in Bytes per second.")
   spark_cgroup_user_rate =                    create_spark_metric_struct("cgroup_cpu_user_rate", "User Space CPU usage of the role's cgroup in Bytes per second.")
   spark_cgroup_write_iosrate =                create_spark_metric_struct("cgroup_write_ios_rate", "Number of write I/O operations to all disks by the role's cgroup.")
   spark_cgroup_write_rate =                   create_spark_metric_struct("cgroup_write_bytes_rate", "Bytes written to all disks by the role's cgroup.")
   spark_mem_rss =                             create_spark_metric_struct("mem_rss", "Resident memory used in Bytes")
   spark_mem_swap =                            create_spark_metric_struct("mem_swap", "Amount of swap memory used by this role's process in Bytes")
   spark_mem_virt =                            create_spark_metric_struct("mem_virtual", "Virtual memory used in Bytes.")
   spark_oomexit =                             create_spark_metric_struct("oom_exits_rate", "The number of times the role's backing process was killed due to an OutOfMemory error. This counter is only incremented if the Cloudera Manager \"Kill When Out of Memory\" option is enabled.")
   spark_query_admission_wait_rate =           create_spark_metric_struct("spark_query_admission_wait_rate", "The time from submission for admission to its completion in milliseconds.")
   spark_query_bytes_hdfs_read_rate =          create_spark_metric_struct("spark_query_hdfs_bytes_read_rate", "The total number of bytes read from HDFS by this Spark query.")
   spark_query_bytes_hdfs_writte_rate =        create_spark_metric_struct("spark_query_hdfs_bytes_written_rate", "The total number of bytes written to HDFS by this Spark query.")
   spark_query_bytes_streamed_rate =           create_spark_metric_struct("spark_query_bytes_streamed_rate", "The total number of bytes sent between Spark Daemons while processing this query.")
   spark_query_cm_cpu =                        create_spark_metric_struct("spark_query_cm_cpu_milliseconds_rate", "spark.analysis.cm_cpu_milliseconds.description.")
   spark_query_duration_rate =                 create_spark_metric_struct("spark_query_query_duration_rate", "The duration of the query in milliseconds.")
   spark_query_ingested_rate =                 create_spark_metric_struct("queries_ingested_rate", "Spark queries ingested by the Service Monitor")
   spark_query_mem_accrual_rate =              create_spark_metric_struct("spark_query_memory_accrual_rate", "The total accrued memory usage by the query. This is computed by multiplaying the average aggregate memory usage of the query by the query's duration.")
   spark_query_mem_spilled_rate =              create_spark_metric_struct("spark_query_memory_spilled_rate", "Amount of memory spilled to disk in Bytes.")
   spark_query_oomrate =                       create_spark_metric_struct("queries_oom_rate", "Number of Spark queries for which memory consumption exceeded what was allowed")
   spark_query_rejected_rate =                 create_spark_metric_struct("queries_rejected_rate", "Number of Spark queries rejected from admission, commonly due to the queue being full or insufficient memory")
   spark_query_spilled_rate =                  create_spark_metric_struct("queries_spilleed_rate", "Number of Spark queries that spilled to disk")
   spark_query_successful_rate =               create_spark_metric_struct("queries_successful_rate", "Number of Spark queries that ran to completion successfully")
   spark_query_thread_cpu_rate =               create_spark_metric_struct("spark_query_thread_cpu_time_rate", "The sum of the CPU time used by all threads of the query.")
   spark_query_time_out_rate =                 create_spark_metric_struct("queries_time_out_rate", "Spark queries that timed out waiting in queue during admission in milliseconds")
   spark_read_rate =                           create_spark_metric_struct("read_bytes_rate", "The number of bytes read from the device.")
   spark_state_store_cache_total_clients =     create_spark_metric_struct("statestore_subscriber_statestore_client_cache_total_clients", "The total number of StateStore subscriber clients in this Spark Daemon's client cache. These clients are for communication from this role to the StateStore.")
   spark_state_store_clients_in_use =          create_spark_metric_struct("statestore_subscriber_statestore_client_cache_clients_in_use", "The number of active StateStore subscriber clients in this Spark Daemon's client cache. These clients are for communication from this role to the StateStore.")
   spark_state_store_heart_beat_last =         create_spark_metric_struct("statestore_subscriber_heartbeat_interval_time_last", "The most recent interval between heartbeats from this Spark Daemon to the StateStore in seconds.")
   spark_state_store_heart_beat_max =          create_spark_metric_struct("statestore_subscriber_heartbeat_interval_time_max", " The maximum interval between heartbeats from this Spark Daemon to the StateStore in seconds. This is calculated over the lifetime of the Spark Daemon.")
   spark_state_store_heart_beat_mean =         create_spark_metric_struct("statestore_subscriber_heartbeat_interval_time_mean", " The average interval between heartbeats from this Spark Daemon to the StateStore in seconds. This is calculated over the lifetime of the Spark Daemon.")
   spark_state_store_heart_beat_min =          create_spark_metric_struct("statestore_subscriber_heartbeat_interval_time_min", "The minimum interval between heartbeats from this Spark Daemon to the StateStore in seconds. This is calculated over the lifetime of the Spark Daemon.")
   spark_state_store_heart_beat_rate =         create_spark_metric_struct("statestore_subscriber_heartbeat_interval_time_rate", "The total number of samples taken of the Spark Daemon's StateStore heartbeat interval in samples per second.")
   spark_state_store_heart_beat_stddev =       create_spark_metric_struct("statestore_subscriber_heartbeat_interval_time_stddev", "The standard deviation in the interval between heartbeats from this Spark Daemon to the StateStore in seconds. This is calculated over the lifetime of the Spark Daemon.")
   spark_state_store_last_recovery_duration =  create_spark_metric_struct("statestore_subscriber_last_recovery_duration", "The amount of time, in seconds, the StateStore subscriber took to recover the connection the last time it was lost.")
   spark_tcmalloc_free_bytes =                 create_spark_metric_struct("tcmalloc_pageheap_free_bytes", "Number of bytes in free, mapped pages in page heap. These bytes can be used to fulfill allocation requests. They always count towards virtual memory usage, and unless the underlying memory is swapped out by the OS, they also count towards physical memory usage.")
   spark_tcmalloc_physical_reserved_bytes =    create_spark_metric_struct("tcmalloc_physical_bytes_reserved", "Derived metric computing the amount of physical memory (in bytes) used by the process, including that actually in use and free bytes reserved by tcmalloc. Does not include the tcmalloc metadata.")
   spark_tcmalloc_total_reserved_bytes =       create_spark_metric_struct("tcmalloc_total_bytes_reserved", "Bytes of system memory reserved by TCMalloc.")
   spark_tcmalloc_unmapped_bytes =             create_spark_metric_struct("tcmalloc_pageheap_unmapped_bytes", "Number of bytes in free, unmapped pages in page heap. These are bytes that have been released back to the OS, possibly by one of the MallocExtension \"Release\" calls. They can be used to fulfill allocation requests, but typically incur a page fault. They always count towards virtual memory usage, and depending on the OS, typically do not count towards physical memory usage.")
   spark_tcmalloc_used_bytes =                 create_spark_metric_struct("tcmalloc_bytes_in_use", "Number of bytes used by the application. This will not typically match the memory use reported by the OS, because it does not include TCMalloc overhead or memory fragmentation.")
   spark_thrift_connections_rate =             create_spark_metric_struct("thrift_server_catalog_service_connections_rate", "The total number of connections made to this Catalog Server's catalog service over its lifetime.")
   spark_thrift_connections_used =             create_spark_metric_struct("thrift_server_catalog_service_connections_in_use", "The number of active catalog service connections to this Catalog Server.")
   spark_write_rate =                          create_spark_metric_struct("write_bytes_rate", "The number of bytes written to the device.")
   spark_yarn_history_server =                          create_spark_metric_struct("yarn_cpu_cores", "The number of bytes written to the device.")
 
 )
 var spark_query_variable_relationship = []relation {
   {SPARK_CATALOG_JVM_COMITTED_BYTES,          *spark_catalog_jvm_comitted_bytes},
   {SPARK_CATALOG_JVM_CURRENT_BYTES,           *spark_catalog_jvm_current_bytes},
   {SPARK_CATALOG_JVM_INIT_BYTES,              *spark_catalog_jvm_init_bytes},
   {SPARK_CATALOG_JVM_MAX_BYTES,               *spark_catalog_jvm_max_bytes},
   {SPARK_CGROUP_MEM_PAGE_CACHE,               *spark_cgroup_mem_page_cache},
   {SPARK_CGROUP_MEM_RSS,                      *spark_cgroup_mem_rss},
   {SPARK_CGROUP_MEM_SWAP,                     *spark_cgroup_mem_swap},
   {SPARK_CGROUP_READ_IOSRATE,                 *spark_cgroup_read_iosrate},
   {SPARK_CGROUP_READ_RATE,                    *spark_cgroup_read_rate},
   {SPARK_CGROUP_SYSTEM_RATE,                  *spark_cgroup_system_rate},
   {SPARK_CGROUP_USER_RATE,                    *spark_cgroup_user_rate},
   {SPARK_CGROUP_WRITE_IOSRATE,                *spark_cgroup_write_iosrate},
   {SPARK_CGROUP_WRITE_RATE,                   *spark_cgroup_write_rate},
   {SPARK_MEM_RSS,                             *spark_mem_rss},
   {SPARK_MEM_SWAP,                            *spark_mem_swap},
   {SPARK_MEM_VIRT,                            *spark_mem_virt},
   {SPARK_OOMEXIT,                             *spark_oomexit},
   {SPARK_QUERY_ADMISSION_WAIT_RATE,           *spark_query_admission_wait_rate},
   {SPARK_QUERY_BYTES_HDFS_READ_RATE,          *spark_query_bytes_hdfs_read_rate},
   {SPARK_QUERY_BYTES_HDFS_WRITTE_RATE,        *spark_query_bytes_hdfs_writte_rate},
   {SPARK_QUERY_BYTES_STREAMED_RATE,           *spark_query_bytes_streamed_rate},
   {SPARK_QUERY_CM_CPU,                        *spark_query_cm_cpu},
   {SPARK_QUERY_DURATION_RATE,                 *spark_query_duration_rate},
   {SPARK_QUERY_INGESTED_RATE,                 *spark_query_ingested_rate},
   {SPARK_QUERY_MEM_ACCRUAL_RATE,              *spark_query_mem_accrual_rate},
   {SPARK_QUERY_MEM_SPILLED_RATE,              *spark_query_mem_spilled_rate},
   {SPARK_QUERY_OOMRATE,                       *spark_query_oomrate},
   {SPARK_QUERY_REJECTED_RATE,                 *spark_query_rejected_rate},
   {SPARK_QUERY_SPILLED_RATE,                  *spark_query_spilled_rate},
   {SPARK_QUERY_SUCCESSFUL_RATE,               *spark_query_successful_rate},
   {SPARK_QUERY_THREAD_CPU_RATE,               *spark_query_thread_cpu_rate},
   {SPARK_QUERY_TIME_OUT_RATE,                 *spark_query_time_out_rate},
   {SPARK_READ_RATE,                           *spark_read_rate},
   {SPARK_STATE_STORE_CACHE_TOTAL_CLIENTS,     *spark_state_store_cache_total_clients},
   {SPARK_STATE_STORE_CLIENTS_IN_USE,          *spark_state_store_clients_in_use},
   {SPARK_STATE_STORE_HEART_BEAT_LAST,         *spark_state_store_heart_beat_last},
   {SPARK_STATE_STORE_HEART_BEAT_MAX,          *spark_state_store_heart_beat_max},
   {SPARK_STATE_STORE_HEART_BEAT_MEAN,         *spark_state_store_heart_beat_mean},
   {SPARK_STATE_STORE_HEART_BEAT_MIN,          *spark_state_store_heart_beat_min},
   {SPARK_STATE_STORE_HEART_BEAT_RATE,         *spark_state_store_heart_beat_rate},
   {SPARK_STATE_STORE_HEART_BEAT_STDDEV,       *spark_state_store_heart_beat_stddev},
   {SPARK_STATE_STORE_LAST_RECOVERY_DURATION,  *spark_state_store_last_recovery_duration},
   {SPARK_TCMALLOC_FREE_BYTES,                 *spark_tcmalloc_free_bytes},
   {SPARK_TCMALLOC_PHYSICAL_RESERVED_BYTES,    *spark_tcmalloc_physical_reserved_bytes},
   {SPARK_TCMALLOC_TOTAL_RESERVED_BYTES,       *spark_tcmalloc_total_reserved_bytes},
   {SPARK_TCMALLOC_UNMAPPED_BYTES,             *spark_tcmalloc_unmapped_bytes},
   {SPARK_TCMALLOC_USED_BYTES,                 *spark_tcmalloc_used_bytes},
   {SPARK_THRIFT_CONNECTIONS_RATE,             *spark_thrift_connections_rate},
   {SPARK_THRIFT_CONNECTIONS_USED,             *spark_thrift_connections_used},
   {SPARK_WRITE_RATE,                          *spark_write_rate},
   {SPARK_YARN_HISTORY_SERVER,                          *spark_yarn_history_server},
 }
 
 
 
 
 /* ======================================================================
  * Functions
  * ====================================================================== */
 // Create and returns a prometheus descriptor for a spark metric. 
 // The "metric_name" parameter its mandatory
 // If the "description" parameter is empty, the function assings it with the
 // value of the name of the metric in uppercase and separated by spaces
 func create_spark_metric_struct(metric_name string, description string) *prometheus.Desc {
   // Correct "description" parameter if is empty
   if len(description) == 0 {
     description = strings.Replace(strings.ToUpper(metric_name), "_", " ", 0)
   }
 
   // return prometheus descriptor
   return prometheus.NewDesc(
     prometheus.BuildFQName(namespace, SPARK_SCRAPER_NAME, metric_name),
     description,
     []string{"cluster", "entityName"},
     nil,
   )
 }
 
 
 // Generic function to extract de metadata associated with the query value
 // Only for Spark metric type
 func create_spark_metric (ctx context.Context, config Collector_connection_data, query string, metric_struct prometheus.Desc, ch chan<- prometheus.Metric) bool {
   if query == "" { return true }
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
     entity_name := jp.Get_timeseries_query_entity_name(json_parsed, ts_index)
     // Get Query LAST value
     value, err := jp.Get_timeseries_query_value(json_parsed, ts_index)
     if err != nil {
       log.Debug_msg("No data for query: %s", query)
       continue
     }
     // Assing the data to the Prometheus descriptor
     ch <- prometheus.MustNewConstMetric(&metric_struct, prometheus.GaugeValue, value, cluster_name, entity_name)
   }
   return true
 }
 
 
 
 
 /* ======================================================================
  * Scrape "Class"
  * ====================================================================== */
 // ScrapeSPARK struct
 type ScrapeSPARK struct{}
 
 // Name of the Scraper. Should be unique.
 func (ScrapeSPARK) Name() string {
   return SPARK_SCRAPER_NAME
 }
 
 // Help describes the role of the Scraper.
 func (ScrapeSPARK) Help() string {
   return "SPARK Metrics"
 }
 
 // Version.
 func (ScrapeSPARK) Version() float64 {
   return 1.0
 }
 
 // Scrape generic function. Override for host module.
 func (ScrapeSPARK) Scrape (ctx context.Context, config *Collector_connection_data, ch chan<- prometheus.Metric) error {
   log.Debug_msg("Ejecutando SPARK Metrics Scraper")
 
   // Queries counters
   success_queries := 0
   error_queries := 0
 
   // Execute the generic funtion for creation of metrics with the pairs (QUERY, PROM:DESCRIPTOR)
   for i:=0 ; i < len(spark_query_variable_relationship) ; i++ {
     if create_spark_metric(ctx, *config, spark_query_variable_relationship[i].Query, spark_query_variable_relationship[i].Metric_struct, ch) {
       success_queries += 1
     } else {
       error_queries += 1
     }
   }
   log.Debug_msg("In the SPARK Module has been executed %d queries. %d success and %d with errors", success_queries + error_queries, success_queries, error_queries)
   return nil
 }
 
 var _ Scraper = ScrapeSPARK{}
 