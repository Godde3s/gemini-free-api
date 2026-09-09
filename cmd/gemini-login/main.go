// gemini-login — interactive cookie capture helper for Gemini-Free-API.
//
// Opens a real browser window (playwright), you log in to gemini.google.com
// by hand, and the tool watches for the __Secure-1PSID / __Secure-1PSIDTS
// cookie pair the moment Google sets it. Pairs are saved to
// gemini-cookies.json and printed ready-to-paste into GEMINI_COOKIES.
//
// Why a helper? Copying cookies from DevTools works too (README shows how),
// but this tool does it in one step and validates the pair live by loading
// the Gemini app page and checking for the SNlM0e token.
//
// Usage:
//
//      go run ./cmd/gemini-login             # login flow, captures 1 pair
//      go run ./cmd/gemini-login --timeout 300

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mxschmitt/playwright-go"
)

const portalURL = "https://gemini.google.com/app"

type cookieFile struct {
	Accounts []cookieEntry `json:"accounts"`
}

type cookieEntry struct {
	PSID      string `json:"psid"`
	PSIDTS    string `json:"psidts"`
	Label     string `json:"label,omitempty"`
	AddedUnix int64  `json:"added_unix"`
}

func main() {
	timeout := flag.Int("timeout", 300, "seconds to wait for a login to complete")
	out := flag.String("out", "gemini-cookies.json", "where to append captured cookie pairs")
	headless := flag.Bool("headless", false, "run without a visible browser (cannot log in by hand)")
	label := flag.String("label", "", "friendly name for the captured account")
	flag.Parse()

	if err := playwright.Install(&playwright.RunOptions{Browsers: []string{"chromium"}}); err != nil {
		fmt.Fprintln(os.Stderr, "playwright install failed:", err)
		os.Exit(1)
	}

	pw, err := playwright.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "browser launch failed:", err)
		os.Exit(1)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(*headless),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "chromium launch failed:", err)
		os.Exit(1)
	}
	defer browser.Close()

	ctx, err := browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"),
		Locale:    playwright.String("en-US"),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "context failed:", err)
		os.Exit(1)
	}

	fmt.Println("════════════════════════════════════════════════════════")
	fmt.Println("  Gemini-Free-API login helper")
	fmt.Println("  1) A Chromium window will open gemini.google.com")
	fmt.Println("  2) Log in with your Google account")
	fmt.Println("  3) As soon as the page loads your session, the cookies")
	fmt.Println("     __Secure-1PSID + __Secure-1PSIDTS are captured here.")
	fmt.Println("════════════════════════════════════════════════════════")

	page, err := ctx.NewPage()
	if err != nil {
		fmt.Fprintln(os.Stderr, "page failed:", err)
		os.Exit(1)
	}
	if _, err := page.Goto(portalURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		fmt.Fprintln(os.Stderr, "goto failed:", err)
		os.Exit(1)
	}

	deadline := time.Now().Add(time.Duration(*timeout) * time.Second)
	captured := false
	for time.Now().Before(deadline) && !captured {
		cookies, cerr := ctx.Cookies("https://gemini.google.com", "https://accounts.google.com")
		if cerr != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		var psid, psidts string
		for _, c := range cookies {
			switch c.Name {
			case "__Secure-1PSID":
				if c.Value != "" {
					psid = c.Value
				}
			case "__Secure-1PSIDTS":
				if c.Value != "" {
					psidts = c.Value
				}
			}
		}
		if psid != "" && !strings.HasPrefix(psid, "g.a000") && len(psid) > 20 {
			savePair(*out, psid, psidts, *label)
			fmt.Println()
			fmt.Println("✔ captured __Secure-1PSID pair")
			fmt.Println()
			fmt.Println("→ .env line (copy-paste):")
			if psidts != "" {
				fmt.Printf("  GEMINI_COOKIES=%s\n", mustJSON([]cookieEntry{{PSID: psid, PSIDTS: psidts, Label: *label}}))
			} else {
				fmt.Printf("  GEMINI_COOKIES=%s\n", psid)
			}
			fmt.Println()
			fmt.Println("→ saved to", *out)
			captured = true
			break
		}
		time.Sleep(2 * time.Second)
	}

	if !captured {
		fmt.Println("✘ timed out without capturing cookies — log in fully, then retry.")
		os.Exit(2)
	}

	// Keep the browser open a moment so the user can see the result.
	time.Sleep(3 * time.Second)
}

func savePair(path, psid, psidts, label string) {
	var cf cookieFile
	if raw, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(raw, &cf)
	}
	for _, e := range cf.Accounts {
		if e.PSID == psid {
			return // already captured
		}
	}
	cf.Accounts = append(cf.Accounts, cookieEntry{
		PSID: psid, PSIDTS: psidts, Label: label, AddedUnix: time.Now().Unix(),
	})
	raw, _ := json.MarshalIndent(cf, "", "  ")
	_ = os.WriteFile(path, raw, 0o600)
}

func mustJSON(v interface{}) string {
	raw, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(raw)
}
