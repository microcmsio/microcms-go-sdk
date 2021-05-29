package microcms

const (
	BaseDomain = "microcms.io"
	ApiVersion = "v1"
)

func Client() string {
	return BaseDomain + ApiVersion
}
