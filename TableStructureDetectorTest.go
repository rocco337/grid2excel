package main

import "testing"

//TestTableStructureDetector ...
func TestTableStructureDetector(t *testing.T) {

	service := TableStructureDetector{}

	t.Run("Should return table structure", func(t *testing.T) {
		service.Detect("test1.jpg")
	})
}
