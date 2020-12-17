/*
 *
 * title           :collector/hue_module.go
 * description     :Submodule Collector for the Cluster HUE metrics
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
const HUE_SCRAPER_NAME = "hue"
var (
  // Agent Queries
  HUE_ALERTS_RATE                                     ="SELECT LAST(alerts_rate) WHERE serviceType=\"HUE\""
  HUE_CGROUP_CPU_SYSTEM_RATE                          ="SELECT LAST(cgroup_cpu_system_rate) WHERE serviceType=\"HUE\""
  HUE_CGROUP_CPU_USER_RATE                            ="SELECT LAST(cgroup_cpu_user_rate) WHERE serviceType=\"HUE\""
  HUE_CGROUP_MEM_PAGE_CACHE                           ="SELECT LAST(cgroup_mem_page_cache) WHERE serviceType=\"HUE\""
  HUE_CGROUP_MEM_RSS                                  ="SELECT LAST(cgroup_mem_rss) WHERE serviceType=\"HUE\""
  HUE_CGROUP_MEM_SWAP                                 ="SELECT LAST(cgroup_mem_swap) WHERE serviceType=\"HUE\""
  HUE_CGROUP_READ_BYTES_RATE                          ="SELECT LAST(cgroup_read_bytes_rate) WHERE serviceType=\"HUE\""
  HUE_CGROUP_READ_IOS_RATE                            ="SELECT LAST(cgroup_read_ios_rate) WHERE serviceType=\"HUE\""
  HUE_CGROUP_WRITE_BYTES_RATE                         ="SELECT LAST(cgroup_write_bytes_rate) WHERE serviceType=\"HUE\""
  HUE_CGROUP_WRITE_IOS_RATE                           ="SELECT LAST(cgroup_write_ios_rate) WHERE serviceType=\"HUE\""
  HUE_CPU_SYSTEM_RATE                                 ="SELECT LAST(cpu_system_rate) WHERE serviceType=\"HUE\""
  HUE_CPU_USER_RATE                                   ="SELECT LAST(cpu_user_rate) WHERE serviceType=\"HUE\""
  HUE_EVENTS_CRITICAL_RATE                            ="SELECT LAST(events_critical_rate) WHERE serviceType=\"HUE\""
  HUE_EVENTS_IMPORTANT_RATE                           ="SELECT LAST(events_important_rate) WHERE serviceType=\"HUE\""
  HUE_EVENTS_INFORMATIONAL_RATE                       ="SELECT LAST(events_informational_rate) WHERE serviceType=\"HUE\""
  HUE_FD_MAX                                          ="SELECT LAST(fd_max) WHERE serviceType=\"HUE\""
  HUE_FD_OPEN                                         ="SELECT LAST(fd_open) WHERE serviceType=\"HUE\""
  HUE_HEALTH_BAD_RATE                                 ="SELECT LAST(health_bad_rate) WHERE serviceType=\"HUE\""
  HUE_HEALTH_CONCERNING_RATE                          ="SELECT LAST(health_concerning_rate) WHERE serviceType=\"HUE\""
  HUE_HEALTH_DISABLED_RATE                            ="SELECT LAST(health_disabled_rate) WHERE serviceType=\"HUE\""
  HUE_HEALTH_GOOD_RATE                                ="SELECT LAST(health_good_rate) WHERE serviceType=\"HUE\""
  HUE_HEALTH_UNKNOWN_RATE                             ="SELECT LAST(health_unknown_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_15M_RATE                ="SELECT LAST(hue_auth_ldap_auth_time_15m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_1M_RATE                 ="SELECT LAST(hue_auth_ldap_auth_time_1m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_5M_RATE                 ="SELECT LAST(hue_auth_ldap_auth_time_5m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_75_PERCENTILE           ="SELECT LAST(hue_auth_ldap_auth_time_75_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_95_PERCENTILE           ="SELECT LAST(hue_auth_ldap_auth_time_95_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_999_PERCENTILE          ="SELECT LAST(hue_auth_ldap_auth_time_999_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_99_PERCENTILE           ="SELECT LAST(hue_auth_ldap_auth_time_99_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_AVG                     ="SELECT LAST(hue_auth_ldap_auth_time_avg) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_MAX                     ="SELECT LAST(hue_auth_ldap_auth_time_max) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_MEDIAN                  ="SELECT LAST(hue_auth_ldap_auth_time_median) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_MIN                     ="SELECT LAST(hue_auth_ldap_auth_time_min) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_RATE                    ="SELECT LAST(hue_auth_ldap_auth_time_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_STD_DEV                 ="SELECT LAST(hue_auth_ldap_auth_time_std_dev) WHERE serviceType=\"HUE\""
  HUE_AUTH_LDAP_AUTH_TIME_SUM_RATE                ="SELECT LAST(hue_auth_ldap_auth_time_sum_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_15M_RATE               ="SELECT LAST(hue_auth_oauth_auth_time_15m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_1M_RATE                ="SELECT LAST(hue_auth_oauth_auth_time_1m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_5M_RATE                ="SELECT LAST(hue_auth_oauth_auth_time_5m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_75_PERCENTILE          ="SELECT LAST(hue_auth_oauth_auth_time_75_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_95_PERCENTILE          ="SELECT LAST(hue_auth_oauth_auth_time_95_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_999_PERCENTILE         ="SELECT LAST(hue_auth_oauth_auth_time_999_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_99_PERCENTILE          ="SELECT LAST(hue_auth_oauth_auth_time_99_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_AVG                    ="SELECT LAST(hue_auth_oauth_auth_time_avg) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_MAX                    ="SELECT LAST(hue_auth_oauth_auth_time_max) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_MEDIAN                 ="SELECT LAST(hue_auth_oauth_auth_time_median) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_MIN                    ="SELECT LAST(hue_auth_oauth_auth_time_min) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_RATE                   ="SELECT LAST(hue_auth_oauth_auth_time_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_STD_DEV                ="SELECT LAST(hue_auth_oauth_auth_time_std_dev) WHERE serviceType=\"HUE\""
  HUE_AUTH_OAUTH_AUTH_TIME_SUM_RATE               ="SELECT LAST(hue_auth_oauth_auth_time_sum_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_15M_RATE              ="SELECT LAST(hue_auth_openid_auth_time_15m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_1M_RATE               ="SELECT LAST(hue_auth_openid_auth_time_1m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_5M_RATE               ="SELECT LAST(hue_auth_openid_auth_time_5m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_75_PERCENTILE         ="SELECT LAST(hue_auth_openid_auth_time_75_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_95_PERCENTILE         ="SELECT LAST(hue_auth_openid_auth_time_95_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_999_PERCENTILE        ="SELECT LAST(hue_auth_openid_auth_time_999_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_99_PERCENTILE         ="SELECT LAST(hue_auth_openid_auth_time_99_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_AVG                   ="SELECT LAST(hue_auth_openid_auth_time_avg) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_MAX                   ="SELECT LAST(hue_auth_openid_auth_time_max) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_MEDIAN                ="SELECT LAST(hue_auth_openid_auth_time_median) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_MIN                   ="SELECT LAST(hue_auth_openid_auth_time_min) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_RATE                  ="SELECT LAST(hue_auth_openid_auth_time_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_STD_DEV               ="SELECT LAST(hue_auth_openid_auth_time_std_dev) WHERE serviceType=\"HUE\""
  HUE_AUTH_OPENID_AUTH_TIME_SUM_RATE              ="SELECT LAST(hue_auth_openid_auth_time_sum_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_15M_RATE                 ="SELECT LAST(hue_auth_pam_auth_time_15m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_1M_RATE                  ="SELECT LAST(hue_auth_pam_auth_time_1m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_5M_RATE                  ="SELECT LAST(hue_auth_pam_auth_time_5m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_75_PERCENTILE            ="SELECT LAST(hue_auth_pam_auth_time_75_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_95_PERCENTILE            ="SELECT LAST(hue_auth_pam_auth_time_95_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_999_PERCENTILE           ="SELECT LAST(hue_auth_pam_auth_time_999_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_99_PERCENTILE            ="SELECT LAST(hue_auth_pam_auth_time_99_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_AVG                      ="SELECT LAST(hue_auth_pam_auth_time_avg) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_MAX                      ="SELECT LAST(hue_auth_pam_auth_time_max) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_MEDIAN                   ="SELECT LAST(hue_auth_pam_auth_time_median) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_MIN                      ="SELECT LAST(hue_auth_pam_auth_time_min) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_RATE                     ="SELECT LAST(hue_auth_pam_auth_time_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_STD_DEV                  ="SELECT LAST(hue_auth_pam_auth_time_std_dev) WHERE serviceType=\"HUE\""
  HUE_AUTH_PAM_AUTH_TIME_SUM_RATE                 ="SELECT LAST(hue_auth_pam_auth_time_sum_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_15M_RATE               ="SELECT LAST(hue_auth_saml2_auth_time_15m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_1M_RATE                ="SELECT LAST(hue_auth_saml2_auth_time_1m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_5M_RATE                ="SELECT LAST(hue_auth_saml2_auth_time_5m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_75_PERCENTILE          ="SELECT LAST(hue_auth_saml2_auth_time_75_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_95_PERCENTILE          ="SELECT LAST(hue_auth_saml2_auth_time_95_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_999_PERCENTILE         ="SELECT LAST(hue_auth_saml2_auth_time_999_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_99_PERCENTILE          ="SELECT LAST(hue_auth_saml2_auth_time_99_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_AVG                    ="SELECT LAST(hue_auth_saml2_auth_time_avg) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_MAX                    ="SELECT LAST(hue_auth_saml2_auth_time_max) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_MEDIAN                 ="SELECT LAST(hue_auth_saml2_auth_time_median) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_MIN                    ="SELECT LAST(hue_auth_saml2_auth_time_min) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_RATE                   ="SELECT LAST(hue_auth_saml2_auth_time_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_STD_DEV                ="SELECT LAST(hue_auth_saml2_auth_time_std_dev) WHERE serviceType=\"HUE\""
  HUE_AUTH_SAML2_AUTH_TIME_SUM_RATE               ="SELECT LAST(hue_auth_saml2_auth_time_sum_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_15M_RATE              ="SELECT LAST(hue_auth_spnego_auth_time_15m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_1M_RATE               ="SELECT LAST(hue_auth_spnego_auth_time_1m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_5M_RATE               ="SELECT LAST(hue_auth_spnego_auth_time_5m_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_75_PERCENTILE         ="SELECT LAST(hue_auth_spnego_auth_time_75_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_95_PERCENTILE         ="SELECT LAST(hue_auth_spnego_auth_time_95_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_999_PERCENTILE        ="SELECT LAST(hue_auth_spnego_auth_time_999_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_99_PERCENTILE         ="SELECT LAST(hue_auth_spnego_auth_time_99_percentile) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_AVG                   ="SELECT LAST(hue_auth_spnego_auth_time_avg) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_MAX                   ="SELECT LAST(hue_auth_spnego_auth_time_max) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_MEDIAN                ="SELECT LAST(hue_auth_spnego_auth_time_median) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_MIN                   ="SELECT LAST(hue_auth_spnego_auth_time_min) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_RATE                  ="SELECT LAST(hue_auth_spnego_auth_time_rate) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_STD_DEV               ="SELECT LAST(hue_auth_spnego_auth_time_std_dev) WHERE serviceType=\"HUE\""
  HUE_AUTH_SPNEGO_AUTH_TIME_SUM_RATE              ="SELECT LAST(hue_auth_spnego_auth_time_sum_rate) WHERE serviceType=\"HUE\""
  HUE_MULTIPROCESSING_PROCESSES_DAEMON            ="SELECT LAST(hue_multiprocessing_processes_daemon) WHERE serviceType=\"HUE\""
  HUE_MULTIPROCESSING_PROCESSES_TOTAL             ="SELECT LAST(hue_multiprocessing_processes_total) WHERE serviceType=\"HUE\""
  HUE_PYTHON_GC_GENERATION_0                      ="SELECT LAST(hue_python_gc_generation_0) WHERE serviceType=\"HUE\""
  HUE_PYTHON_GC_GENERATION_1                      ="SELECT LAST(hue_python_gc_generation_1) WHERE serviceType=\"HUE\""
  HUE_PYTHON_GC_GENERATION_2                      ="SELECT LAST(hue_python_gc_generation_2) WHERE serviceType=\"HUE\""
  HUE_PYTHON_GC_OBJECTS                           ="SELECT LAST(hue_python_gc_objects) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_ACTIVE                             ="SELECT LAST(hue_requests_active) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_EXCEPTIONS_RATE                    ="SELECT LAST(hue_requests_exceptions_rate) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_15M_RATE             ="SELECT LAST(hue_requests_response_time_15m_rate) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_1M_RATE              ="SELECT LAST(hue_requests_response_time_1m_rate) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_5M_RATE              ="SELECT LAST(hue_requests_response_time_5m_rate) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_75_PERCENTILE        ="SELECT LAST(hue_requests_response_time_75_percentile) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_95_PERCENTILE        ="SELECT LAST(hue_requests_response_time_95_percentile) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_999_PERCENTILE       ="SELECT LAST(hue_requests_response_time_999_percentile) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_99_PERCENTILE        ="SELECT LAST(hue_requests_response_time_99_percentile) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_AVG                  ="SELECT LAST(hue_requests_response_time_avg) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_MAX                  ="SELECT LAST(hue_requests_response_time_max) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_MEDIAN               ="SELECT LAST(hue_requests_response_time_median) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_MIN                  ="SELECT LAST(hue_requests_response_time_min) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_RATE                 ="SELECT LAST(hue_requests_response_time_rate) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_STD_DEV              ="SELECT LAST(hue_requests_response_time_std_dev) WHERE serviceType=\"HUE\""
  HUE_REQUESTS_RESPONSE_TIME_SUM_RATE             ="SELECT LAST(hue_requests_response_time_sum_rate) WHERE serviceType=\"HUE\""
  HUE_THREADS_DAEMON                              ="SELECT LAST(hue_threads_daemon) WHERE serviceType=\"HUE\""
  HUE_THREADS_TOTAL                               ="SELECT LAST(hue_threads_total) WHERE serviceType=\"HUE\""
  HUE_USERS                                       ="SELECT LAST(hue_users) WHERE serviceType=\"HUE\""
  HUE_USERS_ACTIVE                                ="SELECT LAST(hue_users_active) WHERE serviceType=\"HUE\""
  HUE_MEM_RSS                                         ="SELECT LAST(mem_rss) WHERE serviceType=\"HUE\""
  HUE_MEM_SWAP                                        ="SELECT LAST(mem_swap) WHERE serviceType=\"HUE\""
  HUE_MEM_VIRTUAL                                     ="SELECT LAST(mem_virtual) WHERE serviceType=\"HUE\""
  HUE_OOM_EXITS_RATE                                  ="SELECT LAST(oom_exits_rate) WHERE serviceType=\"HUE\""
  HUE_READ_BYTES_RATE                                 ="SELECT LAST(read_bytes_rate) WHERE serviceType=\"HUE\""
  HUE_UNEXPECTED_EXITS_RATE                           ="SELECT LAST(unexpected_exits_rate) WHERE serviceType=\"HUE\""
  HUE_UPTIME                                          ="SELECT LAST(uptime) WHERE serviceType=\"HUE\""
  HUE_WEB_METRICS_COLLECTION_DURATION                 ="SELECT LAST(web_metrics_collection_duration) WHERE serviceType=\"HUE\""
  HUE_WRITE_BYTES_RATE                                ="SELECT LAST(write_bytes_rate) WHERE serviceType=\"HUE\""

)




/* ======================================================================
 * Global variables
 * ====================================================================== */
// Prometheus data Descriptors for the metrics to export
var (
  // Agent Metrics
  hue_alerts_rate                               =create_hue_metric_struct("hue_alerts_rate",	"The number of alerts.	events per second")
  hue_cgroup_cpu_system_rate                    =create_hue_metric_struct("hue_cgroup_cpu_system_rate",	"CPU usage of the role's cgroup	seconds per second")
  hue_cgroup_cpu_user_rate                      =create_hue_metric_struct("hue_cgroup_cpu_user_rate",	"User Space CPU usage of the role's cgroup	seconds per second")
  hue_cgroup_mem_page_cache                     =create_hue_metric_struct("hue_cgroup_mem_page_cache",	"Page cache usage of the role's cgroup	bytes")
  hue_cgroup_mem_rss                            =create_hue_metric_struct("hue_cgroup_mem_rss",	"Resident memory of the role's cgroup	bytes")
  hue_cgroup_mem_swap                           =create_hue_metric_struct("hue_cgroup_mem_swap",	"Swap usage of the role's cgroup	bytes")
  hue_cgroup_read_bytes_rate                    =create_hue_metric_struct("hue_cgroup_read_bytes_rate",	"Bytes read from all disks by the role's cgroup	bytes per second")
  hue_cgroup_read_ios_rate                      =create_hue_metric_struct("hue_cgroup_read_ios_rate",	"Number of read I/O operations from all disks by the role's cgroup	ios per second")
  hue_cgroup_write_bytes_rate                   =create_hue_metric_struct("hue_cgroup_write_bytes_rate",	"Bytes written to all disks by the role's cgroup	bytes per second")
  hue_cgroup_write_ios_rate                     =create_hue_metric_struct("hue_cgroup_write_ios_rate",	"Number of write I/O operations to all disks by the role's cgroup	ios per second")
  hue_cpu_system_rate                           =create_hue_metric_struct("hue_cpu_system_rate",	"Total System CPU	seconds per second")
  hue_cpu_user_rate                             =create_hue_metric_struct("hue_cpu_user_rate",	"Total CPU user time	seconds per second")
  hue_events_critical_rate                      =create_hue_metric_struct("hue_events_critical_rate",	"The number of critical events.	events per second")
  hue_events_important_rate                     =create_hue_metric_struct("hue_events_important_rate",	"The number of important events.	events per second")
  hue_events_informational_rate                 =create_hue_metric_struct("hue_events_informational_rate",	"The number of informational events.	events per second")
  hue_fd_max                                    =create_hue_metric_struct("hue_fd_max",	"Maximum number of file descriptors	file descriptors")
  hue_fd_open                                   =create_hue_metric_struct("hue_fd_open",	"Open file descriptors.	file descriptors")
  hue_health_bad_rate                           =create_hue_metric_struct("hue_health_bad_rate",	"Percentage of Time with Bad Health	seconds per second")
  hue_health_concerning_rate                    =create_hue_metric_struct("hue_health_concerning_rate",	"Percentage of Time with Concerning Health	seconds per second")
  hue_health_disabled_rate                      =create_hue_metric_struct("hue_health_disabled_rate",	"Percentage of Time with Disabled Health	seconds per second")
  hue_health_good_rate                          =create_hue_metric_struct("hue_health_good_rate",	"Percentage of Time with Good Health	seconds per second")
  hue_health_unknown_rate                       =create_hue_metric_struct("hue_health_unknown_rate",	"Percentage of Time with Unknown Health	seconds per second")
  hue_auth_ldap_auth_time_15m_rate              =create_hue_metric_struct("hue_auth_ldap_auth_time_15m_rate",	"The time spent waiting for LDAP to authenticate a user: 15-Minute Rate	message.units.authentications per second")
  hue_auth_ldap_auth_time_1m_rate               =create_hue_metric_struct("hue_auth_ldap_auth_time_1m_rate",	"The time spent waiting for LDAP to authenticate a user: 1-Minute moving average rate.	message.units.authentications per second")
  hue_auth_ldap_auth_time_5m_rate               =create_hue_metric_struct("hue_auth_ldap_auth_time_5m_rate",	"The time spent waiting for LDAP to authenticate a user: 5-Minute Rate	message.units.authentications per second")
  hue_auth_ldap_auth_time_75_percentile         =create_hue_metric_struct("hue_auth_ldap_auth_time_75_percentile",	"The time spent waiting for LDAP to authenticate a user: 75th Percentile. This is computed over the past hour.	seconds")
  hue_auth_ldap_auth_time_95_percentile         =create_hue_metric_struct("hue_auth_ldap_auth_time_95_percentile",	"The time spent waiting for LDAP to authenticate a user: 95th Percentile. This is computed over the past hour.	seconds")
  hue_auth_ldap_auth_time_999_percentile        =create_hue_metric_struct("hue_auth_ldap_auth_time_999_percentile",	"The time spent waiting for LDAP to authenticate a user: 999th Percentile. This is computed over the past hour.	seconds")
  hue_auth_ldap_auth_time_99_percentile         =create_hue_metric_struct("hue_auth_ldap_auth_time_99_percentile",	"The time spent waiting for LDAP to authenticate a user: 99th Percentile. This is computed over the past hour.	seconds")
  hue_auth_ldap_auth_time_avg                   =create_hue_metric_struct("hue_auth_ldap_auth_time_avg",	"The time spent waiting for LDAP to authenticate a user: Average. This is computed over the lifetime of the process.	seconds")
  hue_auth_ldap_auth_time_max                   =create_hue_metric_struct("hue_auth_ldap_auth_time_max",	"The time spent waiting for LDAP to authenticate a user: Max. This is computed over the lifetime of the process.	seconds")
  hue_auth_ldap_auth_time_median                =create_hue_metric_struct("hue_auth_ldap_auth_time_median",	"The time spent waiting for LDAP to authenticate a user: 50th Percentile. This is computed over the past hour.	seconds")
  hue_auth_ldap_auth_time_min                   =create_hue_metric_struct("hue_auth_ldap_auth_time_min",	"The time spent waiting for LDAP to authenticate a user: Min. This is computed over the lifetime of the process.	seconds")
  hue_auth_ldap_auth_time_rate                  =create_hue_metric_struct("hue_auth_ldap_auth_time_rate",	"The time spent waiting for LDAP to authenticate a user: Sample Count. This is computed over the lifetime of the process.	message.units.authentications per second")
  hue_auth_ldap_auth_time_std_dev               =create_hue_metric_struct("hue_auth_ldap_auth_time_std_dev",	"The time spent waiting for LDAP to authenticate a user: Standard Deviation. This is computed over the lifetime of the process.	seconds")
  hue_auth_ldap_auth_time_sum_rate              =create_hue_metric_struct("hue_auth_ldap_auth_time_sum_rate",	"The time spent waiting for LDAP to authenticate a user: Sample Sum. This is computed over the lifetime of the process.	message.units.authentications per second")
  hue_auth_oauth_auth_time_15m_rate             =create_hue_metric_struct("hue_auth_oauth_auth_time_15m_rate",	"The time spent waiting for OAUTH to authenticate a user: 15-Minute Rate	message.units.authentications per second")
  hue_auth_oauth_auth_time_1m_rate              =create_hue_metric_struct("hue_auth_oauth_auth_time_1m_rate",	"The time spent waiting for OAUTH to authenticate a user: 1-Minute moving average rate.	message.units.authentications per second")
  hue_auth_oauth_auth_time_5m_rate              =create_hue_metric_struct("hue_auth_oauth_auth_time_5m_rate",	"The time spent waiting for OAUTH to authenticate a user: 5-Minute Rate	message.units.authentications per second")
  hue_auth_oauth_auth_time_75_percentile        =create_hue_metric_struct("hue_auth_oauth_auth_time_75_percentile",	"The time spent waiting for OAUTH to authenticate a user: 75th Percentile. This is computed over the past hour.	seconds")
  hue_auth_oauth_auth_time_95_percentile        =create_hue_metric_struct("hue_auth_oauth_auth_time_95_percentile",	"The time spent waiting for OAUTH to authenticate a user: 95th Percentile. This is computed over the past hour.	seconds")
  hue_auth_oauth_auth_time_999_percentile       =create_hue_metric_struct("hue_auth_oauth_auth_time_999_percentile",	"The time spent waiting for OAUTH to authenticate a user: 999th Percentile. This is computed over the past hour.	seconds")
  hue_auth_oauth_auth_time_99_percentile        =create_hue_metric_struct("hue_auth_oauth_auth_time_99_percentile",	"The time spent waiting for OAUTH to authenticate a user: 99th Percentile. This is computed over the past hour.	seconds")
  hue_auth_oauth_auth_time_avg                  =create_hue_metric_struct("hue_auth_oauth_auth_time_avg",	"The time spent waiting for OAUTH to authenticate a user: Average. This is computed over the lifetime of the process.	seconds")
  hue_auth_oauth_auth_time_max                  =create_hue_metric_struct("hue_auth_oauth_auth_time_max",	"The time spent waiting for OAUTH to authenticate a user: Max. This is computed over the lifetime of the process.	seconds")
  hue_auth_oauth_auth_time_median               =create_hue_metric_struct("hue_auth_oauth_auth_time_median",	"The time spent waiting for OAUTH to authenticate a user: 50th Percentile. This is computed over the past hour.	seconds")
  hue_auth_oauth_auth_time_min                  =create_hue_metric_struct("hue_auth_oauth_auth_time_min",	"The time spent waiting for OAUTH to authenticate a user: Min. This is computed over the lifetime of the process.	seconds")
  hue_auth_oauth_auth_time_rate                 =create_hue_metric_struct("hue_auth_oauth_auth_time_rate",	"The time spent waiting for OAUTH to authenticate a user: Sample Count. This is computed over the lifetime of the process.	message.units.authentications per second")
  hue_auth_oauth_auth_time_std_dev              =create_hue_metric_struct("hue_auth_oauth_auth_time_std_dev",	"The time spent waiting for OAUTH to authenticate a user: Standard Deviation. This is computed over the lifetime of the process.	seconds")
  hue_auth_oauth_auth_time_sum_rate             =create_hue_metric_struct("hue_auth_oauth_auth_time_sum_rate",	"The time spent waiting for OAUTH to authenticate a user: Sample Sum. This is computed over the lifetime of the process.	message.units.authentications per second")
  hue_auth_openid_auth_time_15m_rate            =create_hue_metric_struct("hue_auth_openid_auth_time_15m_rate",	"The time spent waiting for OpenID to authenticate a user: 15-Minute Rate	message.units.authentications per second")
  hue_auth_openid_auth_time_1m_rate             =create_hue_metric_struct("hue_auth_openid_auth_time_1m_rate",	"The time spent waiting for OpenID to authenticate a user: 1-Minute moving average rate.	message.units.authentications per second")
  hue_auth_openid_auth_time_5m_rate             =create_hue_metric_struct("hue_auth_openid_auth_time_5m_rate",	"The time spent waiting for OpenID to authenticate a user: 5-Minute Rate	message.units.authentications per second")
  hue_auth_openid_auth_time_75_percentile       =create_hue_metric_struct("hue_auth_openid_auth_time_75_percentile",	"The time spent waiting for OpenID to authenticate a user: 75th Percentile. This is computed over the past hour.	seconds")
  hue_auth_openid_auth_time_95_percentile       =create_hue_metric_struct("hue_auth_openid_auth_time_95_percentile",	"The time spent waiting for OpenID to authenticate a user: 95th Percentile. This is computed over the past hour.	seconds")
  hue_auth_openid_auth_time_999_percentile      =create_hue_metric_struct("hue_auth_openid_auth_time_999_percentile",	"The time spent waiting for OpenID to authenticate a user: 999th Percentile. This is computed over the past hour.	seconds")
  hue_auth_openid_auth_time_99_percentile       =create_hue_metric_struct("hue_auth_openid_auth_time_99_percentile",	"The time spent waiting for OpenID to authenticate a user: 99th Percentile. This is computed over the past hour.	seconds")
  hue_auth_openid_auth_time_avg                 =create_hue_metric_struct("hue_auth_openid_auth_time_avg",	"The time spent waiting for OpenID to authenticate a user: Average. This is computed over the lifetime of the process.	seconds")
  hue_auth_openid_auth_time_max                 =create_hue_metric_struct("hue_auth_openid_auth_time_max",	"The time spent waiting for OpenID to authenticate a user: Max. This is computed over the lifetime of the process.	seconds")
  hue_auth_openid_auth_time_median              =create_hue_metric_struct("hue_auth_openid_auth_time_median",	"The time spent waiting for OpenID to authenticate a user: 50th Percentile. This is computed over the past hour.	seconds")
  hue_auth_openid_auth_time_min                 =create_hue_metric_struct("hue_auth_openid_auth_time_min",	"The time spent waiting for OpenID to authenticate a user: Min. This is computed over the lifetime of the process.	seconds")
  hue_auth_openid_auth_time_rate                =create_hue_metric_struct("hue_auth_openid_auth_time_rate",	"The time spent waiting for OpenID to authenticate a user: Sample Count. This is computed over the lifetime of the process.	message.units.authentications per second")
  hue_auth_openid_auth_time_std_dev             =create_hue_metric_struct("hue_auth_openid_auth_time_std_dev",	"The time spent waiting for OpenID to authenticate a user: Standard Deviation. This is computed over the lifetime of the process.	seconds")
  hue_auth_openid_auth_time_sum_rate            =create_hue_metric_struct("hue_auth_openid_auth_time_sum_rate",	"The time spent waiting for OpenID to authenticate a user: Sample Sum. This is computed over the lifetime of the process.	message.units.authentications per second")
  hue_auth_pam_auth_time_15m_rate               =create_hue_metric_struct("hue_auth_pam_auth_time_15m_rate",	"The time spent waiting for PAM to authenticate a user: 15-Minute Rate	message.units.authentications per second")
  hue_auth_pam_auth_time_1m_rate                =create_hue_metric_struct("hue_auth_pam_auth_time_1m_rate",	"The time spent waiting for PAM to authenticate a user: 1-Minute moving average rate.	message.units.authentications per second")
  hue_auth_pam_auth_time_5m_rate                =create_hue_metric_struct("hue_auth_pam_auth_time_5m_rate",	"The time spent waiting for PAM to authenticate a user: 5-Minute Rate	message.units.authentications per second")
  hue_auth_pam_auth_time_75_percentile          =create_hue_metric_struct("hue_auth_pam_auth_time_75_percentile",	"The time spent waiting for PAM to authenticate a user: 75th Percentile. This is computed over the past hour.	seconds")
  hue_auth_pam_auth_time_95_percentile          =create_hue_metric_struct("hue_auth_pam_auth_time_95_percentile",	"The time spent waiting for PAM to authenticate a user: 95th Percentile. This is computed over the past hour.	seconds")
  hue_auth_pam_auth_time_999_percentile         =create_hue_metric_struct("hue_auth_pam_auth_time_999_percentile",	"The time spent waiting for PAM to authenticate a user: 999th Percentile. This is computed over the past hour.	seconds")
  hue_auth_pam_auth_time_99_percentile          =create_hue_metric_struct("hue_auth_pam_auth_time_99_percentile",	"The time spent waiting for PAM to authenticate a user: 99th Percentile. This is computed over the past hour.	seconds")
  hue_auth_pam_auth_time_avg                    =create_hue_metric_struct("hue_auth_pam_auth_time_avg",	"The time spent waiting for PAM to authenticate a user: Average. This is computed over the lifetime of the process.	seconds")
  hue_auth_pam_auth_time_max                    =create_hue_metric_struct("hue_auth_pam_auth_time_max",	"The time spent waiting for PAM to authenticate a user: Max. This is computed over the lifetime of the process.	seconds")
  hue_auth_pam_auth_time_median                 =create_hue_metric_struct("hue_auth_pam_auth_time_median",	"The time spent waiting for PAM to authenticate a user: 50th Percentile. This is computed over the past hour.	seconds")
  hue_auth_pam_auth_time_min                    =create_hue_metric_struct("hue_auth_pam_auth_time_min",	"The time spent waiting for PAM to authenticate a user: Min. This is computed over the lifetime of the process.	seconds")
  hue_auth_pam_auth_time_rate                   =create_hue_metric_struct("hue_auth_pam_auth_time_rate",	"The time spent waiting for PAM to authenticate a user: Sample Count. This is computed over the lifetime of the process.	message.units.authentications per second")
  hue_auth_pam_auth_time_std_dev                =create_hue_metric_struct("hue_auth_pam_auth_time_std_dev",	"The time spent waiting for PAM to authenticate a user: Standard Deviation. This is computed over the lifetime of the process.	seconds")
  hue_auth_pam_auth_time_sum_rate               =create_hue_metric_struct("hue_auth_pam_auth_time_sum_rate",	"The time spent waiting for PAM to authenticate a user: Sample Sum. This is computed over the lifetime of the process.	message.units.authentications per second")
  hue_auth_saml2_auth_time_15m_rate             =create_hue_metric_struct("hue_auth_saml2_auth_time_15m_rate",	"The time spent waiting for SAML2 to authenticate a user: 15-Minute Rate	message.units.authentications per second")
  hue_auth_saml2_auth_time_1m_rate              =create_hue_metric_struct("hue_auth_saml2_auth_time_1m_rate",	"The time spent waiting for SAML2 to authenticate a user: 1-Minute moving average rate.	message.units.authentications per second")
  hue_auth_saml2_auth_time_5m_rate              =create_hue_metric_struct("hue_auth_saml2_auth_time_5m_rate",	"The time spent waiting for SAML2 to authenticate a user: 5-Minute Rate	message.units.authentications per second")
  hue_auth_saml2_auth_time_75_percentile        =create_hue_metric_struct("hue_auth_saml2_auth_time_75_percentile",	"The time spent waiting for SAML2 to authenticate a user: 75th Percentile. This is computed over the past hour.	seconds")
  hue_auth_saml2_auth_time_95_percentile        =create_hue_metric_struct("hue_auth_saml2_auth_time_95_percentile",	"The time spent waiting for SAML2 to authenticate a user: 95th Percentile. This is computed over the past hour.	seconds")
  hue_auth_saml2_auth_time_999_percentile       =create_hue_metric_struct("hue_auth_saml2_auth_time_999_percentile",	"The time spent waiting for SAML2 to authenticate a user: 999th Percentile. This is computed over the past hour.	seconds")
  hue_auth_saml2_auth_time_99_percentile        =create_hue_metric_struct("hue_auth_saml2_auth_time_99_percentile",	"The time spent waiting for SAML2 to authenticate a user: 99th Percentile. This is computed over the past hour.	seconds")
  hue_auth_saml2_auth_time_avg                  =create_hue_metric_struct("hue_auth_saml2_auth_time_avg",	"The time spent waiting for SAML2 to authenticate a user: Average. This is computed over the lifetime of the process.	seconds")
  hue_auth_saml2_auth_time_max                  =create_hue_metric_struct("hue_auth_saml2_auth_time_max",	"The time spent waiting for SAML2 to authenticate a user: Max. This is computed over the lifetime of the process.	seconds")
  hue_auth_saml2_auth_time_median               =create_hue_metric_struct("hue_auth_saml2_auth_time_median",	"The time spent waiting for SAML2 to authenticate a user: 50th Percentile. This is computed over the past hour.	seconds")
  hue_auth_saml2_auth_time_min                  =create_hue_metric_struct("hue_auth_saml2_auth_time_min",	"The time spent waiting for SAML2 to authenticate a user: Min. This is computed over the lifetime of the process.	seconds")
  hue_auth_saml2_auth_time_rate                 =create_hue_metric_struct("hue_auth_saml2_auth_time_rate",	"The time spent waiting for SAML2 to authenticate a user: Sample Count. This is computed over the lifetime of the process.	message.units.authentications per second")
  hue_auth_saml2_auth_time_std_dev              =create_hue_metric_struct("hue_auth_saml2_auth_time_std_dev",	"The time spent waiting for SAML2 to authenticate a user: Standard Deviation. This is computed over the lifetime of the process.	seconds")
  hue_auth_saml2_auth_time_sum_rate             =create_hue_metric_struct("hue_auth_saml2_auth_time_sum_rate",	"The time spent waiting for SAML2 to authenticate a user: Sample Sum. This is computed over the lifetime of the process.	message.units.authentications per second")
  hue_auth_spnego_auth_time_15m_rate            =create_hue_metric_struct("hue_auth_spnego_auth_time_15m_rate",	"The time spent waiting for SPNEGO to authenticate a user: 15-Minute Rate	message.units.authentications per second")
  hue_auth_spnego_auth_time_1m_rate             =create_hue_metric_struct("hue_auth_spnego_auth_time_1m_rate",	"The time spent waiting for SPNEGO to authenticate a user: 1-Minute moving average rate.	message.units.authentications per second")
  hue_auth_spnego_auth_time_5m_rate             =create_hue_metric_struct("hue_auth_spnego_auth_time_5m_rate",	"The time spent waiting for SPNEGO to authenticate a user: 5-Minute Rate	message.units.authentications per second")
  hue_auth_spnego_auth_time_75_percentile       =create_hue_metric_struct("hue_auth_spnego_auth_time_75_percentile",	"The time spent waiting for SPNEGO to authenticate a user: 75th Percentile. This is computed over the past hour.	seconds")
  hue_auth_spnego_auth_time_95_percentile       =create_hue_metric_struct("hue_auth_spnego_auth_time_95_percentile",	"The time spent waiting for SPNEGO to authenticate a user: 95th Percentile. This is computed over the past hour.	seconds")
  hue_auth_spnego_auth_time_999_percentile      =create_hue_metric_struct("hue_auth_spnego_auth_time_999_percentile",	"The time spent waiting for SPNEGO to authenticate a user: 999th Percentile. This is computed over the past hour.	seconds")
  hue_auth_spnego_auth_time_99_percentile       =create_hue_metric_struct("hue_auth_spnego_auth_time_99_percentile",	"The time spent waiting for SPNEGO to authenticate a user: 99th Percentile. This is computed over the past hour.	seconds")
  hue_auth_spnego_auth_time_avg                 =create_hue_metric_struct("hue_auth_spnego_auth_time_avg",	"The time spent waiting for SPNEGO to authenticate a user: Average. This is computed over the lifetime of the process.	seconds")
  hue_auth_spnego_auth_time_max                 =create_hue_metric_struct("hue_auth_spnego_auth_time_max",	"The time spent waiting for SPNEGO to authenticate a user: Max. This is computed over the lifetime of the process.	seconds")
  hue_auth_spnego_auth_time_median              =create_hue_metric_struct("hue_auth_spnego_auth_time_median",	"The time spent waiting for SPNEGO to authenticate a user: 50th Percentile. This is computed over the past hour.	seconds")
  hue_auth_spnego_auth_time_min                 =create_hue_metric_struct("hue_auth_spnego_auth_time_min",	"The time spent waiting for SPNEGO to authenticate a user: Min. This is computed over the lifetime of the process.	seconds")
  hue_auth_spnego_auth_time_rate                =create_hue_metric_struct("hue_auth_spnego_auth_time_rate",	"The time spent waiting for SPNEGO to authenticate a user: Sample Count. This is computed over the lifetime of the process.	message.units.authentications per second")
  hue_auth_spnego_auth_time_std_dev             =create_hue_metric_struct("hue_auth_spnego_auth_time_std_dev",	"The time spent waiting for SPNEGO to authenticate a user: Standard Deviation. This is computed over the lifetime of the process.	seconds")
  hue_auth_spnego_auth_time_sum_rate            =create_hue_metric_struct("hue_auth_spnego_auth_time_sum_rate",	"The time spent waiting for SPNEGO to authenticate a user: Sample Sum. This is computed over the lifetime of the process.	message.units.authentications per second")
  hue_multiprocessing_processes_daemon          =create_hue_metric_struct("hue_multiprocessing_processes_daemon",	"Number of daemon multiprocessing processes	message.units.processes")
  hue_multiprocessing_processes_total           =create_hue_metric_struct("hue_multiprocessing_processes_total",	"Number of multiprocessing processes	message.units.processes")
  hue_python_gc_generation_0                    =create_hue_metric_struct("hue_python_gc_generation_0",	"Total number of objects in garbage collection generation 0	message.units.objects")
  hue_python_gc_generation_1                    =create_hue_metric_struct("hue_python_gc_generation_1",	"Total number of objects in garbage collection generation 1	message.units.objects")
  hue_python_gc_generation_2                    =create_hue_metric_struct("hue_python_gc_generation_2",	"Total number of objects in garbage collection generation 2	message.units.objects")
  hue_python_gc_objects                         =create_hue_metric_struct("hue_python_gc_objects",	"Total number of objects in the Python process	message.units.objects")
  hue_requests_active                           =create_hue_metric_struct("hue_requests_active",	"Number of currently active requests	requests")
  hue_requests_exceptions_rate                  =create_hue_metric_struct("hue_requests_exceptions_rate",	"Number requests that resulted in an exception	requests per second")
  hue_requests_response_time_15m_rate           =create_hue_metric_struct("hue_requests_response_time_15m_rate",	"Time taken to respond to requests across all Hue endpoints: 15-Minute Rate	requests per second")
  hue_requests_response_time_1m_rate            =create_hue_metric_struct("hue_requests_response_time_1m_rate",	"Time taken to respond to requests across all Hue endpoints: 1-Minute moving average rate.	requests per second")
  hue_requests_response_time_5m_rate            =create_hue_metric_struct("hue_requests_response_time_5m_rate",	"Time taken to respond to requests across all Hue endpoints: 5-Minute Rate	requests per second")
  hue_requests_response_time_75_percentile      =create_hue_metric_struct("hue_requests_response_time_75_percentile",	"Time taken to respond to requests across all Hue endpoints: 75th Percentile. This is computed over the past hour.	seconds")
  hue_requests_response_time_95_percentile      =create_hue_metric_struct("hue_requests_response_time_95_percentile",	"Time taken to respond to requests across all Hue endpoints: 95th Percentile. This is computed over the past hour.	seconds")
  hue_requests_response_time_999_percentile     =create_hue_metric_struct("hue_requests_response_time_999_percentile",	"Time taken to respond to requests across all Hue endpoints: 999th Percentile. This is computed over the past hour.	seconds")
  hue_requests_response_time_99_percentile      =create_hue_metric_struct("hue_requests_response_time_99_percentile",	"Time taken to respond to requests across all Hue endpoints: 99th Percentile. This is computed over the past hour.	seconds")
  hue_requests_response_time_avg                =create_hue_metric_struct("hue_requests_response_time_avg",	"Time taken to respond to requests across all Hue endpoints: Average. This is computed over the lifetime of the process.	seconds")
  hue_requests_response_time_max                =create_hue_metric_struct("hue_requests_response_time_max",	"Time taken to respond to requests across all Hue endpoints: Max. This is computed over the lifetime of the process.	seconds")
  hue_requests_response_time_median             =create_hue_metric_struct("hue_requests_response_time_median",	"Time taken to respond to requests across all Hue endpoints: 50th Percentile. This is computed over the past hour.	seconds")
  hue_requests_response_time_min                =create_hue_metric_struct("hue_requests_response_time_min",	"Time taken to respond to requests across all Hue endpoints: Min. This is computed over the lifetime of the process.	seconds")
  hue_requests_response_time_rate               =create_hue_metric_struct("hue_requests_response_time_rate",	"Time taken to respond to requests across all Hue endpoints: Sample Count. This is computed over the lifetime of the process.	requests per second")
  hue_requests_response_time_std_dev            =create_hue_metric_struct("hue_requests_response_time_std_dev",	"Time taken to respond to requests across all Hue endpoints: Standard Deviation. This is computed over the lifetime of the process.	seconds")
  hue_requests_response_time_sum_rate           =create_hue_metric_struct("hue_requests_response_time_sum_rate",	"Time taken to respond to requests across all Hue endpoints: Sample Sum. This is computed over the lifetime of the process.	requests per second")
  hue_threads_daemon                            =create_hue_metric_struct("hue_threads_daemon",	"The number of daemon threads	threads")
  hue_threads_total                             =create_hue_metric_struct("hue_threads_total",	"The total number of threads	threads")
  hue_users                                     =create_hue_metric_struct("hue_users",	"Total number of user accounts	Users")
  hue_users_active                              =create_hue_metric_struct("hue_users_active",	"Number of users that were active in the last hour	Users")
  hue_mem_rss                                   =create_hue_metric_struct("hue_mem_rss",	"Resident memory used	bytes")
  hue_mem_swap                                  =create_hue_metric_struct("hue_mem_swap",	"Amount of swap memory used by this role's process.	bytes")
  hue_mem_virtual                               =create_hue_metric_struct("hue_mem_virtual",	"Virtual memory used	bytes")
  hue_oom_exits_rate                            =create_hue_metric_struct("hue_oom_exits_rate",	"The number of times the role's backing process was killed due to an OutOfMemory error. This counter is only incremented if the Cloudera Manager 'Kill When Out of Memory' option is enabled.	exits per second")
  hue_read_bytes_rate                           =create_hue_metric_struct("hue_read_bytes_rate",	"The number of bytes read from the device	bytes per second")
  hue_unexpected_exits_rate                     =create_hue_metric_struct("hue_unexpected_exits_rate",	"The number of times the role's backing process exited unexpectedly.	exits per second")
  hue_uptime                                    =create_hue_metric_struct("hue_uptime",	"For a host, the amount of time since the host was booted. For a role, the uptime of the backing process.	seconds")
  hue_web_metrics_collection_duration           =create_hue_metric_struct("hue_web_metrics_collection_duration",	"Web Server Responsiveness	ms")
  hue_write_bytes_rate                          =create_hue_metric_struct("hue_write_bytes_rate",	"The number of bytes written to the device	bytes per second")
)

// Creation of the structure that relates the queries with the descriptors of the Prometheus metrics
var hue_query_variable_relationship = []relation {
  {HUE_ALERTS_RATE,                               *hue_alerts_rate},
  {HUE_CGROUP_CPU_SYSTEM_RATE,                    *hue_cgroup_cpu_system_rate},
  {HUE_CGROUP_CPU_USER_RATE,                      *hue_cgroup_cpu_user_rate},
  {HUE_CGROUP_MEM_PAGE_CACHE,                     *hue_cgroup_mem_page_cache},
  {HUE_CGROUP_MEM_RSS,                            *hue_cgroup_mem_rss},
  {HUE_CGROUP_MEM_SWAP,                           *hue_cgroup_mem_swap},
  {HUE_CGROUP_READ_BYTES_RATE,                    *hue_cgroup_read_bytes_rate},
  {HUE_CGROUP_READ_IOS_RATE,                      *hue_cgroup_read_ios_rate},
  {HUE_CGROUP_WRITE_BYTES_RATE,                   *hue_cgroup_write_bytes_rate},
  {HUE_CGROUP_WRITE_IOS_RATE,                     *hue_cgroup_write_ios_rate},
  {HUE_CPU_SYSTEM_RATE,                           *hue_cpu_system_rate},
  {HUE_CPU_USER_RATE,                             *hue_cpu_user_rate},
  {HUE_EVENTS_CRITICAL_RATE,                      *hue_events_critical_rate},
  {HUE_EVENTS_IMPORTANT_RATE,                     *hue_events_important_rate},
  {HUE_EVENTS_INFORMATIONAL_RATE,                 *hue_events_informational_rate},
  {HUE_FD_MAX,                                    *hue_fd_max},
  {HUE_FD_OPEN,                                   *hue_fd_open},
  {HUE_HEALTH_BAD_RATE,                           *hue_health_bad_rate},
  {HUE_HEALTH_CONCERNING_RATE,                    *hue_health_concerning_rate},
  {HUE_HEALTH_DISABLED_RATE,                      *hue_health_disabled_rate},
  {HUE_HEALTH_GOOD_RATE,                          *hue_health_good_rate},
  {HUE_HEALTH_UNKNOWN_RATE,                       *hue_health_unknown_rate},
  {HUE_AUTH_LDAP_AUTH_TIME_15M_RATE,              *hue_auth_ldap_auth_time_15m_rate},
  {HUE_AUTH_LDAP_AUTH_TIME_1M_RATE,               *hue_auth_ldap_auth_time_1m_rate},
  {HUE_AUTH_LDAP_AUTH_TIME_5M_RATE,               *hue_auth_ldap_auth_time_5m_rate},
  {HUE_AUTH_LDAP_AUTH_TIME_75_PERCENTILE,         *hue_auth_ldap_auth_time_75_percentile},
  {HUE_AUTH_LDAP_AUTH_TIME_95_PERCENTILE,         *hue_auth_ldap_auth_time_95_percentile},
  {HUE_AUTH_LDAP_AUTH_TIME_999_PERCENTILE,        *hue_auth_ldap_auth_time_999_percentile},
  {HUE_AUTH_LDAP_AUTH_TIME_99_PERCENTILE,         *hue_auth_ldap_auth_time_99_percentile},
  {HUE_AUTH_LDAP_AUTH_TIME_AVG,                   *hue_auth_ldap_auth_time_avg},
  {HUE_AUTH_LDAP_AUTH_TIME_MAX,                   *hue_auth_ldap_auth_time_max},
  {HUE_AUTH_LDAP_AUTH_TIME_MEDIAN,                *hue_auth_ldap_auth_time_median},
  {HUE_AUTH_LDAP_AUTH_TIME_MIN,                   *hue_auth_ldap_auth_time_min},
  {HUE_AUTH_LDAP_AUTH_TIME_RATE,                  *hue_auth_ldap_auth_time_rate},
  {HUE_AUTH_LDAP_AUTH_TIME_STD_DEV,               *hue_auth_ldap_auth_time_std_dev},
  {HUE_AUTH_LDAP_AUTH_TIME_SUM_RATE,              *hue_auth_ldap_auth_time_sum_rate},
  {HUE_AUTH_OAUTH_AUTH_TIME_15M_RATE,             *hue_auth_oauth_auth_time_15m_rate},
  {HUE_AUTH_OAUTH_AUTH_TIME_1M_RATE,              *hue_auth_oauth_auth_time_1m_rate},
  {HUE_AUTH_OAUTH_AUTH_TIME_5M_RATE,              *hue_auth_oauth_auth_time_5m_rate},
  {HUE_AUTH_OAUTH_AUTH_TIME_75_PERCENTILE,        *hue_auth_oauth_auth_time_75_percentile},
  {HUE_AUTH_OAUTH_AUTH_TIME_95_PERCENTILE,        *hue_auth_oauth_auth_time_95_percentile},
  {HUE_AUTH_OAUTH_AUTH_TIME_999_PERCENTILE,       *hue_auth_oauth_auth_time_999_percentile},
  {HUE_AUTH_OAUTH_AUTH_TIME_99_PERCENTILE,        *hue_auth_oauth_auth_time_99_percentile},
  {HUE_AUTH_OAUTH_AUTH_TIME_AVG,                  *hue_auth_oauth_auth_time_avg},
  {HUE_AUTH_OAUTH_AUTH_TIME_MAX,                  *hue_auth_oauth_auth_time_max},
  {HUE_AUTH_OAUTH_AUTH_TIME_MEDIAN,               *hue_auth_oauth_auth_time_median},
  {HUE_AUTH_OAUTH_AUTH_TIME_MIN,                  *hue_auth_oauth_auth_time_min},
  {HUE_AUTH_OAUTH_AUTH_TIME_RATE,                 *hue_auth_oauth_auth_time_rate},
  {HUE_AUTH_OAUTH_AUTH_TIME_STD_DEV,              *hue_auth_oauth_auth_time_std_dev},
  {HUE_AUTH_OAUTH_AUTH_TIME_SUM_RATE,             *hue_auth_oauth_auth_time_sum_rate},
  {HUE_AUTH_OPENID_AUTH_TIME_15M_RATE,            *hue_auth_openid_auth_time_15m_rate},
  {HUE_AUTH_OPENID_AUTH_TIME_1M_RATE,             *hue_auth_openid_auth_time_1m_rate},
  {HUE_AUTH_OPENID_AUTH_TIME_5M_RATE,             *hue_auth_openid_auth_time_5m_rate},
  {HUE_AUTH_OPENID_AUTH_TIME_75_PERCENTILE,       *hue_auth_openid_auth_time_75_percentile},
  {HUE_AUTH_OPENID_AUTH_TIME_95_PERCENTILE,       *hue_auth_openid_auth_time_95_percentile},
  {HUE_AUTH_OPENID_AUTH_TIME_999_PERCENTILE,      *hue_auth_openid_auth_time_999_percentile},
  {HUE_AUTH_OPENID_AUTH_TIME_99_PERCENTILE,       *hue_auth_openid_auth_time_99_percentile},
  {HUE_AUTH_OPENID_AUTH_TIME_AVG,                 *hue_auth_openid_auth_time_avg},
  {HUE_AUTH_OPENID_AUTH_TIME_MAX,                 *hue_auth_openid_auth_time_max},
  {HUE_AUTH_OPENID_AUTH_TIME_MEDIAN,              *hue_auth_openid_auth_time_median},
  {HUE_AUTH_OPENID_AUTH_TIME_MIN,                 *hue_auth_openid_auth_time_min},
  {HUE_AUTH_OPENID_AUTH_TIME_RATE,                *hue_auth_openid_auth_time_rate},
  {HUE_AUTH_OPENID_AUTH_TIME_STD_DEV,             *hue_auth_openid_auth_time_std_dev},
  {HUE_AUTH_OPENID_AUTH_TIME_SUM_RATE,            *hue_auth_openid_auth_time_sum_rate},
  {HUE_AUTH_PAM_AUTH_TIME_15M_RATE,               *hue_auth_pam_auth_time_15m_rate},
  {HUE_AUTH_PAM_AUTH_TIME_1M_RATE,                *hue_auth_pam_auth_time_1m_rate},
  {HUE_AUTH_PAM_AUTH_TIME_5M_RATE,                *hue_auth_pam_auth_time_5m_rate},
  {HUE_AUTH_PAM_AUTH_TIME_75_PERCENTILE,          *hue_auth_pam_auth_time_75_percentile},
  {HUE_AUTH_PAM_AUTH_TIME_95_PERCENTILE,          *hue_auth_pam_auth_time_95_percentile},
  {HUE_AUTH_PAM_AUTH_TIME_999_PERCENTILE,         *hue_auth_pam_auth_time_999_percentile},
  {HUE_AUTH_PAM_AUTH_TIME_99_PERCENTILE,          *hue_auth_pam_auth_time_99_percentile},
  {HUE_AUTH_PAM_AUTH_TIME_AVG,                    *hue_auth_pam_auth_time_avg},
  {HUE_AUTH_PAM_AUTH_TIME_MAX,                    *hue_auth_pam_auth_time_max},
  {HUE_AUTH_PAM_AUTH_TIME_MEDIAN,                 *hue_auth_pam_auth_time_median},
  {HUE_AUTH_PAM_AUTH_TIME_MIN,                    *hue_auth_pam_auth_time_min},
  {HUE_AUTH_PAM_AUTH_TIME_RATE,                   *hue_auth_pam_auth_time_rate},
  {HUE_AUTH_PAM_AUTH_TIME_STD_DEV,                *hue_auth_pam_auth_time_std_dev},
  {HUE_AUTH_PAM_AUTH_TIME_SUM_RATE,               *hue_auth_pam_auth_time_sum_rate},
  {HUE_AUTH_SAML2_AUTH_TIME_15M_RATE,             *hue_auth_saml2_auth_time_15m_rate},
  {HUE_AUTH_SAML2_AUTH_TIME_1M_RATE,              *hue_auth_saml2_auth_time_1m_rate},
  {HUE_AUTH_SAML2_AUTH_TIME_5M_RATE,              *hue_auth_saml2_auth_time_5m_rate},
  {HUE_AUTH_SAML2_AUTH_TIME_75_PERCENTILE,        *hue_auth_saml2_auth_time_75_percentile},
  {HUE_AUTH_SAML2_AUTH_TIME_95_PERCENTILE,        *hue_auth_saml2_auth_time_95_percentile},
  {HUE_AUTH_SAML2_AUTH_TIME_999_PERCENTILE,       *hue_auth_saml2_auth_time_999_percentile},
  {HUE_AUTH_SAML2_AUTH_TIME_99_PERCENTILE,        *hue_auth_saml2_auth_time_99_percentile},
  {HUE_AUTH_SAML2_AUTH_TIME_AVG,                  *hue_auth_saml2_auth_time_avg},
  {HUE_AUTH_SAML2_AUTH_TIME_MAX,                  *hue_auth_saml2_auth_time_max},
  {HUE_AUTH_SAML2_AUTH_TIME_MEDIAN,               *hue_auth_saml2_auth_time_median},
  {HUE_AUTH_SAML2_AUTH_TIME_MIN,                  *hue_auth_saml2_auth_time_min},
  {HUE_AUTH_SAML2_AUTH_TIME_RATE,                 *hue_auth_saml2_auth_time_rate},
  {HUE_AUTH_SAML2_AUTH_TIME_STD_DEV,              *hue_auth_saml2_auth_time_std_dev},
  {HUE_AUTH_SAML2_AUTH_TIME_SUM_RATE,             *hue_auth_saml2_auth_time_sum_rate},
  {HUE_AUTH_SPNEGO_AUTH_TIME_15M_RATE,            *hue_auth_spnego_auth_time_15m_rate},
  {HUE_AUTH_SPNEGO_AUTH_TIME_1M_RATE,             *hue_auth_spnego_auth_time_1m_rate},
  {HUE_AUTH_SPNEGO_AUTH_TIME_5M_RATE,             *hue_auth_spnego_auth_time_5m_rate},
  {HUE_AUTH_SPNEGO_AUTH_TIME_75_PERCENTILE,       *hue_auth_spnego_auth_time_75_percentile},
  {HUE_AUTH_SPNEGO_AUTH_TIME_95_PERCENTILE,       *hue_auth_spnego_auth_time_95_percentile},
  {HUE_AUTH_SPNEGO_AUTH_TIME_999_PERCENTILE,      *hue_auth_spnego_auth_time_999_percentile},
  {HUE_AUTH_SPNEGO_AUTH_TIME_99_PERCENTILE,       *hue_auth_spnego_auth_time_99_percentile},
  {HUE_AUTH_SPNEGO_AUTH_TIME_AVG,                 *hue_auth_spnego_auth_time_avg},
  {HUE_AUTH_SPNEGO_AUTH_TIME_MAX,                 *hue_auth_spnego_auth_time_max},
  {HUE_AUTH_SPNEGO_AUTH_TIME_MEDIAN,              *hue_auth_spnego_auth_time_median},
  {HUE_AUTH_SPNEGO_AUTH_TIME_MIN,                 *hue_auth_spnego_auth_time_min},
  {HUE_AUTH_SPNEGO_AUTH_TIME_RATE,                *hue_auth_spnego_auth_time_rate},
  {HUE_AUTH_SPNEGO_AUTH_TIME_STD_DEV,             *hue_auth_spnego_auth_time_std_dev},
  {HUE_AUTH_SPNEGO_AUTH_TIME_SUM_RATE,            *hue_auth_spnego_auth_time_sum_rate},
  {HUE_MULTIPROCESSING_PROCESSES_DAEMON,          *hue_multiprocessing_processes_daemon},
  {HUE_MULTIPROCESSING_PROCESSES_TOTAL,           *hue_multiprocessing_processes_total},
  {HUE_PYTHON_GC_GENERATION_0,                    *hue_python_gc_generation_0},
  {HUE_PYTHON_GC_GENERATION_1,                    *hue_python_gc_generation_1},
  {HUE_PYTHON_GC_GENERATION_2,                    *hue_python_gc_generation_2},
  {HUE_PYTHON_GC_OBJECTS,                         *hue_python_gc_objects},
  {HUE_REQUESTS_ACTIVE,                           *hue_requests_active},
  {HUE_REQUESTS_EXCEPTIONS_RATE,                  *hue_requests_exceptions_rate},
  {HUE_REQUESTS_RESPONSE_TIME_15M_RATE,           *hue_requests_response_time_15m_rate},
  {HUE_REQUESTS_RESPONSE_TIME_1M_RATE,            *hue_requests_response_time_1m_rate},
  {HUE_REQUESTS_RESPONSE_TIME_5M_RATE,            *hue_requests_response_time_5m_rate},
  {HUE_REQUESTS_RESPONSE_TIME_75_PERCENTILE,      *hue_requests_response_time_75_percentile},
  {HUE_REQUESTS_RESPONSE_TIME_95_PERCENTILE,      *hue_requests_response_time_95_percentile},
  {HUE_REQUESTS_RESPONSE_TIME_999_PERCENTILE,     *hue_requests_response_time_999_percentile},
  {HUE_REQUESTS_RESPONSE_TIME_99_PERCENTILE,      *hue_requests_response_time_99_percentile},
  {HUE_REQUESTS_RESPONSE_TIME_AVG,                *hue_requests_response_time_avg},
  {HUE_REQUESTS_RESPONSE_TIME_MAX,                *hue_requests_response_time_max},
  {HUE_REQUESTS_RESPONSE_TIME_MEDIAN,             *hue_requests_response_time_median},
  {HUE_REQUESTS_RESPONSE_TIME_MIN,                *hue_requests_response_time_min},
  {HUE_REQUESTS_RESPONSE_TIME_RATE,               *hue_requests_response_time_rate},
  {HUE_REQUESTS_RESPONSE_TIME_STD_DEV,            *hue_requests_response_time_std_dev},
  {HUE_REQUESTS_RESPONSE_TIME_SUM_RATE,           *hue_requests_response_time_sum_rate},
  {HUE_THREADS_DAEMON,                            *hue_threads_daemon},
  {HUE_THREADS_TOTAL,                             *hue_threads_total},
  {HUE_USERS,                                     *hue_users},
  {HUE_USERS_ACTIVE,                              *hue_users_active},
  {HUE_MEM_RSS,                                   *hue_mem_rss},
  {HUE_MEM_SWAP,                                  *hue_mem_swap},
  {HUE_MEM_VIRTUAL,                               *hue_mem_virtual},
  {HUE_OOM_EXITS_RATE,                            *hue_oom_exits_rate},
  {HUE_READ_BYTES_RATE,                           *hue_read_bytes_rate},
  {HUE_UNEXPECTED_EXITS_RATE,                     *hue_unexpected_exits_rate},
  {HUE_UPTIME,                                    *hue_uptime},
  {HUE_WEB_METRICS_COLLECTION_DURATION,           *hue_web_metrics_collection_duration},
  {HUE_WRITE_BYTES_RATE,                          *hue_write_bytes_rate},
}




/* ======================================================================
 * Functions
 * ====================================================================== */
// Create and returns a prometheus descriptor for a hue metric. 
// The "metric_name" parameter its mandatory
// If the "description" parameter is empty, the function assings it with the
// value of the name of the metric in uppercase and separated by spaces
func create_hue_metric_struct(metric_name string, description string) *prometheus.Desc {
  // Correct "description" parameter if is empty
  if len(description) == 0 {
    description = strings.Replace(strings.ToUpper(metric_name), "_", " ", 0)
  }

  // return prometheus descriptor
  return prometheus.NewDesc(
    prometheus.BuildFQName(namespace, HUE_SCRAPER_NAME, metric_name),
    description,
    []string{"cluster", "entityName","hostname"},
    nil,
  )
}


// Generic function to extract de metadata associated with the query value
// Only for HUE metric type
func create_hue_metric (ctx context.Context, config Collector_connection_data, query string, metric_struct prometheus.Desc, ch chan<- prometheus.Metric) bool {
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
// ScrapeHUE struct
type ScrapeHUE struct{}

// Name of the Scraper. Should be unique.
func (ScrapeHUE) Name() string {
  return HUE_SCRAPER_NAME
}

// Help describes the role of the Scraper.
func (ScrapeHUE) Help() string {
  return "HUE Metrics"
}

// Version.
func (ScrapeHUE) Version() float64 {
  return 1.0
}

// Scrape generic function. Override for host module.
func (ScrapeHUE) Scrape (ctx context.Context, config *Collector_connection_data, ch chan<- prometheus.Metric) error {
  log.Debug_msg("Executing HUE Metrics Scraper")

  // Queries counters
  success_queries := 0
  error_queries := 0

  // Execute the generic funtion for creation of metrics with the pairs (QUERY, PROM:DESCRIPTOR)
  for i:=0 ; i < len(hue_query_variable_relationship) ; i++ {
    if create_hue_metric(ctx, *config, hue_query_variable_relationship[i].Query, hue_query_variable_relationship[i].Metric_struct, ch) {
      success_queries += 1
    } else {
      error_queries += 1
    }
  }
  log.Info_msg("In the HUE Module has been executed %d queries. %d success and %d with errors", success_queries + error_queries, success_queries, error_queries)
  return nil
}

var _ Scraper = ScrapeHUE{}
