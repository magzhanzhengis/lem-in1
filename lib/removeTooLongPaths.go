package lib

func RemoveTooLongPaths(paths [][]string) [][]string {
	if len(paths) >= 10 {
		paths = append(paths[:5], paths[6:]...)
	}

	var shortestPaths [][]string
	for _, path := range paths {
		if len(paths) >= 10 {
			if isPathWithSameStartExists(shortestPaths, path) {
				continue
			} else {
				isHasDuplicate := false
				for _, p := range path {
					if isPathsHaveThisRoom(shortestPaths, p) {
						isHasDuplicate = true
						break
					}
				}

				if !isHasDuplicate {
					shortestPaths = append(shortestPaths, path)
					continue
				}
			}
		} else {
			if len(path)-2 > len(paths[0]) {
				continue
			} else {
				shortestPaths = append(shortestPaths, path)
			}
		}
	}

	return shortestPaths
}

func isPathWithSameStartExists(paths [][]string, path []string) bool {
	for _, p := range paths {
		if p[1] == path[1] {
			return true
		}
	}
	return false
}

func isPathsHaveThisRoom(paths [][]string, room string) bool {
	for _, path := range paths {
		for _, p := range path {
			if p == room && p != "start" && p != "end" {
				return true
			}
		}
	}
	return false
}
