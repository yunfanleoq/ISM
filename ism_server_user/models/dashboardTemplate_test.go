package models

import "testing"

func TestNormalizeTemplateKindConvergesToThreeRoles(t *testing.T) {
	cases := map[string]string{
		"home":          "home",
		"deviceList":    "deviceList",
		"datapointList": "datapointList",
		"zone":          "deviceList",
		"room":          "deviceList",
		"cabinet":       "deviceList",
		"floor":         "deviceList",
		"device":        "datapointList",
		"unknown":       "",
	}
	for input, want := range cases {
		if got := normalizeTemplateKind(input); got != want {
			t.Fatalf("normalizeTemplateKind(%q)=%q, want %q", input, got, want)
		}
	}
	if got := normalizeTemplateModelUuid("datapointList", "arbitrary-model"); got != "" {
		t.Fatalf("物模型不得参与模板选择，got %q", got)
	}
}

func TestPaginateModelDataPointsIsDatasetIndependent(t *testing.T) {
	points := []map[string]string{
		{"name": "temperature", "uuid": "p-1", "unit": "C"},
		{"name": "pressure", "uuid": "p-2", "unit": "kPa"},
		{"name": "humidity", "uuid": "p-3", "unit": "%"},
	}
	page, total := paginateModelDataPoints(points, "model-any", "", 2, 2)
	if total != 3 || len(page) != 1 || page[0].Name != "humidity" {
		t.Fatalf("分页错误: total=%d page=%+v", total, page)
	}
	if page[0].Muid != "model-any" || page[0].ModelDataUuid != "p-3" {
		t.Fatalf("点位归属错误: %+v", page[0])
	}

	empty, total := paginateModelDataPoints(points, "model-any", "not-found", 1, 20)
	if total != 0 || len(empty) != 0 {
		t.Fatalf("空筛选应返回空页: total=%d page=%+v", total, empty)
	}

	beyond, total := paginateModelDataPoints(points, "model-any", "", 99, 2)
	if total != 3 || len(beyond) != 0 {
		t.Fatalf("越界页应保持总数并返回空页: total=%d page=%+v", total, beyond)
	}
}
