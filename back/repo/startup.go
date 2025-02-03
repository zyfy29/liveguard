package repo

func init() {
	for _, ddlProc := range ddls {
		ddlProc()
	}
}
