package assert

// A simplyfied adaptation of the hamcrest library matcher.
type Matcher interface {
	// Perform the assertion logic against the provided item.
	Matches(item interface{}) bool
	// Generate a description of why the matcher has not accepted
	// the item. The description will be part of a larger description
	// of why a matching failed, so it should be concise.
	// This method assumes that [Matches(interface{})] is false,
	// but will not check this.
	DescribeMismatch(item interface{}) string
	// Generate a descripton of what this matcher expects.
	Describe() string
}

type notMatcher struct {
	proxy Matcher
}

func (n *notMatcher) Matches(item interface{}) bool {
	return false == n.proxy.Matches(item)
}

func (n *notMatcher) DescribeMismatch(item interface{}) string {
	return "[not] " + n.proxy.DescribeMismatch(item)
}

func (n *notMatcher) Describe() string {
	return "[not] " + n.proxy.Describe()
}

func Not(matcher Matcher) Matcher {
	return &notMatcher{matcher}
}
