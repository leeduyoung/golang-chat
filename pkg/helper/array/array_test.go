package array

import "testing"

func TestContains(t *testing.T) {
	t.Run("type is int", func(t *testing.T) {
		intArr := []int{1, 2, 3}

		if !Contains(intArr, 1) {
			t.Errorf("expected true, got false")
		}

		t.Log("[Success] TestContains - type is int")
	})

	t.Run("type is string", func(t *testing.T) {
		strArr := []string{"a", "b", "c"}

		if Contains(strArr, "d") {
			t.Errorf("expected false, got true")
		}

		t.Log("[Success] TestContains - type is string")
	})
}

func TestRemoveDuplcateItem(t *testing.T) {
	intArr := []int{1, 1, 1, 4, 5, 6, 7, 8, 9, 10}

	result := RemoveDuplcateItem(intArr)

	if len(result) != 8 {
		t.Errorf("expected 7, got %d", len(result))
	}

	t.Log("[Success] TestRemoveDuplcateItem")
}
