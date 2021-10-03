package generator

import (
	"strings"
)

func GetMockName(mock string) *mockNameDecl {
	separatorIndex := strings.IndexByte(mock, ':')

	if separatorIndex == -1 {
		return &mockNameDecl{
			typeName: mock,
		}
	} else {
		return &mockNameDecl{
			typeName:      mock[:separatorIndex],
			interfaceName: mock[separatorIndex+1:],
		}
	}
}
