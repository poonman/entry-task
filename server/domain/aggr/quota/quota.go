package quota

type Quota struct {
	Id int
	Username string
	ReadQuota int
	WriteQuota int
}
