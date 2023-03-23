// Code generated by "stringer -type JobStatus -linecomment -output task_string.go"; DO NOT EDIT.

package task

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[JOB_STATUS_FAILED - -1]
	_ = x[JOB_STATUS_TIMEOUT - -2]
	_ = x[JOB_STATUS_CREATED-0]
	_ = x[JOB_STATUS_RUNNING-1]
	_ = x[JOB_STATUS_SUCCESS-2]
}

const _JobStatus_name = "超时失败已创建运行中成功"

var _JobStatus_index = [...]uint8{0, 6, 12, 21, 30, 36}

func (i JobStatus) String() string {
	i -= -2
	if i < 0 || i >= JobStatus(len(_JobStatus_index)-1) {
		return "JobStatus(" + strconv.FormatInt(int64(i+-2), 10) + ")"
	}
	return _JobStatus_name[_JobStatus_index[i]:_JobStatus_index[i+1]]
}
