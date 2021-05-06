package filters

import t "github.com/containrrr/watchtower/pkg/types"

// WatchtowerContainersFilter filters only watchtower containers
func WatchtowerContainersFilter(c t.FilterableContainer) bool { return c.IsWatchtower() }

// NoFilter will not filter out any containers
func NoFilter(t.FilterableContainer) bool { return true }

// FilterByNames returns all containers that match the specified name
func FilterByNames(names []string, baseFilter t.Filter) t.Filter {
	if len(names) == 0 {
		return baseFilter
	}

	return func(c t.FilterableContainer) bool {
		for _, name := range names {
			if (name == c.Name()) || (name == c.Name()[1:]) {
				return baseFilter(c)
			}
		}
		return false
	}
}

// FilterByEnableLabel returns all containers that have the enabled label set
func FilterByEnableLabel(baseFilter t.Filter) t.Filter {
	return func(c t.FilterableContainer) bool {
		// If label filtering is enabled, containers should only be considered
		// if the label is specifically set.
		_, ok := c.Enabled()
		if !ok {
			return false
		}

		return baseFilter(c)
	}
}

// FilterByDisabledLabel returns all containers that have the enabled label set to disable
func FilterByDisabledLabel(baseFilter t.Filter) t.Filter {
	return func(c t.FilterableContainer) bool {
		enabledLabel, ok := c.Enabled()
		if ok && !enabledLabel {
			// If the label has been set and it demands a disable
			return false
		}

		return baseFilter(c)
	}
}

// FilterByScope returns all containers that belongs to a specific scope
func FilterByScope(scope string, baseFilter t.Filter) t.Filter {
	if scope == "" {
		return baseFilter
	}

	return func(c t.FilterableContainer) bool {
		containerScope, ok := c.Scope()
		if ok && containerScope == scope {
			return baseFilter(c)
		}

		return false
	}
}

// FilterByPoll returns all containers that has the poll time label set
func FilterByPollTime(polls []int, baseFilter t.Filter) t.Filter {

	if len(polls) == 0 {
		return baseFilter
	}

	return func(c t.FilterableContainer) bool {
		for _, pollTime := range polls {
			if (pollTime == c.PollTime()) {
				return baseFilter(c)
			}
		}
		return false
	}

}

// BuildFilter creates the needed filter of containers
func BuildFilter(names []string, enableLabel bool, scope string, pollTime []int) t.Filter {
	filter := NoFilter
	filter = FilterByNames(names, filter)
	filter = FilterByPollTime(pollTime, filter)
	
	if enableLabel {
		// If label filtering is enabled, containers should only be considered
		// if the label is specifically set.
		filter = FilterByEnableLabel(filter)
	}
	if scope != "" {
		// If a scope has been defined, containers should only be considered
		// if the scope is specifically set.
		filter = FilterByScope(scope, filter)
	}

	filter = FilterByDisabledLabel(filter)
	return filter
}
