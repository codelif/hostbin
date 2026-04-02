package hosts

func DocumentURL(baseDomain, slug string) string {
	return "https://" + slug + "." + baseDomain + "/"
}
