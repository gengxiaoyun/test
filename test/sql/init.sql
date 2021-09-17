truncate table t_meta_monitor_system_info;
truncate table t_meta_mysql_cluster_info;
truncate table t_meta_mysql_server_info;

insert into t_meta_monitor_system_info(id, system_name, system_type, host_ip, port_num, port_num_slow, base_url, env_id)
values(1, 'pmm-2', 2, '192.168.10.219', 80, 8123, '/prometheus', 6), (2, 'pmm-1', 1, '192.168.10.220', 9090, 33061, '/prometheus', 6);
insert into t_meta_mysql_cluster_info(id, cluster_name, monitor_system_id, env_id)
values(1, 'mysql-cluster-01', 1, 6), (2, 'mysql-cluster-02', 2, 6);
insert into t_meta_mysql_server_info(id, cluster_id, server_name, service_name, host_ip, port_num, deployment_type, version)
values(1, 1, '192-168-10-219-3306', '192-168-10-219:3306', '192.168.10.219', 3306, 2, '5.7'),
      (2, 1, '192-168-10-219-3307', '192-168-10-219:3307', '192.168.10.219', 3307, 2, '5.7'),
      (3, 2, '192-168-10-220-3306', '192-168-10-220:3306', '192.168.10.220', 3306, 2, '5.7'),
       (4, 2, '192-168-10-220-3307', '192-168-10-220:3307', '192.168.10.220', 3307, 2, '5.7');
