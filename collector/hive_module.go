/*
 *
 * title           :collector/hive_module.go
 * description     :Submodule Collector for the Cluster HIVE metrics
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
const HIVE_SCRAPER_NAME = "hive"
var (
  // Agent Queries
  HIVE_ALERTS_RATE                                     ="SELECT LAST(alerts_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_ACQUIREREADWRITELOCKS_AVG                   ="SELECT LAST(hive_api_acquirereadwritelocks_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_ACQUIREREADWRITELOCKS_MAX                   ="SELECT LAST(hive_api_acquirereadwritelocks_max) WHERE serviceType=\"HIVE\""
  HIVE_API_ACQUIREREADWRITELOCKS_MEDIAN                ="SELECT LAST(hive_api_acquirereadwritelocks_median) WHERE serviceType=\"HIVE\""
  HIVE_API_ACQUIREREADWRITELOCKS_MIN                   ="SELECT LAST(hive_api_acquirereadwritelocks_min) WHERE serviceType=\"HIVE\""
  HIVE_API_ACQUIREREADWRITELOCKS_RATE                  ="SELECT LAST(hive_api_acquirereadwritelocks_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_ACQUIREREADWRITELOCKS_STDDEV                ="SELECT LAST(hive_api_acquirereadwritelocks_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_COMPILE_AVG                                 ="SELECT LAST(hive_api_compile_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_COMPILE_MAX                                 ="SELECT LAST(hive_api_compile_max) WHERE serviceType=\"HIVE\""
  HIVE_API_COMPILE_MEDIAN                              ="SELECT LAST(hive_api_compile_median) WHERE serviceType=\"HIVE\""
  HIVE_API_COMPILE_MIN                                 ="SELECT LAST(hive_api_compile_min) WHERE serviceType=\"HIVE\""
  HIVE_API_COMPILE_RATE                                ="SELECT LAST(hive_api_compile_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_COMPILE_STDDEV                              ="SELECT LAST(hive_api_compile_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_DRIVER_EXECUTE_AVG                          ="SELECT LAST(hive_api_driver_execute_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_DRIVER_EXECUTE_MAX                          ="SELECT LAST(hive_api_driver_execute_max) WHERE serviceType=\"HIVE\""
  HIVE_API_DRIVER_EXECUTE_MEDIAN                       ="SELECT LAST(hive_api_driver_execute_median) WHERE serviceType=\"HIVE\""
  HIVE_API_DRIVER_EXECUTE_MIN                          ="SELECT LAST(hive_api_driver_execute_min) WHERE serviceType=\"HIVE\""
  HIVE_API_DRIVER_EXECUTE_RATE                         ="SELECT LAST(hive_api_driver_execute_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_DRIVER_EXECUTE_STDDEV                       ="SELECT LAST(hive_api_driver_execute_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_DRIVER_RUN_AVG                              ="SELECT LAST(hive_api_driver_run_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_DRIVER_RUN_MAX                              ="SELECT LAST(hive_api_driver_run_max) WHERE serviceType=\"HIVE\""
  HIVE_API_DRIVER_RUN_MEDIAN                           ="SELECT LAST(hive_api_driver_run_median) WHERE serviceType=\"HIVE\""
  HIVE_API_DRIVER_RUN_MIN                              ="SELECT LAST(hive_api_driver_run_min) WHERE serviceType=\"HIVE\""
  HIVE_API_DRIVER_RUN_RATE                             ="SELECT LAST(hive_api_driver_run_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_DRIVER_RUN_STDDEV                           ="SELECT LAST(hive_api_driver_run_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_ALL_TABLES_AVG                          ="SELECT LAST(hive_api_get_all_tables_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_ALL_TABLES_MAX                          ="SELECT LAST(hive_api_get_all_tables_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_ALL_TABLES_MEDIAN                       ="SELECT LAST(hive_api_get_all_tables_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_ALL_TABLES_MIN                          ="SELECT LAST(hive_api_get_all_tables_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_ALL_TABLES_RATE                         ="SELECT LAST(hive_api_get_all_tables_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_ALL_TABLES_STDDEV                       ="SELECT LAST(hive_api_get_all_tables_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_CONFIG_VALUE_AVG                        ="SELECT LAST(hive_api_get_config_value_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_CONFIG_VALUE_MAX                        ="SELECT LAST(hive_api_get_config_value_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_CONFIG_VALUE_MEDIAN                     ="SELECT LAST(hive_api_get_config_value_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_CONFIG_VALUE_MIN                        ="SELECT LAST(hive_api_get_config_value_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_CONFIG_VALUE_RATE                       ="SELECT LAST(hive_api_get_config_value_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_CONFIG_VALUE_STDDEV                     ="SELECT LAST(hive_api_get_config_value_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_DATABASE_AVG                            ="SELECT LAST(hive_api_get_database_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_DATABASE_MAX                            ="SELECT LAST(hive_api_get_database_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_DATABASE_MEDIAN                         ="SELECT LAST(hive_api_get_database_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_DATABASE_MIN                            ="SELECT LAST(hive_api_get_database_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_DATABASE_RATE                           ="SELECT LAST(hive_api_get_database_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_DATABASE_STDDEV                         ="SELECT LAST(hive_api_get_database_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_DATABASES_AVG                           ="SELECT LAST(hive_api_get_databases_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_DATABASES_MAX                           ="SELECT LAST(hive_api_get_databases_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_DATABASES_MEDIAN                        ="SELECT LAST(hive_api_get_databases_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_DATABASES_MIN                           ="SELECT LAST(hive_api_get_databases_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_DATABASES_RATE                          ="SELECT LAST(hive_api_get_databases_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_DATABASES_STDDEV                        ="SELECT LAST(hive_api_get_databases_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FIELDS_WITH_ENVIRONMENT_CONTEXT_AVG     ="SELECT LAST(hive_api_get_fields_with_environment_context_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FIELDS_WITH_ENVIRONMENT_CONTEXT_MAX     ="SELECT LAST(hive_api_get_fields_with_environment_context_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FIELDS_WITH_ENVIRONMENT_CONTEXT_MEDIAN  ="SELECT LAST(hive_api_get_fields_with_environment_context_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FIELDS_WITH_ENVIRONMENT_CONTEXT_MIN     ="SELECT LAST(hive_api_get_fields_with_environment_context_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FIELDS_WITH_ENVIRONMENT_CONTEXT_RATE    ="SELECT LAST(hive_api_get_fields_with_environment_context_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FIELDS_WITH_ENVIRONMENT_CONTEXT_STDDEV  ="SELECT LAST(hive_api_get_fields_with_environment_context_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FUNCTION_AVG                            ="SELECT LAST(hive_api_get_function_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FUNCTION_MAX                            ="SELECT LAST(hive_api_get_function_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FUNCTION_MEDIAN                         ="SELECT LAST(hive_api_get_function_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FUNCTION_MIN                            ="SELECT LAST(hive_api_get_function_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FUNCTION_RATE                           ="SELECT LAST(hive_api_get_function_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FUNCTION_STDDEV                         ="SELECT LAST(hive_api_get_function_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FUNCTIONS_AVG                           ="SELECT LAST(hive_api_get_functions_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FUNCTIONS_MAX                           ="SELECT LAST(hive_api_get_functions_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FUNCTIONS_MEDIAN                        ="SELECT LAST(hive_api_get_functions_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FUNCTIONS_MIN                           ="SELECT LAST(hive_api_get_functions_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FUNCTIONS_RATE                          ="SELECT LAST(hive_api_get_functions_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_FUNCTIONS_STDDEV                        ="SELECT LAST(hive_api_get_functions_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEX_BY_NAME_AVG                       ="SELECT LAST(hive_api_get_index_by_name_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEX_BY_NAME_MAX                       ="SELECT LAST(hive_api_get_index_by_name_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEX_BY_NAME_MEDIAN                    ="SELECT LAST(hive_api_get_index_by_name_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEX_BY_NAME_MIN                       ="SELECT LAST(hive_api_get_index_by_name_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEX_BY_NAME_RATE                      ="SELECT LAST(hive_api_get_index_by_name_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEX_BY_NAME_STDDEV                    ="SELECT LAST(hive_api_get_index_by_name_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEX_NAMES_AVG                         ="SELECT LAST(hive_api_get_index_names_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEX_NAMES_MAX                         ="SELECT LAST(hive_api_get_index_names_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEX_NAMES_MEDIAN                      ="SELECT LAST(hive_api_get_index_names_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEX_NAMES_MIN                         ="SELECT LAST(hive_api_get_index_names_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEX_NAMES_RATE                        ="SELECT LAST(hive_api_get_index_names_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEX_NAMES_STDDEV                      ="SELECT LAST(hive_api_get_index_names_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEXES_AVG                             ="SELECT LAST(hive_api_get_indexes_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEXES_MAX                             ="SELECT LAST(hive_api_get_indexes_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEXES_MEDIAN                          ="SELECT LAST(hive_api_get_indexes_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEXES_MIN                             ="SELECT LAST(hive_api_get_indexes_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEXES_RATE                            ="SELECT LAST(hive_api_get_indexes_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_INDEXES_STDDEV                          ="SELECT LAST(hive_api_get_indexes_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_MULTI_TABLE_AVG                         ="SELECT LAST(hive_api_get_multi_table_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_MULTI_TABLE_MAX                         ="SELECT LAST(hive_api_get_multi_table_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_MULTI_TABLE_MEDIAN                      ="SELECT LAST(hive_api_get_multi_table_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_MULTI_TABLE_MIN                         ="SELECT LAST(hive_api_get_multi_table_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_MULTI_TABLE_RATE                        ="SELECT LAST(hive_api_get_multi_table_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_MULTI_TABLE_STDDEV                      ="SELECT LAST(hive_api_get_multi_table_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLE_AVG                               ="SELECT LAST(hive_api_get_table_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLE_MAX                               ="SELECT LAST(hive_api_get_table_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLE_MEDIAN                            ="SELECT LAST(hive_api_get_table_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLE_MIN                               ="SELECT LAST(hive_api_get_table_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLE_NAMES_BY_FILTER_AVG               ="SELECT LAST(hive_api_get_table_names_by_filter_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLE_NAMES_BY_FILTER_MAX               ="SELECT LAST(hive_api_get_table_names_by_filter_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLE_NAMES_BY_FILTER_MEDIAN            ="SELECT LAST(hive_api_get_table_names_by_filter_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLE_NAMES_BY_FILTER_MIN               ="SELECT LAST(hive_api_get_table_names_by_filter_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLE_NAMES_BY_FILTER_RATE              ="SELECT LAST(hive_api_get_table_names_by_filter_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLE_NAMES_BY_FILTER_STDDEV            ="SELECT LAST(hive_api_get_table_names_by_filter_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLE_RATE                              ="SELECT LAST(hive_api_get_table_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLE_STDDEV                            ="SELECT LAST(hive_api_get_table_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLES_AVG                              ="SELECT LAST(hive_api_get_tables_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLES_MAX                              ="SELECT LAST(hive_api_get_tables_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLES_MEDIAN                           ="SELECT LAST(hive_api_get_tables_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLES_MIN                              ="SELECT LAST(hive_api_get_tables_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLES_RATE                             ="SELECT LAST(hive_api_get_tables_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GET_TABLES_STDDEV                           ="SELECT LAST(hive_api_get_tables_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GETINPUTSUMMARY_AVG                         ="SELECT LAST(hive_api_getinputsummary_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GETINPUTSUMMARY_MAX                         ="SELECT LAST(hive_api_getinputsummary_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GETINPUTSUMMARY_MEDIAN                      ="SELECT LAST(hive_api_getinputsummary_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GETINPUTSUMMARY_MIN                         ="SELECT LAST(hive_api_getinputsummary_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GETINPUTSUMMARY_RATE                        ="SELECT LAST(hive_api_getinputsummary_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GETINPUTSUMMARY_STDDEV                      ="SELECT LAST(hive_api_getinputsummary_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_GETSPLITS_AVG                               ="SELECT LAST(hive_api_getsplits_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_GETSPLITS_MAX                               ="SELECT LAST(hive_api_getsplits_max) WHERE serviceType=\"HIVE\""
  HIVE_API_GETSPLITS_MEDIAN                            ="SELECT LAST(hive_api_getsplits_median) WHERE serviceType=\"HIVE\""
  HIVE_API_GETSPLITS_MIN                               ="SELECT LAST(hive_api_getsplits_min) WHERE serviceType=\"HIVE\""
  HIVE_API_GETSPLITS_RATE                              ="SELECT LAST(hive_api_getsplits_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_GETSPLITS_STDDEV                            ="SELECT LAST(hive_api_getsplits_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_MARKPARTITIONFOREVENT_AVG                   ="SELECT LAST(hive_api_markpartitionforevent_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_MARKPARTITIONFOREVENT_MAX                   ="SELECT LAST(hive_api_markpartitionforevent_max) WHERE serviceType=\"HIVE\""
  HIVE_API_MARKPARTITIONFOREVENT_MEDIAN                ="SELECT LAST(hive_api_markpartitionforevent_median) WHERE serviceType=\"HIVE\""
  HIVE_API_MARKPARTITIONFOREVENT_MIN                   ="SELECT LAST(hive_api_markpartitionforevent_min) WHERE serviceType=\"HIVE\""
  HIVE_API_MARKPARTITIONFOREVENT_RATE                  ="SELECT LAST(hive_api_markpartitionforevent_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_MARKPARTITIONFOREVENT_STDDEV                ="SELECT LAST(hive_api_markpartitionforevent_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_INITIALIZED_AVG                   ="SELECT LAST(hive_api_operation_initialized_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_INITIALIZED_MAX                   ="SELECT LAST(hive_api_operation_initialized_max) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_INITIALIZED_MEDIAN                ="SELECT LAST(hive_api_operation_initialized_median) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_INITIALIZED_MIN                   ="SELECT LAST(hive_api_operation_initialized_min) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_INITIALIZED_RATE                  ="SELECT LAST(hive_api_operation_initialized_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_INITIALIZED_STDDEV                ="SELECT LAST(hive_api_operation_initialized_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_PENDING_AVG                       ="SELECT LAST(hive_api_operation_pending_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_PENDING_MAX                       ="SELECT LAST(hive_api_operation_pending_max) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_PENDING_MEDIAN                    ="SELECT LAST(hive_api_operation_pending_median) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_PENDING_MIN                       ="SELECT LAST(hive_api_operation_pending_min) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_PENDING_RATE                      ="SELECT LAST(hive_api_operation_pending_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_PENDING_STDDEV                    ="SELECT LAST(hive_api_operation_pending_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_RUNNING_AVG                       ="SELECT LAST(hive_api_operation_running_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_RUNNING_MAX                       ="SELECT LAST(hive_api_operation_running_max) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_RUNNING_MEDIAN                    ="SELECT LAST(hive_api_operation_running_median) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_RUNNING_MIN                       ="SELECT LAST(hive_api_operation_running_min) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_RUNNING_RATE                      ="SELECT LAST(hive_api_operation_running_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_OPERATION_RUNNING_STDDEV                    ="SELECT LAST(hive_api_operation_running_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_PARSE_AVG                                   ="SELECT LAST(hive_api_parse_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_PARSE_MAX                                   ="SELECT LAST(hive_api_parse_max) WHERE serviceType=\"HIVE\""
  HIVE_API_PARSE_MEDIAN                                ="SELECT LAST(hive_api_parse_median) WHERE serviceType=\"HIVE\""
  HIVE_API_PARSE_MIN                                   ="SELECT LAST(hive_api_parse_min) WHERE serviceType=\"HIVE\""
  HIVE_API_PARSE_RATE                                  ="SELECT LAST(hive_api_parse_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_PARSE_STDDEV                                ="SELECT LAST(hive_api_parse_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_PARTITION_RETRIEVING_AVG                    ="SELECT LAST(hive_api_partition_retrieving_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_PARTITION_RETRIEVING_MAX                    ="SELECT LAST(hive_api_partition_retrieving_max) WHERE serviceType=\"HIVE\""
  HIVE_API_PARTITION_RETRIEVING_MEDIAN                 ="SELECT LAST(hive_api_partition_retrieving_median) WHERE serviceType=\"HIVE\""
  HIVE_API_PARTITION_RETRIEVING_MIN                    ="SELECT LAST(hive_api_partition_retrieving_min) WHERE serviceType=\"HIVE\""
  HIVE_API_PARTITION_RETRIEVING_RATE                   ="SELECT LAST(hive_api_partition_retrieving_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_PARTITION_RETRIEVING_STDDEV                 ="SELECT LAST(hive_api_partition_retrieving_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_RELEASELOCKS_AVG                            ="SELECT LAST(hive_api_releaselocks_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_RELEASELOCKS_MAX                            ="SELECT LAST(hive_api_releaselocks_max) WHERE serviceType=\"HIVE\""
  HIVE_API_RELEASELOCKS_MEDIAN                         ="SELECT LAST(hive_api_releaselocks_median) WHERE serviceType=\"HIVE\""
  HIVE_API_RELEASELOCKS_MIN                            ="SELECT LAST(hive_api_releaselocks_min) WHERE serviceType=\"HIVE\""
  HIVE_API_RELEASELOCKS_RATE                           ="SELECT LAST(hive_api_releaselocks_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_RELEASELOCKS_STDDEV                         ="SELECT LAST(hive_api_releaselocks_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_RUNTASKS_AVG                                ="SELECT LAST(hive_api_runtasks_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_RUNTASKS_MAX                                ="SELECT LAST(hive_api_runtasks_max) WHERE serviceType=\"HIVE\""
  HIVE_API_RUNTASKS_MEDIAN                             ="SELECT LAST(hive_api_runtasks_median) WHERE serviceType=\"HIVE\""
  HIVE_API_RUNTASKS_MIN                                ="SELECT LAST(hive_api_runtasks_min) WHERE serviceType=\"HIVE\""
  HIVE_API_RUNTASKS_RATE                               ="SELECT LAST(hive_api_runtasks_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_RUNTASKS_STDDEV                             ="SELECT LAST(hive_api_runtasks_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_SERIALIZEPLAN_AVG                           ="SELECT LAST(hive_api_serializeplan_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_SERIALIZEPLAN_MAX                           ="SELECT LAST(hive_api_serializeplan_max) WHERE serviceType=\"HIVE\""
  HIVE_API_SERIALIZEPLAN_MEDIAN                        ="SELECT LAST(hive_api_serializeplan_median) WHERE serviceType=\"HIVE\""
  HIVE_API_SERIALIZEPLAN_MIN                           ="SELECT LAST(hive_api_serializeplan_min) WHERE serviceType=\"HIVE\""
  HIVE_API_SERIALIZEPLAN_RATE                          ="SELECT LAST(hive_api_serializeplan_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_SERIALIZEPLAN_STDDEV                        ="SELECT LAST(hive_api_serializeplan_stddev) WHERE serviceType=\"HIVE\""
  HIVE_API_TIMETOSUBMIT_AVG                            ="SELECT LAST(hive_api_timetosubmit_avg) WHERE serviceType=\"HIVE\""
  HIVE_API_TIMETOSUBMIT_MAX                            ="SELECT LAST(hive_api_timetosubmit_max) WHERE serviceType=\"HIVE\""
  HIVE_API_TIMETOSUBMIT_MEDIAN                         ="SELECT LAST(hive_api_timetosubmit_median) WHERE serviceType=\"HIVE\""
  HIVE_API_TIMETOSUBMIT_MIN                            ="SELECT LAST(hive_api_timetosubmit_min) WHERE serviceType=\"HIVE\""
  HIVE_API_TIMETOSUBMIT_RATE                           ="SELECT LAST(hive_api_timetosubmit_rate) WHERE serviceType=\"HIVE\""
  HIVE_API_TIMETOSUBMIT_STDDEV                         ="SELECT LAST(hive_api_timetosubmit_stddev) WHERE serviceType=\"HIVE\""
  HIVE_AUDIT_AGENT_BYTES_LEFT_RATE                     ="SELECT LAST(audit_agent_bytes_left_rate) WHERE serviceType=\"HIVE\""
  HIVE_AUDIT_AGENT_BYTES_SENT_RATE                     ="SELECT LAST(audit_agent_bytes_sent_rate) WHERE serviceType=\"HIVE\""
  HIVE_AUDIT_AGENT_EVENTS_SENT_RATE                    ="SELECT LAST(audit_agent_events_sent_rate) WHERE serviceType=\"HIVE\""
  HIVE_AUDIT_AGENT_SEND_FAILURES_RATE                  ="SELECT LAST(audit_agent_send_failures_rate) WHERE serviceType=\"HIVE\""
  HIVE_AUDIT_PLUGIN_COALESCED_RATE                     ="SELECT LAST(audit_plugin_coalesced_rate) WHERE serviceType=\"HIVE\""
  HIVE_AUDIT_PLUGIN_EVENTS_RATE                        ="SELECT LAST(audit_plugin_events_rate) WHERE serviceType=\"HIVE\""
  HIVE_AUDIT_PLUGIN_EXCEPTIONS_RATE                    ="SELECT LAST(audit_plugin_exceptions_rate) WHERE serviceType=\"HIVE\""
  HIVE_AUDIT_PLUGIN_FILTERED_RATE                      ="SELECT LAST(audit_plugin_filtered_rate) WHERE serviceType=\"HIVE\""
  HIVE_AUDIT_PLUGIN_PARSE_ERRORS_RATE                  ="SELECT LAST(audit_plugin_parse_errors_rate) WHERE serviceType=\"HIVE\""
  HIVE_CANARY_DURATION                                 ="SELECT LAST(canary_duration) WHERE serviceType=\"HIVE\""
  HIVE_CGROUP_CPU_SYSTEM_RATE                          ="SELECT LAST(cgroup_cpu_system_rate) WHERE serviceType=\"HIVE\""
  HIVE_CGROUP_CPU_USER_RATE                            ="SELECT LAST(cgroup_cpu_user_rate) WHERE serviceType=\"HIVE\""
  HIVE_CGROUP_MEM_PAGE_CACHE                           ="SELECT LAST(cgroup_mem_page_cache) WHERE serviceType=\"HIVE\""
  HIVE_CGROUP_MEM_RSS                                  ="SELECT LAST(cgroup_mem_rss) WHERE serviceType=\"HIVE\""
  HIVE_CGROUP_MEM_SWAP                                 ="SELECT LAST(cgroup_mem_swap) WHERE serviceType=\"HIVE\""
  HIVE_CGROUP_READ_BYTES_RATE                          ="SELECT LAST(cgroup_read_bytes_rate) WHERE serviceType=\"HIVE\""
  HIVE_CGROUP_READ_IOS_RATE                            ="SELECT LAST(cgroup_read_ios_rate) WHERE serviceType=\"HIVE\""
  HIVE_CGROUP_WRITE_BYTES_RATE                         ="SELECT LAST(cgroup_write_bytes_rate) WHERE serviceType=\"HIVE\""
  HIVE_CGROUP_WRITE_IOS_RATE                           ="SELECT LAST(cgroup_write_ios_rate) WHERE serviceType=\"HIVE\""
  HIVE_COMPLETED_OPERATION_CANCELED_RATE               ="SELECT LAST(hive_completed_operation_canceled_rate) WHERE serviceType=\"HIVE\""
  HIVE_COMPLETED_OPERATION_CLOSED_RATE                 ="SELECT LAST(hive_completed_operation_closed_rate) WHERE serviceType=\"HIVE\""
  HIVE_COMPLETED_OPERATION_ERROR_RATE                  ="SELECT LAST(hive_completed_operation_error_rate) WHERE serviceType=\"HIVE\""
  HIVE_COMPLETED_OPERATION_FINISHED_RATE               ="SELECT LAST(hive_completed_operation_finished_rate) WHERE serviceType=\"HIVE\""
  HIVE_CPU_SYSTEM_RATE                                 ="SELECT LAST(cpu_system_rate) WHERE serviceType=\"HIVE\""
  HIVE_CPU_USER_RATE                                   ="SELECT LAST(cpu_user_rate) WHERE serviceType=\"HIVE\""
  HIVE_CREATE_TOTAL_COUNT_DBS_RATE                     ="SELECT LAST(hive_create_total_count_dbs_rate) WHERE serviceType=\"HIVE\""
  HIVE_CREATE_TOTAL_COUNT_PARTITIONS_RATE              ="SELECT LAST(hive_create_total_count_partitions_rate) WHERE serviceType=\"HIVE\""
  HIVE_CREATE_TOTAL_COUNT_TABLES_RATE                  ="SELECT LAST(hive_create_total_count_tables_rate) WHERE serviceType=\"HIVE\""
  HIVE_DELETE_TOTAL_COUNT_DBS_RATE                     ="SELECT LAST(hive_delete_total_count_dbs_rate) WHERE serviceType=\"HIVE\""
  HIVE_DELETE_TOTAL_COUNT_PARTITIONS_RATE              ="SELECT LAST(hive_delete_total_count_partitions_rate) WHERE serviceType=\"HIVE\""
  HIVE_DELETE_TOTAL_COUNT_TABLES_RATE                  ="SELECT LAST(hive_delete_total_count_tables_rate) WHERE serviceType=\"HIVE\""
  HIVE_EVENTS_CRITICAL_RATE                            ="SELECT LAST(events_critical_rate) WHERE serviceType=\"HIVE\""
  HIVE_EVENTS_IMPORTANT_RATE                           ="SELECT LAST(events_important_rate) WHERE serviceType=\"HIVE\""
  HIVE_EVENTS_INFORMATIONAL_RATE                       ="SELECT LAST(events_informational_rate) WHERE serviceType=\"HIVE\""
  HIVE_EXEC_ASYNC_POOL_SIZE                            ="SELECT LAST(hive_exec_async_pool_size) WHERE serviceType=\"HIVE\""
  HIVE_EXEC_ASYNC_QUEUE_SIZE                           ="SELECT LAST(hive_exec_async_queue_size) WHERE serviceType=\"HIVE\""
  HIVE_FD_MAX                                          ="SELECT LAST(fd_max) WHERE serviceType=\"HIVE\""
  HIVE_FD_OPEN                                         ="SELECT LAST(fd_open) WHERE serviceType=\"HIVE\""
  HIVE_HEALTH_BAD_RATE                                 ="SELECT LAST(health_bad_rate) WHERE serviceType=\"HIVE\""
  HIVE_HEALTH_CONCERNING_RATE                          ="SELECT LAST(health_concerning_rate) WHERE serviceType=\"HIVE\""
  HIVE_HEALTH_DISABLED_RATE                            ="SELECT LAST(health_disabled_rate) WHERE serviceType=\"HIVE\""
  HIVE_HEALTH_GOOD_RATE                                ="SELECT LAST(health_good_rate) WHERE serviceType=\"HIVE\""
  HIVE_HEALTH_UNKNOWN_RATE                             ="SELECT LAST(health_unknown_rate) WHERE serviceType=\"HIVE\""
  HIVE_INIT_TOTAL_COUNT_DBS                            ="SELECT LAST(hive_init_total_count_dbs) WHERE serviceType=\"HIVE\""
  HIVE_INIT_TOTAL_COUNT_PARTITIONS                     ="SELECT LAST(hive_init_total_count_partitions) WHERE serviceType=\"HIVE\""
  HIVE_INIT_TOTAL_COUNT_TABLES                         ="SELECT LAST(hive_init_total_count_tables) WHERE serviceType=\"HIVE\""
  HIVE_JVM_PAUSE_TIME_RATE                             ="SELECT LAST(hive_jvm_pause_time_rate) WHERE serviceType=\"HIVE\""
  HIVE_JVM_PAUSES_INFO_THRESHOLD_RATE                  ="SELECT LAST(hive_jvm_pauses_info_threshold_rate) WHERE serviceType=\"HIVE\""
  HIVE_JVM_PAUSES_WARN_THRESHOLD_RATE                  ="SELECT LAST(hive_jvm_pauses_warn_threshold_rate) WHERE serviceType=\"HIVE\""
  HIVE_MEM_RSS                                         ="SELECT LAST(mem_rss) WHERE serviceType=\"HIVE\""
  HIVE_MEM_SWAP                                        ="SELECT LAST(mem_swap) WHERE serviceType=\"HIVE\""
  HIVE_MEM_VIRTUAL                                     ="SELECT LAST(mem_virtual) WHERE serviceType=\"HIVE\""
  HIVE_MEMORY_HEAP_COMMITTED                           ="SELECT LAST(hive_memory_heap_committed) WHERE serviceType=\"HIVE\""
  HIVE_MEMORY_HEAP_INIT                                ="SELECT LAST(hive_memory_heap_init) WHERE serviceType=\"HIVE\""
  HIVE_MEMORY_HEAP_MAX                                 ="SELECT LAST(hive_memory_heap_max) WHERE serviceType=\"HIVE\""
  HIVE_MEMORY_HEAP_USED                                ="SELECT LAST(hive_memory_heap_used) WHERE serviceType=\"HIVE\""
  HIVE_MEMORY_NON_HEAP_COMMITTED                       ="SELECT LAST(hive_memory_non_heap_committed) WHERE serviceType=\"HIVE\""
  HIVE_MEMORY_NON_HEAP_INIT                            ="SELECT LAST(hive_memory_non_heap_init) WHERE serviceType=\"HIVE\""
  HIVE_MEMORY_NON_HEAP_MAX                             ="SELECT LAST(hive_memory_non_heap_max) WHERE serviceType=\"HIVE\""
  HIVE_MEMORY_NON_HEAP_USED                            ="SELECT LAST(hive_memory_non_heap_used) WHERE serviceType=\"HIVE\""
  HIVE_MEMORY_TOTAL_COMMITTED                          ="SELECT LAST(hive_memory_total_committed) WHERE serviceType=\"HIVE\""
  HIVE_MEMORY_TOTAL_INIT                               ="SELECT LAST(hive_memory_total_init) WHERE serviceType=\"HIVE\""
  HIVE_MEMORY_TOTAL_MAX                                ="SELECT LAST(hive_memory_total_max) WHERE serviceType=\"HIVE\""
  HIVE_MEMORY_TOTAL_USED                               ="SELECT LAST(hive_memory_total_used) WHERE serviceType=\"HIVE\""
  HIVE_OOM_EXITS_RATE                                  ="SELECT LAST(oom_exits_rate) WHERE serviceType=\"HIVE\""
  HIVE_OPEN_CONNECTIONS                                ="SELECT LAST(hive_open_connections) WHERE serviceType=\"HIVE\""
  HIVE_OPEN_OPERATIONS                                 ="SELECT LAST(hive_open_operations) WHERE serviceType=\"HIVE\""
  HIVE_READ_BYTES_RATE                                 ="SELECT LAST(read_bytes_rate) WHERE serviceType=\"HIVE\""
  HIVE_THREADS_DAEMON_THREAD_COUNT                     ="SELECT LAST(hive_threads_daemon_thread_count) WHERE serviceType=\"HIVE\""
  HIVE_THREADS_DEADLOCKED_THREAD_COUNT                 ="SELECT LAST(hive_threads_deadlocked_thread_count) WHERE serviceType=\"HIVE\""
  HIVE_THREADS_THREAD_COUNT                            ="SELECT LAST(hive_threads_thread_count) WHERE serviceType=\"HIVE\""
  HIVE_UNEXPECTED_EXITS_RATE                           ="SELECT LAST(unexpected_exits_rate) WHERE serviceType=\"HIVE\""
  HIVE_UPTIME                                          ="SELECT LAST(uptime) WHERE serviceType=\"HIVE\""
  HIVE_WAITING_COMPILE_OPS                             ="SELECT LAST(hive_waiting_compile_ops) WHERE serviceType=\"HIVE\""
  HIVE_WRITE_BYTES_RATE                                ="SELECT LAST(write_bytes_rate) WHERE serviceType=\"HIVE\""
  HIVE_ZOOKEEPER_HIVE_EXCLUSIVELOCKS                   ="SELECT LAST(hive_zookeeper_hive_exclusivelocks) WHERE serviceType=\"HIVE\""
  HIVE_ZOOKEEPER_HIVE_SEMISHAREDLOCKS                  ="SELECT LAST(hive_zookeeper_hive_semisharedlocks) WHERE serviceType=\"HIVE\""
  HIVE_ZOOKEEPER_HIVE_SHAREDLOCKS                      ="SELECT LAST(hive_zookeeper_hive_sharedlocks) WHERE serviceType=\"HIVE\""
)




/* ======================================================================
 * Global variables
 * ====================================================================== */
// Prometheus data Descriptors for the metrics to export
var (
  // Agent Metrics
  hive_alerts_rate                                      =create_hive_metric_struct("hive_alerts_rate", "The number of alerts.	events per second")
  hive_api_acquirereadwritelocks_avg                    =create_hive_metric_struct("hive_api_acquirereadwritelocks_avg", "acquireReadWriteLocks method calls: Avg	ms")
  hive_api_acquirereadwritelocks_max                    =create_hive_metric_struct("hive_api_acquirereadwritelocks_max", "acquireReadWriteLocks method calls: Max	ms")
  hive_api_acquirereadwritelocks_median                 =create_hive_metric_struct("hive_api_acquirereadwritelocks_median", "acquireReadWriteLocks method calls: 50th Percentile	ms")
  hive_api_acquirereadwritelocks_min                    =create_hive_metric_struct("hive_api_acquirereadwritelocks_min", "acquireReadWriteLocks method calls: Min	ms")
  hive_api_acquirereadwritelocks_rate                   =create_hive_metric_struct("hive_api_acquirereadwritelocks_rate", "acquireReadWriteLocks method calls: Samples	message.units.executions per second")
  hive_api_acquirereadwritelocks_stddev                 =create_hive_metric_struct("hive_api_acquirereadwritelocks_stddev", "acquireReadWriteLocks method calls: Standard Deviation	ms")
  hive_api_compile_avg                                  =create_hive_metric_struct("hive_api_compile_avg", "compile method calls: Avg	ms")
  hive_api_compile_max                                  =create_hive_metric_struct("hive_api_compile_max", "compile method calls: Max	ms")
  hive_api_compile_median                               =create_hive_metric_struct("hive_api_compile_median", "compile method calls: 50th Percentile	ms")
  hive_api_compile_min                                  =create_hive_metric_struct("hive_api_compile_min", "compile method calls: Min	ms")
  hive_api_compile_rate                                 =create_hive_metric_struct("hive_api_compile_rate", "compile method calls: Samples	message.units.executions per second")
  hive_api_compile_stddev                               =create_hive_metric_struct("hive_api_compile_stddev", "compile method calls: Standard Deviation	ms")
  hive_api_driver_execute_avg                           =create_hive_metric_struct("hive_api_driver_execute_avg", "Driver.execute method calls: Avg	ms")
  hive_api_driver_execute_max                           =create_hive_metric_struct("hive_api_driver_execute_max", "Driver.execute method calls: Max	ms")
  hive_api_driver_execute_median                        =create_hive_metric_struct("hive_api_driver_execute_median", "Driver.execute method calls: 50th Percentile	ms")
  hive_api_driver_execute_min                           =create_hive_metric_struct("hive_api_driver_execute_min", "Driver.execute method calls: Min	ms")
  hive_api_driver_execute_rate                          =create_hive_metric_struct("hive_api_driver_execute_rate", "Driver.execute method calls: Samples	message.units.executions per second")
  hive_api_driver_execute_stddev                        =create_hive_metric_struct("hive_api_driver_execute_stddev", "Driver.execute method calls: Standard Deviation	ms")
  hive_api_driver_run_avg                               =create_hive_metric_struct("hive_api_driver_run_avg", "Driver.run method calls: Avg	ms")
  hive_api_driver_run_max                               =create_hive_metric_struct("hive_api_driver_run_max", "Driver.run method calls: Max	ms")
  hive_api_driver_run_median                            =create_hive_metric_struct("hive_api_driver_run_median", "Driver.run method calls: 50th Percentile	ms")
  hive_api_driver_run_min                               =create_hive_metric_struct("hive_api_driver_run_min", "Driver.run method calls: Min	ms")
  hive_api_driver_run_rate                              =create_hive_metric_struct("hive_api_driver_run_rate", "Driver.run method calls: Samples	message.units.executions per second")
  hive_api_driver_run_stddev                            =create_hive_metric_struct("hive_api_driver_run_stddev", "Driver.run method calls: Standard Deviation	ms")
  hive_api_get_all_tables_avg                           =create_hive_metric_struct("hive_api_get_all_tables_avg", "get_all_tables method calls: Avg	ms")
  hive_api_get_all_tables_max                           =create_hive_metric_struct("hive_api_get_all_tables_max", "get_all_tables method calls: Max	ms")
  hive_api_get_all_tables_median                        =create_hive_metric_struct("hive_api_get_all_tables_median", "get_all_tables method calls: 50th Percentile	ms")
  hive_api_get_all_tables_min                           =create_hive_metric_struct("hive_api_get_all_tables_min", "get_all_tables method calls: Min	ms")
  hive_api_get_all_tables_rate                          =create_hive_metric_struct("hive_api_get_all_tables_rate", "get_all_tables method calls: Samples	message.units.executions per second")
  hive_api_get_all_tables_stddev                        =create_hive_metric_struct("hive_api_get_all_tables_stddev", "get_all_tables method calls: Standard Deviation	ms")
  hive_api_get_config_value_avg                         =create_hive_metric_struct("hive_api_get_config_value_avg", "get_config_value method calls: Avg	ms")
  hive_api_get_config_value_max                         =create_hive_metric_struct("hive_api_get_config_value_max", "get_config_value method calls: Max	ms")
  hive_api_get_config_value_median                      =create_hive_metric_struct("hive_api_get_config_value_median", "get_config_value method calls: 50th Percentile	ms")
  hive_api_get_config_value_min                         =create_hive_metric_struct("hive_api_get_config_value_min", "get_config_value method calls: Min	ms")
  hive_api_get_config_value_rate                        =create_hive_metric_struct("hive_api_get_config_value_rate", "get_config_value method calls: Samples	message.units.executions per second")
  hive_api_get_config_value_stddev                      =create_hive_metric_struct("hive_api_get_config_value_stddev", "get_config_value method calls: Standard Deviation	ms")
  hive_api_get_database_avg                             =create_hive_metric_struct("hive_api_get_database_avg", "get_database method calls: Avg	ms")
  hive_api_get_database_max                             =create_hive_metric_struct("hive_api_get_database_max", "get_database method calls: Max	ms")
  hive_api_get_database_median                          =create_hive_metric_struct("hive_api_get_database_median", "get_database method calls: 50th Percentile	ms")
  hive_api_get_database_min                             =create_hive_metric_struct("hive_api_get_database_min", "get_database method calls: Min	ms")
  hive_api_get_database_rate                            =create_hive_metric_struct("hive_api_get_database_rate", "get_database method calls: Samples	message.units.executions per second")
  hive_api_get_database_stddev                          =create_hive_metric_struct("hive_api_get_database_stddev", "get_database method calls: Standard Deviation	ms")
  hive_api_get_databases_avg                            =create_hive_metric_struct("hive_api_get_databases_avg", "get_databases method calls: Avg	ms")
  hive_api_get_databases_max                            =create_hive_metric_struct("hive_api_get_databases_max", "get_databases method calls: Max	ms")
  hive_api_get_databases_median                         =create_hive_metric_struct("hive_api_get_databases_median", "get_databases method calls: 50th Percentile	ms")
  hive_api_get_databases_min                            =create_hive_metric_struct("hive_api_get_databases_min", "get_databases method calls: Min	ms")
  hive_api_get_databases_rate                           =create_hive_metric_struct("hive_api_get_databases_rate", "get_databases method calls: Samples	message.units.executions per second")
  hive_api_get_databases_stddev                         =create_hive_metric_struct("hive_api_get_databases_stddev", "get_databases method calls: Standard Deviation	ms")
  hive_api_get_fields_with_environment_context_avg      =create_hive_metric_struct("hive_api_get_fields_with_environment_context_avg", "get_fields_with_environment_context method calls: Avg	ms")
  hive_api_get_fields_with_environment_context_max      =create_hive_metric_struct("hive_api_get_fields_with_environment_context_max", "get_fields_with_environment_context method calls: Max	ms")
  hive_api_get_fields_with_environment_context_median   =create_hive_metric_struct("hive_api_get_fields_with_environment_context_median", "get_fields_with_environment_context method calls: 50th Percentile	ms")
  hive_api_get_fields_with_environment_context_min      =create_hive_metric_struct("hive_api_get_fields_with_environment_context_min", "get_fields_with_environment_context method calls: Min	ms")
  hive_api_get_fields_with_environment_context_rate     =create_hive_metric_struct("hive_api_get_fields_with_environment_context_rate", "get_fields_with_environment_context method calls: Samples	message.units.executions per second")
  hive_api_get_fields_with_environment_context_stddev   =create_hive_metric_struct("hive_api_get_fields_with_environment_context_stddev", "get_fields_with_environment_context method calls: Standard Deviation	ms")
  hive_api_get_function_avg                             =create_hive_metric_struct("hive_api_get_function_avg", "get_function method calls: Avg	ms")
  hive_api_get_function_max                             =create_hive_metric_struct("hive_api_get_function_max", "get_function method calls: Max	ms")
  hive_api_get_function_median                          =create_hive_metric_struct("hive_api_get_function_median", "get_function method calls: 50th Percentile	ms")
  hive_api_get_function_min                             =create_hive_metric_struct("hive_api_get_function_min", "get_function method calls: Min	ms")
  hive_api_get_function_rate                            =create_hive_metric_struct("hive_api_get_function_rate", "get_function method calls: Samples	message.units.executions per second")
  hive_api_get_function_stddev                          =create_hive_metric_struct("hive_api_get_function_stddev", "get_function method calls: Standard Deviation	ms")
  hive_api_get_functions_avg                            =create_hive_metric_struct("hive_api_get_functions_avg", "get_functions method calls: Avg	ms")
  hive_api_get_functions_max                            =create_hive_metric_struct("hive_api_get_functions_max", "get_functions method calls: Max	ms")
  hive_api_get_functions_median                         =create_hive_metric_struct("hive_api_get_functions_median", "get_functions method calls: 50th Percentile	ms")
  hive_api_get_functions_min                            =create_hive_metric_struct("hive_api_get_functions_min", "get_functions method calls: Min	ms")
  hive_api_get_functions_rate                           =create_hive_metric_struct("hive_api_get_functions_rate", "get_functions method calls: Samples	message.units.executions per second")
  hive_api_get_functions_stddev                         =create_hive_metric_struct("hive_api_get_functions_stddev", "get_functions method calls: Standard Deviation	ms")
  hive_api_get_index_by_name_avg                        =create_hive_metric_struct("hive_api_get_index_by_name_avg", "get_index_by_name method calls: Avg	ms")
  hive_api_get_index_by_name_max                        =create_hive_metric_struct("hive_api_get_index_by_name_max", "get_index_by_name method calls: Max	ms")
  hive_api_get_index_by_name_median                     =create_hive_metric_struct("hive_api_get_index_by_name_median", "get_index_by_name method calls: 50th Percentile	ms")
  hive_api_get_index_by_name_min                        =create_hive_metric_struct("hive_api_get_index_by_name_min", "get_index_by_name method calls: Min	ms")
  hive_api_get_index_by_name_rate                       =create_hive_metric_struct("hive_api_get_index_by_name_rate", "get_index_by_name method calls: Samples	message.units.executions per second")
  hive_api_get_index_by_name_stddev                     =create_hive_metric_struct("hive_api_get_index_by_name_stddev", "get_index_by_name method calls: Standard Deviation	ms")
  hive_api_get_index_names_avg                          =create_hive_metric_struct("hive_api_get_index_names_avg", "get_index_names method calls: Avg	ms")
  hive_api_get_index_names_max                          =create_hive_metric_struct("hive_api_get_index_names_max", "get_index_names method calls: Max	ms")
  hive_api_get_index_names_median                       =create_hive_metric_struct("hive_api_get_index_names_median", "get_index_names method calls: 50th Percentile	ms")
  hive_api_get_index_names_min                          =create_hive_metric_struct("hive_api_get_index_names_min", "get_index_names method calls: Min	ms")
  hive_api_get_index_names_rate                         =create_hive_metric_struct("hive_api_get_index_names_rate", "get_index_names method calls: Samples	message.units.executions per second")
  hive_api_get_index_names_stddev                       =create_hive_metric_struct("hive_api_get_index_names_stddev", "get_index_names method calls: Standard Deviation	ms")
  hive_api_get_indexes_avg                              =create_hive_metric_struct("hive_api_get_indexes_avg", "get_indexes method calls: Avg	ms")
  hive_api_get_indexes_max                              =create_hive_metric_struct("hive_api_get_indexes_max", "get_indexes method calls: Max	ms")
  hive_api_get_indexes_median                           =create_hive_metric_struct("hive_api_get_indexes_median", "get_indexes method calls: 50th Percentile	ms")
  hive_api_get_indexes_min                              =create_hive_metric_struct("hive_api_get_indexes_min", "get_indexes method calls: Min	ms")
  hive_api_get_indexes_rate                             =create_hive_metric_struct("hive_api_get_indexes_rate", "get_indexes method calls: Samples	message.units.executions per second")
  hive_api_get_indexes_stddev                           =create_hive_metric_struct("hive_api_get_indexes_stddev", "get_indexes method calls: Standard Deviation	ms")
  hive_api_get_multi_table_avg                          =create_hive_metric_struct("hive_api_get_multi_table_avg", "get_multi_table method calls: Avg	ms")
  hive_api_get_multi_table_max                          =create_hive_metric_struct("hive_api_get_multi_table_max", "get_multi_table method calls: Max	ms")
  hive_api_get_multi_table_median                       =create_hive_metric_struct("hive_api_get_multi_table_median", "get_multi_table method calls: 50th Percentile	ms")
  hive_api_get_multi_table_min                          =create_hive_metric_struct("hive_api_get_multi_table_min", "get_multi_table method calls: Min	ms")
  hive_api_get_multi_table_rate                         =create_hive_metric_struct("hive_api_get_multi_table_rate", "get_multi_table method calls: Samples	message.units.executions per second")
  hive_api_get_multi_table_stddev                       =create_hive_metric_struct("hive_api_get_multi_table_stddev", "get_multi_table method calls: Standard Deviation	ms")
  hive_api_get_table_avg                                =create_hive_metric_struct("hive_api_get_table_avg", "get_table method calls: Avg	ms")
  hive_api_get_table_max                                =create_hive_metric_struct("hive_api_get_table_max", "get_table method calls: Max	ms")
  hive_api_get_table_median                             =create_hive_metric_struct("hive_api_get_table_median", "get_table method calls: 50th Percentile	ms")
  hive_api_get_table_min                                =create_hive_metric_struct("hive_api_get_table_min", "get_table method calls: Min	ms")
  hive_api_get_table_names_by_filter_avg                =create_hive_metric_struct("hive_api_get_table_names_by_filter_avg", "get_table_names_by_filter method calls: Avg	ms")
  hive_api_get_table_names_by_filter_max                =create_hive_metric_struct("hive_api_get_table_names_by_filter_max", "get_table_names_by_filter method calls: Max	ms")
  hive_api_get_table_names_by_filter_median             =create_hive_metric_struct("hive_api_get_table_names_by_filter_median", "get_table_names_by_filter method calls: 50th Percentile	ms")
  hive_api_get_table_names_by_filter_min                =create_hive_metric_struct("hive_api_get_table_names_by_filter_min", "get_table_names_by_filter method calls: Min	ms")
  hive_api_get_table_names_by_filter_rate               =create_hive_metric_struct("hive_api_get_table_names_by_filter_rate", "get_table_names_by_filter method calls: Samples	message.units.executions per second")
  hive_api_get_table_names_by_filter_stddev             =create_hive_metric_struct("hive_api_get_table_names_by_filter_stddev", "get_table_names_by_filter method calls: Standard Deviation	ms")
  hive_api_get_table_rate                               =create_hive_metric_struct("hive_api_get_table_rate", "get_table method calls: Samples	message.units.executions per second")
  hive_api_get_table_stddev                             =create_hive_metric_struct("hive_api_get_table_stddev", "get_table method calls: Standard Deviation	ms")
  hive_api_get_tables_avg                               =create_hive_metric_struct("hive_api_get_tables_avg", "get_tables method calls: Avg	ms")
  hive_api_get_tables_max                               =create_hive_metric_struct("hive_api_get_tables_max", "get_tables method calls: Max	ms")
  hive_api_get_tables_median                            =create_hive_metric_struct("hive_api_get_tables_median", "get_tables method calls: 50th Percentile	ms")
  hive_api_get_tables_min                               =create_hive_metric_struct("hive_api_get_tables_min", "get_tables method calls: Min	ms")
  hive_api_get_tables_rate                              =create_hive_metric_struct("hive_api_get_tables_rate", "get_tables method calls: Samples	message.units.executions per second")
  hive_api_get_tables_stddev                            =create_hive_metric_struct("hive_api_get_tables_stddev", "get_tables method calls: Standard Deviation	ms")
  hive_api_getinputsummary_avg                          =create_hive_metric_struct("hive_api_getinputsummary_avg", "getInputSummary method calls: Avg	ms")
  hive_api_getinputsummary_max                          =create_hive_metric_struct("hive_api_getinputsummary_max", "getInputSummary method calls: Max	ms")
  hive_api_getinputsummary_median                       =create_hive_metric_struct("hive_api_getinputsummary_median", "getInputSummary method calls: 50th Percentile	ms")
  hive_api_getinputsummary_min                          =create_hive_metric_struct("hive_api_getinputsummary_min", "getInputSummary method calls: Min	ms")
  hive_api_getinputsummary_rate                         =create_hive_metric_struct("hive_api_getinputsummary_rate", "getInputSummary method calls: Samples	message.units.executions per second")
  hive_api_getinputsummary_stddev                       =create_hive_metric_struct("hive_api_getinputsummary_stddev", "getInputSummary method calls: Standard Deviation	ms")
  hive_api_getsplits_avg                                =create_hive_metric_struct("hive_api_getsplits_avg", "getSplits method calls: Avg	ms")
  hive_api_getsplits_max                                =create_hive_metric_struct("hive_api_getsplits_max", "getSplits method calls: Max	ms")
  hive_api_getsplits_median                             =create_hive_metric_struct("hive_api_getsplits_median", "getSplits method calls: 50th Percentile	ms")
  hive_api_getsplits_min                                =create_hive_metric_struct("hive_api_getsplits_min", "getSplits method calls: Min	ms")
  hive_api_getsplits_rate                               =create_hive_metric_struct("hive_api_getsplits_rate", "getSplits method calls: Samples	message.units.executions per second")
  hive_api_getsplits_stddev                             =create_hive_metric_struct("hive_api_getsplits_stddev", "getSplits method calls: Standard Deviation	ms")
  hive_api_markpartitionforevent_avg                    =create_hive_metric_struct("hive_api_markpartitionforevent_avg", "markPartitionForEvent method calls: Avg	ms")
  hive_api_markpartitionforevent_max                    =create_hive_metric_struct("hive_api_markpartitionforevent_max", "markPartitionForEvent method calls: Max	ms")
  hive_api_markpartitionforevent_median                 =create_hive_metric_struct("hive_api_markpartitionforevent_median", "markPartitionForEvent method calls: 50th Percentile	ms")
  hive_api_markpartitionforevent_min                    =create_hive_metric_struct("hive_api_markpartitionforevent_min", "markPartitionForEvent method calls: Min	ms")
  hive_api_markpartitionforevent_rate                   =create_hive_metric_struct("hive_api_markpartitionforevent_rate", "markPartitionForEvent method calls: Samples	message.units.executions per second")
  hive_api_markpartitionforevent_stddev                 =create_hive_metric_struct("hive_api_markpartitionforevent_stddev", "markPartitionForEvent method calls: Standard Deviation	ms")
  hive_api_operation_initialized_avg                    =create_hive_metric_struct("hive_api_operation_initialized_avg", "HiveServer2 operations in INITIALIZED state: Avg	ms")
  hive_api_operation_initialized_max                    =create_hive_metric_struct("hive_api_operation_initialized_max", "HiveServer2 operations in INITIALIZED state: Max	ms")
  hive_api_operation_initialized_median                 =create_hive_metric_struct("hive_api_operation_initialized_median", "HiveServer2 operations in INITIALIZED state: 50th Percentile	ms")
  hive_api_operation_initialized_min                    =create_hive_metric_struct("hive_api_operation_initialized_min", "HiveServer2 operations in INITIALIZED state: Min	ms")
  hive_api_operation_initialized_rate                   =create_hive_metric_struct("hive_api_operation_initialized_rate", "HiveServer2 operations in INITIALIZED state: Samples	operations per second")
  hive_api_operation_initialized_stddev                 =create_hive_metric_struct("hive_api_operation_initialized_stddev", "HiveServer2 operations in INITIALIZED state: Standard Deviation	ms")
  hive_api_operation_pending_avg                        =create_hive_metric_struct("hive_api_operation_pending_avg", "HiveServer2 operations in PENDING state: Avg	ms")
  hive_api_operation_pending_max                        =create_hive_metric_struct("hive_api_operation_pending_max", "HiveServer2 operations in PENDING state: Max	ms")
  hive_api_operation_pending_median                     =create_hive_metric_struct("hive_api_operation_pending_median", "HiveServer2 operations in PENDING state: 50th Percentile	ms")
  hive_api_operation_pending_min                        =create_hive_metric_struct("hive_api_operation_pending_min", "HiveServer2 operations in PENDING state: Min	ms")
  hive_api_operation_pending_rate                       =create_hive_metric_struct("hive_api_operation_pending_rate", "HiveServer2 operations in PENDING state: Samples	operations per second")
  hive_api_operation_pending_stddev                     =create_hive_metric_struct("hive_api_operation_pending_stddev", "HiveServer2 operations in PENDING state: Standard Deviation	ms")
  hive_api_operation_running_avg                        =create_hive_metric_struct("hive_api_operation_running_avg", "HiveServer2 operations in RUNNING state: Avg	ms")
  hive_api_operation_running_max                        =create_hive_metric_struct("hive_api_operation_running_max", "HiveServer2 operations in RUNNING state: Max	ms")
  hive_api_operation_running_median                     =create_hive_metric_struct("hive_api_operation_running_median", "HiveServer2 operations in RUNNING state: 50th Percentile	ms")
  hive_api_operation_running_min                        =create_hive_metric_struct("hive_api_operation_running_min", "HiveServer2 operations in RUNNING state: Min	ms")
  hive_api_operation_running_rate                       =create_hive_metric_struct("hive_api_operation_running_rate", "HiveServer2 operations in RUNNING state: Samples	operations per second")
  hive_api_operation_running_stddev                     =create_hive_metric_struct("hive_api_operation_running_stddev", "HiveServer2 operations in RUNNING state: Standard Deviation	ms")
  hive_api_parse_avg                                    =create_hive_metric_struct("hive_api_parse_avg", "parse method calls: Avg	ms")
  hive_api_parse_max                                    =create_hive_metric_struct("hive_api_parse_max", "parse method calls: Max	ms")
  hive_api_parse_median                                 =create_hive_metric_struct("hive_api_parse_median", "parse method calls: 50th Percentile	ms")
  hive_api_parse_min                                    =create_hive_metric_struct("hive_api_parse_min", "parse method calls: Min	ms")
  hive_api_parse_rate                                   =create_hive_metric_struct("hive_api_parse_rate", "parse method calls: Samples	message.units.executions per second")
  hive_api_parse_stddev                                 =create_hive_metric_struct("hive_api_parse_stddev", "parse method calls: Standard Deviation	ms")
  hive_api_partition_retrieving_avg                     =create_hive_metric_struct("hive_api_partition_retrieving_avg", "partition-retrieving method calls: Avg	ms")
  hive_api_partition_retrieving_max                     =create_hive_metric_struct("hive_api_partition_retrieving_max", "partition-retrieving method calls: Max	ms")
  hive_api_partition_retrieving_median                  =create_hive_metric_struct("hive_api_partition_retrieving_median", "partition-retrieving method calls: 50th Percentile	ms")
  hive_api_partition_retrieving_min                     =create_hive_metric_struct("hive_api_partition_retrieving_min", "partition-retrieving method calls: Min	ms")
  hive_api_partition_retrieving_rate                    =create_hive_metric_struct("hive_api_partition_retrieving_rate", "partition-retrieving method calls: Samples	message.units.executions per second")
  hive_api_partition_retrieving_stddev                  =create_hive_metric_struct("hive_api_partition_retrieving_stddev", "partition-retrieving method calls: Standard Deviation	ms")
  hive_api_releaselocks_avg                             =create_hive_metric_struct("hive_api_releaselocks_avg", "releaseLocks method calls: Avg	ms")
  hive_api_releaselocks_max                             =create_hive_metric_struct("hive_api_releaselocks_max", "releaseLocks method calls: Max	ms")
  hive_api_releaselocks_median                          =create_hive_metric_struct("hive_api_releaselocks_median", "releaseLocks method calls: 50th Percentile	ms")
  hive_api_releaselocks_min                             =create_hive_metric_struct("hive_api_releaselocks_min", "releaseLocks method calls: Min	ms")
  hive_api_releaselocks_rate                            =create_hive_metric_struct("hive_api_releaselocks_rate", "releaseLocks method calls: Samples	message.units.executions per second")
  hive_api_releaselocks_stddev                          =create_hive_metric_struct("hive_api_releaselocks_stddev", "releaseLocks method calls: Standard Deviation	ms")
  hive_api_runtasks_avg                                 =create_hive_metric_struct("hive_api_runtasks_avg", "runTasks method calls: Avg	ms")
  hive_api_runtasks_max                                 =create_hive_metric_struct("hive_api_runtasks_max", "runTasks method calls: Max	ms")
  hive_api_runtasks_median                              =create_hive_metric_struct("hive_api_runtasks_median", "runTasks method calls: 50th Percentile	ms")
  hive_api_runtasks_min                                 =create_hive_metric_struct("hive_api_runtasks_min", "runTasks method calls: Min	ms")
  hive_api_runtasks_rate                                =create_hive_metric_struct("hive_api_runtasks_rate", "runTasks method calls: Samples	message.units.executions per second")
  hive_api_runtasks_stddev                              =create_hive_metric_struct("hive_api_runtasks_stddev", "runTasks method calls: Standard Deviation	ms")
  hive_api_serializeplan_avg                            =create_hive_metric_struct("hive_api_serializeplan_avg", "serializePlan method calls: Avg	ms")
  hive_api_serializeplan_max                            =create_hive_metric_struct("hive_api_serializeplan_max", "serializePlan method calls: Max	ms")
  hive_api_serializeplan_median                         =create_hive_metric_struct("hive_api_serializeplan_median", "serializePlan method calls: 50th Percentile	ms")
  hive_api_serializeplan_min                            =create_hive_metric_struct("hive_api_serializeplan_min", "serializePlan method calls: Min	ms")
  hive_api_serializeplan_rate                           =create_hive_metric_struct("hive_api_serializeplan_rate", "serializePlan method calls: Samples	message.units.executions per second")
  hive_api_serializeplan_stddev                         =create_hive_metric_struct("hive_api_serializeplan_stddev", "serializePlan method calls: Standard Deviation	ms")
  hive_api_timetosubmit_avg                             =create_hive_metric_struct("hive_api_timetosubmit_avg", "TimeToSubmit method calls: Avg	ms")
  hive_api_timetosubmit_max                             =create_hive_metric_struct("hive_api_timetosubmit_max", "TimeToSubmit method calls: Max	ms")
  hive_api_timetosubmit_median                          =create_hive_metric_struct("hive_api_timetosubmit_median", "TimeToSubmit method calls: 50th Percentile	ms")
  hive_api_timetosubmit_min                             =create_hive_metric_struct("hive_api_timetosubmit_min", "TimeToSubmit method calls: Min	ms")
  hive_api_timetosubmit_rate                            =create_hive_metric_struct("hive_api_timetosubmit_rate", "TimeToSubmit method calls: Samples	message.units.executions per second")
  hive_api_timetosubmit_stddev                          =create_hive_metric_struct("hive_api_timetosubmit_stddev", "TimeToSubmit method calls: Standard Deviation	ms")
  hive_audit_agent_bytes_left_rate                      =create_hive_metric_struct("hive_audit_agent_bytes_left_rate", "Number of bytes of audit events still to be sent to the Audit Server.	bytes per second")
  hive_audit_agent_bytes_sent_rate                      =create_hive_metric_struct("hive_audit_agent_bytes_sent_rate", "Number of bytes of audit events successfully sent by Cloudera Manager Agent to the Audit Server.	bytes per second")
  hive_audit_agent_events_sent_rate                     =create_hive_metric_struct("hive_audit_agent_events_sent_rate", "Number of audit events successfully sent by Cloudera Manager Agent to the Audit Server.	events per second")
  hive_audit_agent_send_failures_rate                   =create_hive_metric_struct("hive_audit_agent_send_failures_rate", "Number of time Cloudera Manager Agent encountered failure when sending audit event(s) to the Audit Server.	errors per second")
  hive_audit_plugin_coalesced_rate                      =create_hive_metric_struct("hive_audit_plugin_coalesced_rate", "Number of audit events coalesced	events per second")
  hive_audit_plugin_events_rate                         =create_hive_metric_struct("hive_audit_plugin_events_rate", "Number of audit events persisted to local filesystem	events per second")
  hive_audit_plugin_exceptions_rate                     =create_hive_metric_struct("hive_audit_plugin_exceptions_rate", "Number of exception encountered when processing audit events	errors per second")
  hive_audit_plugin_filtered_rate                       =create_hive_metric_struct("hive_audit_plugin_filtered_rate", "Number of audit events that were filtered based on user specified criteria	events per second")
  hive_audit_plugin_parse_errors_rate                   =create_hive_metric_struct("hive_audit_plugin_parse_errors_rate", "Number of audit events that cannot be parsed	errors per second")
  hive_canary_duration                                  =create_hive_metric_struct("hive_canary_duration", "Duration of the last or currently running canary job	ms")
  hive_cgroup_cpu_system_rate                           =create_hive_metric_struct("hive_cgroup_cpu_system_rate", "CPU usage of the role's cgroup	seconds per second")
  hive_cgroup_cpu_user_rate                             =create_hive_metric_struct("hive_cgroup_cpu_user_rate", "User Space CPU usage of the role's cgroup	seconds per second")
  hive_cgroup_mem_page_cache                            =create_hive_metric_struct("hive_cgroup_mem_page_cache", "Page cache usage of the role's cgroup	bytes")
  hive_cgroup_mem_rss                                   =create_hive_metric_struct("hive_cgroup_mem_rss", "Resident memory of the role's cgroup	bytes")
  hive_cgroup_mem_swap                                  =create_hive_metric_struct("hive_cgroup_mem_swap", "Swap usage of the role's cgroup	bytes")
  hive_cgroup_read_bytes_rate                           =create_hive_metric_struct("hive_cgroup_read_bytes_rate", "Bytes read from all disks by the role's cgroup	bytes per second")
  hive_cgroup_read_ios_rate                             =create_hive_metric_struct("hive_cgroup_read_ios_rate", "Number of read I/O operations from all disks by the role's cgroup	ios per second")
  hive_cgroup_write_bytes_rate                          =create_hive_metric_struct("hive_cgroup_write_bytes_rate", "Bytes written to all disks by the role's cgroup	bytes per second")
  hive_cgroup_write_ios_rate                            =create_hive_metric_struct("hive_cgroup_write_ios_rate", "Number of write I/O operations to all disks by the role's cgroup	ios per second")
  hive_completed_operation_canceled_rate                =create_hive_metric_struct("hive_completed_operation_canceled_rate", "Number of cancelled HiveServer2 operations	operations per second")
  hive_completed_operation_closed_rate                  =create_hive_metric_struct("hive_completed_operation_closed_rate", "Number of HiveServer2 operations that have been closed	operations per second")
  hive_completed_operation_error_rate                   =create_hive_metric_struct("hive_completed_operation_error_rate", "Number of HiveServer2 operations finished in error	operations per second")
  hive_completed_operation_finished_rate                =create_hive_metric_struct("hive_completed_operation_finished_rate", "Number of successfully finished HiveServer2 operations	operations per second")
  hive_cpu_system_rate                                  =create_hive_metric_struct("hive_cpu_system_rate", "Total System CPU	seconds per second")
  hive_cpu_user_rate                                    =create_hive_metric_struct("hive_cpu_user_rate", "Total CPU user time	seconds per second")
  hive_create_total_count_dbs_rate                      =create_hive_metric_struct("hive_create_total_count_dbs_rate", "Number of databases created by this Hive Metastore Server since startup	databases per second")
  hive_create_total_count_partitions_rate               =create_hive_metric_struct("hive_create_total_count_partitions_rate", "Number of partitions created by this Hive Metastore Server since startup	partitions per second")
  hive_create_total_count_tables_rate                   =create_hive_metric_struct("hive_create_total_count_tables_rate", "Number of tables created by this Hive Metastore Server since startup	table per second")
  hive_delete_total_count_dbs_rate                      =create_hive_metric_struct("hive_delete_total_count_dbs_rate", "Number of databases dropped by this Hive Metastore Server since startup	databases per second")
  hive_delete_total_count_partitions_rate               =create_hive_metric_struct("hive_delete_total_count_partitions_rate", "Number of partitions dropped by this Hive Metastore Server since startup	partitions per second")
  hive_delete_total_count_tables_rate                   =create_hive_metric_struct("hive_delete_total_count_tables_rate", "Number of tables dropped by this Hive Metastore Server since startup	table per second")
  hive_events_critical_rate                             =create_hive_metric_struct("hive_events_critical_rate", "The number of critical events.	events per second")
  hive_events_important_rate                            =create_hive_metric_struct("hive_events_important_rate", "The number of important events.	events per second")
  hive_events_informational_rate                        =create_hive_metric_struct("hive_events_informational_rate", "The number of informational events.	events per second")
  hive_exec_async_pool_size                             =create_hive_metric_struct("hive_exec_async_pool_size", "Current size of HiveServer2 asynchronous thread pool	threads")
  hive_exec_async_queue_size                            =create_hive_metric_struct("hive_exec_async_queue_size", "Current size of HiveServer2 asynchronous operation queue	message.units.runnables")
  hive_fd_max                                           =create_hive_metric_struct("hive_fd_max", "Maximum number of file descriptors	file descriptors")
  hive_fd_open                                          =create_hive_metric_struct("hive_fd_open", "Open file descriptors.	file descriptors")
  hive_health_bad_rate                                  =create_hive_metric_struct("hive_health_bad_rate", "Percentage of Time with Bad Health	seconds per second")
  hive_health_concerning_rate                           =create_hive_metric_struct("hive_health_concerning_rate", "Percentage of Time with Concerning Health	seconds per second")
  hive_health_disabled_rate                             =create_hive_metric_struct("hive_health_disabled_rate", "Percentage of Time with Disabled Health	seconds per second")
  hive_health_good_rate                                 =create_hive_metric_struct("hive_health_good_rate", "Percentage of Time with Good Health	seconds per second")
  hive_health_unknown_rate                              =create_hive_metric_struct("hive_health_unknown_rate", "Percentage of Time with Unknown Health	seconds per second")
  hive_init_total_count_dbs                             =create_hive_metric_struct("hive_init_total_count_dbs", "Number of databases at startup of this Hive Metastore Server	databases")
  hive_init_total_count_partitions                      =create_hive_metric_struct("hive_init_total_count_partitions", "Number of partitions at startup of this Hive Metastore Server	partitions")
  hive_init_total_count_tables                          =create_hive_metric_struct("hive_init_total_count_tables", "Number of tables at startup of this Hive Metastore Server	table")
  hive_jvm_pause_time_rate                              =create_hive_metric_struct("hive_jvm_pause_time_rate", "The amount of extra time the jvm was paused above the requested sleep time. The JVM pause monitor sleeps for 500 milliseconds and any extra time it waited above this is counted in the pause time.	ms per second")
  hive_jvm_pauses_info_threshold_rate                   =create_hive_metric_struct("hive_jvm_pauses_info_threshold_rate", "Number of JVM pauses longer than the info threshold but shorter than the warning threshold. By default the info threshold is set to 1 second. To change use this configuration key JvmPauseMonitorService.info-threshold.ms	pauses per second")
  hive_jvm_pauses_warn_threshold_rate                   =create_hive_metric_struct("hive_jvm_pauses_warn_threshold_rate", "Number of JVM pauses longer than the warning threshold. By default the warning threshold is set to 10 second. To change use this configuration key JvmPauseMonitorService.warn-threshold.ms	pauses per second")
  hive_mem_rss                                          =create_hive_metric_struct("hive_mem_rss", "Resident memory used	bytes")
  hive_mem_swap                                         =create_hive_metric_struct("hive_mem_swap", "Amount of swap memory used by this role's process.	bytes")
  hive_mem_virtual                                      =create_hive_metric_struct("hive_mem_virtual", "Virtual memory used	bytes")
  hive_memory_heap_committed                            =create_hive_metric_struct("hive_memory_heap_committed", "JVM heap committed memory	bytes")
  hive_memory_heap_init                                 =create_hive_metric_struct("hive_memory_heap_init", "JVM heap initial memory	bytes")
  hive_memory_heap_max                                  =create_hive_metric_struct("hive_memory_heap_max", "JVM heap max used memory	bytes")
  hive_memory_heap_used                                 =create_hive_metric_struct("hive_memory_heap_used", "JVM heap used memory	bytes")
  hive_memory_non_heap_committed                        =create_hive_metric_struct("hive_memory_non_heap_committed", "JVM non heap committed memory	bytes")
  hive_memory_non_heap_init                             =create_hive_metric_struct("hive_memory_non_heap_init", "JVM non heap initial memory	bytes")
  hive_memory_non_heap_max                              =create_hive_metric_struct("hive_memory_non_heap_max", "JVM non heap max used memory	bytes")
  hive_memory_non_heap_used                             =create_hive_metric_struct("hive_memory_non_heap_used", "JVM non heap used memory	bytes")
  hive_memory_total_committed                           =create_hive_metric_struct("hive_memory_total_committed", "JVM heap and non-heap committed memory	bytes")
  hive_memory_total_init                                =create_hive_metric_struct("hive_memory_total_init", "JVM heap and non-heap initial memory	bytes")
  hive_memory_total_max                                 =create_hive_metric_struct("hive_memory_total_max", "JVM heap and non-heap max initial memory	bytes")
  hive_memory_total_used                                =create_hive_metric_struct("hive_memory_total_used", "JVM heap and non-heap used memory	bytes")
  hive_oom_exits_rate                                   =create_hive_metric_struct("hive_oom_exits_rate", "The number of times the role's backing process was killed due to an OutOfMemory error. This counter is only incremented if the Cloudera Manager 'Kill When Out of Memory' option is enabled.	exits per second")
  hive_open_connections                                 =create_hive_metric_struct("hive_open_connections", "Number of open connections to the server	connections")
  hive_open_operations                                  =create_hive_metric_struct("hive_open_operations", "Number of open operations to the server	operations")
  hive_read_bytes_rate                                  =create_hive_metric_struct("hive_read_bytes_rate", "The number of bytes read from the device	bytes per second")
  hive_threads_daemon_thread_count                      =create_hive_metric_struct("hive_threads_daemon_thread_count", "JVM daemon thread count	threads")
  hive_threads_deadlocked_thread_count                  =create_hive_metric_struct("hive_threads_deadlocked_thread_count", "JVM deadlocked thread count	threads")
  hive_threads_thread_count                             =create_hive_metric_struct("hive_threads_thread_count", "JVM daemon and non-daemon thread count	threads")
  hive_unexpected_exits_rate                            =create_hive_metric_struct("hive_unexpected_exits_rate", "The number of times the role's backing process exited unexpectedly.	exits per second")
  hive_uptime                                           =create_hive_metric_struct("hive_uptime", "For a host, the amount of time since the host was booted. For a role, the uptime of the backing process.	seconds")
  hive_waiting_compile_ops                              =create_hive_metric_struct("hive_waiting_compile_ops", "The number of queries awaiting compilation	operations")
  hive_write_bytes_rate                                 =create_hive_metric_struct("hive_write_bytes_rate", "The number of bytes written to the device	bytes per second")
  hive_zookeeper_hive_exclusivelocks                    =create_hive_metric_struct("hive_zookeeper_hive_exclusivelocks", "Number of ZooKeeper exclusive locks held.	message.units.locks")
  hive_zookeeper_hive_semisharedlocks                   =create_hive_metric_struct("hive_zookeeper_hive_semisharedlocks", "Number of ZooKeeper semi-shared locks held.	message.units.locks")
  hive_zookeeper_hive_sharedlocks                       =create_hive_metric_struct("hive_zookeeper_hive_sharedlocks", "Number of ZooKeeper shared locks held.	message.units.locks")
)

// Creation of the structure that relates the queries with the descriptors of the Prometheus metrics
var hive_query_variable_relationship = []relation {
  {HIVE_ALERTS_RATE,                                     *hive_alerts_rate},
  {HIVE_API_ACQUIREREADWRITELOCKS_AVG,                   *hive_api_acquirereadwritelocks_avg},
  {HIVE_API_ACQUIREREADWRITELOCKS_MAX,                   *hive_api_acquirereadwritelocks_max},
  {HIVE_API_ACQUIREREADWRITELOCKS_MEDIAN,                *hive_api_acquirereadwritelocks_median},
  {HIVE_API_ACQUIREREADWRITELOCKS_MIN,                   *hive_api_acquirereadwritelocks_min},
  {HIVE_API_ACQUIREREADWRITELOCKS_RATE,                  *hive_api_acquirereadwritelocks_rate},
  {HIVE_API_ACQUIREREADWRITELOCKS_STDDEV,                *hive_api_acquirereadwritelocks_stddev},
  {HIVE_API_COMPILE_AVG,                                 *hive_api_compile_avg},
  {HIVE_API_COMPILE_MAX,                                 *hive_api_compile_max},
  {HIVE_API_COMPILE_MEDIAN,                              *hive_api_compile_median},
  {HIVE_API_COMPILE_MIN,                                 *hive_api_compile_min},
  {HIVE_API_COMPILE_RATE,                                *hive_api_compile_rate},
  {HIVE_API_COMPILE_STDDEV,                              *hive_api_compile_stddev},
  {HIVE_API_DRIVER_EXECUTE_AVG,                          *hive_api_driver_execute_avg},
  {HIVE_API_DRIVER_EXECUTE_MAX,                          *hive_api_driver_execute_max},
  {HIVE_API_DRIVER_EXECUTE_MEDIAN,                       *hive_api_driver_execute_median},
  {HIVE_API_DRIVER_EXECUTE_MIN,                          *hive_api_driver_execute_min},
  {HIVE_API_DRIVER_EXECUTE_RATE,                         *hive_api_driver_execute_rate},
  {HIVE_API_DRIVER_EXECUTE_STDDEV,                       *hive_api_driver_execute_stddev},
  {HIVE_API_DRIVER_RUN_AVG,                              *hive_api_driver_run_avg},
  {HIVE_API_DRIVER_RUN_MAX,                              *hive_api_driver_run_max},
  {HIVE_API_DRIVER_RUN_MEDIAN,                           *hive_api_driver_run_median},
  {HIVE_API_DRIVER_RUN_MIN,                              *hive_api_driver_run_min},
  {HIVE_API_DRIVER_RUN_RATE,                             *hive_api_driver_run_rate},
  {HIVE_API_DRIVER_RUN_STDDEV,                           *hive_api_driver_run_stddev},
  {HIVE_API_GET_ALL_TABLES_AVG,                          *hive_api_get_all_tables_avg},
  {HIVE_API_GET_ALL_TABLES_MAX,                          *hive_api_get_all_tables_max},
  {HIVE_API_GET_ALL_TABLES_MEDIAN,                       *hive_api_get_all_tables_median},
  {HIVE_API_GET_ALL_TABLES_MIN,                          *hive_api_get_all_tables_min},
  {HIVE_API_GET_ALL_TABLES_RATE,                         *hive_api_get_all_tables_rate},
  {HIVE_API_GET_ALL_TABLES_STDDEV,                       *hive_api_get_all_tables_stddev},
  {HIVE_API_GET_CONFIG_VALUE_AVG,                        *hive_api_get_config_value_avg},
  {HIVE_API_GET_CONFIG_VALUE_MAX,                        *hive_api_get_config_value_max},
  {HIVE_API_GET_CONFIG_VALUE_MEDIAN,                     *hive_api_get_config_value_median},
  {HIVE_API_GET_CONFIG_VALUE_MIN,                        *hive_api_get_config_value_min},
  {HIVE_API_GET_CONFIG_VALUE_RATE,                       *hive_api_get_config_value_rate},
  {HIVE_API_GET_CONFIG_VALUE_STDDEV,                     *hive_api_get_config_value_stddev},
  {HIVE_API_GET_DATABASE_AVG,                            *hive_api_get_database_avg},
  {HIVE_API_GET_DATABASE_MAX,                            *hive_api_get_database_max},
  {HIVE_API_GET_DATABASE_MEDIAN,                         *hive_api_get_database_median},
  {HIVE_API_GET_DATABASE_MIN,                            *hive_api_get_database_min},
  {HIVE_API_GET_DATABASE_RATE,                           *hive_api_get_database_rate},
  {HIVE_API_GET_DATABASE_STDDEV,                         *hive_api_get_database_stddev},
  {HIVE_API_GET_DATABASES_AVG,                           *hive_api_get_databases_avg},
  {HIVE_API_GET_DATABASES_MAX,                           *hive_api_get_databases_max},
  {HIVE_API_GET_DATABASES_MEDIAN,                        *hive_api_get_databases_median},
  {HIVE_API_GET_DATABASES_MIN,                           *hive_api_get_databases_min},
  {HIVE_API_GET_DATABASES_RATE,                          *hive_api_get_databases_rate},
  {HIVE_API_GET_DATABASES_STDDEV,                        *hive_api_get_databases_stddev},
  {HIVE_API_GET_FIELDS_WITH_ENVIRONMENT_CONTEXT_AVG,     *hive_api_get_fields_with_environment_context_avg},
  {HIVE_API_GET_FIELDS_WITH_ENVIRONMENT_CONTEXT_MAX,     *hive_api_get_fields_with_environment_context_max},
  {HIVE_API_GET_FIELDS_WITH_ENVIRONMENT_CONTEXT_MEDIAN,  *hive_api_get_fields_with_environment_context_median},
  {HIVE_API_GET_FIELDS_WITH_ENVIRONMENT_CONTEXT_MIN,     *hive_api_get_fields_with_environment_context_min},
  {HIVE_API_GET_FIELDS_WITH_ENVIRONMENT_CONTEXT_RATE,    *hive_api_get_fields_with_environment_context_rate},
  {HIVE_API_GET_FIELDS_WITH_ENVIRONMENT_CONTEXT_STDDEV,  *hive_api_get_fields_with_environment_context_stddev},
  {HIVE_API_GET_FUNCTION_AVG,                            *hive_api_get_function_avg},
  {HIVE_API_GET_FUNCTION_MAX,                            *hive_api_get_function_max},
  {HIVE_API_GET_FUNCTION_MEDIAN,                         *hive_api_get_function_median},
  {HIVE_API_GET_FUNCTION_MIN,                            *hive_api_get_function_min},
  {HIVE_API_GET_FUNCTION_RATE,                           *hive_api_get_function_rate},
  {HIVE_API_GET_FUNCTION_STDDEV,                         *hive_api_get_function_stddev},
  {HIVE_API_GET_FUNCTIONS_AVG,                           *hive_api_get_functions_avg},
  {HIVE_API_GET_FUNCTIONS_MAX,                           *hive_api_get_functions_max},
  {HIVE_API_GET_FUNCTIONS_MEDIAN,                        *hive_api_get_functions_median},
  {HIVE_API_GET_FUNCTIONS_MIN,                           *hive_api_get_functions_min},
  {HIVE_API_GET_FUNCTIONS_RATE,                          *hive_api_get_functions_rate},
  {HIVE_API_GET_FUNCTIONS_STDDEV,                        *hive_api_get_functions_stddev},
  {HIVE_API_GET_INDEX_BY_NAME_AVG,                       *hive_api_get_index_by_name_avg},
  {HIVE_API_GET_INDEX_BY_NAME_MAX,                       *hive_api_get_index_by_name_max},
  {HIVE_API_GET_INDEX_BY_NAME_MEDIAN,                    *hive_api_get_index_by_name_median},
  {HIVE_API_GET_INDEX_BY_NAME_MIN,                       *hive_api_get_index_by_name_min},
  {HIVE_API_GET_INDEX_BY_NAME_RATE,                      *hive_api_get_index_by_name_rate},
  {HIVE_API_GET_INDEX_BY_NAME_STDDEV,                    *hive_api_get_index_by_name_stddev},
  {HIVE_API_GET_INDEX_NAMES_AVG,                         *hive_api_get_index_names_avg},
  {HIVE_API_GET_INDEX_NAMES_MAX,                         *hive_api_get_index_names_max},
  {HIVE_API_GET_INDEX_NAMES_MEDIAN,                      *hive_api_get_index_names_median},
  {HIVE_API_GET_INDEX_NAMES_MIN,                         *hive_api_get_index_names_min},
  {HIVE_API_GET_INDEX_NAMES_RATE,                        *hive_api_get_index_names_rate},
  {HIVE_API_GET_INDEX_NAMES_STDDEV,                      *hive_api_get_index_names_stddev},
  {HIVE_API_GET_INDEXES_AVG,                             *hive_api_get_indexes_avg},
  {HIVE_API_GET_INDEXES_MAX,                             *hive_api_get_indexes_max},
  {HIVE_API_GET_INDEXES_MEDIAN,                          *hive_api_get_indexes_median},
  {HIVE_API_GET_INDEXES_MIN,                             *hive_api_get_indexes_min},
  {HIVE_API_GET_INDEXES_RATE,                            *hive_api_get_indexes_rate},
  {HIVE_API_GET_INDEXES_STDDEV,                          *hive_api_get_indexes_stddev},
  {HIVE_API_GET_MULTI_TABLE_AVG,                         *hive_api_get_multi_table_avg},
  {HIVE_API_GET_MULTI_TABLE_MAX,                         *hive_api_get_multi_table_max},
  {HIVE_API_GET_MULTI_TABLE_MEDIAN,                      *hive_api_get_multi_table_median},
  {HIVE_API_GET_MULTI_TABLE_MIN,                         *hive_api_get_multi_table_min},
  {HIVE_API_GET_MULTI_TABLE_RATE,                        *hive_api_get_multi_table_rate},
  {HIVE_API_GET_MULTI_TABLE_STDDEV,                      *hive_api_get_multi_table_stddev},
  {HIVE_API_GET_TABLE_AVG,                               *hive_api_get_table_avg},
  {HIVE_API_GET_TABLE_MAX,                               *hive_api_get_table_max},
  {HIVE_API_GET_TABLE_MEDIAN,                            *hive_api_get_table_median},
  {HIVE_API_GET_TABLE_MIN,                               *hive_api_get_table_min},
  {HIVE_API_GET_TABLE_NAMES_BY_FILTER_AVG,               *hive_api_get_table_names_by_filter_avg},
  {HIVE_API_GET_TABLE_NAMES_BY_FILTER_MAX,               *hive_api_get_table_names_by_filter_max},
  {HIVE_API_GET_TABLE_NAMES_BY_FILTER_MEDIAN,            *hive_api_get_table_names_by_filter_median},
  {HIVE_API_GET_TABLE_NAMES_BY_FILTER_MIN,               *hive_api_get_table_names_by_filter_min},
  {HIVE_API_GET_TABLE_NAMES_BY_FILTER_RATE,              *hive_api_get_table_names_by_filter_rate},
  {HIVE_API_GET_TABLE_NAMES_BY_FILTER_STDDEV,            *hive_api_get_table_names_by_filter_stddev},
  {HIVE_API_GET_TABLE_RATE,                              *hive_api_get_table_rate},
  {HIVE_API_GET_TABLE_STDDEV,                            *hive_api_get_table_stddev},
  {HIVE_API_GET_TABLES_AVG,                              *hive_api_get_tables_avg},
  {HIVE_API_GET_TABLES_MAX,                              *hive_api_get_tables_max},
  {HIVE_API_GET_TABLES_MEDIAN,                           *hive_api_get_tables_median},
  {HIVE_API_GET_TABLES_MIN,                              *hive_api_get_tables_min},
  {HIVE_API_GET_TABLES_RATE,                             *hive_api_get_tables_rate},
  {HIVE_API_GET_TABLES_STDDEV,                           *hive_api_get_tables_stddev},
  {HIVE_API_GETINPUTSUMMARY_AVG,                         *hive_api_getinputsummary_avg},
  {HIVE_API_GETINPUTSUMMARY_MAX,                         *hive_api_getinputsummary_max},
  {HIVE_API_GETINPUTSUMMARY_MEDIAN,                      *hive_api_getinputsummary_median},
  {HIVE_API_GETINPUTSUMMARY_MIN,                         *hive_api_getinputsummary_min},
  {HIVE_API_GETINPUTSUMMARY_RATE,                        *hive_api_getinputsummary_rate},
  {HIVE_API_GETINPUTSUMMARY_STDDEV,                      *hive_api_getinputsummary_stddev},
  {HIVE_API_GETSPLITS_AVG,                               *hive_api_getsplits_avg},
  {HIVE_API_GETSPLITS_MAX,                               *hive_api_getsplits_max},
  {HIVE_API_GETSPLITS_MEDIAN,                            *hive_api_getsplits_median},
  {HIVE_API_GETSPLITS_MIN,                               *hive_api_getsplits_min},
  {HIVE_API_GETSPLITS_RATE,                              *hive_api_getsplits_rate},
  {HIVE_API_GETSPLITS_STDDEV,                            *hive_api_getsplits_stddev},
  {HIVE_API_MARKPARTITIONFOREVENT_AVG,                   *hive_api_markpartitionforevent_avg},
  {HIVE_API_MARKPARTITIONFOREVENT_MAX,                   *hive_api_markpartitionforevent_max},
  {HIVE_API_MARKPARTITIONFOREVENT_MEDIAN,                *hive_api_markpartitionforevent_median},
  {HIVE_API_MARKPARTITIONFOREVENT_MIN,                   *hive_api_markpartitionforevent_min},
  {HIVE_API_MARKPARTITIONFOREVENT_RATE,                  *hive_api_markpartitionforevent_rate},
  {HIVE_API_MARKPARTITIONFOREVENT_STDDEV,                *hive_api_markpartitionforevent_stddev},
  {HIVE_API_OPERATION_INITIALIZED_AVG,                   *hive_api_operation_initialized_avg},
  {HIVE_API_OPERATION_INITIALIZED_MAX,                   *hive_api_operation_initialized_max},
  {HIVE_API_OPERATION_INITIALIZED_MEDIAN,                *hive_api_operation_initialized_median},
  {HIVE_API_OPERATION_INITIALIZED_MIN,                   *hive_api_operation_initialized_min},
  {HIVE_API_OPERATION_INITIALIZED_RATE,                  *hive_api_operation_initialized_rate},
  {HIVE_API_OPERATION_INITIALIZED_STDDEV,                *hive_api_operation_initialized_stddev},
  {HIVE_API_OPERATION_PENDING_AVG,                       *hive_api_operation_pending_avg},
  {HIVE_API_OPERATION_PENDING_MAX,                       *hive_api_operation_pending_max},
  {HIVE_API_OPERATION_PENDING_MEDIAN,                    *hive_api_operation_pending_median},
  {HIVE_API_OPERATION_PENDING_MIN,                       *hive_api_operation_pending_min},
  {HIVE_API_OPERATION_PENDING_RATE,                      *hive_api_operation_pending_rate},
  {HIVE_API_OPERATION_PENDING_STDDEV,                    *hive_api_operation_pending_stddev},
  {HIVE_API_OPERATION_RUNNING_AVG,                       *hive_api_operation_running_avg},
  {HIVE_API_OPERATION_RUNNING_MAX,                       *hive_api_operation_running_max},
  {HIVE_API_OPERATION_RUNNING_MEDIAN,                    *hive_api_operation_running_median},
  {HIVE_API_OPERATION_RUNNING_MIN,                       *hive_api_operation_running_min},
  {HIVE_API_OPERATION_RUNNING_RATE,                      *hive_api_operation_running_rate},
  {HIVE_API_OPERATION_RUNNING_STDDEV,                    *hive_api_operation_running_stddev},
  {HIVE_API_PARSE_AVG,                                   *hive_api_parse_avg},
  {HIVE_API_PARSE_MAX,                                   *hive_api_parse_max},
  {HIVE_API_PARSE_MEDIAN,                                *hive_api_parse_median},
  {HIVE_API_PARSE_MIN,                                   *hive_api_parse_min},
  {HIVE_API_PARSE_RATE,                                  *hive_api_parse_rate},
  {HIVE_API_PARSE_STDDEV,                                *hive_api_parse_stddev},
  {HIVE_API_PARTITION_RETRIEVING_AVG,                    *hive_api_partition_retrieving_avg},
  {HIVE_API_PARTITION_RETRIEVING_MAX,                    *hive_api_partition_retrieving_max},
  {HIVE_API_PARTITION_RETRIEVING_MEDIAN,                 *hive_api_partition_retrieving_median},
  {HIVE_API_PARTITION_RETRIEVING_MIN,                    *hive_api_partition_retrieving_min},
  {HIVE_API_PARTITION_RETRIEVING_RATE,                   *hive_api_partition_retrieving_rate},
  {HIVE_API_PARTITION_RETRIEVING_STDDEV,                 *hive_api_partition_retrieving_stddev},
  {HIVE_API_RELEASELOCKS_AVG,                            *hive_api_releaselocks_avg},
  {HIVE_API_RELEASELOCKS_MAX,                            *hive_api_releaselocks_max},
  {HIVE_API_RELEASELOCKS_MEDIAN,                         *hive_api_releaselocks_median},
  {HIVE_API_RELEASELOCKS_MIN,                            *hive_api_releaselocks_min},
  {HIVE_API_RELEASELOCKS_RATE,                           *hive_api_releaselocks_rate},
  {HIVE_API_RELEASELOCKS_STDDEV,                         *hive_api_releaselocks_stddev},
  {HIVE_API_RUNTASKS_AVG,                                *hive_api_runtasks_avg},
  {HIVE_API_RUNTASKS_MAX,                                *hive_api_runtasks_max},
  {HIVE_API_RUNTASKS_MEDIAN,                             *hive_api_runtasks_median},
  {HIVE_API_RUNTASKS_MIN,                                *hive_api_runtasks_min},
  {HIVE_API_RUNTASKS_RATE,                               *hive_api_runtasks_rate},
  {HIVE_API_RUNTASKS_STDDEV,                             *hive_api_runtasks_stddev},
  {HIVE_API_SERIALIZEPLAN_AVG,                           *hive_api_serializeplan_avg},
  {HIVE_API_SERIALIZEPLAN_MAX,                           *hive_api_serializeplan_max},
  {HIVE_API_SERIALIZEPLAN_MEDIAN,                        *hive_api_serializeplan_median},
  {HIVE_API_SERIALIZEPLAN_MIN,                           *hive_api_serializeplan_min},
  {HIVE_API_SERIALIZEPLAN_RATE,                          *hive_api_serializeplan_rate},
  {HIVE_API_SERIALIZEPLAN_STDDEV,                        *hive_api_serializeplan_stddev},
  {HIVE_API_TIMETOSUBMIT_AVG,                            *hive_api_timetosubmit_avg},
  {HIVE_API_TIMETOSUBMIT_MAX,                            *hive_api_timetosubmit_max},
  {HIVE_API_TIMETOSUBMIT_MEDIAN,                         *hive_api_timetosubmit_median},
  {HIVE_API_TIMETOSUBMIT_MIN,                            *hive_api_timetosubmit_min},
  {HIVE_API_TIMETOSUBMIT_RATE,                           *hive_api_timetosubmit_rate},
  {HIVE_API_TIMETOSUBMIT_STDDEV,                         *hive_api_timetosubmit_stddev},
  {HIVE_AUDIT_AGENT_BYTES_LEFT_RATE,                     *hive_audit_agent_bytes_left_rate},
  {HIVE_AUDIT_AGENT_BYTES_SENT_RATE,                     *hive_audit_agent_bytes_sent_rate},
  {HIVE_AUDIT_AGENT_EVENTS_SENT_RATE,                    *hive_audit_agent_events_sent_rate},
  {HIVE_AUDIT_AGENT_SEND_FAILURES_RATE,                  *hive_audit_agent_send_failures_rate},
  {HIVE_AUDIT_PLUGIN_COALESCED_RATE,                     *hive_audit_plugin_coalesced_rate},
  {HIVE_AUDIT_PLUGIN_EVENTS_RATE,                        *hive_audit_plugin_events_rate},
  {HIVE_AUDIT_PLUGIN_EXCEPTIONS_RATE,                    *hive_audit_plugin_exceptions_rate},
  {HIVE_AUDIT_PLUGIN_FILTERED_RATE,                      *hive_audit_plugin_filtered_rate},
  {HIVE_AUDIT_PLUGIN_PARSE_ERRORS_RATE,                  *hive_audit_plugin_parse_errors_rate},
  {HIVE_CANARY_DURATION,                                 *hive_canary_duration},
  {HIVE_CGROUP_CPU_SYSTEM_RATE,                          *hive_cgroup_cpu_system_rate},
  {HIVE_CGROUP_CPU_USER_RATE,                            *hive_cgroup_cpu_user_rate},
  {HIVE_CGROUP_MEM_PAGE_CACHE,                           *hive_cgroup_mem_page_cache},
  {HIVE_CGROUP_MEM_RSS,                                  *hive_cgroup_mem_rss},
  {HIVE_CGROUP_MEM_SWAP,                                 *hive_cgroup_mem_swap},
  {HIVE_CGROUP_READ_BYTES_RATE,                          *hive_cgroup_read_bytes_rate},
  {HIVE_CGROUP_READ_IOS_RATE,                            *hive_cgroup_read_ios_rate},
  {HIVE_CGROUP_WRITE_BYTES_RATE,                         *hive_cgroup_write_bytes_rate},
  {HIVE_CGROUP_WRITE_IOS_RATE,                           *hive_cgroup_write_ios_rate},
  {HIVE_COMPLETED_OPERATION_CANCELED_RATE,               *hive_completed_operation_canceled_rate},
  {HIVE_COMPLETED_OPERATION_CLOSED_RATE,                 *hive_completed_operation_closed_rate},
  {HIVE_COMPLETED_OPERATION_ERROR_RATE,                  *hive_completed_operation_error_rate},
  {HIVE_COMPLETED_OPERATION_FINISHED_RATE,               *hive_completed_operation_finished_rate},
  {HIVE_CPU_SYSTEM_RATE,                                 *hive_cpu_system_rate},
  {HIVE_CPU_USER_RATE,                                   *hive_cpu_user_rate},
  {HIVE_CREATE_TOTAL_COUNT_DBS_RATE,                     *hive_create_total_count_dbs_rate},
  {HIVE_CREATE_TOTAL_COUNT_PARTITIONS_RATE,              *hive_create_total_count_partitions_rate},
  {HIVE_CREATE_TOTAL_COUNT_TABLES_RATE,                  *hive_create_total_count_tables_rate},
  {HIVE_DELETE_TOTAL_COUNT_DBS_RATE,                     *hive_delete_total_count_dbs_rate},
  {HIVE_DELETE_TOTAL_COUNT_PARTITIONS_RATE,              *hive_delete_total_count_partitions_rate},
  {HIVE_DELETE_TOTAL_COUNT_TABLES_RATE,                  *hive_delete_total_count_tables_rate},
  {HIVE_EVENTS_CRITICAL_RATE,                            *hive_events_critical_rate},
  {HIVE_EVENTS_IMPORTANT_RATE,                           *hive_events_important_rate},
  {HIVE_EVENTS_INFORMATIONAL_RATE,                       *hive_events_informational_rate},
  {HIVE_EXEC_ASYNC_POOL_SIZE,                            *hive_exec_async_pool_size},
  {HIVE_EXEC_ASYNC_QUEUE_SIZE,                           *hive_exec_async_queue_size},
  {HIVE_FD_MAX,                                          *hive_fd_max},
  {HIVE_FD_OPEN,                                         *hive_fd_open},
  {HIVE_HEALTH_BAD_RATE,                                 *hive_health_bad_rate},
  {HIVE_HEALTH_CONCERNING_RATE,                          *hive_health_concerning_rate},
  {HIVE_HEALTH_DISABLED_RATE,                            *hive_health_disabled_rate},
  {HIVE_HEALTH_GOOD_RATE,                                *hive_health_good_rate},
  {HIVE_HEALTH_UNKNOWN_RATE,                             *hive_health_unknown_rate},
  {HIVE_INIT_TOTAL_COUNT_DBS,                            *hive_init_total_count_dbs},
  {HIVE_INIT_TOTAL_COUNT_PARTITIONS,                     *hive_init_total_count_partitions},
  {HIVE_INIT_TOTAL_COUNT_TABLES,                         *hive_init_total_count_tables},
  {HIVE_JVM_PAUSE_TIME_RATE,                             *hive_jvm_pause_time_rate},
  {HIVE_JVM_PAUSES_INFO_THRESHOLD_RATE,                  *hive_jvm_pauses_info_threshold_rate},
  {HIVE_JVM_PAUSES_WARN_THRESHOLD_RATE,                  *hive_jvm_pauses_warn_threshold_rate},
  {HIVE_MEM_RSS,                                         *hive_mem_rss},
  {HIVE_MEM_SWAP,                                        *hive_mem_swap},
  {HIVE_MEM_VIRTUAL,                                     *hive_mem_virtual},
  {HIVE_MEMORY_HEAP_COMMITTED,                           *hive_memory_heap_committed},
  {HIVE_MEMORY_HEAP_INIT,                                *hive_memory_heap_init},
  {HIVE_MEMORY_HEAP_MAX,                                 *hive_memory_heap_max},
  {HIVE_MEMORY_HEAP_USED,                                *hive_memory_heap_used},
  {HIVE_MEMORY_NON_HEAP_COMMITTED,                       *hive_memory_non_heap_committed},
  {HIVE_MEMORY_NON_HEAP_INIT,                            *hive_memory_non_heap_init},
  {HIVE_MEMORY_NON_HEAP_MAX,                             *hive_memory_non_heap_max},
  {HIVE_MEMORY_NON_HEAP_USED,                            *hive_memory_non_heap_used},
  {HIVE_MEMORY_TOTAL_COMMITTED,                          *hive_memory_total_committed},
  {HIVE_MEMORY_TOTAL_INIT,                               *hive_memory_total_init},
  {HIVE_MEMORY_TOTAL_MAX,                                *hive_memory_total_max},
  {HIVE_MEMORY_TOTAL_USED,                               *hive_memory_total_used},
  {HIVE_OOM_EXITS_RATE,                                  *hive_oom_exits_rate},
  {HIVE_OPEN_CONNECTIONS,                                *hive_open_connections},
  {HIVE_OPEN_OPERATIONS,                                 *hive_open_operations},
  {HIVE_READ_BYTES_RATE,                                 *hive_read_bytes_rate},
  {HIVE_THREADS_DAEMON_THREAD_COUNT,                     *hive_threads_daemon_thread_count},
  {HIVE_THREADS_DEADLOCKED_THREAD_COUNT,                 *hive_threads_deadlocked_thread_count},
  {HIVE_THREADS_THREAD_COUNT,                            *hive_threads_thread_count},
  {HIVE_UNEXPECTED_EXITS_RATE,                           *hive_unexpected_exits_rate},
  {HIVE_UPTIME,                                          *hive_uptime},
  {HIVE_WAITING_COMPILE_OPS,                             *hive_waiting_compile_ops},
  {HIVE_WRITE_BYTES_RATE,                                *hive_write_bytes_rate},
  {HIVE_ZOOKEEPER_HIVE_EXCLUSIVELOCKS,                   *hive_zookeeper_hive_exclusivelocks},
  {HIVE_ZOOKEEPER_HIVE_SEMISHAREDLOCKS,                  *hive_zookeeper_hive_semisharedlocks},
  {HIVE_ZOOKEEPER_HIVE_SHAREDLOCKS,                      *hive_zookeeper_hive_sharedlocks},
}




/* ======================================================================
 * Functions
 * ====================================================================== */
// Create and returns a prometheus descriptor for a hive metric. 
// The "metric_name" parameter its mandatory
// If the "description" parameter is empty, the function assings it with the
// value of the name of the metric in uppercase and separated by spaces
func create_hive_metric_struct(metric_name string, description string) *prometheus.Desc {
  // Correct "description" parameter if is empty
  if len(description) == 0 {
    description = strings.Replace(strings.ToUpper(metric_name), "_", " ", 0)
  }

  // return prometheus descriptor
  return prometheus.NewDesc(
    prometheus.BuildFQName(namespace, HIVE_SCRAPER_NAME, metric_name),
    description,
    []string{"cluster", "entityName","hostname"},
    nil,
  )
}


// Generic function to extract de metadata associated with the query value
// Only for HIVE metric type
func create_hive_metric (ctx context.Context, config Collector_connection_data, query string, metric_struct prometheus.Desc, ch chan<- prometheus.Metric) bool {
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
// ScrapeHIVE struct
type ScrapeHIVE struct{}

// Name of the Scraper. Should be unique.
func (ScrapeHIVE) Name() string {
  return HIVE_SCRAPER_NAME
}

// Help describes the role of the Scraper.
func (ScrapeHIVE) Help() string {
  return "HIVE Metrics"
}

// Version.
func (ScrapeHIVE) Version() float64 {
  return 1.0
}

// Scrape generic function. Override for host module.
func (ScrapeHIVE) Scrape (ctx context.Context, config *Collector_connection_data, ch chan<- prometheus.Metric) error {
  log.Debug_msg("Executing HIVE Metrics Scraper")

  // Queries counters
  success_queries := 0
  error_queries := 0

  // Execute the generic funtion for creation of metrics with the pairs (QUERY, PROM:DESCRIPTOR)
  for i:=0 ; i < len(hive_query_variable_relationship) ; i++ {
    if create_hive_metric(ctx, *config, hive_query_variable_relationship[i].Query, hive_query_variable_relationship[i].Metric_struct, ch) {
      success_queries += 1
    } else {
      error_queries += 1
    }
  }
  log.Info_msg("In the HIVE Module has been executed %d queries. %d success and %d with errors", success_queries + error_queries, success_queries, error_queries)
  return nil
}

var _ Scraper = ScrapeHIVE{}
