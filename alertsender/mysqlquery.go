package alertsender

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
	"xxl_job_alert/alertinit"
	"xxl_job_alert/alertmodel"

	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql"
)

var Pool *MySQLConnectionPools

type MySQLConnectionPools struct {
	pools map[string]*sql.DB
	mu    sync.Mutex
}

func NewMySQLConnectionPools() *MySQLConnectionPools {
	return &MySQLConnectionPools{
		pools: make(map[string]*sql.DB),
		mu:    sync.Mutex{},
	}
}
func InitializeConnectionPools() {
	var config = alertinit.Conf
	connectionPools := NewMySQLConnectionPools()
	for _, instance := range config.Instances {
		if instance.Address != "" {
			if instance.DB == "" {
				instance.DB = "xxl_job"
			}
			db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", instance.Username, instance.Password, instance.Address, instance.DB))
			if err != nil {
				log.Error().Msgf("failed to open MySQL connection for instance '%s': %v", fmt.Sprintf("%s-%s", instance.Address, instance.DB), err)
			}
			// 设置连接池最大连接数为 5
			db.SetMaxOpenConns(5)
			connectionPools.pools[instance.ID] = db
		}
	}
	Pool = connectionPools
}

// 获取连接池对象
func (cp *MySQLConnectionPools) GetMySQLConnectionPool(Id string) (*sql.DB, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	pool, ok := cp.pools[Id]
	if ok {
		return pool, nil
	}

	return nil, fmt.Errorf("connection pool '%s' not found", Id)
}

func QueryList() []alertmodel.XXLAlertResult {
	resultList := []alertmodel.XXLAlertResult{}
	for _, instance := range alertinit.Conf.Instances {
		if instance.Address != "" {
			conn, err := Pool.GetMySQLConnectionPool(instance.ID)
			if err != nil {
				log.Error().Msgf("Failed to get connection:%s", err)
				return resultList
			}
			tm := time.Now().Add(-(time.Duration(alertinit.Conf.Interval)) * time.Second).Format("2006-01-02 15:04:05")
			sql := fmt.Sprintf(`select 
			log.id,
			log.executor_address,
			log.trigger_time,
			log.handle_time,
			log.executor_handler,
			log.executor_param,
			info.job_desc,
			log.trigger_msg 
		from xxl_job_log log 
			LEFT JOIN xxl_job_info info on log.job_id = info.id 
		where 
			log.alarm_status != 0 
			and log.trigger_time > '%s' 
		ORDER BY log.trigger_time DESC 
			`, tm)
			rows, err := conn.Query(sql)
			if err != nil {
				log.Error().Msgf("Failed to execute query:%s", err)
			} else {
				defer rows.Close()
				for rows.Next() {
					var result alertmodel.XXLAlertResult
					err := rows.Scan(&result.Id,
						&result.Executor_address,
						&result.Trigger_time,
						&result.Handle_time,
						&result.Executor_handler,
						&result.Executor_param,
						&result.Job_desc,
						&result.Trigger_msg)
					if err != nil {
						log.Error().Msgf("Failed to scan row:%s", err)
					} else {
						result.Instance = &instance.ID
						resultList = append(resultList, result)
					}
				}
			}
		}
	}
	return resultList
}
