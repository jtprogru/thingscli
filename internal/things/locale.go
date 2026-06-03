package things

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Locale holds the localized names of the seven built-in Things 3 lists.
// We use these to talk to Things from any system locale.
type Locale struct {
	Lang     string `json:"lang"`
	Inbox    string `json:"inbox"`
	Today    string `json:"today"`
	Upcoming string `json:"upcoming"`
	Anytime  string `json:"anytime"`
	Someday  string `json:"someday"`
	Logbook  string `json:"logbook"`
	Trash    string `json:"trash"`
}

// builtinLocales is the table of known list-name translations.
// en and ru are verified against real Things 3 installs.
// The others are best-effort — if Things actually uses a different name
// for one of these languages, the runtime probe will simply not match
// and we will fall back to the next candidate (eventually English).
var builtinLocales = []Locale{
	{Lang: "en", Inbox: "Inbox", Today: "Today", Upcoming: "Upcoming", Anytime: "Anytime", Someday: "Someday", Logbook: "Logbook", Trash: "Trash"},
	{Lang: "ru", Inbox: "Входящие", Today: "Сегодня", Upcoming: "Завтра", Anytime: "В любое время", Someday: "Когда-нибудь", Logbook: "Журнал", Trash: "Корзина"},
	{Lang: "de", Inbox: "Eingang", Today: "Heute", Upcoming: "Demnächst", Anytime: "Jederzeit", Someday: "Irgendwann", Logbook: "Logbuch", Trash: "Papierkorb"},
	{Lang: "fr", Inbox: "Boîte de réception", Today: "Aujourd'hui", Upcoming: "À venir", Anytime: "N'importe quand", Someday: "Un jour", Logbook: "Journal", Trash: "Corbeille"},
	{Lang: "es", Inbox: "Entrada", Today: "Hoy", Upcoming: "Próximamente", Anytime: "En cualquier momento", Someday: "Algún día", Logbook: "Registro", Trash: "Papelera"},
	{Lang: "it", Inbox: "In entrata", Today: "Oggi", Upcoming: "In arrivo", Anytime: "Quando capita", Someday: "Prima o poi", Logbook: "Diario", Trash: "Cestino"},
	{Lang: "ja", Inbox: "受信箱", Today: "今日", Upcoming: "近日予定", Anytime: "いつでも", Someday: "いつか", Logbook: "ログブック", Trash: "ゴミ箱"},
	{Lang: "zh-Hans", Inbox: "收件箱", Today: "今天", Upcoming: "计划", Anytime: "随时", Someday: "某天", Logbook: "日志", Trash: "废纸篓"},
	{Lang: "pt-BR", Inbox: "Caixa de entrada", Today: "Hoje", Upcoming: "Em breve", Anytime: "A qualquer momento", Someday: "Algum dia", Logbook: "Diário", Trash: "Lixo"},
	{Lang: "nl", Inbox: "Postvak IN", Today: "Vandaag", Upcoming: "Binnenkort", Anytime: "Wanneer dan ook", Someday: "Ooit", Logbook: "Logboek", Trash: "Prullenmand"},
}

// Name returns the localized name for a built-in list key.
func (l Locale) Name(k ListKey) (string, error) {
	switch k {
	case ListInbox:
		return l.Inbox, nil
	case ListToday:
		return l.Today, nil
	case ListUpcoming:
		return l.Upcoming, nil
	case ListAnytime:
		return l.Anytime, nil
	case ListSomeday:
		return l.Someday, nil
	case ListLogbook:
		return l.Logbook, nil
	case ListTrash:
		return l.Trash, nil
	}
	return "", fmt.Errorf("unknown list key: %s", k)
}

const localeCacheTTL = 30 * 24 * time.Hour

type cachedLocale struct {
	Saved time.Time `json:"saved"`
	Names Locale `json:"names"`
}

func cachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache", "thingscli", "locale.json"), nil
}

func loadCached() *Locale {
	p, err := cachePath()
	if err != nil {
		return nil
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil
	}
	var c cachedLocale
	if err := json.Unmarshal(data, &c); err != nil {
		return nil
	}
	if time.Since(c.Saved) > localeCacheTTL {
		return nil
	}
	return &c.Names
}

func saveCached(n Locale) {
	p, err := cachePath()
	if err != nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(p), 0o755)
	data, err := json.Marshal(cachedLocale{Saved: time.Now(), Names: n})
	if err != nil {
		return
	}
	_ = os.WriteFile(p, data, 0o600)
}

// detectLocale figures out which language Things 3 is using right now.
// Strategy: ask Things, in one osascript call, which of the candidate
// "Inbox" names actually resolves to a real list. The one that resolves
// determines the language.
//
// Returns the matched Locale, or an error if none match (which would
// mean Things is using a locale we don't yet have a translation for —
// the user can supply names manually via --list).
func detectLocale() (Locale, error) {
	candidates := make([]string, 0, len(builtinLocales))
	for _, l := range builtinLocales {
		candidates = append(candidates, l.Inbox)
	}
	script := buildProbeScript(candidates)
	out, err := runOsa(script)
	if err != nil {
		return Locale{}, fmt.Errorf("probe Things locale: %w", err)
	}
	matched := strings.TrimSpace(out)
	if matched == "" {
		return Locale{}, errors.New("things 3 did not recognize any known Inbox name; pass --lang or open Things to verify it is installed")
	}
	for _, l := range builtinLocales {
		if l.Inbox == matched {
			return l, nil
		}
	}
	return Locale{}, fmt.Errorf("things returned unexpected list name: %q", matched)
}

// ResolveLocale returns the active Things locale, using cache when possible.
// Pass forceRefresh=true to ignore the cache.
// If langOverride is non-empty, the locale for that language is used (if known)
// without probing; useful when auto-detect picks the wrong one.
func ResolveLocale(langOverride string, forceRefresh bool) (Locale, error) {
	if langOverride != "" {
		for _, l := range builtinLocales {
			if strings.EqualFold(l.Lang, langOverride) {
				return l, nil
			}
		}
		return Locale{}, fmt.Errorf("unknown language %q; known: en, ru, de, fr, es, it, ja, zh-Hans, pt-BR, nl", langOverride)
	}
	if !forceRefresh {
		if c := loadCached(); c != nil {
			return *c, nil
		}
	}
	l, err := detectLocale()
	if err != nil {
		return Locale{}, err
	}
	saveCached(l)
	return l, nil
}
