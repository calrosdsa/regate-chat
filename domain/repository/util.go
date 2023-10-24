package repository


type UtilUseCase interface {
	PaginationValues(page int16) int16
	GetNextPage(results int8, pageSize int8, page int16) (nextPage int16)
	LogError(method string, file string, err string)
	LogInfo(method string, file string, err string)
	CustomLog(method string, file string, err string, payload map[string]interface{})
	LogFatal(method string, file string, err string, payload map[string]interface{})
}