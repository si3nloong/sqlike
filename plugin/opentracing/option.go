package opentracing

// WithAllTraceOptions :
func WithAllTraceOptions() TraceOption {
	return func(opt *TraceOptions) {
		opt.Ping = true
		opt.BeginTx = true
		opt.TxCommit = true
		opt.TxRollback = true
		opt.Prepare = true
		opt.Query = true
		opt.Exec = true
		opt.RowsAffected = true
		opt.RowsClose = true
		opt.RowsNext = true
		opt.LastInsertID = true
		opt.Args = true
	}
}

// WithComponent :
func WithComponent(comp string) TraceOption {
	return func(opt *TraceOptions) {
		opt.Component = comp
	}
}

// WithDBInstance :
func WithDBInstance(instance string) TraceOption {
	return func(opt *TraceOptions) {
		opt.DBInstance = instance
	}
}

// WithDBType :
func WithDBType(dbType string) TraceOption {
	return func(opt *TraceOptions) {
		opt.DBType = dbType
	}
}

// WithDBUser :
func WithDBUser(user string) TraceOption {
	return func(opt *TraceOptions) {
		opt.DBUser = user
	}
}

// WithPrepare :
func WithPrepare(flag bool) TraceOption {
	return func(opt *TraceOptions) {
		opt.Prepare = flag
	}
}

// WithExec :
func WithExec(flag bool) TraceOption {
	return func(opt *TraceOptions) {
		opt.Exec = flag
	}
}

// WithQuery :
func WithQuery(flag bool) TraceOption {
	return func(opt *TraceOptions) {
		opt.Query = flag
	}
}

// WithArgs :
func WithArgs(flag bool) TraceOption {
	return func(opt *TraceOptions) {
		opt.Args = flag
	}
}

// WithRowsClose :
func WithRowsClose(flag bool) TraceOption {
	return func(opt *TraceOptions) {
		opt.RowsClose = flag
	}
}
