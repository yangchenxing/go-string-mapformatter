package mapformatter

import (
	"testing"
)

func TestMapFormatter(t *testing.T) {
	var text string
	var err error

	text, err = Format("test %(a|s) %%%(b|d)", map[string]interface{}{
		"a": "foo",
		"b": 1,
	})
	if err != nil {
		t.Error("format fail:", err.Error())
		return
	}
	if text != "test foo %1" {
		t.Error("unexpected format result:", text)
		return
	}

	text = MustFormat("Hello %(name|s), you owe me %(money|.2f) dollar.",
		map[string]interface{}{
			"name":  "anyone",
			"money": 10.3,
		})
	if text != "Hello anyone, you owe me 10.30 dollar." {
		t.Error("unexpected format result:", text)
		return
	}

	// test bad format
	if text = MustFormat("%"); text != "%!(NOVERB)" {
		t.Error("unexpected format result:", text)
		return
	}
	if text = MustFormat("%(TEST"); text != "%!((MISSING)TEST" {
		t.Error("unexpected format result:", text)
		return
	}
	if text = MustFormat("%(TEST|"); text != "%!((MISSING)TEST|" {
		t.Error("unexpected format result:", text)
		return
	}
}
