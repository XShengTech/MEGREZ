package instanceController

const (
	apiPrefix = "/api/v1"
)

const (
	instancePrefix   = "/replicaSet"
	instanceRestart  = "/restart"
	instanceStart    = "/continue"
	instanceStop     = "/stop"
	instancePause    = "/pause"
	instanceHistory  = "/history"
	instanceRollback = "/rollback"
	instanceImage    = "/commit"
	instanceExecute  = "/execute"
)

const (
	volumePrefix  = "/volumes"
	volumeSize    = "/size"
	volumeHistory = "/history"
)

type Action int

const (
	ActionStart   Action = 1
	ActionPause   Action = 2
	ActionStop    Action = 3
	ActionRestart Action = 4
)
