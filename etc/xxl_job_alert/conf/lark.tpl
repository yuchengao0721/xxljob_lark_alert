**XXL_JOB出现{{len .}}个失败调度，详情如下：**
———————————————————{{range .}}
**调  度 器**: {{.Instance}}
{{- if .Executor_address }}
**调度地址**: {{.Executor_address}}
{{- end }}
**失败任务**: {{.Job_desc}} | {{.Executor_handler}} 
**触发时间**: {{.Trigger_time}} 
{{- if .Executor_param }}
**任务参数**: {{.Executor_param}}
{{- end }}
**失败详情**: {{.Trigger_msg}}
———————————————————{{end}}
