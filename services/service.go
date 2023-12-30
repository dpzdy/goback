package services

var (
	FtempServiceMgr *ftempServiceMgr
)

func ServiceInit() {
	FtempServiceMgr = &ftempServiceMgr{}
}
