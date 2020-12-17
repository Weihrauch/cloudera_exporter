/*
 *
 * title           :collector/sentry_module.go
 * description     :Submodule Collector for the Cluster SENTRY metrics
 * author		       :NTUMBA Phinées
 * date            :2020/12/08
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
 const SENTRY_SCRAPER_NAME = "sentry"
 var (
   // Agent Queries
   SENTRY_ALERTS_RATE                       ="SELECT LAST(alerts_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_CGROUP_CPU_SYSTEM_RATE            ="SELECT LAST(cgroup_cpu_system_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_CGROUP_CPU_USER_RATE              ="SELECT LAST(cgroup_cpu_user_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_CGROUP_MEM_PAGE_CACHE             ="SELECT LAST(cgroup_mem_page_cache) WHERE serviceType=\"SENTRY\""
   SENTRY_CGROUP_MEM_RSS                    ="SELECT LAST(cgroup_mem_rss) WHERE serviceType=\"SENTRY\""
   SENTRY_CGROUP_MEM_SWAP                   ="SELECT LAST(cgroup_mem_swap) WHERE serviceType=\"SENTRY\""
   SENTRY_CGROUP_READ_BYTES_RATE            ="SELECT LAST(cgroup_read_bytes_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_CGROUP_READ_IOS_RATE              ="SELECT LAST(cgroup_read_ios_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_CGROUP_WRITE_BYTES_RATE           ="SELECT LAST(cgroup_write_bytes_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_CGROUP_WRITE_IOS_RATE             ="SELECT LAST(cgroup_write_ios_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_CPU_SYSTEM_RATE                   ="SELECT LAST(cpu_system_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_CPU_USER_RATE                     ="SELECT LAST(cpu_user_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_EVENTS_CRITICAL_RATE              ="SELECT LAST(events_critical_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_EVENTS_IMPORTANT_RATE             ="SELECT LAST(events_important_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_EVENTS_INFORMATIONAL_RATE         ="SELECT LAST(events_informational_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_FD_MAX                            ="SELECT LAST(fd_max) WHERE serviceType=\"SENTRY\""
   SENTRY_FD_OPEN                           ="SELECT LAST(fd_open) WHERE serviceType=\"SENTRY\""
   SENTRY_HEALTH_BAD_RATE                   ="SELECT LAST(health_bad_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_HEALTH_CONCERNING_RATE            ="SELECT LAST(health_concerning_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_HEALTH_DISABLED_RATE              ="SELECT LAST(health_disabled_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_HEALTH_GOOD_RATE                  ="SELECT LAST(health_good_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_HEALTH_UNKNOWN_RATE               ="SELECT LAST(health_unknown_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_MEM_RSS                           ="SELECT LAST(mem_rss) WHERE serviceType=\"SENTRY\""
   SENTRY_MEM_SWAP                          ="SELECT LAST(mem_swap) WHERE serviceType=\"SENTRY\""
   SENTRY_MEM_VIRTUAL                       ="SELECT LAST(mem_virtual) WHERE serviceType=\"SENTRY\""
   SENTRY_OOM_EXITS_RATE                    ="SELECT LAST(oom_exits_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_READ_BYTES_RATE                   ="SELECT LAST(read_bytes_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_UNEXPECTED_EXITS_RATE             ="SELECT LAST(unexpected_exits_rate) WHERE serviceType=\"SENTRY\""
   SENTRY_UPTIME                            ="SELECT LAST(uptime) WHERE serviceType=\"SENTRY\""
   SENTRY_WRITE_BYTES_RATE                  ="SELECT LAST(write_bytes_rate) WHERE serviceType=\"SENTRY\""
 )
 
 
 
 
 /* ======================================================================
  * Global variables
  * ====================================================================== */
 // Prometheus data Descriptors for the metrics to export
 var (
   // Agent Metrics
   sentry_alerts_rate                   =create_sentry_metric_struct("sentry_alerts_rate", "The number of alerts.	events per second")
   sentry_cgroup_cpu_system_rate        =create_sentry_metric_struct("sentry_cgroup_cpu_system_rate", "CPU usage of the role's cgroup	seconds per second")
   sentry_cgroup_cpu_user_rate          =create_sentry_metric_struct("sentry_cgroup_cpu_user_rate", "User Space CPU usage of the role's cgroup	seconds per second")
   sentry_cgroup_mem_page_cache         =create_sentry_metric_struct("sentry_cgroup_mem_page_cache", "Page cache usage of the role's cgroup	bytes")
   sentry_cgroup_mem_rss                =create_sentry_metric_struct("sentry_cgroup_mem_rss", "Resident memory of the role's cgroup	bytes")
   sentry_cgroup_mem_swap               =create_sentry_metric_struct("sentry_cgroup_mem_swap", "Swap usage of the role's cgroup	bytes")
   sentry_cgroup_read_bytes_rate        =create_sentry_metric_struct("sentry_cgroup_read_bytes_rate", "Bytes read from all disks by the role's cgroup	bytes per second")
   sentry_cgroup_read_ios_rate          =create_sentry_metric_struct("sentry_cgroup_read_ios_rate", "Number of read I/O operations from all disks by the role's cgroup	ios per second")
   sentry_cgroup_write_bytes_rate       =create_sentry_metric_struct("sentry_cgroup_write_bytes_rate", "Bytes written to all disks by the role's cgroup	bytes per second")
   sentry_cgroup_write_ios_rate         =create_sentry_metric_struct("sentry_cgroup_write_ios_rate", "Number of write I/O operations to all disks by the role's cgroup	ios per second")
   sentry_cpu_system_rate               =create_sentry_metric_struct("sentry_cpu_system_rate", "Total System CPU	seconds per second")
   sentry_cpu_user_rate                 =create_sentry_metric_struct("sentry_cpu_user_rate", "Total CPU user time	seconds per second")
   sentry_events_critical_rate          =create_sentry_metric_struct("sentry_events_critical_rate", "The number of critical events.	events per second")
   sentry_events_important_rate         =create_sentry_metric_struct("sentry_events_important_rate", "The number of important events.	events per second")
   sentry_events_informational_rate     =create_sentry_metric_struct("sentry_events_informational_rate", "The number of informational events.	events per second")
   sentry_fd_max                        =create_sentry_metric_struct("sentry_fd_max", "Maximum number of file descriptors	file descriptors")
   sentry_fd_open                       =create_sentry_metric_struct("sentry_fd_open", "Open file descriptors.	file descriptors")
   sentry_health_bad_rate               =create_sentry_metric_struct("sentry_health_bad_rate", "Percentage of Time with Bad Health	seconds per second")
   sentry_health_concerning_rate        =create_sentry_metric_struct("sentry_health_concerning_rate", "Percentage of Time with Concerning Health	seconds per second")
   sentry_health_disabled_rate          =create_sentry_metric_struct("sentry_health_disabled_rate", "Percentage of Time with Disabled Health	seconds per second")
   sentry_health_good_rate              =create_sentry_metric_struct("sentry_health_good_rate", "Percentage of Time with Good Health	seconds per second")
   sentry_health_unknown_rate           =create_sentry_metric_struct("sentry_health_unknown_rate", "Percentage of Time with Unknown Health	seconds per second")
   sentry_mem_rss                       =create_sentry_metric_struct("sentry_mem_rss", "Resident memory used	bytes")
   sentry_mem_swap                      =create_sentry_metric_struct("sentry_mem_swap", "Amount of swap memory used by this role's process.	bytes")
   sentry_mem_virtual                   =create_sentry_metric_struct("sentry_mem_virtual", "Virtual memory used	bytes")
   sentry_oom_exits_rate                =create_sentry_metric_struct("sentry_oom_exits_rate", "The number of times the role's backing process was killed due to an OutOfMemory error. This counter is only incremented if the Cloudera Manager 'Kill When Out of Memory' option is enabled.	exits per second")
   sentry_read_bytes_rate               =create_sentry_metric_struct("sentry_read_bytes_rate", "The number of bytes read from the device	bytes per second")
   sentry_unexpected_exits_rate         =create_sentry_metric_struct("sentry_unexpected_exits_rate", "The number of times the role's backing process exited unexpectedly.	exits per second")
   sentry_uptime                        =create_sentry_metric_struct("sentry_uptime", "For a host, the amount of time since the host was booted. For a role, the uptime of the backing process.	seconds")
   sentry_write_bytes_rate              =create_sentry_metric_struct("sentry_write_bytes_rate", "The number of bytes written to the device	bytes per second")
 )
 
 // Creation of the structure that relates the queries with the descriptors of the Prometheus metrics
 var sentry_query_variable_relationship = []relation {
   {SENTRY_ALERTS_RATE,                   *sentry_alerts_rate},
   {SENTRY_CGROUP_CPU_SYSTEM_RATE,        *sentry_cgroup_cpu_system_rate},
   {SENTRY_CGROUP_CPU_USER_RATE,          *sentry_cgroup_cpu_user_rate},
   {SENTRY_CGROUP_MEM_PAGE_CACHE,         *sentry_cgroup_mem_page_cache},
   {SENTRY_CGROUP_MEM_RSS,                *sentry_cgroup_mem_rss},
   {SENTRY_CGROUP_MEM_SWAP,               *sentry_cgroup_mem_swap},
   {SENTRY_CGROUP_READ_BYTES_RATE,        *sentry_cgroup_read_bytes_rate},
   {SENTRY_CGROUP_READ_IOS_RATE,          *sentry_cgroup_read_ios_rate},
   {SENTRY_CGROUP_WRITE_BYTES_RATE,       *sentry_cgroup_write_bytes_rate},
   {SENTRY_CGROUP_WRITE_IOS_RATE,         *sentry_cgroup_write_ios_rate},
   {SENTRY_CPU_SYSTEM_RATE,               *sentry_cpu_system_rate},
   {SENTRY_CPU_USER_RATE,                 *sentry_cpu_user_rate},
   {SENTRY_EVENTS_CRITICAL_RATE,          *sentry_events_critical_rate},
   {SENTRY_EVENTS_IMPORTANT_RATE,         *sentry_events_important_rate},
   {SENTRY_EVENTS_INFORMATIONAL_RATE,     *sentry_events_informational_rate},
   {SENTRY_FD_MAX,                        *sentry_fd_max},
   {SENTRY_FD_OPEN,                       *sentry_fd_open},
   {SENTRY_HEALTH_BAD_RATE,               *sentry_health_bad_rate},
   {SENTRY_HEALTH_CONCERNING_RATE,        *sentry_health_concerning_rate},
   {SENTRY_HEALTH_DISABLED_RATE,          *sentry_health_disabled_rate},
   {SENTRY_HEALTH_GOOD_RATE,              *sentry_health_good_rate},
   {SENTRY_HEALTH_UNKNOWN_RATE,           *sentry_health_unknown_rate},
   {SENTRY_MEM_RSS,                       *sentry_mem_rss},
   {SENTRY_MEM_SWAP,                      *sentry_mem_swap},
   {SENTRY_MEM_VIRTUAL,                   *sentry_mem_virtual},
   {SENTRY_OOM_EXITS_RATE,                *sentry_oom_exits_rate},
   {SENTRY_READ_BYTES_RATE,               *sentry_read_bytes_rate},
   {SENTRY_UNEXPECTED_EXITS_RATE,         *sentry_unexpected_exits_rate},
   {SENTRY_UPTIME,                        *sentry_uptime},
   {SENTRY_WRITE_BYTES_RATE,              *sentry_write_bytes_rate},
 }
 
 
 
 
 /* ======================================================================
  * Functions
  * ====================================================================== */
 // Create and returns a prometheus descriptor for a sentry metric. 
 // The "metric_name" parameter its mandatory
 // If the "description" parameter is empty, the function assings it with the
 // value of the name of the metric in uppercase and separated by spaces
 func create_sentry_metric_struct(metric_name string, description string) *prometheus.Desc {
   // Correct "description" parameter if is empty
   if len(description) == 0 {
     description = strings.Replace(strings.ToUpper(metric_name), "_", " ", 0)
   }
 
   // return prometheus descriptor
   return prometheus.NewDesc(
     prometheus.BuildFQName(namespace, SENTRY_SCRAPER_NAME, metric_name),
     description,
     []string{"cluster", "entityName","hostname"},
     nil,
   )
 }
 
 
 // Generic function to extract de metadata associated with the query value
 // Only for SENTRY metric type
 func create_sentry_metric (ctx context.Context, config Collector_connection_data, query string, metric_struct prometheus.Desc, ch chan<- prometheus.Metric) bool {
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
 // ScrapeSENTRY struct
 type ScrapeSENTRY struct{}
 
 // Name of the Scraper. Should be unique.
 func (ScrapeSENTRY) Name() string {
   return SENTRY_SCRAPER_NAME
 }
 
 // Help describes the role of the Scraper.
 func (ScrapeSENTRY) Help() string {
   return "SENTRY Metrics"
 }
 
 // Version.
 func (ScrapeSENTRY) Version() float64 {
   return 1.0
 }
 
 // Scrape generic function. Override for host module.
 func (ScrapeSENTRY) Scrape (ctx context.Context, config *Collector_connection_data, ch chan<- prometheus.Metric) error {
   log.Debug_msg("Executing SENTRY Metrics Scraper")
 
   // Queries counters
   success_queries := 0
   error_queries := 0
 
   // Execute the generic funtion for creation of metrics with the pairs (QUERY, PROM:DESCRIPTOR)
   for i:=0 ; i < len(sentry_query_variable_relationship) ; i++ {
     if create_sentry_metric(ctx, *config, sentry_query_variable_relationship[i].Query, sentry_query_variable_relationship[i].Metric_struct, ch) {
       success_queries += 1
     } else {
       error_queries += 1
     }
   }
   log.Info_msg("In the SENTRY Module has been executed %d queries. %d success and %d with errors", success_queries + error_queries, success_queries, error_queries)
   return nil
 }
 
 var _ Scraper = ScrapeSENTRY{}
 