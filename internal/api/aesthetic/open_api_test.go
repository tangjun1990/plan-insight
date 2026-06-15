package aesthetic

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func TestBuildOpenAccessToken(t *testing.T) {
	pubkey := "test-pubkey"
	prikey := "test-prikey"
	ts := time.Unix(1717660800, 0)

	token := BuildOpenAccessToken(pubkey, prikey, ts)
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		t.Fatalf("decode token failed: %v", err)
	}

	expected := pubkey + "." + buildOpenAccessSign(pubkey, prikey, ts) + ".1717660800"
	if string(decoded) != expected {
		t.Fatalf("unexpected token payload: got %q want %q", string(decoded), expected)
	}
}

func TestBuildOpenAuthorizationHeader(t *testing.T) {
	// 增加随机因子
	rand.Seed(time.Now().UnixNano())
	header := BuildOpenAuthorizationHeader("huazhuhui", "sah1IvQHw0b8E3STlJijOf1RBDMbSbAY", time.Now())
	fmt.Println(header)
}

func TestIsValidPhoneNumber(t *testing.T) {
	cases := []struct {
		name  string
		phone string
		valid bool
	}{
		{name: "valid", phone: "13800138000", valid: true},
		{name: "too short", phone: "1380013800", valid: false},
		{name: "contains letters", phone: "1380013a000", valid: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isValidPhoneNumber(tc.phone); got != tc.valid {
				t.Fatalf("isValidPhoneNumber(%q) = %v, want %v", tc.phone, got, tc.valid)
			}
		})
	}
}

func TestGenerateOpenAuthorizationHeaderForManualRequest(t *testing.T) {
	pubkey := strings.TrimSpace(os.Getenv("PLAN_OPEN_PUBKEY"))
	prikey := strings.TrimSpace(os.Getenv("PLAN_OPEN_PRIKEY"))
	if pubkey == "" || prikey == "" {
		t.Skip("set PLAN_OPEN_PUBKEY and PLAN_OPEN_PRIKEY to print a real Authorization header")
	}

	header := BuildOpenAuthorizationHeader(pubkey, prikey, time.Now())
	t.Logf("Authorization: %s", header)
}
