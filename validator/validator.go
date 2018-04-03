package validator

type Validator interface{
	BuildTree(content []Content) error
	ReBuildTree() error
	VerifyContent(content Content) (bool, error)
}

