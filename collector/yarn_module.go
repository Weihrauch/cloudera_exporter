package collector

import (
	// Go Default libraries
	"context"
	"strings"
	"sync"

	// Own libraries
	"fmt"
	jp "keedio/cloudera_exporter/json_parser"
	log "keedio/cloudera_exporter/logger"
	// Go Prometheus libraries
	"github.com/prometheus/client_golang/prometheus"
)

type YarnMetric struct {
	name          string
	description   string
	query         string
	metric_struct prometheus.Desc
	tags          []string
}

var yarnMetrics = []YarnMetric{
	YarnMetric{
		name:        "yarn_gc_time_across_nodemanagers",
		query:       `SELECT integral(jvm_gc_time_ms_rate_across_nodemanagers) WHERE entityName=yarn AND category=SERVICE`,
		description: "YARN GC time across nodemanagers",
		tags:        DEFAULT_METRICS_TAGS,
	},
	YarnMetric{
		name:        "yarn_total_apps_running",
		query:       `SELECT total_apps_running_across_yarn_pools WHERE entityName = "yarn" AND category = SERVICE`,
		description: "YARN total apps running",
		tags:        DEFAULT_METRICS_TAGS,
	},
	YarnMetric{
		name:        "yarn_allocated_vcores",
		query:       `SELECT allocated_vcores WHERE entityName RLIKE "yarn:root\..*"`,
		description: "YARN allocated vcores",
		tags:        QUEUE_METRIC_TAGS,
	},
	YarnMetric{
		name:        "yarn_allocated_memory_mb",
		query:       `SELECT allocated_memory_mb WHERE entityName RLIKE "yarn:root\..*"`,
		description: "YARN allocated vcores",
		tags:        QUEUE_METRIC_TAGS,
	},
	YarnMetric{
		name:        "available_vcores",
		query:       `SELECT available_vcores WHERE category=YARN_POOL and serviceName=yarn and queueName=root`,
		description: "YARN available vcores",
		tags:        QUEUE_METRIC_TAGS,
	},
	YarnMetric{
		name:        "available_memory_mb",
		query:       `SELECT available_memory_mb WHERE category=YARN_POOL and serviceName=yarn and queueName=root`,
		description: "YARN available mamory",
		tags:        QUEUE_METRIC_TAGS,
	},

	YarnMetric{
		name:        "total_containers_running_across_nodemanagers",
		query:       `SELECT total_containers_running_across_nodemanagers WHERE entityName=yarn AND category=SERVICE`,
		description: "YARN total containers running",
		tags:        DEFAULT_METRICS_TAGS,
	},

	YarnMetric{
		name:        "total_apps_pending",
		query:       `SELECT total_apps_pending_across_yarn_pools WHERE entityName = "yarn" AND category = SERVICE`,
		description: "YARN total_apps_pending_across_yarn_pools",
		tags:        DEFAULT_METRICS_TAGS,
	},

	YarnMetric{
		name:        "total_pending_containers",
		query:       `SELECT total_pending_containers_across_yarn_pools  WHERE entityName = "yarn" AND category = SERVICE`,
		description: "YARN total_pending_containers_across_yarn_pools ",
		tags:        DEFAULT_METRICS_TAGS,
	},
}
var (
	QUEUE_METRIC_TAGS    = []string{"clusterName", "entityName", "queueName"}
	DEFAULT_METRICS_TAGS = []string{"clusterName", "entityName", "hostname"}
)

const YARN_SCRAPER_NAME = "yarn"

func (r *YarnMetric) create_metric_struct() *prometheus.Desc {
	// Correct "description" parameter if is empty
	if len(r.description) == 0 {
		r.description = strings.Replace(strings.ToUpper(r.name), "_", " ", 0)
	}

	// return prometheus descriptor
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, YARN_SCRAPER_NAME, r.name),
		r.description, r.tags,
		nil,
	)
}

// Only for HDFS metric type
func (r *YarnMetric) create_yarn_metric(ctx context.Context, config Collector_connection_data, ch chan<- prometheus.Metric) bool {

	json_parsed, err := make_and_parse_timeseries_query(ctx, config, r.query)
	if err != nil {
		return false
	}

	// Get the num of hosts in the cluster or clusters
	num_ts_series, err := jp.Get_timeseries_num(&json_parsed)
	if err != nil {
		return false
	}
	metric_struct := r.create_metric_struct()
	// Extract Metadata for each TimeSerie
	for ts_index := 0; ts_index < num_ts_series; ts_index++ {
		//		fmt.Printf("json_parsed = %+v\n", json_parsed)
		// Get the Cluster Name
		var tagValues []string
		for _, tag := range r.tags {
			tagValues = append(tagValues, jp.Get_timeseries_attribute(&json_parsed, ts_index, tag))
		}
		fmt.Printf("tagValues = %+v\n", tagValues)
		// Get Query LAST value
		value, err := jp.Get_timeseries_query_value(&json_parsed, ts_index)
		if err != nil {
			continue
		}

		// Assing the data to the Prometheus descriptor
		ch <- prometheus.MustNewConstMetric(metric_struct, prometheus.GaugeValue, value, tagValues...)
	}
	return true

}

// Only for HDFS metric type
func create_yarn_metric(ctx context.Context, config Collector_connection_data, query string, metric_struct prometheus.Desc, ch chan<- prometheus.Metric) bool {
	// Make the query
	json_parsed, err := make_and_parse_timeseries_query(ctx, config, query)
	if err != nil {
		return false
	}

	// Get the num of hosts in the cluster or clusters
	num_ts_series, err := jp.Get_timeseries_num(&json_parsed)
	if err != nil {
		return false
	}

	// Extract Metadata for each TimeSerie
	for ts_index := 0; ts_index < num_ts_series; ts_index++ {
		// Get the Cluster Name
		cluster_name := jp.Get_timeseries_query_cluster(json_parsed, ts_index)
		// Get Host Name
		host_name := jp.Get_timeseries_query_host_name(json_parsed, ts_index)
		// Get the entity Name
		entity_name := jp.Get_timeseries_query_entity_name(json_parsed, ts_index)
		// Get Query LAST value
		value, err := jp.Get_timeseries_query_value(&json_parsed, ts_index)
		if err != nil {
			continue
		}
		// Assing the data to the Prometheus descriptor
		ch <- prometheus.MustNewConstMetric(&metric_struct, prometheus.GaugeValue, value, cluster_name, entity_name, host_name)
	}
	return true
}

// Metric descriptors.
var (
	myNewYARNMetric = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "subsystem", "metric_name"),
		"This is my metrics description.",
		[]string{"label1", "label2", "label3"}, nil)
)

// ScrapeGlobalStatus collects from /clusters/hosts.
type ScrapeYARNMetrics struct{}

// Name of the Scraper. Should be unique.
func (ScrapeYARNMetrics) Name() string {
	return "YARN"
}

// Help describes the role of the Scraper.
func (ScrapeYARNMetrics) Help() string {
	return "Collect YARN Service Metrics"
}

// Version.
func (ScrapeYARNMetrics) Version() float64 {
	return 1.0
}

func (ScrapeYARNMetrics) Scrape(ctx context.Context, config *Collector_connection_data, ch chan<- prometheus.Metric) error {

	log.Debug_msg("Execute YARN Metrics Scraper")

	// Queries counters
	//success_queries := 0
	//error_queries := 0
	var wg sync.WaitGroup
	wg.Add(len(yarnMetrics))
	for _, element := range yarnMetrics {
		go func(element YarnMetric) {
			element.create_yarn_metric(ctx, *config, ch)
			defer wg.Done()
		}(element)
	}
	wg.Wait()
	// Execute the generic funtion for creation of metrics with the pairs (QUERY, PROM:DESCRIPTOR)
	//		for i := 0; i < len(yarn_query_variable_relationship); i++ {
	//			if create_yarn_metric(ctx, *config, yarn_query_variable_relationship[i].Query, yarn_query_variable_relationship[i].Metric_struct, ch) {
	//				success_queries += 1
	//			} else {
	//				error_queries += 1
	//			}
	//		}
	//		log.Debug_msg("In the YARN Module has been executed %d queries. %d success and %d with errors", success_queries+error_queries, success_queries, error_queries)
	return nil
}

// check interface
var _ Scraper = ScrapeYARNMetrics{}
