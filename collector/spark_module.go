/*
 *
 * title           :collector/spark_module.go
 * description     :Submodule Collector for the Cluster SPARK metrics
 * author               :NTUMBA Phinées
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
 )
 
 
 /* ======================================================================
  * Constants with the Host module TSquery sentences
  * ====================================================================== */
  const SPARK_SCRAPER_NAME = "spark"
  var (
    // Agent Queries
    SPARK_ALERTS_RATE                 ="SELECT LAST(alerts_rate)  WHERE serviceType=\"SPARK\""
    SPARK_EVENTS_CRITICAL_RATE        ="SELECT LAST(events_critical_rate)  WHERE serviceType=\"SPARK\""
    SPARK_EVENTS_IMPORTANT_RATE       ="SELECT LAST(events_important_rate)  WHERE serviceType=\"SPARK\""
    SPARK_EVENTS_INFORMATIONAL_RATE   ="SELECT LAST(events_informational_rate)  WHERE serviceType=\"SPARK\""
    SPARK_HEALTH_BAD_RATE             ="SELECT LAST(health_bad_rate)  WHERE serviceType=\"SPARK\""
    SPARK_HEALTH_CONCERNING_RATE      ="SELECT LAST(health_concerning_rate)  WHERE serviceType=\"SPARK\""
    SPARK_HEALTH_DISABLED_RATE        ="SELECT LAST(health_disabled_rate)  WHERE serviceType=\"SPARK\""
    SPARK_HEALTH_GOOD_RATE            ="SELECT LAST(health_good_rate)  WHERE serviceType=\"SPARK\""
    SPARK_HEALTH_UNKNOWN_RATE         ="SELECT LAST(health_unknown_rate)  WHERE serviceType=\"SPARK\""  
  )
 
 
 
 
 /* ======================================================================
  * Global variables
  * ====================================================================== */
 var (
  spark_alerts_rate                 =create_spark_metric_struct("spark_alerts_rate",	"The number of alerts.	events per second	CDH 5, CDH 6")
  spark_events_critical_rate        =create_spark_metric_struct("spark_events_critical_rate",	"The number of critical events.	events per second	CDH 5, CDH 6")
  spark_events_important_rate       =create_spark_metric_struct("spark_events_important_rate",	"The number of important events.	events per second	CDH 5, CDH 6")
  spark_events_informational_rate   =create_spark_metric_struct("spark_events_informational_rate",	"The number of informational events.	events per second	CDH 5, CDH 6")
  spark_health_bad_rate             =create_spark_metric_struct("spark_health_bad_rate",	"Percentage of Time with Bad Health	seconds per second	CDH 5, CDH 6")
  spark_health_concerning_rate      =create_spark_metric_struct("spark_health_concerning_rate",	"Percentage of Time with Concerning Health	seconds per second	CDH 5, CDH 6")
  spark_health_disabled_rate        =create_spark_metric_struct("spark_health_disabled_rate",	"Percentage of Time with Disabled Health	seconds per second	CDH 5, CDH 6")
  spark_health_good_rate            =create_spark_metric_struct("spark_health_good_rate",	"Percentage of Time with Good Health	seconds per second	CDH 5, CDH 6")
  spark_health_unknown_rate         =create_spark_metric_struct("spark_health_unknown_rate",	"Percentage of Time with Unknown Health	seconds per second	CDH 5, CDH 6") 
 )
 var spark_query_variable_relationship = []relation {
  {SPARK_ALERTS_RATE,                 *spark_alerts_rate},
  {SPARK_EVENTS_CRITICAL_RATE,        *spark_events_critical_rate},
  {SPARK_EVENTS_IMPORTANT_RATE,       *spark_events_important_rate},
  {SPARK_EVENTS_INFORMATIONAL_RATE,   *spark_events_informational_rate},
  {SPARK_HEALTH_BAD_RATE,             *spark_health_bad_rate},
  {SPARK_HEALTH_CONCERNING_RATE,      *spark_health_concerning_rate},
  {SPARK_HEALTH_DISABLED_RATE,        *spark_health_disabled_rate},
  {SPARK_HEALTH_GOOD_RATE,            *spark_health_good_rate},
  {SPARK_HEALTH_UNKNOWN_RATE,         *spark_health_unknown_rate},
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
   log.Info_msg("In the SPARK Module has been executed %d queries. %d success and %d with errors", success_queries + error_queries, success_queries, error_queries)
   return nil
 }
 
 var _ Scraper = ScrapeSPARK{}
 