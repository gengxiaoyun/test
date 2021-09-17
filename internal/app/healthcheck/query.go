package healthcheck

const (
	// application mysql
	applicationMySQLVariables = `
		select variable_name, variable_value
		from %s.global_variables
		where variable_name in (%s);
    `
	applicationMySQLTableSize = `
    	select table_schema,
			   table_name,
			   table_rows,
			   truncate((data_length + index_length) / 1024 / 1024 / 1024, 2) as size
		from information_schema.tables
		where table_type = 'BASE TABLE'
		  and table_rows > ?
		order by table_rows desc;
    `
	// prometheus
	PrometheusCPUUsageV1 = `
		clamp_max(sum by () ((avg by (mode) (
		(clamp_max(rate(node_cpu{instance=~"%s",mode!="idle",mode!="iowait"}[5m]),1)) or
		(clamp_max(irate(node_cpu{instance=~"%s",mode!="idle",mode!="iowait"}[5m]),1)) )) *100 or
		sum by () (
		avg_over_time(node_cpu_average{instance=~"%s",mode!="total",mode!="idle"}[5m]) or
		avg_over_time(node_cpu_average{instance=~"%s",mode!="total",mode!="idle"}[5m])) unless
		(avg_over_time(node_cpu_average{instance=~"%s",mode="total",job="rds-basic"}[5m]) or
		avg_over_time(node_cpu_average{instance=~"%s",mode="total",job="rds-basic"}[5m]))
		),100)
    `
	PrometheusCPUUsageV2 = `
		clamp_max(sum by () ((avg by (mode) ( 
		(clamp_max(rate(node_cpu_seconds_total{node_name=~"%s",mode!="idle",mode!="iowait"}[5m]),1)) or 
		(clamp_max(irate(node_cpu_seconds_total{node_name=~"%s",mode!="idle",mode!="iowait"}[5m]),1)) )) *100 or 
		sum by () (
		avg_over_time(node_cpu_average{node_name=~"%s",mode!="total",mode!="idle"}[5m]) or 
		avg_over_time(node_cpu_average{node_name=~"%s",mode!="total",mode!="idle"}[5m])) unless
		(avg_over_time(node_cpu_average{node_name=~"%s",mode="total",job="rds-basic"}[5m]) or 
		avg_over_time(node_cpu_average{node_name=~"%s",mode="total",job="rds-basic"}[5m]))
		),100)
    `
	PrometheusFileSystemV1 = `
		node_filesystem_files{instance=~"%s",fstype!~"rootfs|selinuxfs|autofs|rpc_pipefs|tmpfs"}
    `
	PrometheusFileSystemV2 = `
		node_filesystem_files{node_name=~"%s",fstype!~"rootfs|selinuxfs|autofs|rpc_pipefs|tmpfs"}
    `
	PrometheusIOUtilV1 = `
		avg by (instance) (rate(node_disk_io_time_ms{instance=~"%s"}[20s])/1000 or
		irate(node_disk_io_time_ms{instance=~"%s"}[5m])/1000)
    `
	PrometheusIOUtilV2 = `
		avg by (node_name) (rate(node_disk_io_time_seconds_total{node_name=~"%s"}[20s]) or
		irate(node_disk_io_time_seconds_total{node_name=~"%s"}[5m]) or
		(max_over_time(rdsosmetrics_diskIO_util{node_name=~"%s"}[20s]) or
		max_over_time(rdsosmetrics_diskIO_util{node_name=~"%s"}[5m]))/100)
    `
	PrometheusDiskCapacityV1 = `
		1 - node_filesystem_free{instance=~"%s", mountpoint=~"(%s)", fstype!~"rootfs|selinuxfs|autofs|rpc_pipefs|tmpfs"} /
		node_filesystem_size{instance=~"%s", mountpoint=~"(%s)", fstype!~"rootfs|selinuxfs|autofs|rpc_pipefs|tmpfs"}
    `
	PrometheusDiskCapacityV2 = `
		avg by (node_name, mountpoint) (1 - (max_over_time(node_filesystem_free_bytes{node_name=~"%s", mountpoint=~"(%s)", fstype!~"rootfs|selinuxfs|autofs|rpc_pipefs|tmpfs"}[20s]) or
		max_over_time(node_filesystem_free_bytes{node_name=~"%s", mountpoint=~"(%s)", fstype!~"rootfs|selinuxfs|autofs|rpc_pipefs|tmpfs"}[5m])) /
		(max_over_time(node_filesystem_size_bytes{node_name=~"%s", mountpoint=~"(%s)", fstype!~"rootfs|selinuxfs|autofs|rpc_pipefs|tmpfs"}[20s]) or
		max_over_time(node_filesystem_size_bytes{node_name=~"%s", mountpoint=~"(%s)", fstype!~"rootfs|selinuxfs|autofs|rpc_pipefs|tmpfs"}[5m])))
   `
	PrometheusConnectionUsageV1 = `
		avg by (instance) (max(max_over_time(mysql_global_status_threads_connected{instance=~"%s"}[20s]) or
		max_over_time(mysql_global_status_threads_connected{instance=~"%s"}[5m]))) /
		avg by (instance) (max(max_over_time(mysql_global_variables_max_connections{instance=~"%s"}[20s]) or
		max_over_time(mysql_global_variables_max_connections{instance=~"%s"}[5m])))
    `
	PrometheusConnectionUsageV2 = `
		avg by (service_name) (max(max_over_time(mysql_global_status_threads_connected{service_name=~"%s"}[20s]) or
		max_over_time(mysql_global_status_threads_connected{service_name=~"%s"}[5m]))) /
		avg by (service_name) (max_over_time(mysql_global_variables_max_connections{service_name=~"%s"}[20s]) or
		max_over_time(mysql_global_variables_max_connections{service_name=~"%s"}[5m]))
    `
	PrometheusAverageActiveSessionPercentsV1 = `
		avg by (instance) (avg_over_time(mysql_global_status_threads_running{instance=~"%s"}[20s]) or
		avg_over_time(mysql_global_status_threads_running{instance=~"%s"}[5m]))/
		avg by (instance) (max_over_time(mysql_global_status_threads_connected{instance=~"%s"}[20s]) or
		max_over_time(mysql_global_status_threads_connected{instance=~"%s"}[5m]))
    `
	PrometheusAverageActiveSessionPercentsV2 = `
		avg by (service_name) (avg_over_time(mysql_global_status_threads_running{service_name=~"%s"}[20s]) or
		avg_over_time(mysql_global_status_threads_running{service_name=~"%s"}[5m]))/
		avg by (service_name) (max_over_time(mysql_global_status_threads_connected{service_name=~"%s"}[20s]) or
		max_over_time(mysql_global_status_threads_connected{service_name=~"%s"}[5m]))
    `
	PrometheusCacheMissRatioV1 = `
		avg by (instance) ((rate(mysql_global_status_innodb_buffer_pool_reads{instance=~"%s"}[5m]) or
		irate(mysql_global_status_innodb_buffer_pool_reads{instance=~"%s"}[5m])) /
		(rate(mysql_global_status_innodb_buffer_pool_read_requests{instance=~"%s"}[5m]) or
		irate(mysql_global_status_innodb_buffer_pool_read_requests{instance=~"%s"}[5m])))
    `
	PrometheusCacheMissRatioV2 = `
		avg by (service_name) ((rate(mysql_global_status_innodb_buffer_pool_reads{service_name=~"%s"}[5m]) or
		irate(mysql_global_status_innodb_buffer_pool_reads{service_name=~"%s"}[5m])) /
		(rate(mysql_global_status_innodb_buffer_pool_read_requests{service_name=~"%s"}[5m]) or
		irate(mysql_global_status_innodb_buffer_pool_read_requests{service_name=~"%s"}[5m])))
    `
	// query
	MonitorMySQLQuery = `
		select qc.checksum as sql_id,
			   qc.fingerprint,
			   qe.query    as example,
			   qe.db       as db_name,
			   m.exec_count,
			   m.total_exec_time,
			   m.avg_exec_time,
			   m.rows_examined_max
		from (
				 select qcm.query_class_id,
						sum(qcm.query_count)                                        as exec_count,
						truncate(sum(qcm.query_time_sum), 2)                        as total_exec_time,
						truncate(sum(qcm.query_time_sum) / sum(qcm.query_count), 2) as avg_exec_time,
						qcm.rows_examined_max
				 from query_class_metrics qcm
						  inner join instances i on qcm.instance_id = i.instance_id
				 where i.name = ?
				   and qcm.start_ts >= ?
				   and qcm.start_ts < ?
				   and qcm.rows_examined_max >= ?
				 group by query_class_id
				 order by rows_examined_max desc
				 limit ?) m
				 inner join query_examples qe on m.query_class_id = qe.query_class_id
				 inner join query_classes qc on m.query_class_id = qc.query_class_id;
    `
	MonitorClickhouseQuery = `
		select sm.sql_id,
			   m.fingerprint,
			   m.example,
			   m.db_name,
			   sm.exec_count,
			   sm.total_exec_time,
			   sm.avg_exec_time,
			   sm.rows_examined_max
		
		from (
				 select queryid                                               as sql_id,
						sum(num_queries)                                      as exec_count,
						truncate(sum(m_query_time_sum), 2)                    as total_exec_time,
						truncate(sum(m_query_time_sum) / sum(num_queries), 2) as avg_exec_time,
						max(m_rows_examined_max)                              as rows_examined_max
				 from metrics
				 where service_type = 'mysql'
				   and service_name = ?
				   and period_start >= ?
				   and period_start < ?
				   and m_rows_examined_max >= ?
				 group by queryid
				 order by rows_examined_max desc
				 limit ?) sm
				 left join (select queryid          as sql_id,
								   max(fingerprint) as fingerprint,
								   max(example)     as example,
								   max(database)    as db_name
							from metrics
							where service_type = 'mysql'
							  and service_name = ?
							  and period_start >= ?
							  and period_start < ?
							  and m_rows_examined_max >= ?
							group by queryid) m
						   on sm.sql_id = m.sql_id;
    `
)
