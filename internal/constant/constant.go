package constant

type ContextKey string

const (
	LoggerKey ContextKey = "logger"
	UserIDKey ContextKey = "user_id"
	AdminKey  ContextKey = "is_admin"
)
