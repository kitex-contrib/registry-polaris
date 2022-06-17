package polaris

type Options struct {
	DstNamespace string            `json:"dst_namespace"`
	DstService   string            `json:"dst_service"`
	DstMetadata  map[string]string `json:"dst_metadata"`
	SrcNamespace string            `json:"src_namespace"`
	SrcService   string            `json:"src_service"`
	SrcMetadata  map[string]string `json:"src_metadata"`
}
