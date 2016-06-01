package mapformatter

import (
	"testing"
)

func TestMapFormatter(t *testing.T) {
	var text string

	text = Format("test %(a|s) %%%(b|d)", map[string]interface{}{
		"a": "foo",
		"b": 1,
	})
	if text != "test foo %1" {
		t.Error("unexpected format result:", text)
		return
	}

	text = Format("Hello %(name|s), you owe me %(money|.2f) dollar.",
		map[string]interface{}{
			"name":  "anyone",
			"money": 10.3,
		})
	if text != "Hello anyone, you owe me 10.30 dollar." {
		t.Error("unexpected format result:", text)
		return
	}

	// test bad format
	if mf, _ := newMapFormatter("%"); mf.fmt != "%" {
		t.Error("unexpected format:", mf.fmt)
		return
	}
	if mf, _ := newMapFormatter("%("); mf.fmt != "%(" {
		t.Error("unexpected format:", mf.fmt)
		return
	}
	if mf, _ := newMapFormatter("%(foo"); mf.fmt != "%(foo" {
		t.Error("unexpected format:", mf.fmt)
		return
	}
	if mf, _ := newMapFormatter("%(foo|"); mf.fmt != "%(foo|" {
		t.Error("unexpected format:", mf.fmt)
		return
	}
	if mf, _ := newMapFormatter("%(foo|bar"); mf.fmt != "%(foo|bar" {
		t.Error("unexpected format:", mf.fmt)
		return
	}
	if mf, _ := newMapFormatter("%d"); mf.fmt != "%d" {
		t.Error("unexpected format:", mf.fmt)
		return
	}

	if text = (*mapFormatter)(nil).format(map[string]interface{}{}); text != "!(INVALID_FORMATTER)" {
		t.Error("unexpected format result:", text)
		return
	}
}
